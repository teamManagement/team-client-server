package services

import (
	"github.com/gin-gonic/gin"
	"team-client-server/services/appstore"
	"team-client-server/services/resources"
	"team-client-server/services/userchat"
)

func InitLocalWebServices(engine *gin.Engine) {
	servicesGroup := engine.Group("services")
	{
		resourcesGroup := servicesGroup.Group("resources")
		resources.InitApplicationWebResources(resourcesGroup)
	}

	appstore.InitAppStoreWebService(servicesGroup)
	userchat.InitUserChatWebService(servicesGroup)
}
