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
	"github.com/go-base-lib/goextension"
	"github.com/go-base-lib/logs"
	lockKey "github.com/sjy3/go-keylock"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"team-client-server/config"
	"team-client-server/remoteserver"
	"team-client-server/tools"
	"team-client-server/vos"
)

var keyLock = lockKey.NewKeyLock()

var localWebServerHttpProxy *httputil.ReverseProxy

func init() {
	uri, err := url.Parse(remoteserver.LocalWebServerAddress)
	if err != nil {
		logs.Fatalf("解析web服务器地址格式失败: %s", err.Error())
		os.Exit(10)
	}
	localWebServerHttpProxy = httputil.NewSingleHostReverseProxy(uri)
	localWebServerHttpProxy.Transport = tools.TlsTransport
	localWebServerHttpProxy.ModifyResponse = httpResponseModify
}

func cacheHttpSuccessResponse(response *http.Response) error {
	cacheH := response.Request.Header.Get("_proxy_cache_h")
	if cacheH == "" {
		return nil
	}

	keyLock.Lock(cacheH)

	responseBody := response.Body
	pipeReader, pipeWriter := io.Pipe()
	response.Body = pipeReader

	responseHeaderMarshal, _ := json.Marshal(response.Header)
	responseStatusCode := response.StatusCode
	go func() {
		defer keyLock.Unlock(cacheH)
		defer pipeWriter.Close()

		var (
			cacheFile       *os.File
			cacheFilePath   string
			contentFileHash goextension.Bytes

			err error
		)

		defer func() {
			if err != nil && cacheFilePath != "" {
				_ = os.Remove(cacheFilePath)
			}
		}()

		cacheFile, cacheFilePath, err = config.CreateFileInConfigPath(filepath.Join("proxy", "_cache"), cacheH)
		if err != nil {
			_, _ = io.Copy(pipeWriter, responseBody)
			return
		}
		defer cacheFile.Close()

		writer := io.MultiWriter(cacheFile, pipeWriter)

		_, err = io.Copy(writer, responseBody)

		if _, err = cacheFile.Seek(0, 0); err != nil {
			return
		}

		if contentFileHash, err = coderutils.HashByReader(sha256.New(), cacheFile); err != nil {
			return
		}

		proxyHttpResponseCacheModel := sqlite.Db().Model(&vos.ProxyHttpResponseCache{})

		var httpResponseCacheInfo *vos.ProxyHttpResponseCache
		proxyHttpResponseCacheModel.Select("content_path").Where("request_hash=?", cacheH).First(&httpResponseCacheInfo)
		if httpResponseCacheInfo != nil {
			var stat os.FileInfo
			if stat, err = os.Stat(httpResponseCacheInfo.ContentPath); err == nil && !stat.IsDir() {
				_ = os.Remove(httpResponseCacheInfo.ContentPath)
			}
		} else {
			httpResponseCacheInfo = &vos.ProxyHttpResponseCache{
				RequestHash: cacheH,
			}
		}

		httpResponseCacheInfo.ResponseHeader = string(responseHeaderMarshal)
		httpResponseCacheInfo.ContentPath = cacheFilePath
		httpResponseCacheInfo.ResponseStatusCode = responseStatusCode
		httpResponseCacheInfo.ContentHash = contentFileHash
		proxyHttpResponseCacheModel.Save(httpResponseCacheInfo)
	}()
	return nil
}

func httpResponseModify(response *http.Response) error {
	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return cacheHttpSuccessResponse(response)
	}

	cacheH := response.Request.Header.Get("_proxy_cache_h")
	if cacheH == "" {
		return nil
	}

	keyLock.Lock(cacheH)
	defer keyLock.Unlock(cacheH)

	proxyHttpResponseCache := &vos.ProxyHttpResponseCache{
		RequestHash: cacheH,
	}

	proxyHttpResponseCacheModel := sqlite.Db().Model(&vos.ProxyHttpResponseCache{})
	if err := proxyHttpResponseCacheModel.Where(&proxyHttpResponseCache).First(&proxyHttpResponseCache).Error; err != nil {
		return nil
	}

	contentPath := proxyHttpResponseCache.ContentPath
	f, err := os.OpenFile(contentPath, os.O_RDONLY, 0655)
	if err != nil {
		return nil
	}

	if h, err := coderutils.HashByReader(sha256.New(), f); err != nil || !h.Equal(proxyHttpResponseCache.ContentHash) {
		_ = f.Close()
		return nil
	}

	if _, err = f.Seek(0, 0); err != nil {
		_ = f.Close()
		return nil
	}

	if err = json.Unmarshal([]byte(proxyHttpResponseCache.ResponseHeader), &response.Header); err != nil {
		_ = f.Close()
		return nil
	}

	response.StatusCode = proxyHttpResponseCache.ResponseStatusCode
	_ = response.Body.Close()

	response.Body = f

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

	// proxyCacheForward 目标名称通过后台服务进行转发, 先查询本地缓存，
	//缓存未命中时，将请求转发到后台中再次进行转发, 此项缓存只要url+method存在极为命中
	//不会判断header用时请注意, 没有命中的只要服务器端返回状态码为200~299将会立即缓存数据, 其他则丢弃缓存
	proxyCacheForward gin.HandlerFunc = func(ctx *gin.Context) {
		var (
			token           string
			originUserAgent string

			proxyName       = ctx.Param("name")
			proxyTargetPath = ctx.Param("path")

			//proxyServerInfo = &vos.ProxyHttpServerInfo{
			//	Name: proxyName,
			//}
		)

		//if err := sqlite.Db().Model(proxyServerInfo).First(proxyServerInfo).Error; err != nil {
		//	ctx.Status(405)
		//	return
		//}

		if strings.HasPrefix(proxyTargetPath, "/") {
			proxyTargetPath = proxyTargetPath[1:]
		}

		requestHashBytes := sha512.Sum512([]byte(ctx.Request.RequestURI + ctx.Request.Method))
		requestHash := hex.EncodeToString(requestHashBytes[:])

		token = remoteserver.Token()
		if token != "" {
			ctx.Request.Header.Set("_t", token)
		}

		ctx.Request.URL.Path = fmt.Sprintf("/proxy/%s/%s", proxyName, proxyTargetPath)
		originUserAgent = ctx.Request.Header.Get("User-Agent")
		ctx.Request.Header.Set("User-Agent", "teamManageLocal")
		ctx.Request.Header.Set("Origin-User-Agent", originUserAgent)
		ctx.Request.Header.Set("_proxy_cache_h", requestHash)
		localWebServerHttpProxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
)

func proxyAppWebHandle(ctx *gin.Context) {

}
