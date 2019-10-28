package dht2

import (
	"net"
	"github.com/rickeyliao/ServiceAgent/common"
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










