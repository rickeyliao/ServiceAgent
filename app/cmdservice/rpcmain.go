package cmdservice

import (
	"fmt"
	"github.com/rickeyliao/ServiceAgent/app/cmdservice/api"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/rickeyliao/ServiceAgent/service"
	"github.com/rickeyliao/ServiceAgent/service/shadowsock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var (
	grpcServer *grpc.Server
	stopflag bool
)



func StartCmdService() {

	lis, err := net.Listen("tcp",
		fmt.Sprintf("%s:%d", common.GetSAConfig().CmdListenIP, common.GetSAConfig().CmdListenPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer = grpc.NewServer()

	pb.RegisterDefaultnbssasrvServer(grpcServer, &api.CmdDefaultServer{StopCmdService})
	pb.RegisterConfigchangeServer(grpcServer, &api.CmdConfigServer{})
	pb.RegisterRemotechangeServer(grpcServer, &api.RemoteConfig{})
	pb.RegisterBootstrapCHGServer(grpcServer, &api.CmdBootstrapServer{})
	pb.RegisterLicenseUserChgServer(grpcServer, &api.CmdLicenseUserServer{})
	pb.RegisterHomeIPShowSrvServer(grpcServer, &api.CmdHomeShow{})
	pb.RegisterLicenseSrvServer(grpcServer, &api.CmdOpLicenseSrv{})
	pb.RegisterFileuploadsrvServer(grpcServer, &api.CmdFileUpLoad{})
	pb.RegisterFileudownloadsrvServer(grpcServer, &api.CmdFileDownload{})
	pb.RegisterSSServerServiceServer(grpcServer, &api.CmdSSServer{})
	pb.RegisterHomeIPRemoveSrvServer(grpcServer, &api.HomeIPRemoveSrv{})
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func StopCmdService() {

	if stopflag {
		return
	}

	stopflag = true

	service.Stop()

	if grpcServer == nil {
		return
	}

	shadowsock.StopSS2Server()

	grpcServer.Stop()
	log.Println("Cmd Server Closed")

}
