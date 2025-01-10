// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNormalizeAddr ensures that strings that match either an IP address or URL
// and contain an IPv6 address conform to RFC-5942 §4
// See: https://rfc-editor.org/rfc/rfc5952.html
func TestNormalizeAddr(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		addr            string
		expected        string
		isErrorExpected bool
	}{
		"hostname": {
			addr:     "https://vaultproject.io:8200",
			expected: "https://vaultproject.io:8200",
		},
		"ipv4": {
			addr:     "10.10.1.10",
			expected: "10.10.1.10",
		},
		"ipv4 IP:Port addr": {
			addr:     "10.10.1.10:8500",
			expected: "10.10.1.10:8500",
		},
		"ipv4 URL": {
			addr:     "https://10.10.1.10:8200",
			expected: "https://10.10.1.10:8200",
		},
		"ipv6 IP:Port addr no brackets": {
			addr:     "2001:0db8::0001:8500",
			expected: "2001:db8::1:8500",
		},
		"ipv6 IP:Port addr with brackets": {
			addr:     "[2001:0db8::0001]:8500",
			expected: "[2001:db8::1]:8500",
		},
		"ipv6 RFC-5952 4.1 conformance leading zeroes": {
			addr:     "2001:0db8::0001",
			expected: "2001:db8::1",
		},
		"ipv6 URL RFC-5952 4.1 conformance leading zeroes": {
			addr:     "https://[2001:0db8::0001]:8200",
			expected: "https://[2001:db8::1]:8200",
		},
		"ipv6 RFC-5952 4.2.2 conformance one 16-bit 0 field": {
			addr:     "2001:db8:0:1:1:1:1:1",
			expected: "2001:db8:0:1:1:1:1:1",
		},
		"ipv6 URL RFC-5952 4.2.2 conformance one 16-bit 0 field": {
			addr:     "https://[2001:db8:0:1:1:1:1:1]:8200",
			expected: "https://[2001:db8:0:1:1:1:1:1]:8200",
		},
		"ipv6 RFC-5952 4.2.3 conformance longest run of 0 bits shortened": {
			addr:     "2001:0:0:1:0:0:0:1",
			expected: "2001:0:0:1::1",
		},
		"ipv6 URL RFC-5952 4.2.3 conformance longest run of 0 bits shortened": {
			addr:     "https://[2001:0:0:1:0:0:0:1]:8200",
			expected: "https://[2001:0:0:1::1]:8200",
		},
		"ipv6 RFC-5952 4.2.3 conformance equal runs of 0 bits shortened": {
			addr:     "2001:db8:0:0:1:0:0:1",
			expected: "2001:db8::1:0:0:1",
		},
		"ipv6 URL RFC-5952 4.2.3 conformance equal runs of 0 bits shortened": {
			addr:     "https://[2001:db8:0:0:1:0:0:1]:8200",
			expected: "https://[2001:db8::1:0:0:1]:8200",
		},
		"ipv6 RFC-5952 4.3 conformance downcase hex letters": {
			addr:     "2001:DB8:AC3:FE4::1",
			expected: "2001:db8:ac3:fe4::1",
		},
		"ipv6 URL RFC-5952 4.3 conformance downcase hex letters": {
			addr:     "https://[2001:DB8:AC3:FE4::1]:8200",
			expected: "https://[2001:db8:ac3:fe4::1]:8200",
		},
	}
	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, NormalizeAddr(tc.addr))
		})
	}
}
