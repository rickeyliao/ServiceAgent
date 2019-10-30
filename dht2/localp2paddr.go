package dht2

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
)

type Block struct {
	buf   []byte
	raddr net.Addr
}

type LocalP2pAddr struct {
	addr *P2pAddr
	rcvQ chan *Block
	wrtQ chan *Block
	wg   *sync.WaitGroup
}

var (
	localP2pAddr     *LocalP2pAddr
	localP2pAddrLock sync.Mutex
)

func GetLocalP2pAddr() *LocalP2pAddr {
	if localP2pAddr == nil {
		localP2pAddrLock.Lock()
		defer localP2pAddrLock.Unlock()

		if localP2pAddr == nil {
			localP2pAddr = &LocalP2pAddr{}
			localP2pAddr.addr = NewP2pAddr()
			localP2pAddr.addr.LoadLocalAddr()
			localP2pAddr.rcvQ = make(chan *Block, 2560)
			localP2pAddr.wrtQ = make(chan *Block, 2560)
			localP2pAddr.wg = &sync.WaitGroup{}
		}
	}

	return localP2pAddr
}

func (lp *LocalP2pAddr) Online() {

}

func (lp *LocalP2pAddr) ListenOnCanServicePort() {

	if lp.addr.Port <= 1024 {
		log.Fatal("error P2p Can Service Port must large than 1024")
		return
	}

	laddr := ":" + strconv.Itoa(lp.addr.Port)

	usock, err := net.ListenPacket("udp4", laddr)
	if err != nil {
		log.Fatal("Can't listen on ", laddr, " error:", err)
		return
	}

	defer usock.Close()

	doneChan := make(chan error, 1)
	buf := make([]byte, 20480)
	go func() {
		for {
			n, remoteaddr, err := usock.ReadFrom(buf)
			if err != nil {
				doneChan <- err
				return
			}
			fmt.Println(n, remoteaddr.String())
			//deadline:=time.Now().Add(time.Second*30)
			//err=usock.SetDeadline(deadline)
			//if err!=nil{
			//	doneChan <-err
			//	return
			//}

			//nw,err:=usock.WriteTo(buf[:n],remoteaddr)
			//if err!=nil{
			//	doneChan <- err
			//	return
			//}
			//
			//fmt.Println(nw)
		}
	}()

	select {
	case <-doneChan:
		return
	}

}

func (lp *LocalP2pAddr) StartP2PService() {

}

func (lp *LocalP2pAddr) StopP2pService() {

}
