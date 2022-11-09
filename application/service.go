package application

import (
	"errors"
	"fmt"
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"gorm.io/gorm"
	"strings"
	"team-client-server/remoteserver"
	"team-client-server/vos"
)

func initAppService(engine *gin.RouterGroup) {
	engine.
		POST("/force/refresh", ginmiddleware.WrapperResponseHandle(appForceRefresh)).
		POST("/info/desktop/list", ginmiddleware.WrapperResponseHandle(appInfoDesktopList)).
		POST("/info/:id", ginmiddleware.WrapperResponseHandle(appInfoGetById)).
		POST("/install/:appId", ginmiddleware.WrapperResponseHandle(appInstall)).
		POST("/uninstall/:appId", ginmiddleware.WrapperResponseHandle(appUninstall)).
		POST("/debug/install", ginmiddleware.WrapperResponseHandle(appDebugInstall)).
		POST("/debug/uninstall/:appId", ginmiddleware.WrapperResponseHandle(appDebugUninstall))
}

var (
	// appDebugUninstall 本地调试应用卸载
	appDebugUninstall ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {

		nowUser, err := remoteserver.NowUser()
		if err != nil {
			return err
		}

		appId := ctx.Param("appId")
		if appId == "" {
			return errors.New("未识别的应用ID")
		}
		if !strings.HasSuffix(appId, "-debug") {
			appId += "-debug"
		}
		if err = sqlite.Db().Model(&vos.Application{}).Where("id=? and user_id=?", appId, nowUser.Id).Delete(&vos.Application{}).Error; err != nil {
			return fmt.Errorf("删除应用调试信息失败: %s", err.Error())
		}
		return nil
	}

	// appDebugInstall 本地调试应用安装
	appDebugInstall ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		nowUser, err := remoteserver.NowUser()
		if err != nil {
			return err
		}

		var addAppInfo *vos.Application
		if err = ctx.ShouldBindJSON(&addAppInfo); err != nil {
			return fmt.Errorf("解析待调试应用信息失败: %s", err.Error())
		}

		if addAppInfo.Id == "" {
			return errors.New("应用ID不能为空")
		}

		if !strings.HasSuffix(addAppInfo.Id, "-debug") {
			addAppInfo.Id += "-debug"
		}

		if addAppInfo.Type == vos.ApplicationTypeLocalWeb {
			return errors.New("暂不支持本地文件的调试模式")
		}

		if addAppInfo.RemoteSiteUrl == "" {
			return errors.New("应用地址不能为空")
		}

		addAppInfo.Debugging = true
		addAppInfo.Url = addAppInfo.RemoteSiteUrl
		addAppInfo.UserId = nowUser.Id
		addAppInfo.Status = vos.ApplicationNormal

		return sqlite.Db().Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&vos.Application{}).Where("id=? and user_id=?", addAppInfo.Id, nowUser.Id).Save(&addAppInfo).Error; err != nil {
				return fmt.Errorf("保存应用调试信息失败: %s", err.Error())
			}

			if err := tx.Table(getAppStoreTableName(addAppInfo.Id, nowUser.Id)).AutoMigrate(&vos.Application{}); err != nil {
				return fmt.Errorf("创建应用store存储区失败: %s", err.Error())
			}

			return nil
		})
	}

	// appUninstall 应用卸载
	appUninstall ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		nowUser, err := remoteserver.NowUser()
		if err == nil {
			return err
		}
		appId := ctx.Param("appId")
		if appId == "" {
			return errors.New("获取应用ID失败")
		}

		if err := remoteserver.RequestWebService("/app/uninstall/" + appId); err != nil {
			return err
		}

		return sqlite.Db().Transaction(func(tx *gorm.DB) error {
			if err = tx.Migrator().DropTable(getAppStoreTableName(appId, nowUser.Id)); err != nil {
				return fmt.Errorf("删除应用存储失败: %s", err.Error())
			}
			if err = tx.Model(&vos.Application{}).Where("id=? and user_id=?", appId, nowUser.Id).Delete(&vos.Application{}).Error; err != nil {
				return fmt.Errorf("删除应用数据失败: %s", err.Error())
			}
			return nil
		})
	}
	// appInstall 应用安装
	appInstall ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		nowUser, err := remoteserver.NowUser()
		if err != nil {
			return err
		}
		appId := ctx.Param("appId")
		if appId == "" {
			return errors.New("获取应用ID失败")
		}

		count := int64(0)
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

		nowUser, err := remoteserver.NowUser()
		if err != nil {
			return err
		}
		if err := sqlite.Db().Transaction(func(tx *gorm.DB) error {
			appModel := tx.Model(&vos.Application{})
			appModel.Where("debugging is null or not debugging").UpdateColumns(&vos.Application{
				Status: vos.ApplicationStatusTakeDown,
			})
			for i := range appList {
				appInfo := appList[i]
				if appInfo.Type == vos.ApplicationTypeRemoteWeb {
					appInfo.Url = appInfo.RemoteSiteUrl
				}

				appInfo.UserId = nowUser.Id
				if err := tx.Model(&vos.Application{}).Where("id=?", appInfo.Id).UpdateColumns(appInfo).Error; err != nil {
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
		nowUser, err := remoteserver.NowUser()
		if err != nil {
			return err
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
