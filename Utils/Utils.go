package Utils

import (
	"net"
	"strings"
)

func GetLocalIP() string {
	effaces, _ := net.Interfaces()
	for _, i := range effaces {
		adders, _ := i.Addrs()
		for _, addr := range adders {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ipAddr := ip.String()
			if strings.HasPrefix(ipAddr, "172.") || strings.HasPrefix(ipAddr, "192.") || strings.HasPrefix(ipAddr, "10.") {
				return ipAddr
			}
		}
	}
	return "127.0.0.1"
}
