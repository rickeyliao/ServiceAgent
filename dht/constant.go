package dht

import (
	"encoding/binary"
	"math/big"
)

const (
	PING_REQ  uint32 = 1
	PING_RESP uint32 = 2
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