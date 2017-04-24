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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var routingCmd = &cobra.Command{
	Use:   "routing",
	Short: "Configure the routing options",
	Long:  "Configure the routing options",
	Run: func(_ *cobra.Command, _ []string) {
	},
}

func configRouting() {
	log.Info("configure routing commands")
	//TODO bootstrap := viper.GetStringSlice("bootstrapping.peers")
}

func init() {
	viper.SetDefault("bootstrapping.peers", []string{"127.0.0.1"})
	RootCmd.AddCommand(routingCmd)
}
