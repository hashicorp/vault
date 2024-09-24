// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginutil

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-secure-stdlib/plugincontainer"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginruntimeutil"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMakeConfig(t *testing.T) {
	type testCase struct {
		rc runConfig

		responseWrapInfo      *wrapping.ResponseWrapInfo
		responseWrapInfoErr   error
		responseWrapInfoTimes int

		mlockEnabled      bool
		mlockEnabledTimes int

		expectedConfig       *plugin.ClientConfig
		expectTLSConfig      bool
		expectRunnerFunc     bool
		skipSecureConfig     bool
		useLegacyEnvLayering bool
	}

	tests := map[string]testCase{
		"metadata mode, not-AutoMTLS": {
			rc: runConfig{
				command: "echo",
				args:    []string{"foo", "bar"},
				sha256:  []byte("some_sha256"),
				env:     []string{"initial=true"},
				PluginClientConfig: PluginClientConfig{
					PluginSets: map[int]plugin.PluginSet{
						1: {
							"bogus": nil,
						},
					},
					HandshakeConfig: plugin.HandshakeConfig{
						ProtocolVersion:  1,
						MagicCookieKey:   "magic_cookie_key",
						MagicCookieValue: "magic_cookie_value",
					},
					Logger:         hclog.NewNullLogger(),
					IsMetadataMode: true,
					AutoMTLS:       false,
				},
			},

			responseWrapInfoTimes: 0,

			mlockEnabled:         false,
			mlockEnabledTimes:    1,
			useLegacyEnvLayering: true,

			expectedConfig: &plugin.ClientConfig{
				HandshakeConfig: plugin.HandshakeConfig{
					ProtocolVersion:  1,
					MagicCookieKey:   "magic_cookie_key",
					MagicCookieValue: "magic_cookie_value",
				},
				VersionedPlugins: map[int]plugin.PluginSet{
					1: {
						"bogus": nil,
					},
				},
				Cmd: commandWithEnv(
					"echo",
					[]string{"foo", "bar"},
					append(append([]string{
						"initial=true",
						fmt.Sprintf("%s=%s", PluginVaultVersionEnv, "dummyversion"),
						fmt.Sprintf("%s=%t", PluginMetadataModeEnv, true),
						fmt.Sprintf("%s=%t", PluginAutoMTLSEnv, false),
					}, os.Environ()...), PluginUseLegacyEnvLayering+"=true"),
				),
				SecureConfig: &plugin.SecureConfig{
					Checksum: []byte("some_sha256"),
					// Hash is generated
				},
				AllowedProtocols: []plugin.Protocol{
					plugin.ProtocolNetRPC,
					plugin.ProtocolGRPC,
				},
				Logger:      hclog.NewNullLogger(),
				AutoMTLS:    false,
				SkipHostEnv: true,
			},
			expectTLSConfig: false,
		},
		"non-metadata mode, not-AutoMTLS": {
			rc: runConfig{
				command: "echo",
				args:    []string{"foo", "bar"},
				sha256:  []byte("some_sha256"),
				env:     []string{"initial=true"},
				PluginClientConfig: PluginClientConfig{
					PluginSets: map[int]plugin.PluginSet{
						1: {
							"bogus": nil,
						},
					},
					HandshakeConfig: plugin.HandshakeConfig{
						ProtocolVersion:  1,
						MagicCookieKey:   "magic_cookie_key",
						MagicCookieValue: "magic_cookie_value",
					},
					Logger:         hclog.NewNullLogger(),
					IsMetadataMode: false,
					AutoMTLS:       false,
				},
			},

			responseWrapInfo: &wrapping.ResponseWrapInfo{
				Token: "testtoken",
			},
			responseWrapInfoTimes: 1,

			mlockEnabled:      true,
			mlockEnabledTimes: 1,

			expectedConfig: &plugin.ClientConfig{
				HandshakeConfig: plugin.HandshakeConfig{
					ProtocolVersion:  1,
					MagicCookieKey:   "magic_cookie_key",
					MagicCookieValue: "magic_cookie_value",
				},
				VersionedPlugins: map[int]plugin.PluginSet{
					1: {
						"bogus": nil,
					},
				},
				Cmd: commandWithEnv(
					"echo",
					[]string{"foo", "bar"},
					append(os.Environ(), []string{
						"initial=true",
						fmt.Sprintf("%s=%t", PluginMlockEnabled, true),
						fmt.Sprintf("%s=%s", PluginVaultVersionEnv, "dummyversion"),
						fmt.Sprintf("%s=%t", PluginMetadataModeEnv, false),
						fmt.Sprintf("%s=%t", PluginAutoMTLSEnv, false),
						fmt.Sprintf("%s=%s", PluginUnwrapTokenEnv, "testtoken"),
					}...),
				),
				SecureConfig: &plugin.SecureConfig{
					Checksum: []byte("some_sha256"),
					// Hash is generated
				},
				AllowedProtocols: []plugin.Protocol{
					plugin.ProtocolNetRPC,
					plugin.ProtocolGRPC,
				},
				Logger:      hclog.NewNullLogger(),
				AutoMTLS:    false,
				SkipHostEnv: true,
			},
			expectTLSConfig: true,
		},
		"metadata mode, AutoMTLS": {
			rc: runConfig{
				command: "echo",
				args:    []string{"foo", "bar"},
				sha256:  []byte("some_sha256"),
				env:     []string{"initial=true"},
				PluginClientConfig: PluginClientConfig{
					PluginSets: map[int]plugin.PluginSet{
						1: {
							"bogus": nil,
						},
					},
					HandshakeConfig: plugin.HandshakeConfig{
						ProtocolVersion:  1,
						MagicCookieKey:   "magic_cookie_key",
						MagicCookieValue: "magic_cookie_value",
					},
					Logger:         hclog.NewNullLogger(),
					IsMetadataMode: true,
					AutoMTLS:       true,
				},
			},

			responseWrapInfoTimes: 0,

			mlockEnabled:      false,
			mlockEnabledTimes: 1,

			expectedConfig: &plugin.ClientConfig{
				HandshakeConfig: plugin.HandshakeConfig{
					ProtocolVersion:  1,
					MagicCookieKey:   "magic_cookie_key",
					MagicCookieValue: "magic_cookie_value",
				},
				VersionedPlugins: map[int]plugin.PluginSet{
					1: {
						"bogus": nil,
					},
				},
				Cmd: commandWithEnv(
					"echo",
					[]string{"foo", "bar"},
					append(os.Environ(), []string{
						"initial=true",
						fmt.Sprintf("%s=%s", PluginVaultVersionEnv, "dummyversion"),
						fmt.Sprintf("%s=%t", PluginMetadataModeEnv, true),
						fmt.Sprintf("%s=%t", PluginAutoMTLSEnv, true),
					}...),
				),
				SecureConfig: &plugin.SecureConfig{
					Checksum: []byte("some_sha256"),
					// Hash is generated
				},
				AllowedProtocols: []plugin.Protocol{
					plugin.ProtocolNetRPC,
					plugin.ProtocolGRPC,
				},
				Logger:      hclog.NewNullLogger(),
				AutoMTLS:    true,
				SkipHostEnv: true,
			},
			expectTLSConfig: false,
		},
		"not-metadata mode, AutoMTLS": {
			rc: runConfig{
				command: "echo",
				args:    []string{"foo", "bar"},
				sha256:  []byte("some_sha256"),
				env:     []string{"initial=true"},
				PluginClientConfig: PluginClientConfig{
					PluginSets: map[int]plugin.PluginSet{
						1: {
							"bogus": nil,
						},
					},
					HandshakeConfig: plugin.HandshakeConfig{
						ProtocolVersion:  1,
						MagicCookieKey:   "magic_cookie_key",
						MagicCookieValue: "magic_cookie_value",
					},
					Logger:         hclog.NewNullLogger(),
					IsMetadataMode: false,
					AutoMTLS:       true,
				},
			},

			responseWrapInfoTimes: 0,

			mlockEnabled:      false,
			mlockEnabledTimes: 1,

			expectedConfig: &plugin.ClientConfig{
				HandshakeConfig: plugin.HandshakeConfig{
					ProtocolVersion:  1,
					MagicCookieKey:   "magic_cookie_key",
					MagicCookieValue: "magic_cookie_value",
				},
				VersionedPlugins: map[int]plugin.PluginSet{
					1: {
						"bogus": nil,
					},
				},
				Cmd: commandWithEnv(
					"echo",
					[]string{"foo", "bar"},
					append(os.Environ(), []string{
						"initial=true",
						fmt.Sprintf("%s=%s", PluginVaultVersionEnv, "dummyversion"),
						fmt.Sprintf("%s=%t", PluginMetadataModeEnv, false),
						fmt.Sprintf("%s=%t", PluginAutoMTLSEnv, true),
					}...),
				),
				SecureConfig: &plugin.SecureConfig{
					Checksum: []byte("some_sha256"),
					// Hash is generated
				},
				AllowedProtocols: []plugin.Protocol{
					plugin.ProtocolNetRPC,
					plugin.ProtocolGRPC,
				},
				Logger:      hclog.NewNullLogger(),
				AutoMTLS:    true,
				SkipHostEnv: true,
			},
			expectTLSConfig: false,
		},
		"image set": {
			rc: runConfig{
				command:  "echo",
				args:     []string{"foo", "bar"},
				sha256:   []byte("some_sha256"),
				env:      []string{"initial=true"},
				image:    "some-image",
				imageTag: "0.1.0",
				PluginClientConfig: PluginClientConfig{
					PluginSets: map[int]plugin.PluginSet{
						1: {
							"bogus": nil,
						},
					},
					HandshakeConfig: plugin.HandshakeConfig{
						ProtocolVersion:  1,
						MagicCookieKey:   "magic_cookie_key",
						MagicCookieValue: "magic_cookie_value",
					},
					Logger:         hclog.NewNullLogger(),
					IsMetadataMode: false,
					AutoMTLS:       true,
				},
			},

			responseWrapInfoTimes: 0,

			mlockEnabled:      false,
			mlockEnabledTimes: 2,

			expectedConfig: &plugin.ClientConfig{
				HandshakeConfig: plugin.HandshakeConfig{
					ProtocolVersion:  1,
					MagicCookieKey:   "magic_cookie_key",
					MagicCookieValue: "magic_cookie_value",
				},
				VersionedPlugins: map[int]plugin.PluginSet{
					1: {
						"bogus": nil,
					},
				},
				Cmd:          nil,
				SecureConfig: nil,
				AllowedProtocols: []plugin.Protocol{
					plugin.ProtocolNetRPC,
					plugin.ProtocolGRPC,
				},
				Logger:              hclog.NewNullLogger(),
				AutoMTLS:            true,
				SkipHostEnv:         true,
				GRPCBrokerMultiplex: true,
				UnixSocketConfig: &plugin.UnixSocketConfig{
					Group: strconv.Itoa(os.Getgid()),
				},
			},
			expectTLSConfig:  false,
			expectRunnerFunc: true,
			skipSecureConfig: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockWrapper := new(mockRunnerUtil)
			mockWrapper.On("ResponseWrapData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(test.responseWrapInfo, test.responseWrapInfoErr)
			mockWrapper.On("MlockEnabled").
				Return(test.mlockEnabled)
			test.rc.Wrapper = mockWrapper
			defer mockWrapper.AssertNumberOfCalls(t, "ResponseWrapData", test.responseWrapInfoTimes)
			defer mockWrapper.AssertNumberOfCalls(t, "MlockEnabled", test.mlockEnabledTimes)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if test.useLegacyEnvLayering {
				t.Setenv(PluginUseLegacyEnvLayering, "true")
			}

			config, err := test.rc.makeConfig(ctx)
			if err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			// The following fields are generated, so we just need to check for existence, not specific value
			// The value must be nilled out before performing a DeepEqual check
			if !test.skipSecureConfig {
				hsh := config.SecureConfig.Hash
				if hsh == nil {
					t.Fatalf("Missing SecureConfig.Hash")
				}
				config.SecureConfig.Hash = nil
			}

			if test.expectTLSConfig && config.TLSConfig == nil {
				t.Fatalf("TLS config expected, got nil")
			}
			if !test.expectTLSConfig && config.TLSConfig != nil {
				t.Fatalf("no TLS config expected, got: %#v", config.TLSConfig)
			}
			config.TLSConfig = nil

			if test.expectRunnerFunc != (config.RunnerFunc != nil) {
				t.Fatalf("expected RunnerFunc: %v, actual: %v", test.expectRunnerFunc, config.RunnerFunc != nil)
			}
			config.RunnerFunc = nil

			require.Equal(t, test.expectedConfig, config)
		})
	}
}

func commandWithEnv(cmd string, args []string, env []string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Env = env
	return c
}

var _ RunnerUtil = &mockRunnerUtil{}

type mockRunnerUtil struct {
	mock.Mock
}

func (m *mockRunnerUtil) VaultVersion(ctx context.Context) (string, error) {
	return "dummyversion", nil
}

func (m *mockRunnerUtil) NewPluginClient(ctx context.Context, config PluginClientConfig) (PluginClient, error) {
	args := m.Called(ctx, config)
	return args.Get(0).(PluginClient), args.Error(1)
}

func (m *mockRunnerUtil) ResponseWrapData(ctx context.Context, data map[string]interface{}, ttl time.Duration, jwt bool) (*wrapping.ResponseWrapInfo, error) {
	args := m.Called(ctx, data, ttl, jwt)
	return args.Get(0).(*wrapping.ResponseWrapInfo), args.Error(1)
}

func (m *mockRunnerUtil) MlockEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockRunnerUtil) ClusterID(ctx context.Context) (string, error) {
	return "1234", nil
}

func TestContainerConfig(t *testing.T) {
	dummySHA, err := hex.DecodeString("abc123")
	if err != nil {
		t.Fatal(err)
	}
	myPID := strconv.Itoa(os.Getpid())
	for name, tc := range map[string]struct {
		rc       runConfig
		expected plugincontainer.Config
	}{
		"image set, no runtime": {
			rc: runConfig{
				command:  "echo",
				args:     []string{"foo", "bar"},
				sha256:   dummySHA,
				env:      []string{"initial=true"},
				image:    "some-image",
				imageTag: "0.1.0",
				PluginClientConfig: PluginClientConfig{
					PluginSets: map[int]plugin.PluginSet{
						1: {
							"bogus": nil,
						},
					},
					HandshakeConfig: plugin.HandshakeConfig{
						ProtocolVersion:  1,
						MagicCookieKey:   "magic_cookie_key",
						MagicCookieValue: "magic_cookie_value",
					},
					Logger:     hclog.NewNullLogger(),
					AutoMTLS:   true,
					Name:       "some-plugin",
					PluginType: consts.PluginTypeCredential,
					Version:    "v0.1.0",
				},
			},
			expected: plugincontainer.Config{
				Image:      "some-image",
				Tag:        "0.1.0",
				SHA256:     "abc123",
				Entrypoint: []string{"echo"},
				Args:       []string{"foo", "bar"},
				Env: []string{
					"initial=true",
					fmt.Sprintf("%s=%s", PluginVaultVersionEnv, "dummyversion"),
					fmt.Sprintf("%s=%t", PluginMetadataModeEnv, false),
					fmt.Sprintf("%s=%t", PluginAutoMTLSEnv, true),
				},
				Labels: map[string]string{
					labelVaultPID:           myPID,
					labelVaultClusterID:     "1234",
					labelVaultPluginName:    "some-plugin",
					labelVaultPluginType:    "auth",
					labelVaultPluginVersion: "v0.1.0",
				},
				Runtime:  consts.DefaultContainerPluginOCIRuntime,
				GroupAdd: os.Getgid(),
			},
		},
		"image set, with runtime": {
			rc: runConfig{
				sha256:   dummySHA,
				image:    "some-image",
				imageTag: "0.1.0",
				runtimeConfig: &pluginruntimeutil.PluginRuntimeConfig{
					OCIRuntime:   "some-oci-runtime",
					CgroupParent: "/cgroup/parent",
					CPU:          1000,
					Memory:       2000,
				},
				PluginClientConfig: PluginClientConfig{
					PluginSets: map[int]plugin.PluginSet{
						1: {
							"bogus": nil,
						},
					},
					HandshakeConfig: plugin.HandshakeConfig{
						ProtocolVersion:  1,
						MagicCookieKey:   "magic_cookie_key",
						MagicCookieValue: "magic_cookie_value",
					},
					Logger:     hclog.NewNullLogger(),
					AutoMTLS:   true,
					Name:       "some-plugin",
					PluginType: consts.PluginTypeCredential,
					Version:    "v0.1.0",
				},
			},
			expected: plugincontainer.Config{
				Image:  "some-image",
				Tag:    "0.1.0",
				SHA256: "abc123",
				Env: []string{
					fmt.Sprintf("%s=%s", PluginVaultVersionEnv, "dummyversion"),
					fmt.Sprintf("%s=%t", PluginMetadataModeEnv, false),
					fmt.Sprintf("%s=%t", PluginAutoMTLSEnv, true),
				},
				Labels: map[string]string{
					labelVaultPID:           myPID,
					labelVaultClusterID:     "1234",
					labelVaultPluginName:    "some-plugin",
					labelVaultPluginType:    "auth",
					labelVaultPluginVersion: "v0.1.0",
				},
				Runtime:      "some-oci-runtime",
				GroupAdd:     os.Getgid(),
				CgroupParent: "/cgroup/parent",
				NanoCpus:     1000,
				Memory:       2000,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			mockWrapper := new(mockRunnerUtil)
			mockWrapper.On("ResponseWrapData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(nil, nil)
			mockWrapper.On("MlockEnabled").
				Return(false)
			tc.rc.Wrapper = mockWrapper
			cmd, _, err := tc.rc.generateCmd(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			cfg, err := tc.rc.containerConfig(context.Background(), cmd.Env)
			require.NoError(t, err)
			require.Equal(t, tc.expected, *cfg)
		})
	}
}
