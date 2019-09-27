package license

import (
	"encoding/json"
	"github.com/kprc/nbsnetwork/db"
	"github.com/pkg/errors"
	"github.com/rickeyliao/ServiceAgent/common"
	"path"
	"strconv"
	"sync"
	"time"
)

var (
	licensedb     db.NbsDbInter
	licensedblock sync.Mutex
	quit          chan int
	wg            *sync.WaitGroup
)

type LicenseDesc struct {
	SofaAddress string `json:"-"`
	User        string `json:"user"`
	NDays       int32  `json:"ndays"`
	StrLicense  string `json:"license"`
	Ipaddr      string `json:"ipaddr"`
}

func GetLicenseDB() db.NbsDbInter {
	if licensedb == nil {
		licensedblock.Lock()
		defer licensedblock.Unlock()

		if licensedb == nil {
			licensedb = newLicenseDB()
		}
	}

	return licensedb
}

func newLicenseDB() db.NbsDbInter {
	quit = make(chan int, 0)
	wg = &sync.WaitGroup{}
	cfg := common.GetSAConfig()
	return db.NewFileDb(path.Join(cfg.GetFileDbDir(), cfg.LicenseDBFile)).Load()
}

func Insert(sofaaddr string, user string, ndays int32, ipaddr string, l string) error {
	if sofaaddr == "" || len(sofaaddr) == 0 {
		return errors.New("No sofa address")
	}

	if l == "" {
		return errors.New("No license")
	}

	arlicense := make([]LicenseDesc, 0)
	v, _ := licensedb.Find(sofaaddr)
	if v != "" {
		json.Unmarshal([]byte(v), &arlicense)
	}

	license := LicenseDesc{User: user, NDays: ndays, StrLicense: l, Ipaddr: ipaddr}

	arlicense = append(arlicense, license)

	if bl, err := json.Marshal(arlicense); err != nil {
		return err
	} else {
		GetLicenseDB().Update(sofaaddr, string(bl))
	}

	return nil
}

func CmdLicenseShow(sofaaddress string) string {
	v, err := GetLicenseDB().Find(sofaaddress)
	if err != nil {
		return "Not Found"
	}

	arlicense := make([]LicenseDesc, 0)

	err = json.Unmarshal([]byte(v), &arlicense)
	if err != nil {
		return "Internal error"
	}

	r := ""
	for _, l := range arlicense {
		if r == "" {
			r = "SofaAddress: " + sofaaddress
		}
		r += "\r\nUser: " + l.User
		r += "\r\nIP: " + l.Ipaddr
		r += "\r\nNDays: " + strconv.Itoa(int(l.NDays))
		r += "\r\nLicense: " + l.StrLicense
	}

	return r
}

func CmdShowLicenseAll() string {
	dbcusor := GetLicenseDB().DBIterator()

	alls := ""

	for {
		k, _ := dbcusor.Next()
		if k == "" {
			break
		}
		if alls != "" {
			alls += "\r\n"
		}
		alls += CmdLicenseShow(k)
	}

	return alls
}

func CmdShowLicenseStatistic() string {
	dbcusor := GetLicenseDB().DBIterator()

	cnt := 0
	lcnt := 0
	r := ""

	for {
		k, v := dbcusor.Next()
		if k == "" {
			break
		}
		cnt++

		arlicense := make([]LicenseDesc, 0)
		err := json.Unmarshal([]byte(v), &arlicense)
		if err != nil {
			return "Internal error"
		}

		lcnt += len(arlicense)
	}

	r += "sofaAddress Count: " + strconv.Itoa(cnt)
	r += "\r\nLicense Count: " + strconv.Itoa(lcnt)

	return r
}

func CmdShowLicenseSummary() string {
	dbcusor := GetLicenseDB().DBIterator()

	cnt := 0
	lcnt := 0
	r := ""

	for {
		k, v := dbcusor.Next()
		if k == "" {
			break
		}
		cnt++

		if r != "" {
			r += "\r\n"
		}

		arlicense := make([]LicenseDesc, 0)
		err := json.Unmarshal([]byte(v), &arlicense)
		if err != nil {
			return "Internal error"
		}

		tlcnt := len(arlicense)

		r += "sofaaddress: " + k
		if tlcnt > 0 {
			r += " count: " + strconv.Itoa(tlcnt)
			r += " user: " + arlicense[0].User
			r += " ndays: " + strconv.Itoa(int(arlicense[0].NDays))
		}
		lcnt += tlcnt
	}

	r += "\r\nsofaAddress Count: " + strconv.Itoa(cnt)
	r += "\r\nLicense Count: " + strconv.Itoa(lcnt)

	return r
}
func Save() {
	GetLicenseDB().Save()
}

func IntervalSave() {
	if wg == nil {
		GetLicenseDB()
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
