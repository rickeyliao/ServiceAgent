package dht2

type NCKAResp struct {
	CtrlMsg
	RPort int
}

func NewKAResp(msg *CtrlMsg, port int) *NCKAResp {
	kar := &NCKAResp{}
	kar.CtrlMsg = *msg
	kar.RPort = port

	return kar
}

func BuildRespNCKAMsg(port int) *NCKAResp {
	cm := BuildMsg(Msg_Ka_Resp)

	return NewKAResp(cm, port)
}

func (nk *NCKAResp) Pack(buf []byte) int {
	cm := &nk.CtrlMsg
	offset := PackCtrlMsg(cm, buf)

	offset += putUint16(buf[offset:], uint16(nk.RPort))

	return offset
}

func (nk *NCKAResp) UnPackNCKA(buf []byte) int {

	offset := 0

	nk.RPort = int(toUint16(buf))

	offset += 2

	return offset
}
