package main

import (
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/rickeyliao/ServiceAgent/service"
	"github.com/rickeyliao/ServiceAgent/app/cmd"
)

func main()  {

	if !cmd.CheckProcessCanStarted(){
		return
	}

	cfg:=common.GetSAConfig()
	service.Run(cfg)
}

