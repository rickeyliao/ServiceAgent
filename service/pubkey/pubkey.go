package pubkey

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/rickeyliao/ServiceAgent/common"
	"encoding/json"
	"log"
)

type pubkey struct {
	NbsAddr string `json:"nbsaddr"`
}

func NewHttpPubKey() http.Handler {
	return &pubkey{}
}


func (pk *pubkey)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST"{
		w.WriteHeader(500)
		fmt.Fprintf(w,"{}")
		return
	}

	var err error

	if _,err=ioutil.ReadAll(r.Body);err!=nil{
		w.WriteHeader(500)
		fmt.Fprintf(w,"{}")
		return
	}


	w.WriteHeader(200)
	fmt.Fprintf(w,getNbsAddr())

}

func getNbsAddr() string  {
	addr:=&pubkey{NbsAddr:common.GetSAConfig().NbsRsaAddr}
	if v,err:=json.Marshal(*addr);err!=nil{
		log.Println(err)
		return "{}"
	}else{
		return string(v)
	}
}



