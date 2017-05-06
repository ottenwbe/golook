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

package utils

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"reflect"
)

var mockMutex = &sync.Mutex{}

//Idea based on http://stackoverflow.com/questions/19301742/golang-interface-to-swap-two-numbers
func swap(s1 interface{}, s2 interface{}) {
	ts1 := reflect.ValueOf(s1).Elem()
	ts2 := reflect.ValueOf(s2).Elem()
	tmp := ts1.Interface()

	ts1.Set(ts2)
	ts2.Set(reflect.ValueOf(tmp))
}

/*
Mock allows a function f to be executed in a block. Within this block a mocked value 'mock' is replacing the value of 'orig'.
After the function has executed, Mock ensures that the original value is written back to 'orig'.
Any panic caused by f is captured and ignored.

A prerequisite for using Mock is that orig and mock implement the same interface.
*/
func Mock(orig interface{}, mock interface{}, f func()) {
	mockMutex.Lock()

	// ensure that the original interface is reset after function
	defer func() {
		defer mockMutex.Unlock()

		if rec := recover(); rec != nil {
			log.Errorf("Recovered in Mock: %s", rec)
		}

		log.Debugf("Replacing back %s to %s ", reflect.ValueOf(orig).Elem(), reflect.ValueOf(mock).Elem())
		swap(orig, mock)
		log.Debugf("Replaced back %s to %s ", reflect.ValueOf(orig).Elem(), reflect.ValueOf(mock).Elem())
	}()

	log.Debugf("Will swap %s with %s", reflect.ValueOf(orig).Elem(), reflect.ValueOf(mock).Elem())
	swap(orig, mock)
	log.Debugf("Replaced %s to %s", reflect.ValueOf(orig).Elem(), reflect.ValueOf(mock).Elem())
	f()
}
