package api

import (
	"golang.org/x/net/context"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"time"
)

type StopServer struct {
	Stop func()
}

func (ss *StopServer)Stopnbssa(ctx context.Context, in *pb.StopRequest) (*pb.DefaultResp,error)  {
	 go func() {
	 	time.Sleep(time.Second*2)
		ss.Stop()
	 }()

	 resp := &pb.DefaultResp{}
	 resp.Message = "nbssa server stoped"
	 return resp,nil
}


