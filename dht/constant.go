package dht

import (
	"encoding/binary"
	"math/big"
)

const (
	PING_REQ  uint32 = 1
	PING_RESP uint32 = 2
	FIND_NODE_REQ  uint32 = 3
	FIND_NODE_RESP uint32 = 4
	FIND_VALUE_REQ uint32 = 5
	FIND_VALUE_RESP uint32 = 6
	STORE_REQ uint32 = 7
	STORE_RESP uint32 = 8
)


const(
	DHT_K int32 = 20
	DHT_A int32 = 3
)


func GetDhtHashV(n int) int {
	cnt:=0
	if n==0{
		return cnt
	}

	for {
		n=n>>1
		cnt ++
		if n == 0{
			return cnt
		}

	}
}

func GetDhtHashV1(n int) int  {
	bn:=make([]byte,4)
	binary.BigEndian.PutUint32(bn,uint32(n))

	return big.NewInt(0).SetBytes(bn).BitLen()
}