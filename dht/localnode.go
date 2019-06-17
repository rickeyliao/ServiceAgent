package dht

import (
	"sync"
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/mr-tron/base58"
)

var (
	localNode *NbsNode
	localNode_lock sync.Mutex
)

func GetLocalNode() *NbsNode {
	if localNode == nil{
		localNode_lock.Lock()


		localNode = newLocalNode()

		localNode_lock.Unlock()

	}

	return localNode
}

func newLocalNode() *NbsNode  {

	sac:=common.GetSAConfig()

	nbsAddr:=sac.NbsRsaAddr

	node:=&NbsNode{}

	if addr,err:=base58.Decode(nbsAddr[2:]);err!=nil{
		return nil
	}else{
		node.NbsAddr = addr
		node.Port = NbsAddr2Port(node.NbsAddr)
	}

	return node

}




