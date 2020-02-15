package dht2

const (
	Msg_Online_Req      byte = 1
	Msg_CanSrv_Req      byte = 2
	Msg_Nat_Refresh_Req byte = 3
	Msg_Nat_Conn_Req    byte = 4
	Msg_Nat_Conn_Inform byte = 5
	Msg_Nat_Sess_Create_Req byte = 6

	Msg_BS_Resp          byte = 11
	Msg_Nat_Resp         byte = 12
	Msg_CanService_Resp  byte = 13
	Msg_Nat_Refresh_Resp byte = 14
	Msg_Nat_Conn_Resp    byte = 15
	Msg_Nat_Conn_Reply   byte = 16
	Msg_Nat_Sess_Create_Resp byte = 17

	Msg_CanService_Find byte = 21
	Msg_Dht_Find        byte = 22

	Msg_CanService_Find_Resp byte = 25
	Msg_Dht_Find_Resp byte = 25

	NatServerCount int = 3

	Msg_Ka_Req byte = 31


	Msg_Ka_Resp byte = 41
)
