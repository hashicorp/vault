// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"path"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

func (s *Session) MustWrite(path string, data map[string]any) *api.Secret {
	s.t.Helper()

	secret, err := s.Client.Logical().Write(path, data)
	require.NoError(s.t, err)
	return secret
}

func (s *Session) MustRead(path string) *api.Secret {
	s.t.Helper()

	secret, err := s.Client.Logical().Read(path)
	require.NoError(s.t, err)
	return secret
}

// MustReadRequired is a stricter version of MustRead that fails if a 404/nil is returned
func (s *Session) MustReadRequired(path string) *api.Secret {
	s.t.Helper()

	secret := s.MustRead(path)
	require.NotNil(s.t, secret)
	return secret
}

func (s *Session) MustWriteKV2(mountPath, secretPath string, data map[string]any) {
	s.t.Helper()

	fullPath := path.Join(mountPath, "data", secretPath)
	payload := map[string]any{
		"data": data,
	}
	s.MustWrite(fullPath, payload)
}

func (s *Session) MustReadKV2(mountPath, secretPath string) *api.Secret {
	s.t.Helper()

	fullPath := path.Join(mountPath, "data", secretPath)
	return s.MustRead(fullPath)
}
