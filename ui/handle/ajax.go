package handle

import (
	"net/http"
	"strings"
	"github.com/rickeyliao/ServiceAgent/ui/controller"
	"fmt"
	"reflect"
)

type AjaxHandle struct {

}

func (ah *AjaxHandle)ServeHTTP(w http.ResponseWriter, r *http.Request)   {
	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	fmt.Println(r.URL.Path)
	fmt.Println(pathInfo)
	fmt.Println(parts)
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}

	login := &controller.AjaxController{}
	cls := reflect.ValueOf(login)
	method := cls.MethodByName(action)
	if !method.IsValid() {
		method = cls.MethodByName(strings.Title("Login") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})
}
