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
package helper

import (
	"strings"
	"testing"
)

func TestUniquenessOfUUID(t *testing.T) {
	uuid := createAndCheckUUID(t)
	uuid2 := createAndCheckUUID(t)
	if strings.Compare(uuid, uuid2) == 0 {
		t.Error("UUIDs are equivalent, however, they should be different")
	}
}

func createAndCheckUUID(t *testing.T) string {
	uuid, err := NewUUID()
	if err != nil {
		t.Errorf("Error during uuid creation: %s", err)
	}
	return uuid
}