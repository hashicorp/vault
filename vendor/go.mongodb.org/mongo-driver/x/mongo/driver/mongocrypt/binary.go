// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//go:build cse
// +build cse

package mongocrypt

/*
#include <stdlib.h>
#include <mongocrypt.h>
*/
import "C"
import (
	"unsafe"
)

// binary is a wrapper type around a mongocrypt_binary_t*
type binary struct {
	p       *C.uint8_t
	wrapped *C.mongocrypt_binary_t
}

// newBinary creates an empty binary instance.
func newBinary() *binary {
	return &binary{
		wrapped: C.mongocrypt_binary_new(),
	}
}

// newBinaryFromBytes creates a binary instance from a byte buffer.
func newBinaryFromBytes(data []byte) *binary {
	if len(data) == 0 {
		return newBinary()
	}

	// TODO: Consider using runtime.Pinner to replace the C.CBytes after using go1.21.0.
	addr := (*C.uint8_t)(C.CBytes(data)) // uint8_t*
	dataLen := C.uint32_t(len(data))     // uint32_t
	return &binary{
		p:       addr,
		wrapped: C.mongocrypt_binary_new_from_data(addr, dataLen),
	}
}

// toBytes converts the given binary instance to []byte.
func (b *binary) toBytes() []byte {
	dataPtr := C.mongocrypt_binary_data(b.wrapped) // C.uint8_t*
	dataLen := C.mongocrypt_binary_len(b.wrapped)  // C.uint32_t

	return C.GoBytes(unsafe.Pointer(dataPtr), C.int(dataLen))
}

// close cleans up any resources associated with the given binary instance.
func (b *binary) close() {
	if b.p != nil {
		C.free(unsafe.Pointer(b.p))
	}
	C.mongocrypt_binary_destroy(b.wrapped)
}
