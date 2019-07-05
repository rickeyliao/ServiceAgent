package dhttable

import "math/big"

type DhtNode interface {
	GetBigInt() *big.Int
	Clone() DhtNode
}



