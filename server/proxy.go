package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-base-lib/logs"
	"net/http/httputil"
	"net/url"
	"os"
	"team-client-server/remoteserver"
)

var localWebServerHttpProxy *httputil.ReverseProxy

func init() {
	uri, err := url.Parse(remoteserver.LocalWebServerAddress)
	if err != nil {
		logs.Fatalf("解析web服务器地址格式失败: %s", err.Error())
		os.Exit(10)
	}
	localWebServerHttpProxy = httputil.NewSingleHostReverseProxy(uri)
}

func initProxy(engine *gin.Engine) {
	engine.Any("p/web/*path", proxyLocalWebHandle)
}

func proxyLocalWebHandle(ctx *gin.Context) {
	token := remoteserver.Token()
	header := ctx.Request.Header
	if token != "" {
		header.Set("t", token)
	}
	header.Set("User-Agent", "teamManageLocal")
	ctx.Request.URL.Path = ctx.Request.URL.Path[6:]
	localWebServerHttpProxy.ServeHTTP(ctx.Writer, ctx.Request)
}

func proxyAppWebHandle(ctx *gin.Context) {

}
