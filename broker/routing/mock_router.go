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
	"github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/utils"
)

/*
MockRouter implements a mock of the router interface
*/
type MockRouter struct {
	Visited       int
	VisitedMethod string
	PeerName      string
}

/*
DelPeer increments the Visited counter
*/
func (mockRouter *MockRouter) DelPeer(_ Key) {
	mockRouter.Visited++
}

/*
NewPeer increments the Visited counter
*/
func (mockRouter *MockRouter) NewPeer(_ Key, peerName string) {
	mockRouter.PeerName = peerName
	mockRouter.Visited++
}

/*
BroadCast increments the Visited counter
*/
func (mockRouter *MockRouter) BroadCast(method string, params interface{}) models.EncapsulatedValues {
	mockRouter.Visited++
	mockRouter.VisitedMethod = method
	return nil
}

/*
Route increments the Visited counter
*/
func (mockRouter *MockRouter) Route(key Key, method string, params interface{}) interface{} {
	mockRouter.Visited++
	mockRouter.VisitedMethod = method
	return nil
}

/*
Handle increments the Visited counter
*/
func (mockRouter *MockRouter) Handle(method string, params models.EncapsulatedValues) interface{} {
	mockRouter.Visited++
	mockRouter.VisitedMethod = method
	return nil
}

/*
AddHandler increments the Visited counter
*/
func (mockRouter *MockRouter) AddHandler(name string, handler *Handler) {
	mockRouter.Visited++
	mockRouter.VisitedMethod = name
}

/*
Name returns a generic name: 'mock'
*/
func (mockRouter *MockRouter) Name() string {
	return "mock"
}

/*
NewMockedRouter is a factory for the MockRouter
*/
func NewMockedRouter() Router {
	return &MockRouter{}
}

/*
AccessMockedRouter is an accessor for a router r to the actual MockRouter.
Will panic if r is not a MockRouter.
*/
func AccessMockedRouter(r Router) *MockRouter {
	return r.(*MockRouter)
}

/*
RunWithMockedRouter executes a function f in a block where the given router is hidden  by a MockRouter during the execution.
The router that needs to be hidden has to be given as 'ptrOrig'.
*/
func RunWithMockedRouter(ptrOrig interface{}, f func()) {
	mockedRouter := NewMockedRouter()
	utils.Mock(ptrOrig, &mockedRouter, f)
}
