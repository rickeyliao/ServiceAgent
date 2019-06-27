package api

import (
	"context"
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	pb "github.com/rickeyliao/ServiceAgent/app/pb"
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/spf13/viper"
	"path"
	"strings"
)

type CmdConfigServer struct {
}

func (ccs *CmdConfigServer) ChangeConfig(ctx context.Context, req *pb.DefaultRequestMsg) (*pb.DefaultResp, error) {
	param := strings.Split(req.Message, "=")
	if len(param) != 2 {
		return encResp("config format error"), nil
	}
	//check key,
	sar := common.GetSARootCfg()
	viper.AddConfigPath(path.Join(sar.CfgDir))
	strarr := strings.Split(sar.CfgFileName, ".")
	viper.SetConfigName(strarr[0])

	if err := viper.ReadInConfig(); err != nil {
		return encResp("read config file error"), nil
	}

	if !viper.InConfig(param[0]) {
		return encResp("no set key"), nil
	}

	key := strings.ToLower(param[0])

	viper.Set(key, param[1])

	cfg := &common.SAConfig{}
	viper.Unmarshal(cfg)

	sar.SacInst = cfg

	c := viper.AllSettings()

	s, _ := json.MarshalIndent(c, "", "\t")

	tools.Save2File(s, path.Join(sar.CfgDir, sar.CfgFileName))

	return encResp("success"), nil
}
