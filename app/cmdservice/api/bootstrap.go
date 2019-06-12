package api

import (
	"context"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/common"
	"strings"
	"strconv"
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	"path"
)

type CmdBootstrapServer struct {

}


type bootstrapaddr struct {
	addr string
	port uint16
}

func (bsa *bootstrapaddr)String() string {
	return bsa.addr +":" + strconv.Itoa(int(bsa.port))
}


func (cbs *CmdBootstrapServer)ChangeBootstrap(ctx context.Context, req *pb.BootstrapCHGReq) (*pb.DefaultResp, error){

	arrparam:=strings.Split(req.Address,":")
	if len(arrparam)!=2 {
		return encResp("address error"),nil
	}
	paramip:=arrparam[0]

	var p int
	var e error

	if p,e=strconv.Atoi(arrparam[1]);e!=nil{
		return encResp("address error"),nil
	}
	paramport :=uint16(p)

	sac:=common.GetSAConfig()

	addrarr:=sac.BootstrapIPAddress

	addrs := make([]*bootstrapaddr,0)

	addflag:=false

	for _,v:=range addrarr{
		ipport:=strings.Split(v,":")

		if len(ipport) != 2{
			continue
		}

		bsa:=&bootstrapaddr{}
		bsa.addr = ipport[0]
		if port,err := strconv.Atoi(ipport[1]);err!=nil{
			continue
		}else {
			bsa.port = uint16(port)
		}

		if req.Op {
			if bsa.addr == paramip {
				bsa.port = paramport
				addflag = true
			}
		}else{
			//remove
			if bsa.port == paramport && bsa.addr == paramip{
				continue
			}
		}
		addrs = append(addrs,bsa)
	}
	if req.Op && !addflag{
		bsa:=&bootstrapaddr{}
		bsa.port = paramport
		bsa.addr = paramip
		addrs = append(addrs,bsa)
	}

	straddrs:=make([]string,0)

	for _,bsa:=range addrs{
		straddrs = append(straddrs,bsa.String())
	}

	sac.BootstrapIPAddress = straddrs

	jstr,_:=json.MarshalIndent(sac,"","\t")

	sar:=common.GetSARootCfg()
	tools.Save2File(jstr,path.Join(sar.CfgDir,sar.CfgFileName))

	return encResp("success"),nil

}