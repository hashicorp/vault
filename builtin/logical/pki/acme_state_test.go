// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
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
