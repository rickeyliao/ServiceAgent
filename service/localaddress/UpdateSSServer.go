package localaddress

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/nbsnetwork/tools/crypt/nbscrypt"
	"github.com/rickeyliao/ServiceAgent/common"
	"log"
	"sync"
	"time"
)

var (
	SSAddressList     map[string]Homeipdesc
	SSAddressListLock sync.Mutex
)

type SATime struct {
	time.Time
}

func (t *SATime) MarshalJSON() ([]byte, error) {
	satime := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))

	return []byte(satime), nil
}

func (t *SATime) UnmarshalJSON(data []byte) error {
	var err error

	t.Time, err = time.Parse(`"2006-01-02 15:04:05"`, string(data))
	if err != nil {
		return err
	}

	return nil
}

type SSServerList struct {
	CreateDate  SATime `json:"createdDate"`
	LastModify  SATime `json:"lastModifiedDate"`
	Version     int    `json:"version"`
	NodeId      string `json:"id"`
	Name        string `json:"name"`
	IPAddress   string `json:"ip"`
	SSPort      int    `json:"port"`
	SSPassword  string `json:"password"`
	Location    string `json:"location"`
	LosingTimes int    `json:"losingTimes"`
	Status      int    `json:"status"`
	DeleteFlag  bool   `json:"deleteFlag"`
	Abroad      int    `json:"abroad"`
}

type ServerListPost struct {
	Platform string `json:"platform"`
}

func getServersListPostParam(platform string) (string, error) {
	slp := ServerListPost{Platform: platform}

	if bslp, err := json.Marshal(slp); err != nil {
		return "", err
	} else {
		return string(bslp), nil
	}
}

func GetServerList() []SSServerList {

	var p string
	var err error

	if p, err = getServersListPostParam("proxy"); err != nil {
		return nil
	}

	ret, code, err := common.Post(common.GetRemoteUrlInst().GetHostName(common.GetSAConfig().ListIpsPath), p)
	if err != nil {
		return nil
	}
	if code != 200 {
		return nil
	}

	var ssl []SSServerList

	err = json.Unmarshal([]byte(ret), &ssl)

	return ssl
}

type SSReport struct {
	Nationality int32  `json:"nationality"`
	SSPort      int    `json:"port"`
	SSPassword  string `json:"password"`
	Location    string `json:"location"`
	SSMethod    string `json:"ssmethod"`
}

func GetSSReport() *SSReport {
	ssr := &SSReport{}

	sac := common.GetSAConfig()

	ssr.Nationality = sac.Nationality
	ssr.SSPort = int(sac.ShadowSockPort)
	ssr.SSPassword = sac.GetSSPasswd()
	ssr.SSMethod = sac.GetSSMethod()
	ssr.Location = sac.Location

	return ssr
}

func toSSReport(ssrstr string) *SSReport {

	bssr := base58.Decode(ssrstr)

	bjson, err := nbscrypt.DecryptRsa(bssr, common.GetSAConfig().PrivKey)
	if err != nil {
		log.Println(err)
		return nil
	}

	ssr := &SSReport{}

	err = json.Unmarshal(bjson, ssr)
	if err != nil {
		log.Println(err)
		return nil
	}

	return ssr
}
