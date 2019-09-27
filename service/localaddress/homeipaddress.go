package localaddress

import (
	"encoding/json"
	"github.com/kprc/nbsnetwork/db"
	"github.com/pkg/errors"
	"github.com/rickeyliao/ServiceAgent/common"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	homeipdb     db.NbsDbInter
	homeipdblock sync.Mutex
	quit         chan int
	wg           *sync.WaitGroup
)

func GetHomeIPDB() db.NbsDbInter {
	if homeipdb == nil {
		homeipdblock.Lock()
		defer homeipdblock.Unlock()

		if homeipdb == nil {
			homeipdb = newHomeIPDB()
		}

	}
	return homeipdb
}

func newHomeIPDB() db.NbsDbInter {
	quit = make(chan int, 0)
	wg = &sync.WaitGroup{}
	cfg := common.GetSAConfig()
	return db.NewFileDb(path.Join(cfg.GetFileDbDir(), cfg.HomeIPDBFile)).Load()
}

type Homeipdesc struct {
	MachineName     string   `json:"MachineName"`
	NbsAddress      string   `json:"-"`
	InternetAddress string   `json:"IAddress"`
	NatAddress      []string `json:"nAddress"`
	Nationality     int32    `json:"nationality"`
}

func String2arr(ips string) []string {
	return strings.Split(ips, "=>")
}

func LocalIPArr2string(iparr []string) string {
	ips := ""
	for _, ip := range iparr {
		if ips != "" {
			ips += "=>"
		}

		ips += ip
	}

	return ips
}

func Insert(nbsaddress string, mn string, interAddress string, natAddress string, nationality int32) error {

	if interAddress == "" || len(interAddress) == 0 {
		return errors.New("No Internat address")
	}

	if nbsaddress == "" || len(nbsaddress) == 0 {
		return errors.New("nbsaddress not found")
	}

	hid := &Homeipdesc{MachineName: mn, InternetAddress: interAddress, NatAddress: String2arr(natAddress), Nationality: nationality}

	if bhid, err := json.Marshal(hid); err != nil {
		return err
	} else {
		GetHomeIPDB().Update(nbsaddress, string(bhid))
	}

	return nil
}

func CmdShowAddress(nbsaddr string) string {
	v, err := GetHomeIPDB().Find(nbsaddr)
	if err != nil {
		return "Not found"
	}
	hid := &Homeipdesc{}

	err = json.Unmarshal([]byte(v), hid)

	if err != nil {
		return "Internal error"
	}

	r := "NbsAddr:" + nbsaddr
	r += "\t"
	r += "MachineName:" + hid.MachineName
	r += "\t"
	r += "InternetAddress:" + hid.InternetAddress
	r += "\t"
	r += "Nationality:" + strconv.Itoa(int(hid.Nationality))
	r += "\r\n"

	nataddrs := ""
	for _, nip := range hid.NatAddress {
		if nataddrs != "" {
			nataddrs += "\t"
		}
		nataddrs += nip
	}

	r += "InternalAddress:" + nataddrs

	return r
}

func CmdShowAddressAll() string {
	dbcusor := GetHomeIPDB().DBIterator()

	alls := ""

	for {
		k, _ := dbcusor.Next()
		if k == "" {
			break
		}
		if alls != "" {
			alls += "\r\n"
		}
		alls += CmdShowAddress(k)
	}

	return alls
}

func GetMachineName() string {
	mn := common.GetSAConfig().HostName
	if mn == "" {
		mn, _ = os.Hostname()
		if mn == "" {
			mn = "nbsmachinename"
		}
	}
	return mn
}

func Save() {
	GetHomeIPDB().Save()
}

func IntervalSave() {
	if wg == nil {
		GetHomeIPDB()
	}
	wg.Add(1)
	defer wg.Done()
	var count int64 = 0
	for {

		if count%86400 == 0 {
			Save()
		}
		count++

		select {
		case <-quit:
			return
		default:
		}
		time.Sleep(time.Second * 1)
	}

}

func Destroy() {
	quit <- 1

	Save()

	wg.Wait()
}
