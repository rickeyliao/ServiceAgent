package email

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/rickeyliao/ServiceAgent/common"
)

type emailrecord struct {

}

func NewEmailRecord() http.Handler {
	return &emailrecord{}
}

func (er *emailrecord)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST"{
		fmt.Fprintf(w,"{}")
		return
	}

	var body []byte
	var err error

	if body,err=ioutil.ReadAll(r.Body); err!=nil{
		fmt.Fprintf(w,"{}")
		return
	}

	var ret string
	ret,err=common.Post(common.GetRemoteUrlInst().GetHostName(r.URL.Path),string(body))
	if err!=nil{
		fmt.Fprintf(w,"{}")
		return
	}

	fmt.Fprintf(w,ret)
}
