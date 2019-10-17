package ui

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/rickeyliao/ServiceAgent/ui/asset"
	"net/http"
	"log"
	"time"
	"context"
	"strconv"
	"github.com/rickeyliao/ServiceAgent/common"
)

var webserver *http.Server

func StartWebDaemon()  {
	//if err:=resource.RestoreAssets("./","ui/xadmin");err!=nil{
	//	log.Println("restore asset failed",err)
	//}

	mux:=http.NewServeMux()

	fs:=assetfs.AssetFS{Asset:asset.Asset,AssetDir:asset.AssetDir,AssetInfo:asset.AssetInfo,Prefix:"ui/xadmin"}

	mux.Handle("/",http.FileServer(&fs))

	addr:=":"+strconv.Itoa(int(common.GetSAConfig().WebServerPort))

	log.Println("Web Server Start at",addr)

	webserver = &http.Server{Addr:addr,Handler:mux}

	log.Fatal(webserver.ListenAndServe())

}

func StopWebDaemon()  {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	webserver.Shutdown(ctx)
}