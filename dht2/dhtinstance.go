package dht2

import "sync"

var (
	dhtCanServiceInstance *DhtTable
	dhtCanServiceInstLock sync.Mutex
	dhtAllNodeInstance    *DhtTable
	dhtAllNodeInstLock    sync.Mutex
)

func GetCanServiceDht() *DhtTable {
	if dhtCanServiceInstance == nil {
		dhtCanServiceInstLock.Lock()
		defer dhtCanServiceInstLock.Unlock()

		if dhtCanServiceInstance == nil {
			dhtCanServiceInstance = NewDhtTable()
		}
	}

	return dhtCanServiceInstance
}

func GetAllNodeDht() *DhtTable {
	if dhtAllNodeInstance == nil {
		dhtAllNodeInstLock.Lock()
		defer dhtAllNodeInstLock.Unlock()

		if dhtAllNodeInstance == nil {
			dhtAllNodeInstance = NewDhtTable()
		}
	}

	return dhtAllNodeInstance
}

func DhtRuning() {
	GetCanServiceDht().Run(false)
	GetAllNodeDht().Run(false)
	go GetKAStore().WrapperTimeout()
	go NbsP2PListen()

}

func DhtStop() {
	StopNbsP2pListen()
	GetCanServiceDht().Stop()
	GetAllNodeDht().Stop()
	GetKAStore().TimeoutStop()
}
