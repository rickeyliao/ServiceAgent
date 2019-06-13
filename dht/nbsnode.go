package dht

import (
	"github.com/pkg/errors"
	"github.com/kprc/nbsnetwork/common/list"
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

func (node *NbsNode)Ping() (bool,error)  {

	err:=errors.New("Failed")

	return false,err
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
