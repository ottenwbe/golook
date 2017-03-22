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
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {

	initConfig()
	initDefaults()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Printf("Config file could not be found: %s \n", err)
	}
}

func initConfig() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not determine working directory: %s", err)
	}
	viper.SetConfigName("client.cfg")    // name of config file (without extension)
	viper.AddConfigPath("/etc/golook/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.golook") // call multiple times to add many search paths
	viper.AddConfigPath(wd)              // call multiple times to add many search paths
}

func initDefaults() {
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
}

func Host() string {
	return viper.GetString("server.host")
}

func ServerPort() int {
	return viper.GetInt("server.host")
}
