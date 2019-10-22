package router

import (
	"github.com/rickeyliao/ServiceAgent/ui/controller"
	"net/http"
	"reflect"
	"strings"
)

type AjaxRouter struct {
}

func (ar *AjaxRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")

	cookie, err := r.Cookie("nbsadmin")
	if err != nil || cookie.Value == "" {
		if !(len(parts) > 1 && strings.ToLower(parts[1]) == "login") {
			http.Redirect(w, r, "/login.html", http.StatusFound)
			return
		}
	}

	var action = ""
	if len(parts) > 1 {
		for _, part := range parts[1:] {
			action += strings.Title(part)
		}
		action += "Do"
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
