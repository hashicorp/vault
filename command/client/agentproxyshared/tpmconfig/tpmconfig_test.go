// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package tpmconfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func newValidTPMCert(t *testing.T) *TPMCert {
	t.Helper()

	return &TPMCert{
		Name:      "machine",
		MountPath: "auth/cert",
		CSR: &TPMCertCSR{
			CommonName: "api.internal.example.com",
			Country:    "US",
		},
		TPM:         &TPMCertTPM{DevicePath: "/dev/tpmrm0"},
		Renew:       &TPMCertRenew{},
		Destination: &TPMCertDestination{Path: t.TempDir()},
	}
}

// TestValidate verifies tpm_cert validation shared by Agent and Proxy:
// renewal range checks, duplicate destination detection (including through
// symlinks), and defaults applied to optional fields.
func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(t *testing.T) ([]*TPMCert, string)
		wantErr string
	}{
		{
			name: "valid_defaults",
			setup: func(t *testing.T) ([]*TPMCert, string) {
				return []*TPMCert{newValidTPMCert(t)}, "root/team-a"
			},
		},
		{
			name: "missing_mount_path",
			setup: func(t *testing.T) ([]*TPMCert, string) {
				entry := newValidTPMCert(t)
				entry.MountPath = ""
				return []*TPMCert{entry}, ""
			},
			wantErr: "tpm_cert[0].mount_path must be specified",
		},
		{
			name: "invalid_percent_renew_before_expiry",
			setup: func(t *testing.T) ([]*TPMCert, string) {
				entry := newValidTPMCert(t)
				pct := 100
				entry.Renew.PercentRenewBeforeExpiry = &pct
				return []*TPMCert{entry}, ""
			},
			wantErr: "tpm_cert[0].renew.percent_renew_before_expiry must be between 1 and 99",
		},
		{
			name: "duplicate_destination_path",
			setup: func(t *testing.T) ([]*TPMCert, string) {
				baseDir := t.TempDir()
				linkDir := filepath.Join(t.TempDir(), "linked-dest")
				if err := os.Symlink(baseDir, linkDir); err != nil {
					t.Skipf("symlink unsupported: %v", err)
				}
				entry1 := newValidTPMCert(t)
				entry1.Name = "machine-1"
				entry1.Destination.Path = baseDir
				entry2 := newValidTPMCert(t)
				entry2.Name = "machine-2"
				entry2.Destination.Path = linkDir
				return []*TPMCert{entry1, entry2}, ""
			},
			wantErr: "destination.path duplicates",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			entries, inheritedNS := tc.setup(t)
			err := Validate(entries, inheritedNS)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, "root/team-a/", entries[0].Namespace)
			require.Equal(t, 20, *entries[0].Renew.PercentRenewBeforeExpiry)
			require.Equal(t, "077", entries[0].Destination.Umask)
		})
	}
}

// TestValidateAutoAuthReference verifies that a tpm auto_auth method must
// reference an existing tpm_cert block, and that other method types are
// ignored.
func TestValidateAutoAuthReference(t *testing.T) {
	t.Parallel()

	certs := []*TPMCert{{Name: "machine", Destination: &TPMCertDestination{Path: "/x"}}}

	tests := []struct {
		name       string
		methodType string
		cfg        map[string]interface{}
		certs      []*TPMCert
		wantErr    string
	}{
		{name: "valid_reference", methodType: "tpm", cfg: map[string]interface{}{"tpm_cert": "machine"}, certs: certs},
		{name: "reference_not_found", methodType: "tpm", cfg: map[string]interface{}{"tpm_cert": "nope"}, certs: certs, wantErr: "no such tpm_cert block exists"},
		{name: "reference_absent", methodType: "tpm", cfg: map[string]interface{}{}, certs: certs, wantErr: "missing tpm_cert config"},
		{name: "non_tpm_method_ignored", methodType: "approle", cfg: map[string]interface{}{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateAutoAuthReference(tc.methodType, tc.cfg, tc.certs)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

// TestConfigureAutoAuth verifies that the referenced tpm_cert block's
// destination, device, and role are injected into the auto_auth method config,
// and that a conflicting explicit role_name is rejected.
func TestConfigureAutoAuth(t *testing.T) {
	t.Parallel()

	certs := []*TPMCert{{
		Name:        "machine",
		RoleName:    "web-tier",
		TPM:         &TPMCertTPM{DevicePath: "/dev/tpmrm0"},
		Destination: &TPMCertDestination{Path: "/var/lib/tpm"},
	}}

	t.Run("injects_fields", func(t *testing.T) {
		t.Parallel()
		cfg := map[string]interface{}{"tpm_cert": "machine"}
		require.NoError(t, ConfigureAutoAuth("tpm", cfg, certs))
		require.Equal(t, "/var/lib/tpm", cfg["tpm_cert_destination"])
		require.Equal(t, "/dev/tpmrm0", cfg["tpm_cert_device"])
		require.Equal(t, "web-tier", cfg["role_name"])
	})

	t.Run("matching_role_name_tolerated", func(t *testing.T) {
		t.Parallel()
		cfg := map[string]interface{}{"tpm_cert": "machine", "role_name": "web-tier"}
		require.NoError(t, ConfigureAutoAuth("tpm", cfg, certs))
	})

	t.Run("conflicting_role_name", func(t *testing.T) {
		t.Parallel()
		cfg := map[string]interface{}{"tpm_cert": "machine", "role_name": "other"}
		require.ErrorContains(t, ConfigureAutoAuth("tpm", cfg, certs), "conflicts with role_name")
	})

	t.Run("missing_destination", func(t *testing.T) {
		t.Parallel()
		noDest := []*TPMCert{{Name: "machine", RoleName: "web-tier"}}
		cfg := map[string]interface{}{"tpm_cert": "machine"}
		require.NotPanics(t, func() {
			require.ErrorContains(t, ConfigureAutoAuth("tpm", cfg, noDest), "destination block must be specified")
		})
	})

	t.Run("non_tpm_method_is_noop", func(t *testing.T) {
		t.Parallel()
		cfg := map[string]interface{}{}
		require.NoError(t, ConfigureAutoAuth("approle", cfg, certs))
		require.Empty(t, cfg)
	})
}
