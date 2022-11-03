package entrance

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"github.com/go-base-lib/logs"
	"github.com/kardianos/service"
	"io"
	"os"
	"os/exec"
	"os/user"
	"team-client-server/config"
	loginit "team-client-server/logs"
	"team-client-server/server"
	"team-client-server/tools"
	"team-client-server/updater"
	"time"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	if tools.TelnetHost("127.0.0.1:65528") {
		os.Exit(255)
	}
	go server.Run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func Run() {
	cmd := flag.String("cmd", "no", "要执行的命令")
	configDir := flag.String("configDir", "", "配置存储目录")
	updateInfo := flag.String("updateInfo", "", "更新信息")
	flag.Parse()

	config.LoadConfig(*configDir)

	loginit.InitLog()

	logs.Info("接收到的命令: ", *cmd)
	if *cmd == "updater" {
		logs.Info("进入更新程序...")
		time.Sleep(2 * time.Second)
		updateInfoBytes, err := base64.StdEncoding.DecodeString(*updateInfo)
		if err != nil {
			logs.Infof("更新数据解码失败: %s", err.Error())
			panic(err)
		}

		logs.Infof("更新数据解码成功: %s", updateInfoBytes)

		info := updater.UpdateInfo{}
		if err = json.Unmarshal(updateInfoBytes, &info); err != nil {
			logs.Errorf("更新数据从json转换为object失败: %s", err.Error())
			panic(err)
		}

		logs.Infof("打开 WillUpdateAsarPath 文件: %s", info.WillUpdateAsarPath)
		f, err := os.OpenFile(info.WillUpdateAsarPath, os.O_RDONLY, 0655)
		if err != nil {
			return
		}
		defer f.Close()

		logs.Infof("打开 更新目标 文件: %s", info.Asar)
		destF, err := os.OpenFile(info.Asar, os.O_CREATE|os.O_WRONLY, 0655)
		if err != nil {
			return
		}
		defer destF.Close()

		_, _ = io.Copy(destF, f)
		logs.Info("更新文件拷贝完成")

		_ = os.Remove(info.WillUpdateAsarPath)
		args := make([]string, 0, 3)
		if info.Debug {
			args = append(args, info.WorkDir)
			args = append(args, "__updater_start__")
			args = append(args, "__debug_work_dir__="+info.WorkDir)
		}
		logs.Infof("重新唤醒客户端进程...")
		current, err := user.Current()
		if err != nil {
			logs.Errorf("获取当前用户失败: %s", err.Error())
		}
		logs.Infof("当前用户信息: ", current)
		//command := exec.Command(info.Exec, args...)
		command := exec.Command("C:\\Users\\slx\\AppData\\Local\\Programs\\teamwork\\teamwork.exe")
		command.Dir = info.WorkDir
		output, err := command.CombinedOutput()
		if err != nil {
			logs.Errorf("子进程唤醒发生错误: %s", err.Error())
		}
		logs.Infof("客户端进程唤醒之后的输出: %s", string(output))

		return
	}

	prg := &program{}
	s, err := service.New(prg, generatorServiceConfig(*configDir))
	if err != nil {
		logs.Panicf("创建服务实例失败: %s", err.Error())
	}

	if *cmd == "check" || *cmd == "start" {
		if tools.TelnetHost("127.0.0.1:65528") {
			os.Exit(9)
		}
		if *cmd == "check" {
			os.Exit(0)
			return
		}
	}

	if *cmd == "stop" || *cmd == "install" || *cmd == "start" || *cmd == "uninstall" {
		err = service.Control(s, *cmd)
		if err != nil {
			logs.Panicf("应用命令执行失败: %s", err.Error())
		}
		return
	}

	if err = s.Run(); err != nil {
		logs.Panicf("服务运行失败: %s", err.Error())
	}
}

func generatorServiceConfig(configDir string) *service.Config {
	return &service.Config{
		Name:        "teamLocalServer",
		DisplayName: "team manager application location server",
		Description: "团队协作平台应用程序本地服务",
		Arguments:   []string{"-configDir=" + configDir},
		Option: map[string]interface{}{
			"UserService": true,
			"RunAtLoad":   true,
			"LaunchdConfig": `<?xml version='1.0' encoding='UTF-8'?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN"
"http://www.apple.com/DTDs/PropertyList-1.0.dtd" >
<plist version='1.0'>
 <dict>
   <key>Label</key>
   <string>{{html .Name}}</string>
   <key>ProgramArguments</key>
   <array>
     <string>{{html .Path}}</string>
   {{range .Config.Arguments}}
     <string>{{html .}}</string>
   {{end}}
   </array>
   {{if .UserName}}<key>UserName</key>
   <string>{{html .UserName}}</string>{{end}}
   {{if .ChRoot}}<key>RootDirectory</key>
   <string>{{html .ChRoot}}</string>{{end}}
   {{if .WorkingDirectory}}<key>WorkingDirectory</key>
   <string>{{html .WorkingDirectory}}</string>{{end}}
   <key>SessionCreate</key>
   <{{bool .SessionCreate}}/>
   <key>KeepAlive</key>
   <{{bool .KeepAlive}}/>
   <key>RunAtLoad</key>
   <{{bool .RunAtLoad}}/>
   <key>Disabled</key>
   <false/>
 </dict>
</plist>`,
			"SystemdScript": `[Unit]
Description={{.Description}}
ConditionFileIsExecutable={{.Path|cmdEscape}}
{{range $i, $dep := .Dependencies}}
{{$dep}} {{end}}
[Service]
StartLimitInterval=5
StartLimitBurst=10
ExecStart={{.Path|cmdEscape}}{{range .Arguments}} {{.|cmd}}{{end}}
{{if .ChRoot}}RootDirectory={{.ChRoot|cmd}}{{end}}
{{if .WorkingDirectory}}WorkingDirectory={{.WorkingDirectory|cmdEscape}}{{end}}
{{if .UserName}}User={{.UserName}}{{end}}
{{if .ReloadSignal}}ExecReload=/bin/kill -{{.ReloadSignal}} "$MAINPID"{{end}}
{{if .PIDFile}}PIDFile={{.PIDFile|cmd}}{{end}}
{{if and .LogOutput .HasOutputFileSupport -}}
StandardOutput=file:/var/log/{{.Name}}.out
StandardError=file:/var/log/{{.Name}}.err
{{- end}}
{{if gt .LimitNOFILE -1 }}LimitNOFILE={{.LimitNOFILE}}{{end}}
{{if .Restart}}Restart={{.Restart}}{{end}}
{{if .SuccessExitStatus}}SuccessExitStatus={{.SuccessExitStatus}}{{end}}
RestartSec=120
EnvironmentFile=-/etc/sysconfig/{{.Name}}
[Install]
WantedBy=default.target`,
		},
	}
}
