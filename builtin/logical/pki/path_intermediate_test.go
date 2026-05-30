// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGenerateIntermediate_FormatParam validates that CSR generation
// only supports valid formats and does NOT support PKCS#12 since there's no certificate to bundle
func TestGenerateIntermediate_FormatParam(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)
	testCases := []struct {
		endpoint   string
		format     string
		isValid    bool
		omitFormat bool
	}{
		{endpoint: "intermediate/generate/internal", format: "pkcs12_bundle"},
		{endpoint: "intermediate/generate/exported", format: "pkcs12_bundle"},
		{endpoint: "issuers/generate/intermediate/internal", format: "pkcs12_bundle"},
		{endpoint: "issuers/generate/intermediate/exported", format: "pkcs12_bundle"},

		{endpoint: "intermediate/generate/internal", format: "invalid"},
		{endpoint: "intermediate/generate/exported", format: "invalid"},
		{endpoint: "issuers/generate/intermediate/internal", format: "invalid"},
		{endpoint: "issuers/generate/intermediate/exported", format: "invalid"},

		{endpoint: "intermediate/generate/internal", format: ""},
		{endpoint: "intermediate/generate/exported", format: ""},
		{endpoint: "issuers/generate/intermediate/internal", format: ""},
		{endpoint: "issuers/generate/intermediate/exported", format: ""},

		{endpoint: "intermediate/generate/internal", omitFormat: true},
		{endpoint: "intermediate/generate/exported", omitFormat: true},
		{endpoint: "issuers/generate/intermediate/internal", omitFormat: true},
		{endpoint: "issuers/generate/intermediate/exported", omitFormat: true},
	}
	for _, tc := range testCases {
		name := fmt.Sprintf("endpoint=%s format=%s", tc.endpoint, tc.format)
		t.Run(name, func(t *testing.T) {
			// Attempt to generate intermediate CSR with various formats
			data := map[string]interface{}{
				"common_name": "Intermediate CA",
				"key_type":    "ec",
				"key_bits":    256,
			}
			if !tc.omitFormat {
				data["format"] = tc.format
			}
			_, err := CBWrite(b, s, tc.endpoint, data)

			if tc.omitFormat {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), `the "format" parameter must be "pem", "der" or "pem_bundle"`)
			}
		})
	}
}
