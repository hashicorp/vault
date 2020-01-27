// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// +build cse

package mongocrypt

// #include <mongocrypt.h>
import "C"
import (
	"fmt"
)

// Error represents an error from an operation on a MongoCrypt instance.
type Error struct {
	Code    int32
	Message string
}

// Error implements the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("mongocrypt error %d: %v", e.Code, e.Message)
}

// errorFromStatus builds a Error from a mongocrypt_status_t object.
func errorFromStatus(status *C.mongocrypt_status_t) error {
	cCode := C.mongocrypt_status_code(status) // uint32_t
	// mongocrypt_status_message takes uint32_t* as its second param to store the length of the returned string.
	// pass nil because the length is handled by C.GoString
	cMsg := C.mongocrypt_status_message(status, nil) // const char*
	var msg string
	if cMsg != nil {
		msg = C.GoString(cMsg)
	}

	return Error{
		Code:    int32(cCode),
		Message: msg,
	}
}
