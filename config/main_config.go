//Copyright 2016-2017 Beate Ottenwälder
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

const (
	PROGRAM_NAME = "golook"
)

var RootCmd = &cobra.Command{
	Use:   PROGRAM_NAME,
	Short: "Golook Client/Server",
	Long:  "Golook Client/Server",
}

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Print the version number of %s", PROGRAM_NAME),
	Long:  fmt.Sprintf("All software has versions. This is %s's", PROGRAM_NAME),
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("v0.0.0")
	},
}

func Run() {
	if err := RootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Executing root command failed")
	}
}

func init() {

	initSubCommands()
	initConfig()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		const CFG_FILE string = "golook.yml"
		log.WithError(err).Infof("Config file could not be found, default is created as, %s", CFG_FILE)
		os.Create(CFG_FILE)
	}
}

func initSubCommands() {
	RootCmd.AddCommand(cmdVersion)
}

func initConfig() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not determine working directory: %s", err)
	}
	viper.SetConfigName("golook")        // name of config file (without extension)
	viper.AddConfigPath("/etc/golook/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.golook") // call multiple times to add many search paths
	viper.AddConfigPath(wd)              // call multiple times to add many search paths
}
