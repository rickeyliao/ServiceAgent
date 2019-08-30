package checkip

import (
	"net/http"
	"github.com/rickeyliao/ServiceAgent/common"
	"html/template"
	"path"
	"github.com/kprc/nbsnetwork/tools/privateip"
)

type CheckPrivateIP struct {

}

func NewCheckPrivateIP() *CheckPrivateIP  {
	return &CheckPrivateIP{}
}

func (cpip *CheckPrivateIP)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		sar:=common.GetSARootCfg()
		sac:=common.GetSAConfig()

		t, _ := template.ParseFiles(path.Join(sar.HomeDir,sac.StaticFileDir,sac.CheckIPFile))
		t.Execute(w, nil)
	} else {
		if err:=r.ParseForm();err!=nil{
			w.Write([]byte("parse form error"))
			return
		}

		if r.Form["ipaddr"] == nil || len(r.Form["ipaddr"]) != 1{
			w.Write([]byte("Please Enter correct ipaddr"))
			return
		}

		if privateip.IsPrivateIPStr(r.Form["ipaddr"][0]){
			w.Write([]byte("<h1>"+r.Form["ipaddr"][0] + " is Private IP address</h1>"))
		}else{
			w.Write([]byte("<h1>"+r.Form["ipaddr"][0] + " is not  Private IP address</h1>"))
		}

	}
}


