package common

import (
	"net"
	"strings"
)

func GetAllLocalIpAddr()  []string {
	ips:=make([]string,0)
	if addrs,err:=net.InterfaceAddrs();err==nil{
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil && !strings.Contains(ipnet.IP.String(), "169.254") &&
					!strings.Contains(ipnet.IP.String(),"127.0.0.1"){
					ips = append(ips,ipnet.IP.String())
				}
			}
		}

	}

	return ips
}