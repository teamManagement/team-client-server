package remoteserver

var remoteCache = make(map[CacheType]string)

type CacheType string

const (
	// CacheTypeUserList 用户列表缓存
	CacheTypeUserList CacheType = "/cache/user/list"
	// CacheTypeOrgList 机构列表缓存
	CacheTypeOrgList CacheType = "/cache/org/list"
)

var cachePathList = []CacheType{CacheTypeUserList, CacheTypeOrgList}

func FlushAllCache() error {
	for _, cacheKey := range cachePathList {
		_ = FlushCacheByType(cacheKey)
	}
	return nil
}

func FlushCacheByType(t CacheType) error {
	var res string
	if err := RequestWebServiceWithResponse(string(t), &res); err != nil {
		remoteCache[t] = ""
		return err
	} else {
		remoteCache[t] = res
	}
	return nil
}

func GetCacheByType(t CacheType) string {
	res := remoteCache[t]
	if res == "" {
		res = "[]"
	}
	return res
}
