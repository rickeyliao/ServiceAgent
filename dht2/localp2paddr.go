package dht2

import (
	"fmt"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/kprc/nbsnetwork/tools/privateip"
	"github.com/pkg/errors"
	"github.com/rickeyliao/ServiceAgent/common"
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
	kaQuit  chan struct{}
	wg      *sync.WaitGroup
	sock    *net.UDPConn
	ncsLock sync.Mutex
	ncs     []*NatClient
	ncQuit  chan struct{}
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
	localP2pAddr.addr.Port = int(common.GetSAConfig().DhtListenPort)
	localP2pAddr.addr.LoadLocalAddr()
	localP2pAddr.addr.NbsAddr = GetLocalNAddr()
	localP2pAddr.rcvQ = make(chan *Block, 256)
	localP2pAddr.rcvQuit = make(chan struct{}, 1)
	localP2pAddr.wrtQ = make(chan *Block, 256)
	localP2pAddr.reWrtQ = make(chan *Block, 256)
	localP2pAddr.wrtQuit = make(chan struct{}, 1)
	localP2pAddr.kaQuit = make(chan struct{}, 1)
	localP2pAddr.ncQuit = make(chan struct{}, 1)
	localP2pAddr.wg = &sync.WaitGroup{}

	return localP2pAddr
}

func (lp *LocalP2pAddr) GetP2pAddr() *P2pAddr {
	return lp.addr
}

func SendAndRcv(ip net.IP, port int, b2s []byte) (resp []byte, err error) {
	raddr := &net.UDPAddr{IP: ip, Port: port}
	conn, err := net.DialUDP("udp4", nil, raddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	deadline := time.Now().Add(time.Second * 5)
	conn.SetDeadline(deadline)
	conn.Write(b2s)
	buf := make([]byte, CtrlMsgBufLen)
	nRead, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:nRead], nil

}


func (lp *LocalP2pAddr) Online(naddr NAddr, bsip net.IP, bsport int) error {
	b2s := BuildOnlineReq().Pack()

	res, err := SendAndRcv(bsip, bsport, b2s)
	if err != nil {
		return err
	}

	cm, offset := UnPackCtrlMsg(res)
	if cm.typ == Msg_BS_Resp {
		rbs := &RespBSMsg{}
		rbs.CtrlMsg = *cm
		rbs.UnpackBS(res[offset:])
		BootsTrapFailed(naddr)
		if len(rbs.BSServer) == 0 {
			return errors.New("No Bootstrap server")
		}
		BootsTrapServerAdd(rbs.BSServer)

		return errors.New("Current bootstrap server can't support online service")
	} else if cm.typ == Msg_Nat_Resp {
		rnm := &RespNatMsg{}
		rnm.CtrlMsg = *cm
		rnm.UnpackNatS(res[offset:])

		//fill local address
		lp.addr.CanService = rnm.CanService
		lp.addr.InternetAddr = rnm.ObservrIP
		//lp.addr.NatAddr = rnm.NatServer
		lp.updateNatServer(rnm.NatServer)

		//prepare to saving normal dht and can service dht
		dn := &DTNode{P2pNode: *(rnm.localAddr.Clone()), lastPingTime: tools.GetNowMsTime()}
		//begin to loop searching
		if rnm.CanService {
			//save to can service dht
			GetCanServiceDht().Insert(dn)
			//can service loop searching
			lp.CanServiceLoop(rnm.localAddr)
		} else {
			GetAllNodeDht().Insert(dn)
			//begin to connect to nat server

		}

		//normal dht loop searching
		lp.NormalLoop(rnm.localAddr)
		return nil
	}

	return nil
}

func (lp *LocalP2pAddr) CanServiceLoop(peer *P2pAddr) error {
	lp.FindCanSrvNode(lp.addr.NbsAddr,peer)
	return nil
}

func (lp *LocalP2pAddr) NormalLoop(peer *P2pAddr) error {
	lp.FindNode(lp.addr.NbsAddr,peer)
	return nil
}

func (lp *LocalP2pAddr) FindNodes(node NAddr) (fnode *P2pAddr,nodes []P2pAddr,err error){
	dhtnode:=&DTNode{}
	dhtnode.P2pNode=P2pAddr{NbsAddr:node}
	dtns:=GetAllNodeDht().FindNearest(dhtnode,DHTNearstCount)

	nls:=&NodeAndLens{}

	for i:=0;i<len(dtns);i++{
		dtn:=dtns[i]
		if dtn.P2pNode.NbsAddr.Cmp(node){
			return &dtn.P2pNode,nil,nil
		}
		l,_:=NbsXorLen(node.Bytes(),dtn.P2pNode.NbsAddr.Bytes())
		nls.AddUniq(l,dtn.P2pNode)
	}



	return
}

func (lp *LocalP2pAddr)FindNodeByA(nls *NodeAndLens,node NAddr) (nodes []P2pAddr,err error) {
	if nls.Left() == 0{
		return nil,errors.New("No Nearst DHT Node")
	}

	nls.SortLH()
	nls.Iterator()

	result := make(chan []P2pAddr,128)
	more   := make(chan int,8)
	cnt := 0

	timer1:=time.NewTimer(5*time.Second)

	for{
		if cnt < DHTFindA && nls.Left()>0{
			nl:=nls.Next()
			cnt ++
			go lp.FindNodeAndInform(node,&nl.Node,result,more)
			continue
		}

		select {
		case r:=<-result:
			return r,nil
		case <-more:
			cnt --
			if cnt == 0{
				return nil,errors.New("Not Found at all")
			}
			continue
		case <-timer1.C:
			return nil,errors.New("Time out")

		}

	}
}


func (lp *LocalP2pAddr) FindCanSrvNodes(node NAddr) (fnode *P2pAddr,nodes []P2pAddr,err error) {
	return
}

func (lp *LocalP2pAddr)FindNode(node NAddr,peer *P2pAddr) (nearstNode []P2pAddr,err error)  {

	req:=BuildReqFindMsg(node)
	buf:=make([]byte,1024)

	n:=req.Pack(buf)

	respbuf,err:=lp.NatSendAndRcv(peer,buf[:n])
	if err!=nil{
		return nil,err
	}

	cm,offset:=UnPackCtrlMsg(respbuf)
	frespm:=&FindRespMsg{}
	frespm.CtrlMsg = *cm

	frespm.UnPackFRespMsg(respbuf[offset:])

	nearstNode = frespm.NearestNodes

	return
}

func (lp *LocalP2pAddr)FindNodeAndInform(node NAddr,peer *P2pAddr, result chan []P2pAddr, more chan int) {
	arr,err:=lp.FindNode(node,peer)
	if err!=nil{
		more <- 1
		return
	}

	result <- arr

	return
}


func (lp *LocalP2pAddr)FindCanSrvNode(node NAddr,peer *P2pAddr) (nearstNode []P2pAddr,err error)  {
	req:=BuildReqFindCanServiceMsg(node)
	buf:=make([]byte,1024)
	n:=req.Pack(buf)

	respbuf,err:=lp.NatSendAndRcv(peer,buf[:n])
	if err!=nil{
		return nil,err
	}

	cm,offset:=UnPackCtrlMsg(respbuf)
	resp:=&FindRespMsg{}
	resp.CtrlMsg = *cm
	resp.UnPackFRespMsg(respbuf[offset:])

	nearstNode = resp.NearestNodes

	return

}

func (lp *LocalP2pAddr)NatSendAndRcv(peer *P2pAddr,b2s []byte) (rcvbuf []byte,err error) {
	sess:=lp.CreateConnSession(peer)
	if sess != nil{
		rcvbuf,err = sess.WriteAndRead(b2s)
		if err == nil{
			sess.Socket.Close()
		}
		return
	}

	return nil,errors.New("Can't Connect to peer")
}

type ConnSession struct {
	PeerIP net.IP
	PeerPort int
	LocalIP net.IP
	LocalPort int
	Socket *net.UDPConn
}

func (cs *ConnSession)WriteAndRead(wrt []byte) (rcv []byte,err error) {
	cs.Socket.Write(wrt)
	deadline := time.Now().Add(time.Second * 5)
	cs.Socket.SetDeadline(deadline)
	for {
		buf1 := make([]byte, CtrlMsgBufLen)
		nRead, err := cs.Socket.Read(buf1)
		if err != nil {
			cs.Socket.Close()
			return nil, err
		}
		cm, _ := UnPackCtrlMsg(buf1[:nRead])
		switch cm.typ {
		case Msg_Nat_Conn_Resp:
			continue
		case Msg_Nat_Sess_Create_Req:
			continue
		case Msg_Nat_Sess_Create_Resp:
			continue
		default:
			return buf1,nil
		}
	}

}

func createSess(ip net.IP,port int,addr NAddr) (sess *ConnSession,err error) {

	laddr:=&net.UDPAddr{}
	raddr:=&net.UDPAddr{IP:ip,Port:port}

	conn,err:=net.DialUDP("udp4",laddr,raddr)
	if err!=nil{
		return nil,err
	}

	req:=BuildNCSessCreateReq()
	buf:=make([]byte,1024)
	offset:=req.Pack(buf)

	deadline := time.Now().Add(time.Second * 3)
	conn.SetDeadline(deadline)
	conn.Write(buf[:offset])
	buf1 := make([]byte, CtrlMsgBufLen)
	nRead, err := conn.Read(buf1)
	if err != nil {
		conn.Close()
		return nil, err
	}

	cm,_:=UnPackCtrlMsg(buf1[:nRead])
	if !cm.localAddr.NbsAddr.Cmp(addr) || cm.typ != Msg_Nat_Sess_Create_Resp{
		conn.Close()
		return nil,err
	}
	sess = &ConnSession{}
	sess.LocalIP = laddr.IP
	sess.LocalPort = laddr.Port
	sess.PeerIP = raddr.IP
	sess.PeerPort = raddr.Port
	sess.Socket = conn
	//conn.SetDeadline(time.Time{})

	return sess,nil
}

func createSessAndInform(ip net.IP,port int, nbsaddr NAddr,result chan *ConnSession) error {
	if sess,err:=createSess(ip,port,nbsaddr);err!=nil{
		return err
	}else{
		result <- sess
		return nil
	}
}
func (lp *LocalP2pAddr)createNatSession(peer *P2pAddr) (*ConnSession,error) {
	laddr:=&net.UDPAddr{}
	raddr:=&net.UDPAddr{IP:peer.InternetAddr,Port:peer.Port}

	conn,err:=net.DialUDP("udp4",laddr,raddr)
	if err!=nil{
		return nil,err
	}
	req:=BuildNCConnReq(peer)
	buf:=make([]byte,2048)
	offset:=req.Pack(buf)

	deadline := time.Now().Add(time.Second * 5)
	conn.SetDeadline(deadline)
	conn.Write(buf[:offset])

	for {
		buf1 := make([]byte, CtrlMsgBufLen)
		nRead, err := conn.Read(buf1)
		if err != nil {
			conn.Close()
			return nil, err
		}
		cm, _ := UnPackCtrlMsg(buf1[:nRead])
		if !cm.localAddr.NbsAddr.Cmp(peer.NbsAddr){
			continue
		}
		switch cm.typ {
		case Msg_Nat_Conn_Resp:
			sessreq:=BuildNCSessCreateReq()
			bufreq:=make([]byte,1024)
			of1:=sessreq.Pack(bufreq)
			conn.Write(bufreq[:of1])
		case Msg_Nat_Sess_Create_Req:
			sessresp:=BuildNCSessCreateResp()
			bufresp:=make([]byte,1024)
			of2:=sessresp.Pack(bufresp)
			conn.Write(bufresp[:of2])
			fallthrough
		case Msg_Nat_Sess_Create_Resp:
			sess := &ConnSession{}
			sess.LocalIP = laddr.IP
			sess.LocalPort = laddr.Port
			sess.PeerIP = raddr.IP
			sess.PeerPort = raddr.Port
			sess.Socket = conn

			return sess,nil
		default:
			conn.Close()
			return nil,errors.New("Get Error msg")
		}
	}
	return nil,errors.New("Unknown error")
}

func (lp *LocalP2pAddr)createNatSessionAndInform(peer *P2pAddr, result chan *ConnSession) error {
	if sess,err:=lp.createNatSession(peer);err!=nil{
		return err
	}else{
		result <- sess
		return nil
	}
}


func (lp *LocalP2pAddr)CreateConnSession(peer *P2pAddr) *ConnSession  {
	if peer.CanService {
		sess,err:=createSess(peer.InternetAddr,peer.Port,peer.NbsAddr)
		if err!=nil{
			return nil
		}
		sess.Socket.SetDeadline(time.Time{})
		return sess
	}
	result := make(chan *ConnSession,128)
	if !peer.CanService && !lp.addr.CanService{
		if peer.InternetAddr.Equal(lp.addr.InternetAddr){
			for i:=0;i<len(peer.InternalAddr);i++ {
				addr := peer.InternalAddr[i]
				go createSessAndInform(addr,peer.Port,peer.NbsAddr,result)
			}
		}
	}
	for i:=0;i<len(peer.NatAddr)&&i<3;i++{
		nat:=&peer.NatAddr[i]
		go lp.createNatSessionAndInform(nat,result)

	}

	timer1:=time.NewTimer(5*time.Second)
	select{
	case sess:=<-result:
		sess.Socket.SetDeadline(time.Time{})
		return sess
	case timer1.C:
		return nil
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

func (lp *LocalP2pAddr) TestCanService(remoteIP net.IP) (bool, error) {

	raddr := &net.UDPAddr{IP: remoteIP, Port: int(common.GetSAConfig().DhtListenPort)}
	conn, err := net.DialUDP("udp4", nil, raddr)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	csm := BuildCanServiceReq()
	b2s := csm.Pack()

	deadline := time.Now().Add(time.Second * 1)
	conn.SetDeadline(deadline)
	conn.Write(b2s)
	buf := make([]byte, CtrlMsgBufLen)
	_, err = conn.Read(buf)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (lp *LocalP2pAddr) TestCanServiceTimes(remoteIP net.IP, times int) (bool, error) {
	if times <= 0 {
		times = 1
	}
	count := 0

	for {
		count++
		if count > times {
			return false, errors.New("peer is not a can service node")
		}
		b, _ := lp.TestCanService(remoteIP)
		if b {
			return true, nil
		}

	}

}

func (lp *LocalP2pAddr) doRcv(block *Block) {
	cm, offset := UnPackCtrlMsg(block.buf)
	if cm.typ == Msg_Online_Req {
		if !lp.addr.CanService || privateip.IsPrivateIP(block.raddr.IP) {
			dn := &DTNode{}
			dn.P2pNode = *(cm.localAddr)
			dn.lastPingTime = tools.GetNowMsTime()

			ds := GetCanServiceDht().FindNearest(dn, NatServerCount)
			var bs []P2pAddr
			for i := 0; i < len(ds); i++ {
				bs = append(bs, ds[i].P2pNode)
			}
			rbs := BuildRespBSMsg(bs)
			//reuse buffer
			offset = rbs.Pack(block.buf)
			blocksnd := &Block{buf: block.buf[:offset], raddr: block.raddr}
			lp.Write(blocksnd)
			return
		}

		go func() {
			//Test raddr is or not is a can server
			b, _ := lp.TestCanServiceTimes(block.raddr.IP, 3)
			var nats []P2pAddr
			if !b {
				//if not a can server,send back 3 nat server,
				dn := &DTNode{}
				dn.P2pNode = *(cm.localAddr)
				dn.lastPingTime = tools.GetNowMsTime()
				ds := GetCanServiceDht().FindNearest(dn, NatServerCount)
				for i := 0; i < len(ds); i++ {
					nats = append(nats, ds[i].P2pNode)
				}
			}

			ret := BuildRespNatMsg(b, block.raddr.IP, nats)
			offset = ret.Pack(block.buf)
			blocksnd := &Block{buf: block.buf[:offset], raddr: block.raddr}
			//Send online success
			lp.Write(blocksnd)
		}()

	} else if cm.typ == Msg_CanSrv_Req {
		csresp := BuildCanServiceResp()
		offset := PackCtrlMsg(&csresp.CtrlMsg, block.buf)
		blocksnd := &Block{buf: block.buf[:offset], raddr: block.raddr}
		lp.Write(blocksnd)
	} else if cm.typ == Msg_Ka_Req {
		//insert ka node
		GetKAStore().Insert(block.raddr.IP, block.raddr.Port,cm.localAddr.NbsAddr)

		kar := BuildRespNCKAMsg(block.raddr.Port)

		offset = kar.Pack(block.buf)

		blocksnd := &Block{buf: block.buf[:offset], raddr: block.raddr}

		lp.Write(blocksnd)
	} else if cm.typ == Msg_Nat_Refresh_Req {
		dn := &DTNode{}
		dn.P2pNode = *(cm.localAddr)
		dn.lastPingTime = tools.GetNowMsTime()

		ds := GetCanServiceDht().FindNearest(dn, NatServerCount)
		var bs []P2pAddr
		for i := 0; i < len(ds); i++ {
			bs = append(bs, ds[i].P2pNode)
		}

		rnrm := BuildRespNatRefreshMsg(bs)
		offset = rnrm.Pack(block.buf)
		block2snd := &Block{buf: block.buf[:offset], raddr: block.raddr}

		lp.Write(block2snd)
	} else if cm.typ == Msg_Nat_Conn_Req {
		connreq:=&NCConnReq{}
		connreq.CtrlMsg = *cm
		connreq.UnPack(block.buf[offset:])
		remote:=connreq.Wait4ConnNode

		nc:=GetKAStore().Find(remote.NbsAddr)
		var (
			ip net.IP
			port int
			errCode int
		)
		if nc == nil{
			ip = net.IPv4zero
			errCode = NCConnNotFound
		}else{
			ip = nc.ip
			port = nc.port
			errCode = NCConnSuccess

			ci:=BuildNCConnInform(block.raddr.IP,block.raddr.Port)
			buf:=make([]byte,1024)
			ofs := ci.Pack(buf)
			go SendAndRcv(nc.ip,nc.port,buf[:ofs])
		}
		resp:=BuildNCConnResp(errCode,ip,port)
		offset = resp.Pack(block.buf)
		blocksnd := &Block{buf:block.buf[:offset],raddr:block.raddr}
		lp.Write(blocksnd)
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

func (lp *LocalP2pAddr) stopAllKa() {
	if len(lp.ncs) <= 0 {
		return
	}

	lp.addr.NatAddr = nil

	for i := 0; i < len(lp.ncs); i++ {
		nc := lp.ncs[i]
		nc.Stop()
	}
	lp.ncs = nil
}

func (lp *LocalP2pAddr) updateNatServer(ns []P2pAddr) {
	lp.ncsLock.Lock()
	defer lp.ncsLock.Unlock()
	if len(ns) == 0 && !lp.addr.CanService {
		panic("No nat server for you !")
		return
	}
	if lp.addr.CanService {
		lp.stopAllKa()
		return
	}

	var nc2del []*NatClient
	var ncNews []*NatClient
	var addr2add []P2pAddr

	for i := 0; i < len(ns); i++ {
		addr := ns[i]

		var j int = 0

		for j = 0; j < len(lp.ncs); j++ {
			nc := lp.ncs[j]
			if nc.addr.NbsAddr.Cmp(addr.NbsAddr) {
				break
			}
		}

		if j == len(lp.ncs) {
			addr2add = append(addr2add, addr)
		}
	}

	for i := 0; i < len(lp.ncs); i++ {
		nc := lp.ncs[i]

		var j int = 0
		for j = 0; j < len(ns); j++ {
			addr := ns[j]
			if nc.addr.NbsAddr.Cmp(addr.NbsAddr) {
				break
			}
		}

		if j == len(ns) {
			nc2del = append(nc2del, nc)
		} else {
			ncNews = append(ncNews, nc)
		}

	}

	lp.ncs = ncNews

	for i := 0; i < len(nc2del); i++ {
		nc := nc2del[i]
		nc.Stop()

	}

	for i := 0; i < len(addr2add); i++ {
		nc := NewNatClient(&addr2add[i])
		go nc.Run()
		lp.ncs = append(lp.ncs, nc)
	}
}

func NbsP2PListen() {
	errListenChan = make(chan error, 3)
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

	lp.wg.Add(1)
	go lp.timeDoNCKA()

}

func (lp *LocalP2pAddr) StopP2pService() {

	if lp.sock != nil {
		lp.sock.Close()
		//lp.sock = nil
	}
	lp.ncQuit <- struct{}{}

	lp.kaQuit <- struct{}{}
	lp.rcvQuit <- struct{}{}
	lp.wrtQuit <- struct{}{}

	lp.wg.Wait()
	lp.sock = nil
}

func (lp *LocalP2pAddr) timeDoNCKA() {
	defer lp.wg.Done()

	for {
		select {
		case <-lp.ncQuit:
			lp.stopAllKa()
			return
		default:

		}

		IteratorNCKA(lp.ncs)

		time.Sleep(time.Second)

		lp.ncs = TimeoutNCKA(lp.ncs)

		lp.CheckNCCount()

	}
}

func (lp *LocalP2pAddr) NatClientCnt() int {
	cnt := 0
	for _, nc := range lp.ncs {
		if nc.IsRunning() {
			cnt++
		}
	}

	return cnt
}

func (lp *LocalP2pAddr) CheckNCCount() {
	if len(lp.ncs) < NatServerCount {
		dn := &DTNode{P2pNode: *lp.addr, lastPingTime: tools.GetNowMsTime()}
		dns := GetCanServiceDht().FindNearest(dn, NatServerCount)
		if len(dns) == 0 {
			return
		}

		for _, d := range dns {
			natreq := BuildNatRefreshReq()
			d2s := natreq.Pack()

			res, err := SendAndRcv(d.P2pNode.InternetAddr, d.P2pNode.Port, d2s)
			if err != nil {
				continue
			}
			cm, offset := UnPackCtrlMsg(res)
			if cm.typ == Msg_Nat_Refresh_Resp {
				nrm := &RespNatRefreshMsg{}
				nrm.CtrlMsg = *cm
				nrm.UnpackNatRefreshS(res[offset:])
				lp.updateNatServer(nrm.NatServer)
			}

			if len(lp.ncs) >= NatServerCount {
				return
			}

		}
	}
}
