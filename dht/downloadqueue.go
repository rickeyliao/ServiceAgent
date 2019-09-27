package dht

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/rickeyliao/ServiceAgent/dht/dhtdb"
	"math/big"
	"net"
	"sync"
)

type ddesc struct {
	key string
	ip  string
}

type dqchan struct {
	key  chan *ddesc
	quit chan int
}

type downloadqueue struct {
	q []*dqchan
}

type DownloadQueue interface {
	EnQueue(key []byte, ip net.IP)
	DeQueue(i int) *ddesc
	Run()
	Stop()
}

const (
	DownloadQueueCnt    int = 8
	DonwloadQueueLength int = 1024
)

var (
	dqinst     *downloadqueue
	dqinstlock sync.Mutex
	dqwg       *sync.WaitGroup
)

func GetDownloadQueue() DownloadQueue {
	if dqinst == nil {
		dqinstlock.Lock()
		defer dqinstlock.Unlock()
		if dqinst == nil {
			dqinst = newDownloadQueue()
		}
	}

	return dqinst

}

func newDownloadQueue() *downloadqueue {
	dqueue := &downloadqueue{make([]*dqchan, DownloadQueueCnt)}

	for i := 0; i < len(dqueue.q); i++ {
		dc := &dqchan{key: make(chan *ddesc, DonwloadQueueLength), quit: make(chan int, 1)}
		dqueue.q[i] = dc
	}

	dqwg = &sync.WaitGroup{}

	return dqueue
}

func (dqueue *downloadqueue) EnQueue(key []byte, ip net.IP) {
	bgi := &big.Int{}
	bgi.SetBytes(key)
	i := bgi.BitLen() & (0x07)

	desc := &ddesc{"c1" + base58.Encode(key), ip.String()}
	q := dqueue.q[i]

	select {
	case q.key <- desc:
	default:
	}
}

func (dqueue *downloadqueue) DeQueue(i int) *ddesc {
	if i < 0 || i >= DownloadQueueCnt {
		return nil
	}

	q := dqueue.q[i]

	select {
	case key := <-q.key:
		return key
	case <-q.quit:
		return nil
	}

	return nil
}

func (dqueue *downloadqueue) Run() {
	for i := 0; i < DownloadQueueCnt; i++ {
		dqwg.Add(1)

		go dqueue.doDownloadFile(i)

	}

}

func (dqueue *downloadqueue) doDownloadFile(i int) {
	defer dqwg.Done()

	for {
		d := dqueue.DeQueue(i)

		if d == nil {
			return
		}
		if !dhtdb.GetFileExistFlag(d.key) {
			if err := common.DownloadFile(d.ip, "", d.key); err == nil {
				dhtdb.Update(d.key, false, false, nil, true)
			}
		}

	}
}

func (dqueue *downloadqueue) Stop() {
	for i := 0; i < DownloadQueueCnt; i++ {
		q := dqueue.q[i]
		q.quit <- 1
	}

	dqwg.Wait()
}
