// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package jwt

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// normalizeIssuer normalizes an issuer URL according to RFC 9207 and RFC 3986.
//
// RFC 9207 Section 2.4 requires that issuer identifiers be compared using
// "simple string comparison" as defined in RFC 3986 Section 6.2.1, after
// proper normalization.
//
// RFC 3986 Section 6.2 specifies that:
// - Scheme and host components must be normalized to lowercase
// - Path components are case-sensitive and must NOT be lowercased
// - Default ports should be removed (443 for https, 80 for http)
// - Trailing slashes should be removed
// - Percent-encoding should be normalized
//
// This function performs the following normalizations:
// 1. Trims leading and trailing whitespace
// 2. Parses the URL to separate components
// 3. Lowercases the scheme (e.g., HTTPS → https)
// 4. Lowercases the host (e.g., Example.COM → example.com)
// 5. Preserves path case-sensitivity (e.g., /REALM/Team remains /REALM/Team)
// 6. Removes default ports (:443 for https, :80 for http)
// 7. Removes trailing slashes from the path
// 8. Normalizes percent-encoding in the path (uppercase hex, decode unreserved)
//
// Examples:
//
//	normalizeIssuer("HTTPS://EXAMPLE.COM/REALM/Team/")
//	  → "https://example.com/REALM/Team"
//	normalizeIssuer("https://example.com:443/path")
//	  → "https://example.com/path"
//	normalizeIssuer("  https://Example.COM/  ")
//	  → "https://example.com"
//
// The function is idempotent: normalizing an already-normalized issuer
// returns the same value.
func NormalizeIssuer(issuer string) string {
	// Trim whitespace
	issuer = strings.TrimSpace(issuer)
	if issuer == "" {
		return ""
	}

	// Parse the URL
	u, err := url.Parse(issuer)
	if err != nil {
		// If parsing fails, return the trimmed string as-is
		return issuer
	}

	// Normalize scheme to lowercase
	u.Scheme = strings.ToLower(u.Scheme)

	// Normalize host to lowercase
	u.Host = strings.ToLower(u.Host)

	// Remove default ports
	if (u.Scheme == "https" && strings.HasSuffix(u.Host, ":443")) ||
		(u.Scheme == "http" && strings.HasSuffix(u.Host, ":80")) {
		// Remove the default port
		host := u.Host
		if idx := strings.LastIndex(host, ":"); idx != -1 {
			u.Host = host[:idx]
		}
	}

	// Normalize path: resolve dot segments, preserve case, remove trailing slashes, and normalize percent-encoding
	if u.Path != "" {
		// First resolve dot segments (., .., etc.) per RFC 3986 Section 5.2.4
		u.Path = path.Clean(u.Path)

		// Remove trailing slashes from the decoded path
		u.Path = strings.TrimRight(u.Path, "/")

		// Get the escaped path (percent-encoded) to normalize it
		escapedPath := u.EscapedPath()

		// Normalize percent-encoding in the escaped path
		normalizedPath := normalizePercentEncoding(escapedPath)

		// Set both Path and RawPath to ensure proper encoding
		u.RawPath = normalizedPath
		// Decode the normalized path back to set u.Path correctly
		if decodedPath, err := url.PathUnescape(normalizedPath); err == nil {
			u.Path = decodedPath
		}
	}

	// Reconstruct the URL
	return u.String()
}

// normalizePercentEncoding normalizes percent-encoded characters in a path according to RFC 3986.
// It uppercases hexadecimal digits in percent-encoded triplets and decodes unreserved characters.
func normalizePercentEncoding(path string) string {
	// isUnreserved returns true if the byte is an unreserved character per RFC 3986.
	// Unreserved characters: A-Z, a-z, 0-9, hyphen (-), period (.), underscore (_), tilde (~)
	isUnreserved := func(b byte) bool {
		return (b >= 'A' && b <= 'Z') ||
			(b >= 'a' && b <= 'z') ||
			(b >= '0' && b <= '9') ||
			b == '-' || b == '.' || b == '_' || b == '~'
	}

	var result strings.Builder
	result.Grow(len(path))

	for i := 0; i < len(path); i++ {
		if path[i] == '%' && i+2 < len(path) {
			// Extract the two hex digits
			hex := path[i+1 : i+3]

			// Try to parse the hex value (case-insensitive)
			var b byte
			if n, err := fmt.Sscanf(strings.ToLower(hex), "%02x", &b); err == nil && n == 1 {
				// Check if this is an unreserved character that should be decoded
				if isUnreserved(b) {
					result.WriteByte(b)
					i += 2
					continue
				}

				// Otherwise, normalize to uppercase hex
				result.WriteByte('%')
				result.WriteString(strings.ToUpper(hex))
				i += 2
				continue
			}
		}
		result.WriteByte(path[i])
	}

	return result.String()
}
