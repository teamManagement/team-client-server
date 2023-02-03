package id

import (
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"team-client-server/sfnake"
)

func InitId(engine *gin.Engine) {
	engine.Group("id").
		POST("create", ginmiddleware.WrapperResponseHandle(idCreate))
}

var (
	// idCreate id生成
	idCreate ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		str, err := sfnake.GetIdStr()
		if err != nil {
			return err
		}
		return str
	}
)
