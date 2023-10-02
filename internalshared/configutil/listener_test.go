// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestParseSingleIPTemplate exercises the ParseSingleIPTemplate function to
// ensure that we only attempt to parse templates when the input contains a
// template placeholder (see: go-sockaddr/template).
func TestParseSingleIPTemplate(t *testing.T) {
	tests := map[string]struct {
		arg             string
		want            string
		isErrorExpected bool
		errorMessage    string
	}{
		"test https addr": {
			arg:             "https://vaultproject.io:8200",
			want:            "https://vaultproject.io:8200",
			isErrorExpected: false,
		},
		"test invalid template func": {
			arg:             "{{ FooBar }}",
			want:            "",
			isErrorExpected: true,
			errorMessage:    "unable to parse address template",
		},
		"test partial template": {
			arg:             "{{FooBar",
			want:            "{{FooBar",
			isErrorExpected: false,
		},
	}
	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			got, err := ParseSingleIPTemplate(tc.arg)

			if tc.isErrorExpected {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errorMessage)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.want, got)
		})
	}
}
