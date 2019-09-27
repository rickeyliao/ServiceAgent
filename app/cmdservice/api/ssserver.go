package api

import (
	"context"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/app"
	"github.com/rickeyliao/ServiceAgent/service/localaddress"
	"fmt"
	"strconv"
)

type CmdSSServer struct {

}

func (css *CmdSSServer)SSServerDo(ctx context.Context,req *pb.SSServerReq) (*pb.DefaultResp, error)  {
	switch req.Op {
	case app.CMD_SSSERVER_SHOW:
		return css.show(req)
	case app.CMD_SSSERVER_UPDATE:
		return css.update(req)
	}

	return encResp("Command Not found"),nil
}

func (css *CmdSSServer)show(req *pb.SSServerReq) (*pb.DefaultResp, error) {
	remotessl:=localaddress.GetServerList()

	if len(remotessl) == 0{
		return encResp("No Server List"),nil
	}

	message:=""

	for _,ssl:=range remotessl{
		if ssl.DeleteFlag {
			continue
		}
		if !((req.Nationality == 1 && ssl.Abroad == 1) || (req.Nationality == 0) || (req.Nationality == 86 && ssl.Abroad == 0)) {
			continue
		}

		if len(message) >0{
			message += "\r\n"
		}
		//message += ssl.NodeId
		message += fmt.Sprintf("%-45s",ssl.NodeId)
		//message += ""
		message += fmt.Sprintf("%-16s",ssl.Name)
		//message += "\t"
		message += fmt.Sprintf("%-18s",ssl.IPAddress)
		//message += "\t"
		message += fmt.Sprintf("%-8s",strconv.Itoa(ssl.SSPort))
		//message += "\t"
		message += fmt.Sprintf("%-16s",ssl.SSPassword)
		//message += "\t"
		message += fmt.Sprintf("%-8s",getNodeStatus(ssl.Status))
		//message += "\t"
		message += fmt.Sprintf("%-18s",ssl.Location)
		//message += "\t"
		message += fmt.Sprintf("%-6s",getNodeNationality(ssl.Abroad))
		message += fmt.Sprintf("%-20s",ssl.LastModify.Format("2006-01-02 15:04:05"))
	}
	if message == ""{
		message = "No Server List"
	}

	return encResp(message),nil
}

func getNodeNationality(abroad int) string {
	if abroad == 0{
		return "ML"
	}

	if abroad == 1{
		return "A"
	}

	return ""

}

func getNodeStatus(status int) string {
	if status == 0{
		return "idle   "
	}

	if status == 1{
		return "working"
	}

	return "unknow"
}


func (css *CmdSSServer)update(req *pb.SSServerReq)(*pb.DefaultResp, error) {
	return encResp(""),nil
}