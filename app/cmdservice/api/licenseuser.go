package api

import (
	"context"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"strings"
	"crypto/rand"
	"crypto/sha1"
	"github.com/btcsuite/btcutil/base58"
	"github.com/rickeyliao/ServiceAgent/common"
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	"path"
)

type CmdLicenseUserServer struct {

}

func (clus *CmdLicenseUserServer)ChgLicenseUser(ctx context.Context, req *pb.LicenseUserChgReq) (*pb.DefaultResp, error)  {
	userpair := strings.Split(req.User,":")
	var user string
	var passwd string
	var passwdgen bool
	if len(userpair) == 1{
		user = userpair[0]
		passwdgen = true
		passwd = getRandPasswd()
	}else if len(userpair) == 2{
		user = userpair[0]
		passwd = userpair[1]
	}else {
		return encResp("user error"),nil
	}

	if len(user) < 4 && len(passwd) < 6 {
		return encResp("error: user name length > 4, passwd length > 6"),nil
	}

	sac:=common.GetSAConfig()

	lu:=make([][]string,0)
	addflag:=false

	for _,v:=range sac.LicenseAdminUser{
		up:=make([]string,0)

		if v[0] == user {
			if req.Op {
				if !passwdgen{
					up = append(up,user,passwd)
				}else{
					up = append(up,v[0],v[1])
				}
				addflag = true
			}
		}else{
			up = append(up,v[0],v[1])
		}
		if len(up) == 2{
			lu = append(lu,up)
		}

	}

	if req.Op && !addflag {
		up:=make([]string,0)
		up = append(up,user,passwd)
		lu = append(lu,up)
	}

	sac.LicenseAdminUser = lu

	jstr, _ := json.MarshalIndent(sac, "", "\t")

	sar := common.GetSARootCfg()
	tools.Save2File(jstr, path.Join(sar.CfgDir, sar.CfgFileName))

	return encResp("success"), nil


}

func getRandPasswd() string  {
	buf:=make([]byte,256)
	for{
		n,err:=rand.Read(buf[:])
		if err!=nil{
			continue
		}

		if n <len(buf){
			continue
		}

		break
	}

	s1:=sha1.New()
	s1.Write(buf)
	b:=s1.Sum(nil)


	passwd := base58.Encode(b)
	l := len(passwd)
	if l > 10{
		return passwd[l-10:l]
	}

	return passwd
}