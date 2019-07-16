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

	if !CheckProcessReady(){
		return
	}

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

func LicenseUserCmdSend(op bool, req string) {
	request := &pb.LicenseUserChgReq{Op: op, User: req}
	conn := DialToCmdService()
	defer conn.Close()

	client := pb.NewLicenseUserChgClient(conn.c)

	if response, err := client.ChgLicenseUser(conn.ctx, request); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.Message)
	}

}

func HomeIPShowCmdSend(req string)  {
	request:=&pb.HomeIPShowReq{Nbsaddr:req}
	conn:=DialToCmdService()
	defer conn.Close()

	client:=pb.NewHomeIPShowSrvClient(conn.c)

	if resp,err:=client.ShowHomeIP(conn.ctx,request);err!=nil{
		fmt.Println(err)
	}else {
		fmt.Println(resp.Message)
	}

}

func LicenseCmdSend(op int32,req string)  {
	request:=&pb.LicenseReq{Op:op,Sofaaddress:req}
	conn:=DialToCmdService()
	defer conn.Close()

	client:=pb.NewLicenseSrvClient(conn.c)
	if resp,err:=client.OpLicense(conn.ctx,request);err!=nil{
		fmt.Println(err)
	}else {
		fmt.Println(resp.Message)
	}

}

func UploadFileCmdSend(hostip string,filepath string)  {
	request:=&pb.Fileuploadreq{Hostip:hostip,Filepath:filepath}
	conn:=DialToCmdService()
	defer conn.Close()

	client:=pb.NewFileuploadsrvClient(conn.c)
	if resp,err:=client.Uploadfile(conn.ctx,request);err!=nil{
		fmt.Println(err)
	}else {
		fmt.Println(resp.Message)
	}

}

func DownloadFileCmdSend(hostip string,filename string,filesavepath string)  {
	request:=&pb.Filedownloadreq{Hostip:hostip,Filehash:filename,Savepath:filesavepath}
	conn:=DialToCmdService()
	defer conn.Close()

	client:=pb.NewFileudownloadsrvClient(conn.c)
	if resp,err:=client.Downloadfile(conn.ctx,request);err!=nil{
		fmt.Println(err)
	}else {
		fmt.Println(resp.Message)
	}
}