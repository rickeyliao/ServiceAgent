package dht2

const (
	Msg_Online_Req      byte = 1			//online request
	Msg_CanSrv_Req      byte = 2			//test peer is a can service request
	Msg_Nat_Refresh_Req byte = 3			//if node is a not can service node,
											// the node must have 3 nat servers, if less than 3
											// node will send request to get more nat server
	Msg_Nat_Conn_Req    byte = 4			//Send request to Peer's Nat server,
											// Nat Server will Inform peer to connect back.
	Msg_Nat_Conn_Inform byte = 5			// Nat Server inform the node, a connection come.
	Msg_Nat_Sess_Create_Req byte = 6		//Send request to Create a connection session with peer node

	Msg_BS_Resp          byte = 11          //response Msg_Online_Req ,
	                                        // if peer is not can service node or use private ip address
	Msg_Nat_Resp         byte = 12			//response Msg_Online_Req
											//if request node is not a can service node, tell the node 3 nat server,
											//if request node is a can service node, response it
	Msg_CanService_Resp  byte = 13			//reponse Msg_CanSrv_Req
	Msg_Nat_Refresh_Resp byte = 14			//response Msg_Nat_Refresh_Req
	Msg_Nat_Conn_Resp    byte = 15			//response Msg_Nat_Conn_Req
	Msg_Nat_Conn_Reply   byte = 16			//reply Msg_Nat_Conn_Inform
	Msg_Nat_Sess_Create_Resp byte = 17		//response Msg_Nat_Sess_Create_Req

	Msg_CanService_Find byte = 21
	Msg_Dht_Find        byte = 22

	Msg_CanService_Find_Resp byte = 25
	Msg_Dht_Find_Resp byte = 25



	Msg_Ka_Req byte = 31					//node send ka msg to nat server


	Msg_Ka_Resp byte = 41					//reponse Msg_Ka_Req
)

const(
	NatServerCount int = 3

	DHTNearstCount int = 8
	DHTFindA       int = 3
)