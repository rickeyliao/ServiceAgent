package dhtserver

import (
	"github.com/mr-tron/base58"
	"github.com/rickeyliao/ServiceAgent/common"
	"sync"
	"github.com/rickeyliao/ServiceAgent/dht/dhtimpl"
)

var (
	localNode      *dhtimpl.NbsNode
	localNode_lock sync.Mutex
)

func GetLocalNode() *dhtimpl.NbsNode {
	if localNode == nil {
		localNode_lock.Lock()

		localNode = newLocalNode()

		localNode_lock.Unlock()

	}

	return localNode
}

func newLocalNode() *dhtimpl.NbsNode {

	sac := common.GetSAConfig()

	nbsAddr := sac.NbsRsaAddr

	node := &dhtimpl.NbsNode{}

	if addr, err := base58.Decode(nbsAddr[2:]); err != nil {
		return nil
	} else {
		node.NbsAddr = addr
		node.Port = dhtimpl.NbsAddr2Port(node.NbsAddr)
	}

	return node

}

