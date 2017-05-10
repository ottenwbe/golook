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

import "github.com/ottenwbe/golook/broker/models"

/*
MockClient implements a mock for the interface RPCClient.
MockClient also records the calls to the individual functions of the interface.
*/
type MockClient struct {
	VisitedCall int
	VisitedUrl  int
	Name        string
}

func newMockClient() RPCClient {
	return &MockClient{}
}

/*
URL always returns "test" and ensures that the counter 'VisitedUrl' is increased.
*/
func (client *MockClient) URL() string {
	client.VisitedUrl++
	return "test"
}

/*
Call always returns "nil, nil", but ensures that the counter 'VisitedCall' is increased and the called handler's 'name' is recorded in 'Name'.
*/
func (client *MockClient) Call(handler string, message interface{}) (models.EncapsulatedValues, error) {
	client.Name = handler
	client.VisitedCall++
	return nil, nil
}
