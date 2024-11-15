// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//go:build !cse
// +build !cse

package mongocrypt

import (
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// Context represents a mongocrypt_ctx_t handle
type Context struct{}

// State returns the current State of the Context.
func (c *Context) State() State {
	panic(cseNotSupportedMsg)
}

// NextOperation gets the document for the next database operation to run.
func (c *Context) NextOperation() (bsoncore.Document, error) {
	panic(cseNotSupportedMsg)
}

// AddOperationResult feeds the result of a database operation to mongocrypt.
func (c *Context) AddOperationResult(bsoncore.Document) error {
	panic(cseNotSupportedMsg)
}

// CompleteOperation signals a database operation has been completed.
func (c *Context) CompleteOperation() error {
	panic(cseNotSupportedMsg)
}

// NextKmsContext returns the next KmsContext, or nil if there are no more.
func (c *Context) NextKmsContext() *KmsContext {
	panic(cseNotSupportedMsg)
}

// FinishKmsContexts signals that all KMS contexts have been completed.
func (c *Context) FinishKmsContexts() error {
	panic(cseNotSupportedMsg)
}

// Finish performs the final operations for the context and returns the resulting document.
func (c *Context) Finish() (bsoncore.Document, error) {
	panic(cseNotSupportedMsg)
}

// Close cleans up any resources associated with the given Context instance.
func (c *Context) Close() {
	panic(cseNotSupportedMsg)
}

// ProvideKmsProviders provides the KMS providers when in the NeedKmsCredentials state.
func (c *Context) ProvideKmsProviders(bsoncore.Document) error {
	panic(cseNotSupportedMsg)
}
