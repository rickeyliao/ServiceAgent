package dht2

import (
	"net"
	"encoding/binary"
	"strconv"
)

type NCKAResp struct {
	CtrlMsg
	RPort int
}

func (nckar *NCKAResp)String() string  {
	s:=nckar.CtrlMsg.String()
	s += " Remote Port: "+strconv.Itoa(nckar.RPort)
	return s
}

func NewKAResp(msg *CtrlMsg, port int) *NCKAResp {
	kar := &NCKAResp{}
	kar.CtrlMsg = *msg
	kar.RPort = port

	return kar
}

func BuildRespNCKAMsg(port int,sn []byte) *NCKAResp {
	cm := BuildMsgWithSN(Msg_Ka_Resp,sn)

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


type NCConnReq struct {
	CtrlMsg
	Wait4ConnNode *P2pAddr
}

func (nccr *NCConnReq)String() string{
	s:=nccr.String()
	s += " Want to connect to: "+ nccr.Wait4ConnNode.String()
	return s
}

func NewNCConnReq(cm *CtrlMsg,w4n *P2pAddr) *NCConnReq {
	nccr:=&NCConnReq{}

	nccr.CtrlMsg = *cm
	nccr.Wait4ConnNode = w4n

	return nccr
}

func BuildNCConnReq(w4n *P2pAddr) *NCConnReq {
	cm:=BuildMsg(Msg_Nat_Conn_Req)

	return NewNCConnReq(cm,w4n)
}


func (nc *NCConnReq) Pack(buf []byte) int {
	cm := &nc.CtrlMsg
	offset := PackCtrlMsg(cm, buf)

	offset +=PackP2pAddr(nc.Wait4ConnNode,buf[offset:])

	return offset
}

func (nc *NCConnReq) UnPack(buf []byte) int  {
	offset:=0
	nc.Wait4ConnNode,offset = UnPackP2pAddr(buf[offset:])

	return offset
}

var(
	NCConnSuccess int = 0
	NCConnNotFound int = 1
)

type NCConnResp struct {
	CtrlMsg
	ErrCode int
	RemoteIP net.IP
	RemotePort int
}

func (ncr *NCConnResp)String() string  {
	msg:=ncr.CtrlMsg.String()
	msg += " RemoteIP:" + ncr.RemoteIP.String()
	msg += " RemotePort: " + strconv.Itoa(ncr.RemotePort)
	msg += " ErrCode : " + strconv.Itoa(ncr.ErrCode)

	return msg

}


func NewNCConnResp(cm *CtrlMsg,errCode int,remoteIP net.IP,remotePort int) *NCConnResp {
	ncr:=&NCConnResp{}
	ncr.CtrlMsg = *cm
	ncr.ErrCode = errCode
	ncr.RemoteIP = remoteIP
	ncr.RemotePort = remotePort

	return ncr
}

func BuildNCConnResp(errCode int,remoteIP net.IP,remotePort int,sn []byte) *NCConnResp {
	cm:=BuildMsgWithSN(Msg_Nat_Conn_Resp,sn)

	return NewNCConnResp(cm,errCode,remoteIP,remotePort)
}

func (nc *NCConnResp)Pack(buf []byte) int  {
	cm:=&nc.CtrlMsg

	offset := PackCtrlMsg(cm,buf)

	binary.BigEndian.PutUint32(buf[offset:],uint32(nc.ErrCode))
	offset += 4

	offset += copy(buf[offset:],nc.RemoteIP.To4())

	offset += putUint16(buf[offset:],uint16(nc.RemotePort))
	return offset
}

func (nc *NCConnResp)UnPack(buf []byte) int  {
	offset:=0
	nc.ErrCode = int(binary.BigEndian.Uint32(buf[offset:]))
	offset += 4
	nc.RemoteIP = net.IPv4(buf[offset], buf[offset+1], buf[offset+2], buf[offset+3])
	offset += 4

	nc.RemotePort = int(toUint16(buf[offset:]))

	return offset + 2
}

type NCConnInForm struct {
	CtrlMsg
	Wait4ConnIP net.IP
	Wait4ConnPort int
}

func (nccif *NCConnInForm)String() string  {
	s:=nccif.CtrlMsg.String()
	s+="wait for conn ip:"+nccif.Wait4ConnIP.String()
	s+=" port:"+strconv.Itoa(nccif.Wait4ConnPort)

	return s
}

func NewNCConnInform(cm *CtrlMsg,w4ip net.IP,w4port int) *NCConnInForm {
	nci:=&NCConnInForm{}
	nci.CtrlMsg = *cm
	nci.Wait4ConnIP = w4ip
	nci.Wait4ConnPort = w4port

	return nci
}

func BuildNCConnInform(w4ip net.IP,w4port int)  *NCConnInForm {
	cm:=BuildMsg(Msg_Nat_Conn_Inform)

	return NewNCConnInform(cm,w4ip,w4port)
}

func (nc *NCConnInForm) Pack(buf []byte) int{
	cm:=&nc.CtrlMsg

	offset := PackCtrlMsg(cm,buf)

	offset += copy(buf[offset:],nc.Wait4ConnIP.To4())

	offset += putUint16(buf[offset:],uint16(nc.Wait4ConnPort))
	return offset
}

func (nc *NCConnInForm) UnPack(buf []byte) int{
	offset:=0

	nc.Wait4ConnIP = net.IPv4(buf[offset], buf[offset+1], buf[offset+2], buf[offset+3])
	offset += 4

	nc.Wait4ConnPort = int(toUint16(buf[offset:]))

	return offset + 2
}

type NCConnReply struct {
	CtrlMsg
}

func BuildNCConnReply(sn []byte) *NCConnReply {
	cm:=BuildMsgWithSN(Msg_Nat_Conn_Reply,sn)

	return &NCConnReply{*cm}
}

func (nc *NCConnReply)Pack(buf []byte) int {
	return PackCtrlMsg(&nc.CtrlMsg,buf)
}

func (nc *NCConnReply)UnPack(buf []byte) int {
	return 0
}

type NCSessionCreateReq struct {
	CtrlMsg
}

func BuildNCSessCreateReq() *NCSessionCreateReq {
	cm:=BuildMsg(Msg_Nat_Sess_Create_Req)

	return &NCSessionCreateReq{*cm}
}

func (nc *NCSessionCreateReq)Pack(buf []byte) int {
	return PackCtrlMsg(&nc.CtrlMsg,buf)
}

func (nc *NCSessionCreateReq)UnPack(buf []byte) int {
	return 0
}

type NCSessionCreateResp struct {
	CtrlMsg
}

func BuildNCSessCreateResp(sn []byte) *NCSessionCreateResp {
	cm:=BuildMsgWithSN(Msg_Nat_Sess_Create_Resp,sn)

	return &NCSessionCreateResp{*cm}
}

func (nc *NCSessionCreateResp)Pack(buf []byte) int {
	return PackCtrlMsg(&nc.CtrlMsg,buf)
}

func (nc *NCSessionCreateResp)UnPack(buf []byte) int {
	return 0
}
