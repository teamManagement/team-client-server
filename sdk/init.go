package sdk

import (
	"github.com/gin-gonic/gin"
	"team-client-server/sdk/cache"
)

func InitLocalWebSdk(engine *gin.Engine) {
	cache.InitCache(engine)
}
