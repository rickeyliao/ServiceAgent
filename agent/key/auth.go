package key

import (
	"fmt"
	"github.com/rickeyliao/ServiceAgent/common"
	"io/ioutil"
	"net/http"
)

type keyauth struct {
}

func NewKeyAuth() http.Handler {
	return &keyauth{}
}

func (ka *keyauth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	var body []byte
	var err error

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	var ret string
	var code int
	ret, code, err = common.Post(common.GetRemoteUrlInst().GetHostName(r.URL.Path), string(body))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	w.WriteHeader(code)
	fmt.Fprintf(w, ret)

	return

}
