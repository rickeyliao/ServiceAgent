package dhttable

import (
	"math/big"
	"github.com/kprc/nbsnetwork/common/hashlist"
	"sync"
	"github.com/rickeyliao/ServiceAgent/dht/dhtimpl"
	"github.com/rickeyliao/ServiceAgent/dht"
)

type DhtNode interface {
	GetBigInt() *big.Int
	Clone() DhtNode
	GetLastAccessTime() int64
}



var (
	routetable hashlist.HashList
	routetablelock sync.Mutex
)

func GetRouteTableInst() hashlist.HashList  {
	if routetable == nil{
		routetablelock.Lock()
		defer routetablelock.Unlock()
		if routetable == nil{
			routetable = newRouteTable()
		}
	}


	return routetable
}

func newRouteTable() hashlist.HashList {
	t := hashlist.NewHashList(256, func(v interface{}) uint {
		vbgint:=v.(DhtNode).GetBigInt()
		localBgInt:=dhtimpl.GetLocalNode().GetBgInt()
		z:=&big.Int{}
		bitl:=z.Xor(vbgint,localBgInt).BitLen()
		if bitl>0{
			bitl-=1
		}
		return uint(bitl)
	}, func(v1 interface{}, v2 interface{}) int {
		bgnode1,bgnode2:=v1.(DhtNode),v2.(DhtNode)
		return bgnode1.GetBigInt().Cmp(bgnode2.GetBigInt())
	})
	t.SetLimitCnt(dht.DHT_K)
	t.SetSortFunc(func(v1 interface{}, v2 interface{}) int {
		tm1,tm2:=v1.(DhtNode).GetLastAccessTime(),v2.(DhtNode).GetLastAccessTime()
		
		return int(tm1-tm2)
	})
	
	t.SetCloneFunc(func(v1 interface{}) (r interface{}) {
		r = v1.(DhtNode).Clone()
		return
	})

	return t
}

func Distance(v1 interface{},v2 interface{}) *big.Int{
	bgv1,bgv2:=v1.(DhtNode).GetBigInt(),v2.(DhtNode).GetBigInt()

	z:=&big.Int{}

	return z.Xor(bgv1,bgv2)
}









