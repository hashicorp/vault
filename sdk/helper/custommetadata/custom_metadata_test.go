// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custommetadata

import (
	"strconv"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name       string
		input      map[string]string
		shouldPass bool
	}{
		{
			"valid",
			map[string]string{
				"foo": "abc",
				"bar": "def",
				"baz": "ghi",
			},
			true,
		},
		{
			"too_many_keys",
			func() map[string]string {
				cm := make(map[string]string)

				for i := 0; i < maxKeyLength+1; i++ {
					s := strconv.Itoa(i)
					cm[s] = s
				}

				return cm
			}(),
			false,
		},
		{
			"key_too_long",
			map[string]string{
				strings.Repeat("a", maxKeyLength+1): "abc",
			},
			false,
		},
		{
			"value_too_long",
			map[string]string{
				"foo": strings.Repeat("a", maxValueLength+1),
			},
			false,
		},
		{
			"unprintable_key",
			map[string]string{
				"unprint\u200bable": "abc",
			},
			false,
		},
		{
			"unprintable_value",
			map[string]string{
				"foo": "unprint\u200bable",
			},
			false,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := Validate(tc.input)

			if tc.shouldPass && err != nil {
				t.Fatalf("expected validation to pass, input: %#v, err: %v", tc.input, err)
			}

			if !tc.shouldPass && err == nil {
				t.Fatalf("expected validation to fail, input: %#v, err: %v", tc.input, err)
			}
		})
	}
}
