// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	base "github.com/hashicorp/vault/command/base"
	proxyConfig "github.com/hashicorp/vault/command/client/proxy/config"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProxy_LogFile_CliOverridesConfig tests that the CLI values
// override the config for log files
func TestProxy_LogFile_CliOverridesConfig(t *testing.T) {
	t.Parallel()

	// Create basic config
	configFile := populateCommandTestFile(t, "proxy-config.hcl", basicHCLConfig)
	cfg, err := proxyConfig.LoadConfigFile(configFile.Name())
	require.NoError(t, err)

	// Sanity check that the config value is the current value
	require.Equal(t, "TMPDIR/juan.log", cfg.LogFile)

	// Initialize the command and parse any flags
	cmd := &ProxyCommand{BaseCommand: &base.BaseCommand{}}
	flags := cmd.Flags()
	// Simulate the flag being specified
	require.NoError(t, flags.Parse([]string{"-log-file=/foo/bar/test.log"}))

	// Update the config based on the inputs.
	cmd.applyConfigOverrides(flags, cfg)

	assert.NotEqual(t, "TMPDIR/juan.log", cfg.LogFile)
	assert.NotEqual(t, "/squiggle/logs.txt", cfg.LogFile)
	assert.Equal(t, "/foo/bar/test.log", cfg.LogFile)
}

// TestProxy_LogFile_Config tests log file config when loaded from config
func TestProxy_LogFile_Config(t *testing.T) {
	t.Parallel()

	configFile := populateCommandTestFile(t, "proxy-config.hcl", basicHCLConfig)
	cfg, err := proxyConfig.LoadConfigFile(configFile.Name())
	require.NoError(t, err)

	// Sanity check that the config value is the current value
	require.Equal(t, "TMPDIR/juan.log", cfg.LogFile)
	require.Equal(t, 2, cfg.LogRotateMaxFiles)
	require.Equal(t, 1048576, cfg.LogRotateBytes)

	// Parse the cli flags (but we pass in an empty slice)
	cmd := &ProxyCommand{BaseCommand: &base.BaseCommand{}}
	flags := cmd.Flags()
	require.NoError(t, flags.Parse(nil))

	// Should change nothing...
	cmd.applyConfigOverrides(flags, cfg)

	assert.Equal(t, "TMPDIR/juan.log", cfg.LogFile)
	assert.Equal(t, 2, cfg.LogRotateMaxFiles)
	assert.Equal(t, 1048576, cfg.LogRotateBytes)
}

// TestProxy_EnvVar_Overrides tests that environment variables are properly
// parsed and override defaults.
func TestProxy_EnvVar_Overrides(t *testing.T) {
	configFile := populateCommandTestFile(t, "proxy-config.hcl", basicHCLConfig)
	cfg, err := proxyConfig.LoadConfigFile(configFile.Name())
	require.NoError(t, err)
	require.False(t, cfg.Vault.TLSSkipVerify)

	t.Setenv(api.EnvVaultSkipVerify, "true")
	// Parse the cli flags (but we pass in an empty slice)
	cmd := &ProxyCommand{BaseCommand: &base.BaseCommand{}}
	flags := cmd.Flags()
	require.NoError(t, flags.Parse(nil))
	cmd.applyConfigOverrides(flags, cfg)
	require.True(t, cfg.Vault.TLSSkipVerify)

	t.Setenv(api.EnvVaultSkipVerify, "false")
	cmd.applyConfigOverrides(flags, cfg)
	require.False(t, cfg.Vault.TLSSkipVerify)
}

// TestProxy_Config_NewLogger_Default Tests defaults for log level and
// specifically cmd.newLogger()
func TestProxy_Config_NewLogger_Default(t *testing.T) {
	t.Parallel()

	cmd := &ProxyCommand{BaseCommand: &base.BaseCommand{}, config: proxyConfig.NewConfig()}
	logger, err := cmd.newLogger()

	require.NoError(t, err)
	require.NotNil(t, logger)
	require.Equal(t, hclog.Info.String(), logger.GetLevel().String())
}

// TestProxy_Config_ReloadLogLevel Tests reloading updates the log
// level as expected.
func TestProxy_Config_ReloadLogLevel(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	// Load an initial config
	configFile := populateCommandTestFile(t, "proxy-config.hcl", strings.ReplaceAll(basicHCLConfig, "TMPDIR", tempDir))
	config, err := proxyConfig.LoadConfigFile(configFile.Name())
	require.NoError(t, err)

	// Tweak the loaded config to make sure we can put log files into a temp dir
	// and systemd log attempts work fine, this would usually happen during Run.
	cmd := &ProxyCommand{
		BaseCommand: &base.BaseCommand{},
		config:      config,
		logWriter:   os.Stdout,
	}
	cmd.logger, err = cmd.newLogger()
	require.NoError(t, err)

	// Sanity check
	require.Equal(t, "warn", cmd.config.LogLevel)

	// Load a new config
	configFile = populateCommandTestFile(t, "proxy-config.hcl", strings.ReplaceAll(basicHCLConfig2, "TMPDIR", tempDir))
	require.NoError(t, cmd.reloadConfig([]string{configFile.Name()}))
	require.Equal(t, "debug", cmd.config.LogLevel)
}

// TestProxy_Config_ReloadTls Tests that the TLS certs for the listener are
// correctly reloaded.
func TestProxy_Config_ReloadTls(t *testing.T) {
	var wg sync.WaitGroup
	workingDir := filepath.Join("proxy", "test-fixtures", "reload")
	tempDir := t.TempDir()

	// Set up initial 'foo' certs
	fooCert, err := os.ReadFile(filepath.Join(workingDir, "reload_foo.pem"))
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "reload_cert.pem"), fooCert, 0o600))
	fooKey, err := os.ReadFile(filepath.Join(workingDir, "reload_foo.key"))
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "reload_key.pem"), fooKey, 0o600))

	caPEM, err := os.ReadFile(filepath.Join(workingDir, "reload_ca.pem"))
	require.NoError(t, err)
	certPool := x509.NewCertPool()
	require.True(t, certPool.AppendCertsFromPEM(caPEM))

	configFile := populateCommandTestFile(t, "proxy-config.hcl", strings.ReplaceAll(basicHCLConfig, "TMPDIR", tempDir))

	// Set up Proxy
	ui, cmd, control := MakeTestProxyCommandWithControl(logging.NewVaultLogger(hclog.Trace), nil)

	// Start
	var code int
	wg.Add(1)
	go func() {
		defer wg.Done()
		code = cmd.Run([]string{"-config", configFile.Name()})
	}()

	select {
	case <-control.StartedCh:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out starting Proxy: %s%s", ui.ErrorWriter.String(), ui.OutputWriter.String())
	}

	checkCertificateName := func(commonName string) error {
		conn, err := tls.Dial("tcp", "127.0.0.1:8100", &tls.Config{RootCAs: certPool})
		if err != nil {
			return err
		}
		defer conn.Close()
		if err := conn.Handshake(); err != nil {
			return err
		}
		actual := conn.ConnectionState().PeerCertificates[0].Subject.CommonName
		if actual != commonName {
			return fmt.Errorf("expected %s, got %s", commonName, actual)
		}
		return nil
	}
	require.NoError(t, checkCertificateName("foo.example.com"))

	// Swap out certs
	barCert, err := os.ReadFile(filepath.Join(workingDir, "reload_bar.pem"))
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "reload_cert.pem"), barCert, 0o600))
	barKey, err := os.ReadFile(filepath.Join(workingDir, "reload_bar.key"))
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "reload_key.pem"), barKey, 0o600))

	// Reload
	cmd.SighupCh <- struct{}{}
	select {
	case <-control.ReloadedCh:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out reloading Proxy")
	}
	require.NoError(t, checkCertificateName("bar.example.com"))

	// Shut down
	cmd.ShutdownCh <- struct{}{}
	wg.Wait()
	require.Equalf(t, 0, code, "Proxy output: %s%s", ui.ErrorWriter.String(), ui.OutputWriter.String())
}

// TestProxy_Config_AddrConformance verifies that the vault address is correctly
// normalized to conform to RFC-5942 §4 when configured by a config file,
// environment variables, or CLI flags.
// See: https://rfc-editor.org/rfc/rfc5952.html
func TestProxy_Config_AddrConformance(t *testing.T) {
	for name, test := range map[string]struct {
		args     []string
		envVars  map[string]string
		cfg      string
		expected string
	}{
		"ipv4 config": {
			cfg:      `vault { address = "https://127.0.0.1:8200" }`,
			expected: "https://127.0.0.1:8200",
		},
		"ipv6 config": {
			cfg: `vault { address = "https://[2001:0db8::0001]:8200" }`,
			// Use the normalized version in the config
			expected: "https://[2001:db8::1]:8200",
		},
		"ipv6 cli arg overrides": {
			args: []string{"-address=https://[2001:0:0:1:0:0:0:1]:8200"},
			cfg:  `vault { address = "https://[2001:0db8::0001]:8200" }`,
			// Use a normalized version of the args address
			expected: "https://[2001:0:0:1::1]:8200",
		},
		"ipv6 env var overrides": {
			envVars: map[string]string{api.EnvVaultAddress: "https://[2001:DB8:AC3:FE4::1]:8200"},
			cfg:     `vault { address = "https://[2001:0db8::0001]:8200" }`,
			// Use a normalized version of the env var address
			expected: "https://[2001:db8:ac3:fe4::1]:8200",
		},
		"ipv6 all uses cli overrides": {
			args:     []string{"-address=https://[2001:0:0:1:0:0:0:1]:8200"},
			envVars:  map[string]string{api.EnvVaultAddress: "https://[2001:DB8:AC3:FE4::1]:8200"},
			cfg:      `vault { address = "https://[2001:0db8::0001]:8200" }`,
			expected: "https://[2001:0:0:1::1]:8200",
		},
	} {
		t.Run(name, func(t *testing.T) {
			// In CI our tests are run with VAULT_ADDR=, which will break our tests
			// because it'll default to an unset address. Ensure that's cleared out
			// of the environment.
			// t.Setenv can't represent "unset" - it always leaves the var present,
			// even with an empty value - and code that reads it via os.LookupEnv
			// (rather than checking for an empty string) would treat that as
			// explicitly set. Capture prior state and restore it via Cleanup instead.
			origAddr, hadAddr := os.LookupEnv(api.EnvVaultAddress)
			require.NoError(t, os.Unsetenv(api.EnvVaultAddress))
			t.Cleanup(func() {
				if hadAddr {
					_ = os.Setenv(api.EnvVaultAddress, origAddr)
					return
				}
				_ = os.Unsetenv(api.EnvVaultAddress)
			})
			for key, value := range test.envVars {
				t.Setenv(key, value)
			}

			configFile := populateCommandTestFile(t, "proxy-"+strings.ReplaceAll(name, " ", "-"), test.cfg)
			cfg, err := proxyConfig.LoadConfigFile(configFile.Name())
			require.NoError(t, err)
			require.NotEmpty(t, cfg.Vault.Address)

			cmd := &ProxyCommand{BaseCommand: &base.BaseCommand{}}
			flags := cmd.Flags()
			require.NoError(t, flags.Parse(append([]string{}, test.args...)))
			cmd.applyConfigOverrides(flags, cfg)
			require.Equal(t, test.expected, cfg.Vault.Address)
		})
	}
}
