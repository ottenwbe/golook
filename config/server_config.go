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
	"github.com/ottenwbe/golook/server"
	"github.com/spf13/cobra"
)

var addr string

var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "Start as server",
	Long:  "Start as server",
	Run: func(_ *cobra.Command, _ []string) {
		server.StartServer(addr)
	},
}

func init() {
	cmdServer.Flags().StringVar(&addr, "address", ":8080", "Address of the server (default is :8080)")
	RootCmd.AddCommand(cmdServer)
}
