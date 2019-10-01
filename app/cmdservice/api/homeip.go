package api

import (
	"context"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/service/localaddress"
)

type CmdHomeShow struct {
}

func (chs *CmdHomeShow) ShowHomeIP(cxt context.Context, hi *pb.HomeIPShowReq) (*pb.DefaultResp, error) {



	if hi.Nbsaddr == "" {
		message := localaddress.CmdShowAddressAll(0)
		if message == "" {
			message = "No home ip"
		}
		return encResp(message), nil
	} else {
		message := localaddress.CmdShowAddress(hi.Nbsaddr)
		if message == "" {
			return encResp("not found"), nil
		} else {
			return encResp(message), nil
		}
	}
}
