package main

import (
	"net/http"
	"log"
	"github.com/rickeyliao/ServiceAgent/key"
	"github.com/rickeyliao/ServiceAgent/email"
	"github.com/rickeyliao/ServiceAgent/software"
	"os"
	"github.com/rickeyliao/ServiceAgent/common"
	"strconv"
)

func main()  {
	var remotehost,remoteport,localip string

	if len(os.Args) > 1{
		remotehost = os.Args[1]
		if len(os.Args)>2{
			remoteport = os.Args[2]
			if len(os.Args)>3{
				localip = os.Args[3]
			}
		}
	}

	common.NewRemoteUrl(remotehost,remoteport)


	http.Handle("/public/keys/verify", key.NewKeyAuth())
	http.Handle("/public/keys/consume",key.NewKeyImport())
	http.Handle("/public/key/refresh",email.NewEmailRecord())
	http.Handle("/public/app",software.NewUpdateSoft())

	var localport uint16

	if localip == ""{
		localport = 9527
	}else {
		localport = common.GetPort(localip)
	}

	listenportstr := ":"+strconv.Itoa(int(localport))

	log.Println("Remote Server:",common.GetRemoteUrlInst().GetHostName(""))
	log.Println("Server Listen at:",localport)

	log.Fatal(http.ListenAndServe(listenportstr, nil))
}
