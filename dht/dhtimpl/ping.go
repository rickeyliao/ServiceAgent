package dhtimpl

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/rickeyliao/ServiceAgent/dht/pb"
	"github.com/rickeyliao/ServiceAgent/dht"
	"log"
	"github.com/rickeyliao/ServiceAgent/dht/dhttable"
)

func (node *NbsNode) encPingData() (uint64, []byte) {
	req := &pbdht.Dhtmessage{}

	req.Sn = dht.GetNextMsgCnt()

	req.Msgtyp = dht.PING_REQ

	req.Localnbsaddr = GetLocalNode().NbsAddr
	req.Remotenbsaddr = node.NbsAddr

	if data, err := proto.Marshal(req); err != nil {
		log.Fatal("Marshall Ping Request Message Failed")
		return 0, nil
	} else {
		return req.Sn, data
	}
}

func (node *NbsNode) updateByPingResp(buf []byte, reqsn uint64) error {

	resp := &pbdht.Dhtmessage{}

	if err := proto.Unmarshal(buf, resp); err != nil {
		return err
	}

	if reqsn != resp.Sn {
		return errors.New("SerialNumber not Corrected!")
	}

	if !node.AddrCmp(resp.Localnbsaddr) {

		dhttable.GetRouteTableInst().Del(NewDhtNode(node.NbsAddr,node.Ipv4Addr))
		dhttable.GetRouteTableInst().UpdateOrder(NewDhtNode(resp.Localnbsaddr,node.Ipv4Addr))

		return errors.New("Address not corrected!")
	}else {
		dhttable.GetRouteTableInst().UpdateOrder(NewDhtNode(resp.Localnbsaddr,node.Ipv4Addr))
	}

	return nil
}
