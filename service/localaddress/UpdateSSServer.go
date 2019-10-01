package localaddress

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/nbsnetwork/tools/crypt/nbscrypt"
	"github.com/pkg/errors"
	"github.com/rickeyliao/ServiceAgent/app"
	"github.com/rickeyliao/ServiceAgent/common"
	"log"
	"strconv"
	"strings"
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

type SSServerListNode struct {
	CreateDate  *SATime `json:"createdDate,omitempty"`
	LastModify  *SATime `json:"lastModifiedDate,omitempty"`
	Version     int     `json:"version,omitempty"`
	NodeId      string  `json:"id,omitempty"`
	Name        string  `json:"name"`
	IPAddress   string  `json:"ip"`
	SSPort      int     `json:"port"`
	SSPassword  string  `json:"password"`
	Location    string  `json:"location"`
	LosingTimes int     `json:"losingTimes,omitempty"`
	Status      int     `json:"status"`
	DeleteFlag  bool    `json:"deleteFlag,omitempty"`
	Abroad      int     `json:"abroad"`
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

func GetServerList() []SSServerListNode {

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

	var ssl []SSServerListNode

	err = json.Unmarshal([]byte(ret), &ssl)

	return ssl
}

type SSReport struct {
	Nationality int32  `json:"nationality"`
	SSPort      int    `json:"port"`
	SSPassword  string `json:"password"`
	SSMethod    string `json:"ssmethod"`
}

func GetSSReport() *SSReport {
	ssr := &SSReport{}

	sac := common.GetSAConfig()

	ssr.Nationality = sac.Nationality
	ssr.SSPort = int(sac.ShadowSockPort)
	ssr.SSPassword = sac.GetSSPasswd()
	ssr.SSMethod = sac.GetSSMethod()

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

func DeleteServer(ips []string) error {

	bjip, err := json.Marshal(ips)
	if err != nil {
		log.Println(err)
		return errors.New("Internal Error")
	}

	fmt.Println(string(bjip))

	ret, code, err := common.Post(common.GetRemoteUrlInst().GetHostName(common.GetSAConfig().ServerDelete), string(bjip))
	if err != nil {
		log.Println(err)
		return errors.New("Can't interactive with remote server")
	}
	if code != 200 {
		log.Println("ServerDelete Post response code", code)
		return errors.New("Remote Server Response error code:" + strconv.Itoa(code))
	}

	if !strings.Contains(strings.ToUpper(ret), "OK") {
		log.Println("Delete Server Internal error", ips)
		return errors.New("Remote Server Response a error message:" + ret)
	}

	return nil

}

func AddServer(nbsaddr string, hi *Homeipdesc) error {

	ssnode := &SSServerListNode{}
	ssnode.SSPort = hi.SSPort
	ssnode.SSPassword = hi.SSPassword
	ssnode.IPAddress = hi.InternetAddress
	ssnode.Location = hi.MachineName
	if hi.Nationality == app.NATIONALITY_CHINA_MAINLAND {
		ssnode.Abroad = app.ABROAD_CHINA_MAINLAND
	} else {
		ssnode.Abroad = app.ABROAD_AMERICAN
	}
	ssnode.Name = nbsaddr
	ssnode.Status = 1

	hdarr := []*SSServerListNode{ssnode}
	bjhda, err := json.Marshal(hdarr)
	if err != nil {
		log.Println(err)
		return errors.New("Internal Error")
	}

	fmt.Println(string(bjhda))

	ret, code, err := common.Post(common.GetRemoteUrlInst().GetHostName(common.GetSAConfig().ServerAdd), string(bjhda))
	if err != nil {
		log.Println(err)
		return errors.New("Can't interactive with remote server")
	}
	if code != 200 {
		log.Println("ServerAdd Post response code", code)
		return errors.New("Remote Server Response error code:" + strconv.Itoa(code))
	}
	if !strings.Contains(strings.ToUpper(ret), "OK") {
		log.Println("ServerAdd Server Internal error", hi.InternetAddress)
		return errors.New("Remote Server Response a error message:" + ret)
	}

	return nil

}

func CmdDeleteServer(nationality int32) string {
	l := GetServerList()

	ips := make([]string, 0)

	for _, n := range l {
		if nationality == app.NATIONALITY_AMERICAN && n.Abroad == app.ABROAD_AMERICAN {
			ips = append(ips, n.IPAddress)
		}
		if nationality != 0 && nationality > app.NATIONALITY_AMERICAN && n.Abroad == app.ABROAD_CHINA_MAINLAND {
			ips = append(ips, n.IPAddress)
		}
		if nationality == 0 {
			ips = append(ips, n.IPAddress)
		}

	}

	err := DeleteServer(ips)
	if err != nil {
		log.Println(err)
		return "Internal error,Please check the error log"
	} else {
		return deleteCmdMsg(ips)
	}

}

func CmdDeleteServerByIP(ip string) string {
	ips := []string{ip}

	err := DeleteServer(ips)
	if err != nil {
		log.Println(err)
		return "Internal error,Please check the error log"
	} else {
		return deleteCmdMsg(ips)
	}

}

func deleteCmdMsg(ips []string) string {
	message := ""

	if len(ips) > 0 {

		for _, ip := range ips {
			if message == "" {
				message = "Delete Server List:"
			}
			message += "\r\n"

			message += fmt.Sprintf("        %-20s", ip)
		}

	} else {
		message = "no IP delete"
	}

	return message
}

func getServer(srvl []SSServerListNode,nbsaddr string) *SSServerListNode  {
	for _,ssl:=range srvl{
		if ssl.Name == nbsaddr{
			return &ssl
		}
	}

	return nil
}


func UpdateServer(nas int32, ip string, nbsaddr string) string {
	srvl := GetServerList()

	delsrv := make([]string, 0)
	addsrv := make([]*Homeipdesc, 0)

	if ip != "" {
		nbsaddr = GetNbsAddrByIP(ip)
		if nbsaddr == "" {
			return "Not found the Server by ip"
		}
	}

	if nbsaddr != ""{
		ssl:=getServer(srvl,nbsaddr)
		if ssl != nil{
			d,a,hid:=UpdateToServer(ssl)
			if d{
				delsrv = append(delsrv,ssl.IPAddress)
			}
			if a{
				addsrv = append(addsrv,hid)
			}
		}else{
			hid:=GetHomeIPDescByNbsaddr(nbsaddr)
			if hid.SSPassword != ""{
				addsrv = append(addsrv,hid)
			}
		}
	}else{
		arrssl:=make([]*SSServerListNode,0)

		for _,ssl:=range srvl{
			if (nas == 0) ||
				(nas==app.NATIONALITY_CHINA_MAINLAND && ssl.Abroad == app.ABROAD_CHINA_MAINLAND) ||
				((nas ==app.NATIONALITY_AMERICAN ||
					nas == app.NATIONALITY_JAPANESE ||
					nas == app.NATIONALITY_SINGAPORE ||
					nas == app.NATIONALITY_ENGLAND) && ssl.Abroad == app.ABROAD_AMERICAN){
				arrssl = append(arrssl,&ssl)
			}
		}

		UpdateToServers(arrssl,delsrv,addsrv,nas)

	}

	messageDel := ""

	delfault := false

	if len(delsrv) > 0{
		if err:=DeleteServer(delsrv);err==nil{
			messageDel = "Delete ips:\r\n" + deleteCmdMsg(delsrv)
		}else{
			log.Println("del ips failed",delsrv)
			delfault = true
		}
	}

	messageAdd := ""

	addfault := false

	if len(addsrv)>0{

		for _,add:=range addsrv{
			if err:=AddServer(add.NbsAddress,add);err==nil{
				if messageAdd == ""{
					messageAdd += "Add Server List:\r\n"
				}
				messageAdd += fmt.Sprintf("%-48s",add.NbsAddress)
				messageAdd += fmt.Sprintf("%-18s",add.InternetAddress)
				messageAdd += "\r\n"
			}else{
				addfault = true
			}
		}
	}

	if messageDel != ""{
		messageDel += "\r\n"
	}

	message:=messageDel + messageAdd

	if message == ""{
		if delfault || addfault {
			message = "Something wrong"
		}
	}


	return message
}
