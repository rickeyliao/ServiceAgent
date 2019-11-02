package wifiap

import (
	"log"
	"github.com/rickeyliao/ServiceAgent/wifiap/control"
	"github.com/rickeyliao/ServiceAgent/common"
)

func ExtractWifiAPFiles()  {
	if err:=control.RestoreAssets(common.GetSAConfig().GetWifiDir(),"wifiap/staticfile");err!=nil{
		log.Println("restore asset failed",err)
	}
}

