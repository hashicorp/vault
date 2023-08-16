package config

import (
	"os"
	"testing"
	"time"

	"github.com/go-test/deep"
	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
)

func TestLoadConfigFile_AgentCache(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-cache.hcl")
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
					Type:        "tcp",
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
		Cache: &Cache{
			UseAutoAuthToken:    true,
			UseAutoAuthTokenRaw: true,
			ForceAutoAuthToken:  false,
			Persist: &Persist{
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

	config, err = LoadConfig("./test-fixtures/config-cache-embedded-type.hcl")
	if err != nil {
		t.Fatal(err)
	}
	expected.Vault.TLSSkipVerifyRaw = interface{}(true)

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_AgentCache_NoListeners(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-cache-no-listeners.hcl")
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
		Cache: &Cache{
			UseAutoAuthToken:    true,
			UseAutoAuthTokenRaw: true,
			ForceAutoAuthToken:  false,
			Persist: &Persist{
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

	config, err := LoadConfig("./test-fixtures/config.hcl")
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
		Vault: &Vault{
			Retry: &Retry{
				NumRetries: 12,
			},
		},
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}

	config, err = LoadConfig("./test-fixtures/config-embedded-type.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestLoadConfigFile_Method_Wrapping(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-method-wrapping.hcl")
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
		Vault: &Vault{
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

func TestLoadConfigFile_Method_InitialBackoff(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-method-initial-backoff.hcl")
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
		Vault: &Vault{
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

func TestLoadConfigFile_Method_ExitOnErr(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-method-exit-on-err.hcl")
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
		Vault: &Vault{
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

func TestLoadConfigFile_AgentCache_NoAutoAuth(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-cache-no-auto_auth.hcl")
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
		Vault: &Vault{
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

func TestLoadConfigFile_Bad_AgentCache_InconsisentAutoAuth(t *testing.T) {
	_, err := LoadConfig("./test-fixtures/bad-config-cache-inconsistent-auto_auth.hcl")
	if err == nil {
		t.Fatal("LoadConfig should return an error when use_auto_auth_token=true and no auto_auth section present")
	}
}

func TestLoadConfigFile_Bad_AgentCache_ForceAutoAuthNoMethod(t *testing.T) {
	_, err := LoadConfig("./test-fixtures/bad-config-cache-inconsistent-auto_auth.hcl")
	if err == nil {
		t.Fatal("LoadConfig should return an error when use_auto_auth_token=true and no auto_auth section present")
	}
}

func TestLoadConfigFile_Bad_AgentCache_NoListeners(t *testing.T) {
	_, err := LoadConfig("./test-fixtures/bad-config-cache-no-listeners.hcl")
	if err == nil {
		t.Fatal("LoadConfig should return an error when cache section present and no listeners present and no templates defined")
	}
}

func TestLoadConfigFile_Bad_AutoAuth_Wrapped_Multiple_Sinks(t *testing.T) {
	_, err := LoadConfig("./test-fixtures/bad-config-auto_auth-wrapped-multiple-sinks")
	if err == nil {
		t.Fatal("LoadConfig should return an error when auth_auth.method.wrap_ttl nonzero and multiple sinks defined")
	}
}

func TestLoadConfigFile_Bad_AutoAuth_Nosinks_Nocache_Notemplates(t *testing.T) {
	_, err := LoadConfig("./test-fixtures/bad-config-auto_auth-nosinks-nocache-notemplates.hcl")
	if err == nil {
		t.Fatal("LoadConfig should return an error when auto_auth configured and there are no sinks, caches or templates")
	}
}

func TestLoadConfigFile_Bad_AutoAuth_Both_Wrapping_Types(t *testing.T) {
	_, err := LoadConfig("./test-fixtures/bad-config-method-wrapping-and-sink-wrapping.hcl")
	if err == nil {
		t.Fatal("LoadConfig should return an error when auth_auth.method.wrap_ttl nonzero and sinks.wrap_ttl nonzero")
	}
}

func TestLoadConfigFile_Bad_AgentCache_AutoAuth_Method_wrapping(t *testing.T) {
	_, err := LoadConfig("./test-fixtures/bad-config-cache-auto_auth-method-wrapping.hcl")
	if err == nil {
		t.Fatal("LoadConfig should return an error when auth_auth.method.wrap_ttl nonzero and cache.use_auto_auth_token=true")
	}
}

func TestLoadConfigFile_AgentCache_AutoAuth_NoSink(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-cache-auto_auth-no-sink.hcl")
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
		Cache: &Cache{
			UseAutoAuthToken:    true,
			UseAutoAuthTokenRaw: true,
			ForceAutoAuthToken:  false,
		},
		Vault: &Vault{
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

func TestLoadConfigFile_AgentCache_AutoAuth_Force(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-cache-auto_auth-force.hcl")
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
		Cache: &Cache{
			UseAutoAuthToken:    true,
			UseAutoAuthTokenRaw: "force",
			ForceAutoAuthToken:  true,
		},
		Vault: &Vault{
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

func TestLoadConfigFile_AgentCache_AutoAuth_True(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-cache-auto_auth-true.hcl")
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
		Cache: &Cache{
			UseAutoAuthToken:    true,
			UseAutoAuthTokenRaw: "true",
			ForceAutoAuthToken:  false,
		},
		Vault: &Vault{
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

func TestLoadConfigFile_AgentCache_AutoAuth_False(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-cache-auto_auth-false.hcl")
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
		Vault: &Vault{
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

func TestLoadConfigFile_AgentCache_Persist(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-cache-persist-false.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Cache: &Cache{
			Persist: &Persist{
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
		Vault: &Vault{
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

func TestLoadConfigFile_AgentCache_PersistMissingType(t *testing.T) {
	_, err := LoadConfig("./test-fixtures/config-cache-persist-empty-type.hcl")
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
			},
		},
		"empty": {
			"./test-fixtures/config-template_config-empty.hcl",
			TemplateConfig{
				ExitOnRetryFailure: false,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			config, err := LoadConfig(tc.fixturePath)
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
			config, err := LoadConfig(tc.fixturePath)
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
				Vault: &Vault{
					Retry: &Retry{
						NumRetries: 12,
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
			config, err := LoadConfig(tc.fixturePath)
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
				Vault: &Vault{
					Retry: &Retry{
						NumRetries: 12,
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

func TestLoadConfigFile_Vault_Retry(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-vault-retry.hcl")
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
	config, err := LoadConfig("./test-fixtures/config-vault-retry-empty.hcl")
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
	config, err := LoadConfig("./test-fixtures/config-consistency.hcl")
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
		Vault: &Vault{
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

func TestLoadConfigFile_Disable_Idle_Conns_All(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-disable-idle-connections-all.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableIdleConns:           []string{"auto-auth", "caching", "templating"},
		DisableIdleConnsCaching:    true,
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
	config, err := LoadConfig("./test-fixtures/config-disable-idle-connections-auto-auth.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableIdleConns:           []string{"auto-auth"},
		DisableIdleConnsCaching:    false,
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
	config, err := LoadConfig("./test-fixtures/config-disable-idle-connections-templating.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableIdleConns:           []string{"templating"},
		DisableIdleConnsCaching:    false,
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
	config, err := LoadConfig("./test-fixtures/config-disable-idle-connections-caching.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableIdleConns:           []string{"caching"},
		DisableIdleConnsCaching:    true,
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
	config, err := LoadConfig("./test-fixtures/config-disable-idle-connections-empty.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableIdleConns:           []string{},
		DisableIdleConnsCaching:    false,
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
	config, err := LoadConfig("./test-fixtures/config-disable-idle-connections-empty.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableIdleConns:           []string{"auto-auth", "caching", "templating"},
		DisableIdleConnsCaching:    true,
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
	_, err := LoadConfig("./test-fixtures/bad-config-disable-idle-connections.hcl")
	if err == nil {
		t.Fatal("should have error, it didn't")
	}
}

func TestLoadConfigFile_Disable_Keep_Alives_All(t *testing.T) {
	config, err := LoadConfig("./test-fixtures/config-disable-keep-alives-all.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableKeepAlives:           []string{"auto-auth", "caching", "templating"},
		DisableKeepAlivesCaching:    true,
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
	config, err := LoadConfig("./test-fixtures/config-disable-keep-alives-auto-auth.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableKeepAlives:           []string{"auto-auth"},
		DisableKeepAlivesCaching:    false,
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
	config, err := LoadConfig("./test-fixtures/config-disable-keep-alives-templating.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableKeepAlives:           []string{"templating"},
		DisableKeepAlivesCaching:    false,
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
	config, err := LoadConfig("./test-fixtures/config-disable-keep-alives-caching.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableKeepAlives:           []string{"caching"},
		DisableKeepAlivesCaching:    true,
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
	config, err := LoadConfig("./test-fixtures/config-disable-keep-alives-empty.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableKeepAlives:           []string{},
		DisableKeepAlivesCaching:    false,
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
	config, err := LoadConfig("./test-fixtures/config-disable-keep-alives-empty.hcl")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile: "./pidfile",
		},
		DisableKeepAlives:           []string{"auto-auth", "caching", "templating"},
		DisableKeepAlivesCaching:    true,
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
	_, err := LoadConfig("./test-fixtures/bad-config-disable-keep-alives.hcl")
	if err == nil {
		t.Fatal("should have error, it didn't")
	}
}
