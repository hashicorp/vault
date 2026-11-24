// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package tpmutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestIDFromDER verifies that a TPM ID is "sha256-" followed by the lowercase
// hex SHA-256 of the DER bytes.  The identity store enrolls TPMs under this ID
// and `vault tpm ek` prints it, so the format must not change.
func TestIDFromDER(t *testing.T) {
	t.Parallel()

	require.Equal(t, "sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", IDFromDER(nil))
	require.Equal(t, "sha256-ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad", IDFromDER([]byte("abc")))
}
