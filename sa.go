package main

import (
	"net/http"
	"io"
	"log"
	"sync"
	"fmt"
)




type countHandler struct {
	mu sync.Mutex // guards n
	n  int
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	fmt.Fprintf(w, "count is %d\n", h.n)
}

func main() {
	http.Handle("/count", new(countHandler))
	log.Fatal(http.ListenAndServe(":33221", nil))
}


//func main(){
//	http.HandleFunc("/hello",HelloServer)
//	//http.Handle("/foo",nil)
//
//
//	err:=http.ListenAndServe(":33112",nil)
//	if err!=nil{
//		log.Fatal("Listen and service error,",err)
//	}
//
//}

func HelloServer(w http.ResponseWriter,req *http.Request)  {
	io.WriteString(w,"hello, world\n")
}


