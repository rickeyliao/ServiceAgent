package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"github.com/rickeyliao/ServiceAgent/common"
	"path"
	"github.com/kprc/nbsnetwork/tools"
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

	f, err := os.OpenFile(getSaveFilePath(h.Filename), os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {

		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	w.Write([]byte("success"))
	fmt.Println("save file", h.Filename, "success")
}

func getSaveFilePath(filename string) string  {
	sac:=common.GetSAConfig()

	arrpath:=getArrPath(filename)

	filepath := path.Join(sac.Root.HomeDir,sac.FileStoreDir)

	for i:=0;i<len(arrpath); i++{
		filepath = path.Join(filepath,arrpath[i])
	}

	if !tools.FileExists(filepath){
		os.MkdirAll(filepath,0755)
	}

	absfilename := path.Join(filepath,filename)

	return absfilename
}

func getArrPath(filename string) []string  {
	arrpath:=make([]string,0)

	s:=[]byte(filename)

	for i:=len(filename);i>0;i=i-2 {
		s:=s[:i]
		if len(s) >=2{
			arrpath = append(arrpath,string(s[i-2:i]))

		}else{
			break
		}
		if len(arrpath) >= 4{
			break
		}
	}

	if len(arrpath)>0{
		arrret:=make([]string,0)
		for i:=len(arrpath)-1; i>=0; i--{
			arrret = append(arrret,arrpath[i])
		}
		return arrret
	}

	return nil
}