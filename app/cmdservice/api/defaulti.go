package api

import (
	"encoding/json"
	"github.com/rickeyliao/ServiceAgent/app"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/common"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

type CmdDefaultServer struct {
	Stop func()
}

func (ss *CmdDefaultServer) DefaultNbssa(ctx context.Context, in *pb.DefaultRequest) (*pb.DefaultResp, error) {
	if in.Reqid == app.CMD_STOP_REQ {
		return ss.stop()
	}

	if in.Reqid == app.CMD_CONFIG_SHOW_REQ {
		return ss.configshow()
	}

	if in.Reqid == app.CMD_REMOTE_SHOW_REQ {
		return ss.remoteshow()
	}

	if in.Reqid == app.CMD_BOOTSTRAP_SHOW_REQ {
		return ss.bootstrapshow()
	}

	resp := &pb.DefaultResp{}
	resp.Message = "no cmd found"

	return resp, nil
}

func (ss *CmdDefaultServer) stop() (*pb.DefaultResp, error) {
	go func() {
		time.Sleep(time.Second * 2)
		ss.Stop()
	}()

	resp := &pb.DefaultResp{}
	resp.Message = "nbssa server stoped"
	return resp, nil
}

func (ss *CmdDefaultServer) configshow() (*pb.DefaultResp, error) {
	sac := common.GetSAConfig()

	resp := &pb.DefaultResp{}

	if j, err := json.MarshalIndent(sac, "", "\t"); err != nil {
		resp.Message = "Marshal json failed"
	} else {
		resp.Message = string(j)
	}

	return resp, nil

}

func (ss *CmdDefaultServer) remoteshow() (*pb.DefaultResp, error) {
	sac := common.GetSAConfig()

	resp := &pb.DefaultResp{}

	resp.Message = sac.RemoteServerIP + ":" + strconv.Itoa(int(sac.RemoteServerPort))

	return resp, nil
}

func (ss *CmdDefaultServer) bootstrapshow() (*pb.DefaultResp, error) {
	sac := common.GetSAConfig()
	resp := &pb.DefaultResp{}

	message := ""
	for _, v := range sac.BootstrapIPAddress {
		if message != "" {
			message += "\r\n"
		}
		message += v

	}

	resp.Message = message

	return resp, nil
}
