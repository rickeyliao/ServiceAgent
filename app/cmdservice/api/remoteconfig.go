package api

import (
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"context"
	"strings"
	"net"
	"github.com/rickeyliao/ServiceAgent/common"
	"strconv"
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	"path"
)

type RemoteConfig struct {

}


func (ccs *RemoteConfig) RemoteChange(ctx context.Context,req *pb.DefaultRequestMsg) (*pb.DefaultResp, error) {
	
	param:=strings.Split(req.Message,":")
	if len(param) != 2 {
		return encResp("config format error"),nil
	}
	if param[0] != ""{
		addr,err:=net.ResolveIPAddr("ip4",param[0])
		if err!=nil || addr.String() == "0.0.0.0"{
			return encResp("ip address format error"),nil
		}
	}

	var port int
	var err error

	port,err=strconv.Atoi(param[1])

	if param[1]==""|| err!=nil{
		return encResp("please set remote port"),nil
	}
	sac:=common.GetSAConfig()
	sac.RemoteServerPort = uint16(port)
	sac.RemoteServerIP = param[0]

	common.NewRemoteUrl1(req.Message)

	jstr,_:=json.MarshalIndent(sac,"","\t")

	sar:=common.GetSARootCfg()
	tools.Save2File(jstr,path.Join(sar.CfgDir,sar.CfgFileName))

	return encResp("success"),nil
}