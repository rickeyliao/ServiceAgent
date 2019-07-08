package localaddress

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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

	log.Println(r.RemoteAddr,arr[0],r.Header.Get("nbsaddress"))

	w.Header().Add("Connection","close")
	fmt.Fprintf(w, arr[0])


}
