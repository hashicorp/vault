// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package aead

import (
	"context"
	"crypto/rand"
	"testing"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/hashicorp/go-kms-wrapping/v2/extras/multi"
	"github.com/stretchr/testify/require"
)

// TestPooledWrapper returns a pooled aead wrapper for testing
func TestPooledWrapper(t *testing.T) wrapping.Wrapper {
	t.Helper()
	require := require.New(t)
	testCtx := context.Background()
	w1Key := make([]byte, 32)
	n, err := rand.Read(w1Key)
	require.NoError(err)
	require.Equal(32, n)
	w1 := NewWrapper()
	_, err = w1.SetConfig(testCtx, wrapping.WithKeyId("w1"))
	require.NoError(err)
	err = w1.SetAesGcmKeyBytes(w1Key)
	require.NoError(err)

	w2Key := make([]byte, 32)
	n, err = rand.Read(w2Key)
	require.NoError(err)
	require.Equal(32, n)

	w2 := NewWrapper()
	_, err = w2.SetConfig(testCtx, wrapping.WithKeyId("w2"))
	require.NoError(err)
	err = w2.SetAesGcmKeyBytes(w2Key)
	require.NoError(err)

	p, err := multi.NewPooledWrapper(testCtx, w1)
	require.NoError(err)
	return p
}

// TestWrapper returns a test aead wrapper for testing
func TestWrapper(t *testing.T) wrapping.Wrapper {
	t.Helper()
	require := require.New(t)
	testCtx := context.Background()
	w1Key := make([]byte, 32)
	n, err := rand.Read(w1Key)
	require.NoError(err)
	require.Equal(32, n)
	w := NewWrapper()
	_, err = w.SetConfig(testCtx, wrapping.WithKeyId("w1"))
	require.NoError(err)
	err = w.SetAesGcmKeyBytes(w1Key)
	require.NoError(err)
	return w
}
