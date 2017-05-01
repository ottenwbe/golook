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
	//com "github.com/ottenwbe/golook/broker/communication"
	//golook "github.com/ottenwbe/golook/broker/runtime"

	"github.com/spf13/viper"
)

func GetConfiguration() map[string]map[string]interface{} {

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
