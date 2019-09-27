package localaddress

import (
	"fmt"
	"net/http"
	"strings"
	"strconv"
)

type localaddress struct {

}

func NewLocalAddress() http.Handler {
	return &localaddress{}
}

func (la *localaddress) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "{}")
		return
	}

	ra := r.RemoteAddr

	arr := strings.Split(ra, ":")

	//log.Println(r.RemoteAddr,r.Header.Get("nbsaddress"))

	nataddrs:=r.Header.Get("nataddrs")
	hostname:=r.Header.Get("hostname")

	nbsaddr := r.Header.Get("nbsaddress")


	inas,_:=strconv.Atoi(r.Header.Get("nationality"))

	if len(nbsaddr) >0{
		Insert(nbsaddr,hostname,arr[0],nataddrs,int32(inas))


	}


	w.Header().Add("Connection","close")
	fmt.Fprintf(w, arr[0])


}
