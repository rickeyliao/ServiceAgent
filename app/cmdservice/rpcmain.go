package cmdservice

import (
	"net"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"github.com/rickeyliao/ServiceAgent/app/cmdservice/api"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"google.golang.org/grpc/reflection"
	"github.com/rickeyliao/ServiceAgent/service"
	"github.com/rickeyliao/ServiceAgent/common"
)

var grpcServer *grpc.Server

func StartCmdService()  {

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", common.GetSAConfig().CmdListenIP,common.GetSAConfig().CmdListenPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer = grpc.NewServer()

	pb.RegisterDefaultnbssasrvServer(grpcServer,&api.CmdDefaultServer{StopCmdService})
	pb.RegisterConfigchangeServer(grpcServer,&api.CmdConfigServer{})
	pb.RegisterRemotechangeServer(grpcServer,&api.RemoteConfig{})
	pb.RegisterBootstrapCHGServer(grpcServer,&api.CmdBootstrapServer{})
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}


func StopCmdService()  {

	service.Stop()

	if grpcServer == nil{
		return
	}

	grpcServer.Stop()
	log.Println("Cmd Server Closed")

}
