package application

import (
	"errors"
	"fmt"
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"team-client-server/remoteserver"
	"team-client-server/vos"
)

func initAppService(engine *gin.RouterGroup) {
	engine.
		POST("/force/refresh", ginmiddleware.WrapperResponseHandle(appForceRefresh)).
		POST("/info/desktop/list", ginmiddleware.WrapperResponseHandle(appInfoDesktopList)).
		POST("/info/:id", ginmiddleware.WrapperResponseHandle(appInfoGetById))
}

var (
	// appForceRefresh 应用列表强制刷新
	appForceRefresh ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		var appList []*vos.Application
		if err := remoteserver.RequestWebServiceWithResponse("/app/list", &appList); err != nil {
			return fmt.Errorf("获取应用列表失败: %s", err.Error())
		}

		nowUser := remoteserver.NowUser()
		if nowUser == nil {
			return errors.New("获取当前用户的登录信息失败")
		}
		if err := sqlite.Db().Transaction(func(tx *gorm.DB) error {
			appModel := tx.Model(&vos.Application{})
			appModel.Where("debugging is null or not debugging").Delete(&vos.Application{})
			for i := range appList {
				appInfo := appList[i]
				if appInfo.Type == vos.ApplicationTypeRemoteWeb {
					appInfo.Url = appInfo.RemoteSiteUrl
				}

				if err := appModel.Create(appInfo).Error; err != nil {
					return fmt.Errorf("保存应用信息失败: %s", err.Error())
				}
				tx.Table("app-" + nowUser.Id + "-" + appInfo.Id).AutoMigrate(&vos.Setting{})
			}
			return nil
		}); err != nil {
			return err
		}

		return appList
	}

	// appInfoDesktopList 应用桌面信息列表获取
	appInfoDesktopList ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		appList := make([]*vos.Application, 0)
		if err := sqlite.Db().Where("status=?", vos.ApplicationNormal).Find(&appList).Error; err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("查询应用列表失败: %s", err.Error())
		}
		return appList
	}

	// appInfoGetById 应用信息获取通过ID
	appInfoGetById ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		return nil
	}
)
