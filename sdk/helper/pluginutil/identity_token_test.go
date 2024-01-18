// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIdentityToken_Stringer ensures that plugin identity tokens that
// are printed in formatted strings or errors are redacted and getters
// return expected values.
func TestIdentityToken_Stringer(t *testing.T) {
	contents := "header.payload.signature"
	tk := IdentityToken(contents)

	// token getters
	assert.Equal(t, contents, tk.Token())
	assert.Equal(t, redactedTokenString, tk.String())

	// formatted strings and errors
	assert.NotContains(t, fmt.Sprintf("%v", tk), tk.Token())
	assert.NotContains(t, fmt.Sprintf("%s", tk), tk.Token())
	assert.NotContains(t, fmt.Errorf("%v", tk).Error(), tk.Token())
	assert.NotContains(t, fmt.Errorf("%s", tk).Error(), tk.Token())
}
