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

//TODO: expire keys

import "sync"

var maxDuplicateMapLen = 100

type duplicateFilter map[int]bool
type duplicateMap struct {
	filters      map[string]duplicateFilter
	duplicateMtx sync.Mutex
}

func newDuplicateMap() *duplicateMap {
	return &duplicateMap{
		filters: make(map[string]duplicateFilter, 0),
	}
}

func (m *duplicateMap) watchForDuplicatesFrom(system string) {
	if _, ok := m.filters[system]; !ok {
		m.filters[system] = make(duplicateFilter, 0)
	}
}

func (m *duplicateMap) isDuplicate(source Source) bool {
	ok := m.filters[source.System][source.ID]
	return ok
}

func (m *duplicateMap) add(source Source) {
	m.filters[source.System][source.ID] = true
}

/*
prune removes entries from the duplicate map.
This implementation is very simple and restricts the history of potential duplicates by length.
In turn, this does not ensure that each and every duplicate is detected.
*/
func (m *duplicateMap) prune(source Source) {
	if len(m.filters[source.System]) >= maxDuplicateMapLen {
		for k := range m.filters[source.System] {
			delete(m.filters[source.System], k)
			break
		}
	}
}

/*
CheckForDuplicates returns false if called multiple times the same (id and system-uuid) of a source. Otherwise it will return true.
*/
func (m *duplicateMap) CheckForDuplicates(source Source) bool {
	m.duplicateMtx.Lock()
	defer m.duplicateMtx.Unlock()

	m.watchForDuplicatesFrom(source.System)
	m.prune(source)
	result := m.isDuplicate(source)
	m.add(source)

	return result
}
