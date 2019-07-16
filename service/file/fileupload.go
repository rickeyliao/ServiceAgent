package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"github.com/rickeyliao/ServiceAgent/common"
)

const (
	maxUploadSize = 1 << 20 //1M
)

type fileupload struct {
}

func NewFileUpLoad() http.Handler {
	return &fileupload{}
}

func (fu *fileupload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(maxUploadSize)

	defer r.MultipartForm.RemoveAll()
	file, h, err := r.FormFile("FileHash")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	f, err := os.OpenFile(common.GetSaveFilePath(h.Filename), os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {

		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	w.Write([]byte("success"))
	fmt.Println("save file", h.Filename, "success")
}

