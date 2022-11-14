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

		return sqlite.Db().Transaction(func(tx *gorm.DB) error {
			return appHandlerUninstallInfo(tx, nowUser.Id, appId)
		})

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
			return appHandlerInstallInfo(tx, addAppInfo, nowUser.Id)
		})
	}

	// appUninstall 应用卸载
	appUninstall ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		nowUser, err := remoteserver.NowUser()
		if err != nil {
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
			return appHandlerUninstallInfo(tx, nowUser.Id, appId)
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

		return sqlite.Db().Transaction(func(tx *gorm.DB) error {
			return appHandlerInstallInfo(tx, appInfo, nowUser.Id)
		})
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
				if err := appHandlerInstallInfo(tx, appInfo, nowUser.Id); err != nil {
					return err
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

// appHandlerUninstallInfo 处理应用
func appHandlerUninstallInfo(db *gorm.DB, userId, appId string) error {
	if userId == "" {
		return errors.New("缺失当前用户信息")
	}

	if appId == "" {
		return errors.New("缺失应用信息")
	}

	if err := db.Migrator().DropTable(getAppStoreTableName(appId, userId)); err != nil {
		return fmt.Errorf("删除应用存储失败: %s", err.Error())
	}

	if err := db.Model(&vos.Application{}).Where("id=? and user_id=?", appId, userId).Delete(&vos.Application{}).Error; err != nil {
		return fmt.Errorf("删除应用信息失败: %s", err.Error())
	}
	return nil
}

// appHandlerInstallInfo 处理应用的安装信息
func appHandlerInstallInfo(db *gorm.DB, appInfo *vos.Application, userId string) error {
	if appInfo == nil || appInfo.Id == "" {
		return errors.New("缺失应用信息或应用ID")
	}

	if userId == "" {
		return errors.New("缺失当前用户信息")
	}

	appStoreTableName := getAppStoreTableName(appInfo.Id, userId)
	appModel := db.Model(&vos.Application{})
	switch appInfo.Type {
	case vos.ApplicationTypeRemoteWeb:
		appInfo.Url = appInfo.RemoteSiteUrl
	case vos.ApplicationTypeLocalWeb:
		return errors.New("暂不支持本地应用模式")
	default:
		return errors.New("未知的应用模式")
	}
	appInfo.UserId = userId

	if err := db.Table(appStoreTableName).AutoMigrate(&vos.Setting{}); err != nil {
		return fmt.Errorf("创建应用存储失败: %s", err.Error())
	}

	if err := appModel.Where("id=? and user_id=?", appInfo.Id, userId).Save(appInfo).Error; err != nil {
		return fmt.Errorf("应用信息保存失败: %s", err.Error())
	}

	return nil
}
