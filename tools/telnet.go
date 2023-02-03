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

func TelnetHostRangeNetInterfaces(host string) (net.IP, bool) {
	inters, err := net.Interfaces()
	if err != nil {
		return nil, false
	}

	raddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return nil, false
	}

	var ip net.IP

	for i := range inters {
		inter := inters[i]
		addrs, err := inter.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			switch t := addr.(type) {
			case *net.IPNet:
				ip = t.IP
			case *net.IPAddr:
				ip = t.IP
			default:
				continue
			}

			//if ip.IsLoopback() {
			//	continue
			//}

			laddr, err := net.ResolveTCPAddr("tcp", ip.String()+":0")
			if err != nil {
				continue
			}

			conn, err := net.DialTCP("tcp", laddr, raddr)
			if err != nil {
				continue
			}
			_ = conn.Close()
			return ip, true
		}
	}
	return nil, false
}
