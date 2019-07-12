package dhtimpl

import (
	"github.com/rickeyliao/ServiceAgent/dht/pb"
	"github.com/rickeyliao/ServiceAgent/dht"
	"github.com/rickeyliao/ServiceAgent/dht/dhttable"
	"github.com/golang/protobuf/proto"
	"log"
	"errors"
	"github.com/kprc/nbsnetwork/common/list"
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

func (node *NbsNode)updateByFindNode(buf []byte, reqsn uint64) (list.List,error) {

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



	return nil,nil
}

