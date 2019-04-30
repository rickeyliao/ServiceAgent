package common

import (
	"net/http"
	"io/ioutil"
	"log"
)

func Get(url string) (jsonret string,err error){

	log.Println(url)

	req,err:=http.NewRequest("GET",url,nil)
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
