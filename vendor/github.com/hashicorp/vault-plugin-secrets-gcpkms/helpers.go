// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/jeffchao/backoff"
)

// withFieldValidator wraps an OperationFunc and validates the user-supplied
// fields match the schema.
func withFieldValidator(f framework.OperationFunc) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		if err := validateFields(req, d); err != nil {
			return nil, logical.CodedError(400, err.Error())
		}
		return f(ctx, req, d)
	}
}

// validateFields verifies that no bad arguments were given to the request.
func validateFields(req *logical.Request, data *framework.FieldData) error {
	var unknownFields []string
	for k := range req.Data {
		if _, ok := data.Schema[k]; !ok {
			unknownFields = append(unknownFields, k)
		}
	}

	switch len(unknownFields) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("unknown field: %s", unknownFields[0])
	default:
		sort.Strings(unknownFields)
		return fmt.Errorf("unknown fields: %s", strings.Join(unknownFields, ","))
	}
}

// errMissingFields is a helper to return an error when required fields are
// missing.
func errMissingFields(f ...string) error {
	return logical.CodedError(400, fmt.Sprintf(
		"missing required field(s): %q", f))
}

// retryFib accepts a function and retries using a fibonacci algorithm.
func retryFib(op func() error) error {
	f := backoff.Fibonacci()
	f.Interval = 100 * time.Millisecond
	f.MaxRetries = 10
	return f.Retry(op)
}

// retryExp accepts a function and retries using an exponential backoff
// algorithm.
func retryExp(op func() error) error {
	f := backoff.Exponential()
	f.Interval = 100 * time.Millisecond
	f.MaxRetries = 10
	return f.Retry(op)
}
