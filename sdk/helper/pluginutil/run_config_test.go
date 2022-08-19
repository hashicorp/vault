package pluginutil

import (
	"context"
	"fmt"
	"os/exec"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/version"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/stretchr/testify/mock"
)

func TestMakeConfig(t *testing.T) {
	type testCase struct {
		rc runConfig

		responseWrapInfo      *wrapping.ResponseWrapInfo
		responseWrapInfoErr   error
		responseWrapInfoTimes int

		mlockEnabled      bool
		mlockEnabledTimes int

		expectedConfig  *plugin.ClientConfig
		expectTLSConfig bool
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
					[]string{
						"initial=true",
						fmt.Sprintf("%s=%s", PluginVaultVersionEnv, version.GetVersion().Version),
						fmt.Sprintf("%s=%t", PluginMetadataModeEnv, true),
					},
				),
				SecureConfig: &plugin.SecureConfig{
					Checksum: []byte("some_sha256"),
					// Hash is generated
				},
				AllowedProtocols: []plugin.Protocol{
					plugin.ProtocolNetRPC,
					plugin.ProtocolGRPC,
				},
				Logger:   hclog.NewNullLogger(),
				AutoMTLS: false,
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
					[]string{
						"initial=true",
						fmt.Sprintf("%s=%t", PluginMlockEnabled, true),
						fmt.Sprintf("%s=%s", PluginVaultVersionEnv, version.GetVersion().Version),
						fmt.Sprintf("%s=%t", PluginMetadataModeEnv, false),
						fmt.Sprintf("%s=%s", PluginUnwrapTokenEnv, "testtoken"),
					},
				),
				SecureConfig: &plugin.SecureConfig{
					Checksum: []byte("some_sha256"),
					// Hash is generated
				},
				AllowedProtocols: []plugin.Protocol{
					plugin.ProtocolNetRPC,
					plugin.ProtocolGRPC,
				},
				Logger:   hclog.NewNullLogger(),
				AutoMTLS: false,
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
					[]string{
						"initial=true",
						fmt.Sprintf("%s=%s", PluginVaultVersionEnv, version.GetVersion().Version),
						fmt.Sprintf("%s=%t", PluginMetadataModeEnv, true),
					},
				),
				SecureConfig: &plugin.SecureConfig{
					Checksum: []byte("some_sha256"),
					// Hash is generated
				},
				AllowedProtocols: []plugin.Protocol{
					plugin.ProtocolNetRPC,
					plugin.ProtocolGRPC,
				},
				Logger:   hclog.NewNullLogger(),
				AutoMTLS: true,
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
					[]string{
						"initial=true",
						fmt.Sprintf("%s=%s", PluginVaultVersionEnv, version.GetVersion().Version),
						fmt.Sprintf("%s=%t", PluginMetadataModeEnv, false),
					},
				),
				SecureConfig: &plugin.SecureConfig{
					Checksum: []byte("some_sha256"),
					// Hash is generated
				},
				AllowedProtocols: []plugin.Protocol{
					plugin.ProtocolNetRPC,
					plugin.ProtocolGRPC,
				},
				Logger:   hclog.NewNullLogger(),
				AutoMTLS: true,
			},
			expectTLSConfig: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockWrapper := new(mockRunnerUtil)
			mockWrapper.On("ResponseWrapData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(test.responseWrapInfo, test.responseWrapInfoErr)
			mockWrapper.On("MlockEnabled").
				Return(test.mlockEnabled)
			test.rc.wrapper = mockWrapper
			defer mockWrapper.AssertNumberOfCalls(t, "ResponseWrapData", test.responseWrapInfoTimes)
			defer mockWrapper.AssertNumberOfCalls(t, "MlockEnabled", test.mlockEnabledTimes)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			config, err := test.rc.makeConfig(ctx)
			if err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			// The following fields are generated, so we just need to check for existence, not specific value
			// The value must be nilled out before performing a DeepEqual check
			hsh := config.SecureConfig.Hash
			if hsh == nil {
				t.Fatalf("Missing SecureConfig.Hash")
			}
			config.SecureConfig.Hash = nil

			if test.expectTLSConfig && config.TLSConfig == nil {
				t.Fatalf("TLS config expected, got nil")
			}
			if !test.expectTLSConfig && config.TLSConfig != nil {
				t.Fatalf("no TLS config expected, got: %#v", config.TLSConfig)
			}
			config.TLSConfig = nil

			if !reflect.DeepEqual(config, test.expectedConfig) {
				t.Fatalf("Actual config: %#v\nExpected config: %#v", config, test.expectedConfig)
			}
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
