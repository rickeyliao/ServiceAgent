package postsocks5

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"log"
)

type postsocks5 struct {

}

func NewPostSocks5() http.Handler {
	return &postsocks5{}
}

func (ka *postsocks5)ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	log.Println(string(body))

	w.WriteHeader(200)
	fmt.Fprintf(w,"{}")

	return

}


