package dhtserver

import (
	"github.com/rickeyliao/ServiceAgent/dht/pb"
	"net"
	"github.com/rickeyliao/ServiceAgent/dht"

	"errors"
	"github.com/golang/protobuf/proto"
)

func respStore(dm pbdht.Dhtmessage,addr *net.UDPAddr,conn *net.UDPConn) (err error)  {
	if dm.Msgtyp != dht.STORE_REQ{
		return errors.New("store func receive a error type")
	}

	resp:=&pbdht.Dhtmessage{}
	resp.Msgtyp = dht.STORE_RESP
	resp.Sn = dm.Sn
	resp.Localnbsaddr = dht.GetLocalNode().NbsAddr
	resp.Remotenbsaddr = dm.Localnbsaddr

	dht.GetRouteTableInst().UpdateOrder(dht.NewDhtNode(dm.Localnbsaddr,addr.IP))

	var data []byte
	data,err=proto.Marshal(resp)
	if err!=nil{
		return
	}


	_,err=conn.WriteToUDP(data,addr)


	//todo inform to download the file

	return

}

func storeContent(dm pbdht.Dhtmessage) error  {

	sv:=&pbdht.Dhtstore{}

	if err:=proto.Unmarshal(dm.Data,sv);err!=nil{
		return err
	}

	if sv.Share {
		if len(sv.Key)==0 || len(sv.Value)==0{
			return errors.New("sotre share file error")
		}






	}else{
		// go to download the key
	}

	return nil

}