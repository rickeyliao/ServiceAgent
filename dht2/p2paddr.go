package dht2

import (
	"net"
	"github.com/rickeyliao/ServiceAgent/common"
	"fmt"
)

var (
	CanServiceTue = true
	CanServiceFalse = false
)

type P2pAddr struct {
	NbsAddr NAddr
	CanService bool
	InternetAddr net.IP
	Port int
	InternalAddr []net.IP
	NatAddr []P2pAddr
}

func (pa *P2pAddr)Clone() *P2pAddr  {
	pa1:=&P2pAddr{}

	pa1.NbsAddr = pa.NbsAddr

	pa1.CanService = pa.CanService

	pa1.InternalAddr = pa.InternalAddr

	pa1.Port = pa.Port
	pa1.InternalAddr = make([]net.IP,0)
	pa1.InternalAddr = append(pa1.InternalAddr,pa.InternalAddr...)

	pa1.NatAddr = make([]P2pAddr,0)

	for _,nata:=range pa.NatAddr{
		pa1.NatAddr = append(pa1.NatAddr,*(nata.Clone()))
	}

	return pa1

}


func (pa *P2pAddr)LoadLocalAddr()  {
	pa.InternalAddr = GetAllLocalIps()
}

func NewP2pAddr() *P2pAddr {

	pa:=&P2pAddr{Port:int(common.GetSAConfig().DhtListenPort)}
	return pa
}

func (pa *P2pAddr)Ping() bool  {



	return false
}



func (pa *P2pAddr)String() string  {
	var s string
	s += fmt.Sprintf("NbsAddr: %s ",pa.NbsAddr.ID())
	s += fmt.Sprintf("CanService: %t ",pa.CanService)
	s += fmt.Sprintf("InternetAddr: %s ",pa.InternetAddr.To4().String())
	s += fmt.Sprintf("Port: %d ",pa.Port)
	if len(pa.InternetAddr) > 0{
		var internalips string
		for _,ip:=range pa.InternalAddr{
			if internalips != ""{
				internalips += " "
			}
			internalips += ip.To4().String()
		}
		s += fmt.Sprintf("Internal Address: %s ",internalips)
	}

	if len(pa.NatAddr)>0{
		var nataddrs string

		for _,pa1:=range pa.NatAddr{
			s += pa1.String()
		}

		s += fmt.Sprintf("Nat addr: { %s } ", nataddrs)
	}

	return s
}








