// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

func (s *Session) MustEnableSecretsEngine(path string, input *api.MountInput) {
	s.t.Helper()

	err := s.Client.Sys().Mount(path, input)
	require.NoError(s.t, err)
}

func (s *Session) MustDisableSecretsEngine(path string) {
	s.t.Helper()

	err := s.Client.Sys().Unmount(path)
	require.NoError(s.t, err)
}

func (s *Session) MustEnableAuth(path string, options *api.EnableAuthOptions) {
	s.t.Helper()

	err := s.Client.Sys().EnableAuthWithOptions(path, options)
	require.NoError(s.t, err)
}

func (s *Session) MustWritePolicy(name, rules string) {
	s.t.Helper()

	err := s.Client.Sys().PutPolicy(name, rules)
	require.NoError(s.t, err)
}
