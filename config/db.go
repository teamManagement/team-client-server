package config

import (
	"github.com/byzk-worker/go-db-utils/sqlite"
	"sync"
	"team-client-server/vos"
)

var lock = sync.Mutex{}
var dbSettingMap = make(map[string]string)

func GetDbSetting(key string) (string, bool) {
	lock.Lock()
	defer lock.Unlock()
	if val, ok := dbSettingMap[key]; ok {
		return val, true
	}

	var dbSetting *vos.Setting
	if err := sqlite.Db().Model(&dbSetting).Where("key=?", key).First(&dbSetting).Error; err != nil {
		return "", false
	}

	val := string(dbSetting.Value)
	dbSettingMap[key] = val
	return val, true
}

func SetDbSetting(key string, val string) error {
	lock.Lock()
	defer lock.Unlock()

	setting := &vos.Setting{
		Name:  key,
		Value: vos.EncryptValue(val),
	}
	if err := sqlite.Db().Model(&setting).Save(setting).Error; err != nil {
		return err
	}

	dbSettingMap[key] = val
	return nil
}
