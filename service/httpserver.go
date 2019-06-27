package service

import (
	"context"
	"github.com/rickeyliao/ServiceAgent/agent/email"
	"github.com/rickeyliao/ServiceAgent/agent/key"
	"github.com/rickeyliao/ServiceAgent/agent/listallip"
	"github.com/rickeyliao/ServiceAgent/agent/software"
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/rickeyliao/ServiceAgent/service/file"
	"github.com/rickeyliao/ServiceAgent/service/localaddress"
	"github.com/rickeyliao/ServiceAgent/service/postsocks5"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"
)

var (
	httpserver *http.Server
)

func Run(cfg *common.SAConfig) {

	remotehost := cfg.RemoteServerIP
	remoteport := cfg.RemoteServerPort

	if remotehost == "" && remoteport == 0 {
		log.Println("Please set remote host and port")
		return
	}

	mux := http.NewServeMux()

	ips := remotehost + ":" + strconv.Itoa(int(remoteport))
	common.NewRemoteUrl1(ips)

	mux.Handle(cfg.VerifyPath, key.NewKeyAuth())
	mux.Handle(cfg.ConsumePath, key.NewKeyImport())
	mux.Handle(cfg.EmailPath, email.NewEmailRecord())
	mux.Handle(cfg.UpdateClientSoftwarePath, software.NewUpdateSoft())
	mux.Handle(cfg.TestIPAddress, localaddress.NewLocalAddress())
	mux.Handle(cfg.ListIpsPath, listallip.NewListAllIps())
	mux.Handle(cfg.PostSocks5Path, postsocks5.NewPostSocks5())
	mux.Handle(path.Join("/", cfg.UploadDir), file.NewFileUpLoad())

	listenportstr := ":" + strconv.Itoa(int(cfg.HttpListenPort))

	log.Println("Remote Server:", common.GetRemoteUrlInst().GetHostName(""))
	log.Println("Server Listen at:", listenportstr)
	log.Println("LocalNbsAddress:", cfg.NbsRsaAddr)

	httpserver = &http.Server{Addr: listenportstr, Handler: mux}

	log.Fatal(httpserver.ListenAndServe())
}

func Stop() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	httpserver.Shutdown(ctx)
}
