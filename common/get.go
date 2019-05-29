package common

import (
	"io/ioutil"
	"log"
	"net/http"
)

func Get(url string) (jsonret string, code int, err error) {

	log.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, errresp := client.Do(req)

	if errresp != nil {
		return "", 0, errresp
	}

	defer resp.Body.Close()

	body, errbody := ioutil.ReadAll(resp.Body)
	if errbody != nil {
		return "", 0, errbody
	}

	return string(body), resp.StatusCode, nil
}
