// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSignIntermediate_MaxPathLengthValidation verifies that when signing an
// intermediate CA, the requested max_path_length is validated against the
// signing CA's BasicConstraints pathLenConstraint per RFC 5280 4.2.1.9.
//
// If the signing CA has a pathLenConstraint of N, any intermediate it signs
// must have a pathLenConstraint strictly less than N.
//
// When max_path_length is not provided at all (omitted from the request), the
// validation block is skipped entirely and the request always succeeds.
//
// When max_path_length is explicitly set to -1 (meaning "no constraint on the
// intermediate"), and the signing CA already has a pathLenConstraint, the
// request is rejected because an unconstrained intermediate would violate the
// CA's constraint.
//
// When the generated certificate has MaxPathLen == 0 (pathLenConstraint=0),
// Vault adds a warning that the certificate cannot be used to issue
// intermediate CAs.
func TestSignIntermediate_MaxPathLengthValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		caMaxPathLength  int  // max_path_length set on the root CA (-1 means not set / unconstrained)
		intMaxPathLength int  // max_path_length requested when signing the intermediate (-1 = "no constraint")
		omitParam        bool // if true, do not include max_path_length in the sign request at all
		expectError      bool
		errorContains    string
		// Expected values on the issued certificate (only checked when expectError==false).
		expectedMaxPathLen int // expected cert.MaxPathLen; -1 means "no constraint"
	}{
		// --- Explicit positive values: validation is active ---
		{
			name:               "ca_path_len_2_int_path_len_1_allowed",
			caMaxPathLength:    2,
			intMaxPathLength:   1,
			expectError:        false,
			expectedMaxPathLen: 1,
		},
		{
			name:             "ca_path_len_2_int_path_len_2_rejected",
			caMaxPathLength:  2,
			intMaxPathLength: 2,
			expectError:      true,
			errorContains:    "requested max_path_length 2 is not allowed",
		},
		{
			name:             "ca_path_len_2_int_path_len_3_rejected",
			caMaxPathLength:  2,
			intMaxPathLength: 3,
			expectError:      true,
			errorContains:    "requested max_path_length 3 is not allowed",
		},
		{
			// pathLenConstraint=0 on the issued cert; Go represents this with
			// MaxPathLen==0 and MaxPathLenZero==true.
			name:               "ca_path_len_1_int_path_len_0_allowed",
			caMaxPathLength:    1,
			intMaxPathLength:   0,
			expectError:        false,
			expectedMaxPathLen: 0,
		},
		{
			name:             "ca_path_len_1_int_path_len_1_rejected",
			caMaxPathLength:  1,
			intMaxPathLength: 1,
			expectError:      true,
			errorContains:    "requested max_path_length 1 is not allowed",
		},
		{
			// pathLenConstraint=0 on the CA means it may not issue further CAs
			// with any pathLenConstraint (including 0).
			name:             "ca_path_len_0_int_path_len_0_rejected",
			caMaxPathLength:  0,
			intMaxPathLength: 0,
			expectError:      true,
			errorContains:    "requested max_path_length 0 is not allowed",
		},
		{
			// CA has no pathLenConstraint; any explicit positive value is allowed.
			name:               "ca_no_path_len_constraint_int_path_len_5_allowed",
			caMaxPathLength:    -1, // no constraint on CA
			intMaxPathLength:   5,
			expectError:        false,
			expectedMaxPathLen: 5,
		},

		// --- Explicit -1 ("no constraint on intermediate"): rejected only when CA is constrained ---
		{
			// CA has no constraint; explicit -1 is allowed and results in no
			// pathLenConstraint on the issued certificate.
			name:               "ca_no_constraint_int_explicit_neg1_allowed",
			caMaxPathLength:    -1,
			intMaxPathLength:   -1,
			expectError:        false,
			expectedMaxPathLen: -1,
		},
		{
			// CA has a constraint; explicit -1 means "no constraint on the
			// intermediate", which violates the CA's pathLenConstraint.
			name:             "ca_path_len_5_int_explicit_neg1_rejected",
			caMaxPathLength:  5,
			intMaxPathLength: -1,
			expectError:      true,
			errorContains:    "requested max_path_length -1 is not allowed",
		},
		{
			// CA has pathLenConstraint=0; explicit -1 is also rejected.
			name:             "ca_path_len_0_int_explicit_neg1_rejected",
			caMaxPathLength:  0,
			intMaxPathLength: -1,
			expectError:      true,
			errorContains:    "requested max_path_length -1 is not allowed",
		},

		// --- Parameter omitted entirely: validation is skipped ---
		// When omitted, Vault auto-derives the pathLenConstraint as (CA pathLen - 1).
		// When the CA has no constraint, the issued cert also has no constraint (-1).
		{
			// CA has no constraint; omitting the parameter results in no
			// pathLenConstraint on the issued certificate.
			name:               "ca_no_constraint_int_param_omitted_allowed",
			caMaxPathLength:    -1,
			omitParam:          true,
			expectError:        false,
			expectedMaxPathLen: -1,
		},
		{
			// CA has pathLenConstraint=5; omitting the parameter causes Vault to
			// auto-derive pathLenConstraint = 5 - 1 = 4 on the issued certificate.
			name:               "ca_path_len_5_int_param_omitted_allowed",
			caMaxPathLength:    5,
			omitParam:          true,
			expectError:        false,
			expectedMaxPathLen: 4,
		},
		{
			// CA has pathLenConstraint=1; omitting the parameter causes Vault to
			// auto-derive pathLenConstraint = 1 - 1 = 0 on the issued certificate.
			name:               "ca_path_len_1_int_param_omitted_allowed",
			caMaxPathLength:    1,
			omitParam:          true,
			expectError:        false,
			expectedMaxPathLen: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			b, s := CreateBackendWithStorage(t)

			// Generate root CA with the specified max_path_length.
			rootParams := map[string]interface{}{
				"common_name": "root.example.com",
				"ttl":         "87600h",
			}
			if tc.caMaxPathLength >= 0 {
				rootParams["max_path_length"] = tc.caMaxPathLength
			}

			resp, err := CBWrite(b, s, "root/generate/internal", rootParams)
			requireSuccessNonNilResponse(t, resp, err, "failed to generate root CA")

			// Generate an intermediate CSR.
			resp, err = CBWrite(b, s, "intermediate/generate/internal", map[string]interface{}{
				"common_name": "int.example.com",
			})
			requireSuccessNonNilResponse(t, resp, err, "failed to generate intermediate CSR")
			csr := resp.Data["csr"].(string)

			// Build the sign-intermediate request.
			signParams := map[string]interface{}{
				"common_name": "int.example.com",
				"csr":         csr,
				"ttl":         "43800h",
			}
			if !tc.omitParam {
				signParams["max_path_length"] = tc.intMaxPathLength
			}

			resp, err = CBWrite(b, s, "root/sign-intermediate", signParams)

			if tc.expectError {
				require.Error(t, err, "expected sign-intermediate to fail but it succeeded")
				if tc.errorContains != "" {
					require.Contains(t, err.Error(), tc.errorContains,
						"error message did not contain expected text")
				}
			} else {
				require.NoError(t, err, "expected sign-intermediate to succeed but it failed")
				require.NotNil(t, resp, "expected non-nil response from sign-intermediate")
				require.False(t, resp.IsError(), "sign-intermediate returned error: %v", resp.Error())

				cert := parseCert(t, resp.Data["certificate"].(string))
				require.True(t, cert.BasicConstraintsValid, "expected BasicConstraints to be set")
				require.True(t, cert.IsCA, "expected IsCA to be true")

				require.Equal(t, tc.expectedMaxPathLen == 0, cert.MaxPathLenZero,
					"issued certificate has unexpected MaxPathLenZero")
				require.Equal(t, tc.expectedMaxPathLen, cert.MaxPathLen,
					"issued certificate has unexpected MaxPathLen")
				if tc.expectedMaxPathLen == 0 {
					requireMaxPathLengthZeroWarning(t, resp.Warnings)
				}
			}
		})
	}
}
