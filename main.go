package main

import (
	"net/http"
	"log"
	"github.com/rickeyliao/ServiceAgent/key"
	"github.com/rickeyliao/ServiceAgent/email"
	"github.com/rickeyliao/ServiceAgent/software"
	"os"
	"github.com/rickeyliao/ServiceAgent/common"
	"fmt"
)

func main()  {
	var host,port string

	if len(os.Args) > 1{
		host = os.Args[1]
		if len(os.Args)>2{
			port = os.Args[2]
		}
	}

	common.NewRemoteUrl(host,port)
	fmt.Println("Remote Server:",common.GetRemoteUrlInst().GetHostName(""))

	http.Handle("/public/keys/verify", key.NewKeyAuth())
	http.Handle("/public/keys/consume",key.NewKeyImport())
	http.Handle("/public/key/refresh",email.NewEmailRecord())
	http.Handle("/public/app",software.NewUpdateSoft())

	log.Fatal(http.ListenAndServe(":9527", nil))
}
