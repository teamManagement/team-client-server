//go:build !windows

package updater

import (
	"github.com/go-base-lib/logs"
	"os/exec"
	"strings"
	"syscall"
)

func updateStart(updateInfo *UpdateInfo, args []string) {
	cmdPath := updateInfo.Exec
	//workDir := updateInfo.WorkDir
	logs.Info("UNIX程序更新UI重启...")
	logs.Infof("更新执行命令: %s, 参数: %s", cmdPath, strings.Join(args, " "))
	cmd := exec.Command(cmdPath, args...)
	cmd.Env = append(cmd.Env, "DISPLAY="+updateInfo.Display)
	logs.Infof("环境变量: %s", strings.Join(cmd.Env, " "))
	//if workDir != "" {
	//	cmd.Dir = workDir
	//}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(updateInfo.Uid),
			Gid: uint32(updateInfo.Gid),
		},
	}
	output, err := cmd.CombinedOutput()
	logs.Infof("electron重启之后的输出: %s", string(output))
	if err != nil {
		logs.Warningf("electron应用重启失败: %s", err.Error())
	}
}
