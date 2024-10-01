// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePrefixFilters(t *testing.T) {
	t.Parallel()
	cases := []struct {
		inputFilters            []string
		expectedErrStr          string
		expectedAllowedPrefixes []string
		expectedBlockedPrefixes []string
	}{
		{
			[]string{""},
			"Cannot have empty filter rule in prefix_filter",
			[]string(nil),
			[]string(nil),
		},
		{
			[]string{"vault.abc"},
			"Filter rule must begin with either '+' or '-': \"vault.abc\"",
			[]string(nil),
			[]string(nil),
		},
		{
			[]string{"+vault.abc", "-vault.bcd"},
			"",
			[]string{"vault.abc"},
			[]string{"vault.bcd"},
		},
	}
	t.Run("validate metric filter configs", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {

			allowedPrefixes, blockedPrefixes, err := parsePrefixFilter(tc.inputFilters)

			if err != nil {
				assert.EqualError(t, err, tc.expectedErrStr)
			} else {
				assert.Equal(t, "", tc.expectedErrStr)
				assert.Equal(t, tc.expectedAllowedPrefixes, allowedPrefixes)

				assert.Equal(t, tc.expectedBlockedPrefixes, blockedPrefixes)
			}
		}
	})
}
