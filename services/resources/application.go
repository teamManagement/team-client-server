package resources

import (
	"fmt"
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"gorm.io/gorm"
	"team-client-server/db"
	"team-client-server/remoteserver"
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

		if nowUser.Id == "0" {
			return "[]"
		}

		var applicationList []db.Application
		if err := sqlite.Db().Model(&db.Application{}).Where("user_id=? and not debugging and status=?", nowUser.Id, db.ApplicationStatusNormal).Find(&applicationList).Error; err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("查询应用列表失败: %s", err.Error())
		}

		return applicationList
	}
)
