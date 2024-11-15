// Copyright (C) MongoDB, Inc. 2023-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package logger

import "context"

// contextKey is a custom type used to prevent key collisions when using the
// context package.
type contextKey string

const (
	contextKeyOperation   contextKey = "operation"
	contextKeyOperationID contextKey = "operationID"
)

// WithOperationName adds the operation name to the context.
func WithOperationName(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, contextKeyOperation, operation)
}

// WithOperationID adds the operation ID to the context.
func WithOperationID(ctx context.Context, operationID int32) context.Context {
	return context.WithValue(ctx, contextKeyOperationID, operationID)
}

// OperationName returns the operation name from the context.
func OperationName(ctx context.Context) (string, bool) {
	operationName := ctx.Value(contextKeyOperation)
	if operationName == nil {
		return "", false
	}

	return operationName.(string), true
}

// OperationID returns the operation ID from the context.
func OperationID(ctx context.Context) (int32, bool) {
	operationID := ctx.Value(contextKeyOperationID)
	if operationID == nil {
		return 0, false
	}

	return operationID.(int32), true
}
