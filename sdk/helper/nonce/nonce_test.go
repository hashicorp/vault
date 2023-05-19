package nonce

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNonceService(t *testing.T) {
	t.Parallel()

	s := NewNonceService()
	err := s.Initialize()
	require.NoError(t, err)

	nonce, _, err := s.Get()
	require.NoError(t, err)
	require.NotEmpty(t, nonce)

	require.True(t, s.Redeem(nonce))
	require.False(t, s.Redeem(nonce))

	// Redeeming in opposite order should work.
	var nonces []string
	numNonces := 100
	for i := 0; i < numNonces; i++ {
		nonce, _, err = s.Get()
		require.NoError(t, err)
		require.NotEmpty(t, nonce)

		nonces = append(nonces, nonce)
	}

	for i := len(nonces) - 1; i >= 0; i-- {
		nonce = nonces[i]
		require.True(t, s.Redeem(nonce))
	}

	for i := 0; i < len(nonces); i++ {
		nonce = nonces[i]
		require.False(t, s.Redeem(nonce))
	}

    status := s.Tidy()
    require.NotNil(t, status)
    require.Equal(t, uint64(1 + numNonces), status.Issued)
    require.Equal(t, uint64(0), status.Outstanding)
}
