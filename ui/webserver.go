package ui

import (
	"context"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/rickeyliao/ServiceAgent/ui/asset"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/rickeyliao/ServiceAgent/ui/router"
)

var webserver *http.Server

func StartWebDaemon() {
	//if err:=resource.RestoreAssets("./","ui/xadmin");err!=nil{
	//	log.Println("restore asset failed",err)
	//}

	mux := http.NewServeMux()

	mux.Handle("/ajax/",&router.AjaxRouter{})


	fs := assetfs.AssetFS{Asset: asset.Asset, AssetDir: asset.AssetDir, AssetInfo: asset.AssetInfo, Prefix: "ui/xadmin"}

	mux.Handle("/", http.FileServer(&fs))


	addr := ":" + strconv.Itoa(int(common.GetSAConfig().WebServerPort))

	log.Println("Web Server Start at", addr)

	webserver = &http.Server{Addr: addr, Handler: mux}

	log.Fatal(webserver.ListenAndServe())

}

func StopWebDaemon() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	webserver.Shutdown(ctx)
}
