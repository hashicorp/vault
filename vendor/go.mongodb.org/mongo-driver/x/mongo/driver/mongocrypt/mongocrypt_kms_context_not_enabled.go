// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//go:build !cse
// +build !cse

package mongocrypt

// KmsContext represents a mongocrypt_kms_ctx_t handle.
type KmsContext struct{}

// HostName gets the host name of the KMS.
func (kc *KmsContext) HostName() (string, error) {
	panic(cseNotSupportedMsg)
}

// Message returns the message to send to the KMS.
func (kc *KmsContext) Message() ([]byte, error) {
	panic(cseNotSupportedMsg)
}

// KMSProvider gets the KMS provider of the KMS context.
func (kc *KmsContext) KMSProvider() string {
	panic(cseNotSupportedMsg)
}

// BytesNeeded returns the number of bytes that should be received from the KMS.
// After sending the message to the KMS, this message should be called in a loop until the number returned is 0.
func (kc *KmsContext) BytesNeeded() int32 {
	panic(cseNotSupportedMsg)
}

// FeedResponse feeds the bytes received from the KMS to mongocrypt.
func (kc *KmsContext) FeedResponse([]byte) error {
	panic(cseNotSupportedMsg)
}
