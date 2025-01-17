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
	"github.com/rickeyliao/ServiceAgent/app"
	"github.com/spf13/cobra"
	"log"
	"net"
)

// removeCmd represents the remove command
var ssremoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove remote server",
	Long:  `remove remote server`,
	Run: func(cmd *cobra.Command, args []string) {
		if !CheckProcessReady() {
			return
		}

		ip := ""

		if len(args) > 0 {
			if len(args[0]) > 0 {
				if _, err := net.ResolveIPAddr("ip", args[0]); err != nil {
					log.Println("ip address error", args[0], err)
					return
				}
				ip = args[0]
			}
		}

		SServerCmdSend(app.CMD_SSSERVER_REMOVE, ssserverNationality, false, ip, "")
	},
}

func init() {
	ssserverCmd.AddCommand(ssremoveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
