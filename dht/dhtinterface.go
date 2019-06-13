package dht

import "github.com/kprc/nbsnetwork/common/list"

type DhtInter interface {
	Ping() (bool,error)
	Store(key []byte,v []byte) error
	FindNode(key []byte) (list.List,error)
	FindValue(key []byte) (list.List,[]byte,error)
}



