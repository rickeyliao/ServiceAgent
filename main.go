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

package main

import (
	"errors"
	"github.com/rickeyliao/ServiceAgent/app/cmd"
	"github.com/rickeyliao/ServiceAgent/common"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	Version   string
	Build     string
	BuildTime string
	WifiRes   string
)

func main() {
	cmd.Version = Version
	cmd.Build = Build
	cmd.BuildTime = BuildTime
	cmd.RunPath, _ = getCurrentPath()

	common.WifiRes = WifiRes

	cmd.Execute()
}
func getCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New("error: can't find / or \\ ")
	}
	return string(path[0 : i+1]), nil
}
