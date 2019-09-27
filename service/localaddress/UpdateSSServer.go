package localaddress

import (
	"sync"
	"time"
	"fmt"
	"github.com/rickeyliao/ServiceAgent/common"
	"encoding/json"
)

var (
	SSAddressList map[string]Homeipdesc
	SSAddressListLock sync.Mutex
)


type SATime struct{
	time.Time
}

func (t *SATime)MarshalJSON() ([]byte,error)  {
	satime:=fmt.Sprintf("\"%s\"",t.Format("2006-01-02 15:04:05"))

	return []byte(satime),nil
}

func (t *SATime)UnmarshalJSON(data []byte) error {
	var err error

	t.Time, err = time.Parse(`"2006-01-02 15:04:05"`, string(data))
	if err != nil {
		return err
	}

	return nil
}

type SSServerList struct {
	CreateDate SATime `json:"createdDate"`
	LastModify SATime `json:"lastModifiedDate"`
	Version int       `json:"version"`
	NodeId  string    `json:"id"`
	Name    string    `json:"name"`
	IPAddress string  `json:"ip"`
	SSPort    int     `json:"port"`
	SSPassword string `json:"password"`
	Location   string `json:"location"`
	LosingTimes int   `json:"losingTimes"`
	Status      int   `json:"status"`
	DeleteFlag  bool  `json:"deleteFlag"`
	Abroad      int   `json:"abroad"`
}

type ServerListPost struct {
	Platform string `json:"platform"`
}

func getServersListPostParam(platform string) (string,error){
	slp:=ServerListPost{Platform:platform}

	if bslp,err:=json.Marshal(slp);err!=nil{
		return "",err
	}else{
		return string(bslp),nil
	}
}

func GetServerList() []SSServerList {

	var p string
	var err error

	if p,err=getServersListPostParam("proxy");err!=nil{
		return nil
	}

	//fmt.Println(p)

	ret, code, err := common.Post(common.GetRemoteUrlInst().GetHostName(common.GetSAConfig().ListIpsPath), p)
	if err != nil {
		return nil
	}
	if code != 200{
		return nil
	}

	//fmt.Println(ret)


	//ssl:=make([]SSServerList,0)

	var ssl []SSServerList

	err = json.Unmarshal([]byte(ret),&ssl)

	//fmt.Println(err)

	//fmt.Println(ssl)

	return ssl
}

