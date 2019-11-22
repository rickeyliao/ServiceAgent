package dht2

import (
	"net"
)

const (
	MsgTypeLen             int = 1
	InternalAddrCountLen   int = 2
	NatAddrCountLen        int = 2
	SerialNumberBytesCount     = 32
)

type OnlineCtrlMsg struct {
	typ       byte
	sn        [SerialNumberBytesCount]byte
	localAddr *P2pAddr
}

func PackCtrlMsg(ocm *OnlineCtrlMsg, buf []byte) int {
	offset := 0
	buf[0] = ocm.typ
	offset += MsgTypeLen
	copy(buf[offset:], ocm.sn[:])
	offset += len(ocm.sn)

	p2paddr := &P2pAddr{}

	offset += PackP2pAddr(p2paddr, buf[offset:])

	ocm.localAddr = p2paddr

	return offset
}

func putUint16(buf []byte, v uint16) {
	buf[0] = byte(v >> 8)
	buf[1] = byte(v)
}

func toUint16(buf []byte) uint16 {
	var n uint16
	n = uint16(buf[0])
	n <<= 8
	n |= uint16(buf[1])

	return n
}

func PackP2pAddr(addr *P2pAddr, buf []byte) int {
	offset := 0
	copy(buf, addr.NbsAddr.Bytes())
	offset += addr.NbsAddr.Len()
	flag := byte(0)
	if addr.CanService {
		flag = 1
	}
	buf[offset] = flag
	offset += 1
	copy(buf[offset:], addr.InternetAddr.To4())
	offset += 4

	putUint16(buf[offset:], uint16(addr.Port))
	offset += 2

	putUint16(buf[offset:], uint16(len(addr.InternalAddr)))
	offset += InternalAddrCountLen

	for _, naddr := range addr.InternalAddr {
		copy(buf[offset:], naddr.To4())
		offset += 4
	}

	putUint16(buf[offset:], uint16(len(addr.NatAddr)))
	offset += NatAddrCountLen

	for i := 0; i < len(addr.NatAddr); i++ {
		nataddr := &addr.NatAddr[i]
		offset += PackP2pAddr(nataddr, buf[offset:])
	}

	return offset

}

func UnPackCtrlMsg(buf []byte) (*OnlineCtrlMsg, int) {
	offset := 0
	ocm := &OnlineCtrlMsg{}
	ocm.typ = buf[offset]
	offset += 1
	copy(ocm.sn[:], buf[offset:])
	offset += SerialNumberBytesCount

	var nxtof int
	ocm.localAddr, nxtof = UnPackP2pAddr(buf[offset:])

	offset += nxtof

	return ocm, offset

}

func UnPackP2pAddr(buf []byte) (*P2pAddr, int) {
	offset := 0
	addr := &P2pAddr{}
	copy(addr.NbsAddr[0:], buf[offset:])
	offset += addr.NbsAddr.Len()
	flg := buf[offset]
	if flg == 1 {
		addr.CanService = true
	}
	offset += 1
	addr.InternetAddr = net.IPv4(buf[offset], buf[offset+1], buf[offset+2], buf[offset+3])
	offset += 4
	addr.Port = int(toUint16(buf[offset:]))
	offset += 2
	internalCnt := toUint16(buf[offset:])
	offset += 2
	if internalCnt > 0 {
		for i := 0; i < int(internalCnt); i++ {
			ip := net.IPv4(buf[offset], buf[offset+1], buf[offset+2], buf[offset+3])
			offset += 4
			addr.InternalAddr = append(addr.InternalAddr, ip)
		}
	}
	natCnt := toUint16(buf[offset:])
	if natCnt > 0 {
		for i := 0; i < int(natCnt); i++ {
			p2paddr, nxtoffset := UnPackP2pAddr(buf[offset:])
			offset += nxtoffset
			addr.NatAddr = append(addr.NatAddr, *p2paddr)
		}
	}

	return addr, offset
}
