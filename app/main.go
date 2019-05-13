package app

import (
	"net/http"
	"log"
	"os"
	"strconv"
	"github.com/rickeyliao/ServiceAgent/key"
	"github.com/rickeyliao/ServiceAgent/email"
	"github.com/rickeyliao/ServiceAgent/software"
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/rickeyliao/ServiceAgent/localaddress"
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
	http.Handle("/localipaddress",localaddress.NewLocalAddress())

	var localport uint16

	if localip == ""{
		localport = 50810
	}else {
		localport = common.GetPort(localip)
	}

	listenportstr := ":"+strconv.Itoa(int(localport))

	log.Println("Remote Server:",common.GetRemoteUrlInst().GetHostName(""))
	log.Println("Server Listen at:",localport)

	log.Fatal(http.ListenAndServe(listenportstr, nil))
}