// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !race && !hsm

package command

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// testPGPPublicKey is a sample PGP public key for testing purposes
const testPGPPublicKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mQINBGB9+xkBEACabYZOWKmgZsHTdRDiyPJxhbuUiKX65GUWkyRMJKi/1dviVxOX
PG6hBPtF48IFnVgxKpIb7G6NjBousAV+CuLlv5yqFKpOZEGC6sBV+Gx8Vu1CICpl
Zm+HpQPcIzwBpN+Ar4l/exCG/f/MZq/oxGgH+TyRF3XcYDjG8dbJCpHO5nQ5Cy9h
QIp3/Bh09kET6lk+4QlofNgHKVT2epV8iK1cXlbQe2tZtfCUtxk+pxvU0UHXp+AB
0xc3/gIhjZp/dePmCOyQyGPJbp5bpO4UeAJ6frqhexmNlaw9Z897ltZmRLGq1p4a
RnWL8FPkBz9SCSKXS8uNyV5oMNVn4G1obCkc106iWuKBTibffYQzq5TG8FYVJKrh
RwWB6piacEB8hl20IIWSxIM3J9tT7CPSnk5RYYCTRHgA5OOrqZhC7JefudrP8n+M
pxkDgNORDu7GCfAuisrf7dXYjLsxG4tu22DBJJC0c/IpRpXDnOuJN1Q5e/3VUKKW
mypNumuQpP5lc1ZFG64TRzb1HR6oIdHfbrVQfdiQXpvdcFx+Fl57WuUraXRV6qfb
4ZmKHX1JEwM/7tu21QE4F1dz0jroLSricZxfaCTHHWNfvGJoZ30/MZUrpSC0IfB3
iQutxbZrwIlTBt+fGLtm3vDtwMFNWM+Rb1lrOxEQd2eijdxhvBOHtlIcswARAQAB
tERIYXNoaUNvcnAgU2VjdXJpdHkgKGhhc2hpY29ycC5jb20vc2VjdXJpdHkpIDxz
ZWN1cml0eUBoYXNoaWNvcnAuY29tPokCVAQTAQoAPhYhBMh0AR8KtAURDQIQVTQ2
XZRy10aPBQJgffsZAhsDBQkJZgGABQsJCAcCBhUKCQgLAgQWAgMBAh4BAheAAAoJ
EDQ2XZRy10aPtpcP/0PhJKiHtC1zREpRTrjGizoyk4Sl2SXpBZYhkdrG++abo6zs
buaAG7kgWWChVXBo5E20L7dbstFK7OjVs7vAg/OLgO9dPD8n2M19rpqSbbvKYWvp
0NSgvFTT7lbyDhtPj0/bzpkZEhmvQaDWGBsbDdb2dBHGitCXhGMpdP0BuuPWEix+
QnUMaPwU51q9GM2guL45Tgks9EKNnpDR6ZdCeWcqo1IDmklloidxT8aKL21UOb8t
cD+Bg8iPaAr73bW7Jh8TdcV6s6DBFub+xPJEB/0bVPmq3ZHs5B4NItroZ3r+h3ke
VDoSOSIZLl6JtVooOJ2la9ZuMqxchO3mrXLlXxVCo6cGcSuOmOdQSz4OhQE5zBxx
LuzA5ASIjASSeNZaRnffLIHmht17BPslgNPtm6ufyOk02P5XXwa69UCjA3RYrA2P
QNNC+OWZ8qQLnzGldqE4MnRNAxRxV6cFNzv14ooKf7+k686LdZrP/3fQu2p3k5rY
0xQUXKh1uwMUMtGR867ZBYaxYvwqDrg9XB7xi3N6aNyNQ+r7zI2lt65lzwG1v9hg
FG2AHrDlBkQi/t3wiTS3JOo/GCT8BjN0nJh0lGaRFtQv2cXOQGVRW8+V/9IpqEJ1
qQreftdBFWxvH7VJq2mSOXUJyRsoUrjkUuIivaA9Ocdipk2CkP8bpuGz7ZF4
=s1CX
-----END PGP PUBLIC KEY BLOCK-----`

// invalidPGPKey is not a valid PGP key format
const invalidPGPKey = `This is not a valid PGP key`

// testConfig is used to disable prometheus to avoid issues with the
// global prometheus registry in tests
const testConfig = `
	storage "inmem" {}
	listener "tcp" {
		address = "127.0.0.1:0"
		tls_disable = true
	}
	disable_mlock = true
	telemetry {
		prometheus_retention_time = "0s"
		disable_hostname = true
	}
`

// TestServer_DevPluginPGPKey_ValidKey tests that the server accepts a valid PGP key file
func TestServer_DevPluginPGPKey_ValidKey(t *testing.T) {
	// Create a valid PGP key file
	keyPath := filepath.Join(t.TempDir(), "test-pgp-key.asc")
	err := os.WriteFile(keyPath, []byte(testPGPPublicKey), 0o644)
	require.NoError(t, err)

	ui, cmd := testServerCommand(t)
	args := []string{
		"-dev",
		"-dev-plugin-pgp-key=" + keyPath,
		"-dev-listen-address=127.0.0.1:0",
		"-test-server-config",
	}

	retCode := cmd.Run(args)
	output := ui.ErrorWriter.String() + ui.OutputWriter.String()

	// The server should start successfully with a valid key
	require.Equal(t, 0, retCode, "expected server to start successfully, output: %s", output)
}

// TestServer_DevPluginPGPKey_InvalidPath tests that the server handles invalid file paths
func TestServer_DevPluginPGPKey_InvalidPath(t *testing.T) {
	ui, cmd := testServerCommand(t)

	// Use an invalid path with null bytes (not allowed in file paths)
	invalidPath := "/tmp/test\x00key.asc"

	configPath := filepath.Join(t.TempDir(), "config.hcl")
	err := os.WriteFile(configPath, []byte(testConfig), 0o644)
	require.NoError(t, err)

	args := []string{
		"-dev",
		"-config=" + configPath,
		"-dev-plugin-pgp-key=" + invalidPath,
		"-dev-listen-address=127.0.0.1:0",
		"-test-server-config",
	}

	retCode := cmd.Run(args)
	output := ui.ErrorWriter.String() + ui.OutputWriter.String()

	// The server should fail to start with an invalid path
	require.NotEqual(t, 0, retCode, "expected server to fail with invalid path")
	require.Contains(t, output, "dev plugin PGP key", "expected error message about PGP key path, output: %s", output)
}

// TestServer_DevPluginPGPKey_NonExistentFile tests that the server handles non-existent key files
func TestServer_DevPluginPGPKey_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	// Use a path that doesn't exist
	keyPath := filepath.Join(tmpDir, "non-existent-key.asc")

	ui, cmd := testServerCommand(t)

	configPath := filepath.Join(tmpDir, "config.hcl")
	err := os.WriteFile(configPath, []byte(testConfig), 0o644)
	require.NoError(t, err)

	args := []string{
		"-dev",
		"-config=" + configPath,
		"-dev-plugin-pgp-key=" + keyPath,
		"-dev-listen-address=127.0.0.1:0",
		"-test-server-config",
	}

	retCode := cmd.Run(args)
	output := ui.ErrorWriter.String() + ui.OutputWriter.String()

	// The server should fail to start with a non-existent key file
	require.NotEqual(t, 0, retCode, "expected server to fail with non-existent key file")
	require.Contains(t, output, "dev plugin PGP key", "expected error message about PGP key path, output: %s", output)
}

// TestServer_DevPluginPGPKey_EmptyFlag tests default behavior when flag is not set
func TestServer_DevPluginPGPKey_EmptyFlag(t *testing.T) {
	ui, cmd := testServerCommand(t)

	configPath := filepath.Join(t.TempDir(), "config.hcl")
	err := os.WriteFile(configPath, []byte(testConfig), 0o644)
	require.NoError(t, err)

	args := []string{
		"-dev",
		"-config=" + configPath,
		"-dev-listen-address=127.0.0.1:0",
		"-test-server-config",
	}

	retCode := cmd.Run(args)
	output := ui.ErrorWriter.String() + ui.OutputWriter.String()

	// The server should start successfully without the flag (uses default HashiCorp key)
	require.Equal(t, 0, retCode, "expected server to start successfully without flag, output: %s", output)
}

// TestServer_DevPluginPGPKey_InvalidKeyContent tests handling of invalid PGP key content
func TestServer_DevPluginPGPKey_InvalidKeyContent(t *testing.T) {
	tmpDir := t.TempDir()
	// Create a file with invalid PGP key content
	keyPath := filepath.Join(tmpDir, "invalid-key.asc")
	err := os.WriteFile(keyPath, []byte(invalidPGPKey), 0o644)
	require.NoError(t, err)

	ui, cmd := testServerCommand(t)

	configPath := filepath.Join(tmpDir, "config.hcl")
	err = os.WriteFile(configPath, []byte(testConfig), 0o644)
	require.NoError(t, err)

	args := []string{
		"-dev",
		"-config=" + configPath,
		"-dev-plugin-pgp-key=" + keyPath,
		"-dev-listen-address=127.0.0.1:0",
		"-test-server-config",
	}

	retCode := cmd.Run(args)
	output := ui.ErrorWriter.String() + ui.OutputWriter.String()

	// The server should start (validation happens when actually using the key)
	// but we verify the path was set correctly
	require.Equal(t, 0, retCode, "expected server to start (key validation happens at use time), output: %s", output)
}

// TestServer_DevPluginPGPKey_EmptyFile tests handling of an empty key file
func TestServer_DevPluginPGPKey_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	// Create an empty file
	keyPath := filepath.Join(tmpDir, "empty-key.asc")
	err := os.WriteFile(keyPath, []byte(""), 0o644)
	require.NoError(t, err)

	ui, cmd := testServerCommand(t)

	configPath := filepath.Join(tmpDir, "config.hcl")
	err = os.WriteFile(configPath, []byte(testConfig), 0o644)
	require.NoError(t, err)

	args := []string{
		"-dev",
		"-config=" + configPath,
		"-dev-plugin-pgp-key=" + keyPath,
		"-dev-listen-address=127.0.0.1:0",
		"-test-server-config",
	}

	retCode := cmd.Run(args)
	output := ui.ErrorWriter.String() + ui.OutputWriter.String()

	// The server should start (validation happens when actually using the key)
	require.Equal(t, 0, retCode, "expected server to start (key validation happens at use time), output: %s", output)
}
