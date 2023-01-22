package db

import (
	"fmt"
	dbutils "github.com/byzk-worker/go-db-utils"
	"github.com/byzk-worker/go-db-utils/sqlite"
	"github.com/go-base-lib/logs"
	"os"
	"path/filepath"
	"team-client-server/config"
)

var dbFilePath = ""

func InitDb() {
	config.AddWatchAndNowExec(configChange)
}

func configChange(config *config.Info) {
	dbPath := config.Database.Path
	if dbFilePath == dbPath {
		return
	}

	dbFilePath = dbPath

	_ = os.MkdirAll(filepath.Dir(dbFilePath), 0755)

	if err := sqlite.Init(fmt.Sprintf("file:%s?auto_vacuum=1", dbFilePath), dbutils.DefaultGetContextFn); err != nil {
		logs.Errorf("初始化数据存储文件失败: %s", err.Error())
		os.Exit(9)
	}

	sqlite.EnableDebug()

	initDataTable()
}

func initDataTable() {
	if err := sqlite.Db().
		AutoMigrate(&Setting{},
			&Setting{},
			&ProxyHttpServerInfo{},
			&ProxyHttpResponseCache{},
			&UserChatMsg{},
			&ChatGroupInfo{},
			&QueueChannelMsgInfo{},
			&Application{}); err != nil {
		logs.Fatalf("初始化数据库信息失败: %s")
		os.Exit(10)
	}
}
