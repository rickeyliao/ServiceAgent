package api

import (
	"context"
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/common"
	"path"
	"strings"
)

type CmdLicenseUserServer struct {
}

func (clus *CmdLicenseUserServer) ChgLicenseUser(ctx context.Context, req *pb.LicenseUserChgReq) (*pb.DefaultResp, error) {
	userpair := strings.Split(req.User, ":")
	var user string
	var passwd string
	var passwdgen bool
	if len(userpair) == 1 {
		user = userpair[0]
		passwdgen = true
		passwd = common.GetRandPasswd(10)
	} else if len(userpair) == 2 {
		user = userpair[0]
		passwd = userpair[1]
	} else {
		return encResp("user error"), nil
	}

	if len(user) < 4 || len(passwd) < 6 {
		return encResp("error: user name length > 4, passwd length > 6"), nil
	}

	sac := common.GetSAConfig()

	lu := make([][]string, 0)
	addflag := false

	for _, v := range sac.LicenseAdminUser {
		up := make([]string, 0)

		if v[0] == user {
			if req.Op {
				if !passwdgen {
					up = append(up, user, passwd)
				} else {
					up = append(up, v[0], v[1])
				}
				addflag = true
			}
		} else {
			up = append(up, v[0], v[1])
		}
		if len(up) == 2 {
			lu = append(lu, up)
		}

	}

	if req.Op && !addflag {
		up := make([]string, 0)
		up = append(up, user, passwd)
		lu = append(lu, up)
	}

	sac.LicenseAdminUser = lu

	jstr, _ := json.MarshalIndent(sac, "", "\t")

	sar := common.GetSARootCfg()
	tools.Save2File(jstr, path.Join(sar.CfgDir, sar.CfgFileName))

	return encResp("success"), nil

}

