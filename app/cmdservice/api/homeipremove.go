package api


import (
	"context"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/common"
	"net"
	"github.com/rickeyliao/ServiceAgent/service/localaddress"
)

type HomeIPRemoveSrv struct {

}

func (hirs *HomeIPRemoveSrv)RemoveHomeIP(ctx context.Context, req *pb.HomeIPRemoveReq) (*pb.DefaultResp, error)  {
	if req.Ipaddr == "" && req.Nbsaddr == ""{
		return encResp("nbsaddr or ip address must set"),nil
	}

	if req.Nbsaddr != ""{
		if !common.CheckNbsNodeHash(req.Nbsaddr){
			return encResp("nbsaddr error"),nil
		}
	}

	if req.Nbsaddr == "" && req.Ipaddr != ""{
		if ip:=net.ParseIP(req.Ipaddr);ip==nil{
			return encResp("ip address error"),nil
		}
	}

	var hid *localaddress.Homeipdesc
	var addr string

	if req.Ipaddr != ""{
		hidb:=localaddress.GetHomeIPDB()

		hid,addr=hidb.FindByIP(req.Ipaddr)

		if hid == nil{
			return encResp("not found node by ip address"),nil
		}

	}

	//begin to delete




	return encResp("Delete Successfully!")


}