package main

import (
	"net/http"
	"sync"
	"fmt"
	"github.com/rickeyliao/ServiceAgent/common"
)


type countHandler struct {
	mu sync.Mutex // guards n
	n  int
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	fmt.Println(r.URL.Path,r.Method)
	if r.Method != "POST"{
		fmt.Fprintf(w,"terst")
		return
	}
	h.n++
	fmt.Fprintf(w, "count is %d\n", h.n)

}

//func main() {
//	http.Handle("/count", new(countHandler))
//	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
//
//		body,_:=ioutil.ReadAll(request.Body)
//
//		res,code,err:=common.Post("http://39.98.40.7:8078/public/keys/consume",string(body))
//
//
//		fmt.Println(res)
//		fmt.Println(code)
//		fmt.Println(err)
//
//	})
//	log.Fatal(http.ListenAndServe(":33221", nil))
//}

func main()  {
	cfg:=common.GetSARootCfg()

	fmt.Println(cfg)


	cfg.InitConfig()

	fmt.Println(cfg)

	fmt.Println(*cfg.SacInst)

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
//
//func HelloServer(w http.ResponseWriter,req *http.Request)  {
//	io.WriteString(w,"hello, world\n")
//}


