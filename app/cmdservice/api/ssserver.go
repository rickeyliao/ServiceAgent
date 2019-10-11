package api

import (
	"bytes"
	"context"
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/rickeyliao/ServiceAgent/app"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/service/localaddress"
)

type CmdSSServer struct {
}

func (css *CmdSSServer) SSServerDo(ctx context.Context, req *pb.SSServerReq) (*pb.DefaultResp, error) {
	switch req.Op {
	case app.CMD_SSSERVER_SHOW:
		return css.show(req)
	case app.CMD_SSSERVER_UPDATE:
		return css.update(req)
	case app.CMD_SSSERVER_REMOVE:
		return css.remove(req)

	}

	return encResp("Command Not found"), nil
}

func sortlist(srvl []localaddress.SSServerListNode) list.List {
	l := list.NewList(func(v1 interface{}, v2 interface{}) int {
		s1, s2 := v1.(localaddress.SSServerListNode), v2.(localaddress.SSServerListNode)
		return bytes.Compare([]byte(s1.IPAddress), []byte(s2.IPAddress))
	})
	l.SetSortFunc(func(v1 interface{}, v2 interface{}) int {
		s1, s2 := v1.(localaddress.SSServerListNode), v2.(localaddress.SSServerListNode)
		if s1.Abroad >= s2.Abroad {
			return 1
		} else {
			return -1
		}
	})

	for _, ssl := range srvl {
		l.AddValueOrder(ssl)
	}

	return l
}

func (css *CmdSSServer) show(req *pb.SSServerReq) (*pb.DefaultResp, error) {
	if !req.Local {
		return css.showremote(req)
	} else {
		return css.showlocal(req)
	}
}

func (css *CmdSSServer) showlocal(req *pb.SSServerReq) (*pb.DefaultResp, error) {

	message := localaddress.CmdShowAddressAll(req.Nationality)

	if message == "" {
		message = "No Server List"
	}

	return encResp(message), nil
}
func (css *CmdSSServer) showremote(req *pb.SSServerReq) (*pb.DefaultResp, error) {

	remotessl := localaddress.GetServerList()

	if len(remotessl) == 0 {
		return encResp("No Server List"), nil
	}

	message := ""

	l := sortlist(remotessl)
	cursor := l.ListIterator(0)

	for {

		v := cursor.Next()
		if v == nil {
			break
		}

		ssl := v.(localaddress.SSServerListNode)

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
		message += ssl.String()
	}

	if message == "" {
		message = "No Server List"
	}

	return encResp(message), nil
}

func (css *CmdSSServer) update(req *pb.SSServerReq) (*pb.DefaultResp, error) {
	message := localaddress.UpdateServer(req.Nationality, req.Ip, req.Nbsaddr)
	if message == "" {
		message = "Nothing to do..."
	}
	return encResp(message), nil
}

func (css *CmdSSServer) remove(req *pb.SSServerReq) (*pb.DefaultResp, error) {
	msg := ""
	if req.Ip != "" {
		msg = localaddress.CmdDeleteServerByIP(req.Ip)
	} else {
		msg = localaddress.CmdDeleteServer(req.Nationality)
	}

	return encResp(msg), nil
}
