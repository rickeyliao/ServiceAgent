package file

import (
	"net/http"
	"github.com/rickeyliao/ServiceAgent/common"
)

type filedownload struct {
}

func NewFileDownLoad() http.Handler {
	return &filedownload{}
}

func (fdl *filedownload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	filehash:=r.Header.Get("FileHash")

	//w.Header().Add("Content-Disposition", "Attachment")
	//begin to download...
	//fmt.Println("begin to downloading")
	filename:=common.GetSaveFilePath(filehash)
	if filename == ""{
		w.Header().Add("message","FileNotFound")
		w.Write([]byte("Failed"))
		return
	}

	http.ServeFile(w,r,filename)

}

