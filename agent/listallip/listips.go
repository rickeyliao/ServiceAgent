package listallip

import (
	"net/http"
	"fmt"
	"github.com/rickeyliao/ServiceAgent/common"
	"io/ioutil"
)

type listallips struct {

}

func NewListAllIps() http.Handler  {
	return &listallips{}
}

func (us *listallips)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST"{
		w.WriteHeader(500)
		fmt.Fprintf(w,"{}")
		return
	}

	var body []byte
	var err error

	if body,err=ioutil.ReadAll(r.Body); err!=nil{
		w.WriteHeader(500)
		fmt.Fprintf(w,"{}")
		return
	}

	var ret string
	var code int
	ret,code,err=common.Post(common.GetRemoteUrlInst().GetHostName(r.URL.Path),string(body))
	if err!=nil{
		w.WriteHeader(500)
		fmt.Fprintf(w,"{}")
		return
	}

	w.WriteHeader(code)
	fmt.Fprintf(w,ret)

}