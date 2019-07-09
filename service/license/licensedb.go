package license

import (
	"github.com/rickeyliao/ServiceAgent/db"
	"sync"
)

var (
	licensedb db.NbsDbInter
	licensedblock sync.Mutex
	quit chan int
	wg *sync.WaitGroup
)



