package sdk

import (
	"github.com/gin-gonic/gin"
	"team-client-server/sdk/cache"
	"team-client-server/sdk/id"
)

func InitLocalWebSdk(engine *gin.Engine) {
	cache.InitCache(engine)
	id.InitId(engine)
}
