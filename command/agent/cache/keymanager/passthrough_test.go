// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keymanager

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyManager_PassthrougKeyManager(t *testing.T) {
	tests := []struct {
		name    string
		key     []byte
		wantErr bool
	}{
		{
			"new key",
			nil,
			false,
		},
		{
			"existing valid key",
			[]byte("e679e2f3d8d0e489d408bc617c6890d6"),
			false,
		},
		{
			"invalid key length",
			[]byte("foobar"),
			true,
		},
	}

	ctx := context.Background()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := NewPassthroughKeyManager(ctx, tc.key)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if w := m.Wrapper(); w == nil {
				t.Fatalf("expected non-nil wrapper from the key manager")
			}

			token, err := m.RetrievalToken(ctx)
			if err != nil {
				t.Fatalf("unable to retrieve token: %s", err)
			}

			if len(tc.key) != 0 && !bytes.Equal(tc.key, token) {
				t.Fatalf("expected key bytes: %x, got: %x", tc.key, token)
			}
		})
	}
}
