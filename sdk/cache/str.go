package cache

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"strings"
	"time"
)

var strCache = cache.New(10*time.Minute, 5*time.Minute)

type strCacheInfo struct {
	Key    string        `json:"key,omitempty"`
	Value  string        `json:"value,omitempty"`
	Expire time.Duration `json:"expire,omitempty"`
}

func initStrCache(engine *gin.RouterGroup) {
	engine.Group("str").
		POST("set/:appId", ginmiddleware.WrapperResponseHandle(cacheStrSet)).
		POST("get/:appId", ginmiddleware.WrapperResponseHandle(cacheStrGet)).
		POST("delete/:appId", ginmiddleware.WrapperResponseHandle(cacheStrDelete)).
		POST("has/:appId", ginmiddleware.WrapperResponseHandle(cacheStrHas)).
		POST("clear/:appId", ginmiddleware.WrapperResponseHandle(cacheStrClear)).
		POST("delay/:appId", ginmiddleware.WrapperResponseHandle(cacheStrDelay))
}

var (
	// cacheStrSet 缓存设置
	cacheStrSet ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		var cacheInfo *strCacheInfo
		if err := ctx.ShouldBindJSON(&cacheInfo); err != nil {
			return fmt.Errorf("解析缓存设置失败: %s", err.Error())
		}

		appId := ctx.Param("appId")
		if appId == "" {
			return errors.New("缺失应用ID")
		}

		if cacheInfo.Key == "" {
			return errors.New("缓存要设置的key名称不能为空")
		}

		key, err := strCacheKeyGet(ctx, cacheInfo.Key)
		if err != nil {
			return err
		}

		if cacheInfo.Expire < 0 {
			cacheInfo.Expire = 0
		}

		cacheInfo.Expire = time.Millisecond * cacheInfo.Expire

		strCache.Set(key, cacheInfo.Value, cacheInfo.Expire)
		return nil
	}

	// cacheStrGet 缓存获取
	cacheStrGet ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		key, err := strCacheKeyGet(ctx, ctx.Query("k"))
		if err != nil {
			return err
		}

		v, ok := strCache.Get(key)
		if !ok {
			return nil
		}
		return v.(string)
	}

	// cacheStrDelete 缓存删除
	cacheStrDelete ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		key, err := strCacheKeyGet(ctx, ctx.Query("k"))
		if err != nil {
			return err
		}

		strCache.Delete(key)
		return nil
	}

	// cacheStrHash 缓存KEY是否存在
	cacheStrHas ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		key, err := strCacheKeyGet(ctx, ctx.Query("k"))
		if err != nil {
			return err
		}

		_, ok := strCache.Get(key)

		return ok
	}

	// cacheStrClear 缓存清除
	cacheStrClear ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		appId := ctx.Param("appId")
		if appId == "" {
			return errors.New("获取应用ID失败")
		}

		keyPrefix := appId + "-"

		items := strCache.Items()
		for k := range items {
			if strings.HasPrefix(k, keyPrefix) {
				strCache.Delete(k)
			}
		}
		return nil
	}

	// cacheStrDelay 缓存延期
	cacheStrDelay ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		var cacheInfo *strCacheInfo
		if err := ctx.ShouldBindJSON(&cacheInfo); err != nil {
			return fmt.Errorf("解析缓存设置失败: %s", err.Error())
		}

		key, err := strCacheKeyGet(ctx, cacheInfo.Key)
		if err != nil {
			return err
		}
		v, ok := strCache.Get(key)
		if !ok {
			return false
		}

		if cacheInfo.Expire < 0 {
			cacheInfo.Expire = 0
		}
		cacheInfo.Expire = cacheInfo.Expire * time.Millisecond

		strCache.Set(key, v, cacheInfo.Expire)

		return true
	}
)

func strCacheKeyGet(ctx *gin.Context, key string) (string, error) {
	appId := ctx.Param("appId")
	if appId == "" {
		return "", errors.New("获取应用ID失败")
	}
	if key == "" {
		return "", errors.New("缓存KEY不能为空")
	}
	return appId + "-" + key, nil
}
