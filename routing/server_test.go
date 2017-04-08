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
	"testing"
)

func TestRouterCreation(t *testing.T) {
	router := createRouter()

	if router == nil {
		t.Error("Router is nil after creation")
	}

	/* TODO: check if all routes are registered */
	/*if router.Get("/") == nil {
		t.Error("Route / does not exists")
	}*/
}

//TODO: test registered routes by control
// 1.) start server (in go routine)
// 2.) start a "routing" (in go routine) testing all routes
