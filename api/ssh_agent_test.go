// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSSH_CanLoadDuplicateKeys verifies that during the deprecation process of duplicate HCL attributes this function
// will still allow them.
// TODO (HCL_DUP_KEYS_DEPRECATION): on full removal change this test to ensure that duplicate attributes cannot be parsed
// under any circumstances.
func TestSSH_CanLoadDuplicateKeys(t *testing.T) {
	t.Run("fail parsing without env var", func(t *testing.T) {
		_, err := LoadSSHHelperConfig("./test-fixtures/agent_config_duplicate_keys.hcl")
		require.Error(t, err)
		require.Contains(t, err.Error(), "Each argument can only be defined once")
	})
	t.Run("fail parsing with env var set to false", func(t *testing.T) {
		t.Setenv(allowHclDuplicatesEnvVar, "false")
		_, err := LoadSSHHelperConfig("./test-fixtures/agent_config_duplicate_keys.hcl")
		require.Error(t, err)
		require.Contains(t, err.Error(), "Each argument can only be defined once")
	})
	t.Run("succeed parsing with env var set to true", func(t *testing.T) {
		t.Setenv(allowHclDuplicatesEnvVar, "true")
		_, err := LoadSSHHelperConfig("./test-fixtures/agent_config_duplicate_keys.hcl")
		require.NoError(t, err)
	})
}

func TestSSH_CreateTLSClient(t *testing.T) {
	// load the default configuration
	config, err := LoadSSHHelperConfig("./test-fixtures/agent_config.hcl")
	if err != nil {
		panic(fmt.Sprintf("error loading agent's config file: %s", err))
	}

	client, err := config.NewClient()
	if err != nil {
		panic(fmt.Sprintf("error creating the client: %s", err))
	}

	// Provide a certificate and enforce setting of transport
	config.CACert = "./test-fixtures/vault.crt"

	client, err = config.NewClient()
	if err != nil {
		panic(fmt.Sprintf("error creating the client: %s", err))
	}
	if client.config.HttpClient.Transport == nil {
		panic(fmt.Sprintf("error creating client with TLS transport"))
	}
}

func TestSSH_CreateTLSClient_tlsServerName(t *testing.T) {
	// Ensure that the HTTP client is associated with the configured TLS server name.
	tlsServerName := "tls.server.name"

	config, err := ParseSSHHelperConfig(fmt.Sprintf(`
vault_addr = "1.2.3.4"
tls_server_name = "%s"
`, tlsServerName))
	if err != nil {
		panic(fmt.Sprintf("error loading config: %s", err))
	}

	client, err := config.NewClient()
	if err != nil {
		panic(fmt.Sprintf("error creating the client: %s", err))
	}

	actualTLSServerName := client.config.HttpClient.Transport.(*http.Transport).TLSClientConfig.ServerName
	if actualTLSServerName != tlsServerName {
		panic(fmt.Sprintf("incorrect TLS server name. expected: %s actual: %s", tlsServerName, actualTLSServerName))
	}
}

func TestParseSSHHelperConfig(t *testing.T) {
	config, err := ParseSSHHelperConfig(`
		vault_addr = "1.2.3.4"
`)
	if err != nil {
		t.Fatal(err)
	}

	if config.SSHMountPoint != SSHHelperDefaultMountPoint {
		t.Errorf("expected %q to be %q", config.SSHMountPoint, SSHHelperDefaultMountPoint)
	}
}

func TestParseSSHHelperConfig_missingVaultAddr(t *testing.T) {
	_, err := ParseSSHHelperConfig("")
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), `missing config "vault_addr"`) {
		t.Errorf("bad error: %s", err)
	}
}

func TestParseSSHHelperConfig_badKeys(t *testing.T) {
	_, err := ParseSSHHelperConfig(`
vault_addr = "1.2.3.4"
nope = "bad"
`)
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), `ssh_helper: invalid key "nope" on line 3`) {
		t.Errorf("bad error: %s", err)
	}
}

func TestParseSSHHelperConfig_tlsServerName(t *testing.T) {
	tlsServerName := "tls.server.name"

	config, err := ParseSSHHelperConfig(fmt.Sprintf(`
vault_addr = "1.2.3.4"
tls_server_name = "%s"
`, tlsServerName))
	if err != nil {
		t.Fatal(err)
	}

	if config.TLSServerName != tlsServerName {
		t.Errorf("incorrect TLS server name. expected: %s actual: %s", tlsServerName, config.TLSServerName)
	}
}
