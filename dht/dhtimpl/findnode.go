package dhtimpl

import (
	"github.com/rickeyliao/ServiceAgent/dht/pb"
	"github.com/rickeyliao/ServiceAgent/dht"
	"github.com/rickeyliao/ServiceAgent/dht/dhttable"
	"github.com/golang/protobuf/proto"
	"log"
	"errors"
	"github.com/kprc/nbsnetwork/common/list"
	"math/big"
)



func (node *NbsNode) encFindNode(key []byte) (uint64, []byte) {
	req := &pbdht.Dhtmessage{}

	req.Sn = dht.GetNextMsgCnt()

	req.Msgtyp = dht.FIND_NODE_REQ

	req.Localnbsaddr = GetLocalNode().NbsAddr
	req.Remotenbsaddr = node.NbsAddr
	req.Data = key

	if data, err := proto.Marshal(req); err != nil {
		log.Fatal("Marshall Ping Request Message Failed")
		return 0, nil
	} else {
		return req.Sn, data
	}
}

func (node *NbsNode)updateByFindNode(key []byte,buf []byte, reqsn uint64) (list.List,error) {

	resp := &pbdht.Dhtmessage{}

	if err := proto.Unmarshal(buf, resp); err != nil {
		return nil,err
	}

	if reqsn != resp.Sn {
		return nil,errors.New("SerialNumber not Corrected!")
	}

	if !node.AddrCmp(resp.Localnbsaddr) {

		dhttable.GetRouteTableInst().Del(NewDhtNode(node.NbsAddr,node.Ipv4Addr))
		dhttable.GetRouteTableInst().UpdateOrder(NewDhtNode(resp.Localnbsaddr,node.Ipv4Addr))

		return nil,errors.New("Address not corrected!")
	}

	dhttable.GetRouteTableInst().UpdateOrder(NewDhtNode(resp.Localnbsaddr,node.Ipv4Addr))

	nl:=&pbdht.NbsNodeList{}
	if err:=proto.Unmarshal(resp.Data,nl);err!=nil{
		return nil,err
	}

	l:=list.NewList(func(v1 interface{}, v2 interface{}) int {
		bg1,bg2:=v1.(dhttable.DhtNode).GetBigInt(),v2.(dhttable.DhtNode).GetBigInt()
		return bg1.Cmp(bg2)
	})
	
	l.SetCloneFunc(func(v1 interface{}) (r interface{}) {
		return v1.(dhttable.DhtNode).Clone()
	})
	l.SetSortFunc(func(v1 interface{}, v2 interface{}) int {
		bg1,bg2:=v1.(dhttable.DhtNode).GetBigInt(),v2.(dhttable.DhtNode).GetBigInt()
		z1,z2:=&big.Int{},&big.Int{}
		bgk:=(&big.Int{}).SetBytes(key)
		d1,d2:=z1.Xor(bg1,bgk),z2.Xor(bg2,bgk)

		return d1.Cmp(d2)
	})

	for _,nn:=range nl.Nodes{
		dn:=NewDhtNode(nn.Nbsaddr,nn.INetAddr)
		l.AddValueOrder(dn)
	}

	return l,nil
}

