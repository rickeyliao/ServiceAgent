package common

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"log"
)

func Post(url string,jsonstr string) (jsonret string, code int,err error) {
	log.Println(url)
	log.Println(jsonstr)

	bjson := []byte(jsonstr)
	req,err:=http.NewRequest("POST",url,bytes.NewBuffer(bjson))
	if err!= nil{
		return "",0,err
	}

	req.Header.Set("Content-Type","application/json")

	client:=&http.Client{}
	resp,errresp:=client.Do(req)

	if errresp!=nil{
		return "", 0,errresp
	}

	defer resp.Body.Close()



	body,errbody:=ioutil.ReadAll(resp.Body)
	if errbody != nil{
		return "",0,errbody
	}

	return string(body),resp.StatusCode,nil
}
