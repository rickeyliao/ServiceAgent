package api

import (
	"context"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"net"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/rickeyliao/ServiceAgent/common"
	"fmt"
	"strings"
	"strconv"
)

type CmdFileDownload struct {

}

func (cfd *CmdFileDownload)Downloadfile(ctx context.Context, req *pb.Filedownloadreq) (*pb.DefaultResp, error)  {
	if req.Hostip == "" || req.Filehash == "" || req.Savepath==""{
		return encResp("Parameter error"),nil
	}

	hiparr:=strings.Split(req.Hostip,":")
	if len(hiparr)!=2{

		return encResp("host ip address error"),nil
	}

	if net.ParseIP(hiparr[0]) == nil{

		return encResp("host ip address error"),nil
	}

	if rport,err:=strconv.Atoi(hiparr[1]);err!=nil{

		return encResp("host ip address error"),nil
	}else {
		if rport<1024 || rport >65535{

			return encResp("host ip address error"),nil
		}
	}

	if !tools.FileExists(req.Savepath){
		return encResp("Save Path not found"),nil
	}

	if !common.CheckNbsCotentHash(req.Filehash){
		return encResp("File Hash error"),nil
	}

	if err:=common.DownloadFile(req.Hostip,req.Savepath,req.Filehash);err!=nil{
		msg:="Download failed\r\n"+err.Error()
		return encResp(msg),nil
	}

	message:=fmt.Sprintf("Success!!!\r\nFile: %s\r\nDownload From: %s\r\nSave to: %s",
		req.Filehash,req.Hostip,req.Savepath)
	return encResp(message),nil
}

