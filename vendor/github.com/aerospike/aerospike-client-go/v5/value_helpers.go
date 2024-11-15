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

package aerospike

// MapIter allows to define general maps of your own type to be used in the Go client
// without the use of reflection.
// function PackMap should be exactly Like the following (Do not change, just copy/paste and adapt PackXXX methods):
//  func (cm *CustomMap) PackMap(buf aerospike.BufferEx) (int, error) {
//  	size := 0
//  	for k, v := range cm {
//  		n, err := PackXXX(buf, k)
//  		size += n
//  		if err != nil {
//  			return size, err
//  		}
//
//  		n, err = PackXXX(buf, v)
//  		size += n
//  		if err != nil {
//  			return size, err
//  		}
//  	}
//  	return size, nil
//  }
type MapIter interface {
	PackMap(buf BufferEx) (int, error)
	Len() int
}

// ListIter allows to define general maps of your own type to be used in the Go client
// without the use of reflection.
// function PackList should be exactly Like the following (Do not change, just copy/paste and adapt PackXXX methods):
//  func (cs *CustomSlice) PackList(buf aerospike.BufferEx) (int, error) {
//  	size := 0
//  	for _, elem := range cs {
//  		n, err := PackXXX(buf, elem)
//  		size += n
//  		if err != nil {
//  			return size, err
//  		}
//  	}
//  	return size, nil
//  }
type ListIter interface {
	PackList(buf BufferEx) (int, error)
	Len() int
}
