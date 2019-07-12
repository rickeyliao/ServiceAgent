package dht

import (
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/rickeyliao/ServiceAgent/dht/dhtimpl"
)

type DhtInter interface {
	Ping() (bool, error)
	Store(key []byte) error
	FindNode(key []byte) (list.List, error)
	FindValue(key []byte) (list.List, *dhtimpl.DhtValue, error)
}
