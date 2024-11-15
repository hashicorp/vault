// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package mallocator

// #include <stdlib.h>
// #include <string.h>
//
// void* realloc_and_initialize(void* ptr, size_t old_len, size_t new_len) {
//   void* new_ptr = realloc(ptr, new_len);
//   if (new_ptr && new_len > old_len) {
//     memset(new_ptr + old_len, 0, new_len - old_len);
//   }
//   return new_ptr;
// }
import "C"

import (
	"reflect"
	"sync/atomic"
	"unsafe"
)

// Mallocator is an allocator which defers to libc malloc.
//
// The primary reason to use this is when exporting data across the C Data
// Interface. CGO requires that pointers to Go memory are not stored in C
// memory, which is exactly what the C Data Interface would otherwise
// require. By allocating with Mallocator up front, we can safely export the
// buffers in Arrow arrays without copying buffers or violating CGO rules.
//
// The build tag 'mallocator' will also make this the default allocator.
type Mallocator struct {
	allocatedBytes uint64
}

func NewMallocator() *Mallocator { return &Mallocator{} }

func (alloc *Mallocator) Allocate(size int) []byte {
	// Use calloc to zero-initialize memory.
	// > ...the current implementation may sometimes cause a runtime error if the
	// > contents of the C memory appear to be a Go pointer. Therefore, avoid
	// > passing uninitialized C memory to Go code if the Go code is going to store
	// > pointer values in it. Zero out the memory in C before passing it to Go.
	if size < 0 {
		panic("mallocator: negative size")
	}
	ptr, err := C.calloc(C.size_t(size), 1)
	if err != nil {
		panic(err)
	} else if ptr == nil {
		panic("mallocator: out of memory")
	}
	atomic.AddUint64(&alloc.allocatedBytes, uint64(size))
	return unsafe.Slice((*byte)(ptr), size)
}

func (alloc *Mallocator) Free(b []byte) {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	C.free(unsafe.Pointer(sh.Data))
	// Subtract sh.Len via two's complement (since atomic doesn't offer subtract)
	atomic.AddUint64(&alloc.allocatedBytes, ^(uint64(sh.Len) - 1))
}

func (alloc *Mallocator) Reallocate(size int, b []byte) []byte {
	if size < 0 {
		panic("mallocator: negative size")
	}
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	ptr, err := C.realloc_and_initialize(unsafe.Pointer(sh.Data), C.size_t(sh.Cap), C.size_t(size))
	if err != nil {
		panic(err)
	} else if ptr == nil && size != 0 {
		panic("mallocator: out of memory")
	}
	delta := size - len(b)
	if delta >= 0 {
		atomic.AddUint64(&alloc.allocatedBytes, uint64(delta))
	} else {
		atomic.AddUint64(&alloc.allocatedBytes, ^(uint64(-delta) - 1))
	}
	return unsafe.Slice((*byte)(ptr), size)
}

func (alloc *Mallocator) AllocatedBytes() int64 {
	return int64(alloc.allocatedBytes)
}

// Duplicate interface to avoid circular import
type TestingT interface {
	Errorf(format string, args ...interface{})
	Helper()
}

func (alloc *Mallocator) AssertSize(t TestingT, sz int) {
	cur := alloc.AllocatedBytes()
	if int64(sz) != cur {
		t.Helper()
		t.Errorf("invalid memory size exp=%d, got=%d", sz, cur)
	}
}
