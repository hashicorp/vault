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

import "sync/atomic"

// Bool implements a synchronized boolean value
type Bool struct {
	val int32
}

// NewBool generates a new Boolean instance.
func NewBool(value bool) *Bool {
	var i int32
	if value {
		i = 1
	}
	return &Bool{
		val: i,
	}
}

// Get atomically retrieves the boolean value.
func (ab *Bool) Get() bool {
	return atomic.LoadInt32(&(ab.val)) != 0
}

// Set atomically sets the boolean value.
func (ab *Bool) Set(newVal bool) {
	var i int32
	if newVal {
		i = 1
	}
	atomic.StoreInt32(&(ab.val), int32(i))
}

// Or atomically applies OR operation to the boolean value.
func (ab *Bool) Or(newVal bool) bool {
	if !newVal {
		return ab.Get()
	}
	atomic.StoreInt32(&(ab.val), int32(1))
	return true
}

//CompareAndToggle atomically sets the boolean value if the current value is equal to updated value.
func (ab *Bool) CompareAndToggle(expect bool) bool {
	updated := 1
	if expect {
		updated = 0
	}
	res := atomic.CompareAndSwapInt32(&ab.val, int32(1-updated), int32(updated))
	return res
}
