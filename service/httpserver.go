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
	"github.com/rickeyliao/ServiceAgent/service/login"
	"sync"
)

var (
	httpserver *http.Server
	quit chan int
	wg sync.WaitGroup
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
	mux.Handle(cfg.TestIPAddressPath, localaddress.NewLocalAddress())
	mux.Handle(cfg.ListIpsPath, listallip.NewListAllIps())
	mux.Handle(cfg.PostSocks5Path, postsocks5.NewPostSocks5())
	mux.Handle(cfg.Uploadpath, file.NewFileUpLoad())
	mux.Handle(cfg.DownloadPath, file.NewFileDownLoad())
	mux.Handle(path.Join("/",cfg.LoginDir),login.NewLoginInfo())

	listenportstr := ":" + strconv.Itoa(int(cfg.HttpListenPort))

	log.Println("Remote Server:", common.GetRemoteUrlInst().GetHostName(""))
	log.Println("Server Listen at:", listenportstr)
	log.Println("LocalNbsAddress:", cfg.NbsRsaAddr)

	wg.Add(1)
	go reportAddress()


	httpserver = &http.Server{Addr: listenportstr, Handler: mux}

	log.Fatal(httpserver.ListenAndServe())
}

func Stop() {
	quit<-1
	wg.Wait()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	httpserver.Shutdown(ctx)
}

func reportAddress()  {
	var count int64

	for {
		count++
		if count %100 == 0{
			c:=&http.Client{}

			if req,err:=http.NewRequest("GET","http://103.45.98.72:50810/localipaddress",nil);err!=nil{
				//fmt.Println(err)
				continue
			}else{

				req.Header.Add("nbsaddress",common.GetSAConfig().NbsRsaAddr)
				//req.Header.Add("FileName","foo.txt")

				if _,errresp:=c.Do(req);errresp != nil{
					log.Println(errresp)
					continue
				}else {
					//log.Println(resp)
				}



			}
		}
		time.Sleep(time.Second*1)
		select {
		case <-quit:
			wg.Done()
			return
		default:
			//todo...
		}
	}
}