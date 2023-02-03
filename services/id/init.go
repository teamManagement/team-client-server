package id

import (
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"team-client-server/remoteserver"
	"team-client-server/sfnake"
)

func InitIdWebServices(engine *gin.RouterGroup) {
	engine.Group("id").
		POST("create", ginmiddleware.WrapperResponseHandle(idCreate))

}

var (
	// idCreate id创建
	idCreate ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		currentUser, err := remoteserver.NowUser()
		if err != nil {
			return err
		}

		idStr, err := sfnake.GetIdStr()
		if err != nil {
			return err
		}

		return currentUser.Id + "_" + idStr
	}
)
