package localaddress

import (
	"encoding/json"
	"fmt"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/pkg/errors"
	"github.com/rickeyliao/ServiceAgent/app"
	"github.com/rickeyliao/ServiceAgent/common"
	"log"
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
		if k == "" {
			return
		}
		hid := &Homeipdesc{}
		err := json.Unmarshal([]byte(v), hid)
		if err != nil {
			continue
		}

		hid.NbsAddress = k

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

func (hidb *HomeIPDB) Find(nbsaddr string) *Homeipdesc {
	hidb.memdblock.Lock()
	defer hidb.memdblock.Unlock()

	if v, ok := hidb.memdb[nbsaddr]; !ok {
		return nil
	} else {
		return v.Clone()
	}

}

func (hidb *HomeIPDB) FindByIP(ipaddr string) (hid *Homeipdesc, nbsaddr string) {
	hidb.memdblock.Lock()
	defer hidb.memdblock.Unlock()

	for k, v := range hidb.memdb {
		if v.InternetAddress == ipaddr {
			return v.Clone(), k
		}
	}

	return nil, ""
}

type Homeipdesc struct {
	MachineName     string   `json:"MachineName,omitempty"`
	NbsAddress      string   `json:"-"`
	InternetAddress string   `json:"IAddress,omitempty"`
	NatAddress      []string `json:"nAddress,omitempty"`
	Nationality     int32    `json:"nationality,omitempty"`
	SSPort          int      `json:"port,omitempty"`
	SSPassword      string   `json:"password,omitempty"`
	SSMethod        string   `json:"ssmethod,omitempty"`
	LastUpdate      int64    `json:"-"`
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

func GetNbsAddrByIP(ip string) string {

	hi := GetHomeIPDB()

	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	for k, v := range hi.memdb {
		if v.InternetAddress == ip {
			return k
		}
	}

	return ""

}

func GetHomeIPDescByNbsaddr(nbsaddr string) *Homeipdesc {

	hi := GetHomeIPDB()
	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	v, _ := hi.memdb[nbsaddr]

	return v.Clone()
}

func (hid *Homeipdesc) Clone() *Homeipdesc {
	hid1 := &Homeipdesc{}

	*hid1 = *hid

	nataddrs := make([]string, 0)

	for _, addr := range hid.NatAddress {
		nataddrs = append(nataddrs, addr)
	}

	hid1.NatAddress = nataddrs

	return hid1
}

func (hid *Homeipdesc) String() string {
	r := fmt.Sprintf("%-48s", trim(hid.NbsAddress, 46))

	r += fmt.Sprintf("%-20s", trim(hid.MachineName, 18))
	r += fmt.Sprintf("%-18s", hid.InternetAddress)
	r += fmt.Sprintf("%-6s", strconv.Itoa(int(hid.Nationality)))
	r += fmt.Sprintf("%-6s", strconv.Itoa(int(hid.SSPort)))
	r += fmt.Sprintf("%-20s", trim(hid.SSPassword, 18))
	r += fmt.Sprintf("%-18s", trim(hid.SSMethod, 16))

	for _, nip := range hid.NatAddress {

		r += fmt.Sprintf("%-16s", nip)
	}

	return r
}

func getHomeIPDesc(memdb map[string]*Homeipdesc, srvName string) (hid *Homeipdesc, ok bool) {
	for k, v := range memdb {
		if k[:16] == srvName {
			return v, true
		}
	}

	return nil, false
}

func UpdateToServer(sn *SSServerListNode) (del, add bool, hid *Homeipdesc) {

	hi := GetHomeIPDB()
	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	delflag := false
	addflag := false

	v, ok := getHomeIPDesc(hi.memdb, sn.Name)
	if !ok {
		delflag = true

	} else {
		if v.SSPassword == "" {
			return
		}

		if sn.IPAddress != v.InternetAddress || sn.SSPassword != v.SSPassword || sn.SSPort != v.SSPort {
			delflag = true
			addflag = true
			hid = v.Clone()
		}
	}

	return delflag, addflag, hid
}

func UpdateToServers(srvl []*SSServerListNode, delsrv *[]string, addsrv *[]*Homeipdesc, nas int32) {

	hi := GetHomeIPDB()
	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	memhids := make(map[string]*Homeipdesc, 0)

	for k, v := range hi.memdb {
		if nas == 0 ||
			(nas == app.NATIONALITY_CHINA_MAINLAND && v.Nationality == nas) ||
			(nas == app.NATIONALITY_AMERICAN && v.Nationality == app.NATIONALITY_AMERICAN) ||
			(nas == app.NATIONALITY_AMERICAN && v.Nationality == app.NATIONALITY_JAPANESE) ||
			(nas == app.NATIONALITY_AMERICAN && v.Nationality == app.NATIONALITY_SINGAPORE) ||
			(nas == app.NATIONALITY_AMERICAN && v.Nationality == app.NATIONALITY_ENGLAND) {
			if v.SSPassword != "" {
				memhids[k] = v
			}
		}
	}

	keys := make(map[string]struct{}, 0)

	for _, ssl := range srvl {
		keys[ssl.Name] = struct{}{}
		v, ok := memhids[ssl.Name]
		if !ok {
			*delsrv = append(*delsrv, ssl.IPAddress)
			log.Println(ssl.IPAddress)
		} else {
			if ssl.IPAddress != v.InternetAddress || ssl.SSPassword != v.SSPassword || ssl.SSPort != v.SSPort {
				*delsrv = append(*delsrv, ssl.IPAddress)
				v.NbsAddress = ssl.Name
				*addsrv = append(*addsrv, v)
			}
		}
	}

	for k, v := range memhids {
		if _, ok := keys[k]; !ok {
			v.NbsAddress = k
			*addsrv = append(*addsrv, v)
		}
	}

	return
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
		hid = &Homeipdesc{MachineName: mn, InternetAddress: interAddress,
			NatAddress: String2arr(natAddress), LastUpdate: tools.GetNowMsTime()}
	} else {
		hid = &Homeipdesc{MachineName: mn, InternetAddress: interAddress, NatAddress: String2arr(natAddress),
			Nationality: ssr.Nationality,
			SSPassword:  ssr.SSPassword, SSPort: ssr.SSPort, SSMethod: ssr.SSMethod,
			LastUpdate: tools.GetNowMsTime()}
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

func (hi *HomeIPDB) Delete(nbsaddr string) error {
	if len(nbsaddr) == 0 {
		return errors.New("nbs address is empty")
	}

	hid := hi.memdbDelete(nbsaddr)
	if hid == nil {
		return nil
	}

	hi.homeipdblock.Lock()
	defer hi.homeipdblock.Unlock()

	hi.homeipdb.Delete(nbsaddr)

	return nil
}

func (hi *HomeIPDB) memdbInsert(nbsaddr string, hid *Homeipdesc) {
	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	hid.NbsAddress = nbsaddr
	hi.memdb[nbsaddr] = hid
}

func (hi *HomeIPDB) memdbDelete(nbsaddr string) *Homeipdesc {
	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	v, ok := hi.memdb[nbsaddr]
	if ok {
		delete(hi.memdb, nbsaddr)
	}
	return v
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

	r := hid.String()

	return r
}

func CmdShowAddressAll(nas int32) string {

	hi := GetHomeIPDB()
	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	return hi.CmdShowAddressAll(nas)

}

func (hi *HomeIPDB) CmdShowAddressAll(nas int32) string {

	alls := ""
	for k, v := range hi.memdb {
		if nas != 0 && v.Nationality != nas {
			continue
		}
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

func timeout2deletes() []string {

	keys := make([]string, 0)

	now := tools.GetNowMsTime()
	hi := GetHomeIPDB()
	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	for k, v := range hi.memdb {
		if v.LastUpdate == 0 || now-v.LastUpdate < 36000 {
			continue
		}
		delete(hi.memdb, k)
		keys = append(keys, k)

	}

	return keys
}

func TimeOut() {

	keys := timeout2deletes()

	hi := GetHomeIPDB()
	hi.homeipdblock.Lock()
	defer hi.homeipdblock.Unlock()
	for _, key := range keys {
		hi.homeipdb.Delete(key)
	}
}

func IntervalSave() {
	hi := GetHomeIPDB()

	hi.wg.Add(1)
	defer hi.wg.Done()
	var count int64 = 0
	for {

		if count%86400 == 0 {
			TimeOut()
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
