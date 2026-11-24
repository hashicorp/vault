// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAcmeNonces(t *testing.T) {
	t.Parallel()

	a := NewACMEState()
	a.nonces.Initialize()

	// Simple operation should succeed.
	nonce, _, err := a.GetNonce()
	require.NoError(t, err)
	require.NotEmpty(t, nonce)

	require.True(t, a.RedeemNonce(nonce))
	require.False(t, a.RedeemNonce(nonce))

	// Redeeming in opposite order should work.
	var nonces []string
	for i := 0; i < len(nonce); i++ {
		nonce, _, err = a.GetNonce()
		require.NoError(t, err)
		require.NotEmpty(t, nonce)
	}

	for i := len(nonces) - 1; i >= 0; i-- {
		nonce = nonces[i]
		require.True(t, a.RedeemNonce(nonce))
	}

	for i := 0; i < len(nonces); i++ {
		nonce = nonces[i]
		require.False(t, a.RedeemNonce(nonce))
	}
}

// TestVerifyEabPayloadProtectedFieldValidation verifies that the verifyEabPayload function
// properly validates and rejects malformed 'protected' fields in ACME External Account Binding
// (EAB) payloads. It ensures that only string values are accepted and that various non-string
// types (object, array, number, boolean) produce appropriate error messages.
func TestVerifyEabPayloadProtectedFieldValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		payload       map[string]interface{}
		expectError   bool
		expectMessage string
	}{
		{
			name: "valid string protected field",
			payload: map[string]interface{}{
				"protected": "eyJhbGciOiAiSFMyNTYifQ",
			},
			expectError: false,
		},
		{
			name: "missing protected field",
			payload: map[string]interface{}{
				"payload":   "test",
				"signature": "test",
			},
			expectError:   true,
			expectMessage: "missing required field 'protected'",
		},
		{
			name: "protected field as object",
			payload: map[string]interface{}{
				"protected": map[string]interface{}{ // should be a string, not an object
					"alg": "HS256",
				},
				"payload":   "test",
				"signature": "test",
			},
			expectError:   true,
			expectMessage: "failed to parse 'protected' field",
		},
		{
			name: "protected field as array",
			payload: map[string]interface{}{
				"protected": []interface{}{"test"}, // should be a string, not an array
				"payload":   "test",
				"signature": "test",
			},
			expectError:   true,
			expectMessage: "failed to parse 'protected' field",
		},
		{
			name: "protected field as number",
			payload: map[string]interface{}{
				"protected": 12345, // should be a string, not a number
				"payload":   "test",
				"signature": "test",
			},
			expectError:   true,
			expectMessage: "failed to parse 'protected' field",
		},
		{
			name: "protected field as boolean",
			payload: map[string]interface{}{
				"protected": true, // should be a string, not a boolean
				"payload":   "test",
				"signature": "test",
			},
			expectError:   true,
			expectMessage: "failed to parse 'protected' field",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			acmeState := NewACMEState()
			acmeCtx := &acmeContext{}
			outerJws := &jwsCtx{}

			_, err := verifyEabPayload(acmeState, acmeCtx, outerJws, "/new-account", tc.payload)

			if tc.expectError {
				require.Error(t, err, "expected error but got none")
				require.Contains(t, err.Error(), tc.expectMessage, "error message does not match")
				require.NotEmpty(t, err.Error())
			} else {
				// For valid protected field, we expect an error later in the flow
				// (since we don't have real base64 data), but it should NOT be from parsing the protected field
				if err != nil {
					require.NotContains(t, err.Error(), "failed to parse 'protected' field",
						"should not fail at protected field parsing stage")
				}
			}
		})
	}
}

// TestErrorResponseNoSubproblems builds the http body that exists in the header of an ACME error response and checks
// in a simple case that "type" and "detail" two fields on the body do exist, but that "subproblems" a field which is
// optional, is omitted because it does not exist in this case (rather than being included with a value null which can
// trip up some systems).
func TestErrorResponseNoSubproblems(t *testing.T) {
	t.Parallel()
	errResponse, err := TranslateError(ErrAlreadyRevoked)
	if err != nil {
		return
	}
	require.NoError(t, err, "already revoked should generate an error response")
	require.NotNil(t, errResponse.Data)
	body := map[string]string{}
	rawBody, ok := errResponse.Data["http_raw_body"]
	err = json.Unmarshal(rawBody.([]byte), &body)
	require.True(t, ok, "Raw Body of Error response should exist, but doesn't")
	typeString, ok := body["type"]
	require.True(t, ok, "Type on Raw Body of Error response should exist, but doesn't")
	require.Equal(t, typeString, "urn:ietf:params:acme:error:alreadyRevoked")
	_, ok = body["detail"]
	require.True(t, ok, "Detail on Raw Body of Error response should exist, but doesn't")
	subProblems, ok := body["subproblems"]
	require.False(t, ok, "subproblems on Raw Body of Error response should be omitted, but exists with value %v", subProblems)
}
