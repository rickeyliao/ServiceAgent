package dht2

import (
	"bytes"
	"errors"
	"github.com/kprc/nbsnetwork/tools/privateip"
	"github.com/rickeyliao/ServiceAgent/common"
	"log"
	"net"
	"strings"
	"sync"
)

type BootsTrapServer struct {
	BSServer P2pAddr
	IsFailed bool
}

var (
	bsservers     []BootsTrapServer
	bsserversLock sync.Mutex
)

func bootsTrapAdd(bs string) error {
	cfg := common.GetSAConfig()
	bsarr := strings.Split(bs, "@")
	if len(bsarr) != 2 {
		return errors.New("bootstrap server string error")
	}
	nbsid := NID(bsarr[0])
	addr, err := nbsid.Addr()
	if err != nil {
		return err
	}

	ip := net.ParseIP(bsarr[1])
	p2paddr := &P2pAddr{}
	p2paddr.Port = int(cfg.DhtListenPort)
	p2paddr.NbsAddr = addr
	if privateip.IsPrivateIP(ip) {
		p2paddr.InternalAddr = append(p2paddr.InternalAddr, ip)
	} else {
		p2paddr.InternetAddr = ip
	}

	bts := BootsTrapServer{BSServer: *p2paddr}
	bsservers = append(bsservers, bts)

	return nil
}

func InitBS() {
	cfg := common.GetSAConfig()
	if cfg == nil{
		return
	}

	for _, bs := range cfg.BootstrapIPAddress {
		err := bootsTrapAdd(bs)
		if err != nil {
			log.Println(err.Error())
		}
	}

}

func BootsTrapServerAddFrmCmd(bs string) error {
	bsserversLock.Lock()
	defer bsserversLock.Unlock()

	return bootsTrapAdd(bs)
}

func BootsTrapServerAdd(bss []P2pAddr) {

	bsserversLock.Lock()
	defer bsserversLock.Unlock()

	for i := 0; i < len(bss); i++ {
		bts := BootsTrapServer{BSServer: bss[i]}
		bsservers = append(bsservers, bts)
	}
}

func BootsTrapGetNxtBS() (naddr NAddr, ip net.IP, port int, err error) {
	bsserversLock.Lock()
	defer bsserversLock.Unlock()

	for i := 0; i < len(bsservers); i++ {
		bts := bsservers[i]
		if bts.IsFailed {
			continue
		}
		if bts.BSServer.CanService {
			return bts.BSServer.NbsAddr, bts.BSServer.InternetAddr, bts.BSServer.Port, nil
		}
	}

	for i := 0; i < len(bsservers); i++ {
		bts := bsservers[i]
		if bts.IsFailed {
			continue
		}
		var rip net.IP
		if bytes.Compare(bts.BSServer.InternetAddr.To4(), []byte{0, 0, 0, 0}) == 0 {
			if len(bts.BSServer.InternalAddr) == 0 {
				bts.IsFailed = true
				continue
			}
			rip = bts.BSServer.InternalAddr[0]
		} else {
			rip = bts.BSServer.InternetAddr
		}
		return bts.BSServer.NbsAddr, rip, bts.BSServer.Port, nil
	}

	return NAddr{}, nil, 0, errors.New("No BootsTrap for use")
}

func BootsTrapFailed(addr NAddr) {
	bsserversLock.Lock()
	defer bsserversLock.Unlock()

	for i := 0; i < len(bsservers); i++ {
		s := bsservers[i].BSServer
		if s.NbsAddr.Cmp(addr) {
			bsservers[i].IsFailed = true
		}
	}

}
