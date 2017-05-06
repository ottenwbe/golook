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

package service

import (
	"github.com/ottenwbe/golook/broker/runtime"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

/*
ConfigurationService implements the means to read the configuration
*/
type ConfigurationService interface {
	GetConfiguration() map[string]map[string]interface{}
}

type viperConfiguration struct{}

/*
NewConfigurationService returns by default a configuration service based on the viper configuration solution.
*/
func NewConfigurationService() ConfigurationService {
	return &viperConfiguration{}
}

func (*viperConfiguration) GetConfiguration() map[string]map[string]interface{} {

	var (
		configurations = map[string]map[string]interface{}{}
		configGeneral  = map[string]interface{}{}
	)

	configGeneral["file"] = viper.ConfigFileUsed()
	configGeneral["settings"] = viper.AllSettings()
	configGeneral["keys"] = viper.AllKeys()
	configurations["config"] = configGeneral

	return configurations
}

func ApplyServiceConfiguration() {
	runtime.Schedule(viper.GetString("service.informer.specification"), *newSystemService())
	OpenFileServices(FileServiceType(viper.GetString("service.type")))
}

func InitServiceConfiguration() {
	viper.SetDefault("service.type", string(BroadcastFiles))
	viper.SetDefault("service.informer.specification", "@every 5m0s")

	wd, err := os.Getwd()
	if err != nil {
		logrus.WithError(err).Fatalf("Cannot determine working directory")
	}
	viper.SetConfigName("golook")        // name of cmd file (without extension)
	viper.AddConfigPath("/etc/golook/")  // path to look for the cmd file in
	viper.AddConfigPath("$HOME/.golook") // call multiple times to add many search paths
	viper.AddConfigPath(wd)              // call multiple times to add many search paths
	viper.AddConfigPath(wd + "/config")  // call multiple times to add many search paths

	err = viper.ReadInConfig() // Find and read the cmd file
	if err != nil {            // Handle errors reading the cmd file
		log.WithError(err).Infof("Config file cannot be found, falling back to default parameters")
	}
}
