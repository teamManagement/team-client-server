package entrance

import (
	"flag"
	"github.com/go-base-lib/logs"
	"github.com/kardianos/service"
	"team-client-server/config"
	loginit "team-client-server/logs"
	"team-client-server/server"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go server.Run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func Run() {
	cmd := flag.String("c", "no", "要执行的命令")
	configDir := flag.String("configDir", "", "配置存储目录")

	config.LoadConfig(*configDir)

	loginit.InitLog()

	prg := &program{}
	s, err := service.New(prg, generatorServiceConfig())
	if err != nil {
		logs.Panicf("创建服务实例失败: %s", err.Error())
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

func generatorServiceConfig() *service.Config {
	return &service.Config{
		Name:        "teamLocalServer",
		DisplayName: "team manager application location server",
		Description: "团队协作平台应用程序本地服务",
		Arguments:   []string{},
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
