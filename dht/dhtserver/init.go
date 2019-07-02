package dhtserver

import "github.com/rickeyliao/ServiceAgent/dht"

func init()  {
	dh:=GetDhtHandlerInst()
	dh.Reg(dht.PING_REQ,respPing)
}

