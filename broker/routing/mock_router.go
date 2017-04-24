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
	"github.com/ottenwbe/golook/broker/communication"
	. "github.com/ottenwbe/golook/broker/models"
	. "github.com/ottenwbe/golook/broker/utils"
)

type MockedLookRouter struct {
	Visited       int
	VisitedMethod string
}

func (*MockedLookRouter) NewNeighbor(key Key, neighbor communication.LookupClient) {
	panic("implement me")
}

func (lr *MockedLookRouter) BroadCast(method string, params interface{}) interface{} {
	lr.Visited += 1
	lr.VisitedMethod = method
	return nil
}

func (lr *MockedLookRouter) Route(key Key, method string, params interface{}) interface{} {
	lr.Visited += 1
	lr.VisitedMethod = method
	return nil
}

func (lr *MockedLookRouter) Handle(method string, params MsgParams) interface{} {
	lr.Visited += 1
	lr.VisitedMethod = method
	return nil
}

func (lr *MockedLookRouter) HandlerFunction(name string, handler func(params interface{}) interface{}) {
	lr.Visited += 1
	lr.VisitedMethod = name
}

func (lr *MockedLookRouter) Name() string {
	return "mock"
}

func NewMockedRouter() Router {
	return &MockedLookRouter{}
}

func AccessMockedRouter(r Router) *MockedLookRouter {
	return r.(*MockedLookRouter)
}

func RunWithMockedRouter(ptrOrig interface{}, f func()) {
	mockedRouter := NewMockedRouter()
	Mock(ptrOrig, &mockedRouter, f)
}
