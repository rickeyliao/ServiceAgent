package api

import (
	"context"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"net"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
	"os"
	"github.com/minio/sha256-simd"
	"io"
	"github.com/btcsuite/btcutil/base58"
	"mime/multipart"
	"net/http"
	"io/ioutil"
	"fmt"
)

type CmdFileUpLoad struct {
	
}

func (cfu *CmdFileUpLoad)Uploadfile(ctx context.Context,req  *pb.Fileuploadreq) (*pb.DefaultResp, error)  {
	if req.Filepath=="" || req.Hostip == ""{
		return encResp("Param error"),nil
	}
	if net.ParseIP(req.Hostip) == nil{
		return encResp("host ip address error"),nil
	}

	if !tools.FileExists(req.Filepath){
		return encResp("file not found"),nil
	}

	if hv,err:=uploadfile(req.Hostip,req.Filepath);err!=nil{
		return encResp("upload failed"),nil
	}else{
		msg:=fmt.Sprintf("Upload file:%s\r\nTo host:%s\r\nRename To HashCode:%s",req.Filepath,req.Hostip,hv)
		return encResp(msg),nil
	}
}

func uploadfile(hostip string,filepath string) (string,error) {

	resultchan:=make(chan error,1)

	wg:=&sync.WaitGroup{}
	piper,pipew:=io.Pipe()
	defer piper.Close()


	file, err := os.Open(filepath)
	if err != nil {
		return "",err
	}

	hv:=genrsa256hash(file)

	file.Close()

	m:=multipart.NewWriter(pipew)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer pipew.Close()

		part,err:=m.CreateFormFile("FileHash",hv)
		if err!=nil{
			resultchan <- err
			return
		}
		defer m.Close()

		file, err := os.Open(filepath)
		if err != nil {
			resultchan <- err
			return
		}
		defer file.Close()

		if _, err = io.Copy(part, file); err != nil {
			resultchan <- err
			return
		}
		resultchan <- err
		return
	}()

	posturl:="http://"+hostip+":50810/upload"


	resp,err:=http.Post(posturl, m.FormDataContentType(), piper)
	wg.Wait()

	select{
	    case err=<-resultchan:
		//todo...
	    default:
	}

	if err!=nil{
		fmt.Println(err)
	}

	if resp!=nil && resp.Body != nil{
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}

	if err!=nil {
		return "", err
	}else{
		return hv,nil
	}
}

func genrsa256hash(f *os.File) string {
	s:=sha256.New()

	if _,err:=io.Copy(s,f);err!=nil{
		return ""
	}

	hv := s.Sum(nil)

	 return "c1"+base58.Encode(hv)
}
