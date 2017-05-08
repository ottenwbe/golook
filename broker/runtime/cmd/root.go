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
	"github.com/ottenwbe/golook/broker/api"
	"github.com/ottenwbe/golook/broker/communication"
	"github.com/ottenwbe/golook/broker/repository"
	"github.com/ottenwbe/golook/broker/routing"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	"github.com/ottenwbe/golook/broker/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

/*
Run is executed to start the program
*/
func Run() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Panic("Executing root command failed")
	}
}

var rootCmd = &cobra.Command{
	Use:   golook.GolookName,
	Short: "Golook Broker",
	Long:  "Golook Broker which implements a Servent (Client/Server) for a distributed file search",
	Run: func(_ *cobra.Command, _ []string) {

		applyConfiguration()

		log.Infof("Starting up Golook System: %s", golook.GolookSystem.UUID)

		golook.RunServer()

		log.Info("Shutting down server...")
	},
}

func init() {
	service.GetLogService().Init()

	api.InitConfiguration()
	repositories.InitConfiguration()
	service.InitServiceConfiguration()
	routing.InitConfiguration()
	communication.InitCommunicationConfiguration()
}

func applyConfiguration() {

	service.GetLogService().Apply()

	api.ApplyConfiguration()
	repositories.ApplyConfiguration()
	service.ApplyServiceConfiguration()
	routing.ApplyConfiguration()
	communication.ApplyCommunicationConfiguration()
}
