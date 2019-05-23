package listallip

import (
	"net/http"
	"fmt"
	"github.com/rickeyliao/ServiceAgent/common"
)

type listallips struct {

}

func NewListAllIps() http.Handler  {
	return &listallips{}
}

func (us *listallips)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "GET"{
		w.WriteHeader(500)
		fmt.Fprintf(w,"{}")
		return
	}
	ret, code, err := common.Get(common.GetRemoteUrlInst().GetHostName(r.URL.Path))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	w.WriteHeader(code)
	fmt.Fprintf(w, ret)

}