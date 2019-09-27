package api

import (
	"context"
	"fmt"
	"github.com/rickeyliao/ServiceAgent/app"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/service/localaddress"
	"strconv"
)

type CmdSSServer struct {
}

func (css *CmdSSServer) SSServerDo(ctx context.Context, req *pb.SSServerReq) (*pb.DefaultResp, error) {
	switch req.Op {
	case app.CMD_SSSERVER_SHOW:
		return css.show(req)
	case app.CMD_SSSERVER_UPDATE:
		return css.update(req)
	}

	return encResp("Command Not found"), nil
}

func (css *CmdSSServer) show(req *pb.SSServerReq) (*pb.DefaultResp, error) {
	remotessl := localaddress.GetServerList()

	if len(remotessl) == 0 {
		return encResp("No Server List"), nil
	}

	message := ""

	for _, ssl := range remotessl {
		if ssl.DeleteFlag {
			continue
		}
		if !(((req.Nationality == app.NATIONALITY_AMERICAN ||
			req.Nationality == app.NATIONALITY_JAPANESE ||
			req.Nationality == app.NATIONALITY_SINGAPORE ||
			req.Nationality == app.NATIONALITY_ENGLAND) && ssl.Abroad == app.ABROAD_AMERICAN) ||
			(req.Nationality == 0) ||
			(req.Nationality == app.NATIONALITY_CHINA_MAINLAND && ssl.Abroad == app.ABROAD_CHINA_MAINLAND)) {
			continue
		}

		if len(message) > 0 {
			message += "\r\n"
		}
		message += fmt.Sprintf("%-45s", ssl.NodeId)
		message += fmt.Sprintf("%-16s", ssl.Name)
		message += fmt.Sprintf("%-18s", ssl.IPAddress)
		message += fmt.Sprintf("%-8s", strconv.Itoa(ssl.SSPort))
		message += fmt.Sprintf("%-16s", ssl.SSPassword)
		message += fmt.Sprintf("%-8s", getNodeStatus(ssl.Status))
		message += fmt.Sprintf("%-18s", ssl.Location)
		message += fmt.Sprintf("%-6s", getNodeNationality(ssl.Abroad))
		message += fmt.Sprintf("%-20s", ssl.LastModify.Format("2006-01-02 15:04:05"))
	}
	if message == "" {
		message = "No Server List"
	}

	return encResp(message), nil
}

func getNodeNationality(abroad int) string {
	if abroad == 0 {
		return "ML"
	}

	if abroad == 1 {
		return "A"
	}

	return ""
}

func getNodeStatus(status int) string {
	if status == 0 {
		return "idle   "
	}

	if status == 1 {
		return "working"
	}

	return "unknow"
}

func (css *CmdSSServer) update(req *pb.SSServerReq) (*pb.DefaultResp, error) {
	return encResp(""), nil
}
