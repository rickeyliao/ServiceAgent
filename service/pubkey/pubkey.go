package pubkey

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/rickeyliao/ServiceAgent/common"
	"io/ioutil"
	"log"
	"net/http"
)

type pubkey struct {
	NbsAddr string `json:"nbsaddr"`
	Pubkey  string `json:"pubkey"`
}

func NewHttpPubKey() http.Handler {
	return &pubkey{}
}

func (pk *pubkey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	var err error

	if _, err = ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	w.WriteHeader(200)
	w.Header().Add("Connection", "close")
	fmt.Fprintf(w, GetNbsPubkey())

}

func GetNbsPubkey() string {
	sac := common.GetSAConfig()
	addr := &pubkey{NbsAddr: sac.NbsRsaAddr, Pubkey: sac.GetPubKey()}
	fmt.Println("send to my pk:", addr.Pubkey)
	fmt.Println("nbsaddr:", sac.NbsRsaAddr)
	if v, err := json.Marshal(*addr); err != nil {
		log.Println(err)
		return "{}"
	} else {
		return string(v)
	}
}

func UnMarshalPubKey(pkjson []byte) (addr string, pk *rsa.PublicKey) {
	p := &pubkey{}

	err := json.Unmarshal(pkjson, p)
	if err != nil {
		return "", nil
	}

	pk = common.ToPubKey(p.Pubkey)

	fmt.Println("----->",p.NbsAddr,p.Pubkey)

	return p.NbsAddr, pk
}
