// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package vault_test

import (
	"testing"

	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// TestIsEphemeralEntityID verifies that ephemeral entity IDs are recognized by
// their reserved prefix alone, without depending on a temporary MemDB record.
func TestIsEphemeralEntityID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   string
		want bool
	}{
		{name: "with identifier", id: "eid.abc", want: true},
		{name: "prefix only", id: "eid.", want: true},
		{name: "missing separator", id: "eid", want: false},
		{name: "uppercase prefix", id: "EID.abc", want: false},
		{name: "empty", id: "", want: false},
		{name: "UUID", id: "550e8400-e29b-41d4-a716-446655440000", want: false},
		{name: "leading space", id: " eid.abc", want: false},
		{name: "prefix not at start", id: "xeid.abc", want: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.want {
				require.True(t, vault.IsEphemeralEntityID(tt.id), "ID %q should be ephemeral", tt.id)
				return
			}
			require.False(t, vault.IsEphemeralEntityID(tt.id), "ID %q should not be ephemeral", tt.id)
		})
	}
}
