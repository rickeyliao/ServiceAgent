package api
import pb "github.com/rickeyliao/ServiceAgent/app/pb"

func encResp(msg string) *pb.DefaultResp {
	resp:=&pb.DefaultResp{}
	resp.Message = msg

	return resp
}