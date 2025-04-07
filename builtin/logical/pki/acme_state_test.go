// Copyright (c) HashiCorp, Inc.
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
