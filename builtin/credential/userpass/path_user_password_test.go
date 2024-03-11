// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package userpass

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestUserPass_ParseHash ensures that we correctly validate password hashes that
// conform to the bcrypt standard based on the prefix of the hash.
func TestUserPass_ParseHash(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input           string
		isErrorExpected bool
	}{
		"spaces": {
			input:           "         ",
			isErrorExpected: true,
		},
		"jibberish": {
			input:           "jibberfish",
			isErrorExpected: true,
		},
		"non-ascii": {
			input:           "$2a$qwerty",
			isErrorExpected: false,
		},
		"truncation": {
			input:           "$2b$qwerty",
			isErrorExpected: false,
		},
		"php-only-fixed": {
			input:           "$2y$qwerty",
			isErrorExpected: false,
		},
		"php-only-existing": {
			input:           "$2x$qwerty",
			isErrorExpected: true,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := parsePasswordHash(tc.input)
			switch {
			case tc.isErrorExpected:
				require.EqualError(t, err, "\"password_hash\" doesn't appear to be a valid bcrypt hash")
			default:
				require.NoError(t, err)
				require.Equal(t, tc.input, string(got))
			}
		})
	}
}
