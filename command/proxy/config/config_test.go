// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/command/agentproxyshared"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/stretchr/testify/require"
)

// TestLoadConfigFile_ProxyCache tests loading a config file containing a cache
// as well as a valid proxy config.
func TestLoadConfigFile_ProxyCache(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-cache.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
			Listeners: []*configutil.Listener{
				{
					Type:        "unix",
					Address:     "/path/to/socket",
					TLSDisable:  true,
					SocketMode:  "configmode",
					SocketUser:  "configuser",
					SocketGroup: "configgroup",
				},
				{
					Type:       "tcp",
					Address:    "127.0.0.1:8300",
					TLSDisable: true,
				},
				{
					Type:       "tcp",
					Address:    "127.0.0.1:3000",
					Role:       "metrics_only",
					TLSDisable: true,
				},
				{
					Type:        "tcp",
					Role:        "default",
					Address:     "127.0.0.1:8400",
					TLSKeyFile:  "/path/to/cakey.pem",
					TLSCertFile: "/path/to/cacert.pem",
				},
			},
		},
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
			Sinks: []*Sink{
				{
					Type:   "file",
					DHType: "curve25519",
					DHPath: "/tmp/file-foo-dhpath",
					AAD:    "foobar",
					Config: map[string]interface{}{
						"path": "/tmp/file-foo",
					},
				},
			},
		},
		APIProxy: &APIProxy{
			EnforceConsistency:  "always",
			WhenInconsistent:    "retry",
			UseAutoAuthTokenRaw: true,
			UseAutoAuthToken:    true,
			ForceAutoAuthToken:  false,
		},
		Cache: &Cache{
			Persist: &agentproxyshared.PersistConfig{
				Type:                    "kubernetes",
				Path:                    "/vault/agent-cache/",
				KeepAfterImport:         true,
				ExitOnErr:               true,
				ServiceAccountTokenFile: "/tmp/serviceaccount/token",
			},
		},
		Vault: &Vault{
			Address:          "http://127.0.0.1:1111",
			CACert:           "config_ca_cert",
			CAPath:           "config_ca_path",
			TLSSkipVerifyRaw: interface{}("true"),
			TLSSkipVerify:    true,
			ClientCert:       "config_client_cert",
			ClientKey:        "config_client_key",
			Retry: &Retry{
				NumRetries: 12,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}

	config, err = LoadConfigFile("./test-fixtures/config-cache-embedded-type.hcl")
	if err != nil {
		t.Fatal(err)
	}
	expected.Vault.TLSSkipVerifyRaw = interface{}(true)

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

// TestLoadConfigFile_NoCachingEnabled tests that you cannot enable a cache
// without either of the options to enable caching secrets
func TestLoadConfigFile_NoCachingEnabled(t *testing.T) {
	cfg, err := LoadConfigFile("./test-fixtures/config-cache-but-no-secrets.hcl")
	if err != nil {
		t.Fatal(err)
	}

	if err := cfg.ValidateConfig(); err == nil {
		t.Fatalf("expected error, as you cannot configure a cache without caching secrets")
	}
}

// TestLoadConfigFile_StaticSecretCachingWithoutAutoAuth tests that loading
// a config file with static secret caching enabled but no auto auth will fail.
func TestLoadConfigFile_StaticSecretCachingWithoutAutoAuth(t *testing.T) {
	cfg, err := LoadConfigFile("./test-fixtures/config-cache-static-no-auto-auth.hcl")
	if err != nil {
		t.Fatal(err)
	}

	if err := cfg.ValidateConfig(); err == nil {
		t.Fatalf("expected error, as static secret caching requires auto-auth")
	}
}

// TestLoadConfigFile_ProxyCacheStaticSecrets tests loading a config file containing a cache
// as well as a valid proxy config with static secret caching enabled
func TestLoadConfigFile_ProxyCacheStaticSecrets(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-cache-static-secret-cache.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
			Listeners: []*configutil.Listener{
				{
					Type:       "tcp",
					Address:    "127.0.0.1:8300",
					TLSDisable: true,
				},
			},
		},
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
			Sinks: []*Sink{
				{
					Type:   "file",
					DHType: "curve25519",
					DHPath: "/tmp/file-foo-dhpath",
					AAD:    "foobar",
					Config: map[string]interface{}{
						"path": "/tmp/file-foo",
					},
				},
			},
		},
		Cache: &Cache{
			CacheStaticSecrets:                         true,
			StaticSecretTokenCapabilityRefreshInterval: 1 * time.Hour,
		},
		Vault: &Vault{
			Address:          "http://127.0.0.1:1111",
			TLSSkipVerify:    true,
			TLSSkipVerifyRaw: interface{}("true"),
			Retry: &Retry{
				NumRetries: 12,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

// Test_LoadConfigFile_AutoAuth_AddrConformance verifies basic config file
// loading in addition to RFC-5942 §4 normalization of auto-auth methods.
// See: https://rfc-editor.org/rfc/rfc5952.html
func Test_LoadConfigFile_AutoAuth_AddrConformance(t *testing.T) {
	t.Parallel()

	for name, method := range map[string]*Method{
		"aws": {
			Type:      "aws",
			MountPath: "auth/aws",
			Namespace: "aws-namespace/",
			Config: map[string]any{
				"role": "foobar",
			},
		},
		"azure": {
			Type:      "azure",
			MountPath: "auth/azure",
			Namespace: "azure-namespace/",
			Config: map[string]any{
				"authenticate_from_environment": true,
				"role":                          "dev-role",
				"resource":                      "https://[2001:0:0:1::1]",
			},
		},
		"gcp": {
			Type:      "gcp",
			MountPath: "auth/gcp",
			Namespace: "gcp-namespace/",
			Config: map[string]any{
				"role":            "dev-role",
				"service_account": "https://[2001:db8:ac3:fe4::1]",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			config, err := LoadConfigFile("./test-fixtures/config-auto-auth-" + name + ".hcl")
			require.NoError(t, err)

			expected := &Config{
				SharedConfig: &configutil.SharedConfig{
					PidFile: "./pidfile",
					Listeners: []*configutil.Listener{
						{
							Type:        "unix",
							Address:     "/path/to/socket",
							TLSDisable:  true,
							SocketMode:  "configmode",
							SocketUser:  "configuser",
							SocketGroup: "configgroup",
						},
						{
							Type:       "tcp",
							Address:    "2001:db8::1:8200", // Normalized
							TLSDisable: true,
						},
						{
							Type:       "tcp",
							Address:    "[2001:0:0:1::1]:3000", // Normalized
							Role:       "metrics_only",
							TLSDisable: true,
						},
						{
							Type:        "tcp",
							Role:        "default",
							Address:     "2001:db8:0:1:1:1:1:1:8400", // Normalized
							TLSKeyFile:  "/path/to/cakey.pem",
							TLSCertFile: "/path/to/cacert.pem",
						},
					},
				},
				Vault: &Vault{
					Address:          "https://[2001:db8::1]:8200", // Normalized
					CACert:           "config_ca_cert",
					CAPath:           "config_ca_path",
					TLSSkipVerifyRaw: interface{}("true"),
					TLSSkipVerify:    true,
					ClientCert:       "config_client_cert",
					ClientKey:        "config_client_key",
					Retry: &Retry{
						NumRetries: 12,
					},
				},
				AutoAuth: &AutoAuth{
					Method: method, // Method properties are normalized correctly
					Sinks: []*Sink{
						{
							Type:   "file",
							DHType: "curve25519",
							DHPath: "/tmp/file-foo-dhpath",
							AAD:    "foobar",
							Config: map[string]interface{}{
								"path": "/tmp/file-foo",
							},
						},
					},
				},
				APIProxy: &APIProxy{
					EnforceConsistency:  "always",
					WhenInconsistent:    "retry",
					UseAutoAuthTokenRaw: true,
					UseAutoAuthToken:    true,
					ForceAutoAuthToken:  false,
				},
				Cache: &Cache{
					Persist: &agentproxyshared.PersistConfig{
						Type:                    "kubernetes",
						Path:                    "/vault/agent-cache/",
						KeepAfterImport:         true,
						ExitOnErr:               true,
						ServiceAccountTokenFile: "/tmp/serviceaccount/token",
					},
				},
			}

			config.Prune()
			require.EqualValues(t, expected, config)
		})
	}
}
