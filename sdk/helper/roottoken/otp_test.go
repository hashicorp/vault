// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package roottoken

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64OTPGeneration(t *testing.T) {
	token, err := GenerateOTP(0)
	assert.Len(t, token, 24)
	assert.Nil(t, err)
}

func TestBase62OTPGeneration(t *testing.T) {
	token, err := GenerateOTP(20)
	assert.Len(t, token, 20)
	assert.Nil(t, err)
}
