package cache

import "github.com/gin-gonic/gin"

func InitCache(engine *gin.Engine) {
	engineGroup := engine.Group("cache")
	initFileCache(engineGroup)
	initStrCache(engineGroup)
	initRemoteServerCache(engineGroup)

}
