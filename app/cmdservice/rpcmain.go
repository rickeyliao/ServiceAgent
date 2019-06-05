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
)

var grpcServer *grpc.Server

func StartCmdService()  {

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", 50811))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer = grpc.NewServer()

	pb.RegisterStopnbssaServer(grpcServer,&api.StopServer{StopCmdService})
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
