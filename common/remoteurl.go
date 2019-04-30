package common

import (
	"strconv"
	"sync"
)

type remoteurl struct {
	host string
	port uint16
}

type RemoteUrl interface {
	GetHostName(path string) string
	GetHostNameSSL(path string) string
	SetHost(host string)
	GetHost() string
	SetPort(port uint16)
	GetPort() uint16
}


var (
	remoteurlinst RemoteUrl
	remoteurllock sync.Mutex

)

func GetRemoteUrlInst() RemoteUrl {
	return remoteurlinst
}

func NewRemoteUrl(host string,port string) RemoteUrl {
	if port == ""{
		port = "80"
	}

	if host == ""{
		host = "localhost"
	}

	if remoteurlinst == nil {
		remoteurllock.Lock()
		defer remoteurllock.Unlock()

		if remoteurlinst == nil{
			ru:=&remoteurl{}
			ru.host = host
			p,_ := strconv.Atoi(port)
			ru.port = uint16(p)

			remoteurlinst = ru
		}
	}

	return remoteurlinst
}


func (ru *remoteurl)getHostName(path string) string  {
	var port string
	if ru.port == 80{
		port = ""
	}else{
		port = ":"+strconv.Itoa(int(ru.port))
	}

	return ru.host+port + path
}

func (ru *remoteurl)GetHostName(path string) string  {
	return "http://"+ru.getHostName(path)
}

func (ru *remoteurl)GetHostNameSSL(path string) string  {
	return "https://"+ru.getHostName(path)
}


func (ru *remoteurl)SetHost(host string){
	ru.host = host
}
func (ru *remoteurl)GetHost() string{
	return ru.host
}
func (ru *remoteurl)SetPort(port uint16){
	ru.port = port
}
func (ru *remoteurl)GetPort() uint16{
	return ru.port
}
