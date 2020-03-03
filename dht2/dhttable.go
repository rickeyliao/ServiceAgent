package dht2

import (
	"fmt"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
	"sort"
	"github.com/pkg/errors"
	"time"
)

//max node in bukcet
var MaxKBucket int = 16

//max node in backup bucket
var MaxKBucketBackup int = 16

type DTNode struct {
	P2pNode      P2pAddr
	lastPingTime int64
	RefCnt       int
	Next         *DTNode
}

type PingNode struct {
	Wait2Ping   *DTNode
	Wait2Insert *DTNode
	Dht         *DhtTable
}

type DTBucket struct {
	RootLock   sync.Mutex
	RootCnt    int
	Root       *DTNode
	BackupLock sync.Mutex
	BackupCnt  int
	Backup     *DTNode
	Dht        *DhtTable
}

type DhtTable struct {
	HashTable      [257]DTBucket
	DTLock         sync.Mutex
	PingNodeChan   chan PingNode
	TimeOutChan    chan PingNode
	Wg             *sync.WaitGroup
	PingQuit       chan int
	TimeQuitCreate chan int
	TimeQuit       chan int
}

type NodeAndLen struct {
	Len int
	Node P2pAddr
}

type NodeAndLens struct {
	Nls []*NodeAndLen
	Indicator int
}

func (nal *NodeAndLens)Add(l int,node P2pAddr)  {
	nl:=&NodeAndLen{Len:l,Node:node}
	nal.Nls = append(nal.Nls,nl)
}

func (nal *NodeAndLens)AddUniq(l int,node P2pAddr) error {
	for i:=0;i<len(nal.Nls);i++{
		if node.NbsAddr.Cmp(nal.Nls[i].Node.NbsAddr){
			return errors.New("dup")
		}
	}

	nal.Add(l,node)

	return nil
}

func (nal *NodeAndLens)Concat(nal2 *NodeAndLens)  {

	for i:=0;i<len(nal2.Nls);i++{
		nal.AddUniq(nal2.Nls[i].Len,nal2.Nls[i].Node)
	}
}

func (nal *NodeAndLens)Count() int {
	return len(nal.Nls)
}

func (nal *NodeAndLens)SortLH() {

	sort.Slice(nal.Nls, func(i, j int) bool {
		if nal.Nls[i].Len <= nal.Nls[j].Len{
			return true
		}else {
			return false
		}
	})
}

func (nal *NodeAndLens)SortHL()  {
	sort.Slice(nal.Nls, func(i, j int) bool {
		if nal.Nls[i].Len >= nal.Nls[j].Len{
			return true
		}else {
			return false
		}
	})
}

func (nal *NodeAndLens)Iterator()  {
	nal.Indicator = 0
}

func (nal *NodeAndLens)Next() (nl *NodeAndLen) {
	if nal.Indicator >= len(nal.Nls){
		return nil
	}
	nl = nal.Nls[nal.Indicator]
	nal.Indicator ++
	return
}

func (nal *NodeAndLens)Left() int  {
	return len(nal.Nls) - nal.Indicator
}

func DTNS2Addrs(dtns []*DTNode) (arr []P2pAddr){
	for i:=0;i<len(dtns);i++{
		arr = append(arr,dtns[i].P2pNode)
	}

	return
}

func minV(a,b int) int {
	if a > b{
		return b
	}else{
		return a
	}
}

func (nal *NodeAndLens)Equals(nal2 *NodeAndLens,cnt int) bool {
	sort.Slice(nal.Nls, func(i, j int) bool {
		if nal.Nls[i].Len <= nal.Nls[j].Len{
			return true
		}
		return false
	})
	sort.Slice(nal2.Nls, func(i, j int) bool {
		if nal.Nls[i].Len <= nal.Nls[j].Len{
			return true
		}
		return false
	})
	min:=cnt
	min = minV(min,len(nal.Nls))
	min = minV(min,len(nal2.Nls))
	for i:=0;i<min;i++{
		if !nal.Nls[i].Node.NbsAddr.Cmp(nal2.Nls[i].Node.NbsAddr){
			return false
		}
	}

	if cnt > min{
		return false
	}

	return true
}

func NewDhtTable() *DhtTable {
	dt := &DhtTable{}

	dt.DTLock.Lock()
	defer dt.DTLock.Unlock()

	dt.PingNodeChan = make(chan PingNode, 2560)
	dt.TimeOutChan = make(chan PingNode, 2560)

	for i := 0; i < len(dt.HashTable); i++ {
		dt.HashTable[i].Dht = dt
	}

	dt.PingQuit = make(chan int, 0)
	dt.TimeQuit = make(chan int, 0)
	dt.TimeQuitCreate = make(chan int, 0)

	dt.Wg = &sync.WaitGroup{}

	return dt
}

func (dtn *DTNode) String() string {
	s := (&(dtn.P2pNode)).String()

	s += fmt.Sprintf("lastPingTime %d ", dtn.lastPingTime)
	s += fmt.Sprintf("refCnt %d ", dtn.RefCnt)

	return s
}

func (dtn *DTNode) Clone() *DTNode {
	dtn1 := &DTNode{}

	dtn1.P2pNode = *((&(dtn.P2pNode)).Clone())

	dtn1.lastPingTime = dtn.lastPingTime
	dtn1.RefCnt = dtn.RefCnt
	dtn1.Next = nil

	return dtn1
}

func (dt *DhtTable)String() string  {
	s:=""
	for i:=0;i<len(dt.HashTable);i++{
		dtb:=&dt.HashTable[i]
		dtb.RootLock.Lock()
		root:=dtb.Root
		nxt:=root
		for{
			if  nxt == nil{
				break
			}
			s += "main:   "+nxt.String() + "\r\n"
			nxt = nxt.Next
		}

		dtb.RootLock.Unlock()
		dtb.BackupLock.Lock()
		root=dtb.Backup
		nxt=root
		for{
			if  nxt == nil{
				break
			}
			s += "backup: "+nxt.String() + "\r\n"
			nxt = nxt.Next
		}
		dtb.BackupLock.Unlock()
	}

	return s
}

func (dt *DhtTable) Find(node *DTNode) *DTNode {
	laddr := GetLocalNAddr()
	dtaddr := node.P2pNode.NbsAddr

	bucketidx, _ := NbsXorLen(laddr.Bytes(), dtaddr.Bytes())

	bucket := &dt.HashTable[bucketidx]

	bucket.RootLock.Lock()
	defer bucket.RootLock.Unlock()

	return bucket.Find(node).Clone()
}

func (dtb *DTBucket) Find(node *DTNode) *DTNode {
	root := dtb.Root

	if root == nil {
		return nil
	}

	nxt := root

	for {

		if nxt == nil {
			break
		}

		if nxt.P2pNode.NbsAddr.Cmp(node.P2pNode.NbsAddr) {
			break
		}
		nxt = nxt.Next

	}

	return nxt
}

func (dtb *DTBucket) CloneAllNotes() []*DTNode {
	dtnodes := make([]*DTNode, 0)
	nxt := dtb.Root

	for {
		if nxt == nil {
			break
		}

		cpy := nxt.Clone()
		dtnodes = append(dtnodes, cpy)

		nxt = nxt.Next
	}

	return dtnodes
}

func (dt *DhtTable) FindNearest(node *DTNode, cnt int) []*DTNode {
	laddr := GetLocalNAddr()
	dtaddr := node.P2pNode.NbsAddr

	bucketidx, _ := NbsXorLen(laddr.Bytes(), dtaddr.Bytes())

	startbucketidx := bucketidx

	curcnt := 0

	dtnodes := make([]*DTNode, 0)

	for i := startbucketidx; i < len(dt.HashTable); i++ {
		bucket := &dt.HashTable[i]
		bucket.RootLock.Lock()

		if bucket.RootCnt > 0 {
			dtnodes = append(dtnodes, bucket.CloneAllNotes()...)
			curcnt += bucket.RootCnt

			if curcnt >= cnt {
				bucket.RootLock.Unlock()
				return dtnodes
			}
		}

		bucket.RootLock.Unlock()
	}

	if startbucketidx > 0 {
		for i := startbucketidx - 1; i >= 0; i-- {
			bucket := &dt.HashTable[i]
			bucket.RootLock.Lock()

			if bucket.RootCnt > 0 {
				dtnodes = append(dtnodes, bucket.CloneAllNotes()...)
				curcnt += bucket.RootCnt

				if curcnt >= cnt {
					bucket.RootLock.Unlock()
					return dtnodes
				}
			}

			bucket.RootLock.Unlock()
		}
	}

	return dtnodes
}

func (dtb *DTBucket) GetLast() *DTNode {
	prev := dtb.Root
	nxt := prev
	for {
		if nxt == nil {
			return prev
		}
		nxt = nxt.Next
	}
}

func (dtb *DTBucket) GetLastBackup() *DTNode {
	prev := dtb.Backup
	nxt := prev
	for {
		if nxt == nil {
			return prev
		}
		nxt = nxt.Next
	}
}

func (dtb *DTBucket) GetFirstBackup() *DTNode {

	root := dtb.Root

	if root != nil {
		dtb.Root = root.Next
	}

	return root

}

func (dtb *DTBucket) Add(node *DTNode) {
	nxt := dtb.Root
	dtb.Root = node
	node.Next = nxt
	dtb.RootCnt++
}

func (dtb *DTBucket) AddBackup(node *DTNode) {
	nxt := dtb.Backup
	dtb.Backup = node
	node.Next = nxt
	dtb.BackupCnt++
}

func (dtb *DTBucket) Insert(node *DTNode) {
	//node.lastPingTime = tools.GetNowMsTime()

	n1 := dtb.Find(node)
	if n1 == nil {
		if dtb.RootCnt < MaxKBucket {
			dtb.Add(node)
			return
		}

	} else {
		dtb.Remove(n1)
		dtb.Add(node)
		return
	}

	pingnode := PingNode{}
	pingnode.Wait2Ping = dtb.GetLast().Clone()
	pingnode.Wait2Insert = node

	pingnode.Dht = dtb.Dht

	dtb.Dht.PingNodeChan <- pingnode

}

func (dtb *DTBucket) FindBackup(node *DTNode) *DTNode {
	root := dtb.Backup

	if root == nil {
		return nil
	}

	nxt := root

	for {

		if nxt == nil {
			break
		}

		if nxt.P2pNode.NbsAddr.Cmp(node.P2pNode.NbsAddr) {
			break
		}
		nxt = nxt.Next

	}

	return nxt
}

func (dtb *DTBucket) InsertBackup(node *DTNode) {
	n1 := dtb.FindBackup(node)
	if n1 == nil {
		if dtb.BackupCnt < MaxKBucketBackup {
			dtb.AddBackup(node)
			return
		}
	} else {
		dtb.RemoveBackup(n1)
		dtb.AddBackup(node)
		return
	}

	n2 := dtb.GetLastBackup()
	dtb.RemoveBackup(n2)
	dtb.AddBackup(node)
}

func (dt *DhtTable) Insert(node *DTNode) {

	fmt.Println("insert:",node.String())

	laddr := GetLocalNAddr()
	dtaddr := node.P2pNode.NbsAddr

	bucketidx, _ := NbsXorLen(laddr.Bytes(), dtaddr.Bytes())

	fmt.Println("bucketidx:",bucketidx)

	bucket := &dt.HashTable[bucketidx]

	bucket.RootLock.Lock()
	defer bucket.RootLock.Unlock()
	bucket.Insert(node)

}

func (dt *DhtTable) Update(pingNode PingNode) {
	laddr := GetLocalNAddr()
	dtaddr := pingNode.Wait2Ping.P2pNode.NbsAddr

	bucketidx, _ := NbsXorLen(laddr.Bytes(), dtaddr.Bytes())

	bucket := &dt.HashTable[bucketidx]

	bucket.RootLock.Lock()
	defer bucket.RootLock.Unlock()

	bucket.Remove(pingNode.Wait2Ping)
	if bucket.RootCnt < MaxKBucket {
		bucket.Add(pingNode.Wait2Insert)
	}

}

func (dt *DhtTable) UpdateBackup(pingNode PingNode) {
	laddr := GetLocalNAddr()
	dtaddr := pingNode.Wait2Ping.P2pNode.NbsAddr

	bucketidx, _ := NbsXorLen(laddr.Bytes(), dtaddr.Bytes())

	bucket := &dt.HashTable[bucketidx]

	bucket.RootLock.Lock()
	bucket.Remove(pingNode.Wait2Ping)
	bucket.RootLock.Unlock()

	bucket.BackupLock.Lock()
	bucket.InsertBackup(pingNode.Wait2Insert)
	bucket.BackupLock.Unlock()

}

func (dtb *DTBucket) Remove(node *DTNode) {
	nxt := dtb.Root
	prev := nxt
	for {
		if nxt == nil {
			return
		}

		if nxt.P2pNode.NbsAddr.Cmp(node.P2pNode.NbsAddr) {
			if nxt == prev {
				dtb.Root = nxt.Next
			} else {
				prev.Next = nxt.Next
			}
			nxt.Next = nil //for quick recycle
			dtb.RootCnt--
			return
		}
		prev = nxt
		nxt = nxt.Next
	}
}

func (dtb *DTBucket) RemoveBackup(node *DTNode) {

	if node == nil {
		return
	}

	nxt := dtb.Backup
	prev := nxt
	for {
		if nxt == nil {
			return
		}

		if nxt.P2pNode.NbsAddr.Cmp(node.P2pNode.NbsAddr) {
			if nxt == prev {
				dtb.Root = nxt.Next
			} else {
				prev.Next = nxt.Next
			}
			dtb.BackupCnt--
			nxt.Next = nil //for quick recycle
			return
		}
		prev = nxt
		nxt = nxt.Next
	}
}

func (dt *DhtTable) Remove(node *DTNode) {
	laddr := GetLocalNAddr()
	dtaddr := node.P2pNode.NbsAddr

	bucketidx, _ := NbsXorLen(laddr.Bytes(), dtaddr.Bytes())

	bucket := &dt.HashTable[bucketidx]

	bucket.RootLock.Lock()
	defer bucket.RootLock.Unlock()

	bucket.Remove(node)
}

func (dt *DhtTable) RemoveBackup(node *DTNode) {
	laddr := GetLocalNAddr()
	dtaddr := node.P2pNode.NbsAddr

	bucketidx, _ := NbsXorLen(laddr.Bytes(), dtaddr.Bytes())

	bucket := &dt.HashTable[bucketidx]

	bucket.BackupLock.Lock()
	defer bucket.BackupLock.Unlock()

	bucket.RemoveBackup(node)
}

func (dt *DhtTable)doPingTimes(pingNode PingNode,times int)  {
	flag := false
	for i:=0;i<times;i++{
		if pingNode.Wait2Ping.P2pNode.Ping(){
			flag = true
			break
		}else{
			time.Sleep(time.Second)
		}
	}

	if flag{
		dt.Insert(pingNode.Wait2Ping)
	}else {
		dt.Update(pingNode)
	}
}

func (dt *DhtTable) DoPing(pingNode PingNode) {
	go dt.doPingTimes(pingNode,3)
}

func (dt *DhtTable) DoTimeOut(node PingNode) {
	if node.Wait2Ping.P2pNode.Ping() {
		dt.Insert(node.Wait2Ping)
	} else {
		dt.TimeOutUpdate(node.Wait2Ping)
	}
}

func (dt *DhtTable) TimeOutUpdate(node *DTNode) {
	laddr := GetLocalNAddr()
	dtaddr := node.P2pNode.NbsAddr

	bucketidx, _ := NbsXorLen(laddr.Bytes(), dtaddr.Bytes())

	bucket := &dt.HashTable[bucketidx]

	bucket.RootLock.Lock()
	bucket.Remove(node)
	bucket.RootLock.Unlock()

	bucket.BackupLock.Lock()
	n := bucket.GetFirstBackup()
	if n != nil {
		bucket.RemoveBackup(n)
		dt.TimeOutChan <- PingNode{Wait2Ping: n, Dht: dt}
	}
	bucket.BackupLock.Unlock()
}

func (dt *DhtTable) RunTimeOut() {
	defer func() {
		dt.Wg.Done()
	}()

	for {
		select {
		case pn := <-dt.TimeOutChan:
			pn.Dht.DoTimeOut(pn)
		case <-dt.TimeQuit:
			return
		}
	}

}

func (dtb *DTBucket) TimeOut(tv int) {
	now := tools.GetNowMsTime()

	nxt := dtb.Root
	for {
		if nxt == nil {
			return
		}

		if now-nxt.lastPingTime > int64(tv) {
			pn := PingNode{}
			pn.Wait2Ping = nxt
			pn.Dht = dtb.Dht

			dtb.Dht.TimeOutChan <- pn
		}
		nxt = nxt.Next
	}
}

func (dt *DhtTable) TimeOut(tv int) {
	if tv == 0 {
		tv = 3600000 //ms,1hour
	}

	for idx := 0; idx < len(dt.HashTable); idx++ {

		bucket := &dt.HashTable[idx]
		bucket.RootLock.Lock()
		bucket.TimeOut(tv)
		bucket.RootLock.Unlock()

	}

}

func (dt *DhtTable) WrapperTimeOut(tv int) {

	ticker := tools.GetNbsTickerInstance()
	c := make(chan int64, 1)
	ticker.RegWithTimeOut(&c, int64(tv)/2)

	defer func() {
		dt.Wg.Done()
		ticker.UnReg(&c)
	}()

	for {
		select {
		case <-c:
			go dt.TimeOut(tv)
		case <-dt.TimeQuitCreate:
			return
		}
	}

}

func (dt *DhtTable) RunPing() {

	defer func() {
		dt.Wg.Done()
	}()

	for {
		select {
		case node := <-dt.PingNodeChan:
			node.Dht.DoPing(node)
		case <-dt.PingQuit:
			return
		}
	}
}

func (dt *DhtTable) Run(wait bool) {
	go tools.GetNbsTickerInstance().Run()

	dt.Wg.Add(1)
	go dt.RunPing()

	dt.Wg.Add(1)

	go dt.WrapperTimeOut(0)

	dt.Wg.Add(1)
	if wait {
		dt.RunTimeOut()
	} else {
		go dt.RunTimeOut()
	}
}

func (dt *DhtTable) Stop() {
	dt.PingQuit <- 1
	dt.TimeQuitCreate <- 1
	dt.TimeQuit <- 1

	tools.GetNbsTickerInstance().Stop()

	dt.Wg.Wait()
}
