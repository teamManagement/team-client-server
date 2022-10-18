package tools

import (
	"net"
	"time"
)

func TelnetHost(host string) bool {
	conn, err := net.DialTimeout("tcp", host, 3*time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
