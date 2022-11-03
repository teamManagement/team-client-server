package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/go-base-lib/logs"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"team-client-server/errors"
	"team-client-server/sysuser"
)

type WatchRunFn func(config *Info)

var (
	configDirPath = filepath.Join(sysuser.HomeDir(), ".teamwork")
)

var (
	configFilePath   string
	databaseFilePath string
	logDirPath       string
	cacheFilePath    string
)

var (
	configWatchRunList = make([]WatchRunFn, 0, 16)

	CurrentConfig *Info
)

func initConfigAboutPath() {
	configFilePath = filepath.Join(configDirPath, "config.toml")
	databaseFilePath = filepath.Join(configDirPath, "db", "local.data")
	logDirPath = filepath.Join(configDirPath, "logs")
	cacheFilePath = filepath.Join(configDirPath, ".cache", "file")
}

func LoadConfig(configDir string) {
	if configDir != "" {
		configDirPath = configDir
	}

	initConfigAboutPath()

	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("toml")
	viper.SetDefault("database.path", databaseFilePath)
	viper.SetDefault("logs.level", "info")
	viper.SetDefault("logs.path", logDirPath)
	viper.SetDefault("cache.file", cacheFilePath)

	stat, err := os.Stat(configFilePath)
	if err != nil || stat.IsDir() {
		if err = os.MkdirAll(configDirPath, 0755); err != nil {
			errors.ExitConfigFileCreatEmpty.Println("创建空的默认配置文件失败: %s", err.Error())
		}

		if err = viper.WriteConfigAs(configFilePath); err != nil {
			errors.ExitConfigFileWriteToEmpty.Println("写出默认配置失败: %s", err.Error())
		}
	}

	if err = viper.ReadInConfig(); err != nil {
		errors.ExitConfigFileRead.Println("读取配置文件内容失败: %s", err.Error())
	}

	if err = viper.Unmarshal(&CurrentConfig); err != nil {
		errors.ExitConfigFileParser.Println("配置文件解析失败: %s", err.Error())
	}

	viper.OnConfigChange(onConfigChange)
	viper.WatchConfig()
}

// onConfigChange 配置文件发生改变
func onConfigChange(in fsnotify.Event) {
	if in.Op != fsnotify.Write {
		return
	}

	if err := viper.Unmarshal(&CurrentConfig); err != nil {
		logs.Warnf("配置文件已更新，但解析失败: %s", err.Error())
		return
	}

	for i := range configWatchRunList {
		fn := configWatchRunList[i]
		fn(CurrentConfig)
	}
}

func AddWatchFn(fn WatchRunFn) {
	configWatchRunList = append(configWatchRunList, fn)
}

func AddWatchAndNowExec(fn WatchRunFn) {
	fn(CurrentConfig)
	AddWatchFn(fn)
}

func CreateDirInConfigPath(dir string) (string, error) {
	dir = filepath.Join(configDirPath, dir)
	return dir, os.MkdirAll(dir, 0755)
}

func CreateFileInConfigPath(dir string, filename string) (*os.File, string, error) {
	dirPath := filepath.Join(configDirPath, dir)
	_ = os.MkdirAll(dirPath, 0755)

	fPath := filepath.Join(dirPath, filename)
	file, err := os.OpenFile(fPath, os.O_CREATE|os.O_RDWR, 0655)
	if err != nil {
		return nil, "", err
	}
	return file, fPath, nil
}
