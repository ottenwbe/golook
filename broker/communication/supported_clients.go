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
package communication

import log "github.com/sirupsen/logrus"

var LookupClients *SupportedClients = newCommunicationClients()

type LookupClientBuilder (func(url string) LookupClient)

type SupportedClients struct {
	defaultClient string
	clients       map[string]LookupClientBuilder
}

func (c *SupportedClients) BuildDefault(url string) LookupClient {
	if result, ok := c.clients[c.defaultClient]; !ok {
		log.Panic("Default communication client has not been selected")
		return nil
	} else {
		return result(url)
	}
}

func (c *SupportedClients) Build(url string, name string) (LookupClient, bool) {
	if result, ok := c.clients[name]; ok {
		return result(url), ok
	} else {
		return nil, ok
	}
}

func (c *SupportedClients) Add(name string, client LookupClientBuilder, isDefault bool) {

	if isDefault {
		if c.defaultClient != "" {
			log.WithField("existing", c.defaultClient).WithField("new", name).Panicf("Only one default client allowed. %s.", name)
		}
		c.defaultClient = name
	}

	if _, ok := c.clients[name]; ok {
		log.Panicf("Duplicate communication client, denoted %s, is not allowed.", name)
	} else {
		c.clients[name] = client
	}
}

func SupportedCommunicationClients() []string {
	result := make([]string, 0)
	for key := range LookupClients.clients {
		result = append(result, key)
	}
	return result
}

func newCommunicationClients() *SupportedClients {
	return &SupportedClients{
		defaultClient: "",
		clients:       make(map[string]LookupClientBuilder),
	}
}
