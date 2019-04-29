package email

import (
	"fmt"
	"net/http"
)

type emailrecord struct {

}

func NewEmailRecord() http.Handler {
	return &emailrecord{}
}

func (er *emailrecord)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w,"email record")
}
