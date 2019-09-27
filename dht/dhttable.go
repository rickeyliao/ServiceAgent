package dht

import (
	"github.com/kprc/nbsnetwork/common/hashlist"
	"math/big"
	"sync"
)

type IDhtNode interface {
	GetBigInt() *big.Int
	Clone() IDhtNode
	GetLastAccessTime() int64
}

var (
	routetable     hashlist.HashList
	routetablelock sync.Mutex
)

func GetRouteTableInst() hashlist.HashList {
	if routetable == nil {
		routetablelock.Lock()
		defer routetablelock.Unlock()
		if routetable == nil {
			routetable = newRouteTable()
		}
	}

	return routetable
}

func newRouteTable() hashlist.HashList {
	t := hashlist.NewHashList(256, func(v interface{}) uint {
		vbgint := v.(IDhtNode).GetBigInt()
		localBgInt := GetLocalNode().GetBgInt()
		z := &big.Int{}
		bitl := z.Xor(vbgint, localBgInt).BitLen()
		if bitl > 0 {
			bitl -= 1
		}
		return uint(bitl)
	}, func(v1 interface{}, v2 interface{}) int {
		bgnode1, bgnode2 := v1.(IDhtNode), v2.(IDhtNode)
		return bgnode1.GetBigInt().Cmp(bgnode2.GetBigInt())
	})
	t.SetLimitCnt(DHT_K)
	t.SetSortFunc(func(v1 interface{}, v2 interface{}) int {
		tm1, tm2 := v1.(IDhtNode).GetLastAccessTime(), v2.(IDhtNode).GetLastAccessTime()

		return int(tm1 - tm2)
	})

	t.SetCloneFunc(func(v1 interface{}) (r interface{}) {
		r = v1.(IDhtNode).Clone()
		return
	})

	return t
}

func Distance(v1 interface{}, v2 interface{}) *big.Int {
	bgv1, bgv2 := v1.(IDhtNode).GetBigInt(), v2.(IDhtNode).GetBigInt()

	z := &big.Int{}

	return z.Xor(bgv1, bgv2)
}
