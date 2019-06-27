package cmd

import (
	"context"
	"fmt"
	"github.com/kprc/nbsnetwork/tools"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/common"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

type CmdConnection struct {
	c      *grpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

func DialToCmdService() *CmdConnection {
	var address = "127.0.0.1:" + strconv.Itoa(int(common.GetSAConfig().CmdListenPort))

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("can not connect rpc server:", err)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &CmdConnection{
		c:      conn,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (conn *CmdConnection) Close() {
	conn.c.Close()
	conn.cancel()
}

func DefaultCmdSend(cmd int32) {
	request := &pb.DefaultRequest{}
	request.Reqid = cmd

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewDefaultnbssasrvClient(conn.c)

	if response, err := client.DefaultNbssa(conn.ctx, request); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.Message)
	}

}

func DefaultCmdSendStr(msg string) {
	request := &pb.DefaultRequestMsg{}
	request.Message = msg

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewConfigchangeClient(conn.c)

	if response, err := client.ChangeConfig(conn.ctx, request); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.Message)
	}
}

func RemoteCmdSendStr(msg string) {
	request := &pb.DefaultRequestMsg{}
	request.Message = msg

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewRemotechangeClient(conn.c)

	if response, err := client.RemoteChange(conn.ctx, request); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.Message)
	}

}

func CheckProcessReady() bool {
	sar := common.GetSARootCfg()
	if !sar.IsInitialized() {
		log.Println("Please Initialize First")
		return false
	}
	//load config
	sar.LoadCfg()
	sar.LoadRsaKey()
	cfg := sar.SacInst
	//if the program started, quit
	if !tools.CheckPortUsed(cfg.ListenTyp, cfg.CmdListenPort) {
		log.Println("nbssa not started")
		return false
	}

	return true
}

func CheckProcessCanStarted() bool {
	sar := common.GetSARootCfg()
	if !sar.IsInitialized() {
		log.Println("Please Initialize First")
		return false
	}
	//load config
	sar.LoadCfg()
	sar.LoadRsaKey()
	cfg := sar.SacInst
	//if the program started, quit
	if tools.CheckPortUsed(cfg.ListenTyp, cfg.CmdListenPort) {
		log.Println("nbssa have started")
		return false
	}

	return true
}

func BootstrapCmdSend(op bool, req string) {
	request := &pb.BootstrapCHGReq{Op: op, Address: req}
	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewBootstrapCHGClient(conn.c)

	if response, err := client.ChangeBootstrap(conn.ctx, request); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.Message)
	}

}
