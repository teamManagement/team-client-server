package server

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-base-lib/logs"
	"github.com/sirupsen/logrus"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"team-client-server/db"
	"team-client-server/remoteserver"
	"team-client-server/sdk"
	cache2 "team-client-server/sdk/cache"
	"team-client-server/services"
	"team-client-server/tools"
	"team-client-server/updater"
	"team-client-server/website"
	"time"
)

func Run() {

	db.InitDb()

	engine := gin.New()

	initProxyLocal443Config()
	ginmiddleware.UseNotFoundHandle(engine)
	InitIcons(engine)
	website.InitAppWebSite(engine)
	cache2.InitHttpDownloadFileHand(engine)
	// cache
	engine.Any("/c/forward/:name/*path", proxyCacheForward)

	engine.Use(ginmiddleware.UseLogs(), ginmiddleware.UseVerifyUserAgent("teamManagerLocalView"))
	initProxy(engine)
	initProxy443Route(engine)

	engine.Use(ginmiddleware.UseRecover2HttpResult())
	initWs(engine)
	initLocalService(engine)
	sdk.InitLocalWebSdk(engine)
	updater.InitUpdaterHttpRestful(engine)
	remoteserver.InitLocalService(engine)
	services.InitLocalWebServices(engine)

	keyPair, err := tls.X509KeyPair(tools.ClientCertBytes, tools.ClientKeyBytes)
	if err != nil {
		logs.Panicf("SSL密钥解析失败")
	}

	listener, err := tls.Listen("tcp", "127.0.0.1:65528", &tls.Config{
		Certificates: []tls.Certificate{keyPair},
	})
	if err = engine.RunListener(listener); err != nil {
		logs.Panicf("服务发生异常, 程序将要终止, 错误信息: %s", err.Error())
	}
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := fmt.Sprintf("%6v", endTime.Sub(startTime))

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		//日志格式
		logs.WithFields(logrus.Fields{
			"http_status": statusCode,
			"total_time":  latencyTime,
			"ip":          clientIP,
			"method":      reqMethod,
			"uri":         reqUri,
		}).Info("access")
	}
}
