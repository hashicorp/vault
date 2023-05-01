// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build testonly

package vault

import (
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestSystemBackend_handleActivityWriteData calls the activity log write endpoint and confirms that the inputs are
// correctly validated
func TestSystemBackend_handleActivityWriteData(t *testing.T) {
	testCases := []struct {
		name      string
		operation logical.Operation
		input     map[string]interface{}
		wantError error
	}{
		{
			name:      "read fails",
			operation: logical.ReadOperation,
			wantError: logical.ErrUnsupportedOperation,
		},
		{
			name:      "empty write fails",
			operation: logical.CreateOperation,
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "wrong key fails",
			operation: logical.CreateOperation,
			input:     map[string]interface{}{"other": "data"},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "incorrectly formatted data fails",
			operation: logical.CreateOperation,
			input:     map[string]interface{}{"input": "data"},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "incorrect json data fails",
			operation: logical.CreateOperation,
			input:     map[string]interface{}{"input": `{"other":"json"}`},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "empty write value fails",
			operation: logical.CreateOperation,
			input:     map[string]interface{}{"input": `{"write":[],"data":[]}`},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "empty data value fails",
			operation: logical.CreateOperation,
			input:     map[string]interface{}{"input": `{"write":["WRITE_PRECOMPUTED_QUERIES"],"data":[]}`},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "correctly formatted data succeeds",
			operation: logical.CreateOperation,
			input:     map[string]interface{}{"input": `{"write":["WRITE_PRECOMPUTED_QUERIES"],"data":[{"current_month":true,"all":{"clients":[{"count":5}]}}]}`},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := testSystemBackend(t)
			req := logical.TestRequest(t, tc.operation, "internal/counters/activity/write")
			req.Data = tc.input
			resp, err := b.HandleRequest(namespace.RootContext(nil), req)
			if tc.wantError != nil {
				require.Equal(t, tc.wantError, err, resp.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
