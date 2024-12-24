// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// TestNormalizeTelemetryAddresses ensures that any telemetry configuration that
// can be a URL, IP Address, or host:port address is conformant with RFC-5942 §4
// See: https://rfc-editor.org/rfc/rfc5952.html
func TestNormalizeTelemetryAddresses(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		given    *Telemetry
		expected *Telemetry
	}{
		"ipv6-conformance": {
			given: &Telemetry{
				// RFC-5952 4.1 leading zeroes
				CirconusAPIURL: "https://[2001:0db8::0001]:443",
				// RFC-5952 4.2.3 longest run of 0 bits shortened
				CirconusCheckSubmissionURL: "https://[2001:0:0:1:0:0:0:1]:443",
				// RFC-5952 4.2.3 equal runs of 0 bits shortened
				DogStatsDAddr: "https://[2001:db8:0:0:1:0:0:1]:443",
				// 	RFC-5952 4.3 downcase hex letters
				StatsdAddr:   "https://[2001:DB8:AC3:FE4::1]:443",
				StatsiteAddr: "https://[2001:DB8:AC3:FE4::1]:443",
			},
			expected: &Telemetry{
				CirconusAPIURL:             "https://[2001:db8::1]:443",
				CirconusCheckSubmissionURL: "https://[2001:0:0:1::1]:443",
				DogStatsDAddr:              "https://[2001:db8::1:0:0:1]:443",
				StatsdAddr:                 "https://[2001:db8:ac3:fe4::1]:443",
				StatsiteAddr:               "https://[2001:db8:ac3:fe4::1]:443",
			},
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			normalizeTelemetryAddresses(tc.given)
			require.EqualValues(t, tc.expected, tc.given)
		})
	}
}
