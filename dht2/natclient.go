package dht2

import (
	"fmt"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/pkg/errors"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	MsgRCVBufferLen int = 256
	MsgWrtBufferLen int = 256

	NC_UnInitialed int = 0
	NC_IsRunning   int = 1
	NC_IsStoped    int = 2
	NC_IsStopping  int = 3
	NC_IsQuiting   int = 4
	NC_InternalErr int = 5
)

type NatClient struct {
	addr         *P2pAddr
	listenPort   int
	listenSock   *net.UDPConn
	internetPort int
	rcvQ         chan *Block
	rcvQuit      chan struct{}
	wrtQ         chan *Block
	wrtQuit      chan struct{}
	errC         chan error
	quit         chan struct{}
	lastRcvTime  int64
	running      int
	runLock      sync.Mutex
	wg           sync.WaitGroup
}

func NewNatClient(addr *P2pAddr) *NatClient {
	nc := &NatClient{}

	nc.addr = addr
	nc.rcvQ = make(chan *Block, MsgRCVBufferLen)
	nc.rcvQuit = make(chan struct{}, 1)
	nc.wrtQ = make(chan *Block, MsgWrtBufferLen)
	nc.wrtQuit = make(chan struct{}, 1)
	nc.errC = make(chan error, 3)
	nc.quit = make(chan struct{}, 1)
	nc.running = NC_UnInitialed

	return nc
}

func (nc *NatClient) ncinit() (err error) {
	laddr := net.UDPAddr{IP: net.ParseIP("0.0.0.0")}
	if nc.listenSock, err = net.ListenUDP("udp4", &laddr); err != nil {
		return
	}
	nc.listenPort = laddr.Port
	fmt.Println("Start Ka at local port:", laddr.Port)

	return
}

func (nc *NatClient) start() (err error) {
	defer nc.wg.Done()

	for {
		buf := make([]byte, 2048)
		var (
			n          int
			remoteaddr *net.UDPAddr
		)
		n, remoteaddr, err = nc.listenSock.ReadFromUDP(buf)
		if err != nil {
			nc.errC <- err
			return
		}

		fmt.Println(n, remoteaddr.String())
		blk := &Block{buf: buf[:n], raddr: remoteaddr}
		nc.lastRcvTime = tools.GetNowMsTime()

		nc.rcvQ <- blk
	}

}

func (nc *NatClient) Run() (err error, satus int) {
	if nc.running == NC_IsRunning || nc.running == NC_IsStopping || nc.running == NC_IsQuiting {
		return errors.New("nc: " + nc.addr.NbsAddr.String() + " is running"), NC_IsRunning
	}

	nc.runLock.Lock()
	if nc.running == NC_IsRunning || nc.running == NC_IsStopping || nc.running == NC_IsQuiting {
		return errors.New("nc: " + nc.addr.NbsAddr.String() + " is running"), NC_IsRunning
	}
	nc.running = NC_IsRunning
	nc.runLock.Unlock()

	if err = nc.ncinit(); err != nil {
		return
	}
	nc.wg.Add(1)
	go nc.start()
	nc.wg.Add(1)
	go nc.doRcv()
	nc.wg.Add(1)
	go nc.doWrt()

	select {
	case err = <-nc.errC:
		return err, NC_InternalErr
	case <-nc.quit:
		return nil, NC_IsQuiting
	}

}

func (nc *NatClient) IsRunning() bool {
	nc.runLock.Lock()
	defer nc.runLock.Unlock()
	if nc.running == NC_IsRunning {
		return true
	}
	return false
}

func (nc *NatClient) Stop() {

	nc.runLock.Lock()
	nc.running = NC_IsStopping
	nc.runLock.Unlock()

	nc.rcvQuit <- struct{}{}
	nc.wrtQuit <- struct{}{}
	nc.quit <- struct{}{}

	if nc.listenSock != nil {
		nc.listenSock.Close()
		nc.listenSock = nil
	}

	nc.wg.Wait()
	nc.runLock.Lock()
	nc.running = NC_IsStoped
	nc.runLock.Unlock()
}

func (nc *NatClient) doRcv() {
	defer nc.wg.Done()
	for {
		select {
		case blk := <-nc.rcvQ:
			nc.rcv(blk)
		case <-nc.rcvQuit:
			return
		}
	}
}

func (nc *NatClient) WrtKa() error {
	ka := BuildNCKAReq()

	buf := ka.Pack()

	return nc.Wrt(buf)
}

func (nc *NatClient) Wrt(buf []byte) error {
	blk := &Block{buf: buf}
	raddr := &net.UDPAddr{IP: nc.addr.InternetAddr, Port: nc.addr.Port}
	blk.raddr = raddr

	select {
	case nc.wrtQ <- blk:
		return nil
	default:
		return errors.New("Buffer overflow")
	}
}

func (nc *NatClient) rcv(blk *Block) {
	cm, offset := UnPackCtrlMsg(blk.buf)

	switch cm.typ {
	case Msg_Ka_Resp:
		kar := &NCKAResp{}
		kar.CtrlMsg = *cm
		offset = kar.UnPackNCKA(blk.buf[offset:])
		nc.internetPort = kar.RPort
		return
	case Msg_Nat_Conn_Inform:
		inform:=&NCConnInForm{}
		inform.CtrlMsg = *cm
		offset += inform.UnPack(blk.buf[offset:])
		reply:=BuildNCConnReply()
		buf:=make([]byte,1024)
		offset = reply.Pack(buf)
		nc.Wrt(buf[:offset])

		sess:=BuildNCSessCreateReq()
		buf = make([]byte,1024)
		offset = sess.Pack(buf)
		go SendAndRcv(inform.Wait4ConnIP,inform.Wait4ConnPort,buf[:offset])
		return
	case Msg_Nat_Sess_Create_Resp:
		//nothing to do ...
		return
	case Msg_Nat_Sess_Create_Req:
		sess:=&NCSessionCreateReq{}
		sess.CtrlMsg = *cm
		sessResp:=BuildNCSessCreateResp()
		buf := make([]byte,1024)
		offset = sessResp.Pack(buf)
		go SendAndRcv(blk.raddr.IP,blk.raddr.Port,buf[:offset])
		return
	case Msg_Dht_Find:
		req:=&FindReqMsg{}
		req.CtrlMsg = *cm
		offset+=req.UnPackFRM(blk.buf[offset:])

		dhtnode:=&DTNode{P2pNode:P2pAddr{NbsAddr:req.NodeToFind}}
		dtns:=GetAllNodeDht().FindNearest(dhtnode,DHTNearstCount)
		dn := &DTNode{P2pNode: *(cm.localAddr.Clone()), lastPingTime: tools.GetNowMsTime()}
		GetAllNodeDht().Insert(dn)
		resp:=BuildRespFindMsg(req.NodeToFind,DTNS2Addrs(dtns))
		buf := make([]byte,1024)
		offset = resp.Pack(buf)
		go SendAndRcv(blk.raddr.IP,blk.raddr.Port,buf[:offset])
		return
	}
}

func (nc *NatClient) doWrt() {
	defer nc.wg.Done()
	for {
		select {
		case wblk := <-nc.wrtQ:
			_, err := nc.listenSock.WriteToUDP(wblk.buf, wblk.raddr)
			if err != nil {
				if strings.Contains(err.Error(), "no buffer space available") {
					time.Sleep(time.Millisecond * 200)
					continue
				}
				nc.errC <- err
				return
			}
		case <-nc.wrtQuit:
			return
		}
	}
}

func IteratorNCKA(ncs []*NatClient) {

	if len(ncs) == 0 {
		return
	}

	now := tools.GetNowMsTime()

	for i := 0; i < len(ncs); i++ {
		nc := ncs[i]

		if !nc.IsRunning() {
			continue
		}

		if now-nc.lastRcvTime > 3000 {
			nc.WrtKa()
		}

	}
}

func TimeoutNCKA(ncs []*NatClient) (newncs []*NatClient) {
	if len(ncs) == 0 {
		return
	}

	now := tools.GetNowMsTime()

	for i := 0; i < len(ncs); i++ {
		nc := ncs[i]
		if !nc.IsRunning() {
			continue
		}

		if now-nc.lastRcvTime > 30000 {
			nc.Stop()
			continue
		}

		newncs = append(newncs, nc)

	}

	return
}
