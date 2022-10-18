package logs

import (
	"github.com/go-base-lib/logs"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"team-client-server/config"
	"team-client-server/errors"
	"time"
)

var (
	dir   string
	level logrus.Level
)

func InitLog() {
	config.AddWatchAndNowExec(configChange)
}

func configChange(config *config.Info) {
	logLevel := config.Logs.Level.Level()
	if config.Logs.Path != dir {
		dir = config.Logs.Path
		if stat, err := os.Stat(dir); err != nil || !stat.IsDir() {
			if err = os.MkdirAll(dir, 0755); err != nil {
				errors.ExitLogDirCreate.Println("创建日志目录失败: %s", err.Error())
			}
		}

		level = logLevel
		logFile := filepath.Join(dir, "app.log")
		if err := logs.InitDefault(&logs.Config{
			CurrentLevel: level,
			Formatter: &logrus.TextFormatter{
				DisableQuote:    true,
				TimestampFormat: "2006-01-02 15:04:05",
			},
			PathConfig: &logs.PathConfig{
				LogPath: logFile + ".%Y%m%d",
				//LinkName:     logFile,
				MaxAge:       24 * 30 * time.Hour,
				RotationTime: 24 * time.Hour,
			},
		}); err != nil {
			logs.Error("日志初始化失败")
		}
	} else if level != logLevel {
		level = logLevel
		logs.SetCurrentLevel(level)
	}
}
