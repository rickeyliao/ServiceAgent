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
	"os"

	"github.com/spf13/cobra"
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"net/http"
	"github.com/rickeyliao/ServiceAgent/agent/key"
	"github.com/rickeyliao/ServiceAgent/agent/email"
	"github.com/rickeyliao/ServiceAgent/agent/software"
	"github.com/rickeyliao/ServiceAgent/service/localaddress"
	"strconv"
	"github.com/rickeyliao/ServiceAgent/agent/listallip"
	"github.com/rickeyliao/ServiceAgent/service/postsocks5"
)


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sa",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			sar:=common.GetSARootCfg()
			if !sar.IsInitialized() {
				log.Println("Please Initialize First")
				return
			}
			//load config
			sar.LoadCfg()
			cfg:=sar.SacInst
			//if the program started, quit
			if tools.CheckPortUsed(cfg.ListenTyp,cfg.LocalListenPort){
				log.Println("sa have started")
				return
			}
			remotehost:=cfg.RemoteServerIP
			remoteport:=cfg.RemoteServerPort

			if remotehost == "" && remoteport == 0{
				log.Println("Please set remote host and port")
				return
			}

			//fmt.Println(remotehost)

			ips:=remotehost + ":" + strconv.Itoa(int(remoteport))
			common.NewRemoteUrl1(ips)

			http.Handle(cfg.VerifyPath, key.NewKeyAuth())
			http.Handle(cfg.ConsumePath,key.NewKeyImport())
			http.Handle(cfg.EmailPath,email.NewEmailRecord())
			http.Handle(cfg.UpdateClientSoftwarePath,software.NewUpdateSoft())
			http.Handle(cfg.TestIPAddress,localaddress.NewLocalAddress())
			http.Handle(cfg.ListIpsPath,listallip.NewListAllIps())
			http.Handle(cfg.PostSocks5Path,postsocks5.NewPostSocks5())

			listenportstr := ":"+strconv.Itoa(int(cfg.LocalListenPort))

			log.Println("Remote Server:",common.GetRemoteUrlInst().GetHostName(""))
			log.Println("Server Listen at:",listenportstr)

			log.Fatal(http.ListenAndServe(listenportstr, nil))

		},
}



// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "","config file (default is $HOME/.sa/config/sa.json)")
	//rootCmd.Flags().StringVarP(&sahome,"homedir","d","","home directory (default is $HOME/.sa/)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	//fmt.Println("cfgfile :",cfgFile)
	//if cfgFile != "" {
	//	// Use config file from the flag.
	//	viper.SetConfigFile(cfgFile)
	//} else {
	//	// Find home directory.
	//	_, err := homedir.Dir()
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//
	//	// Search config in home directory with name ".console" (without extension).
	//	viper.AddConfigPath("/Users/rickey/.sa/config")
	//	viper.SetConfigName("sa")
	//}
	//
	//viper.AutomaticEnv() // read in environment variables that match
	//
	//
	//fmt.Println("=====>",viper.GetString("sahome"))
	//// If a config file is found, read it in.
	//if err := viper.ReadInConfig(); err == nil {
	//	fmt.Println("Using config file:", viper.ConfigFileUsed())
	//}
	//
	//fmt.Println(viper.GetStringSlice("bootstrapipaddress"))
	//
	//var cfg common.SAConfig
	//
	//viper.Unmarshal(&cfg)
	//
	//
	//fmt.Println(cfg)
	//
	//
	//viper.Set("UploadDir","/uploaddir")
	//
	//var cfg1 common.SAConfig
	//
	//viper.Unmarshal(&cfg1)
	//
	//
	//fmt.Println(cfg1)
	//
	//
	//c:=viper.AllSettings()
	//
	//s,_:=json.Marshal(c)
	//
	//fmt.Println(string(s))
	//
	//tools.Save2File(s,"/Users/rickey/.sa/config/sa.config")


}





