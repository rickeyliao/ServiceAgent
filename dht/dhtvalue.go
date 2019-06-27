package dht

import "github.com/gogo/protobuf/io"

type DhtValue struct {
	Buf []byte
	io.Reader
	IsStream bool
}
