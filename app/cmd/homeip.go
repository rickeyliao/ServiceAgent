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
)

// homeipCmd represents the homeip command
var homeipCmd = &cobra.Command{
	Use:   "homeip",
	Short: "show home ip",
	Long: `show home ip`,
	Run: func(cmd *cobra.Command, args []string) {
		if !CheckProcessReady() {
			return
		}
		nbsaddr:=""
		if len(args) > 0{
			nbsaddr=args[0]
		}
		HomeIPShowCmdSend(nbsaddr)
	},
}

func init() {
	rootCmd.AddCommand(homeipCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// homeipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// homeipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
