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
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop a nbssa daemon",
	Long:  `stop a nbssa daemon`,
	Run: func(cmd *cobra.Command, args []string) {
		if !CheckProcessReady() {
			return
		}
		DefaultCmdSend(app.CMD_STOP_REQ)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
