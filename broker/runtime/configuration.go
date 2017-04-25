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
package runtime

type (
	ConfigHandler    func()
	ConfigHandlerArr []func()
)

var (
	ConfigurationHandler = ConfigHandlerArr{}
)

func (c *ConfigHandlerArr) RegisterConfig(handler ConfigHandler) {
	*c = append(*c, handler)
}

func (c *ConfigHandlerArr) Execute() {
	for _, config := range *c {
		config()
	}
}
