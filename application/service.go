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
		POST("/info/:id", ginmiddleware.WrapperResponseHandle(appInfoGetById)).
		POST("/install/:appId", ginmiddleware.WrapperResponseHandle(appInstall)).
		POST("/uninstall/:appId", ginmiddleware.WrapperResponseHandle(appUninstall))
}

var (
	appUninstall ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		nowUser := remoteserver.NowUser()
		if nowUser == nil {
			return errors.New("获取当前登录用户信息失败")
		}
		appId := ctx.Param("appId")
		if appId == "" {
			return errors.New("获取应用ID失败")
		}

		if err := remoteserver.RequestWebService("/app/uninstall/" + appId); err != nil {
			return err
		}

		sqlite.Db().Model(&vos.Application{}).Where("id=? and user_id=?", appId, nowUser.Id).Delete(&vos.Application{})
		return nil
	}
	// appInstall 应用安装
	appInstall ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		nowUser := remoteserver.NowUser()
		if nowUser == nil {
			return errors.New("获取当前登录用户信息失败")
		}
		appId := ctx.Param("appId")
		if appId == "" {
			return errors.New("获取应用ID失败")
		}

		count := 0
		if err := sqlite.Db().Model(&vos.Application{}).Where("id=? and user_id=?", appId, nowUser.Id).Count(&count).Error; err != nil {
			return fmt.Errorf("查询应用安装情况失败: %s", err.Error())
		}

		if count != 0 {
			return errors.New("应用信息已安装")
		}

		var appInfo *vos.Application
		if err := remoteserver.RequestWebServiceWithResponse("/app/install/"+appId, &appInfo); err != nil {
			return err
		}

		appInfo.UserId = nowUser.Id
		if err := sqlite.Db().Model(&vos.Application{}).Where("id=? and user_id=?", appId, nowUser.Id).Delete(&vos.Application{}).
			Create(appInfo).Error; err != nil {
			return fmt.Errorf("保存应用信息失败: %s", err.Error())
		}

		return nil
	}

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

				appInfo.UserId = nowUser.Id
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
		nowUser := remoteserver.NowUser()
		if nowUser == nil {
			return fmt.Errorf("获取当前用户信息失败")
		}
		appList := make([]*vos.Application, 0)
		if err := sqlite.Db().Where("status=? and user_id=?", vos.ApplicationNormal, nowUser.Id).Find(&appList).Error; err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("查询应用列表失败: %s", err.Error())
		}
		return appList
	}

	// appInfoGetById 应用信息获取通过ID
	appInfoGetById ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		return nil
	}
)
