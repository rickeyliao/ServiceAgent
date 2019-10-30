package dht2

import (
	"github.com/kprc/nbsnetwork/tools"
	"net"
	"sync"
)

type KANode struct {
	nbsaddr        NAddr
	sn             int64
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
}

func NewKAStore() *KAStore {
	return &KAStore{}
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
func (ks *KAStore) Insert(ip net.IP, port int, nbsaddr NAddr, sn int64) {
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

func (kb *KABucket) findBySn(nbsaddr NAddr, sn int64) *KANode {
	r := kb.root
	for {
		if r == nil {
			return nil
		}

		if nbsaddr.Cmp(r.nbsaddr) && sn == r.sn {
			return r
		}

		r = r.next

	}
}

func (ks *KAStore) FindBySn(nbsaddr NAddr, sn int64) *KANode {
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
