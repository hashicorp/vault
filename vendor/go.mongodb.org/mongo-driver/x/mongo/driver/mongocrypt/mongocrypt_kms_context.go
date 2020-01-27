// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// +build cse

package mongocrypt

// #include <mongocrypt.h>
import "C"

// KmsContext represents a mongocrypt_kms_ctx_t handle.
type KmsContext struct {
	wrapped *C.mongocrypt_kms_ctx_t
}

// newKmsContext creates a KmsContext wrapper around the given C type.
func newKmsContext(wrapped *C.mongocrypt_kms_ctx_t) *KmsContext {
	return &KmsContext{
		wrapped: wrapped,
	}
}

// HostName gets the host name of the KMS.
func (kc *KmsContext) HostName() (string, error) {
	var hostname *C.char // out param for mongocrypt function to fill in hostname
	if ok := C.mongocrypt_kms_ctx_endpoint(kc.wrapped, &hostname); !ok {
		return "", kc.createErrorFromStatus()
	}
	return C.GoString(hostname), nil
}

// Message returns the message to send to the KMS.
func (kc *KmsContext) Message() ([]byte, error) {
	msgBinary := newBinary()
	defer msgBinary.close()

	if ok := C.mongocrypt_kms_ctx_message(kc.wrapped, msgBinary.wrapped); !ok {
		return nil, kc.createErrorFromStatus()
	}
	return msgBinary.toBytes(), nil
}

// BytesNeeded returns the number of bytes that should be received from the KMS.
// After sending the message to the KMS, this message should be called in a loop until the number returned is 0.
func (kc *KmsContext) BytesNeeded() int32 {
	return int32(C.mongocrypt_kms_ctx_bytes_needed(kc.wrapped))
}

// FeedResponse feeds the bytes received from the KMS to mongocrypt.
func (kc *KmsContext) FeedResponse(response []byte) error {
	responseBinary := newBinaryFromBytes(response)
	defer responseBinary.close()

	if ok := C.mongocrypt_kms_ctx_feed(kc.wrapped, responseBinary.wrapped); !ok {
		return kc.createErrorFromStatus()
	}
	return nil
}

// createErrorFromStatus creates a new Error from the status of the KmsContext instance.
func (kc *KmsContext) createErrorFromStatus() error {
	status := C.mongocrypt_status_new()
	defer C.mongocrypt_status_destroy(status)
	C.mongocrypt_kms_ctx_status(kc.wrapped, status)
	return errorFromStatus(status)
}
