// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestResponse_ErrorResponse validates that our helper functions produce responses
// that we consider errors.
func TestResponse_ErrorResponse(t *testing.T) {
	simpleResp := ErrorResponse("a test %s", "error")
	assert.True(t, simpleResp.IsError())

	dataMap := map[string]string{
		"test1": "testing",
	}

	withDataResp := ErrorResponseWithData(dataMap, "a test %s", "error")
	assert.True(t, withDataResp.IsError())
}
