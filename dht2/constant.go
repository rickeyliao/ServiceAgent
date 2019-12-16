package dht2

const (
	Msg_Online_Req byte = 1
	Msg_CanSrv_Req byte = 2
	Msg_Nat_Refresh_Req    byte = 3

	Msg_BS_Resp         byte = 11
	Msg_Nat_Resp        byte = 12
	Msg_CanService_Resp byte = 13
	Msg_Nat_Refresh_Resp byte = 14

	Msg_CanService_Loop byte = 21
	Msg_Dht_Loop        byte = 22

	NatServerCount int = 3

	Msg_Ka_Req byte = 31

	Msg_Ka_Resp byte = 41
)
