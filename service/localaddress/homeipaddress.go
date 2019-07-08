package localaddress

import (
	"github.com/rickeyliao/ServiceAgent/db"
	"sync"
	"github.com/rickeyliao/ServiceAgent/common"
	"path"
)

var (
	homeipdb db.NbsDbInter
	homeipdblock sync.Mutex
)

func GetHomeIPDB() db.NbsDbInter {
	if homeipdb == nil{
		homeipdblock.Lock()
		defer homeipdblock.Unlock()

		if homeipdb == nil{
			homeipdb = newHomeIPDB()
		}

	}
	return homeipdb
}

func newHomeIPDB() db.NbsDbInter {
	cfg:=common.GetSAConfig()
	return db.NewFileDb(path.Join(cfg.FileDBDir,cfg.HomeIPDBFile)).Load()
}

type Homeipdesc struct {
	MachineName string `json:"MachineName"`
	NbsAddress  string `json:"-"`
	InternetAddress string `json:"IAddress"`
	NatAddress string `json:"nAddress"`
}

func Insert(nbsaddress string,machineName string,InterAddress string,natAddress string)  {
	//hid:=&Homeipdesc{MachineName:machineName,InternetAddress:InterAddress,NatAddress:natAddress}

	//if bhid,err:=json.Marshal(hid);err!=nil{
	//
	//}
}
