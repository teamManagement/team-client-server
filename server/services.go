package server

import (
	"fmt"
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"team-client-server/remoteserver"
	"team-client-server/vos"
)

func initLocalService(engine *gin.Engine) {
	{
		engine.
			POST("/app/force/refresh", ginmiddleware.WrapperResponseHandle(appForceRefresh)).
			POST("/app/info/desktop/list", ginmiddleware.WrapperResponseHandle(appInfoDesktopList)).
			POST("/app/info/:id", ginmiddleware.WrapperResponseHandle(appInfoGetById))
	}

	{
		engine.
			POST("/user/now", ginmiddleware.WrapperResponseHandle(userNowInfo)).
			POST("/user/status", ginmiddleware.WrapperResponseHandle(userNowStatus))
	}

}

var (
	userNowInfo ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		return remoteserver.NowUser()
	}

	userNowStatus ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		if remoteserver.LoginOk() {
			return "online"
		}
		return "offline"
	}
)

var (
	// appForceRefresh 应用列表强制刷新
	appForceRefresh ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		var appList []*vos.Application
		if err := remoteserver.RequestWebServiceWithResponse("/app/list", &appList); err != nil {
			return fmt.Errorf("获取应用列表失败: %s", err.Error())
		}

		if err := sqlite.Db().Transaction(func(tx *gorm.DB) error {
			appModel := tx.Model(&vos.Application{})
			appModel.Delete(&vos.Application{})
			for i := range appList {
				if err := appModel.Create(appList[i]).Error; err != nil {
					return fmt.Errorf("保存应用信息失败: %s", err.Error())
				}
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
