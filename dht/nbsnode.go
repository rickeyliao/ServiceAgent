package dht

import (
	"github.com/pkg/errors"
	"github.com/kprc/nbsnetwork/common/list"
	"net"
	"time"
)

type NbsNode struct {
	Ipv4Addr []byte
	Port uint16    //caculate from NbsAddr
	NbsAddr []byte
}

func NbsAddr2Port(nbsaddr []byte)  uint16 {
	return uint16(50810)
}

func NewNbsNode(ipaddr []byte,nbsaddr []byte) *NbsNode {
	nn:=&NbsNode{}

	nn.Ipv4Addr = ipaddr
	nn.NbsAddr = nbsaddr

	nn.Port = NbsAddr2Port(nbsaddr)

	return nn
}

func (node *NbsNode)AddrCmp(addr []byte) bool  {
	if len(node.NbsAddr) != len(addr){
		return false
	}

	for idx,b:=range node.NbsAddr{
		if b != addr[idx]{
			return false
		}
	}
	return true
}

func (node *NbsNode)Ping() (bool,error)  {
	remoteaddr:=&net.UDPAddr{
		IP:node.Ipv4Addr,
		Port:int(node.Port),
	}

	localaddr:=&net.UDPAddr{
		IP:net.ParseIP("0.0.0.0"),
		Port:0,
	}

	conn,err:=net.DialUDP("udp",localaddr,remoteaddr)
	if err!=nil{
		return false,errors.New("Dial UDP Error")
	}
	defer conn.Close()

	data:=node.encPingData()
	if data == nil{
		return false,errors.New("enc Ping Request Failed")
	}

	var n int
	n,err=conn.Write(data)
	if err!=nil || n != len(data){
		return false,errors.New("Send Ping Request Failed")
	}


	conn.SetReadDeadline(time.Now().Add(time.Second*2))

	buf :=make([]byte,1024)

	n,err=conn.Read(buf)
	if err!=nil{
		return false,err
	}

	return true,nil
}

func (node *NbsNode) Store(key []byte,v []byte) error {
	

	return nil
}

func (node *NbsNode)FindNode(key []byte) (list.List,error)  {
	return nil,nil
}

func (node *NbsNode)FindValue(key []byte) (list.List,[]byte,error)  {
	return nil,nil,nil
}
