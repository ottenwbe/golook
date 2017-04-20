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
	. "github.com/ottenwbe/golook/broker/runtime"

	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Print the version number of %s.", GOLOOK_NAME),
	Long:  fmt.Sprintf("All software has versions. This is %s's version.", GOLOOK_NAME),
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Print(VERSION)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
