package localaddress

import (
	"net/http"
	"fmt"
	"log"
	"strings"
)

type localaddress struct {

}

func NewLocalAddress() http.Handler  {
	return &localaddress{}
}

func (la *localaddress)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "GET"{
		fmt.Fprintf(w,"{}")
		return
	}

	log.Println(r.RemoteAddr)

	ra := r.RemoteAddr

	arr:=strings.Split(ra,":")

	fmt.Fprintf(w,arr[0])

}

