// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"software.sslmate.com/src/go-pkcs12"
)

// TestGenerateRoot_InvalidCountry validates that ISO 3166 is followed for the Country field
func TestGenerateRoot_InvalidCountry(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)
	params := map[string]interface{}{
		"common_name": "root.example.com",
		"ttl":         "87600h",
		"key_type":    "ec",
		"country":     "Japan",
	}

	resp, err := CBWrite(b, s, "root/generate/internal", params)
	require.NoError(t, err)

	require.True(t, stringSliceContainsAny(resp.Warnings, "3166"))
}

func TestGenerateRoot_MaxPathLengthValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		maxPathLength int
		omitParam     bool

		expectError   bool
		errorContains string

		// Expected fields on the issued certificate (only checked when
		// expectError == false).
		expectedMaxPathLen int // cert.MaxPathLen; -1 means no constraint
	}{
		// ----------------------------------------------------------------
		// Parameter omitted: Vault generates a root with no pathLenConstraint.
		// ----------------------------------------------------------------
		{
			name:               "param_omitted_no_constraint",
			omitParam:          true,
			expectError:        false,
			expectedMaxPathLen: -1,
		},

		// ----------------------------------------------------------------
		// Explicit -1: "no constraint" — identical outcome to omitting the
		// parameter; the generated certificate carries no pathLenConstraint.
		// ----------------------------------------------------------------
		{
			name:               "explicit_neg1_no_constraint",
			maxPathLength:      -1,
			expectError:        false,
			expectedMaxPathLen: -1,
		},

		// ----------------------------------------------------------------
		// Explicit 0: pathLenConstraint=0 — the CA may not issue further
		// intermediate CAs.  Vault should succeed but add a warning.
		// ----------------------------------------------------------------
		{
			name:               "explicit_0_zero_constraint_with_warning",
			maxPathLength:      0,
			expectError:        false,
			expectedMaxPathLen: 0,
		},

		// ----------------------------------------------------------------
		// Explicit positive values: pathLenConstraint is set as requested.
		// ----------------------------------------------------------------
		{
			name:               "explicit_1",
			maxPathLength:      1,
			expectError:        false,
			expectedMaxPathLen: 1,
		},
		{
			name:               "explicit_2",
			maxPathLength:      2,
			expectError:        false,
			expectedMaxPathLen: 2,
		},
		{
			name:               "explicit_5",
			maxPathLength:      5,
			expectError:        false,
			expectedMaxPathLen: 5,
		},

		// ----------------------------------------------------------------
		// Invalid values (< -1): the handler must reject these immediately.
		// ----------------------------------------------------------------
		{
			name:          "invalid_neg2",
			maxPathLength: -2,
			expectError:   true,
			errorContains: "max_path_length -2 is invalid",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			b, s := CreateBackendWithStorage(t)

			params := map[string]interface{}{
				"common_name": "root.example.com",
				"ttl":         "87600h",
				"key_type":    "ec",
			}
			if !tc.omitParam {
				params["max_path_length"] = tc.maxPathLength
			}

			resp, err := CBWrite(b, s, "root/generate/internal", params)

			if tc.expectError {
				require.Error(t, err, "expected root/generate/internal to fail but it succeeded")
				if tc.errorContains != "" {
					require.Contains(t, err.Error(), tc.errorContains,
						"error message did not contain expected text")
				}
				return
			}

			// Success path.
			require.NoError(t, err, "expected root/generate/internal to succeed but it failed")
			require.NotNil(t, resp, "expected non-nil response from root/generate/internal")
			require.False(t, resp.IsError(), "root/generate/internal returned error response: %v", resp.Error())

			// Parse the returned PEM certificate and verify BasicConstraints.
			certPEM, ok := resp.Data["certificate"].(string)
			require.True(t, ok, "response missing 'certificate' field")

			cert := parseCert(t, certPEM)
			require.True(t, cert.BasicConstraintsValid, "expected BasicConstraints to be set on root CA")
			require.True(t, cert.IsCA, "expected IsCA to be true on root CA")

			require.Equal(t, tc.expectedMaxPathLen, cert.MaxPathLen,
				"certificate has unexpected MaxPathLen")
			require.Equal(t, tc.expectedMaxPathLen == 0, cert.MaxPathLenZero,
				"certificate has unexpected MaxPathLenZero")

			// Check for the zero-path-length warning.
			if tc.expectedMaxPathLen == 0 {
				requireMaxPathLengthZeroWarning(t, resp.Warnings)
			}
		})
	}
}

func requireMaxPathLengthZeroWarning(t testing.TB, warnings []string) {
	t.Helper()
	require.NotEmpty(t, warnings, "expected at least one warning for zero max_path_length")
	foundWarning := false
	for _, w := range warnings {
		if strings.Contains(w, "Max path length of the generated certificate is zero") {
			foundWarning = true
			break
		}
	}
	require.True(t, foundWarning, "expected warning about zero max path length, got warnings: %v", warnings)
}

// TestGenerateAndRotateRoot_PKCS12Format validates PKCS12 support for (internal and exported) root generation and rotation.
// PKCS12 archives from exported endpoints should contain a private key and root certificate while
// internal endpoints should only have the root certificate.
func TestGenerateAndRotateRoot_PKCS12Format(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)
	var password string
	buildData := func(omitPassword bool, encoder string) map[string]interface{} {
		password = pkcs12.DefaultPassword
		data := map[string]interface{}{
			"format":      "pkcs12_bundle",
			"common_name": "Root CA",
			"ttl":         "87600h",
			"key_type":    "ec",
			"key_bits":    256,
		}
		if !omitPassword {
			password = "secure-root-password"
			data["pkcs12_password"] = password
		}
		if encoder != "" {
			data["pkcs12_encoder"] = encoder
		}
		return data
	}

	for _, p := range []string{"root/generate/", "root/rotate/"} {
		testCases := []struct {
			name         string
			endpoint     string
			encoder      string
			omitPassword bool
			shouldError  bool
		}{
			// exported
			{name: "custom password", endpoint: "exported"},
			{name: "default password", endpoint: "exported", omitPassword: true},
			{name: "with modern2026 encoder", endpoint: "exported", encoder: "modern2026"},
			{name: "with modern2023 encoder", endpoint: "exported", encoder: "modern2023"},
			{name: "with invalid encoder", endpoint: "exported", encoder: "modern2020", shouldError: true},
			// internal
			{name: "custom password", endpoint: "internal"},
			{name: "default password", endpoint: "internal", omitPassword: true},
			{name: "with modern2026 encoder", endpoint: "internal", encoder: "modern2026"},
			{name: "with modern2023 encoder", endpoint: "internal", encoder: "modern2023"},
			{name: "with invalid encoder", endpoint: "internal", encoder: "modern2020", shouldError: true},
		}

		for _, tc := range testCases {
			path := p + tc.endpoint
			t.Run(tc.name, func(t *testing.T) {
				data := buildData(tc.omitPassword, tc.encoder)
				pkcs12Bytes := requestAndVerifyPKCS12(t, tc.shouldError, requestPKCS12params{b, s, path, data})
				if tc.shouldError {
					return
				}

				if tc.endpoint == "exported" {
					_, cert, caCerts := requireDecodesPKCS12Chain(t, pkcs12Bytes, password)
					require.Equal(t, "Root CA", cert.Subject.CommonName)
					require.True(t, cert.IsCA, "should be a CA certificate")
					// The root certificate itself is in 'cert', not in 'caCerts'
					require.Len(t, caCerts, 0, "should have no CA chain because root is self-signed")
				} else {
					// Validate PKCS12 trust store for internal root: contains only certificates (no private key)
					certs := requireDecodesPKCS12TrustStore(t, pkcs12Bytes, password)
					require.Len(t, certs, 1, "should have no additional certs because root is self-signed")
					// First cert should be the root
					require.True(t, certs[0].IsCA, "cert should be a CA")
					require.Equal(t, "Root CA", certs[0].Subject.CommonName)
					requireSignedBy(t, certs[0], certs[0])
				}
			})
		}
	}
}

// TestGenerateAndRotateRoot_FormatValidation validates that empty format string and invalid encoders are rejected
func TestGenerateAndRotateRoot_FormatValidation(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	tcEmptyFormat := []struct {
		name     string
		endpoint string
		format   string
	}{
		{name: "generate empty format", endpoint: "root/generate/exported", format: ""},
		{name: "rotate empty format", endpoint: "root/rotate/exported", format: ""},
	}

	for _, tc := range tcEmptyFormat {
		t.Run(tc.endpoint, func(t *testing.T) {
			_, err := CBWrite(b, s, tc.endpoint, map[string]interface{}{
				"common_name": "Root CA",
				"format":      "",
				"key_type":    "ec",
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), `the "format" parameter must be "pem", "der", "pem_bundle" or "pkcs12_bundle"`)
		})
	}

	tcInvalidEncoder := []struct {
		name     string
		endpoint string
		encoder  string
	}{
		{name: "generate invalid format", endpoint: "root/generate/exported", encoder: "invalid"},
		{name: "rotate invalid format", endpoint: "root/rotate/exported", encoder: "invalid"},
	}

	for _, tc := range tcInvalidEncoder {
		t.Run(tc.endpoint, func(t *testing.T) {
			_, err := CBWrite(b, s, tc.endpoint, map[string]interface{}{
				"common_name":    "Root CA",
				"format":         "pkcs12_bundle",
				"pkcs12_encoder": tc.encoder,
				"key_type":       "ec",
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), `invalid "pkcs12_encoder" parameter: encoder must be "modern2026" or "modern2023"; received: "invalid"`)
		})
	}
}
