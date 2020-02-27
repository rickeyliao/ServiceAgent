package api

import (
	"context"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/app"
	"github.com/rickeyliao/ServiceAgent/dht2"
)

type DhtCmdSrv struct {

}

func (dcs *DhtCmdSrv)DhtCmdDo(ctx context.Context,cmd *pb.DhtCmdPb) (*pb.DefaultResp, error)  {
	switch cmd.Op {
	case app.CMD_DHT_ONLINE:
		dht2.Online()
	    return encResp("Node is start Online"),nil
	case app.CMD_DHT_START:
	case app.CMD_DHT_RESTART:
	case app.CMD_DHT_STOP:
	default:
		return encResp("command line not found"),nil
	}
}