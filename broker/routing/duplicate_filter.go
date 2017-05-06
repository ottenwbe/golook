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

import "sync"

type DuplicateFilter map[int]bool
type DuplicateMap struct {
	filters      map[string]DuplicateFilter
	duplicateMtx sync.Mutex
}

func newDuplicateMap() *DuplicateMap {
	return &DuplicateMap{
		filters: make(map[string]DuplicateFilter, 0),
	}
}

func (m *DuplicateMap) watchForDuplicatesFrom(system string) {
	if _, ok := m.filters[system]; !ok {
		m.filters[system] = make(DuplicateFilter, 0)
	}
}

func (m *DuplicateMap) isDuplicate(source Source) bool {
	ok := m.filters[source.System][source.Id]
	return ok
}

func (m *DuplicateMap) add(source Source) {
	m.filters[source.System][source.Id] = true
}

/*
CheckForDuplicates returns false if called multiple times the same (id and system-uuid) of a source. Otherwise it will return true.
*/
func (m *DuplicateMap) CheckForDuplicates(source Source) bool {
	m.duplicateMtx.Lock()
	defer m.duplicateMtx.Unlock()

	m.watchForDuplicatesFrom(source.System)
	result := m.isDuplicate(source)
	m.add(source)

	return result
}
