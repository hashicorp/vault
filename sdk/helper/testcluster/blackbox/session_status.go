// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

func (s *Session) AssertUnsealed(expectedType string) {
	s.t.Helper()

	status, err := s.Client.Sys().SealStatus()
	require.NoError(s.t, err)

	if status.Sealed {
		s.t.Fatal("Vault is sealed")
	}

	if expectedType != "" {
		require.Equal(s.t, expectedType, status.Type, "unexpected seal type")
	}
}

// AssertUnsealedAny verifies that the cluster is unsealed regardless of seal type.
// This is useful for environments where the seal type may vary (e.g., HCP uses awskms, Docker uses shamir).
func (s *Session) AssertUnsealedAny() {
	s.t.Helper()

	status, err := s.Client.Sys().SealStatus()
	require.NoError(s.t, err)

	if status.Sealed {
		s.t.Fatal("Vault is sealed")
	}

	s.t.Logf("Vault is unsealed (seal type: %s)", status.Type)
}

func (s *Session) AssertCLIVersion(version, sha, buildDate, edition string) {
	s.t.Helper()

	// make sure the binary exists first
	_, err := exec.LookPath("vault")
	require.NoError(s.t, err)

	cmd := exec.Command("vault", "version")
	out, err := cmd.CombinedOutput()
	require.NoError(s.t, err)

	output := string(out)

	expectedVersion := fmt.Sprintf("Vault v%s ('%s'), built %s", version, sha, buildDate)

	switch edition {
	case "ce", "ent":
	case "ent.hsm", "ent.fips1403", "ent.hsm.fips1403":
		expectedVersion += " (cgo)"
	default:
		s.t.Fatalf("unknown Vault edition: %s", edition)
	}

	if !strings.Contains(output, expectedVersion) {
		s.t.Fatalf("CLI version mismatch. expected %s. got %s", expectedVersion, output)
	}
}

func (s *Session) AssertServerVersion(version string) {
	s.t.Helper()

	// strip off any version metadata
	b, _, _ := strings.Cut(version, "+")
	expectedVersion, _, _ := strings.Cut(b, "-")

	secret, err := s.Client.Logical().List("sys/version-history")
	require.NoError(s.t, err)

	keysRaw, ok := secret.Data["keys"].([]any)
	if !ok {
		s.t.Fatal("sys/version-history missing 'keys'")
	}

	found := false
	for _, k := range keysRaw {
		if kStr, ok := k.(string); ok && kStr == expectedVersion {
			found = true
			break
		}
	}

	if !found {
		s.t.Fatalf("expected to find %s in version history but didn't", expectedVersion)
	}
}

func (s *Session) AssertReplicationDisabled() {
	s.assertReplicationStatus("ce", "disabled")
}

func (s *Session) AssertDRReplicationStatus(expectedMode string) {
	s.assertReplicationStatus("dr", expectedMode)
}

func (s *Session) AssertPerformanceReplicationStatus(expectedMode string) {
	s.assertReplicationStatus("performance", expectedMode)
}

func (s *Session) assertReplicationStatus(which, expectedMode string) {
	s.t.Helper()

	secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
		return s.Client.Logical().Read("sys/replication/status")
	})

	require.NoError(s.t, err)
	require.NotNil(s.t, secret)

	data := s.AssertSecret(secret).Data()

	if which == "ce" {
		data.HasKey("mode", "disabled")
	} else {
		data.GetMap(which).HasKey("mode", expectedMode)
	}
}
