package api

import (
	"golang.org/x/net/context"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"time"
	"github.com/rickeyliao/ServiceAgent/app"
	"github.com/rickeyliao/ServiceAgent/common"
	"encoding/json"
)

type CmdDefaultServer struct {
	Stop func()
}

func (ss *CmdDefaultServer)DefaultNbssa(ctx context.Context, in *pb.DefaultRequest) (*pb.DefaultResp,error)  {
	if in.Reqid == app.CMD_STOP_REQ {
		return ss.stop()
	}

	if in.Reqid == app.CMD_CONFIG_SHOW_REQ{
		return ss.configshow()
	}
	resp := &pb.DefaultResp{}
	resp.Message = "no cmd found"

	return resp,nil
}

func (ss *CmdDefaultServer)stop() (*pb.DefaultResp,error)  {
	go func() {
		time.Sleep(time.Second*2)
		ss.Stop()
	}()

	resp := &pb.DefaultResp{}
	resp.Message = "nbssa server stoped"
	return resp,nil
}

func (ss *CmdDefaultServer)configshow() (*pb.DefaultResp,error) {
	sac:=common.GetSAConfig()

	resp:=&pb.DefaultResp{}

	if j,err:=json.MarshalIndent(sac,"","\t");err!=nil{
		resp.Message = "Marshal json failed"
	}else{
		resp.Message = string(j)
	}

	return resp,nil

}
