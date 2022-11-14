//go:build !windows

package updater

import (
	"github.com/go-base-lib/logs"
	"os/exec"
)

func updateStart(cmdPath string, workDir string, args []string) {
	logs.Info("UNIX程序更新UI重启...")
	cmd := exec.Command(cmdPath, args...)
	if workDir != "" {
		cmd.Dir = workDir
	}
	output, _ := cmd.CombinedOutput()
	logs.Infof("electron重启之后的输出: %s", string(output))
}
