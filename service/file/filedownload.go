package file

import "net/http"

type filedownload struct {
}

func NewFileDownLoad() http.Handler {
	return &filedownload{}
}

func (fdl *filedownload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

}
