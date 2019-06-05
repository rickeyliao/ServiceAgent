package main

import (
	"net/http"
	"sync"
	"fmt"
	"github.com/kprc/nbsnetwork/tools/crypt/nbscrypt"
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
	//cfg:=common.GetSARootCfg()
	//
	//fmt.Println(cfg)
	//
	//cfg.InitConfig(false)
	//
	//fmt.Println(cfg)
	//
	//fmt.Println(*cfg.SacInst)

	//os.RemoveAll("/Users/rickey/xxsa")

	//priv,_:= nbscrypt.GenerateKeyPair(1024)
	//
	//err:= nbscrypt.Save2FileRSAKey(".rsa/key",priv)
	//if err!=nil{
	//	fmt.Println(err)
	//}

	priv,pub,_:=nbscrypt.LoadRSAKey("/Users/rickey/.rsa/key")

	cipertext,_:=nbscrypt.EncryptRSA([]byte("Hello World,I'm a Goland Programer"),pub)
	fmt.Println(len(cipertext))
	fmt.Println(string(cipertext))

	plaintext,_:=nbscrypt.DecryptRsa(cipertext,priv)
	fmt.Println(string(plaintext))

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

