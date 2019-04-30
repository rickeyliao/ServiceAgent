package common

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"log"
)

func Post(url string,jsonstr string) (jsonret string, err error) {
	log.Println(url,jsonstr)

	bjson := []byte(jsonstr)
	req,err:=http.NewRequest("POST",url,bytes.NewBuffer(bjson))
	if err!= nil{
		return "",err
	}

	req.Header.Set("Content-Type","application/json")

	client:=&http.Client{}
	resp,errresp:=client.Do(req)

	if errresp!=nil{
		return "", errresp
	}

	defer resp.Body.Close()

	body,errbody:=ioutil.ReadAll(resp.Body)
	if errbody != nil{
		return "",errbody
	}

	return string(body),nil
}
