package wifiap

import (
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/rickeyliao/ServiceAgent/wifiap/control"
	"log"
)

func ExtractWifiAPFiles() {

	res := "wifiap/staticfile"

	if common.WifiRes != "" {
		res = common.WifiRes
	}

	if err := control.RestoreAssets(common.GetSAConfig().GetWifiDir(), res); err != nil {
		log.Println("restore asset failed", err)
	}
}
