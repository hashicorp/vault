// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
