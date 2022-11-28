package application

import (
	"github.com/gin-gonic/gin"
)

func InitAppService(engine *gin.Engine) {
	appEngineGroup := engine.Group("app")
	initAppService(appEngineGroup)
	initAppLocalStore(appEngineGroup)
}
