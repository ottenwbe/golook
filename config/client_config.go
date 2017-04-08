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
	Use:   "routing",
	Short: "Start as routing",
	Long:  "Start as a simple routing",
	Run: func(_ *cobra.Command, _ []string) {
		//TODO
	},
}

func Host() string {
	return viper.GetString("server.uplink")
}

func ServerPort() int {
	return viper.GetInt("server.port")
}

func RunDetatched() bool {
	return viper.GetBool("routing.detatch")
}

func init() {
	cmdClient.Flags().StringP("uplink", "u", "127.0.0.1", "Url of the uplink server (default is 127.0.0.1)")
	viper.BindPFlag("server.uplink", cmdClient.Flags().Lookup("uplink"))
	viper.SetDefault("server.uplink", "http://127.0.0.1")

	cmdClient.Flags().IntP("port", "p", 8080, "Port of the uplink server (default is 8080)")
	viper.BindPFlag("server.port", cmdClient.Flags().Lookup("port"))
	viper.SetDefault("server.port", 8080)

	cmdClient.Flags().BoolP("detach", "d", false, "Run as background process")
	viper.BindPFlag("routing.detatch", cmdClient.Flags().Lookup("detach"))
	viper.SetDefault("routing.detatch", false)

	RootCmd.AddCommand(cmdClient)
}
