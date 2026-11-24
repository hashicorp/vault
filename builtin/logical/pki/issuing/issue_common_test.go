// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_wildHostnameRegex tests the wildHostnameRegex regular expression against a variety of
// valid and invalid hostnames.
func Test_wildHostnameRegex(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		want     bool
	}{
		// Valid hostnames without wildcards - with trailing dot
		{
			name:     "simple domain with dot",
			hostname: "example.com.",
			want:     true,
		},
		{
			name:     "subdomain with dot",
			hostname: "www.example.com.",
			want:     true,
		},
		{
			name:     "multi-level subdomain with dot",
			hostname: "api.v1.example.com.",
			want:     true,
		},
		{
			name:     "single label with dot",
			hostname: "localhost.",
			want:     true,
		},

		// Valid hostnames without wildcards - without trailing dot
		{
			name:     "simple domain without dot",
			hostname: "example.com",
			want:     true,
		},
		{
			name:     "subdomain without dot",
			hostname: "www.example.com",
			want:     true,
		},
		{
			name:     "multi-level subdomain without dot",
			hostname: "api.v1.example.com",
			want:     true,
		},
		{
			name:     "single label without dot",
			hostname: "localhost",
			want:     true,
		},
		{
			name:     "hyphenated domain",
			hostname: "my-domain.example.com",
			want:     true,
		},
		{
			name:     "numeric in domain",
			hostname: "server123.example.com",
			want:     true,
		},

		// Valid wildcards - entire label
		{
			name:     "wildcard entire first label",
			hostname: "*.example.com",
			want:     true,
		},
		{
			name:     "wildcard entire first label with dot",
			hostname: "*.example.com.",
			want:     true,
		},
		{
			name:     "wildcard entire label multi-level",
			hostname: "*.api.example.com",
			want:     true,
		},

		// Valid wildcards - start of label
		{
			name:     "wildcard at start of label",
			hostname: "*test.example.com",
			want:     true,
		},
		{
			name:     "wildcard at start with numbers",
			hostname: "*123.example.com",
			want:     true,
		},
		{
			name:     "wildcard at start alphanumeric",
			hostname: "*abc123.example.com",
			want:     true,
		},

		// Valid wildcards - end of label
		{
			name:     "wildcard at end of label",
			hostname: "test*.example.com",
			want:     true,
		},
		{
			name:     "wildcard at end with numbers",
			hostname: "123*.example.com",
			want:     true,
		},
		{
			name:     "wildcard at end alphanumeric",
			hostname: "abc123*.example.com",
			want:     true,
		},

		// Valid wildcards - middle of label
		{
			name:     "wildcard in middle of label",
			hostname: "test*server.example.com",
			want:     true,
		},
		{
			name:     "wildcard in middle with numbers",
			hostname: "api*v1.example.com",
			want:     true,
		},
		{
			name:     "wildcard in middle complex",
			hostname: "my*server.example.com",
			want:     true,
		},

		// Invalid cases - hyphens adjacent to wildcards (labels can't start/end with hyphen)
		{
			name:     "wildcard at start with hyphen after",
			hostname: "*-server.example.com",
			want:     false,
		},
		{
			name:     "wildcard at end with hyphen before",
			hostname: "server-*.example.com",
			want:     false,
		},
		{
			name:     "wildcard in middle with hyphens",
			hostname: "my-*-server.example.com",
			want:     false,
		},

		// Invalid cases - wildcard in wrong position
		{
			name:     "wildcard in second label",
			hostname: "example.*.com",
			want:     false,
		},
		{
			name:     "wildcard in last label",
			hostname: "example.com*",
			want:     false,
		},
		{
			name:     "wildcard in TLD",
			hostname: "example.*",
			want:     false,
		},

		// Invalid cases - multiple wildcards
		{
			name:     "multiple wildcards in same label",
			hostname: "*test*",
			want:     false,
		},
		{
			name:     "multiple wildcard labels",
			hostname: "*.*.example.com",
			want:     false,
		},

		// Invalid cases - empty or malformed
		{
			name:     "empty string",
			hostname: "",
			want:     false,
		},
		{
			name:     "just dot",
			hostname: ".",
			want:     false,
		},
		{
			name:     "double dots",
			hostname: "example..com",
			want:     false,
		},

		// Invalid cases - invalid characters
		{
			name:     "underscore in domain",
			hostname: "test_server.example.com",
			want:     false,
		},
		{
			name:     "space in domain",
			hostname: "test server.example.com",
			want:     false,
		},
		{
			name:     "special char in domain",
			hostname: "test@server.example.com",
			want:     false,
		},

		// Edge cases - label boundaries
		{
			name:     "label starting with hyphen",
			hostname: "-test.example.com",
			want:     false,
		},
		{
			name:     "label ending with hyphen",
			hostname: "test-.example.com",
			want:     false,
		},
		{
			name:     "single character label",
			hostname: "a.example.com",
			want:     true,
		},
		{
			name:     "single digit label",
			hostname: "1.example.com",
			want:     true,
		},

		// Edge cases - wildcard combinations
		{
			name:     "wildcard with single char after",
			hostname: "*a.example.com",
			want:     true,
		},
		{
			name:     "wildcard with single char before",
			hostname: "a*.example.com",
			want:     true,
		},
		{
			name:     "wildcard between single chars",
			hostname: "a*b.example.com",
			want:     true,
		},

		// Complex valid cases
		{
			name:     "long subdomain chain",
			hostname: "service.api.v2.prod.example.com",
			want:     true,
		},
		{
			name:     "wildcard in long chain",
			hostname: "*service.api.v2.prod.example.com",
			want:     true,
		},
		{
			name:     "alphanumeric with hyphens",
			hostname: "my-api-v2-test.example.com",
			want:     true,
		},

		// Invalid - wildcard not in leftmost label
		{
			name:     "wildcard in middle of hostname",
			hostname: "api.*.example.com",
			want:     false,
		},
		{
			name:     "wildcard at end of hostname",
			hostname: "api.example.*.com",
			want:     false,
		},

		// Edge cases - only wildcard
		{
			name:     "only wildcard",
			hostname: "*",
			want:     false,
		},
		{
			name:     "wildcard with dot",
			hostname: "*.",
			want:     false,
		},

		// Edge cases - wildcard without domain
		{
			name:     "wildcard label only",
			hostname: "*..",
			want:     false,
		},

		// Valid - wildcard with minimal domain
		{
			name:     "wildcard with single label domain",
			hostname: "*.com",
			want:     true,
		},
		{
			name:     "wildcard pattern with single label domain",
			hostname: "*test.com",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wildHostnameRegex.MatchString(tt.hostname)
			require.Equal(t, tt.want, got, "wildHostnameRegex.MatchString(%q) = %v, want %v", tt.hostname, got, tt.want)
		})
	}
}

// Made with Bob
