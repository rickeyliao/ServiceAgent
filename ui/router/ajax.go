package router

import (
	"net/http"
	"strings"
	"github.com/rickeyliao/ServiceAgent/test/oldui/controller"
	"reflect"
)

type AjaxRouter struct {

}

func (ar *AjaxRouter)ServeHTTP(w http.ResponseWriter,r *http.Request)  {

	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")

	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Do"
	}

	login := &controller.AjaxController{}
	cls := reflect.ValueOf(login)
	method := cls.MethodByName(action)
	if !method.IsValid() {
		method = cls.MethodByName(strings.Title("login") + "Do")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})

}
