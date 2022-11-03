package cache

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"team-client-server/config"
	"team-client-server/sfnake"
	"team-client-server/tools"
	"time"
)

var (
	fileStoreCache = cache.New(10*time.Minute, 5*time.Minute)
)

type FileInfo struct {
	// Path 文件路径
	Path string `json:"path,omitempty"`
	// Expire 有效期
	Expire time.Duration `json:"expire,omitempty"`
	// AppId 应用ID
	AppId string `json:"appId,omitempty"`
}

type FileUploadType string

const (
	FileUploadTypeHttp FileUploadType = "http"
)

var (
	fileCacheLocalRecord string

	fileCacheTempDir  string
	filecacheStoreDir string
)

// FileUploadInfo 文件上传信息
type FileUploadInfo struct {
	// URL 上传到的路径
	URL string `json:"url,omitempty"`
	// Method 上传方式, 默认POST
	Method string `json:"method,omitempty"`
	// Type 上传类型目前仅仅支持Http
	Type FileUploadType `json:"type,omitempty"`
	// Data 上传时携带的数据
	Data map[string]string `json:"data,omitempty"`
	// Header 请求头
	Header map[string]string `json:"header,omitempty"`
	// FileList 文件名称: 文件ID
	FileList map[string]string `json:"fileIdList,omitempty"`
	// AppId 应用ID
	AppId string `json:"appId,omitempty"`
}

func InitHttpDownloadFileHand(engine *gin.Engine) {
	engine.GET("/cache/file/download/:appId/:cacheId", fileCacheDownload)
}

func initFileCache(engine *gin.RouterGroup) {
	_ = os.MkdirAll(config.CurrentConfig.Cache.File, 0755)
	fileCacheLocalRecord = filepath.Join(config.CurrentConfig.Cache.File, ".record")
	if stat, err := os.Stat(fileCacheLocalRecord); err == nil && !stat.IsDir() {
		_ = fileStoreCache.LoadFile(fileCacheLocalRecord)
	}

	filecacheStoreDir = filepath.Join(config.CurrentConfig.Cache.File, "store")
	_ = os.MkdirAll(filecacheStoreDir, 0755)

	fileCacheTempDir = filepath.Join(config.CurrentConfig.Cache.File, "temp")
	tempDir, err := os.ReadDir(fileCacheTempDir)
	if err == nil {
		for _, d := range tempDir {
			fileStoreCache.Delete(d.Name())
			_ = os.Remove(filepath.Join(fileCacheTempDir, d.Name()))
		}
	}

	_ = os.Remove(fileCacheTempDir)
	_ = os.MkdirAll(fileCacheTempDir, 0755)

	fileStoreCache.OnEvicted(func(s string, i interface{}) {
		p, ok := i.(string)
		if !ok {
			return
		}
		_ = os.Remove(p)
		_ = fileStoreCache.SaveFile(fileCacheLocalRecord)
	})

	engine.Group("/file").
		POST("store", ginmiddleware.WrapperResponseHandle(fileCacheStore)).
		POST("del/:appId/:cacheId", ginmiddleware.WrapperResponseHandle(fileCacheDel)).
		POST("clear/:appId", ginmiddleware.WrapperResponseHandle(fileCacheClear)).
		POST("delay/:cacheId", ginmiddleware.WrapperResponseHandle(fileCacheDelay)).
		POST("upload/:appId/:cacheId", fileCacheUpload)
}

func getCacheId(appId string, cacheId string) (string, error) {
	var err error
	if cacheId == "" {
		cacheId, err = sfnake.GetIdStr()
		if err != nil {
			return "", fmt.Errorf("生成缓存ID失败: %s", err.Error())
		}
	}

	return appId + "-" + cacheId, nil
}

var (
	// fileCacheDownload 文件缓存下载
	fileCacheDownload gin.HandlerFunc = func(ctx *gin.Context) {
		appId := ctx.Param("appId")
		cacheId := ctx.Param("cacheId")
		if cacheId == "" {
			ctx.Status(http.StatusNotFound)
			return
		}

		id, _ := getCacheId(appId, cacheId)
		fPath, ok := fileStoreCache.Get(id)
		if !ok {
			ctx.Status(http.StatusNotFound)
			return
		}

		ctx.File(fPath.(string))
	}
	// fileCacheStore 文件缓存存储
	fileCacheStore ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		var storeInfo *FileInfo
		if err := ctx.ShouldBindJSON(&storeInfo); err != nil {
			return fmt.Errorf("解析缓存文件信息失败: %s", err.Error())
		}

		if storeInfo.AppId == "" {
			return errors.New("缺失应用标识")
		}

		p := storeInfo.Path
		pStat, err := os.Stat(p)
		if err != nil {
			return fmt.Errorf("[%s]路径不存在", p)
		}

		if pStat.IsDir() {
			return errors.New("不支持缓存目录")
		}

		if storeInfo.Expire < -1 {
			storeInfo.Expire = -1
		}

		if storeInfo.Expire > 0 {
			storeInfo.Expire = time.Millisecond * storeInfo.Expire
		}

		idStr, err := getCacheId(storeInfo.AppId, "")
		if err != nil {
			return err
		}

		fPath := filepath.Join(fileCacheTempDir, idStr)
		if storeInfo.Expire == -1 {
			fPath = filepath.Join(filecacheStoreDir, idStr)
		}

		f, err := os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY, 0655)
		if err != nil {
			return fmt.Errorf("创建缓存文件失败: %s", err.Error())
		}
		defer f.Close()

		srcFile, err := os.OpenFile(p, os.O_RDONLY, 0655)
		if err != nil {
			return fmt.Errorf("打开[%s]失败: %s", p, err.Error())
		}

		if _, err = io.Copy(f, srcFile); err != nil {
			return fmt.Errorf("拷贝[%s]至缓存中失败: %s", p, err.Error())
		}

		fileStoreCache.Set(idStr, fPath, storeInfo.Expire)
		if storeInfo.Expire == -1 {
			if err := fileStoreCache.SaveFile(fileCacheLocalRecord); err != nil {
				fileStoreCache.Delete(idStr)
				return fmt.Errorf("永久性存储缓存文件失败: %s", err.Error())
			}
		}
		return strings.Split(idStr, "-")[1]
	}
	// fileCacheDel 文件缓存删除
	fileCacheDel ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		appId := ctx.Param("appId")
		cacheId := ctx.Param("cacheId")
		if cacheId == "" {
			return errors.New("缓存ID不能为空")
		}

		id, err := getCacheId(appId, cacheId)
		if err != nil {
			return err
		}

		fileStoreCache.Delete(id)
		return nil
	}
	// fileCacheClear 文件缓存清除
	fileCacheClear ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		appId := ctx.Param("appId")
		if appId == "" {
			return fmt.Errorf("appId不能为空")
		}

		cachePrefix := appId + "-"
		items := fileStoreCache.Items()
		for key := range items {
			if strings.HasPrefix(key, cachePrefix) {
				fileStoreCache.Delete(key)
			}
		}

		return nil
	}
	// fileCacheDelay 文件缓存延期
	fileCacheDelay ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		var storeInfo *FileInfo
		if err := ctx.ShouldBindJSON(&storeInfo); err != nil {
			return fmt.Errorf("解析缓存文件信息失败: %s", err.Error())
		}

		if storeInfo.AppId == "" {
			return errors.New("缺失应用标识")
		}

		cacheId := ctx.Param("cacheId")
		if cacheId == "" {
			return errors.New("缺失缓存ID")
		}

		id, err := getCacheId(storeInfo.AppId, cacheId)
		if err != nil {
			return err
		}

		v, ok := fileStoreCache.Get(id)
		if !ok {
			return false
		}

		if storeInfo.Expire < -1 {
			storeInfo.Expire = -1
		}

		storeInfo.Expire = time.Millisecond * storeInfo.Expire

		fileStoreCache.Set(id, v, storeInfo.Expire)
		return nil
	}
	// fileCacheUpload 文件缓存上传
	fileCacheUpload gin.HandlerFunc = func(ctx *gin.Context) {
		var param *FileUploadInfo
		if err := ctx.ShouldBindJSON(&param); err != nil {
			ctx.String(http.StatusBadGateway, "请求参数解析失败: %s", err.Error())
			return
		}

		if param.AppId == "" {
			ctx.String(http.StatusBadGateway, "应用ID不能为空")
			return
		}

		if param.FileList == nil || len(param.FileList) == 0 {
			ctx.String(http.StatusBadGateway, "文件列表不能为空")
			return
		}

		if param.Type != FileUploadTypeHttp {
			ctx.String(http.StatusBadGateway, "上传方式目前只支持http")
			return
		}

		if param.Method == "" {
			param.Method = "POST"
		}

		param.Method = strings.ToUpper(param.Method)

		if param.URL == "" {
			ctx.String(http.StatusBadGateway, "转发的URL不能为空")
			return
		}

		pipeR, pipeW := io.Pipe()
		m := multipart.NewWriter(pipeW)
		go func() {
			defer pipeW.Close()
			defer m.Close()

			for key, val := range param.FileList {
				if err := fileCacheWriteTo(m, key, param.AppId, val); err != nil {
					ctx.String(http.StatusBadGateway, err.Error())
					return
				}
			}

			for k, v := range param.Data {
				if err := m.WriteField(k, v); err != nil {
					ctx.String(http.StatusBadGateway, "写出表单数据失败: %s", err.Error())
					return
				}
			}
		}()

		req, err := http.NewRequest(param.Method, param.URL, pipeR)
		if err != nil {
			ctx.String(http.StatusBadGateway, "创建请求对象失败: %s", err.Error())
			return
		}

		for k, v := range param.Header {
			req.Header.Set(k, v)
		}

		req.Header.Add("Content-Type", m.FormDataContentType())

		res, err := tools.SkipVerifyCertHttpClient.Do(req)
		if err != nil {
			ctx.String(http.StatusBadGateway, "获取请求响应失败: %s", err.Error())
			return
		}
		defer res.Body.Close()

		if err = res.Write(ctx.Writer); err != nil {
			ctx.String(http.StatusBadGateway, "转发响应数据失败: %s", err.Error())
			return
		}

	}
)

// fileCacheWriteTo 文件缓存写出
func fileCacheWriteTo(m *multipart.Writer, fieldName string, appId string, cacheId string) error {

	id, _ := getCacheId(appId, cacheId)
	p, ok := fileStoreCache.Get(id)
	if !ok {
		return errors.New("获取缓存文件失败")
	}
	f, err := os.OpenFile(p.(string), os.O_RDONLY, 0655)
	if err != nil {
		return fmt.Errorf("打开缓存文件失败: %s", err.Error())
	}
	defer f.Close()

	dest, err := m.CreateFormField(fieldName)
	if err != nil {
		return fmt.Errorf("生成文件上传字段失败: %s", err.Error())
	}
	if _, err = io.Copy(dest, f); err != nil {
		return fmt.Errorf("拷贝缓存至http数据流失败: %s", err.Error())
	}

	return nil
}
