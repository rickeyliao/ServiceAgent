package handle

import (
	"net/http"
	"strings"
	"reflect"
	"github.com/rickeyliao/ServiceAgent/ui/controller"
	"fmt"
)

type LoginHandle struct {

}

func (lh *LoginHandle)ServeHTTP(w http.ResponseWriter, r *http.Request)  {

	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	fmt.Println(r.URL.Path)
	fmt.Println(pathInfo)
	fmt.Println(parts)
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}

	login := &controller.LoginController{}
	cls := reflect.ValueOf(login)
	method := cls.MethodByName(action)
	if !method.IsValid() {
		method = cls.MethodByName(strings.Title("index") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})
}