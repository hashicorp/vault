// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !noasm

package memory

import "unsafe"

//go:noescape
func _memset_sse4(buf unsafe.Pointer, len, c uintptr)

func memory_memset_sse4(buf []byte, c byte) {
	if len(buf) == 0 {
		return
	}
	_memset_sse4(unsafe.Pointer(&buf[0]), uintptr(len(buf)), uintptr(c))
}
