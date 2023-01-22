package updater

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-base-lib/coderutils"
	"github.com/go-base-lib/logs"
	ginmiddleware "github.com/teamManagement/gin-middleware"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"team-client-server/config"
	"team-client-server/tools"
	"time"
)

type UpdateInfo struct {
	ServerExePath      string `json:"serverExePath,omitempty"`
	Exec               string `json:"exec,omitempty"`
	Asar               string `json:"asar,omitempty"`
	WillUpdateAsarPath string `json:"willUpdateAsarPath,omitempty"`
	Uid                int    `json:"uid,omitempty"`
	Gid                int    `json:"gid,omitempty"`
	WorkDir            string `json:"workDir,omitempty"`
	Debug              bool   `json:"debug,omitempty"`
	Display            string `json:"display,omitempty"`
}

func InitUpdaterHttpRestful(engine *gin.Engine) {
	engine.Group("updater").
		POST("/check/:version", ginmiddleware.WrapperResponseHandle(updaterCheck)).
		POST("/update", ginmiddleware.WrapperResponseHandle(updaterUpdate))
}

var (
	// updaterCheck 更新检查
	updaterCheck ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		currentVersion := ctx.Param("version")
		if currentVersion == "" {
			return errors.New("当前版本号不能为空")
		}

		serverVersion := ""

		if err := requestToRemoteServer("/release/current/version", func(resp *http.Response) error {
			all, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			serverRes := &ginmiddleware.HttpResult{}
			if err := json.Unmarshal(all, &serverRes); err != nil || serverRes.Error {
				return errors.New("服务器查询版本失败: %s" + err.Error())
			}

			ok := false
			serverVersion, ok = serverRes.Result.(string)
			if !ok {
				return errors.New("数据转换失败")
			}
			return nil
		}); err != nil {
			return err
		}

		if serverVersion == "" || serverVersion == currentVersion {
			return nil
		}

		updaterStoreDir, _ := config.CreateDirInConfigPath(".updater")

		if err := downloadReleasePackage("asar", filepath.Join(updaterStoreDir, "asar")); err != nil {
			return err
		}

		p := filepath.Join(updaterStoreDir, "clientServer")
		if err := downloadReleasePackage("clientServer", p); err != nil {
			return err
		}

		if runtime.GOOS != "windows" {
			targetCopyPath := ctx.Query("_uPath")
			f, err := os.OpenFile(p, os.O_RDONLY, 0777)
			if err != nil {
				return fmt.Errorf("文件打开失败: %s", err.Error())
			}
			defer f.Close()

			tf, err := os.OpenFile(targetCopyPath, os.O_WRONLY|os.O_CREATE, 0777)
			if err != nil {
				return fmt.Errorf("打开目的地址失败: %s", err.Error())
			}
			defer tf.Close()

			if _, err = io.Copy(tf, f); err != nil {
				return fmt.Errorf("拷贝更新文件到目标失败: %s", err.Error())
			}

			if err = os.Chmod(p, 0755); err != nil {
				return fmt.Errorf("重新付权失败: %s", err.Error())
			}

			resourceDir := filepath.Dir(targetCopyPath)
			if err = os.Chmod(filepath.Join(resourceDir, "teamClientServer"), 0777); err != nil {
				return fmt.Errorf("更改源teamClientServer文件权限失败: %s", err.Error())
			}

			if err = os.Chmod(targetCopyPath, 0777); err != nil {
				return fmt.Errorf("更改更新服务包失败: %s", err.Error())
			}

			if err = os.Chmod(filepath.Dir(targetCopyPath), 0777); err != nil {
				return fmt.Errorf("更改资源目录权限失败: %s", err.Error())
			}

		}
		return p
	}
	// updaterUpdate 执行更新
	updaterUpdate ginmiddleware.ServiceFun = func(ctx *gin.Context) interface{} {
		logs.Info("进入热更新方法")
		var updateInfo *UpdateInfo
		if err := ctx.ShouldBindJSON(&updateInfo); err != nil {
			logs.Errorf("解析更新数据失败: %s", err.Error())
			return err
		}

		if execStat, err := os.Stat(updateInfo.Exec); err != nil || execStat.IsDir() {
			logs.Errorf("解析更新可执行文件失败: %s", err.Error())
			return err
		}

		if updateInfo.Asar == "" {
			return errors.New("缺失asar目录")
		}

		updaterStoreDir, _ := config.CreateDirInConfigPath(".updater")
		willUpdateAsarFilePath := filepath.Join(updaterStoreDir, "asar")
		stat, err := os.Stat(willUpdateAsarFilePath)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return errors.New("获取更新文件失败")
		}

		if runtime.GOOS != "windows" {
			if err = os.Chown(willUpdateAsarFilePath, updateInfo.Uid, updateInfo.Gid); err != nil {
				return err
			}
		}

		if updateInfo.Debug {
			if stat, err := os.Stat(updateInfo.WorkDir); err != nil || !stat.IsDir() {
				return errors.New("工作目录不正确")
			}
		}

		logs.Info("更新数据检查完成, 准备进行热更替。。。")
		//if runtime.GOOS != "windows" {
		go func() {
			time.Sleep(3 * time.Second)
			defer func() {
				logs.Info("资源文件拷贝完成, 开始调用程序启动")
				_ = os.Remove(willUpdateAsarFilePath)
				args := make([]string, 0, 3)
				args = append(args, updateInfo.WorkDir)
				args = append(args, "__updater_start__")
				if updateInfo.Debug {
					args = append(args, "__debug_work_dir__="+updateInfo.WorkDir)
				}

				updateStart(updateInfo, args)
				//cmd := exec.Command(updateInfo.Exec, args...)
				//if updateInfo.Debug {
				//	cmd.Dir = updateInfo.WorkDir
				//}
				//output, _ := cmd.CombinedOutput()
				//logs.Infof("electron重启之后的输出: %s", string(output))
			}()

			logs.Info("开始拷贝资源文件...")
			f, err := os.OpenFile(willUpdateAsarFilePath, os.O_RDONLY, 0777)
			if err != nil {
				return
			}
			defer f.Close()

			destF, err := os.OpenFile(updateInfo.Asar, os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				return
			}
			defer destF.Close()

			_, _ = io.Copy(destF, f)
		}()
		//} else {
		//	return willUpdateAsarFilePath
		//}
		//else {
		//	if updateInfo.ServerExePath == "" {
		//		return errors.New("缺失可执行文件路径")
		//	}
		//	updateInfo.WillUpdateAsarPath = willUpdateAsarFilePath
		//
		//	marshal, _ := json.Marshal(updateInfo)
		//	fmt.Println("创建文件拷贝子进程...")
		//	_, _ = os.StartProcess(updateInfo.ServerExePath, []string{"-cmd=updater", "-updateInfo=" + base64.StdEncoding.EncodeToString(marshal)}, nil)
		//}
		return nil
	}
)

func downloadReleasePackage(t string, localSavePath string) error {
	if err := requestToRemoteServer(fmt.Sprintf("/release/download/r/%s/%s/%s", runtime.GOOS, runtime.GOARCH, t), func(resp *http.Response) error {
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return errors.New(resp.Status)
		}
		f, err := os.OpenFile(localSavePath, os.O_CREATE|os.O_RDWR, 0777)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err = io.Copy(f, resp.Body); err != nil {
			return err
		}

		if _, err = f.Seek(0, 0); err != nil {
			return err
		}

		h, err := coderutils.HashByReader(sha256.New(), f)
		if err != nil {
			return err
		}

		if h.ToHexStr() != resp.Header.Get("_h") {
			return errors.New("HASH不匹配")
		}

		return nil
	}); err != nil {
		_ = os.Remove(localSavePath)
		return err
	}
	return nil
}

func requestToRemoteServer(url string, fn func(resp *http.Response) error) error {

	resp, err := tools.DefaultHttpClient.Get(config.LocalWebServerAddress() + url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return fn(resp)
}
