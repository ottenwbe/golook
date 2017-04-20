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

import "github.com/sirupsen/logrus"

type RouteLayerCallbackClient interface {
	Handle(method string, message interface{}) interface{}
	//BroadCast(v interface{}) //future work
}

type RegisteredRouteLayers struct {
	callback RouteLayerCallbackClient
}

func (r *RegisteredRouteLayers) RegisterDefaultClient(client RouteLayerCallbackClient) {
	r.callback = client
}

func ToRouteLayer(method string, message interface{}) interface{} {
	if RouteLayerRegistrar.callback != nil {
		return RouteLayerRegistrar.callback.Handle(method, message)
	} else {
		logrus.Error("Method dropped before handing over to route layer. No route layer defined.")
		return nil
	}

}

var RouteLayerRegistrar *RegisteredRouteLayers = &RegisteredRouteLayers{callback: nil}
