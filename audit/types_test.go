// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTypes_IsAllowedType is used to check if a type of device should be allowed
// for audit. e.g. file/socket/syslog.
func TestTypes_IsAllowedType(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    string
		expected bool
	}{
		"file-upper": {
			input:    "FILE",
			expected: true,
		},
		"file-lower": {
			input:    "file",
			expected: true,
		},
		"file-mixed": {
			input:    "fILe",
			expected: true,
		},
		"socket-upper": {
			input:    "SOCKET",
			expected: true,
		},
		"socket-lower": {
			input:    "socket",
			expected: true,
		},
		"socket-mixed": {
			input:    "sOcKeT",
			expected: true,
		},
		"syslog-upper": {
			input:    "SYSLOG",
			expected: true,
		},
		"syslog-lower": {
			input:    "syslog",
			expected: true,
		},
		"syslog-mixed": {
			input:    "sYsLoG",
			expected: true,
		},
		"something-else": {
			input:    "squiggly",
			expected: false,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := IsAllowedAuditType(tc.input)
			require.Equal(t, tc.expected, res)
		})
	}
}
