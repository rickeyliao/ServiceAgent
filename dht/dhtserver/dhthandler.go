package dhtserver

import (
	"github.com/pkg/errors"
	"github.com/rickeyliao/ServiceAgent/dht/pb"
	"log"
	"net"
	"sync"
)

type DhtFunc func(dhtmessage pbdht.Dhtmessage, addr *net.UDPAddr, conn *net.UDPConn) (err error)

type DhtHandler struct {
	h map[uint32]DhtFunc
}

var (
	dhthandlerInst     *DhtHandler
	dhthandlerInstLock sync.Mutex
)

func GetDhtHandlerInst() *DhtHandler {
	if dhthandlerInst != nil {
		return dhthandlerInst
	}

	dhthandlerInstLock.Lock()
	defer dhthandlerInstLock.Unlock()

	if dhthandlerInst != nil {
		return dhthandlerInst
	}

	dhthandlerInst = &DhtHandler{}

	dhthandlerInst.h = make(map[uint32]DhtFunc)

	return dhthandlerInst
}

func (dhr *DhtHandler) Reg(msgtyp uint32, f DhtFunc) {
	if _, ok := dhr.h[msgtyp]; !ok {
		dhr.h[msgtyp] = f
	} else {
		log.Println("Reg msg type:", msgtyp, "Duplicated")
	}
}

func (dhr *DhtHandler) Dispatch(dm pbdht.Dhtmessage, addr *net.UDPAddr, conn *net.UDPConn) (err error) {
	if f, ok := dhr.h[dm.Msgtyp]; ok {
		if f != nil {
			return f(dm, addr, conn)
		}
	} else {
		return errors.Errorf("Message %d not found", dm.Msgtyp)
	}

	return
}
