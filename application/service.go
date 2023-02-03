package application

import (
	"errors"
	"fmt"
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"gorm.io/gorm"
	"strings"
	"team-client-server/db"
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
		POST("/debug/uninstall/:appId", ginmiddleware.WrapperResponseHandle(appDebugUninstall)).
		POST("/try/create/db/again/:appId", ginmiddleware.WrapperResponseHandle(appTryCreateRemoteDbAgain))

	remoteserver.RegistryInsideEvent(remoteserver.InsideEventNameFlushAppList, func(userInfo *vos.UserInfo) error {
		_, err := appForceFlushDesktop(userInfo)
		return err
	})
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

		var addAppInfo *db.Application
		if err = ctx.ShouldBindJSON(&addAppInfo); err != nil {
			return fmt.Errorf("解析待调试应用信息失败: %s", err.Error())
		}

		if addAppInfo.Id == "" {
			return errors.New("应用ID不能为空")
		}

		if !strings.HasSuffix(addAppInfo.Id, "-debug") {
			addAppInfo.Id += "-debug"
		}

		if addAppInfo.Type == db.ApplicationTypeLocalWeb {
			return errors.New("暂不支持本地文件的调试模式")
		}

		if addAppInfo.RemoteSiteUrl == "" {
			return errors.New("应用地址不能为空")
		}

		addAppInfo.Debugging = true
		addAppInfo.Url = addAppInfo.RemoteSiteUrl
		addAppInfo.UserId = nowUser.Id
		addAppInfo.Status = db.ApplicationStatusNormal

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
		if err := sqlite.Db().Model(&db.Application{}).Where("id=? and user_id=?", appId, nowUser.Id).Count(&count).Error; err != nil {
			return fmt.Errorf("查询应用安装情况失败: %s", err.Error())
		}

		if count != 0 {
			return errors.New("应用信息已安装")
		}

		var appInfo *db.Application
		if err := remoteserver.RequestWebServiceWithResponse("/app/install/"+appId, &appInfo); err != nil {
			return err
		}

		return sqlite.Db().Transaction(func(tx *gorm.DB) error {
			return appHandlerInstallInfo(tx, appInfo, nowUser.Id)
		})
	}

	// appForceRefresh 应用列表强制刷新
	appForceRefresh ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {

		nowUser, err := remoteserver.NowUser()
		if err != nil {
			return err
		}

		if appList, err := appForceFlushDesktop(nowUser); err != nil {
			return err
		} else {
			return appList
		}

	}

	// appInfoDesktopList 应用桌面信息列表获取
	appInfoDesktopList ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		nowUser, err := remoteserver.NowUser()
		if err != nil {
			return err
		}
		appList := make([]*db.Application, 0)
		if err := sqlite.Db().Where("status=? and user_id=?", db.ApplicationStatusNormal, nowUser.Id).Find(&appList).Error; err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("查询应用列表失败: %s", err.Error())
		}

		for i := range appList {
			appInfo := appList[i]
			if !appInfo.HaveRemoteDb && !appInfo.Debugging {
				_ = sqlite.Db().Transaction(func(tx *gorm.DB) error {
					if err = tx.Model(&db.Application{}).Where("id=? and user_id=? and not debugging", appInfo.Id, appInfo.UserId).Update("have_remote_db", true).Error; err != nil {
						return err
					}
					if err := remoteserver.RequestWebService("/app/db/create/" + appInfo.Id); err != nil {
						return err
					}
					return nil
				})

			}
		}

		return appList
	}

	// appInfoGetById 应用信息获取通过ID
	appInfoGetById ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		return nil
	}

	appTryCreateRemoteDbAgain ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		nowUser, err := remoteserver.NowUser()
		if err != nil {
			return err
		}

		appId := ctx.Param("appId")
		return sqlite.Db().Transaction(func(tx *gorm.DB) error {
			appModalWhere := tx.Model(&db.Application{}).Where("id=? and user_id=? and status='4' and not debugging", appId, nowUser.Id)

			var count int64
			if err := appModalWhere.Count(&count).Error; err != nil {
				return fmt.Errorf("查询应用信息失败: %s", err.Error())
			}

			if err := remoteserver.RequestWebService("/app/db/create/" + appId); err != nil {
				return err
			}

			if err := appModalWhere.UpdateColumn("have_remote_db", true).Error; err != nil {
				return fmt.Errorf("更新应用状态失败: %s", err.Error())
			}

			return nil
		})
	}
)

func appForceFlushDesktop(nowUser *vos.UserInfo) ([]*db.Application, error) {
	var appList []*db.Application
	if err := remoteserver.RequestWebServiceWithResponse("/app/list", &appList); err != nil {
		return nil, fmt.Errorf("获取应用列表失败: %s", err.Error())
	}

	if err := sqlite.Db().Transaction(func(tx *gorm.DB) error {
		appModel := tx.Model(&db.Application{})
		appModel.Where("debugging is null or not debugging").UpdateColumns(&db.Application{
			Status: db.ApplicationStatusTakeDown,
		})
		for i := range appList {
			appInfo := appList[i]
			if err := appHandlerInstallInfo(tx, appInfo, nowUser.Id); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return appList, nil
}

// appHandlerUninstallInfo 处理应用
func appHandlerUninstallInfo(gormDb *gorm.DB, userId, appId string) error {
	if userId == "" {
		return errors.New("缺失当前用户信息")
	}

	if appId == "" {
		return errors.New("缺失应用信息")
	}

	if err := gormDb.Migrator().DropTable(getAppStoreTableName(appId, userId)); err != nil {
		return fmt.Errorf("删除应用存储失败: %s", err.Error())
	}

	if err := gormDb.Model(&db.Application{}).Where("id=? and user_id=?", appId, userId).Delete(&db.Application{}).Error; err != nil {
		return fmt.Errorf("删除应用信息失败: %s", err.Error())
	}
	return nil
}

// appHandlerInstallInfo 处理应用的安装信息
func appHandlerInstallInfo(gormDb *gorm.DB, appInfo *db.Application, userId string) (err error) {
	//defer func() {
	//	if err != nil {
	//		_ = appHandlerUninstallInfo(gormDb, userId, appInfo.Id)
	//	}
	//}()
	if appInfo == nil || appInfo.Id == "" {
		return errors.New("缺失应用信息或应用ID")
	}

	if userId == "" {
		return errors.New("缺失当前用户信息")
	}

	appStoreTableName := getAppStoreTableName(appInfo.Id, userId)
	appModel := gormDb.Model(&db.Application{})
	switch appInfo.Type {
	case db.ApplicationTypeRemoteWeb:
		appInfo.Url = appInfo.RemoteSiteUrl
	case db.ApplicationTypeLocalWeb:
		return errors.New("暂不支持本地应用模式")
	default:
		return errors.New("未知的应用模式")
	}
	appInfo.UserId = userId

	if err := gormDb.Table(appStoreTableName).AutoMigrate(&db.Setting{}); err != nil {
		return fmt.Errorf("创建应用存储失败: %s", err.Error())
	}

	if err := appModel.Where("id=? and user_id=?", appInfo.Id, userId).Save(appInfo).Error; err != nil {
		return fmt.Errorf("应用信息保存失败: %s", err.Error())
	}

	//if len(appInfo.LocalFileHash) > 0 {
	//	req, err := remoteserver.RequestWebServiceToRawReq("/app/download/file/"+appInfo.Id, nil)
	//	if err != nil {
	//		return fmt.Errorf("请求程序文件失败: %w", err)
	//	}
	//
	//	req.ResponseWithHandler(func(res *http.Response) error {
	//
	//	})
	//
	//}

	return nil
}
