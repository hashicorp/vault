// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package random

import (
	"testing"
)

func TestCharset(t *testing.T) {
	type testCase struct {
		charset  string
		minChars int
		input    string
		expected bool
	}

	tests := map[string]testCase{
		"0 minimum, empty input": {
			charset:  LowercaseCharset,
			minChars: 0,
			input:    "",
			expected: true,
		},
		"0 minimum, many matching": {
			charset:  LowercaseCharset,
			minChars: 0,
			input:    LowercaseCharset,
			expected: true,
		},
		"0 minimum, no matching": {
			charset:  LowercaseCharset,
			minChars: 0,
			input:    "0123456789",
			expected: true,
		},
		"1 minimum, empty input": {
			charset:  LowercaseCharset,
			minChars: 1,
			input:    "",
			expected: false,
		},
		"1 minimum, no matching": {
			charset:  LowercaseCharset,
			minChars: 1,
			input:    "0123456789",
			expected: false,
		},
		"1 minimum, exactly 1 matching": {
			charset:  LowercaseCharset,
			minChars: 1,
			input:    "a",
			expected: true,
		},
		"1 minimum, many matching": {
			charset:  LowercaseCharset,
			minChars: 1,
			input:    "abcdefhaaaa",
			expected: true,
		},
		"2 minimum, 1 matching": {
			charset:  LowercaseCharset,
			minChars: 2,
			input:    "f",
			expected: false,
		},
		"2 minimum, 2 matching": {
			charset:  LowercaseCharset,
			minChars: 2,
			input:    "fz",
			expected: true,
		},
		"2 minimum, many matching": {
			charset:  LowercaseCharset,
			minChars: 2,
			input:    "joixnbonxd",
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cr := CharsetRule{
				Charset:  []rune(test.charset),
				MinChars: test.minChars,
			}
			actual := cr.Pass([]rune(test.input))
			if actual != test.expected {
				t.FailNow()
			}
		})
	}
}
