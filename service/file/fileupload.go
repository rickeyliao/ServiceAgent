package file

import (
	"net/http"
	"fmt"
	"os"
	"io"
)

const (
	maxUploadSize=10*1024*1024   //10M
)


type fileupload struct {

}

func NewFileUpLoad()  http.Handler {
	return &fileupload{}
}

func (fu *fileupload)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	r.ParseMultipartForm(32<<20)
	file,h,err:=r.FormFile("filename")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", h.Header)
	f, err := os.OpenFile("./test/"+h.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	w.Write([]byte("success"))
}