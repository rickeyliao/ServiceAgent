package handle

import (
	"net/http"
	"log"
	"html/template"
)

type NotFound404 struct {

}

func (nf *NotFound404)ServeHTTP(w http.ResponseWriter, r *http.Request)   {
	
	//fmt.Println(r.URL.Path)
	
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	type myurl struct {
		MyUrl string
	}


	u:=myurl{MyUrl:r.URL.Path}


	t, err := template.ParseFiles("staticfile/html/404.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, &u)
}