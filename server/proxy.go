package server

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/gin-gonic/gin"
	"github.com/go-base-lib/coderutils"
	"github.com/go-base-lib/logs"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"team-client-server/config"
	"team-client-server/remoteserver"
	"team-client-server/vos"
)

var localWebServerHttpProxy *httputil.ReverseProxy

var lock = &sync.Mutex{}

func init() {
	uri, err := url.Parse(remoteserver.LocalWebServerAddress)
	if err != nil {
		logs.Fatalf("解析web服务器地址格式失败: %s", err.Error())
		os.Exit(10)
	}
	localWebServerHttpProxy = httputil.NewSingleHostReverseProxy(uri)
	localWebServerHttpProxy.ModifyResponse = httpResponseModify
}

func httpResponseModify(response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil
	}
	cacheH := response.Header.Get("_proxy_cache_result")
	if cacheH == "" {
		return nil
	}

	lock.Lock()
	defer lock.Unlock()

	f, cacheFilePath, err := config.CreateFileInConfigPath(filepath.Join("proxy", "_cache"), cacheH)
	if err != nil {
		return nil
	}

	if _, err = io.Copy(f, response.Body); err != nil {
		return err
	}

	if _, err = f.Seek(0, 0); err != nil {
		return fmt.Errorf("移动文件指针失败: %s", err.Error())
	}

	contentFileHash, err := coderutils.HashByReader(sha256.New(), f)
	if err != nil {
		return errors.New("获取文件响应体内容HASH失败")
	}

	if _, err = f.Seek(0, 0); err != nil {
		return fmt.Errorf("移动最终响应文件指针失败: %s", err.Error())
	}

	response.Body = f

	proxyHttpResponseCache := &vos.ProxyHttpResponseCache{
		RequestHash: cacheH,
	}

	proxyHttpResponseCacheModel := sqlite.Db().Model(&vos.ProxyHttpResponseCache{})
	if err = proxyHttpResponseCacheModel.Where(&proxyHttpResponseCache).First(&proxyHttpResponseCacheModel).Error; err == nil && contentFileHash.Equal(proxyHttpResponseCache.ContentHash) {
		return nil
	}

	//if _, err = f.Seek(0, 0); err != nil {
	//	return fmt.Errorf("还原响应体失败: %s", err.Error())
	//}

	marshal, _ := json.Marshal(response.Header)

	proxyHttpResponseCache.ResponseHeader = string(marshal)
	proxyHttpResponseCache.ResponseStatusCode = response.StatusCode
	proxyHttpResponseCache.ContentHash = contentFileHash
	proxyHttpResponseCache.ContentPath = cacheFilePath

	proxyHttpResponseCacheModel.Save(&proxyHttpResponseCache)

	return nil
}

func initProxy(engine *gin.Engine) {
	engine.Group("/p").
		POST("/register/http/proxy", ginmiddleware.WrapperResponseHandle(proxyRegisterHttpName)).
		Any("/web/*path", proxyLocalWebHandle)
}

var (
	proxyLocalWebHandle gin.HandlerFunc = func(ctx *gin.Context) {
		token := remoteserver.Token()
		header := ctx.Request.Header
		if token != "" {
			header.Set("_t", token)
		}
		header.Set("User-Agent", "teamManageLocal")
		ctx.Request.URL.Path = ctx.Request.URL.Path[6:]
		localWebServerHttpProxy.ServeHTTP(ctx.Writer, ctx.Request)
	}

	proxyRegisterHttpName ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		var proxyHttpServerInfo *vos.ProxyHttpServerInfo
		if err := ctx.ShouldBindJSON(&proxyHttpServerInfo); err != nil {
			return fmt.Errorf("解析请求参数失败: %s", err.Error())
		}

		if proxyHttpServerInfo.Name == "" {
			return fmt.Errorf("要注册的代理名称不能为空")
		}

		proxyHttpServerInfo.Schema = strings.ToLower(proxyHttpServerInfo.Schema)
		if proxyHttpServerInfo.Schema == "" {
			proxyHttpServerInfo.Schema = "http"
		}

		if proxyHttpServerInfo.Schema != "http" && proxyHttpServerInfo.Schema != "https" {
			return errors.New("代理的协议只支持http及https")
		}

		proxyHttpServerInfo.AllowReturnType = strings.ToLower(proxyHttpServerInfo.AllowReturnType)
		proxyHttpServerInfo.NotAllowReturnType = strings.ToLower(proxyHttpServerInfo.NotAllowReturnType)

		if err := sqlite.Db().Model(&proxyHttpServerInfo).Save(proxyHttpServerInfo).Error; err != nil {
			return fmt.Errorf("代理信息保存失败: %s", err.Error())
		}

		return nil
	}

	// proxyCacheForward 目标地址通过后台服务进行转发, 先查询本地缓存，
	//缓存未命中时，将请求转发到后台中再次进行转发, 此项缓存只要url+method存在极为命中
	//不会判断header用时请注意, 没有命中的只要服务器端返回状态码为200~299将会立即缓存数据, 其他则丢弃缓存
	proxyCacheForward gin.HandlerFunc = func(ctx *gin.Context) {
		var (
			token           string
			originUserAgent string

			proxyName       = ctx.Param("name")
			proxyTargetPath = ctx.Param("path")

			proxyServerInfo = &vos.ProxyHttpServerInfo{
				Name: proxyName,
			}
		)

		if err := sqlite.Db().Model(proxyServerInfo).First(proxyServerInfo).Error; err != nil {
			ctx.Status(405)
			return
		}

		if strings.HasPrefix(proxyTargetPath, "/") {
			proxyTargetPath = proxyTargetPath[1:]
		}

		requestHash := hex.EncodeToString(sha512.New().Sum([]byte(ctx.Request.RequestURI + ctx.Request.Method)))
		proxyHttpResponseCacheModel := sqlite.Db().Model(&vos.ProxyHttpResponseCache{})

		proxyResponseCache := &vos.ProxyHttpResponseCache{
			RequestHash: requestHash,
		}

		if err := proxyHttpResponseCacheModel.Where(&proxyResponseCache).First(&proxyResponseCache).Error; err == nil && proxyResponseCache.ContentHash != nil && proxyResponseCache.ContentPath != "" {
			h, err := coderutils.HashByFilePath(sha256.New(), proxyResponseCache.ContentPath)
			if err != nil || !h.Equal(proxyResponseCache.ContentHash) {
				goto StartForward
			}

			f, err := os.OpenFile(proxyResponseCache.ContentPath, os.O_RDONLY, 0655)
			if err != nil {
				goto StartForward
			}
			defer f.Close()

			if proxyResponseCache.ResponseHeader != "" {
				var header http.Header
				if err = json.Unmarshal([]byte(proxyResponseCache.ResponseHeader), &header); err != nil {
					goto StartForward
				}
				for k := range header {
					val := header[k]
					v := ""
					if len(val) > 0 {
						v = val[1]
					}
					ctx.Writer.Header().Set(k, v)
				}

				ctx.Writer.WriteHeader(proxyResponseCache.ResponseStatusCode)
			}

			_, _ = io.Copy(ctx.Writer, f)
			return
		}

	StartForward:
		token = remoteserver.Token()
		if token != "" {
			ctx.Request.Header.Set("_t", token)
		}

		ctx.Request.URL.Path = fmt.Sprintf("/proxy/c/forward/%s/%s/%s/%s", proxyServerInfo.Schema, ctx.Param("name"), proxyServerInfo.Host)
		originUserAgent = ctx.Request.Header.Get("User-Agent")
		ctx.Request.Header.Set("User-Agent", "teamManageLocal")
		ctx.Request.Header.Set("Origin-User-Agent", originUserAgent)
		ctx.Request.Header.Set("_proxy_cache_h", requestHash)
		localWebServerHttpProxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
)

func proxyAppWebHandle(ctx *gin.Context) {

}
