package dht

import (
	"github.com/rickeyliao/ServiceAgent/dht/pb"
	"github.com/gogo/protobuf/proto"
	"log"
	"github.com/pkg/errors"
)

func (node *NbsNode)encPingData() []byte {
	req:=&pbdht.Pingreq{}

	req.Sn = GetNextMsgCnt()

	req.Msgtyp = PING_REQ

	req.Nbsaddr = GetLocalNode().NbsAddr

	if data,err:=proto.Marshal(req);err!=nil{
		log.Fatal("Marshall Ping Request Message Failed")
		return nil
	}else {
		return data
	}
}

func (node *NbsNode)updateByPingResp(buf []byte,reqsn uint64) error  {

	resp:=&pbdht.Pingresp{}

	if err:=proto.Unmarshal(buf,resp);err!=nil{
		return err
	}

	if reqsn != resp.Rcvsn{
		return errors.New("SerialNumber not Corrected!")
	}

	if !node.AddrCmp(resp.Nbsaddr){
		//todo update
	}

	return nil
}