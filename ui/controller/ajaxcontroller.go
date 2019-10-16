package controller

import (
	"net/http"
	"encoding/json"
)

type Result struct{
	Ret int
	Reason string
	Data interface{}
}

type AjaxController struct {
}

func (this *AjaxController)LoginAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	err := r.ParseForm()
	if err != nil {
		OutputJson(w, 0, "参数错误", nil)
		return
	}

	admin_name := r.FormValue("admin_name")
	admin_password := r.FormValue("admin_password")

	if admin_name == "" || admin_password == ""{
		OutputJson(w, 0, "参数错误", nil)
		return
	}

	if admin_name != "admin" || admin_password != "passwd"{
		OutputJson(w,0,"User name or password error",nil)
	}


	// 存入cookie,使用cookie存储
	//expiration := time.Unix(1, 0)
	cookie := http.Cookie{Name: "admin_name", Value: admin_name, Path: "/"}
	http.SetCookie(w, &cookie)

	OutputJson(w, 1, "操作成功", nil)
	return
}

func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := &Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	w.Write(b)
}
