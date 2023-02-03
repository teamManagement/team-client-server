package cache

import (
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"team-client-server/remoteserver"
)

func initRemoteServerCache(routes *gin.RouterGroup) {
	routes.Group("/remote").
		POST("/user/list", ginmiddleware.WrapperResponseHandle(remoteUserList)).
		POST("/org/list", ginmiddleware.WrapperResponseHandle(remoteOrgList))
}

var (
	// remoteUserList 获取远程用户列表
	remoteUserList ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		return remoteserver.GetCacheByType(remoteserver.CacheTypeUserList)
	}

	// remoteOrgList 远程机构列表
	remoteOrgList ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		return remoteserver.GetCacheByType(remoteserver.CacheTypeOrgList)
	}
)
