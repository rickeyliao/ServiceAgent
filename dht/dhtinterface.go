package dht

import (
	"github.com/kprc/nbsnetwork/common/list"
)

type DhtValue struct {
	HaveData bool
	data     []byte
}

type DhtInter interface {
	Ping() (bool, error)
	Store(key []byte) error
	FindNode(key []byte) (list.List, error)
	FindValue(key []byte) (list.List, *DhtValue, error)
}
