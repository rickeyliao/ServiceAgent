package dht2

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"
	"github.com/rickeyliao/ServiceAgent/common"
	"sync"
)

type NAddr [32]byte

var (
	localNAddrLock sync.Mutex
	localNAddr     *NAddr
	localNID       NID
)

func (na NAddr) ID() NID {
	return NID("91" + base58.Encode(na[:]))
}

func (na NAddr) String() string {
	return base58.Encode(na[:])
}

type NID string

var EmptyId = NID("")

func (nid NID) Addr() (NAddr, error) {
	id := string(nid)

	na := NAddr{}

	if id[:2] != "91" || nid == EmptyId {
		return na, errors.New("ID is empty or ID first 2 byte must be 91")
	}
	bid := base58.Decode(id[2:])
	copy(na[:], bid)

	return na, nil
}

func (nid NID) Addr2() NAddr {

	id := string(nid)
	na := NAddr{}
	bid := base58.Decode(id[2:])

	copy(na[:], bid)

	return na
}

func PubKey2NAddr(pk *rsa.PublicKey) NAddr {
	pubkeybytes := x509.MarshalPKCS1PublicKey(pk)

	s := sha256.New()
	s.Write(pubkeybytes)
	sum := s.Sum(nil)
	na := NAddr{}

	copy(na[:], sum)

	return na
}

func PubKey2ID(pk *rsa.PublicKey) NID {
	na := PubKey2NAddr(pk)

	return na.ID()
}

func GetLocalNAddr() NAddr {
	if localNID == EmptyId {
		localNAddrLock.Lock()
		defer localNAddrLock.Unlock()
		if localNID == EmptyId {
			pk := &common.GetSAConfig().PrivKey.PublicKey
			na := PubKey2NAddr(pk)

			localNAddr = &na

			localNID = na.ID()
		}
	}

	return *localNAddr
}

func GetLocalNID() NID {
	if localNID == EmptyId {
		GetLocalNAddr()
	}

	return localNID
}

func (na NAddr) Bytes() []byte {
	return na[:]
}

func (na NAddr) Array() [32]byte {
	return na
}

func (na NAddr) Cmp(na2 NAddr) bool {
	if 0 == bytes.Compare(na.Bytes(), na2.Bytes()) {
		return true
	} else {
		return false
	}
}

func (na NAddr) Len() int {
	return len(na)
}
