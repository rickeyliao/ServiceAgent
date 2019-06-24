package main

import (
	"io"
	"mime/multipart"

	"os"
	"net/http"
	"sync"
)

func main()  {

	wg := &sync.WaitGroup{}
	name:=os.Args[0]

	piper,pipew:=io.Pipe()

	m:=multipart.NewWriter(pipew)
	//defer m.Close()
	wg.Add(1)
	go func() {
		defer pipew.Close()
		defer piper.Close()
		part,err:=m.CreateFormFile("filename","foo.txt")
		if err!=nil{
			return
		}

		file, err := os.Open(name)
		if err != nil {
			return
		}
		defer file.Close()
		if _, err = io.Copy(part, file); err != nil {
			return
		}

		wg.Done()
	}()

	resp,_:=http.Post("http://192.168.20.178:50810/upload", m.FormDataContentType(), piper)
	wg.Wait()
	if resp!=nil && resp.Body != nil{
		resp.Body.Close()
	}


}
