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
	"github.com/rickeyliao/ServiceAgent/app"
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
		if k == ""{
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

type Homeipdesc struct {
	MachineName     string   `json:"MachineName,omitempty"`
	NbsAddress      string   `json:"-"`
	InternetAddress string   `json:"IAddress,omitempty"`
	NatAddress      []string `json:"nAddress,omitempty"`
	Nationality     int32    `json:"nationality,omitempty"`
	SSPort          int      `json:"port,omitempty"`
	SSPassword      string   `json:"password,omitempty"`
	SSMethod        string   `json:"ssmethod,omitempty"`
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

	fmt.Println("GetNbsAddrByIP......")
	hi := GetHomeIPDB()

	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	fmt.Println("GetNbsAddrByIP------")

	for k, v := range hi.memdb {
		if v.InternetAddress == ip {
			return k
		}
	}

	return ""

}

func GetHomeIPDescByNbsaddr(nbsaddr string) *Homeipdesc {

	fmt.Println("GetHomeIPDescByNbsaddr......")
	hi := GetHomeIPDB()
	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()
	fmt.Println("GetHomeIPDescByNbsaddr-----")

	v,_:=hi.memdb[nbsaddr]

	return v.Clone()
}


func (hid *Homeipdesc)Clone() *Homeipdesc {
	hid1:=&Homeipdesc{}

	*hid1 = *hid

	nataddrs:=make([]string,0)

	for _,addr:=range hid.NatAddress{
		nataddrs = append(nataddrs,addr)
	}

	hid1.NatAddress = nataddrs

	return hid1
}


func UpdateToServer(sn *SSServerListNode) (del,add bool, hid *Homeipdesc) {

	fmt.Println("UpdateToServer......")
	hi:=GetHomeIPDB()
	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()
	fmt.Println("UpdateToServer-----")

	delflag:=false
	addflag:=false



	v,ok:=hi.memdb[sn.Name]
	if v.SSPassword == ""{
		return
	}
	if !ok{
		delflag = true
	}else{
		if sn.IPAddress != v.InternetAddress || sn.SSPassword != v.SSPassword || sn.SSPort != v.SSPort {
			delflag = true
			addflag = true
		}
	}

	return delflag,addflag,v.Clone()
}

func UpdateToServers(srvl []*SSServerListNode,delsrv []string,addsrv []*Homeipdesc,nas int32)  {

	fmt.Println("UpdateToServers......")
	hi:=GetHomeIPDB()
	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	fmt.Println("UpdateToServers------")

	memhids:=make(map[string]*Homeipdesc,0)

	for k,v:=range hi.memdb{
		if nas == 0 ||
			(nas == app.NATIONALITY_CHINA_MAINLAND && v.Nationality == nas) ||
			(nas >0 && v.Nationality == app.NATIONALITY_AMERICAN) ||
			(nas >0 && v.Nationality == app.NATIONALITY_JAPANESE) ||
			(nas >0 && v.Nationality == app.NATIONALITY_SINGAPORE) ||
			(nas >0 && v.Nationality == app.NATIONALITY_ENGLAND){
				if v.SSPassword != ""{
					memhids[k] = v
				}
		}
	}

	keys:=make(map[string]struct{},0)

	for _,ssl:=range srvl{
		keys[ssl.Name]= struct{}{}
		v,ok:=memhids[ssl.Name]
		if !ok{
			delsrv = append(delsrv,v.NbsAddress)
		}else{
			if ssl.IPAddress != v.InternetAddress || ssl.SSPassword != v.SSPassword || ssl.SSPort !=  v.SSPort{
				delsrv = append(delsrv,ssl.Name)
				v.NbsAddress = ssl.Name
				addsrv = append(addsrv,v)
			}
		}
	}

	for k,v:=range memhids{
		if _,ok:=keys[k];!ok{
			v.NbsAddress = k
			addsrv = append(addsrv,v)
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
		hid = &Homeipdesc{MachineName: mn, InternetAddress: interAddress, NatAddress: String2arr(natAddress)}
	} else {
		hid = &Homeipdesc{MachineName: mn, InternetAddress: interAddress, NatAddress: String2arr(natAddress),
			Nationality: ssr.Nationality,
			SSPassword:  ssr.SSPassword, SSPort: ssr.SSPort, SSMethod: ssr.SSMethod}
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
	fmt.Println("memdbInsert......")
	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()
	fmt.Println("memdbInsert.-----")

	hid.NbsAddress = nbsaddr
	hi.memdb[nbsaddr] = hid
}

func CmdShowAddress(nbsaddr string) string {

	fmt.Println("CmdShowAddress...")

	hi := GetHomeIPDB()

	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	fmt.Println("CmdShowAddress---")

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

	for _, nip := range hid.NatAddress {

		r += fmt.Sprintf("%-16s", nip)
	}

	return r
}

func CmdShowAddressAll(nas int32) string {

	fmt.Println("CmdShowAddressAll......")
	hi := GetHomeIPDB()
	fmt.Println("CmdShowAddressAll.=======")
	hi.memdblock.Lock()
	defer hi.memdblock.Unlock()

	fmt.Println("CmdShowAddressAll----")

	return hi.CmdShowAddressAll(nas)

}

func (hi *HomeIPDB) CmdShowAddressAll(nas int32) string {

	alls := ""
	for k, v := range hi.memdb {
		if v.Nationality != 0 && v.Nationality != nas {
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
