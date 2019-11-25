package dht2

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
