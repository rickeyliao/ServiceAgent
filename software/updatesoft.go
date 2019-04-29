package software

import (
	"fmt"
	"net/http"
)

type updatesoft struct {

}

func NewUpdateSoft() http.Handler  {
	return &updatesoft{}
}

func (us *updatesoft)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w,"update software")
}