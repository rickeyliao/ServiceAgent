package software

import (
	"net/http"
	"fmt"
	"github.com/rickeyliao/ServiceAgent/common"
)

type updatesoft struct {

}

func NewUpdateSoft() http.Handler  {
	return &updatesoft{}
}

func (us *updatesoft)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "POST"{
		var err error

		var ret string
		var code int
		ret,code,err=common.Post(common.GetRemoteUrlInst().GetHostName(r.URL.Path),"")
		if err!=nil{
			w.WriteHeader(500)
			fmt.Fprintf(w,"{}")
			return
		}
		w.WriteHeader(code)
		fmt.Fprintf(w,ret)
	}else if r.Method == "GET"{
		ret, code, err := common.Get(common.GetRemoteUrlInst().GetHostName(r.URL.Path))
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "{}")
			return
		}
		w.WriteHeader(code)
		fmt.Fprintf(w, ret)
	}else{
		w.WriteHeader(500)
		fmt.Fprintf(w,"{}")
		return
	}
}