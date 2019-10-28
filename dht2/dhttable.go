package dht2

import (
	"sync"
	"fmt"
)

//max node in bukcet
var MaxKBucket int = 16
//max node in backup bucket
var MaxKBucketBackup int = 16

type DTNode struct {
	P2pNode P2pAddr
	lastPingTime int64
	RefCnt int
	Next *DTNode
}

type PingNode struct {
	Wait2Ping *DTNode
	Dht *DhtTable
}

type DTBucket struct {
	RootLock sync.Mutex
	RootCnt int
	Root *DTNode
	BackupLock sync.Mutex
	BackupCnt int
	Backup *DTNode
	Dht *DhtTable
}

type DhtTable struct {
	HashTable [257]DTBucket
	DTLock sync.Mutex
	PingNodeChan chan PingNode
	Quit chan int
}

func NewDhtTable() *DhtTable  {
	dt:=&DhtTable{}

	dt.DTLock.Lock()
	defer dt.DTLock.Unlock()

	dt.PingNodeChan = make(chan PingNode,2560)

	for i:=0;i<len(dt.HashTable);i++{
		dt.HashTable[i].Dht = dt
	}

	dt.Quit = make(chan int,0)


	return dt
}

func (dtn *DTNode)String() string {
	s := (&(dtn.P2pNode)).String()

	s += fmt.Sprintf("lastPingTime %d ",dtn.lastPingTime)
	s += fmt.Sprintf("refCnt %d ",dtn.RefCnt)

	return s
}


func (dtn *DTNode)Clone() *DTNode  {
	dtn1:=&DTNode{}

	dtn1.P2pNode = *((&(dtn.P2pNode)).Clone())

	dtn1.lastPingTime = dtn.lastPingTime
	dtn1.RefCnt = dtn.RefCnt
	dtn1.Next = nil

	return dtn1
}


func (dt *DhtTable)Find(node *DTNode) *DTNode {
	laddr:=GetLocalNAddr()
	dtaddr:=node.P2pNode.NbsAddr

	bucketidx,_:=NbsXorLen(laddr.Bytes(),dtaddr.Bytes())

	bucket := dt.HashTable[bucketidx]

	bucket.RootLock.Lock()
	defer bucket.RootLock.Unlock()


	return (&bucket).Find(node).Clone()
}

func (dtb *DTBucket)Find(node *DTNode) *DTNode  {
	root := dtb.Root

	if root == nil{
		return nil
	}

	nxt:=root

	for{

		if nxt == nil {
			break
		}

		if nxt.P2pNode.NbsAddr.Cmp(node.P2pNode.NbsAddr){
			break
		}
		nxt = nxt.Next

	}

	return nxt
}

func (dtb *DTBucket)CloneAllNotes() []*DTNode {
	dtnodes:=make([]*DTNode,0)
	nxt:=dtb.Root

	for{
		if nxt == nil{
			break
		}

		cpy:=nxt.Clone()
		dtnodes = append(dtnodes,cpy)

		nxt=nxt.Next
	}

	return dtnodes
}

func (dt *DhtTable)FindNearest(node *DTNode,cnt int) ([]*DTNode,int)  {
	laddr:=GetLocalNAddr()
	dtaddr:=node.P2pNode.NbsAddr

	bucketidx,_:=NbsXorLen(laddr.Bytes(),dtaddr.Bytes())

	startbucketidx := bucketidx

	curcnt:=0

	dtnodes:=make([]*DTNode,0)

	for i:=startbucketidx;i<len(dt.HashTable);i++{
		bucket := dt.HashTable[i]
		bucket.RootLock.Lock()

		if bucket.RootCnt>0 {
			dtnodes= append(dtnodes,bucket.CloneAllNotes()...)
			curcnt += bucket.RootCnt

			if curcnt >=cnt{
				bucket.RootLock.Unlock()
				return dtnodes,curcnt
			}
		}

		bucket.RootLock.Unlock()
	}

	if startbucketidx>1{
		for i:=startbucketidx-1;i>0;i--{
			bucket := dt.HashTable[i]
			bucket.RootLock.Lock()

			if bucket.RootCnt>0 {
				dtnodes= append(dtnodes,bucket.CloneAllNotes()...)
				curcnt += bucket.RootCnt

				if curcnt >=cnt{
					bucket.RootLock.Unlock()
					return dtnodes,curcnt
				}
			}

			bucket.RootLock.Unlock()
		}
	}

	return dtnodes,curcnt
}

func (dtb *DTBucket)GetLast() *DTNode {
	prev:=dtb.Root
	nxt:=prev
	for{
		if nxt == nil{
			return prev
		}
		nxt = nxt.Next
	}
}

func (dtb *DTBucket)GetLastBackup() *DTNode {
	prev:=dtb.Backup
	nxt:=prev
	for{
		if nxt == nil{
			return prev
		}
		nxt = nxt.Next
	}
}

func (dtb *DTBucket)GetFirstBackup() *DTNode  {

	root:=dtb.Root

	if root != nil{
		dtb.Root = root.Next
	}

	return root

}


func (dtb *DTBucket)Add(node *DTNode)  {
	nxt:=dtb.Root
	dtb.Root = node
	node.Next = nxt
}

func (dtb *DTBucket)AddBackup(node *DTNode)  {
	nxt:=dtb.Backup
	dtb.Backup = node
	node.Next = nxt
}

func (dtb *DTBucket)Insert(node *DTNode) bool {
	//node.lastPingTime = tools.GetNowMsTime()

	n1:=dtb.Find(node)
	if n1==nil{
		if dtb.RootCnt < MaxKBucket{
			dtb.Add(node)
			dtb.RootCnt ++
			return false
		}

	}else{
		dtb.Remove(n1)
		dtb.Add(node)
		return false
	}

	//rootcnt >= maxkbucket
	pingnode:=PingNode{}
	pingnode.Wait2Ping = dtb.GetLast().Clone()
	pingnode.Dht = dtb.Dht

	dtb.Dht.PingNodeChan <- pingnode

	return true
}

func (dtb *DTBucket)FindBackup(node *DTNode) *DTNode  {
	root := dtb.Backup

	if root == nil{
		return nil
	}

	nxt:=root

	for{

		if nxt == nil {
			break
		}

		if nxt.P2pNode.NbsAddr.Cmp(node.P2pNode.NbsAddr){
			break
		}
		nxt = nxt.Next

	}

	return nxt
}

func (dtb *DTBucket)InsertBackup(node *DTNode)  {
	n1:=dtb.FindBackup(node)
	if n1 == nil{
		if dtb.BackupCnt < MaxKBucketBackup{
			dtb.AddBackup(node)
			dtb.BackupCnt ++
			return
		}
	}else{
		dtb.RemoveBackup(n1)
		dtb.AddBackup(node)
		return
	}

	n2:=dtb.GetLastBackup()
	dtb.RemoveBackup(n2)
	dtb.AddBackup(node)
}

func (dt *DhtTable)Insert(node *DTNode) error  {

	laddr:=GetLocalNAddr()
	dtaddr:=node.P2pNode.NbsAddr

	bucketidx,_:=NbsXorLen(laddr.Bytes(),dtaddr.Bytes())

	bucket := dt.HashTable[bucketidx]

	bucket.RootLock.Lock()
	flag:=bucket.Insert(node)
	bucket.RootLock.Unlock()

	if flag{
		bucket.BackupLock.Lock()
		bucket.InsertBackup(node)
		bucket.BackupLock.Unlock()
	}

	return nil
}

func (dt *DhtTable)Update(node *DTNode)  {
	laddr:=GetLocalNAddr()
	dtaddr:=node.P2pNode.NbsAddr

	bucketidx,_:=NbsXorLen(laddr.Bytes(),dtaddr.Bytes())

	bucket := dt.HashTable[bucketidx]

	bucket.BackupLock.Lock()

	n1:=bucket.GetFirstBackup()

	bucket.RemoveBackup(n1)
	bucket.BackupCnt --

	bucket.BackupLock.Unlock()

	bucket.RootLock.Lock()

	bucket.Remove(node)
	if n1!=nil{
		bucket.Add(n1)
	}
	bucket.RootLock.Unlock()

}


func (dtb *DTBucket)Remove(node *DTNode) {
	nxt:=dtb.Root
	prev:=nxt
	for{
		if nxt == nil{
			return
		}

		if nxt.P2pNode.NbsAddr.Cmp(node.P2pNode.NbsAddr){
			if nxt == prev{
				dtb.Root = nxt.Next
			}else{
				prev.Next = nxt.Next
			}
			return
		}
		prev = nxt
		nxt = nxt.Next
	}
}


func (dtb *DTBucket)RemoveBackup(node *DTNode) {

	if node == nil{
		return
	}

	nxt:=dtb.Backup
	prev:=nxt
	for{
		if nxt == nil{
			return
		}

		if nxt.P2pNode.NbsAddr.Cmp(node.P2pNode.NbsAddr){
			if nxt == prev{
				dtb.Root = nxt.Next
			}else{
				prev.Next = nxt.Next
			}
			return
		}
		prev = nxt
		nxt = nxt.Next
	}
}


func (dt *DhtTable)Remove(node *DTNode) {
	laddr:=GetLocalNAddr()
	dtaddr:=node.P2pNode.NbsAddr

	bucketidx,_:=NbsXorLen(laddr.Bytes(),dtaddr.Bytes())

	bucket := dt.HashTable[bucketidx]

	bucket.RootLock.Lock()
	defer bucket.RootLock.Unlock()

	(&bucket).Remove(node)
}

func (dt *DhtTable)DoPing(node *DTNode)  {
	if node.P2pNode.Ping(){
		dt.Insert(node)
	}else{
		dt.Update(node)
	}
}

func (dt *DhtTable)RunPing(wg *sync.WaitGroup)  {

	defer func() {
		wg.Done()
	}()

	for{
		select {
		case node:=<-dt.PingNodeChan:
			node.Dht.DoPing(node.Wait2Ping)
		case <-dt.Quit:
			return
		}
	}


}

func (dt *DhtTable)Run() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go dt.RunPing(wg)


	wg.Wait()
}

func (dt *DhtTable)Stop()  {
	dt.Quit <- 1

}

