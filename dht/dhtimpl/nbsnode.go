package dhtimpl

import (
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/pkg/errors"
	"net"
	"time"
	"github.com/rickeyliao/ServiceAgent/dht/dhttable"
	"github.com/rickeyliao/ServiceAgent/common"
)

type NbsNode struct {
	Ipv4Addr []byte
	Port     uint16 //caculate from NbsAddr
	NbsAddr  []byte
}

func NbsAddr2Port(nbsaddr []byte) uint16 {
	return common.GetSAConfig().DhtListenPort
}

func NewNbsNode(ipaddr []byte, nbsaddr []byte) *NbsNode {
	nn := &NbsNode{}

	nn.Ipv4Addr = ipaddr
	nn.NbsAddr = nbsaddr

	nn.Port = NbsAddr2Port(nbsaddr)

	return nn
}

func (node *NbsNode) AddrCmp(addr []byte) bool {
	if len(node.NbsAddr) != len(addr) {
		return false
	}

	for idx, b := range node.NbsAddr {
		if b != addr[idx] {
			return false
		}
	}
	return true
}

func (node *NbsNode)connect() (*net.UDPConn,error)  {
	remoteaddr := &net.UDPAddr{
		IP:   node.Ipv4Addr,
		Port: int(node.Port),
	}

	localaddr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 0,
	}

	return net.DialUDP("udp", localaddr, remoteaddr)
}

func (node *NbsNode) Ping() (bool, error) {

	conn, err := node.connect()
	if err != nil {
		return false, errors.New("Dial UDP Error")
	}
	defer conn.Close()

	sn, data := node.encPingData()
	if data == nil {
		return false, errors.New("enc Ping Request Failed")
	}

	var n int
	n, err = conn.Write(data)
	if err != nil || n != len(data) {
		return false, errors.New("Send Ping Request Failed")
	}

	conn.SetReadDeadline(time.Now().Add(time.Second * 2))

	buf := make([]byte, 1024)

	n, err = conn.Read(buf)
	if err != nil {
		return false, err
	}

	if err = node.updateByPingResp(buf, sn); err != nil {
		return false, err
	}

	return true, nil
}

func (node *NbsNode) Store(key []byte, dv *dhttable.DhtNode) error {
	//if node.AddrCmp(dhtserver.GetLocalNode().NbsAddr) {
	//	//save local db
	//
	//	return nil
	//}

	//send to remote

	return nil
}

func (node *NbsNode) FindNode(key []byte) (list.List, error) {

	conn,err:=node.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	sn, data := node.encFindNode(key)
	if data == nil {
		return nil, errors.New("enc FindNode Request Failed")
	}

	var n int
	n, err = conn.Write(data)
	if err != nil || n != len(data) {
		return nil, errors.New("Send FindNode Request Failed")
	}

	conn.SetReadDeadline(time.Now().Add(time.Second * 2))

	buf := make([]byte, 2048)

	n, err = conn.Read(buf)
	if err != nil {
		return nil, err
	}

	var l list.List
	if l,err = node.updateByFindNode(key,buf, sn); err != nil {
		return nil, err
	}

	return l, nil
}


func (node *NbsNode) FindValue(key []byte) (list.List, []byte, error) {
	return nil, nil, nil
}


