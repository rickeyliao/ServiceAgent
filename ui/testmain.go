package main

import (
	"net/http"
	"log"
	"github.com/rickeyliao/ServiceAgent/ui/handle"
	"fmt"
)

func main()  {

	mux:=http.NewServeMux()

	mux.Handle("/css/", http.FileServer(http.Dir("staticfile")))
	mux.Handle("/js/", http.FileServer(http.Dir("staticfile")))


	mux.Handle("/login/",&handle.LoginHandle{})
	mux.Handle("/ajax/",&handle.AjaxHandle{})
	mux.Handle("/admin/",&handle.AdminHandle{})

	mux.Handle("/",&handle.NotFound404{})

	fmt.Println("Http Server Start at: 9527")

	httpserver := &http.Server{Addr: ":9527", Handler: mux}

	log.Fatal(httpserver.ListenAndServe())
}
