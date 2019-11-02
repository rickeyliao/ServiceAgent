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
	"github.com/rickeyliao/ServiceAgent/app/cmdservice"
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/rickeyliao/ServiceAgent/dht2"
	"github.com/rickeyliao/ServiceAgent/service"
	"github.com/rickeyliao/ServiceAgent/service/shadowsock"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	"log"
	"path"
	"github.com/rickeyliao/ServiceAgent/wifiap"
)

//var (
//	hostname *string
//	role *int64
//)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "nbssa start as daemon",
	Long:  `nbssa start as daemon`,
	Run: func(cmd *cobra.Command, args []string) {
		if !CheckProcessCanStarted() {
			return
		}

		cfg := common.GetSAConfig()
		cfg.SetNBSVersion(Version)
		daemondir := cfg.GetPidDir()
		cntxt := daemon.Context{
			PidFileName: path.Join(daemondir, "nbssa.pid"),
			PidFilePerm: 0644,
			LogFileName: path.Join(daemondir, "nbssa.log"),
			LogFilePerm: 0640,
			WorkDir:     daemondir,
			Umask:       027,
			Args:        []string{},
		}
		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if d != nil {
			log.Println("nbssa starting, please check log at:", path.Join(daemondir, "nbssa.log"))
			return
		}
		defer cntxt.Release()

		log.Println("nbssa daemon started")

		go service.Run(cfg)

		if cfg.ShadowSockServerSwitch {
			go shadowsock.StartSS2Server()
		}
		dht2.DhtRuning()
		wifiap.ExtractWifiAPFiles()
		cmdservice.StartCmdService()
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// daemonCmd.PersistentFlags().String("foo", "", "A help for foo")
	//hostname=daemonCmd.Flags().StringP("username","u","","Set User Name, Usually set to email address")
	//role=daemonCmd.Flags().Int64P("role","r",0,"Set nbssa role")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daemonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
