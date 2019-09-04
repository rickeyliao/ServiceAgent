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
	"github.com/spf13/cobra"
	"github.com/rickeyliao/ServiceAgent/common"
	"log"
)

// updateCmd represents the update command
var cfgupdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update configuration",
	Long: `"If a secondary version flag changed, Please run this command to update the configuration. For example"`,
	Run: func(cmd *cobra.Command, args []string) {
		sar := common.GetSARootCfg()

		sar.InitConfig(false,common.InitNone,"")

		sar.SacInst.Save()

		log.Println("Config Update success!")

	},
}

func init() {
	configCmd.AddCommand(cfgupdateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
