package login

import (
	"net/http"
	"html/template"
	"github.com/rickeyliao/ServiceAgent/common"
	"path"
	"github.com/ribencong/go-manager/kingkey"
)

type LoginInfo struct {
	username string
	password string
}

func NewLoginInfo() http.Handler {
	return &LoginInfo{}
}

func (li *LoginInfo)ServeHTTP(w http.ResponseWriter, r *http.Request)   {

	if r.Method == "GET" {
		sar:=common.GetSARootCfg()
		sac:=common.GetSAConfig()

		t, _ := template.ParseFiles(path.Join(sar.HomeDir,sac.StaticFileDir,sac.Loginfile))
		t.Execute(w, nil)
	} else {
		if err:=r.ParseForm();err!=nil{
			w.Write([]byte("parse form error"))
			return
		}

		if r.Form["username"] == nil || len(r.Form["username"]) != 1{
			w.Write([]byte("Please Enter correct username & password"))
			return
		}

		if r.Form["password"] == nil || len(r.Form["password"]) != 1{
			w.Write([]byte("Please Enter correct username & password"))
			return
		}

		if r.Form["pubaddr"] == nil || len(r.Form["pubaddr"]) != 1{
			w.Write([]byte("Please Enter correct address"))
			return
		}

		if !common.CheckUserPassword(r.Form["username"][0],r.Form["password"][0]){
			w.Write([]byte("Please Enter correct username & password"))
			return
		}

		if l,err:=kingkey.GenLicense("YouPipe2019",r.Form["pubaddr"][0],"",30);err!=nil{
			w.Write([]byte("Address not correct"))
		}else {
			w.Write([]byte(l))
		}

	}

}


