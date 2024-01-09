// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginutil

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestMultiplexingSupported(t *testing.T) {
	type args struct {
		ctx  context.Context
		cc   grpc.ClientConnInterface
		name string
	}

	type testCase struct {
		name    string
		args    args
		env     string
		want    bool
		wantErr bool
	}

	tests := []testCase{
		{
			name: "multiplexing is supported if plugin is not opted out",
			args: args{
				ctx:  context.Background(),
				cc:   &MockClientConnInterfaceNoop{},
				name: "plugin",
			},
			env:  "",
			want: true,
		},
		{
			name: "multiplexing is not supported if plugin is opted out",
			args: args{
				ctx:  context.Background(),
				cc:   &MockClientConnInterfaceNoop{},
				name: "optedOutPlugin",
			},
			env:  "optedOutPlugin",
			want: false,
		},
		{
			name: "multiplexing is not supported if plugin among one of the opted out",
			args: args{
				ctx:  context.Background(),
				cc:   &MockClientConnInterfaceNoop{},
				name: "optedOutPlugin",
			},
			env:  "firstPlugin,optedOutPlugin,otherPlugin",
			want: false,
		},
		{
			name: "multiplexing is supported if different plugin is opted out",
			args: args{
				ctx:  context.Background(),
				cc:   &MockClientConnInterfaceNoop{},
				name: "plugin",
			},
			env:  "optedOutPlugin",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(PluginMultiplexingOptOut, tt.env)
			got, err := MultiplexingSupported(tt.args.ctx, tt.args.cc, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MultiplexingSupported() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MultiplexingSupported() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMultiplexIDFromContext(t *testing.T) {
	type testCase struct {
		ctx          context.Context
		expectedResp string
		expectedErr  error
	}

	tests := map[string]testCase{
		"missing plugin multiplexing metadata": {
			ctx:          context.Background(),
			expectedResp: "",
			expectedErr:  fmt.Errorf("missing plugin multiplexing metadata"),
		},
		"unexpected number of IDs in metadata": {
			ctx:          idCtx(t, "12345", "67891"),
			expectedResp: "",
			expectedErr:  fmt.Errorf("unexpected number of IDs in metadata: (2)"),
		},
		"empty multiplex ID in metadata": {
			ctx:          idCtx(t, ""),
			expectedResp: "",
			expectedErr:  fmt.Errorf("empty multiplex ID in metadata"),
		},
		"happy path, id is returned from metadata": {
			ctx:          idCtx(t, "12345"),
			expectedResp: "12345",
			expectedErr:  nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := GetMultiplexIDFromContext(test.ctx)

			if test.expectedErr != nil && test.expectedErr.Error() != "" && err == nil {
				t.Fatalf("err expected, got nil")
			} else if !reflect.DeepEqual(err, test.expectedErr) {
				t.Fatalf("Actual error: %#v\nExpected error: %#v", err, test.expectedErr)
			}

			if test.expectedErr != nil && test.expectedErr.Error() == "" && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(resp, test.expectedResp) {
				t.Fatalf("Actual response: %#v\nExpected response: %#v", resp, test.expectedResp)
			}
		})
	}
}

// idCtx is a test helper that will return a context with the IDs set in its
// metadata
func idCtx(t *testing.T, ids ...string) context.Context {
	// Context doesn't need to timeout since this is just passed through
	ctx := context.Background()
	md := metadata.MD{}
	for _, id := range ids {
		md.Append(MultiplexingCtxKey, id)
	}
	return metadata.NewIncomingContext(ctx, md)
}

type MockClientConnInterfaceNoop struct{}

func (m *MockClientConnInterfaceNoop) Invoke(_ context.Context, _ string, _ interface{}, reply interface{}, _ ...grpc.CallOption) error {
	reply.(*MultiplexingSupportResponse).Supported = true
	return nil
}

func (m *MockClientConnInterfaceNoop) NewStream(_ context.Context, _ *grpc.StreamDesc, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, nil
}
