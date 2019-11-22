package api

import (
	"context"
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/common"
	"net"
	"path"
	"strings"
)

type CmdBootstrapServer struct {
}

type bootstrapaddr struct {
	nbsaddr string
	ipaddr  string
}

func (bsa *bootstrapaddr) String() string {
	return bsa.nbsaddr + "@" + bsa.ipaddr
}

func (cbs *CmdBootstrapServer) ChangeBootstrap(ctx context.Context, req *pb.BootstrapCHGReq) (*pb.DefaultResp, error) {

	arrparam := strings.Split(req.Address, "@")
	if len(arrparam) != 2 {
		return encResp("address format error"), nil
	}
	nbsaddr := arrparam[0]
	if !common.CheckNbsNodeHash(nbsaddr) {
		return encResp("node address error"), nil
	}

	ipaddr := arrparam[1]
	if net.ParseIP(ipaddr) == nil {
		return encResp("ip address error"), nil
	}

	sac := common.GetSAConfig()

	addrarr := sac.BootstrapIPAddress

	addrs := make([]*bootstrapaddr, 0)

	addflag := false

	for _, v := range addrarr {
		addrarr := strings.Split(v, "@")

		if len(addrarr) != 2 {
			continue
		}

		bsa := &bootstrapaddr{}
		bsa.nbsaddr = addrarr[0]
		bsa.ipaddr = addrarr[1]

		if req.Op {
			if bsa.nbsaddr == nbsaddr {
				bsa.ipaddr = ipaddr
				addflag = true
			}
		} else {
			//remove
			if bsa.nbsaddr == nbsaddr && bsa.ipaddr == ipaddr {
				continue
			}
		}
		addrs = append(addrs, bsa)
	}
	if req.Op && !addflag {
		bsa := &bootstrapaddr{}
		bsa.nbsaddr = nbsaddr
		bsa.ipaddr = ipaddr
		addrs = append(addrs, bsa)
	}

	straddrs := make([]string, 0)

	for _, bsa := range addrs {
		straddrs = append(straddrs, bsa.String())
	}

	sac.BootstrapIPAddress = straddrs

	jstr, _ := json.MarshalIndent(sac, "", "\t")

	sar := common.GetSARootCfg()
	tools.Save2File(jstr, path.Join(sar.CfgDir, sar.CfgFileName))

	return encResp("success"), nil

}
