// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package jwt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNormalizeIssuer verifies RFC-compliant issuer normalization per RFC 9207 and RFC 3986.
// Only scheme and host are lowercased; path case is preserved.
func TestNormalizeIssuer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		issuer   string
		expected string
	}{
		// Basic normalization
		{
			name:     "scheme and host normalization",
			issuer:   "HTTPS://EXAMPLE.COM",
			expected: "https://example.com",
		},
		{
			name:     "path case preservation",
			issuer:   "https://example.com/REALM/Team",
			expected: "https://example.com/REALM/Team",
		},
		{
			name:     "scheme host and path normalization",
			issuer:   "HTTPS://EXAMPLE.COM/realm/TEAM",
			expected: "https://example.com/realm/TEAM",
		},
		{
			name:     "trims whitespace and trailing slash",
			issuer:   "  https://issuer.example.com/  ",
			expected: "https://issuer.example.com",
		},
		{
			name:     "normalizes scheme and host but preserves path case",
			issuer:   "HTTPS://ISSUER.EXAMPLE.COM/REALM/Team/",
			expected: "https://issuer.example.com/REALM/Team",
		},
		{
			name:     "trims multiple trailing slashes",
			issuer:   "https://issuer.example.com///",
			expected: "https://issuer.example.com",
		},
		{
			name:     "trims trailing slashes from path",
			issuer:   "https://example.com/path///",
			expected: "https://example.com/path",
		},
		{
			name:     "whitespace only becomes empty",
			issuer:   " \t\n ",
			expected: "",
		},
		{
			name:     "empty string remains empty",
			issuer:   "",
			expected: "",
		},

		// Port normalization
		{
			name:     "removes default https port",
			issuer:   "https://example.com:443",
			expected: "https://example.com",
		},
		{
			name:     "removes default https port with path",
			issuer:   "https://example.com:443/path",
			expected: "https://example.com/path",
		},
		{
			name:     "removes default http port",
			issuer:   "http://example.com:80",
			expected: "http://example.com",
		},
		{
			name:     "preserves non-default port",
			issuer:   "https://example.com:8443",
			expected: "https://example.com:8443",
		},
		{
			name:     "preserves non-default port with path",
			issuer:   "https://example.com:8443/path",
			expected: "https://example.com:8443/path",
		},

		// Percent-encoding normalization
		{
			name:     "uppercases percent-encoded hex digits",
			issuer:   "https://example.com/%2a",
			expected: "https://example.com/%2A",
		},
		{
			name:     "decodes unreserved characters - tilde",
			issuer:   "https://example.com/%7Euser",
			expected: "https://example.com/~user",
		},
		{
			name:     "decodes unreserved characters - hyphen period underscore",
			issuer:   "https://example.com/%2D%2E%5F",
			expected: "https://example.com/-._",
		},
		{
			name:     "preserves reserved percent-encoded characters",
			issuer:   "https://example.com/path%2Fto",
			expected: "https://example.com/path%2Fto",
		},
		{
			name:     "alpha/numeric",
			issuer:   "https://example.com/%41%31",
			expected: "https://example.com/A1",
		},
		{
			name:     "hyphen/period",
			issuer:   "https://example.com/%2D%2E",
			expected: "https://example.com/-.",
		},
		{
			name:     "underscore/tilde",
			issuer:   "https://example.com/%5F%7E",
			expected: "https://example.com/_~",
		},

		// Complex cases
		{
			name:     "combined normalization with all features",
			issuer:   "  HTTPS://EXAMPLE.COM:443/REALM/Team/  ",
			expected: "https://example.com/REALM/Team",
		},
		{
			name:     "complex path with mixed case and trailing slashes",
			issuer:   "HTTPS://ISSUER.EXAMPLE.COM:443/realm/TEAM///",
			expected: "https://issuer.example.com/realm/TEAM",
		},
		{
			name:     "path with percent-encoding and case preservation",
			issuer:   "https://example.com/Path%2FWith%7ESpecial",
			expected: "https://example.com/Path%2FWith~Special",
		},

		// Root path
		{
			name:     "root path with trailing slash",
			issuer:   "https://example.com/",
			expected: "https://example.com",
		},
		{
			name:     "root path uppercase with trailing slash",
			issuer:   "HTTPS://EXAMPLE.COM/",
			expected: "https://example.com",
		},

		// Path segment normalization
		{
			name:     "single dot",
			issuer:   "http://example.com/a/./b",
			expected: "http://example.com/a/b",
		},
		{
			name:     "double dot",
			issuer:   "http://example.com/a/b/../c",
			expected: "http://example.com/a/c",
		},
		{
			name:     "leading dots",
			issuer:   "http://example.com/../a",
			expected: "http://example.com/a",
		},
		{
			name:     "trailing dot",
			issuer:   "http://example.com/a/.",
			expected: "http://example.com/a",
		},
		{
			name:     "complex path",
			issuer:   "http://example.com/a/b/c/./../../g",
			expected: "http://example.com/a/g",
		},

		// Edge cases
		{
			name:     "invalid URL returns trimmed input",
			issuer:   "not-a-url",
			expected: "not-a-url",
		},
		{
			name:     "URL without scheme",
			issuer:   "example.com/path",
			expected: "example.com/path",
		},

		// Idempotency
		{
			name:     "already normalized remains unchanged",
			issuer:   "https://example.com/path",
			expected: "https://example.com/path",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NormalizeIssuer(tc.issuer)
			require.Equal(t, tc.expected, result)

			// Verify idempotency: normalizing twice should produce same result
			result2 := NormalizeIssuer(result)
			require.Equal(t, result, result2, "normalizeIssuer should be idempotent")
		})
	}
}
