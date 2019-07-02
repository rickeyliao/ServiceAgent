package file

import (
	"net/http"
	"fmt"
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
	filename:=r.Header.Get("FileName")

	fmt.Println(filehash,filename)

	//w.Header().Add("Content-Disposition", "Attachment")
	//begin to download...
	fmt.Println("begin to downloading")

	http.ServeFile(w,r,"/Users/rickey/Downloads/android-studio-ide-173.4819257-mac.dmg")

	//http.ServeContent()

	fmt.Println("end download..")

}
