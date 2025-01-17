package dhtserver

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/rickeyliao/ServiceAgent/dht"
	"github.com/rickeyliao/ServiceAgent/dht/pb"
	"math/big"
	"net"
)

func respFindNode(dm pbdht.Dhtmessage, addr *net.UDPAddr, conn *net.UDPConn) (err error) {
	if dm.Msgtyp != dht.FIND_NODE_REQ {
		return errors.New("findnode func receive a error type")
	}

	resp := &pbdht.Dhtmessage{}
	resp.Msgtyp = dht.FIND_NODE_RESP
	resp.Sn = dm.Sn
	resp.Localnbsaddr = dht.GetLocalNode().NbsAddr
	resp.Remotenbsaddr = dm.Localnbsaddr

	dht.GetRouteTableInst().UpdateOrder(dht.NewDhtNode(dm.Localnbsaddr, addr.IP))

	if len(dm.Data) == 0 {
		return errors.New("no key found")
	}

	knode := dht.NewDhtNode(dm.Data, nil)

	l := dht.GetRouteTableInst().GetNodes(knode, dht.DHT_K, func(v1 interface{}, v2 interface{}) int {
		bg1, bg2 := v1.(dht.IDhtNode).GetBigInt(), v2.(dht.IDhtNode).GetBigInt()
		z1, z2 := &big.Int{}, &big.Int{}
		bgk := knode.GetBigInt()
		d1, d2 := z1.Xor(bg1, bgk), z2.Xor(bg2, bgk)

		return d1.Cmp(d2)
	})

	if l != nil && l.Count() > 0 {
		dhtnodes := &pbdht.NbsNodeList{}
		dhtnodes.Nodes = make([]*pbdht.NbsNode, 0)
		it := l.ListIterator(int(dht.DHT_K))

		for {
			v := it.Next()
			if v == nil {
				return nil
			}
			pbnn := &pbdht.NbsNode{}
			node := v.(*dht.DhtNode)
			pbnn.INetAddr = node.Ipv4Addr
			pbnn.Nbsaddr = node.NbsAddr

			dhtnodes.Nodes = append(dhtnodes.Nodes, pbnn)

		}
		var nodesdata []byte
		nodesdata, err = proto.Marshal(dhtnodes)
		if err != nil {
			return
		}

		resp.Data = nodesdata
	}

	var data []byte
	data, err = proto.Marshal(resp)
	if err != nil {
		return
	}

	_, err = conn.WriteToUDP(data, addr)

	return

}
