package wifiap

import (
	"log"
	"github.com/rickeyliao/ServiceAgent/wifiap/control"
	"github.com/rickeyliao/ServiceAgent/common"
)

func ExtractWifiAPFiles()  {

	res:="wifiap/staticfile"

	if common.WifiRes != ""{
		res = common.WifiRes
	}

	if err:=control.RestoreAssets(common.GetSAConfig().GetWifiDir(),res);err!=nil{
		log.Println("restore asset failed",err)
	}
}

