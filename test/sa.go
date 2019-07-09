package main

import (
	"net/http"
	"sync"
	"fmt"
	"crypto/sha1"
	"gx/ipfs/QmWFAMPqsEyUX7gDUsRVmMWz59FxSpJ1b2v6bJ1yYzo7jY/go-base58-fast/base58"
	"time"
	"os"
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

	hn,_:=os.Hostname()

	fmt.Println(hn)

	return

	testlist()

	return

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

	//priv,pub,_:=nbscrypt.LoadRSAKey("/Users/rickey/.rsa/key")
	//
	//cipertext,_:=nbscrypt.EncryptRSA([]byte("Hello World,I'm a Goland Programer"),pub)
	//fmt.Println(len(cipertext))
	//fmt.Println(string(cipertext))
	//
	//plaintext,_:=nbscrypt.DecryptRsa(cipertext,priv)
	//fmt.Println(string(plaintext))

	s:=sha1.New()
	s.Write([]byte("hello world"))
	result := s.Sum(nil)

	fmt.Println(len(result),base58.Encode(result))



	var i int
	var j int

	wg :=&sync.WaitGroup{}

	wg.Add(1)

	go func() {
		for{
			i++
			j++
			time.Sleep(time.Second*1)
			if i==20{
				break
			}
		}
		wg.Done()
	}()


	//http.ServeFile()
	//http.ServeContent()

	//fmt.Println(len(strings.Split(":11223",":")))
	//
	//s,err:=net.ResolveIPAddr("ip4", "")
	//if err!=nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(s.String())
	//}
	//
	//if _,err1:=strconv.Atoi("5544");err1!=nil{
	//	fmt.Println(err1)
	//}

	for{
		fmt.Printf("\033[1ALeft Second:%d,%d\033[K\n",i,j)

		if i>=20{
			break
		}
		time.Sleep(time.Second*2)
	}

	wg.Wait()





}

func testlist()  {
	//l:=list.NewList(list.Cmp)
	//l.SetCloneFunc(func(v1 interface{}) (r interface{}) {
	//	var v2 int
	//	v2 =v1.(int)
	//	return v2
	//})
	//l.SetDistance(func(v1 interface{}) *big.Int {
	//	v2:=v1.(int)
	//
	//	z:=big.NewInt(1000)
	//
	//	return z.Sub(big.NewInt(int64(v2)),z)
	//
	//})
	//
	//
	//l.AddValueSort(int(1))
	//l.AddValueSort(int(2))
	//l.AddValueSort(int(100))
	//l.AddValueSort(int(100))
	//l.AddValueSort(int(1002))
	//l.AddValueSort(int(10002))
	//l.AddValueSort(int(1009))
	//l.AddValueSort(int(89))
	//l.AddValueSort(int(64))
	//l.AddValueSort(int(77))
	//
	////l.Traverse("==>",list.Print)
	//lnew,_:=l.DeepClone()
	////lnew.Traverse("++>",list.Print)
	//
	//l.DeepConCat(lnew,true)
	//
	////l.Traverse("==++>",list.Print)
	//
	//
	//
	//
	//
	//fmt.Println(dht.GetDhtHashV(0))
	//fmt.Println(dht.GetDhtHashV(1))
	//fmt.Println(dht.GetDhtHashV(2))
	//fmt.Println(dht.GetDhtHashV(3))
	//fmt.Println(dht.GetDhtHashV(4))
	//fmt.Println(dht.GetDhtHashV(6))
	//fmt.Println(dht.GetDhtHashV(8))
	//fmt.Println(dht.GetDhtHashV(31))
	//fmt.Println(dht.GetDhtHashV(32))
	//fmt.Println(dht.GetDhtHashV(69))
	//fmt.Println(dht.GetDhtHashV(112))
	//fmt.Println(dht.GetDhtHashV(123))
	//fmt.Println(dht.GetDhtHashV(233))
	//fmt.Println(dht.GetDhtHashV(0x20151031))
	//
	//
	//fmt.Println(dht.GetDhtHashV1(0))
	//fmt.Println(dht.GetDhtHashV1(1))
	//fmt.Println(dht.GetDhtHashV1(2))
	//fmt.Println(dht.GetDhtHashV1(3))
	//fmt.Println(dht.GetDhtHashV1(4))
	//fmt.Println(dht.GetDhtHashV1(6))
	//fmt.Println(dht.GetDhtHashV1(8))
	//fmt.Println(dht.GetDhtHashV1(31))
	//fmt.Println(dht.GetDhtHashV1(32))
	//fmt.Println(dht.GetDhtHashV1(69))
	//fmt.Println(dht.GetDhtHashV1(112))
	//fmt.Println(dht.GetDhtHashV1(123))
	//fmt.Println(dht.GetDhtHashV1(233))
	//fmt.Println(dht.GetDhtHashV1(0x20151031))


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


