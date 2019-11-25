package dht2

import (
	"fmt"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/pkg/errors"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Block struct {
	buf   []byte
	raddr *net.UDPAddr
}

type LocalP2pAddr struct {
	addr    *P2pAddr
	rcvQ    chan *Block
	rcvQuit chan struct{}
	wrtQ    chan *Block
	reWrtQ  chan *Block
	wrtQuit chan struct{}
	wg      *sync.WaitGroup
	sock    *net.UDPConn
}

var (
	localP2pAddr     *LocalP2pAddr
	localP2pAddrLock sync.Mutex
	quitListen       chan struct{}
	errListenChan    chan error
)

func GetLocalP2pAddr() *LocalP2pAddr {
	if localP2pAddr == nil {
		localP2pAddrLock.Lock()
		defer localP2pAddrLock.Unlock()

		if localP2pAddr == nil {

			localP2pAddr = NewLocalP2pAddr()
		}
	}

	return localP2pAddr
}

func ResetLocalP2p() {
	localP2pAddr = NewLocalP2pAddr()
}

func NewLocalP2pAddr() *LocalP2pAddr {
	localP2pAddr = &LocalP2pAddr{}
	localP2pAddr.addr = NewP2pAddr()
	localP2pAddr.addr.LoadLocalAddr()
	localP2pAddr.addr.NbsAddr = GetLocalNAddr()
	localP2pAddr.rcvQ = make(chan *Block, 256)
	localP2pAddr.rcvQuit = make(chan struct{}, 1)
	localP2pAddr.wrtQ = make(chan *Block, 256)
	localP2pAddr.reWrtQ = make(chan *Block, 256)
	localP2pAddr.wrtQuit = make(chan struct{}, 1)
	localP2pAddr.wg = &sync.WaitGroup{}

	return localP2pAddr
}

func (lp *LocalP2pAddr) GetP2pAddr() *P2pAddr {
	return lp.addr
}

func SendAndRcv(ip string, port int, b2s []byte) (resp []byte, err error) {
	raddr := &net.UDPAddr{IP: net.ParseIP(ip), Port: port}
	conn, err := net.DialUDP("udp4", nil, raddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	deadline := time.Now().Add(time.Second * 2)
	conn.SetDeadline(deadline)
	conn.Write(b2s)
	buf := make([]byte, OnlineBufLen)
	nRead, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:nRead], nil

}

func (lp *LocalP2pAddr) Online(bsip string, bsport int) {
	b2s := BuildOnlineReq().Pack()

	res, err := SendAndRcv(bsip, bsport, b2s)
	if err != nil {
		return
	}

}

func (lp *LocalP2pAddr) ListenOnCanServicePort() {
	defer lp.wg.Done()

	if lp.addr.Port <= 1024 {
		log.Fatal("error P2p Can Service Port must large than 1024")
		return
	}

	log.Println("Can service start at:", strconv.Itoa(lp.addr.Port))

	laddr := &net.UDPAddr{IP: net.IPv4zero, Port: lp.addr.Port}

	usock, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		log.Fatal("Can't listen on ", laddr, " error:", err)
		return
	}

	lp.sock = usock

	//defer usock.Close()

	for {
		buf := make([]byte, 2048)
		n, remoteaddr, err := usock.ReadFromUDP(buf)
		if err != nil {
			return
		}
		fmt.Println(n, remoteaddr.String())
		blk := &Block{buf: buf[:n], raddr: remoteaddr}
		lp.rcvQ <- blk
	}

}

func (lp *LocalP2pAddr) doRcv(block *Block) {
	cm, offset := UnPackCtrlMsg(block.buf)
	if cm.typ == Msg_Online_Req {
		if !lp.addr.CanService {
			dn := &DTNode{}
			dn.P2pNode = *(cm.localAddr.Clone())
			dn.lastPingTime = tools.GetNowMsTime()

			ds, cnt := GetCanServiceDht().FindNearest(dn, NatServerCount)
			//if cnt == 0{
			//
			//}
			//send back Can Server Node addr as bootstrap server

			return
		}
		//Test raddr is or not is a can server
		//Send online success
		//if not a can server,send back 3 nat server,
	} else {

		fmt.Println("offset", offset)
	}

}

func (lp *LocalP2pAddr) Write(block *Block) error {
	select {
	case lp.wrtQ <- block:
	default:
		return errors.New("no enough space to write")
	}

	return nil
}

func (lp *LocalP2pAddr) DoRcv() {
	defer lp.wg.Done()
	for {
		select {
		case rblk := <-lp.rcvQ:
			lp.doRcv(rblk)
		case <-lp.rcvQuit:
			return
		}
	}
}

func (lp *LocalP2pAddr) DoWrt() {
	defer lp.wg.Done()

	cnt := 0
	for {

		select {
		case wblk := <-lp.reWrtQ:
			cnt++
			time.Sleep(time.Millisecond * 200 * time.Duration(cnt))
			_, err := lp.sock.WriteToUDP(wblk.buf, wblk.raddr)
			if err != nil {
				if strings.Contains(err.Error(), "no buffer space available") {
					lp.reWrtQ <- wblk
				} else {
					errListenChan <- err
					return
				}

			}
			continue

		default:

		}

		select {
		case wblk := <-lp.wrtQ:
			_, err := lp.sock.WriteToUDP(wblk.buf, wblk.raddr)
			if err != nil {
				if strings.Contains(err.Error(), "no buffer space available") {
					lp.reWrtQ <- wblk
					continue
				}
				log.Println("send to :", wblk.raddr.String(), "failed, err:", err)
				errListenChan <- err
				return
			} else {
				if cnt > 0 {
					cnt = 0
				}
			}
		case <-lp.wrtQuit:
			return
		}
	}
}

func NbsP2PListen() {
	errListenChan = make(chan error, 1)
	quitListen = make(chan struct{}, 1)

	for {

		for {
			select {
			case <-errListenChan:
			default:
				break
			}
		}

		lp := GetLocalP2pAddr()
		lp.StartP2PService()

		select {
		case err := <-errListenChan:
			log.Println("P2P Listen Receive a err", err)
			lp.StopP2pService()
		case <-quitListen:
			lp.StopP2pService()
			return
		}
		ResetLocalP2p()
	}
}

func StopNbsP2pListen() {
	quitListen <- struct{}{}
}

func (lp *LocalP2pAddr) StartP2PService() {
	lp.wg.Add(1)
	go lp.ListenOnCanServicePort()
	lp.wg.Add(1)
	go lp.DoRcv()
	lp.wg.Add(1)
	go lp.DoWrt()

}

func (lp *LocalP2pAddr) StopP2pService() {

	if lp.sock != nil {
		lp.sock.Close()
		//lp.sock = nil
	}

	lp.rcvQuit <- struct{}{}
	lp.wrtQuit <- struct{}{}

	lp.wg.Wait()
	lp.sock = nil
}
