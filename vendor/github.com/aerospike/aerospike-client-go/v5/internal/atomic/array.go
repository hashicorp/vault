// Copyright 2014-2021 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atomic

import (
	"fmt"
	"sync"
)

// Array implement a fixed width array with atomic semantics
type Array struct {
	items  []interface{}
	length int
	mutex  sync.RWMutex
}

// NewArray generates a new Array instance.
func NewArray(length int) *Array {
	return &Array{
		length: length,
		items:  make([]interface{}, length),
	}
}

// Get atomically retrieves an element from the Array.
// If idx is out of range, it will return nil
func (aa *Array) Get(idx int) interface{} {
	// do not lock if not needed
	if idx < 0 || idx >= aa.length {
		return nil
	}

	aa.mutex.RLock()
	res := aa.items[idx]
	aa.mutex.RUnlock()
	return res
}

// Set atomically sets an element in the Array.
// If idx is out of range, it will return an error
func (aa *Array) Set(idx int, node interface{}) error {
	// do not lock if not needed
	if idx < 0 || idx >= aa.length {
		return fmt.Errorf("index %d is larger than array size (%d)", idx, aa.length)
	}

	aa.mutex.Lock()
	aa.items[idx] = node
	aa.mutex.Unlock()
	return nil
}

// Length returns the array size.
func (aa *Array) Length() int {
	aa.mutex.RLock()
	res := aa.length
	aa.mutex.RUnlock()

	return res
}
