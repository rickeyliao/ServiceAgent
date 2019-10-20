package controller

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/rickeyliao/ServiceAgent/common"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/app/cmdservice/api"
	"fmt"
)



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
	fmt.Println(string(formjson),err)
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

	fmt.Println(bj)

	ac.cookie = lr

	http.SetCookie(w,&cookie)

	w.Write([]byte("true"))

	return

}

//change password

func (ac *AjaxController)ChgpwdDo(w http.ResponseWriter,r *http.Request)  {
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