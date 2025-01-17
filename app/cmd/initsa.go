// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/spf13/cobra"
	"log"
)

var sahome string
var force bool
var iscoord bool
var machinename string
var shadowsockParam string
var nationality int32
var iscanservice bool
var internetipstr string

//var sslocation string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init nbssa start environment",
	Long:  `init nbssa start environment`,
	Run: func(cmd *cobra.Command, args []string) {
		if nationality == 0 {
			log.Println("Please Set Nationality, 1 American, 86 China Main Land")
			return
		}

		sar := common.GetSARootCfgHdir(sahome, force)
		cmbool := common.InitFalse
		if iscoord {
			cmbool = common.InitTrue
		}

		cip := &common.ConfigInitParam{}
		cip.Force = force
		cip.IsCoord = cmbool
		cip.Hostname = machinename
		cip.Nationality = nationality
		cip.SS = shadowsockParam
		cip.Ip = internetipstr
		cip.IsCanService = iscanservice

		//cip.Location = sslocation

		sar.InitConfig(cip)
		sar.InitRSAKey(force)
		sar.LoadRsaKey()

		if shadowsockParam != "" {
			sar.SetShadowSockParam(shadowsockParam)
		} else {
			if sar.SacInst.ShadowSockServerSwitch {
				sar.SacInst.ShadowSockServerSwitch = false
				sar.SacInst.Save()
			}
		}

		log.Println("Config initialized")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&sahome, "homedir", "d", "", "home directory (default is $HOME/.sa/)")
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "force init a default configuration ( default false)")
	initCmd.Flags().StringVarP(&shadowsockParam, "ss", "s", "", "configuration shadowsock passwd and method")
	initCmd.Flags().BoolVarP(&iscoord, "coordinator", "c", false, "current machine is a coordinator")
	initCmd.Flags().StringVarP(&machinename, "machinename", "m", "", "set machine name")
	//initCmd.Flags().StringVarP(&sslocation, "location", "l", "", "set shadowsock server vps location")
	initCmd.Flags().Int32VarP(&nationality, "nationality", "n", 0, "set nationality")
	initCmd.Flags().StringVarP(&internetipstr, "ipaddress", "i", "", "set internet ip address")
	initCmd.Flags().BoolVarP(&iscanservice, "canservice", "r", false, "set can service value")
}
