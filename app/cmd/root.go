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

	"github.com/rickeyliao/ServiceAgent/app/cmdservice"
	"github.com/rickeyliao/ServiceAgent/common"
	"github.com/rickeyliao/ServiceAgent/service"
	"github.com/spf13/cobra"

	"github.com/rickeyliao/ServiceAgent/dht2"
	"github.com/rickeyliao/ServiceAgent/service/shadowsock"
	"github.com/rickeyliao/ServiceAgent/wifiap"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nbssa",
	Short: "start nbssa in command line",
	Long:  `start nbssa in command line`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if !CheckProcessCanStarted() {
			return
		}

		cfg := common.GetSAConfig()
		cfg.SetNBSVersion(Version)

		go service.Run(cfg)
		if cfg.ShadowSockServerSwitch {
			go shadowsock.StartSS2Server()
		}
		dht2.DhtRuning()
		wifiap.ExtractWifiAPFiles()
		cmdservice.StartCmdService()
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
