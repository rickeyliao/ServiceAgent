package cmd

import (
	"google.golang.org/grpc"
	"context"
	"log"
	"strconv"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"fmt"
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/kprc/nbsnetwork/tools"
)


type CmdConnection struct {
	c      *grpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

func DialToCmdService() *CmdConnection {
	var address = "127.0.0.1:" + strconv.Itoa(50811)

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

func DefaultCmdSend(cmd int32)  {
	request := &pb.DefaultRequest{}
	request.Reqid = cmd

	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewDefaultnbssasrvClient(conn.c)

	if response, err := client.DefaultNbssa(conn.ctx, request);err!=nil{
		fmt.Println(err)
	}else {
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
	if !tools.CheckPortUsed(cfg.ListenTyp, cfg.LocalListenPort) {
		log.Println("nbssa not started")
		return false
	}

	return true
}
