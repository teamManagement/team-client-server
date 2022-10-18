package errors

import (
	"fmt"
	"os"
)

type ExitCode uint8

const (
	// ExitGetSysCurrentUser 获取系统当前用户失败
	ExitGetSysCurrentUser ExitCode = iota + 1
	// ExitConfigFileCreatEmpty 创建空的配置文件失败
	ExitConfigFileCreatEmpty
	// ExitConfigFileRead 配置文件读取失败
	ExitConfigFileRead
	// ExitConfigFileWriteToEmpty 默认配置写出失败
	ExitConfigFileWriteToEmpty
	// ExitConfigFileParser 配置解析失败
	ExitConfigFileParser
	// ExitDatabaseOpen 数据库打开失败
	ExitDatabaseOpen
	// ExitDatabaseCreateTable 创建数据库表失败
	ExitDatabaseCreateTable
	// ExitLogDirCreate 创建日志目录失败
	ExitLogDirCreate
	// ExitUnixSocketFileCreate 创建Unix通信文件失败
	ExitUnixSocketFileCreate
	// ExitUnixSocketListener Unix通信文件监听失败
	ExitUnixSocketListener
	// ExitServerListenerExit 服务监听退出
	ExitServerListenerExit
)

func (e ExitCode) Println(formatStr string, args ...any) {
	_, _ = os.Stderr.Write([]byte(fmt.Sprintf(formatStr, args...)))
	e.Exit()
}

func (e ExitCode) Exit() {
	os.Exit(int(e))
}
