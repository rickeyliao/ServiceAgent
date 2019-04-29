package main

import (
	"net/http"
	"log"
	"github.com/rickeyliao/ServiceAgent/key"
	"github.com/rickeyliao/ServiceAgent/email"
	"github.com/rickeyliao/ServiceAgent/software"
)

func main()  {
	http.Handle("/public/keys/verify", key.NewKeyAuth())
	http.Handle("/public/keys/consume",key.NewKeyImport())
	http.Handle("/public/key/refresh",email.NewEmailRecord())
	http.Handle("/public/app",software.NewUpdateSoft())

	log.Fatal(http.ListenAndServe(":80", nil))
}
