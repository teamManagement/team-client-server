package server

import (
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"os/exec"
	"team-client-server/application"
	"team-client-server/remoteserver"
)

func initLocalService(engine *gin.Engine) {
	application.InitAppService(engine)
	{
		engine.Group("/user").
			POST("/now", ginmiddleware.WrapperResponseHandle(userNowInfo)).
			POST("/status", ginmiddleware.WrapperResponseHandle(userNowStatus))
	}

	{
		engine.Group("/exec").
			POST("/lookPath/:name", ginmiddleware.WrapperResponseHandle(execLookPath))
	}

}

var (
	userNowInfo ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		nowUserInfo, err := remoteserver.NowUser()
		if err != nil {
			return err
		}
		return nowUserInfo
	}

	userNowStatus ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		if remoteserver.LoginOk() {
			return "online"
		}
		return "offline"
	}
)

var (
	// execLookPath 查找二进制文件路径
	execLookPath ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		name := ctx.Param("name")
		p, err := exec.LookPath(name)
		if err != nil {
			return err
		}
		return p
	}
)
