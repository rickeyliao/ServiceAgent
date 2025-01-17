package dht

import (
	"github.com/mr-tron/base58"
	"github.com/rickeyliao/ServiceAgent/common"
	"math/big"
	"sync"
)

type LocalNode struct {
	NbsNode
	Bgint *big.Int
}

var (
	localNode      *LocalNode
	localNode_lock sync.Mutex
)

func GetLocalNode() *LocalNode {
	if localNode == nil {
		localNode_lock.Lock()
		defer localNode_lock.Unlock()

		if localNode == nil {
			localNode = newLocalNode()
		}
	}

	return localNode
}

func newLocalNode() *LocalNode {

	sac := common.GetSAConfig()

	nbsAddr := sac.NbsRsaAddr

	node := &LocalNode{}

	if addr, err := base58.Decode(nbsAddr[2:]); err != nil {
		return nil
	} else {
		node.NbsAddr = addr
		node.Port = NbsAddr2Port(node.NbsAddr)
	}
	bgi := &big.Int{}
	node.Bgint = bgi.SetBytes(node.NbsAddr)

	return node

}

func (ln *LocalNode) GetBgInt() *big.Int {
	return ln.Bgint
}
