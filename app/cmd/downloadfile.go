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

	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/spf13/cobra"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	downloadip *string
	outputpath *string
)

// downloadfileCmd represents the downloadfile command
var downloadfileCmd = &cobra.Command{
	Use:   "downloadfile",
	Short: "download a file from server",
	Long:  `download a file from server`,
	Run: func(cmd *cobra.Command, args []string) {

		if !CheckProcessReady() {
			return
		}

		if len(args) == 0 || len(args[0]) == 0 {
			fmt.Println("Please Set download file")
			return
		}

		if !common.CheckNbsCotentHash(args[0]) {
			fmt.Println("File hash not Correct")
			return
		}

		if downloadip == nil || *downloadip == "" {
			fmt.Println("Please Set Host IP Address")
			return
		}

		hiparr := strings.Split(*downloadip, ":")
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

		outputpathstr := ""
		curdir, err := os.Getwd()
		if err != nil {
			fmt.Println("Internal error")
			return
		}
		if outputpath != nil && len(*outputpath) > 0 {
			outputpathstr = *outputpath
			if outputpathstr[0] != '/' {
				outputpathstr = path.Join(curdir, outputpathstr)
			}
		} else {
			outputpathstr = curdir
		}

		DownloadFileCmdSend(*downloadip, args[0], outputpathstr)

	},
}

func init() {
	rootCmd.AddCommand(downloadfileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadfileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	outputpath = downloadfileCmd.Flags().StringP("outputpath", "o", "", "Output path")
	downloadip = downloadfileCmd.Flags().StringP("host", "s", "", "server ip")
}
