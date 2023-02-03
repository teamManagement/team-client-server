package remoteserver

import (
	"fmt"
)

var serverTcpAddress = ""

func GetServerTcpAddress() (string, error) {
	if serverTcpAddress != "" {
		return serverTcpAddress, nil
	}

	if err := RequestWebServiceWithResponse("/config/tcp/address", &serverTcpAddress); err != nil {
		return "", fmt.Errorf("从远程服务器获取通信通道失败: %w", err)
	}

	return serverTcpAddress, nil
}

func init() {
	_, _ = GetServerTcpAddress()
}
