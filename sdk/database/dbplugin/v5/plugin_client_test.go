// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dbplugin

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestNewPluginClient(t *testing.T) {
	type testCase struct {
		config       pluginutil.PluginClientConfig
		pluginClient pluginutil.PluginClient
		expectedResp *DatabasePluginClient
		expectedErr  error
	}

	tests := map[string]testCase{
		"happy path": {
			config: testPluginClientConfig(),
			pluginClient: &fakePluginClient{
				connResp:     nil,
				dispenseResp: gRPCClient{client: fakeClient{}},
				dispenseErr:  nil,
			},
			expectedResp: &DatabasePluginClient{
				client: &fakePluginClient{
					connResp:     nil,
					dispenseResp: gRPCClient{client: fakeClient{}},
					dispenseErr:  nil,
				},
				Database: gRPCClient{client: proto.NewDatabaseClient(nil), versionClient: logical.NewPluginVersionClient(nil), doneCtx: context.Context(nil)},
			},
			expectedErr: nil,
		},
		"dispense error": {
			config: testPluginClientConfig(),
			pluginClient: &fakePluginClient{
				connResp:     nil,
				dispenseResp: gRPCClient{},
				dispenseErr:  errors.New("dispense error"),
			},
			expectedResp: nil,
			expectedErr:  errors.New("dispense error"),
		},
		"error unsupported client type": {
			config: testPluginClientConfig(),
			pluginClient: &fakePluginClient{
				connResp:     nil,
				dispenseResp: nil,
				dispenseErr:  nil,
			},
			expectedResp: nil,
			expectedErr:  errors.New("unsupported client type"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			mockWrapper := new(mockRunnerUtil)
			mockWrapper.On("NewPluginClient", ctx, mock.Anything).
				Return(test.pluginClient, nil)
			defer mockWrapper.AssertNumberOfCalls(t, "NewPluginClient", 1)

			resp, err := NewPluginClient(ctx, mockWrapper, test.config)
			if test.expectedErr != nil && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if test.expectedErr == nil && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
			if test.expectedErr == nil && !reflect.DeepEqual(resp, test.expectedResp) {
				t.Fatalf("Actual response: %#v\nExpected response: %#v", resp, test.expectedResp)
			}
		})
	}
}

func testPluginClientConfig() pluginutil.PluginClientConfig {
	return pluginutil.PluginClientConfig{
		Name:            "test-plugin",
		PluginSets:      PluginSets,
		PluginType:      consts.PluginTypeDatabase,
		HandshakeConfig: HandshakeConfig,
		Logger:          log.NewNullLogger(),
		IsMetadataMode:  true,
		AutoMTLS:        true,
	}
}

var _ pluginutil.PluginClient = &fakePluginClient{}

type fakePluginClient struct {
	connResp grpc.ClientConnInterface

	dispenseResp interface{}
	dispenseErr  error
}

func (f *fakePluginClient) Conn() grpc.ClientConnInterface {
	return nil
}

func (f *fakePluginClient) Reload() error {
	return nil
}

func (f *fakePluginClient) Dispense(name string) (interface{}, error) {
	return f.dispenseResp, f.dispenseErr
}

func (f *fakePluginClient) Ping() error {
	return nil
}

func (f *fakePluginClient) Close() error {
	return nil
}

var _ pluginutil.RunnerUtil = &mockRunnerUtil{}

type mockRunnerUtil struct {
	mock.Mock
}

func (m *mockRunnerUtil) VaultVersion(ctx context.Context) (string, error) {
	return "dummyversion", nil
}

func (m *mockRunnerUtil) NewPluginClient(ctx context.Context, config pluginutil.PluginClientConfig) (pluginutil.PluginClient, error) {
	args := m.Called(ctx, config)
	return args.Get(0).(pluginutil.PluginClient), args.Error(1)
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
	return "clusterid", nil
}
