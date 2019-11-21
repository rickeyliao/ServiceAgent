package controller

import (
	"net/http"
	"os/exec"
	"github.com/rickeyliao/ServiceAgent/common"
	"path"
	"log"
	"github.com/kprc/flowsharectrl/config"
	"encoding/json"
	"runtime"
)

func (ac *AjaxController) ChangeUplinkDo(w http.ResponseWriter, r *http.Request) {
	ups,ok:=r.URL.Query()["uplink"]
	if !ok || len(ups)==0{
		w.WriteHeader(500)
		w.Write([]byte("{}"))
		return
	}

	if runtime.GOARCH != "arm"{
		w.Write([]byte("only arm platform can do this callback"))
		w.WriteHeader(200)
		return
	}

	up := ups[0]

	wifidir:=common.GetSAConfig().GetWifiDir()

	res:="wifiap/staticfile"

	if common.WifiRes != ""{
		res = common.WifiRes
	}

	wifidir = path.Join(wifidir,res)

	shcmd := ""

	if up == "eth"{
		shcmd = path.Join(wifidir,"change2eth.sh")
	}

	if up == "ppp"{
		shcmd = path.Join(wifidir,"change2ppp.sh")
	}

	log.Println(shcmd)


	cmd := exec.Command("/bin/sh",shcmd,wifidir)

	err := cmd.Run()

	w.WriteHeader(200)
	if err!=nil{
		w.Write([]byte(err.Error()))
	}else{
		w.Write([]byte("success"))
	}

	return

}

func (ac *AjaxController)RebootDo(w http.ResponseWriter, r *http.Request)  {

	cmd:=exec.Command("/sbin/reboot")
	err:=cmd.Run()
	if err!=nil{
		w.Write([]byte(err.Error()))
	}else{
		w.Write([]byte("wooooooo"))
	}


	return

}

type CurUplink struct {
	Cur string  `json:"curuplink"`
}

func (ac *AjaxController)CurrentUplinkDo(w http.ResponseWriter, r *http.Request) {

	if runtime.GOARCH != "arm"{
		w.Write([]byte("only arm platform can do this callback"))
		w.WriteHeader(200)
		return
	}

	cfg:=config.Reload()

	cu:=CurUplink{}

	if cfg.Flag4g {
		cu.Cur = "4g"
	}else{
		cu.Cur = "eth"
	}

	bdi, err := json.Marshal(cu)
	if err != nil {
		w.Write([]byte("error"))
	} else {
		w.Write(bdi)
	}

}