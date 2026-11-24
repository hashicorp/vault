// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// loadTPMFixture writes the TPM fixture to a temp file with TMPDIR replaced by
// a real destination directory, then loads it.
func loadTPMFixture(t *testing.T, destDir string) (*Config, error) {
	t.Helper()

	raw, err := os.ReadFile("./test-fixtures/config-auto-auth-tpm.hcl")
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "config.hcl")
	require.NoError(t, os.WriteFile(path, []byte(strings.ReplaceAll(string(raw), "TMPDIR", destDir)), 0o600))

	return LoadConfigFile(path)
}

// TestLoadConfigFile_AutoAuthTPM is the regression test for VAULT-51583: Proxy
// used to reject the "tpm" auto-auth method with "unknown auto-auth method
// type tpm". It verifies the method is accepted, that the tpm_cert block is
// parsed, and that the destination, device, and role from the tpm_cert block
// are injected into the auto-auth method config the way Agent does.
func TestLoadConfigFile_AutoAuthTPM(t *testing.T) {
	t.Parallel()

	destDir := t.TempDir()
	config, err := loadTPMFixture(t, destDir)
	require.NoError(t, err)
	require.NoError(t, config.ValidateConfig())

	require.Len(t, config.TPMCerts, 1)
	entry := config.TPMCerts[0]
	require.Equal(t, "machine", entry.Name)
	require.Equal(t, "auth/tpm", entry.MountPath)
	require.Equal(t, "web-tier", entry.RoleName)
	require.Equal(t, "platform/", entry.Namespace)
	require.Equal(t, destDir, entry.Destination.Path)
	require.Equal(t, 20, *entry.Renew.PercentRenewBeforeExpiry)

	require.Equal(t, "tpm", config.AutoAuth.Method.Type)
	require.Equal(t, "machine", config.AutoAuth.Method.Config["tpm_cert"])
	require.Equal(t, destDir, config.AutoAuth.Method.Config["tpm_cert_destination"])
	require.Equal(t, "/dev/tpmrm0", config.AutoAuth.Method.Config["tpm_cert_device"])
	require.Equal(t, "web-tier", config.AutoAuth.Method.Config["role_name"])
}

// TestValidateConfig_AutoAuthTPMMissingBlock verifies that a tpm auto-auth
// method referencing a tpm_cert block that doesn't exist is rejected.
func TestValidateConfig_AutoAuthTPMMissingBlock(t *testing.T) {
	t.Parallel()

	config := &Config{
		AutoAuth: &AutoAuth{
			Method: &Method{Type: "tpm", Config: map[string]interface{}{"tpm_cert": "nonexistent"}},
			Sinks:  []*Sink{{Type: "file"}},
		},
	}
	require.ErrorContains(t, config.ValidateConfig(), "no such tpm_cert block exists")
}

// TestConfigMerge_TPMCerts verifies tpm_cert blocks from every merged config
// file are retained.
func TestConfigMerge_TPMCerts(t *testing.T) {
	t.Parallel()

	a, err := loadTPMFixture(t, t.TempDir())
	require.NoError(t, err)
	b := NewConfig()

	require.Len(t, a.Merge(b).TPMCerts, 1)
	require.Len(t, b.Merge(a).TPMCerts, 1)
}

// TestLoadConfigFile_AutoAuthTPMMissingDestination verifies that a tpm_cert
// block without a destination block produces a configuration error rather
// than a nil pointer panic, both when loading the file and when validating it.
func TestLoadConfigFile_AutoAuthTPMMissingDestination(t *testing.T) {
	t.Parallel()

	hcl := `
auto_auth {
  method "tpm" {
    mount_path = "auth/tpm"
    config = {
      tpm_cert = "machine"
    }
  }
}

api_proxy {
  use_auto_auth_token = "force"
}

listener "tcp" {
  address     = "127.0.0.1:18100"
  tls_disable = true
}

tpm_cert "machine" {
  mount_path = "auth/tpm"
  role_name  = "web-tier"
}
`
	path := filepath.Join(t.TempDir(), "config.hcl")
	require.NoError(t, os.WriteFile(path, []byte(hcl), 0o600))

	require.NotPanics(t, func() {
		config, err := LoadConfigFile(path)
		require.NoError(t, err)
		require.ErrorContains(t, config.ValidateConfig(), "destination block must be specified")
	})
}

// TestLoadConfigDir_AutoAuthTPMSplitAcrossFiles verifies that an auto_auth tpm
// method and the tpm_cert block it references can live in different config
// files, in either order. The reference is resolved and the destination,
// device, and role are injected only after the files are merged.
func TestLoadConfigDir_AutoAuthTPMSplitAcrossFiles(t *testing.T) {
	t.Parallel()

	autoAuthHCL := `
auto_auth {
  method "tpm" {
    mount_path = "auth/tpm"
    config = {
      tpm_cert = "machine"
    }
  }
}

api_proxy {
  use_auto_auth_token = "force"
}

listener "tcp" {
  address     = "127.0.0.1:18100"
  tls_disable = true
}
`
	tpmCertHCL := func(destDir string) string {
		return `
tpm_cert "machine" {
  mount_path = "auth/tpm"
  role_name  = "web-tier"

  tpm {
    device_path = "/dev/tpmrm0"
  }

  destination {
    path = "` + destDir + `"
  }
}
`
	}

	// Files in a config directory are loaded in directory order, so the file
	// names control which one is loaded first.
	for name, files := range map[string][2]string{
		"auto_auth_first": {"a-auto-auth.hcl", "b-tpm-cert.hcl"},
		"tpm_cert_first":  {"b-auto-auth.hcl", "a-tpm-cert.hcl"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			destDir := t.TempDir()
			dir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(dir, files[0]), []byte(autoAuthHCL), 0o600))
			require.NoError(t, os.WriteFile(filepath.Join(dir, files[1]), []byte(tpmCertHCL(destDir)), 0o600))

			config, err := LoadConfigDir(dir)
			require.NoError(t, err)
			require.NoError(t, config.ValidateConfig())

			require.Len(t, config.TPMCerts, 1)
			require.Equal(t, destDir, config.AutoAuth.Method.Config["tpm_cert_destination"])
			require.Equal(t, "/dev/tpmrm0", config.AutoAuth.Method.Config["tpm_cert_device"])
			require.Equal(t, "web-tier", config.AutoAuth.Method.Config["role_name"])
		})
	}
}
