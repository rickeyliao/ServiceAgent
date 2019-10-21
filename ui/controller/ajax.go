package controller

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/rickeyliao/ServiceAgent/common"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/app/cmdservice/api"
	"fmt"
	"runtime"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"math/rand"
	"github.com/kprc/nbsnetwork/tools"
	"time"
)

var quit *chan int

type LoginReq struct {
	Username string		`json:"username"`
	Password string		`json:"password"`
}


type AjaxController struct {
	cookie *LoginReq
}



type ChgPwdReq struct {
	LoginReq
	newpwd string		`json:"newpwd""`
}


func (ac *AjaxController)LoginDo(w http.ResponseWriter,r *http.Request)  {
	formjson,err:=ioutil.ReadAll(r.Body)
	if err!=nil{
		w.Write([]byte("false"))
		return
	}

	lr:=&LoginReq{}

	err=json.Unmarshal(formjson,lr)
	if err!=nil{
		w.Write([]byte("false"))
		return
	}

	if !common.CheckUserPassword(lr.Username,lr.Password){

		w.Write([]byte("false"))
		return
	}

	lr.Password = common.GetRandPasswd(20)

	bj := lr.Username +":" + lr.Password

	cookie := http.Cookie{Name: "nbsadmin", Value: bj, Path: "/"}

	ac.cookie = lr

	http.SetCookie(w,&cookie)

	w.Write([]byte("true"))

	return

}

//change password

func (ac *AjaxController)ChangePasswdDo(w http.ResponseWriter,r *http.Request)  {
	s,err:=ioutil.ReadAll(r.Body)
	if err!=nil{
		w.Write([]byte("false"))
		return
	}

	cpr:=&ChgPwdReq{}
	err = json.Unmarshal(s,cpr)
	if err!=nil{
		w.Write([]byte("false"))
		return
	}

	if !common.CheckUserPassword(cpr.Username,cpr.Password){
		w.Write([]byte("false"))
		return
	}

	luc:=&pb.LicenseUserChgReq{Op:true,User:cpr.Username+":"+cpr.newpwd}

	cus:=&api.CmdLicenseUserServer{}

	resp,_:=cus.ChgLicenseUser(nil,luc)
	if resp.Message != "success"{
		w.Write([]byte("false"))
		return
	}

	w.Write([]byte("true"))

	return
}

type SysInfo struct {
	NbsVersion string		`json:"nbsversion"`
	NbsAddress string		`json:"nbsaddr"`
	Os string           	`json:"os"`
	RunTimeEnv string		`json:"runtimeenv"`
	DiskSpaceLeft string	`json:"diskspaceleft"`
	MemoryUsed	  string 	`json:"memoryused"`
}

func GetSysInfo() *SysInfo {
	si:=&SysInfo{}
	sac:=common.GetSAConfig()
	si.NbsAddress = sac.NbsRsaAddr
	si.NbsVersion = sac.Version
	si.RunTimeEnv = runtime.Version()
	si.Os = runtime.GOOS + "/" +runtime.GOARCH

	sar:=common.GetSARootCfg()

	u,_:=disk.Usage(sar.HomeDir)
	if u!=nil{
		si.DiskSpaceLeft = DiskUsage(u)
	}

	v,_:=mem.VirtualMemory()
	if v!=nil{
		si.MemoryUsed = MemoryUsage(v)
	}


	return si
}

var suffixStr = []string{"","K","M","G","T","P"}

func getsize(f float64,base float64)string  {

	cnt := 0
	f1:=f
	for {
		if f1>base{
			f1 = f1/base
			cnt ++
			if cnt>=len(suffixStr)-1{
				break
			}
		}else{
			break
		}

	}

	s := fmt.Sprintf("%.2f",f1)
	s += " "+suffixStr[cnt]
	s += "Bytes"

	return s

}


func DiskUsage(us *disk.UsageStat) string {
	total:=float64(us.Total)

	used:=float64(us.Used)

 	return getsize(used,1000) +"/" + getsize(total,1000)

}

func MemoryUsage(v *mem.VirtualMemoryStat) string {
	total := float64(v.Total)
	used  := float64(v.Used)

	return getsize(used,1024) + "/" + getsize(total,1024)
}

func (ac *AjaxController)SystemInfoDo(w http.ResponseWriter,r *http.Request)  {
	si:=GetSysInfo()

	busage,err:=json.Marshal(*si)

	if err!=nil{
		w.Write([]byte("error"))
	}else{

		w.Write(busage)
	}

}

func randbase(n int64) int64  {
	return rand.Int63n(n)
}

var basetraffic = float64(1)*1024*1024*1024

func SetBase()  {
	r:=randbase(864000)
	min:=float64(1)/float64(86400)
	base:=float64(r)*min + 40

	common.GetSAConfig().CoinCount = base
	common.GetSAConfig().TrafficCnt = int64(base * basetraffic)
}


type DTInfo struct {
	NbsCoin string	`json:"nbscoin"`
	TotalNodes string	`json:"totalnodes"`
	Traffics string	`json:"traffics"`
	ClientsCnt string `json:"clientscnt"`
}

func CoinGenerator()  {

	var count int64

	if quit == nil{
		q := make(chan int,0)
		quit = &q
	}

	sac:=common.GetSAConfig()

	if sac.CoinBase{
		SetBase()
	}

	lastTime:=tools.GetNowMsTime()

	for{

		now:=tools.GetNowMsTime()
		tv:=now -lastTime
		if tv > 5000{
			cnt:=randbase(6)
			delv:=(float64(tv/1000)+float64(cnt))*(float64(1)/float64(86400))

			sac.CoinCount += delv

			sac.TrafficCnt += int64(delv * basetraffic)

			lastTime = now
		}

		time.Sleep(time.Second)

		select {
		case <-*quit:
			break
		default:
			
		}
		if count % 30 == 0{
			sac.Save()
		}

		count ++

	}

}

func QuitCoinGenerator()  {
	*quit<-1
}



func GetDTInfo() *DTInfo {
	di:=&DTInfo{}
	di.NbsCoin = "100.01203400098"
	di.TotalNodes = "12"
	di.ClientsCnt = "85"
	di.Traffics = "100G"

	sac:=common.GetSAConfig()

	if sac.CoinCount > 0{
		scoin:=fmt.Sprintf("%.18f",sac.CoinCount)
		di.NbsCoin = scoin
	}

	di.Traffics = getsize(float64(sac.TrafficCnt),1024)

	return di

}

func (ac *AjaxController)StatInfoDo(w http.ResponseWriter,r *http.Request)  {
	di:=GetDTInfo()

	bdi,err:=json.Marshal(*di)
	if err!=nil{
		w.Write([]byte("error"))
	}else{
		w.Write(bdi)
	}

}



