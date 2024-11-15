// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//go:build cse
// +build cse

package mongocrypt

// #include <mongocrypt.h>
import "C"
import (
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// Context represents a mongocrypt_ctx_t handle
type Context struct {
	wrapped *C.mongocrypt_ctx_t
}

// newContext creates a Context wrapper around the given C type.
func newContext(wrapped *C.mongocrypt_ctx_t) *Context {
	return &Context{
		wrapped: wrapped,
	}
}

// State returns the current State of the Context.
func (c *Context) State() State {
	return State(int(C.mongocrypt_ctx_state(c.wrapped)))
}

// NextOperation gets the document for the next database operation to run.
func (c *Context) NextOperation() (bsoncore.Document, error) {
	opDocBinary := newBinary() // out param for mongocrypt_ctx_mongo_op to fill in operation
	defer opDocBinary.close()

	if ok := C.mongocrypt_ctx_mongo_op(c.wrapped, opDocBinary.wrapped); !ok {
		return nil, c.createErrorFromStatus()
	}
	return opDocBinary.toBytes(), nil
}

// AddOperationResult feeds the result of a database operation to mongocrypt.
func (c *Context) AddOperationResult(result bsoncore.Document) error {
	resultBinary := newBinaryFromBytes(result)
	defer resultBinary.close()

	if ok := C.mongocrypt_ctx_mongo_feed(c.wrapped, resultBinary.wrapped); !ok {
		return c.createErrorFromStatus()
	}
	return nil
}

// CompleteOperation signals a database operation has been completed.
func (c *Context) CompleteOperation() error {
	if ok := C.mongocrypt_ctx_mongo_done(c.wrapped); !ok {
		return c.createErrorFromStatus()
	}
	return nil
}

// NextKmsContext returns the next KmsContext, or nil if there are no more.
func (c *Context) NextKmsContext() *KmsContext {
	ctx := C.mongocrypt_ctx_next_kms_ctx(c.wrapped)
	if ctx == nil {
		return nil
	}
	return newKmsContext(ctx)
}

// FinishKmsContexts signals that all KMS contexts have been completed.
func (c *Context) FinishKmsContexts() error {
	if ok := C.mongocrypt_ctx_kms_done(c.wrapped); !ok {
		return c.createErrorFromStatus()
	}
	return nil
}

// Finish performs the final operations for the context and returns the resulting document.
func (c *Context) Finish() (bsoncore.Document, error) {
	docBinary := newBinary() // out param for mongocrypt_ctx_finalize to fill in resulting document
	defer docBinary.close()

	if ok := C.mongocrypt_ctx_finalize(c.wrapped, docBinary.wrapped); !ok {
		return nil, c.createErrorFromStatus()
	}
	return docBinary.toBytes(), nil
}

// Close cleans up any resources associated with the given Context instance.
func (c *Context) Close() {
	C.mongocrypt_ctx_destroy(c.wrapped)
}

// createErrorFromStatus creates a new Error based on the status of the MongoCrypt instance.
func (c *Context) createErrorFromStatus() error {
	status := C.mongocrypt_status_new()
	defer C.mongocrypt_status_destroy(status)
	C.mongocrypt_ctx_status(c.wrapped, status)
	return errorFromStatus(status)
}

// ProvideKmsProviders provides the KMS providers when in the NeedKmsCredentials state.
func (c *Context) ProvideKmsProviders(kmsProviders bsoncore.Document) error {
	kmsProvidersBinary := newBinaryFromBytes(kmsProviders)
	defer kmsProvidersBinary.close()

	if ok := C.mongocrypt_ctx_provide_kms_providers(c.wrapped, kmsProvidersBinary.wrapped); !ok {
		return c.createErrorFromStatus()
	}
	return nil
}
