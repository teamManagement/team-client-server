package updater

import (
	"github.com/go-base-lib/logs"
	"path/filepath"
	"strings"
	"team-client-server/tools"
)

func updateStart(cmdPath string, workDir string, args []string) {
	logs.Info("Windows程序更新UI重启...")
	if workDir == "" {
		workDir = cmdPath
	}

	workDir = filepath.Dir(workDir)
	logs.Info("windows更新程序工作路径: ", workDir)
	err := tools.StartProcessAsCurrentUser(cmdPath, strings.Join(args, " "), workDir, true)
	errMsg := "无错误"
	if err != nil {
		errMsg = err.Error()
	}
	logs.Info("windows下程序调用结束, 错误信息: ", errMsg)
}
