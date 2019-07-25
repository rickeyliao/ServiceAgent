package dhtserver

import (
	"github.com/rickeyliao/ServiceAgent/dht/pb"
	"net"
	"github.com/rickeyliao/ServiceAgent/dht"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/btcsuite/btcutil/base58"
	"github.com/rickeyliao/ServiceAgent/dht/dhtdb"
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

	storeContent(dm,addr.IP)

	_,err=conn.WriteToUDP(data,addr)

	return

}

func storeContent(dm pbdht.Dhtmessage,ip net.IP) error  {

	sv:=&pbdht.Dhtstore{}

	if err:=proto.Unmarshal(dm.Data,sv);err!=nil{
		return err
	}
	if len(sv.Key)==0 {
		return errors.New("sotre file error")
	}

	key:="c1"+base58.Encode(sv.Key)
	if sv.Share {

		if len(sv.Value)==0{
			return errors.New("sotre share file error")
		}

		nbsaddrs:=make([]string,0)

		for _,v:=range sv.Value{
			nbsaddrs = append(nbsaddrs,"91"+base58.Encode(v))
		}

		if _,err:=dhtdb.Find(key);err!=nil{
			dhtdb.Insert(key,false,true,nbsaddrs,false)
		}else{
			dhtdb.Update(key,false,true,nbsaddrs,false)
		}
	}else{
		if _,err:=dhtdb.Find(key);err!=nil{
			dhtdb.Insert(key,false,false,nil,false)
			dht.GetDownloadQueue().EnQueue(sv.Key,ip)
		}

	}

	return nil

}