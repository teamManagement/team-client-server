package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"team-client-server/tools"
)

type proxy443HttpHandlerImpl struct {
}

func (p *proxy443HttpHandlerImpl) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	serverName := request.TLS.ServerName
	if serverName == "" {
		writer.WriteHeader(403)
		return
	}

	if _, ok := proxy443AllowServerName[serverName]; !ok {
		writer.WriteHeader(401)
		return
	}

	if !strings.HasPrefix(request.URL.Path, "/") {
		request.URL.Path += "/"
	}

	request.URL.Path = fmt.Sprintf("/c/forward/%s%s", serverName, request.URL.Path)
	proxy443ToLockHttpProxy.ServeHTTP(writer, request)
}

var (
	proxy443Keypair        tls.Certificate
	proxy443KeypairLoadErr error

	proxy443IsRunning bool
	proxy443Lock      = sync.Mutex{}

	proxy443HttpServer      *http.Server
	proxy443HttpHandler     = &proxy443HttpHandlerImpl{}
	proxy443ToLockHttpProxy *httputil.ReverseProxy

	proxy443AllowServerName = make(map[string]struct{})
)

func init() {
	proxy443Keypair, proxy443KeypairLoadErr = tls.X509KeyPair(tools.Proxy443CertBytes, tools.Proxy443KeyBytes)
	uri, _ := url.Parse("https://127.0.0.1:65528")
	proxy443ToLockHttpProxy = httputil.NewSingleHostReverseProxy(uri)

	block, _ := pem.Decode(tools.Proxy443CertBytes)
	if block == nil {
		proxy443KeypairLoadErr = errors.New("转换443代理客户端证书编码格式失败")
		return
	}

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		proxy443KeypairLoadErr = fmt.Errorf("解析443代理客户端证书内容失败")
		return
	}

	for _, name := range certificate.DNSNames {
		if name == "127.0.0.1" || name == "apps.byzk.cn" {
			continue
		}

		proxy443AllowServerName[name] = struct{}{}
	}

}

func initProxy443Route(engine *gin.Engine) {
	engine.Group("/proxy/config").
		POST("is/running", ginmiddleware.WrapperResponseHandle(proxy443IsRunningCheck)).
		POST("start", ginmiddleware.WrapperResponseHandle(proxy443Start)).
		POST("shutdown", ginmiddleware.WrapperResponseHandle(proxy443Shutdown))
}

var (
	// proxy443IsRunningCheck 检查443端口是否已经启动
	proxy443IsRunningCheck ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		proxy443Lock.Lock()
		defer proxy443Lock.Unlock()
		return proxy443IsRunning
	}

	// proxy443Start 启动443端口代理
	proxy443Start ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		proxy443Lock.Lock()
		defer proxy443Lock.Unlock()

		if proxy443IsRunning {
			return nil
		}

		if proxy443KeypairLoadErr != nil {
			return fmt.Errorf("加载TLS密钥对失败: %s", proxy443KeypairLoadErr.Error())
		}

		if tools.TelnetHost("127.0.0.1:443") {
			return errors.New("本机443端口被占用, 无法开启请求代理, 请检查本机端口占用情况")
		}

		tlsListener, err := tls.Listen("tcp", "127.0.0.1:443", &tls.Config{
			Certificates: []tls.Certificate{proxy443Keypair},
		})
		if err != nil {
			return fmt.Errorf("监听443端口失败: %s", err.Error())
		}

		proxy443HttpServer = &http.Server{
			Handler: proxy443HttpHandler,
		}
		proxy443IsRunning = true
		go func() {
			_ = proxy443HttpServer.Serve(tlsListener)
			lock.Lock()
			defer lock.Unlock()
			proxy443IsRunning = false
			_ = proxy443HttpServer.Close()
			proxy443HttpServer = nil
		}()
		return nil
	}

	// proxy443Shutdown 443端口停止代理
	proxy443Shutdown ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		lock.Lock()
		defer lock.Unlock()

		if !proxy443IsRunning {
			return nil
		}

		proxy443IsRunning = false
		if proxy443HttpServer == nil {
			return nil
		}

		_ = proxy443HttpServer.Close()
		return nil
	}
)
