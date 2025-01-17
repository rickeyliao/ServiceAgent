package dhtserver

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/rickeyliao/ServiceAgent/dht"
	"github.com/rickeyliao/ServiceAgent/dht/pb"
	"net"
)

func respPing(dm pbdht.Dhtmessage, addr *net.UDPAddr, conn *net.UDPConn) (err error) {
	if dm.Msgtyp != dht.PING_REQ {
		return errors.New("ping func receive a error type")
	}

	resp := &pbdht.Dhtmessage{}
	resp.Msgtyp = dht.PING_RESP
	resp.Sn = dm.Sn
	resp.Localnbsaddr = dht.GetLocalNode().NbsAddr
	resp.Remotenbsaddr = dm.Localnbsaddr

	dht.GetRouteTableInst().UpdateOrder(dht.NewDhtNode(dm.Localnbsaddr, addr.IP))

	var data []byte
	data, err = proto.Marshal(resp)
	if err != nil {
		return
	}

	_, err = conn.WriteToUDP(data, addr)

	return

}
