// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package pki

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSecureImportKey_CEError verifies that supplying wrapped_key on CE returns an enterprise error.
func TestSecureImportKey_CEError(t *testing.T) {
	t.Parallel()

	b := &backend{}
	_, err := b.unwrapAndImportKey(context.Background(), nil, "wrapped-blob", "sha256:1234")
	require.Error(t, err)
	require.Contains(t, err.Error(), "secure key import with wrapped_key is only supported on Vault Enterprise")
}
