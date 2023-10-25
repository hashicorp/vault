// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/command/agentproxyshared"
	"github.com/hashicorp/vault/internalshared/configutil"
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
