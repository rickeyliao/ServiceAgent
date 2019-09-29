package localaddress

import (
	"encoding/json"
	"fmt"
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

type HomeIPDB struct {
	homeipdb     db.NbsDbInter
	homeipdblock sync.Mutex
	quit         chan int
	wg           *sync.WaitGroup
	memdb        map[string]*Homeipdesc
	memdblock    sync.Mutex
}

var (
	homeipdbinst     *HomeIPDB
	homeipdbinstlock sync.Mutex
)

func GetHomeIPDB() *HomeIPDB {
	if homeipdbinst == nil {
		homeipdbinstlock.Lock()
		defer homeipdbinstlock.Unlock()

		if homeipdbinst == nil {
			homeipdbinst = newHomeIPDB()
		}

	}
	return homeipdbinst
}

func memDBLoad(memdb map[string]*Homeipdesc, fdb db.NbsDbInter) {
	dbcusor := fdb.DBIterator()

	if dbcusor == nil {
		return
	}

	for {
		k, v := dbcusor.Next()
		hid := &Homeipdesc{}
		err := json.Unmarshal([]byte(v), hid)
		if err != nil {
			continue
		}

		memdb[k] = hid
	}
}

func newHomeIPDB() *HomeIPDB {

	hi := &HomeIPDB{}
	hi.quit = make(chan int, 0)
	hi.wg = &sync.WaitGroup{}

	cfg := common.GetSAConfig()
	hi.homeipdb = db.NewFileDb(path.Join(cfg.GetFileDbDir(), cfg.HomeIPDBFile)).Load()
	hi.memdb = make(map[string]*Homeipdesc)
	memDBLoad(hi.memdb, hi.homeipdb)

	return hi
}

type Homeipdesc struct {
	MachineName     string   `json:"MachineName"`
	NbsAddress      string   `json:"-"`
	InternetAddress string   `json:"IAddress"`
	NatAddress      []string `json:"nAddress"`
	Nationality     int32    `json:"nationality"`
	SSPort          int      `json:"port"`
	SSPassword      string   `json:"password"`
	Location        string   `json:"location"`
	SSMethod        string   `json:"ssmethod"`
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

func Insert(nbsaddress string, mn string, interAddress string, natAddress string, ssr *SSReport) error {

	if interAddress == "" || len(interAddress) == 0 {
		return errors.New("No Internat address")
	}

	if nbsaddress == "" || len(nbsaddress) == 0 {
		return errors.New("nbsaddress not found")
	}

	var hid *Homeipdesc

	if ssr == nil {
		hid = &Homeipdesc{MachineName: mn, InternetAddress: interAddress, NatAddress: String2arr(natAddress)}
	} else {
		hid = &Homeipdesc{MachineName: mn, InternetAddress: interAddress, NatAddress: String2arr(natAddress),
			Nationality: ssr.Nationality,
			SSPassword:  ssr.SSPassword, SSPort: ssr.SSPort, Location: ssr.Location, SSMethod: ssr.SSMethod}
	}

	GetHomeIPDB().Insert(nbsaddress, hid)

	return nil
}

func (hi *HomeIPDB) Insert(nbsaddr string, hid *Homeipdesc) error {
	if len(nbsaddr) == 0 {
		return errors.New("nbs address not found")
	}

	hi.memdbInsert(nbsaddr, hid)

	if bhid, err := json.Marshal(*hid); err != nil {
		return err
	} else {
		hi.homeipdblock.Lock()
		defer hi.homeipdblock.Unlock()

		hi.homeipdb.Update(nbsaddr, string(bhid))
	}

	return nil
}

func (hi *HomeIPDB) memdbInsert(nbsaddr string, hid *Homeipdesc) {
	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	hi.memdb[nbsaddr] = hid
}

func CmdShowAddress(nbsaddr string) string {
	hi := GetHomeIPDB()

	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	return hi.CmdShowAddress(nbsaddr)
}

func (hi *HomeIPDB) CmdShowAddress(nbsaddr string) string {

	hid, ok := hi.memdb[nbsaddr]
	if !ok {
		return "Not found"
	}

	r := fmt.Sprintf("%-48s", nbsaddr)

	r += fmt.Sprintf("%-16s", hid.MachineName)
	r += fmt.Sprintf("%-18s", hid.InternetAddress)
	r += fmt.Sprintf("%-6s", strconv.Itoa(int(hid.Nationality)))
	r += fmt.Sprintf("%-6s", strconv.Itoa(int(hid.SSPort)))
	r += fmt.Sprintf("%-16s", hid.SSPassword)
	r += fmt.Sprintf("%-16s", hid.SSMethod)
	r += fmt.Sprintf("%-16s", hid.Location)

	for _, nip := range hid.NatAddress {

		r += fmt.Sprintf("%-16s", nip)
	}

	return r
}

func CmdShowAddressAll() string {
	hi := GetHomeIPDB()

	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	return hi.CmdShowAddressAll()

}

func (hi *HomeIPDB) CmdShowAddressAll() string {

	alls := ""
	for k, _ := range hi.memdb {
		if alls != "" {
			alls += "\r\n"
		}
		alls += hi.CmdShowAddress(k)
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
	hi := GetHomeIPDB()
	hi.homeipdblock.Lock()
	defer hi.homeipdblock.Unlock()
	hi.homeipdb.Save()
}

func IntervalSave() {
	hi := GetHomeIPDB()

	hi.wg.Add(1)
	defer hi.wg.Done()
	var count int64 = 0
	for {

		if count%86400 == 0 {
			Save()
		}
		count++

		select {
		case <-hi.quit:
			return
		default:
		}
		time.Sleep(time.Second * 1)
	}

}

func Destroy() {
	hi := GetHomeIPDB()
	hi.quit <- 1

	Save()

	hi.wg.Wait()
}
