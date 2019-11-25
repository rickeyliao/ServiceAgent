package dht2

import "net"

type RespBSMsg struct {
	CtrlMsg
	BSServer []P2pAddr
}

func NewRespBSMsg(msg *CtrlMsg, bs []P2pAddr) *RespBSMsg {
	rbsm := &RespBSMsg{}
	rbsm.CtrlMsg = *msg
	rbsm.BSServer = bs

	return rbsm
}

func BuildRespBSMsg(bs []P2pAddr) *RespBSMsg {
	cm := BuildMsg(Msg_BS_Resp)

	return NewRespBSMsg(cm, bs)
}

func (rbsm *RespBSMsg) Pack(buf []byte) int {
	cm := &rbsm.CtrlMsg

	offset := PackCtrlMsg(cm, buf)

	cnt := len(rbsm.BSServer)
	offset += putUint16(buf[offset:], uint16(cnt))

	for i := 0; i < cnt; i++ {
		offset += PackP2pAddr(&rbsm.BSServer[i], buf[offset:])
	}

	return offset
}

func (rbsm *RespBSMsg) UnpackBS(buf []byte) int {
	offset := 0

	cnt := toUint16(buf)
	offset += 2

	if cnt > 0 {
		for i := 0; i < int(cnt); i++ {
			addr, of1 := UnPackP2pAddr(buf[offset:])
			offset += of1
			rbsm.BSServer = append(rbsm.BSServer, *addr)
		}
	}

	return offset
}

type RespNatMsg struct {
	CtrlMsg
	CanService bool
	ObservrIP  net.IP
	NatServer  []P2pAddr
}

func NewRespNatMsg(msg *CtrlMsg, can bool, obsip net.IP, nats []P2pAddr) *RespNatMsg {
	rnm := &RespNatMsg{}
	rnm.CtrlMsg = *msg
	rnm.CanService = can
	rnm.ObservrIP = obsip
	rnm.NatServer = nats

	return rnm
}

func BuildRespNatMsg(can bool, obsip net.IP, nats []P2pAddr) *RespNatMsg {

	cm := BuildMsg(Msg_Nat_Resp)

	return NewRespNatMsg(cm, can, obsip, nats)
}

func (rnm *RespNatMsg) Pack(buf []byte) int {
	cm := &rnm.CtrlMsg

	offset := PackCtrlMsg(cm, buf)

	buf[offset] = func() byte {
		if rnm.CanService {
			return 1
		} else {
			return 0
		}
	}()

	offset += 1

	copy(buf[offset:], rnm.ObservrIP.To4())
	offset += 4

	cnt := len(rnm.NatServer)
	offset += putUint16(buf[offset:], uint16(cnt))

	for i := 0; i < cnt; i++ {
		offset += PackP2pAddr(&rnm.NatServer[i], buf[offset:])
	}

	return offset
}

func (rnm *RespNatMsg) UnpackNatS(buf []byte) int {
	offset := 0

	rnm.CanService = func() bool {
		if buf[offset] == 1 {
			return true
		} else {
			return false
		}
	}()

	offset += 1

	rnm.ObservrIP = net.IPv4(buf[offset], buf[offset+1], buf[offset+2], buf[offset+3])
	offset += 4

	cnt := toUint16(buf)
	offset += 2

	if cnt > 0 {
		for i := 0; i < int(cnt); i++ {
			addr, of1 := UnPackP2pAddr(buf[offset:])
			offset += of1
			rnm.NatServer = append(rnm.NatServer, *addr)
		}
	}

	return offset
}
