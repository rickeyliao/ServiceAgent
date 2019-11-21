package dht2

import (
	"github.com/kprc/nbsnetwork/tools"
	"net"
	"sync"
	"time"
	"bytes"
)

type KANode struct {
	nbsaddr        NAddr
	sn             []byte
	ip             net.IP
	port           int
	lastAccessTime int64
	next           *KANode
}

type KABucket struct {
	lock sync.Mutex
	root *KANode
}

type KAStore struct {
	HashTable [256]KABucket
	lock      sync.Mutex
	quit      chan int
	wg        *sync.WaitGroup
}

var (
	kastore     *KAStore
	kastoreLock sync.Mutex
)

func NewKAStore() *KAStore {
	return &KAStore{quit: make(chan int, 1), wg: &sync.WaitGroup{}}
}

func GetKAStore() *KAStore {
	if kastore == nil {
		kastoreLock.Lock()
		defer kastoreLock.Unlock()
		if kastore == nil {
			kastore = NewKAStore()
		}
	}
	return kastore
}

func (na NAddr) KAHash() int {
	return int(na[31])
}

func (kn *KANode) clone() *KANode {
	n1 := &KANode{}

	n1.nbsaddr = kn.nbsaddr
	n1.port = kn.port
	n1.ip = kn.ip
	n1.lastAccessTime = kn.lastAccessTime

	return n1

}

func (kb *KABucket) find(nbsaddr NAddr, port int, ip net.IP) *KANode {
	r := kb.root

	for {
		if r == nil {
			return nil
		}

		if r.nbsaddr.Cmp(nbsaddr) && r.port == port && r.ip.Equal(ip) {
			return r
		}

		r = r.next
	}
}

func (kb *KABucket) delete(nbsaddr NAddr, port int, ip net.IP) {
	r := kb.root
	prev := r
	for {
		if r == nil {
			return
		}

		if r.nbsaddr.Cmp(nbsaddr) && r.port == port && r.ip.Equal(ip) {

			if r == kb.root {
				kb.root = r.next
			} else {
				prev.next = r.next
			}

			return
		}
		prev = r
		r = r.next
	}
}

func (kb *KABucket) deleteall(nbsaddr NAddr) {
	r := kb.root
	prev := r
	for {
		if r == nil {
			return
		}

		if r.nbsaddr.Cmp(nbsaddr) {

			if r == kb.root {
				kb.root = r.next
			} else {
				prev.next = r.next
			}

			r = r.next
		} else {
			prev = r
			r = r.next
		}

	}
}

func (kb *KABucket) findall(nbsaddr NAddr) []*KANode {
	r := kb.root

	arr := make([]*KANode, 0)

	for {
		if r == nil {
			break
		}

		if nbsaddr.Cmp(r.nbsaddr) {
			arr = append(arr, r)
		}

		r = r.next

	}

	return arr
}

func (kb *KABucket) insert(n *KANode) {
	nxt := kb.root
	kb.root = n
	n.next = nxt
}

//if node have been existed, refresh access time, if not, insert it
func (ks *KAStore) Insert(ip net.IP, port int, nbsaddr NAddr, sn []byte) {
	h := nbsaddr.KAHash()

	b := ks.HashTable[h]

	b.lock.Lock()
	defer b.lock.Unlock()

	n := b.find(nbsaddr, port, ip)
	if n != nil {
		n.lastAccessTime = tools.GetNowMsTime()
		return
	}

	n = &KANode{nbsaddr: nbsaddr, ip: ip, port: port, lastAccessTime: tools.GetNowMsTime(), sn: sn}

	b.insert(n)
}

func (ks *KAStore) Find(nbsaddr NAddr) []*KANode {
	h := nbsaddr.KAHash()
	b := ks.HashTable[h]

	b.lock.Lock()
	defer b.lock.Unlock()

	ns := b.findall(nbsaddr)

	arr := make([]*KANode, 0)

	for _, n := range ns {
		arr = append(arr, n.clone())
	}

	return arr

}

func (kb *KABucket) findBySn(nbsaddr NAddr, sn []byte) *KANode {
	r := kb.root
	for {
		if r == nil {
			return nil
		}

		if nbsaddr.Cmp(r.nbsaddr) && bytes.Compare(sn,r.sn)==0 {
			return r
		}

		r = r.next

	}
}

func (ks *KAStore) FindBySn(nbsaddr NAddr, sn []byte) *KANode {
	h := nbsaddr.KAHash()
	b := ks.HashTable[h]

	b.lock.Lock()
	defer b.lock.Unlock()

	return b.findBySn(nbsaddr, sn)
}

func (ks *KAStore) Delete(nbsaddr NAddr, port int, ip net.IP) {
	h := nbsaddr.KAHash()
	b := ks.HashTable[h]

	b.lock.Lock()
	defer b.lock.Unlock()
	b.delete(nbsaddr, port, ip)
}

func (ks *KAStore) DeleteAll(nbsaddr NAddr) {
	h := nbsaddr.KAHash()
	b := ks.HashTable[h]

	b.lock.Lock()
	defer b.lock.Unlock()

	b.deleteall(nbsaddr)
}

func (kb *KABucket) timeout() {
	now := tools.GetNowMsTime()

	r := kb.root
	prev := r
	for {
		if r == nil {
			return
		}

		if now-r.lastAccessTime > 3600000 {
			if r == kb.root {
				kb.root = r.next
			} else {
				prev.next = r.next
			}

			r = r.next
		} else {
			prev = r
			r = r.next
		}
	}
}

func (ks *KAStore) Timeout() {
	for i := 0; i < len(ks.HashTable); i++ {
		b := ks.HashTable[i]
		b.lock.Lock()
		b.timeout()
		b.lock.Unlock()
	}
}

func (ks *KAStore) WrapperTimeout() {

	ks.wg.Add(1)

	defer func() {
		ks.wg.Done()
	}()

	starttime := tools.GetNowMsTime()
	for {

		select {
		case <-ks.quit:
			return
		default:

		}

		if tools.GetNowMsTime()-starttime < 1800000 {
			time.Sleep(time.Second * 1)
			continue
		}
		ks.Timeout()

		starttime = tools.GetNowMsTime()
	}
}

func (ks *KAStore) TimeoutStop() {
	ks.quit <- 1
	ks.wg.Wait()
}
