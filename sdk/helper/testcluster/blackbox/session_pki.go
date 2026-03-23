// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"path"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

// MustSetupPKIRoot bootstraps a PKI engine as a Root CA in one shot.
// It returns the role name you can use to issue certs immediately.
func (s *Session) MustSetupPKIRoot(mountPath string) string {
	s.t.Helper()

	s.MustEnableSecretsEngine(mountPath, &api.MountInput{Type: "pki"})

	// Root CA generation often fails if MaxTTL < requested TTL
	err := s.Client.Sys().TuneMount(mountPath, api.MountConfigInput{
		MaxLeaseTTL: "87600h",
	})
	require.NoError(s.t, err)

	s.MustWrite(path.Join(mountPath, "root/generate/internal"), map[string]any{
		"common_name": "vault-test-root",
		"ttl":         "8760h",
	})

	roleName := "server-cert"
	s.MustWrite(path.Join(mountPath, "roles", roleName), map[string]any{
		"allowed_domains":  "example.com",
		"allow_subdomains": true,
		"max_ttl":          "72h",
	})

	return roleName
}
