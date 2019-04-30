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

	if r.Method != "GET"{
		fmt.Fprintf(w,"{}")
		return
	}

	ret,err:=common.Get(common.GetRemoteUrlInst().GetHostName(r.URL.Path))
	if err!=nil{
		fmt.Fprintf(w,"{}")
		return
	}
	fmt.Fprintf(w,ret)
}