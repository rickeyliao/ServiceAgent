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
	"fmt"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/spf13/cobra"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	uploadhostip *string
)

// uploadCmd represents the upload command
var uploadFileCmd = &cobra.Command{
	Use:     "uploadfile",
	Aliases: []string{"uf", "upload"},
	Short:   "upload a file to a host",
	Long:    `upload a file to a host`,
	Run: func(cmd *cobra.Command, args []string) {

		if !CheckProcessReady() {
			return
		}
		if uploadhostip == nil || *uploadhostip == "" {
			fmt.Println("Please Set Host IP Address")
			return
		}

		hiparr := strings.Split(*uploadhostip, ":")
		if len(hiparr) != 2 {
			fmt.Println("Please Set Host IP Address like 1.1.1.1:50810")
			return
		}

		if net.ParseIP(hiparr[0]) == nil {
			fmt.Println("Please Set Correct Host Ip Address")
			return
		}

		if rport, err := strconv.Atoi(hiparr[1]); err != nil {
			fmt.Println("Please Set Correct Remote Port")
			return
		} else {
			if rport < 1024 || rport > 65535 {
				fmt.Println("Please Set Correct Remote Port")
				return
			}
		}

		filepath := ""

		if len(args) == 0 || len(args[0]) == 0 {
			fmt.Println("Please Set upload file")
			return
		} else {
			filepath = args[0]
		}

		if filepath[0] != '/' {
			curdir, err := os.Getwd()
			if err != nil {
				fmt.Println("Internal Error")
				return
			}

			filepath = path.Join(curdir, filepath)
		}

		if !tools.FileExists(filepath) {
			fmt.Println("File not found")
			return
		}

		UploadFileCmdSend(*uploadhostip, filepath)

	},
}

func init() {
	rootCmd.AddCommand(uploadFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	uploadhostip = uploadFileCmd.Flags().StringP("host", "s", "", "Set Peer IP Address")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//uploadFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
