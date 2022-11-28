package appstore

import (
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"team-client-server/remoteserver"
)

func InitAppStoreWebService(engine *gin.RouterGroup) {
	engine.Group("appstore").
		POST("/manager/add/:userId", ginmiddleware.WrapperResponseHandle(appstoreAddManager)).
		POST("/manager/del/:userId", ginmiddleware.WrapperResponseHandle(appstoreDelManager))
}

var (
	// appstoreAddManager 应用商店增加管理员
	appstoreAddManager ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		if err := remoteserver.RequestWebService("/appstore/manager/add/" + ctx.Param("userId")); err != nil {
			return err
		}
		_ = remoteserver.FlushCacheByType(remoteserver.CacheTypeUserList)
		return nil
	}

	// appstoreDelManager 删除应用商店管理员
	appstoreDelManager ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		if err := remoteserver.RequestWebService("/appstore/manager/del/" + ctx.Param("userId")); err != nil {
			return err
		}
		_ = remoteserver.FlushCacheByType(remoteserver.CacheTypeUserList)
		return nil
	}
)
