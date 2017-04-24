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
package cmd

import (
	"os"

	. "github.com/ottenwbe/golook/broker/runtime"

	"github.com/ottenwbe/golook/broker/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Run() {

	//TODO refactor
	configRouting()
	api.RegisterApi()

	if err := RootCmd.Execute(); err != nil {
		log.WithError(err).Panic("Executing root command failed")
	}
}

var RootCmd = &cobra.Command{
	Use:   GOLOOK_NAME,
	Short: "Golook Broker",
	Long:  "Golook Broker which implements a Servent (Client/Server) for a distributed file search",
	Run: func(_ *cobra.Command, _ []string) {
		log.Info("Starting up Golook...")
		HttpServer.StartServer()
		log.Info("Shutting down server...")
	},
}

func init() {

	initConfig()

	err := viper.ReadInConfig() // Find and read the cmd file
	if err != nil {             // Handle errors reading the cmd file
		log.WithError(err).Infof("Config file could not be found, falling back to default parameters")
	}
}

func initConfig() {
	wd, err := os.Getwd()
	if err != nil {
		log.WithError(err).Fatalf("Could not determine working directory")
	}
	viper.SetConfigName("golook")        // name of cmd file (without extension)
	viper.AddConfigPath("/etc/golook/")  // path to look for the cmd file in
	viper.AddConfigPath("$HOME/.golook") // call multiple times to add many search paths
	viper.AddConfigPath(wd)              // call multiple times to add many search paths
}
