package cache

import (
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"team-client-server/remoteserver"
)

func initRemoteServerCache(routes *gin.RouterGroup) {
	routes.Group("/remote").
		POST("/user/list", ginmiddleware.WrapperResponseHandle(remoteUserList))
}

var (
	// remoteUserList 获取远程用户列表
	remoteUserList ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		res := remoteserver.GetCacheByType(remoteserver.CacheTypeUserList)
		if res == "" {
			res = "[]"
		}
		return res
	}
)
