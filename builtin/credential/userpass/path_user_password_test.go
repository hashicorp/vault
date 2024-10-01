// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package userpass

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// TestUserPass_ParseHash ensures that we correctly validate password hashes that
// conform to the bcrypt standard based on the prefix of the hash.
func TestUserPass_ParseHash(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input                string
		isErrorExpected      bool
		expectedErrorMessage string
	}{
		"too-short": {
			input:                "too short",
			isErrorExpected:      true,
			expectedErrorMessage: "password hash has incorrect length",
		},
		"60-spaces": {
			input:                "                                                            ",
			isErrorExpected:      true,
			expectedErrorMessage: "password hash has incorrect prefix",
		},
		"jibberish": {
			input:                "jibberfishjibberfishjibberfishjibberfishjibberfishjibberfish",
			isErrorExpected:      true,
			expectedErrorMessage: "password hash has incorrect prefix",
		},
		"non-ascii-prefix": {
			input:           "$2a$qwertyjibberfishjibberfishjibberfishjibberfishjibberfish",
			isErrorExpected: false,
		},
		"truncation-prefix": {
			input:           "$2b$qwertyjibberfishjibberfishjibberfishjibberfishjibberfish",
			isErrorExpected: false,
		},
		"php-only-fixed-prefix": {
			input:           "$2y$qwertyjibberfishjibberfishjibberfishjibberfishjibberfish",
			isErrorExpected: false,
		},
		"php-only-existing": {
			input:                "$2x$qwertyjibberfishjibberfishjibberfishjibberfishjibberfish",
			isErrorExpected:      true,
			expectedErrorMessage: "password hash has incorrect prefix",
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
				require.EqualError(t, err, tc.expectedErrorMessage)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.input, string(got))
			}
		})
	}
}

// TestUserPass_BcryptHashLength ensures that using the bcrypt library to generate
// a hash from a password always produces the same length.
func TestUserPass_BcryptHashLength(t *testing.T) {
	t.Parallel()

	tests := []string{
		"",
		"    ",
		"foo",
		"this is a long password woo",
	}

	for _, input := range tests {
		hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
		require.NoError(t, err)
		require.Len(t, hash, bcryptHashLength)
	}
}
