package key

import (
	"net/http"
	"fmt"
)

type keyimport struct {

}

func NewKeyImport() http.Handler  {
	return &keyimport{}
}

func (ki *keyimport)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w,"key import")
}