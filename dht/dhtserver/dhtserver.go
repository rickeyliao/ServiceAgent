package dhtserver

import (

	"net"
	"sync"
	"log"
	"github.com/rickeyliao/ServiceAgent/dht/pb"
	"github.com/gogo/protobuf/proto"

	"github.com/rickeyliao/ServiceAgent/dht"
)

type DhtServer struct {
	servernode *dht.LocalNode
	conn *net.UDPConn
}

var (
	dhtserver *DhtServer
	dhtserverlock sync.Mutex
)


func GetDhtServer() *DhtServer {

	if dhtserver !=nil{
		return dhtserver
	}

	dhtserverlock.Lock()
	defer dhtserverlock.Unlock()

	if dhtserver != nil{
		return dhtserver
	}


	ds:=&DhtServer{}

	ds.servernode = dht.GetLocalNode()

	ds.conn = nil

	dhtserver = ds

	return ds
}


func (ds *DhtServer)Run()  {
	if ds.conn != nil{
		return
	}

	lddr:=&net.UDPAddr{
		IP:ds.servernode.Ipv4Addr,
		Port:int(ds.servernode.Port),
	}
	conn,err:=net.ListenUDP("udp4",lddr)
	if err !=nil{
		log.Fatal("Can't Start DHT SERVER")
	}
	for{
		buf:=make([]byte,1024)
		n,addr,err:=conn.ReadFromUDP(buf)
		if err!=nil{
			log.Fatal("dht service failed")
			return
		}
		handleBuf(buf[:n],addr,conn)
	}

}

func handleBuf(buf []byte,addr *net.UDPAddr,conn *net.UDPConn)  {
	dm:=pbdht.Dhtmessage{}

	if err:=proto.Unmarshal(buf,&dm);err!=nil{
		return
	}

	if err:=GetDhtHandlerInst().Dispatch(dm,addr,conn);err!=nil{
		log.Println(err)
	}
}


func (ds *DhtServer)Stop()  {
	if ds.conn != nil{
		ds.conn.Close()
		ds.conn = nil
	}
}


