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

func initAppLocalStore(engine *gin.RouterGroup) {
	engine.Group("store").
		POST("set/:appId", ginmiddleware.WrapperResponseHandle(appDataStoreSet)).
		POST("get/:appId", ginmiddleware.WrapperResponseHandle(appDataStoreGet)).
		POST("delete/:appId", ginmiddleware.WrapperResponseHandle(appDataStoreDelete)).
		POST("has/:appId", ginmiddleware.WrapperResponseHandle(appDataStoreHas)).
		POST("clear/:appId", ginmiddleware.WrapperResponseHandle(appDataStoreClear))
}

var (
	// appDataStoreSet 应用数据设置
	appDataStoreSet ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		var setting *vos.Setting
		if err := ctx.ShouldBindJSON(&setting); err != nil {
			return fmt.Errorf("解析配置数据失败: %s", err.Error())
		}
		if err := appDataStoreModel(ctx).Save(&setting).Error; err != nil {
			return fmt.Errorf("保存应用数据失败: %s", err.Error())
		}
		return nil
	}

	// appDataStoreGet 获取数据
	appDataStoreGet ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		name := ctx.Query("name")
		if name == "" {
			return errors.New("要查询的key不能为空")
		}

		valueSetting := &vos.Setting{}
		if err := appDataStoreModel(ctx).Select("value").Where("name=?", name).First(&valueSetting).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return fmt.Errorf("数据查询失败: %s", err.Error())
		}

		return valueSetting.Value
	}

	// appDataStoreDelete 数据删除
	appDataStoreDelete ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		name := ctx.Query("name")
		if name == "" {
			return errors.New("要查询的key不能为空")
		}
		if err := appDataStoreModel(ctx).Where("name=?", name).Delete(&vos.Setting{}).Error; err != nil {
			return fmt.Errorf("删除存储内容失败: %s", err.Error)
		}
		return nil
	}

	// appDataStoreHas 数据是否存在
	appDataStoreHas ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		name := ctx.Query("name")
		if name == "" {
			return errors.New("要查询的key不能为空")
		}

		count := 0
		appDataStoreModel(ctx).Where("name=?", name).Count(&count)
		return count > 0
	}

	// appDataStoreClear 数据清除
	appDataStoreClear ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		if err := appDataStoreModel(ctx).Delete(&vos.Setting{}).Error; err != nil {
			return fmt.Errorf("清空数据失败: %s", err.Error())
		}
		return nil
	}
)

func appDataStoreModel(ctx *gin.Context) *gorm.DB {
	userId := ""
	nowUserInfo := remoteserver.NowUser()
	if nowUserInfo != nil {
		userId = nowUserInfo.Id
	}
	return sqlite.Db().Table("app-" + userId + "-" + ctx.Param("appId"))
}
