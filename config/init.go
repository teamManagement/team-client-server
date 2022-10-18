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
	configDirPath = filepath.Join(sysuser.HomeDir(), ".teamManager")
)

var (
	configFilePath   string
	databaseFilePath string
	logDirPath       string
)

var (
	configWatchRunList = make([]WatchRunFn, 0, 16)

	currentConfig *Info
)

func initConfigAboutPath() {
	configFilePath = filepath.Join(configDirPath, "config.toml")
	databaseFilePath = filepath.Join(configDirPath, "db", "local.data")
	logDirPath = filepath.Join(configDirPath, "logs")
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

	if err = viper.Unmarshal(&currentConfig); err != nil {
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

	if err := viper.Unmarshal(&currentConfig); err != nil {
		logs.Warnf("配置文件已更新，但解析失败: %s", err.Error())
		return
	}

	for i := range configWatchRunList {
		fn := configWatchRunList[i]
		fn(currentConfig)
	}
}

func AddWatchFn(fn WatchRunFn) {
	configWatchRunList = append(configWatchRunList, fn)
}

func AddWatchAndNowExec(fn WatchRunFn) {
	fn(currentConfig)
	AddWatchFn(fn)
}
