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

//go:build !go1.20 && !tinygo

package hashing

import (
	"reflect"
	"unsafe"
)

func hashString(val string, alg uint64) uint64 {
	if val == "" {
		return Hash([]byte{}, alg)
	}
	// highly efficient way to get byte slice without copy before
	// the introduction of unsafe.StringData in go1.20
	// (https://stackoverflow.com/questions/59209493/how-to-use-unsafe-get-a-byte-slice-from-a-string-without-memory-copy)
	const MaxInt32 = 1<<31 - 1
	buf := (*[MaxInt32]byte)(unsafe.Pointer((*reflect.StringHeader)(
		unsafe.Pointer(&val)).Data))[: len(val)&MaxInt32 : len(val)&MaxInt32]
	return Hash(buf, alg)
}
