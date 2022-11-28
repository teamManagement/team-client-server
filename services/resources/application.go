package resources

import (
	"fmt"
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"gorm.io/gorm"
	"team-client-server/remoteserver"
	"team-client-server/vos"
)

func InitApplicationWebResources(engine *gin.RouterGroup) {
	engine.POST("/app/list", ginmiddleware.WrapperResponseHandle(applicationListNoDebugging))
}

var (
	// applicationListNoDebugging 获取正常应用列表
	applicationListNoDebugging ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		nowUser, err := remoteserver.NowUser()
		if err != nil {
			return err
		}

		var applicationList []vos.Application
		if err := sqlite.Db().Model(&vos.Application{}).Where("user_id=? and not debugging and status=?", nowUser.Id, vos.ApplicationStatusNormal).Find(&applicationList).Error; err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("查询应用列表失败: %s", err.Error())
		}

		return applicationList
	}
)
