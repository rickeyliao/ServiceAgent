package api

import (
	"context"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/service/license"
)

type CmdOpLicenseSrv struct {
}

func (chs *CmdOpLicenseSrv) OpLicense(cxt context.Context, hi *pb.LicenseReq) (*pb.DefaultResp, error) {
	if hi.Op == 0 && hi.Sofaaddress != "" {
		return encResp(license.CmdLicenseShow(hi.Sofaaddress)), nil
	}
	switch hi.Op {
	case 0:
		return encResp(license.CmdShowLicenseStatistic()), nil
	case 1:
		return encResp(license.CmdShowLicenseAll()), nil
	case 2:
		license.Save()
		return encResp("save success"), nil
	case 3:
		return encResp(license.CmdShowLicenseSummary()), nil
	}

	return encResp("not found cmd"), nil
}
