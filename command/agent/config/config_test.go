// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/go-test/deep"
	ctconfig "github.com/hashicorp/consul-template/config"
	"golang.org/x/exp/slices"

	"github.com/hashicorp/vault/command/agentproxyshared"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
)

func TestLoadConfigFile_AgentCache(t *testing.T) {
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
			UseAutoAuthToken:   true,
			ForceAutoAuthToken: false,
		},
		Cache: &Cache{
			UseAutoAuthToken:    true,
			UseAutoAuthTokenRaw: true,
			ForceAutoAuthToken:  false,
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

func TestLoadConfigDir_AgentCache(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-dir-cache/")
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
			UseAutoAuthToken:   true,
			ForceAutoAuthToken: false,
		},
		Cache: &Cache{
			UseAutoAuthToken:    true,
			UseAutoAuthTokenRaw: true,
			ForceAutoAuthToken:  false,
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

	config, err = LoadConfigFile("./test-fixtures/config-dir-cache/config-cache1.hcl")
	if err != nil {
		t.Fatal(err)
	}
	config2, err := LoadConfigFile("./test-fixtures/config-dir-cache/config-cache2.hcl")

	mergedConfig := config.Merge(config2)

	mergedConfig.Prune()
	if diff := deep.Equal(mergedConfig, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigDir_AutoAuthAndListener(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-dir-auto-auth-and-listener/")
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
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}

	config, err = LoadConfigFile("./test-fixtures/config-dir-auto-auth-and-listener/config1.hcl")
	if err != nil {
		t.Fatal(err)
	}
	config2, err := LoadConfigFile("./test-fixtures/config-dir-auto-auth-and-listener/config2.hcl")

	mergedConfig := config.Merge(config2)

	mergedConfig.Prune()
	if diff := deep.Equal(mergedConfig, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigDir_VaultBlock(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-dir-vault-block/")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
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
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}

	config, err = LoadConfigFile("./test-fixtures/config-dir-vault-block/config1.hcl")
	if err != nil {
		t.Fatal(err)
	}
	config2, err := LoadConfigFile("./test-fixtures/config-dir-vault-block/config2.hcl")

	mergedConfig := config.Merge(config2)

	mergedConfig.Prune()
	if diff := deep.Equal(mergedConfig, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_AgentCache_NoListeners(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-cache-no-listeners.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
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
			UseAutoAuthToken:   true,
			ForceAutoAuthToken: false,
		},
		Cache: &Cache{
			UseAutoAuthToken:    true,
			UseAutoAuthTokenRaw: true,
			ForceAutoAuthToken:  false,
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
		Templates: []*ctconfig.TemplateConfig{
			{
				Source:      pointerutil.StringPtr("/path/on/disk/to/template.ctmpl"),
				Destination: pointerutil.StringPtr("/path/on/disk/where/template/will/render.txt"),
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile(t *testing.T) {
	if err := os.Setenv("TEST_AAD_ENV", "aad"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Unsetenv("TEST_AAD_ENV"); err != nil {
			t.Fatal(err)
		}
	}()

	config, err := LoadConfigFile("./test-fixtures/config.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
			LogFile: "/var/log/vault/vault-agent.log",
		},
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
				Config: map[string]interface{}{
					"role": "foobar",
				},
				MaxBackoff: 0,
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
				{
					Type:      "file",
					WrapTTL:   5 * time.Minute,
					DHType:    "curve25519",
					DHPath:    "/tmp/file-foo-dhpath2",
					AAD:       "aad",
					DeriveKey: true,
					Config: map[string]interface{}{
						"path": "/tmp/file-bar",
					},
				},
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}

	config, err = LoadConfigFile("./test-fixtures/config-embedded-type.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Method_Wrapping(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-method-wrapping.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:        "aws",
				MountPath:   "auth/aws",
				ExitOnError: false,
				WrapTTL:     5 * time.Minute,
				MaxBackoff:  2 * time.Minute,
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
			Sinks: []*Sink{
				{
					Type: "file",
					Config: map[string]interface{}{
						"path": "/tmp/file-foo",
					},
				},
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Method_InitialBackoff(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-method-initial-backoff.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:        "aws",
				MountPath:   "auth/aws",
				ExitOnError: false,
				WrapTTL:     5 * time.Minute,
				MinBackoff:  5 * time.Second,
				MaxBackoff:  2 * time.Minute,
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
			Sinks: []*Sink{
				{
					Type: "file",
					Config: map[string]interface{}{
						"path": "/tmp/file-foo",
					},
				},
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Method_ExitOnErr(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-method-exit-on-err.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:        "aws",
				MountPath:   "auth/aws",
				ExitOnError: true,
				WrapTTL:     5 * time.Minute,
				MinBackoff:  5 * time.Second,
				MaxBackoff:  2 * time.Minute,
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
			Sinks: []*Sink{
				{
					Type: "file",
					Config: map[string]interface{}{
						"path": "/tmp/file-foo",
					},
				},
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_AgentCache_NoAutoAuth(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-cache-no-auto_auth.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Cache: &Cache{},
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
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Bad_AgentCache_InconsisentAutoAuth(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/bad-config-cache-inconsistent-auto_auth.hcl")
	if err != nil {
		t.Fatalf("LoadConfigFile should not return an error for this config, err: %v", err)
	}
	if config == nil {
		t.Fatal("config was nil")
	}
	err = config.ValidateConfig()
	if err == nil {
		t.Fatal("ValidateConfig should return an error when use_auto_auth_token=true and no auto_auth section present")
	}
}

func TestLoadConfigFile_Bad_AgentCache_ForceAutoAuthNoMethod(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/bad-config-cache-force-token-no-auth-method.hcl")
	if err != nil {
		t.Fatalf("LoadConfigFile should not return an error for this config, err: %v", err)
	}
	if config == nil {
		t.Fatal("config was nil")
	}
	err = config.ValidateConfig()
	if err == nil {
		t.Fatal("ValidateConfig should return an error when use_auto_auth_token=force and no auto_auth section present")
	}
}

func TestLoadConfigFile_Bad_AgentCache_NoListeners(t *testing.T) {
	_, err := LoadConfigFile("./test-fixtures/bad-config-cache-no-listeners.hcl")
	if err != nil {
		t.Fatalf("LoadConfigFile should return an error for this config")
	}
}

func TestLoadConfigFile_Bad_AutoAuth_Wrapped_Multiple_Sinks(t *testing.T) {
	_, err := LoadConfigFile("./test-fixtures/bad-config-auto_auth-wrapped-multiple-sinks.hcl")
	if err == nil {
		t.Fatalf("LoadConfigFile should return an error for this config, err: %v", err)
	}
}

func TestLoadConfigFile_Bad_AutoAuth_Nosinks_Nocache_Notemplates(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/bad-config-auto_auth-nosinks-nocache-notemplates.hcl")
	if err != nil {
		t.Fatalf("LoadConfigFile should not return an error for this config, err: %v", err)
	}
	if config == nil {
		t.Fatal("config was nil")
	}
	err = config.ValidateConfig()
	if err == nil {
		t.Fatal("ValidateConfig should return an error when auto_auth configured and there are no sinks, caches or templates")
	}
}

func TestLoadConfigFile_Bad_AutoAuth_Both_Wrapping_Types(t *testing.T) {
	_, err := LoadConfigFile("./test-fixtures/bad-config-method-wrapping-and-sink-wrapping.hcl")
	if err == nil {
		t.Fatalf("LoadConfigFile should return an error for this config")
	}
}

func TestLoadConfigFile_Bad_AgentCache_AutoAuth_Method_wrapping(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/bad-config-cache-auto_auth-method-wrapping.hcl")
	if err != nil {
		t.Fatalf("LoadConfigFile should not return an error for this config, err: %v", err)
	}
	if config == nil {
		t.Fatal("config was nil")
	}
	err = config.ValidateConfig()
	if err == nil {
		t.Fatal("ValidateConfig should return an error when auth_auth.method.wrap_ttl nonzero and cache.use_auto_auth_token=true")
	}
}

func TestLoadConfigFile_Bad_APIProxy_And_Cache_Same_Config(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/bad-config-api_proxy-cache.hcl")
	if err != nil {
		t.Fatalf("LoadConfigFile should not return an error for this config, err: %v", err)
	}
	if config == nil {
		t.Fatal("config was nil")
	}
	err = config.ValidateConfig()
	if err == nil {
		t.Fatal("ValidateConfig should return an error when cache and api_proxy try and configure the same value")
	}
}

func TestLoadConfigFile_AgentCache_AutoAuth_NoSink(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-cache-auto_auth-no-sink.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:       "tcp",
					Address:    "127.0.0.1:8300",
					TLSDisable: true,
				},
			},
			PidFile: "./pidfile",
		},
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
		},
		APIProxy: &APIProxy{
			UseAutoAuthToken:   true,
			ForceAutoAuthToken: false,
		},
		Cache: &Cache{
			UseAutoAuthToken:    true,
			UseAutoAuthTokenRaw: true,
			ForceAutoAuthToken:  false,
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_AgentCache_AutoAuth_Force(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-cache-auto_auth-force.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:       "tcp",
					Address:    "127.0.0.1:8300",
					TLSDisable: true,
				},
			},
			PidFile: "./pidfile",
		},
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
		},
		APIProxy: &APIProxy{
			UseAutoAuthToken:   true,
			ForceAutoAuthToken: true,
		},
		Cache: &Cache{
			UseAutoAuthToken:    true,
			UseAutoAuthTokenRaw: "force",
			ForceAutoAuthToken:  true,
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_AgentCache_AutoAuth_True(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-cache-auto_auth-true.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:       "tcp",
					Address:    "127.0.0.1:8300",
					TLSDisable: true,
				},
			},
			PidFile: "./pidfile",
		},
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
		},
		APIProxy: &APIProxy{
			UseAutoAuthToken:   true,
			ForceAutoAuthToken: false,
		},
		Cache: &Cache{
			UseAutoAuthToken:    true,
			UseAutoAuthTokenRaw: "true",
			ForceAutoAuthToken:  false,
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Agent_AutoAuth_APIProxyAllConfig(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-api_proxy-auto_auth-all-api_proxy-config.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:       "tcp",
					Address:    "127.0.0.1:8300",
					TLSDisable: true,
				},
			},
			PidFile: "./pidfile",
		},
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
		},
		APIProxy: &APIProxy{
			UseAutoAuthToken:    true,
			UseAutoAuthTokenRaw: "force",
			ForceAutoAuthToken:  true,
			EnforceConsistency:  "always",
			WhenInconsistent:    "forward",
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_AgentCache_AutoAuth_False(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-cache-auto_auth-false.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:       "tcp",
					Address:    "127.0.0.1:8300",
					TLSDisable: true,
				},
			},
			PidFile: "./pidfile",
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
			UseAutoAuthToken:    false,
			UseAutoAuthTokenRaw: "false",
			ForceAutoAuthToken:  false,
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_AgentCache_Persist(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-cache-persist-false.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Cache: &Cache{
			Persist: &agentproxyshared.PersistConfig{
				Type:                    "kubernetes",
				Path:                    "/vault/agent-cache/",
				KeepAfterImport:         false,
				ExitOnErr:               false,
				ServiceAccountTokenFile: "",
			},
		},
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
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_AgentCache_PersistMissingType(t *testing.T) {
	_, err := LoadConfigFile("./test-fixtures/config-cache-persist-empty-type.hcl")
	if err == nil || os.IsNotExist(err) {
		t.Fatal("expected error or file is missing")
	}
}

func TestLoadConfigFile_TemplateConfig(t *testing.T) {
	testCases := map[string]struct {
		fixturePath            string
		expectedTemplateConfig TemplateConfig
	}{
		"set-true": {
			"./test-fixtures/config-template_config.hcl",
			TemplateConfig{
				ExitOnRetryFailure:    true,
				StaticSecretRenderInt: 1 * time.Minute,
				MaxConnectionsPerHost: 100,
			},
		},
		"empty": {
			"./test-fixtures/config-template_config-empty.hcl",
			TemplateConfig{
				ExitOnRetryFailure:    false,
				MaxConnectionsPerHost: 10,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			config, err := LoadConfigFile(tc.fixturePath)
			if err != nil {
				t.Fatal(err)
			}

			expected := &Config{
				SharedConfig: &configutil.SharedConfig{},
				Vault: &Vault{
					Address: "http://127.0.0.1:1111",
					Retry: &Retry{
						NumRetries: 5,
					},
				},
				TemplateConfig: &tc.expectedTemplateConfig,
				Templates: []*ctconfig.TemplateConfig{
					{
						Source:      pointerutil.StringPtr("/path/on/disk/to/template.ctmpl"),
						Destination: pointerutil.StringPtr("/path/on/disk/where/template/will/render.txt"),
					},
				},
			}

			config.Prune()
			if diff := deep.Equal(config, expected); diff != nil {
				t.Fatal(diff)
			}
		})
	}
}

// TestLoadConfigFile_Template tests template definitions in Vault Agent
func TestLoadConfigFile_Template(t *testing.T) {
	testCases := map[string]struct {
		fixturePath       string
		expectedTemplates []*ctconfig.TemplateConfig
	}{
		"min": {
			fixturePath: "./test-fixtures/config-template-min.hcl",
			expectedTemplates: []*ctconfig.TemplateConfig{
				{
					Source:      pointerutil.StringPtr("/path/on/disk/to/template.ctmpl"),
					Destination: pointerutil.StringPtr("/path/on/disk/where/template/will/render.txt"),
				},
			},
		},
		"full": {
			fixturePath: "./test-fixtures/config-template-full.hcl",
			expectedTemplates: []*ctconfig.TemplateConfig{
				{
					Backup:         pointerutil.BoolPtr(true),
					Command:        []string{"restart service foo"},
					CommandTimeout: pointerutil.TimeDurationPtr("60s"),
					Contents:       pointerutil.StringPtr("{{ keyOrDefault \"service/redis/maxconns@east-aws\" \"5\" }}"),
					CreateDestDirs: pointerutil.BoolPtr(true),
					Destination:    pointerutil.StringPtr("/path/on/disk/where/template/will/render.txt"),
					ErrMissingKey:  pointerutil.BoolPtr(true),
					LeftDelim:      pointerutil.StringPtr("<<"),
					Perms:          pointerutil.FileModePtr(0o655),
					RightDelim:     pointerutil.StringPtr(">>"),
					SandboxPath:    pointerutil.StringPtr("/path/on/disk/where"),
					Exec: &ctconfig.ExecConfig{
						Command: []string{"foo"},
						Timeout: pointerutil.TimeDurationPtr("10s"),
					},

					Wait: &ctconfig.WaitConfig{
						Min: pointerutil.TimeDurationPtr("10s"),
						Max: pointerutil.TimeDurationPtr("40s"),
					},
				},
			},
		},
		"many": {
			fixturePath: "./test-fixtures/config-template-many.hcl",
			expectedTemplates: []*ctconfig.TemplateConfig{
				{
					Source:         pointerutil.StringPtr("/path/on/disk/to/template.ctmpl"),
					Destination:    pointerutil.StringPtr("/path/on/disk/where/template/will/render.txt"),
					ErrMissingKey:  pointerutil.BoolPtr(false),
					CreateDestDirs: pointerutil.BoolPtr(true),
					Command:        []string{"restart service foo"},
					Perms:          pointerutil.FileModePtr(0o600),
				},
				{
					Source:      pointerutil.StringPtr("/path/on/disk/to/template2.ctmpl"),
					Destination: pointerutil.StringPtr("/path/on/disk/where/template/will/render2.txt"),
					Backup:      pointerutil.BoolPtr(true),
					Perms:       pointerutil.FileModePtr(0o755),
					Wait: &ctconfig.WaitConfig{
						Min: pointerutil.TimeDurationPtr("2s"),
						Max: pointerutil.TimeDurationPtr("10s"),
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			config, err := LoadConfigFile(tc.fixturePath)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			expected := &Config{
				SharedConfig: &configutil.SharedConfig{
					PidFile: "./pidfile",
				},
				AutoAuth: &AutoAuth{
					Method: &Method{
						Type:      "aws",
						MountPath: "auth/aws",
						Namespace: "my-namespace/",
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
				Templates: tc.expectedTemplates,
			}

			config.Prune()
			if diff := deep.Equal(config, expected); diff != nil {
				t.Fatal(diff)
			}
		})
	}
}

// TestLoadConfigFile_Template_NoSinks tests template definitions without sinks in Vault Agent
func TestLoadConfigFile_Template_NoSinks(t *testing.T) {
	testCases := map[string]struct {
		fixturePath       string
		expectedTemplates []*ctconfig.TemplateConfig
	}{
		"min": {
			fixturePath: "./test-fixtures/config-template-min-nosink.hcl",
			expectedTemplates: []*ctconfig.TemplateConfig{
				{
					Source:      pointerutil.StringPtr("/path/on/disk/to/template.ctmpl"),
					Destination: pointerutil.StringPtr("/path/on/disk/where/template/will/render.txt"),
				},
			},
		},
		"full": {
			fixturePath: "./test-fixtures/config-template-full-nosink.hcl",
			expectedTemplates: []*ctconfig.TemplateConfig{
				{
					Backup:         pointerutil.BoolPtr(true),
					Command:        []string{"restart service foo"},
					CommandTimeout: pointerutil.TimeDurationPtr("60s"),
					Contents:       pointerutil.StringPtr("{{ keyOrDefault \"service/redis/maxconns@east-aws\" \"5\" }}"),
					CreateDestDirs: pointerutil.BoolPtr(true),
					Destination:    pointerutil.StringPtr("/path/on/disk/where/template/will/render.txt"),
					ErrMissingKey:  pointerutil.BoolPtr(true),
					LeftDelim:      pointerutil.StringPtr("<<"),
					Perms:          pointerutil.FileModePtr(0o655),
					RightDelim:     pointerutil.StringPtr(">>"),
					SandboxPath:    pointerutil.StringPtr("/path/on/disk/where"),

					Wait: &ctconfig.WaitConfig{
						Min: pointerutil.TimeDurationPtr("10s"),
						Max: pointerutil.TimeDurationPtr("40s"),
					},
				},
			},
		},
		"many": {
			fixturePath: "./test-fixtures/config-template-many-nosink.hcl",
			expectedTemplates: []*ctconfig.TemplateConfig{
				{
					Source:         pointerutil.StringPtr("/path/on/disk/to/template.ctmpl"),
					Destination:    pointerutil.StringPtr("/path/on/disk/where/template/will/render.txt"),
					ErrMissingKey:  pointerutil.BoolPtr(false),
					CreateDestDirs: pointerutil.BoolPtr(true),
					Command:        []string{"restart service foo"},
					Perms:          pointerutil.FileModePtr(0o600),
				},
				{
					Source:      pointerutil.StringPtr("/path/on/disk/to/template2.ctmpl"),
					Destination: pointerutil.StringPtr("/path/on/disk/where/template/will/render2.txt"),
					Backup:      pointerutil.BoolPtr(true),
					Perms:       pointerutil.FileModePtr(0o755),
					Wait: &ctconfig.WaitConfig{
						Min: pointerutil.TimeDurationPtr("2s"),
						Max: pointerutil.TimeDurationPtr("10s"),
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			config, err := LoadConfigFile(tc.fixturePath)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			expected := &Config{
				SharedConfig: &configutil.SharedConfig{
					PidFile: "./pidfile",
				},
				AutoAuth: &AutoAuth{
					Method: &Method{
						Type:      "aws",
						MountPath: "auth/aws",
						Namespace: "my-namespace/",
						Config: map[string]interface{}{
							"role": "foobar",
						},
					},
					Sinks: nil,
				},
				Templates: tc.expectedTemplates,
			}

			config.Prune()
			if diff := deep.Equal(config, expected); diff != nil {
				t.Fatal(diff)
			}
		})
	}
}

// TestLoadConfigFile_Template_WithCache tests ensures that cache {} stanza is
// permitted in vault agent configuration with template(s)
func TestLoadConfigFile_Template_WithCache(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-template-with-cache.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
		},
		Cache: &Cache{},
		Templates: []*ctconfig.TemplateConfig{
			{
				Source:      pointerutil.StringPtr("/path/on/disk/to/template.ctmpl"),
				Destination: pointerutil.StringPtr("/path/on/disk/where/template/will/render.txt"),
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Vault_Retry(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-vault-retry.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				NumRetries: 5,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Vault_Retry_Empty(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-vault-retry-empty.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_EnforceConsistency(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-consistency.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:       "tcp",
					Address:    "127.0.0.1:8300",
					TLSDisable: true,
				},
			},
			PidFile: "",
		},
		Cache: &Cache{
			EnforceConsistency: "always",
			WhenInconsistent:   "retry",
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_EnforceConsistency_APIProxy(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-consistency-apiproxy.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:       "tcp",
					Address:    "127.0.0.1:8300",
					TLSDisable: true,
				},
			},
			PidFile: "",
		},
		APIProxy: &APIProxy{
			EnforceConsistency: "always",
			WhenInconsistent:   "retry",
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Disable_Idle_Conns_All(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-disable-idle-connections-all.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableIdleConns:           []string{"auto-auth", "caching", "templating", "proxying"},
		DisableIdleConnsAPIProxy:   true,
		DisableIdleConnsAutoAuth:   true,
		DisableIdleConnsTemplating: true,
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Disable_Idle_Conns_Auto_Auth(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-disable-idle-connections-auto-auth.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableIdleConns:           []string{"auto-auth"},
		DisableIdleConnsAPIProxy:   false,
		DisableIdleConnsAutoAuth:   true,
		DisableIdleConnsTemplating: false,
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Disable_Idle_Conns_Templating(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-disable-idle-connections-templating.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableIdleConns:           []string{"templating"},
		DisableIdleConnsAPIProxy:   false,
		DisableIdleConnsAutoAuth:   false,
		DisableIdleConnsTemplating: true,
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Disable_Idle_Conns_Caching(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-disable-idle-connections-caching.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableIdleConns:           []string{"caching"},
		DisableIdleConnsAPIProxy:   true,
		DisableIdleConnsAutoAuth:   false,
		DisableIdleConnsTemplating: false,
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Disable_Idle_Conns_Proxying(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-disable-idle-connections-proxying.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableIdleConns:           []string{"proxying"},
		DisableIdleConnsAPIProxy:   true,
		DisableIdleConnsAutoAuth:   false,
		DisableIdleConnsTemplating: false,
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Disable_Idle_Conns_Empty(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-disable-idle-connections-empty.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableIdleConns:           []string{},
		DisableIdleConnsAPIProxy:   false,
		DisableIdleConnsAutoAuth:   false,
		DisableIdleConnsTemplating: false,
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Disable_Idle_Conns_Env(t *testing.T) {
	err := os.Setenv(DisableIdleConnsEnv, "auto-auth,caching,templating")
	defer os.Unsetenv(DisableIdleConnsEnv)

	if err != nil {
		t.Fatal(err)
	}
	config, err := LoadConfigFile("./test-fixtures/config-disable-idle-connections-empty.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableIdleConns:           []string{"auto-auth", "caching", "templating"},
		DisableIdleConnsAPIProxy:   true,
		DisableIdleConnsAutoAuth:   true,
		DisableIdleConnsTemplating: true,
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Bad_Value_Disable_Idle_Conns(t *testing.T) {
	_, err := LoadConfigFile("./test-fixtures/bad-config-disable-idle-connections.hcl")
	if err == nil {
		t.Fatal("should have error, it didn't")
	}
}

func TestLoadConfigFile_Disable_Keep_Alives_All(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-disable-keep-alives-all.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableKeepAlives:           []string{"auto-auth", "caching", "templating", "proxying"},
		DisableKeepAlivesAPIProxy:   true,
		DisableKeepAlivesAutoAuth:   true,
		DisableKeepAlivesTemplating: true,
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Disable_Keep_Alives_Auto_Auth(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-disable-keep-alives-auto-auth.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableKeepAlives:           []string{"auto-auth"},
		DisableKeepAlivesAPIProxy:   false,
		DisableKeepAlivesAutoAuth:   true,
		DisableKeepAlivesTemplating: false,
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Disable_Keep_Alives_Templating(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-disable-keep-alives-templating.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableKeepAlives:           []string{"templating"},
		DisableKeepAlivesAPIProxy:   false,
		DisableKeepAlivesAutoAuth:   false,
		DisableKeepAlivesTemplating: true,
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Disable_Keep_Alives_Caching(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-disable-keep-alives-caching.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableKeepAlives:           []string{"caching"},
		DisableKeepAlivesAPIProxy:   true,
		DisableKeepAlivesAutoAuth:   false,
		DisableKeepAlivesTemplating: false,
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Disable_Keep_Alives_Proxying(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-disable-keep-alives-proxying.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableKeepAlives:           []string{"proxying"},
		DisableKeepAlivesAPIProxy:   true,
		DisableKeepAlivesAutoAuth:   false,
		DisableKeepAlivesTemplating: false,
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Disable_Keep_Alives_Empty(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config-disable-keep-alives-empty.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableKeepAlives:           []string{},
		DisableKeepAlivesAPIProxy:   false,
		DisableKeepAlivesAutoAuth:   false,
		DisableKeepAlivesTemplating: false,
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Disable_Keep_Alives_Env(t *testing.T) {
	err := os.Setenv(DisableKeepAlivesEnv, "auto-auth,caching,templating")
	defer os.Unsetenv(DisableKeepAlivesEnv)

	if err != nil {
		t.Fatal(err)
	}
	config, err := LoadConfigFile("./test-fixtures/config-disable-keep-alives-empty.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableKeepAlives:           []string{"auto-auth", "caching", "templating"},
		DisableKeepAlivesAPIProxy:   true,
		DisableKeepAlivesAutoAuth:   true,
		DisableKeepAlivesTemplating: true,
		AutoAuth: &AutoAuth{
			Method: &Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Namespace: "my-namespace/",
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
		Vault: &Vault{
			Address: "http://127.0.0.1:1111",
			Retry: &Retry{
				ctconfig.DefaultRetryAttempts,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Bad_Value_Disable_Keep_Alives(t *testing.T) {
	_, err := LoadConfigFile("./test-fixtures/bad-config-disable-keep-alives.hcl")
	if err == nil {
		t.Fatal("should have error, it didn't")
	}
}

// TestLoadConfigFile_EnvTemplates_Simple loads and validates an env_template config
func TestLoadConfigFile_EnvTemplates_Simple(t *testing.T) {
	cfg, err := LoadConfigFile("./test-fixtures/config-env-templates-simple.hcl")
	if err != nil {
		t.Fatalf("error loading config file: %s", err)
	}

	if err := cfg.ValidateConfig(); err != nil {
		t.Fatalf("validation error: %s", err)
	}

	expectedKey := "MY_DATABASE_USER"
	found := false
	for _, envTemplate := range cfg.EnvTemplates {
		if *envTemplate.MapToEnvironmentVariable == expectedKey {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected environment variable name to be populated")
	}
}

// TestLoadConfigFile_EnvTemplates_Complex loads and validates an env_template config
func TestLoadConfigFile_EnvTemplates_Complex(t *testing.T) {
	cfg, err := LoadConfigFile("./test-fixtures/config-env-templates-complex.hcl")
	if err != nil {
		t.Fatalf("error loading config file: %s", err)
	}

	if err := cfg.ValidateConfig(); err != nil {
		t.Fatalf("validation error: %s", err)
	}

	expectedKeys := []string{
		"FOO_PASSWORD",
		"FOO_USER",
	}

	envExists := func(key string) bool {
		for _, envTmpl := range cfg.EnvTemplates {
			if *envTmpl.MapToEnvironmentVariable == key {
				return true
			}
		}
		return false
	}

	for _, expected := range expectedKeys {
		if !envExists(expected) {
			t.Fatalf("expected environment variable %s", expected)
		}
	}
}

// TestLoadConfigFile_EnvTemplates_WithSource loads and validates an
// env_template config with "source" instead of "contents"
func TestLoadConfigFile_EnvTemplates_WithSource(t *testing.T) {
	cfg, err := LoadConfigFile("./test-fixtures/config-env-templates-with-source.hcl")
	if err != nil {
		t.Fatalf("error loading config file: %s", err)
	}

	if err := cfg.ValidateConfig(); err != nil {
		t.Fatalf("validation error: %s", err)
	}
}

// TestLoadConfigFile_EnvTemplates_NoName ensures that env_template with no name triggers an error
func TestLoadConfigFile_EnvTemplates_NoName(t *testing.T) {
	_, err := LoadConfigFile("./test-fixtures/bad-config-env-templates-no-name.hcl")
	if err == nil {
		t.Fatalf("expected error")
	}
}

// TestLoadConfigFile_EnvTemplates_ExecInvalidSignal ensures that an invalid signal triggers an error
func TestLoadConfigFile_EnvTemplates_ExecInvalidSignal(t *testing.T) {
	_, err := LoadConfigFile("./test-fixtures/bad-config-env-templates-invalid-signal.hcl")
	if err == nil {
		t.Fatalf("expected error")
	}
}

// TestLoadConfigFile_EnvTemplates_ExecSimple validates the exec section with default parameters
func TestLoadConfigFile_EnvTemplates_ExecSimple(t *testing.T) {
	cfg, err := LoadConfigFile("./test-fixtures/config-env-templates-simple.hcl")
	if err != nil {
		t.Fatalf("error loading config file: %s", err)
	}

	if err := cfg.ValidateConfig(); err != nil {
		t.Fatalf("validation error: %s", err)
	}

	expectedCmd := []string{"/path/to/my/app", "arg1", "arg2"}
	if !slices.Equal(cfg.Exec.Command, expectedCmd) {
		t.Fatal("exec.command does not have expected value")
	}

	// check defaults
	if cfg.Exec.RestartOnSecretChanges != "always" {
		t.Fatalf("expected cfg.Exec.RestartOnSecretChanges to be 'always', got '%s'", cfg.Exec.RestartOnSecretChanges)
	}

	if cfg.Exec.RestartStopSignal != syscall.SIGTERM {
		t.Fatalf("expected cfg.Exec.RestartStopSignal to be 'syscall.SIGTERM', got '%s'", cfg.Exec.RestartStopSignal)
	}
}

// TestLoadConfigFile_EnvTemplates_ExecComplex validates the exec section with non-default parameters
func TestLoadConfigFile_EnvTemplates_ExecComplex(t *testing.T) {
	cfg, err := LoadConfigFile("./test-fixtures/config-env-templates-complex.hcl")
	if err != nil {
		t.Fatalf("error loading config file: %s", err)
	}

	if err := cfg.ValidateConfig(); err != nil {
		t.Fatalf("validation error: %s", err)
	}

	if !slices.Equal(cfg.Exec.Command, []string{"env"}) {
		t.Fatal("exec.command does not have expected value")
	}

	if cfg.Exec.RestartOnSecretChanges != "never" {
		t.Fatalf("expected cfg.Exec.RestartOnSecretChanges to be 'never', got %q", cfg.Exec.RestartOnSecretChanges)
	}

	if cfg.Exec.RestartStopSignal != syscall.SIGINT {
		t.Fatalf("expected cfg.Exec.RestartStopSignal to be 'syscall.SIGINT', got %q", cfg.Exec.RestartStopSignal)
	}
}

// TestLoadConfigFile_Bad_EnvTemplates_MissingExec ensures that ValidateConfig
// errors when "env_template" stanza(s) are specified but "exec" is missing
func TestLoadConfigFile_Bad_EnvTemplates_MissingExec(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/bad-config-env-templates-missing-exec.hcl")
	if err != nil {
		t.Fatalf("error loading config file: %s", err)
	}

	if err := config.ValidateConfig(); err == nil {
		t.Fatal("expected an error from ValidateConfig: exec section is missing")
	}
}

// TestLoadConfigFile_Bad_EnvTemplates_WithProxy ensures that ValidateConfig
// errors when both env_template and api_proxy stanzas are present
func TestLoadConfigFile_Bad_EnvTemplates_WithProxy(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/bad-config-env-templates-with-proxy.hcl")
	if err != nil {
		t.Fatalf("error loading config file: %s", err)
	}

	if err := config.ValidateConfig(); err == nil {
		t.Fatal("expected an error from ValidateConfig: listener / api_proxy are not compatible with env_template")
	}
}

// TestLoadConfigFile_Bad_EnvTemplates_WithFileTemplates ensures that
// ValidateConfig errors when both env_template and template stanzas are present
func TestLoadConfigFile_Bad_EnvTemplates_WithFileTemplates(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/bad-config-env-templates-with-file-templates.hcl")
	if err != nil {
		t.Fatalf("error loading config file: %s", err)
	}

	if err := config.ValidateConfig(); err == nil {
		t.Fatal("expected an error from ValidateConfig: file template stanza is not compatible with env_template")
	}
}

// TestLoadConfigFile_Bad_EnvTemplates_DisalowedFields ensure that
// ValidateConfig errors for disalowed env_template fields
func TestLoadConfigFile_Bad_EnvTemplates_DisalowedFields(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/bad-config-env-templates-disalowed-fields.hcl")
	if err != nil {
		t.Fatalf("error loading config file: %s", err)
	}

	if err := config.ValidateConfig(); err == nil {
		t.Fatal("expected an error from ValidateConfig: disallowed fields specified in env_template")
	}
}
