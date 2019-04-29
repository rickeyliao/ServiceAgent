package key

import (
	"net/http"
	"fmt"
)

type keyauth struct {

}

func NewKeyAuth() http.Handler {
	return &keyauth{}
}

func (ka *keyauth)ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "test")
}


