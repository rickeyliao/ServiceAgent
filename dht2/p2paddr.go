package dht2

import (
	"fmt"
	"github.com/pkg/errors"
	"net"
	"time"
	"github.com/kprc/nbsnetwork/tools/privateip"
	"github.com/rickeyliao/ServiceAgent/common"
)

var (
	CanServiceTue   = true
	CanServiceFalse = false
)

type P2pAddr struct {
	NbsAddr      NAddr
	CanService   bool
	InternetAddr net.IP
	Port         int
	InternalAddr []net.IP
	NatAddr      []P2pAddr
}

func (pa *P2pAddr) Clone() *P2pAddr {
	pa1 := &P2pAddr{}

	pa1.NbsAddr = pa.NbsAddr

	pa1.CanService = pa.CanService

	pa1.InternalAddr = pa.InternalAddr

	pa1.Port = pa.Port
	pa1.InternalAddr = make([]net.IP, 0)
	pa1.InternalAddr = append(pa1.InternalAddr, pa.InternalAddr...)

	pa1.NatAddr = make([]P2pAddr, 0)

	for _, nata := range pa.NatAddr {
		pa1.NatAddr = append(pa1.NatAddr, *(nata.Clone()))
	}

	return pa1

}

func (pa *P2pAddr) LoadLocalAddr() {
	addrs := GetAllLocalIps()
	var privAddrs []net.IP
	for i:=0;i<len(addrs);i++{
		if privateip.IsPrivateIP(addrs[i]){
			privAddrs = append(privAddrs,addrs[i])
		}
	}

	pa.InternalAddr = privAddrs

	cfg:=common.GetSAConfig()

	if cfg.DhtInternetIp != ""{
		dip := net.ParseIP(cfg.DhtInternetIp)
		pa.InternetAddr = dip
	}

}

func (pa *P2pAddr) sendAndRcv(b2s []byte, mstimeout int) (bfr []byte, err error) {
	raddr := &net.UDPAddr{IP: pa.InternetAddr, Port: pa.Port}
	conn, err := net.DialUDP("udp4", nil, raddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if mstimeout <= 0 {
		mstimeout = 1000 //1 second
	}

	deadline := time.Now().Add(time.Duration(int64(time.Millisecond) * int64(mstimeout)))
	conn.SetDeadline(deadline)
	conn.Write(b2s)
	buf := make([]byte, CtrlMsgBufLen)
	var n int
	n, err = conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func (pa *P2pAddr) SendAndRcv(b2s []byte, times int, mstimeout int) (b2r []byte, err error) {

	if times <= 0 {
		times = 1
	}

	for i := 0; i < times; i++ {
		b2r, err = pa.sendAndRcv(b2s, mstimeout)
		if err != nil {
			continue
		} else {
			return b2r, nil
		}
	}
	return nil, errors.New("Can't connect")
}

func NewP2pAddr() *P2pAddr {

	pa := &P2pAddr{}
	return pa
}

func (pa *P2pAddr) Ping() bool {
	//local -> pa
	if pa.CanService {
		//ping direct
	} else {
		//
	}

	return false
}

func (pa *P2pAddr) String() string {
	var s string
	s += fmt.Sprintf("NbsAddr: %s ", pa.NbsAddr.ID())
	s += fmt.Sprintf("CanService: %t ", pa.CanService)
	s += fmt.Sprintf("InternetAddr: %s ", pa.InternetAddr.To4().String())
	s += fmt.Sprintf("Port: %d ", pa.Port)
	if len(pa.InternetAddr) > 0 {
		var internalips string
		for _, ip := range pa.InternalAddr {
			if internalips != "" {
				internalips += " "
			}
			internalips += ip.To4().String()
		}
		s += fmt.Sprintf("Internal Address: %s ", internalips)
	}

	if len(pa.NatAddr) > 0 {
		var nataddrs string

		for _, pa1 := range pa.NatAddr {
			s += pa1.String()
		}

		s += fmt.Sprintf("Nat addr: { %s } ", nataddrs)
	}

	return s
}
