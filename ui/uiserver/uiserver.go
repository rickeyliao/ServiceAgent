package uiserver

import (
	"net/http"
	"github.com/rickeyliao/ServiceAgent/ui/handle"
	"fmt"
	"log"
	"time"
	"context"
)

var uiserver *http.Server

func StartUIServer() {
	mux:=http.NewServeMux()

	mux.Handle("/css/", http.FileServer(http.Dir("staticfile")))
	mux.Handle("/js/", http.FileServer(http.Dir("staticfile")))


	mux.Handle("/login/",&handle.LoginHandle{})
	mux.Handle("/ajax/",&handle.AjaxHandle{})
	mux.Handle("/admin/",&handle.AdminHandle{})

	mux.Handle("/",&handle.NotFound404{})

	fmt.Println("Http Server Start at: 9527")

	uiserver = &http.Server{Addr: ":9527", Handler: mux}

	log.Fatal(uiserver.ListenAndServe())

}

func StopUIServer()  {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	uiserver.Shutdown(ctx)
}





