package dhtserver

import "github.com/rickeyliao/ServiceAgent/dht"

func init()  {
	dh:=GetDhtHandlerInst()
	dh.Reg(dht.PING_REQ,respPing)
	dh.Reg(dht.FIND_NODE_REQ,respFindNode)
	dh.Reg(dht.STORE_REQ,respStore)
}

