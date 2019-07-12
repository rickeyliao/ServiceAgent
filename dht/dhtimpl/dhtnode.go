package dhtimpl

import (
	"math/big"
	"github.com/kprc/nbsnetwork/tools"
)

type DhtNode struct {
	NbsNode
	TimeStamp int64
}

func (dn *DhtNode)Clone() *DhtNode  {
	tdn:=&DhtNode{}

	tdn.TimeStamp = dn.TimeStamp

	tdn.NbsAddr = make([]byte,len(dn.NbsAddr))
	copy(tdn.NbsAddr,dn.NbsAddr)
	tdn.Port = dn.Port
	tdn.Ipv4Addr = make([]byte,len(dn.Ipv4Addr))
	copy(tdn.Ipv4Addr,dn.Ipv4Addr)

	return tdn
}

func (dn *DhtNode)GetBigInt() *big.Int  {
	bgi:=&big.Int{}

	return bgi.SetBytes(dn.NbsAddr)
}

func (dn *DhtNode)GetLastAccessTime() int64  {
	return dn.TimeStamp
}

func (dn *DhtNode)SetTimeStamp(ts int64)  {
	dn.TimeStamp = ts
}

func NewDhtNodeTimeStamp(nbsaddr,ipaddr []byte,timestamp int64) *DhtNode {
	dn:=&DhtNode{}

	dn.NbsAddr = nbsaddr
	dn.Ipv4Addr = ipaddr
	dn.Port = NbsAddr2Port(dn.NbsAddr)
	dn.TimeStamp = timestamp

	return dn
}

func NewDhtNode(nbsaddr,ipaddr []byte) *DhtNode  {
	dn:=&DhtNode{}

	dn.NbsAddr = nbsaddr
	dn.Ipv4Addr = ipaddr
	dn.Port = NbsAddr2Port(dn.NbsAddr)
	dn.TimeStamp = tools.GetNowMsTime()

	return dn
}