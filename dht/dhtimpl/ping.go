package dhtimpl

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/rickeyliao/ServiceAgent/dht/pb"
	"github.com/rickeyliao/ServiceAgent/dht"
	"log"
)

func (node *NbsNode) encPingData() (uint64, []byte) {
	req := &pbdht.Dhtmessage{}

	req.Sn = dht.GetNextMsgCnt()

	req.Msgtyp = dht.PING_REQ

	//req.Nbsaddr = dhtserver.GetLocalNode().NbsAddr

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

	if !node.AddrCmp(resp.Nbsaddr) {
		//todo update
	}

	return nil
}
