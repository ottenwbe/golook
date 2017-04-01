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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdClient = &cobra.Command{
	Use:   "client",
	Short: "Start as client",
	Long:  "Start as client",
	Run: func(_ *cobra.Command, _ []string) {
		//TODO
	},
}

func Host() string {
	return viper.GetString("server.host")
}

func ServerPort() int {
	return viper.GetInt("server.port")
}

func init() {
	cmdClient.Flags().StringP("host", "h", "127.0.0.1", "Url of the server (default is 127.0.0.1)")
	viper.BindPFlag("server.host", cmdClient.Flags().Lookup("host"))
	viper.SetDefault("server.host", "http://127.0.0.1")

	cmdClient.Flags().IntP("port", "p", 8080, "Port of the server (default is 8080)")
	viper.BindPFlag("server.port", cmdClient.Flags().Lookup("port"))
	viper.SetDefault("server.port", 8080)

	RootCmd.AddCommand(cmdClient)
}
