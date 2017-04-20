////Copyright 2016-2017 Beate Ottenwälder
////
////Licensed under the Apache License, Version 2.0 (the "License");
////you may not use this file except in compliance with the License.
////You may obtain a copy of the License at
////
////http://www.apache.org/licenses/LICENSE-2.0
////
////Unless required by applicable law or agreed to in writing, software
////distributed under the License is distributed on an "AS IS" BASIS,
////WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
////See the License for the specific language governing permissions and
////limitations under the License.
package routing

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type MockedLookRouter struct {
	Visited       int
	VisitedMethod string
}

func (lr *MockedLookRouter) Route(key Key, method string, params interface{}) interface{} {
	lr.Visited += 1
	lr.VisitedMethod = method
	return nil
}

func (lr *MockedLookRouter) Handle(method string, params interface{}) interface{} {
	lr.Visited += 1
	lr.VisitedMethod = method
	return nil
}

func (lr *MockedLookRouter) HandlerFunction(name string, handler func(params interface{}) interface{}) {
	lr.Visited += 1
	lr.VisitedMethod = name
}

func NewMockedRouter() Router {
	return &MockedLookRouter{}
}

func AccessMockedRouter() *MockedLookRouter {
	return GoLookRouter.(*MockedLookRouter)
}

var mockMutex = &sync.Mutex{}

func RunWithMockedRouter(f func()) Router {
	mockMutex.Lock()
	// ensure that router is reset
	defer func(tmpRouter Router) {
		GoLookRouter = tmpRouter
		mockMutex.Unlock()
	}(GoLookRouter)
	GoLookRouter = NewMockedRouter()

	f()

	if r := recover(); r != nil {
		logrus.Errorf("Recovered in RunWithMockedRouter: %s", r)
	}

	return GoLookRouter
}
