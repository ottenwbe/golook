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

package routing

import (
	com "github.com/ottenwbe/golook/broker/communication"
	log "github.com/sirupsen/logrus"
)

/*
RouterType represents the type of a router
*/
type RouterType string

const (
	/*BroadcastRouterType represent the type of the default broadcast router*/
	BroadcastRouterType RouterType = "broadcast"
	/*MockRouterType represent the type of the mocked router*/
	MockRouterType RouterType = "mock"
)

var (
	defaultPeers = []string{}
)

/*
NewRouter creates a new router of type "routerType" with a unique human readable name to identify the router.
Note, however, that a new router cannot receive messages, yet. To this end, the router needs to activated; see ActivateRouter(r Router).
*/
func NewRouter(name string, routerType RouterType) (result Router) {

	switch routerType {
	case MockRouterType:
		result = NewMockedRouter()
	case BroadcastRouterType:
		result = newBroadcastRouter(name)
	default:
		log.WithField("type", routerType).Error("Cannot instantiate new router for unknown type.")
		return nil
	}

	for _, peer := range defaultPeers {
		result.NewPeer(NewKey(peer), peer)
	}

	return result
}

/*
ActivateRouter enables the router to receive messages.
NOTE: To clean up after the router is no longer used DeactivateRouter(r Router) has to be called
*/
func ActivateRouter(r Router) {
	com.MessageDispatcher.RegisterHandler(r.Name(), r, RequestMessage{}, ResponseMessage{})
}

/*
DeactivateRouter performs all cleanup operations for the router
*/
func DeactivateRouter(r Router) {
	com.MessageDispatcher.RemoveHandler(r.Name())
}
