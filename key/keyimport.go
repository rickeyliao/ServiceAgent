package key

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/rickeyliao/ServiceAgent/common"
)

type keyimport struct {

}

func NewKeyImport() http.Handler  {
	return &keyimport{}
}

func (ki *keyimport)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
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