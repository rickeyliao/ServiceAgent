package api

import (
	"context"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"net"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/rickeyliao/ServiceAgent/common"
	"net/http"
	"os"
	"io"
	"path"
	"fmt"
	"github.com/pkg/errors"
)

type CmdFileDownload struct {

}

func (cfd *CmdFileDownload)Downloadfile(ctx context.Context, req *pb.Filedownloadreq) (*pb.DefaultResp, error)  {
	if req.Hostip == "" || req.Filehash == "" || req.Savepath==""{
		return encResp("Parameter error"),nil
	}

	if net.ParseIP(req.Hostip) == nil{
		return encResp("host ip address error"),nil
	}

	if !tools.FileExists(req.Savepath){
		return encResp("Save Path not found"),nil
	}

	if !common.CheckNbsCotentHash(req.Filehash){
		return encResp("File Hash error"),nil
	}

	if err:=downloadFile(req.Hostip,req.Savepath,req.Filehash);err!=nil{
		msg:="Download failed\r\n"+err.Error()
		return encResp(msg),nil
	}

	message:=fmt.Sprintf("Success!!!\r\nFile: %s\r\nDownload From: %s\r\nSave to: %s",
		req.Filehash,req.Hostip,req.Savepath)
	return encResp(message),nil
}

func downloadFile(hostip string,savepath string,filehash string) error  {
	tp:=http.Transport{DisableKeepAlives:true}
	c:=&http.Client{Transport:&tp}

	geturl:="http://"+hostip+":50810/download"

	if req,err:=http.NewRequest("GET",geturl,nil);err!=nil{

		return err
	}else{

		req.Header.Add("FileHash",filehash)


		if resp,errresp:=c.Do(req);errresp != nil{

			return errresp
		}else{
			defer resp.Body.Close()
			message:=resp.Header.Get("message")
			if message == "FileNotFound"{
				return errors.New("File Not Found")
			}
			f, err := os.OpenFile(path.Join(savepath,filehash), os.O_WRONLY|os.O_CREATE, 0755)
			if err != nil {

				return err
			}
			defer f.Close()
			io.Copy(f, resp.Body)

		}

		return nil

	}
}