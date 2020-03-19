package dht2

import "strconv"

type FindReqMsg struct {
	CtrlMsg
	NodeToFind NAddr
}

func (frm *FindReqMsg)String() string  {
	s:=frm.CtrlMsg.String()

	s += " " + string(frm.NodeToFind.ID())

	return s
}

func NewReqFindMsg(msg *CtrlMsg,addr NAddr) *FindReqMsg  {
	fnm:=&FindReqMsg{}
	fnm.CtrlMsg = *msg
	fnm.NodeToFind = addr

	return fnm
}

func BuildReqFindMsg(addr NAddr) *FindReqMsg  {
	cm:=BuildMsg(Msg_Dht_Find)

	return NewReqFindMsg(cm,addr)
}

func BuildReqFindCanServiceMsg(addr NAddr) *FindReqMsg  {
	cm:=BuildMsg(Msg_CanService_Find)

	return NewReqFindMsg(cm,addr)
}

func (frm *FindReqMsg)Pack(buf []byte) int {
	cm:=&frm.CtrlMsg

	offset := PackCtrlMsg(cm,buf)

	copy(buf[offset:],frm.NodeToFind.Bytes())

	offset += frm.NodeToFind.Len()

	return offset
}

func (frm *FindReqMsg)UnPackFRM(buf []byte) int  {

	pna:=&frm.NodeToFind
	pna.Set(buf[:frm.NodeToFind.Len()])

	return frm.NodeToFind.Len()
}



type FindRespMsg struct {
	CtrlMsg
	NodeToFind NAddr
	NearestNodes []P2pAddr
}

func (frm *FindRespMsg)String() string  {
	s:=frm.CtrlMsg.String()
	s+=" "+string(frm.NodeToFind.ID())
	s+=" Total Nearest Count:"+strconv.Itoa(len(frm.NearestNodes))
	for i:=0;i<len(frm.NearestNodes);i++{
		s += "" + frm.NearestNodes[i].String()
	}

	return s
}

func NewFindRespMsg(cm *CtrlMsg,addr NAddr,nodes []P2pAddr) *FindRespMsg {
	frm:=&FindRespMsg{}
	frm.CtrlMsg = *cm
	frm.NodeToFind = addr
	frm.NearestNodes = nodes

	return frm
}

func BuildRespFindMsg(addr NAddr,nodes []P2pAddr) *FindRespMsg {
	cm:=BuildMsg(Msg_Dht_Find_Resp)

	return NewFindRespMsg(cm,addr,nodes)
}

func BuildRespFindCanServiceMsg(addr NAddr,nodes []P2pAddr) *FindRespMsg {
	cm:=BuildMsg(Msg_CanService_Find_Resp)

	return NewFindRespMsg(cm,addr,nodes)
}

func (frm *FindRespMsg)Pack(buf []byte) int{
	cm:=&frm.CtrlMsg

	offset := PackCtrlMsg(cm,buf)

	copy(buf[offset:],frm.NodeToFind.Bytes())

	offset += frm.NodeToFind.Len()

	cnt := len(frm.NearestNodes)
	offset += putUint16(buf[offset:], uint16(cnt))

	for i := 0; i < cnt; i++ {
		offset += PackP2pAddr(&frm.NearestNodes[i], buf[offset:])
	}

	return offset
}

func (frm *FindRespMsg)UnPackFRespMsg(buf []byte) int {
	pna:=&frm.NodeToFind
	pna.Set(buf[:frm.NodeToFind.Len()])

	offset:=frm.NodeToFind.Len()

	cnt := toUint16(buf[offset:])
	offset += 2

	if cnt > 0 {
		for i := 0; i < int(cnt); i++ {
			addr, of1 := UnPackP2pAddr(buf[offset:])
			offset += of1
			frm.NearestNodes = append(frm.NearestNodes, *addr)
		}
	}

	return offset
}
