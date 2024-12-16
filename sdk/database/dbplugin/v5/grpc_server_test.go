// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dbplugin

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// Before minValidSeconds in ptypes package
var invalidExpiration = time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)

func TestGRPCServer_Initialize(t *testing.T) {
	type testCase struct {
		db            Database
		req           *proto.InitializeRequest
		expectedResp  *proto.InitializeResponse
		expectErr     bool
		expectCode    codes.Code
		grpcSetupFunc func(*testing.T, Database) (context.Context, gRPCServer)
	}

	tests := map[string]testCase{
		"database errored": {
			db: fakeDatabase{
				initErr: errors.New("initialization error"),
			},
			req:           &proto.InitializeRequest{},
			expectedResp:  &proto.InitializeResponse{},
			expectErr:     true,
			expectCode:    codes.Internal,
			grpcSetupFunc: testGrpcServer,
		},
		"newConfig can't marshal to JSON": {
			db: fakeDatabase{
				initResp: InitializeResponse{
					Config: map[string]interface{}{
						"bad-data": badJSONValue{},
					},
				},
			},
			req:           &proto.InitializeRequest{},
			expectedResp:  &proto.InitializeResponse{},
			expectErr:     true,
			expectCode:    codes.Internal,
			grpcSetupFunc: testGrpcServer,
		},
		"happy path with config data for multiplexed plugin": {
			db: fakeDatabase{
				initResp: InitializeResponse{
					Config: map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			req: &proto.InitializeRequest{
				ConfigData: marshal(t, map[string]interface{}{
					"foo": "bar",
				}),
			},
			expectedResp: &proto.InitializeResponse{
				ConfigData: marshal(t, map[string]interface{}{
					"foo": "bar",
				}),
			},
			expectErr:     false,
			expectCode:    codes.OK,
			grpcSetupFunc: testGrpcServer,
		},
		"happy path with config data for non-multiplexed plugin": {
			db: fakeDatabase{
				initResp: InitializeResponse{
					Config: map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			req: &proto.InitializeRequest{
				ConfigData: marshal(t, map[string]interface{}{
					"foo": "bar",
				}),
			},
			expectedResp: &proto.InitializeResponse{
				ConfigData: marshal(t, map[string]interface{}{
					"foo": "bar",
				}),
			},
			expectErr:     false,
			expectCode:    codes.OK,
			grpcSetupFunc: testGrpcServerSingleImpl,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			idCtx, g := test.grpcSetupFunc(t, test.db)
			resp, err := g.Initialize(idCtx, test.req)

			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			actualCode := status.Code(err)
			if actualCode != test.expectCode {
				t.Fatalf("Actual code: %s Expected code: %s", actualCode, test.expectCode)
			}

			if !reflect.DeepEqual(resp, test.expectedResp) {
				t.Fatalf("Actual response: %#v\nExpected response: %#v", resp, test.expectedResp)
			}
		})
	}
}

func TestCoerceFloatsToInt(t *testing.T) {
	type testCase struct {
		input    map[string]interface{}
		expected map[string]interface{}
	}

	tests := map[string]testCase{
		"no numbers": {
			input: map[string]interface{}{
				"foo": "bar",
			},
			expected: map[string]interface{}{
				"foo": "bar",
			},
		},
		"raw integers": {
			input: map[string]interface{}{
				"foo": 42,
			},
			expected: map[string]interface{}{
				"foo": 42,
			},
		},
		"floats ": {
			input: map[string]interface{}{
				"foo": 42.2,
			},
			expected: map[string]interface{}{
				"foo": 42.2,
			},
		},
		"floats coerced to ints": {
			input: map[string]interface{}{
				"foo": float64(42),
			},
			expected: map[string]interface{}{
				"foo": int64(42),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := copyMap(test.input)
			coerceFloatsToInt(actual)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("Actual: %#v\nExpected: %#v", actual, test.expected)
			}
		})
	}
}

func copyMap(m map[string]interface{}) map[string]interface{} {
	newMap := map[string]interface{}{}
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

func TestGRPCServer_NewUser(t *testing.T) {
	type testCase struct {
		db           Database
		req          *proto.NewUserRequest
		expectedResp *proto.NewUserResponse
		expectErr    bool
		expectCode   codes.Code
	}

	tests := map[string]testCase{
		"missing username config": {
			db:           fakeDatabase{},
			req:          &proto.NewUserRequest{},
			expectedResp: &proto.NewUserResponse{},
			expectErr:    true,
			expectCode:   codes.InvalidArgument,
		},
		"bad expiration": {
			db: fakeDatabase{},
			req: &proto.NewUserRequest{
				UsernameConfig: &proto.UsernameConfig{
					DisplayName: "dispname",
					RoleName:    "rolename",
				},
				Expiration: &timestamp.Timestamp{
					Seconds: invalidExpiration.Unix(),
				},
			},
			expectedResp: &proto.NewUserResponse{},
			expectErr:    true,
			expectCode:   codes.InvalidArgument,
		},
		"database error": {
			db: fakeDatabase{
				newUserErr: errors.New("new user error"),
			},
			req: &proto.NewUserRequest{
				UsernameConfig: &proto.UsernameConfig{
					DisplayName: "dispname",
					RoleName:    "rolename",
				},
				Expiration: ptypes.TimestampNow(),
			},
			expectedResp: &proto.NewUserResponse{},
			expectErr:    true,
			expectCode:   codes.Internal,
		},
		"happy path with expiration": {
			db: fakeDatabase{
				newUserResp: NewUserResponse{
					Username: "someuser_foo",
				},
			},
			req: &proto.NewUserRequest{
				UsernameConfig: &proto.UsernameConfig{
					DisplayName: "dispname",
					RoleName:    "rolename",
				},
				Expiration: ptypes.TimestampNow(),
			},
			expectedResp: &proto.NewUserResponse{
				Username: "someuser_foo",
			},
			expectErr:  false,
			expectCode: codes.OK,
		},
		"happy path without expiration": {
			db: fakeDatabase{
				newUserResp: NewUserResponse{
					Username: "someuser_foo",
				},
			},
			req: &proto.NewUserRequest{
				UsernameConfig: &proto.UsernameConfig{
					DisplayName: "dispname",
					RoleName:    "rolename",
				},
			},
			expectedResp: &proto.NewUserResponse{
				Username: "someuser_foo",
			},
			expectErr:  false,
			expectCode: codes.OK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			idCtx, g := testGrpcServer(t, test.db)
			resp, err := g.NewUser(idCtx, test.req)

			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			actualCode := status.Code(err)
			if actualCode != test.expectCode {
				t.Fatalf("Actual code: %s Expected code: %s", actualCode, test.expectCode)
			}

			if !reflect.DeepEqual(resp, test.expectedResp) {
				t.Fatalf("Actual response: %#v\nExpected response: %#v", resp, test.expectedResp)
			}
		})
	}
}

func TestGRPCServer_UpdateUser(t *testing.T) {
	type testCase struct {
		db           Database
		req          *proto.UpdateUserRequest
		expectedResp *proto.UpdateUserResponse
		expectErr    bool
		expectCode   codes.Code
	}

	tests := map[string]testCase{
		"missing username": {
			db:           fakeDatabase{},
			req:          &proto.UpdateUserRequest{},
			expectedResp: &proto.UpdateUserResponse{},
			expectErr:    true,
			expectCode:   codes.InvalidArgument,
		},
		"missing changes": {
			db: fakeDatabase{},
			req: &proto.UpdateUserRequest{
				Username: "someuser",
			},
			expectedResp: &proto.UpdateUserResponse{},
			expectErr:    true,
			expectCode:   codes.InvalidArgument,
		},
		"database error": {
			db: fakeDatabase{
				updateUserErr: errors.New("update user error"),
			},
			req: &proto.UpdateUserRequest{
				Username: "someuser",
				Password: &proto.ChangePassword{
					NewPassword: "90ughaino",
				},
			},
			expectedResp: &proto.UpdateUserResponse{},
			expectErr:    true,
			expectCode:   codes.Internal,
		},
		"bad expiration date": {
			db: fakeDatabase{},
			req: &proto.UpdateUserRequest{
				Username: "someuser",
				Expiration: &proto.ChangeExpiration{
					NewExpiration: &timestamp.Timestamp{
						// Before minValidSeconds in ptypes package
						Seconds: invalidExpiration.Unix(),
					},
				},
			},
			expectedResp: &proto.UpdateUserResponse{},
			expectErr:    true,
			expectCode:   codes.InvalidArgument,
		},
		"change password happy path": {
			db: fakeDatabase{},
			req: &proto.UpdateUserRequest{
				Username: "someuser",
				Password: &proto.ChangePassword{
					NewPassword: "90ughaino",
				},
			},
			expectedResp: &proto.UpdateUserResponse{},
			expectErr:    false,
			expectCode:   codes.OK,
		},
		"change expiration happy path": {
			db: fakeDatabase{},
			req: &proto.UpdateUserRequest{
				Username: "someuser",
				Expiration: &proto.ChangeExpiration{
					NewExpiration: ptypes.TimestampNow(),
				},
			},
			expectedResp: &proto.UpdateUserResponse{},
			expectErr:    false,
			expectCode:   codes.OK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			idCtx, g := testGrpcServer(t, test.db)
			resp, err := g.UpdateUser(idCtx, test.req)

			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			actualCode := status.Code(err)
			if actualCode != test.expectCode {
				t.Fatalf("Actual code: %s Expected code: %s", actualCode, test.expectCode)
			}

			if !reflect.DeepEqual(resp, test.expectedResp) {
				t.Fatalf("Actual response: %#v\nExpected response: %#v", resp, test.expectedResp)
			}
		})
	}
}

func TestGRPCServer_DeleteUser(t *testing.T) {
	type testCase struct {
		db           Database
		req          *proto.DeleteUserRequest
		expectedResp *proto.DeleteUserResponse
		expectErr    bool
		expectCode   codes.Code
	}

	tests := map[string]testCase{
		"missing username": {
			db:           fakeDatabase{},
			req:          &proto.DeleteUserRequest{},
			expectedResp: &proto.DeleteUserResponse{},
			expectErr:    true,
			expectCode:   codes.InvalidArgument,
		},
		"database error": {
			db: fakeDatabase{
				deleteUserErr: errors.New("delete user error"),
			},
			req: &proto.DeleteUserRequest{
				Username: "someuser",
			},
			expectedResp: &proto.DeleteUserResponse{},
			expectErr:    true,
			expectCode:   codes.Internal,
		},
		"happy path": {
			db: fakeDatabase{},
			req: &proto.DeleteUserRequest{
				Username: "someuser",
			},
			expectedResp: &proto.DeleteUserResponse{},
			expectErr:    false,
			expectCode:   codes.OK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			idCtx, g := testGrpcServer(t, test.db)
			resp, err := g.DeleteUser(idCtx, test.req)

			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			actualCode := status.Code(err)
			if actualCode != test.expectCode {
				t.Fatalf("Actual code: %s Expected code: %s", actualCode, test.expectCode)
			}

			if !reflect.DeepEqual(resp, test.expectedResp) {
				t.Fatalf("Actual response: %#v\nExpected response: %#v", resp, test.expectedResp)
			}
		})
	}
}

func TestGRPCServer_Type(t *testing.T) {
	type testCase struct {
		db           Database
		expectedResp *proto.TypeResponse
		expectErr    bool
		expectCode   codes.Code
	}

	tests := map[string]testCase{
		"database error": {
			db: fakeDatabase{
				typeErr: errors.New("type error"),
			},
			expectedResp: &proto.TypeResponse{},
			expectErr:    true,
			expectCode:   codes.Internal,
		},
		"happy path": {
			db: fakeDatabase{
				typeResp: "fake database",
			},
			expectedResp: &proto.TypeResponse{
				Type: "fake database",
			},
			expectErr:  false,
			expectCode: codes.OK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			idCtx, g := testGrpcServer(t, test.db)
			resp, err := g.Type(idCtx, &proto.Empty{})

			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			actualCode := status.Code(err)
			if actualCode != test.expectCode {
				t.Fatalf("Actual code: %s Expected code: %s", actualCode, test.expectCode)
			}

			if !reflect.DeepEqual(resp, test.expectedResp) {
				t.Fatalf("Actual response: %#v\nExpected response: %#v", resp, test.expectedResp)
			}
		})
	}
}

func TestGRPCServer_Close(t *testing.T) {
	type testCase struct {
		db            Database
		expectErr     bool
		expectCode    codes.Code
		grpcSetupFunc func(*testing.T, Database) (context.Context, gRPCServer)
		assertFunc    func(t *testing.T, g gRPCServer)
	}

	tests := map[string]testCase{
		"database error": {
			db: fakeDatabase{
				closeErr: errors.New("close error"),
			},
			expectErr:     true,
			expectCode:    codes.Internal,
			grpcSetupFunc: testGrpcServer,
			assertFunc:    nil,
		},
		"happy path for multiplexed plugin": {
			db:            fakeDatabase{},
			expectErr:     false,
			expectCode:    codes.OK,
			grpcSetupFunc: testGrpcServer,
			assertFunc: func(t *testing.T, g gRPCServer) {
				if len(g.instances) != 0 {
					t.Fatalf("err expected instances map to be empty")
				}
			},
		},
		"error path for multiplexed plugin": {
			db: fakeDatabase{
				closeErr: errors.New("close error"),
			},
			expectErr:     true,
			expectCode:    codes.Internal,
			grpcSetupFunc: testGrpcServer,
			assertFunc: func(t *testing.T, g gRPCServer) {
				if len(g.instances) != 1 {
					t.Fatalf("err expected instances map to contain exactly 1 element")
				}
			},
		},
		"happy path for non-multiplexed plugin": {
			db:            fakeDatabase{},
			expectErr:     false,
			expectCode:    codes.OK,
			grpcSetupFunc: testGrpcServerSingleImpl,
			assertFunc:    nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			idCtx, g := test.grpcSetupFunc(t, test.db)
			_, err := g.Close(idCtx, &proto.Empty{})

			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			actualCode := status.Code(err)
			if actualCode != test.expectCode {
				t.Fatalf("Actual code: %s Expected code: %s", actualCode, test.expectCode)
			}

			if test.assertFunc != nil {
				test.assertFunc(t, g)
			}
		})
	}
}

func TestGRPCServer_Version(t *testing.T) {
	type testCase struct {
		db           Database
		expectedResp string
		expectErr    bool
		expectCode   codes.Code
	}

	tests := map[string]testCase{
		"backend that does not implement version": {
			db:           fakeDatabase{},
			expectedResp: "",
			expectErr:    false,
			expectCode:   codes.OK,
		},
		"backend with version": {
			db: fakeDatabaseWithVersion{
				version: "v123",
			},
			expectedResp: "v123",
			expectErr:    false,
			expectCode:   codes.OK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			idCtx, g := testGrpcServer(t, test.db)
			resp, err := g.Version(idCtx, &logical.Empty{})

			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			actualCode := status.Code(err)
			if actualCode != test.expectCode {
				t.Fatalf("Actual code: %s Expected code: %s", actualCode, test.expectCode)
			}

			if !reflect.DeepEqual(resp.PluginVersion, test.expectedResp) {
				t.Fatalf("Actual response: %#v\nExpected response: %#v", resp, test.expectedResp)
			}
		})
	}
}

// testGrpcServer is a test helper that returns a context with an ID set in its
// metadata and a gRPCServer instance for a multiplexed plugin
func testGrpcServer(t *testing.T, db Database) (context.Context, gRPCServer) {
	t.Helper()
	g := gRPCServer{
		factoryFunc: func() (interface{}, error) {
			return db, nil
		},
		instances: make(map[string]Database),
	}

	id := "12345"
	idCtx := idCtx(t, id)
	g.instances[id] = db

	return idCtx, g
}

// testGrpcServerSingleImpl is a test helper that returns a context and a
// gRPCServer instance for a non-multiplexed plugin
func testGrpcServerSingleImpl(t *testing.T, db Database) (context.Context, gRPCServer) {
	t.Helper()
	return context.Background(), gRPCServer{
		singleImpl: db,
	}
}

// idCtx is a test helper that will return a context with the IDs set in its
// metadata
func idCtx(t *testing.T, ids ...string) context.Context {
	t.Helper()
	// Context doesn't need to timeout since this is just passed through
	ctx := context.Background()
	md := metadata.MD{}
	for _, id := range ids {
		md.Append(pluginutil.MultiplexingCtxKey, id)
	}
	return metadata.NewIncomingContext(ctx, md)
}

func marshal(t *testing.T, m map[string]interface{}) *structpb.Struct {
	t.Helper()

	strct, err := mapToStruct(m)
	if err != nil {
		t.Fatalf("unable to marshal to protobuf: %s", err)
	}
	return strct
}

type badJSONValue struct{}

func (badJSONValue) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("this cannot be marshalled to JSON")
}

func (badJSONValue) UnmarshalJSON([]byte) error {
	return fmt.Errorf("this cannot be unmarshalled from JSON")
}

var _ Database = fakeDatabase{}

type fakeDatabase struct {
	initResp InitializeResponse
	initErr  error

	newUserResp NewUserResponse
	newUserErr  error

	updateUserResp UpdateUserResponse
	updateUserErr  error

	deleteUserResp DeleteUserResponse
	deleteUserErr  error

	typeResp string
	typeErr  error

	closeErr error
}

func (e fakeDatabase) Initialize(ctx context.Context, req InitializeRequest) (InitializeResponse, error) {
	return e.initResp, e.initErr
}

func (e fakeDatabase) NewUser(ctx context.Context, req NewUserRequest) (NewUserResponse, error) {
	return e.newUserResp, e.newUserErr
}

func (e fakeDatabase) UpdateUser(ctx context.Context, req UpdateUserRequest) (UpdateUserResponse, error) {
	return e.updateUserResp, e.updateUserErr
}

func (e fakeDatabase) DeleteUser(ctx context.Context, req DeleteUserRequest) (DeleteUserResponse, error) {
	return e.deleteUserResp, e.deleteUserErr
}

func (e fakeDatabase) Type() (string, error) {
	return e.typeResp, e.typeErr
}

func (e fakeDatabase) Close() error {
	return e.closeErr
}

var _ Database = &recordingDatabase{}

type recordingDatabase struct {
	initializeCalls int
	newUserCalls    int
	updateUserCalls int
	deleteUserCalls int
	typeCalls       int
	closeCalls      int

	// recordingDatabase can act as middleware so we can record the calls to other test Database implementations
	next Database
}

func (f *recordingDatabase) Initialize(ctx context.Context, req InitializeRequest) (InitializeResponse, error) {
	f.initializeCalls++
	if f.next == nil {
		return InitializeResponse{}, nil
	}
	return f.next.Initialize(ctx, req)
}

func (f *recordingDatabase) NewUser(ctx context.Context, req NewUserRequest) (NewUserResponse, error) {
	f.newUserCalls++
	if f.next == nil {
		return NewUserResponse{}, nil
	}
	return f.next.NewUser(ctx, req)
}

func (f *recordingDatabase) UpdateUser(ctx context.Context, req UpdateUserRequest) (UpdateUserResponse, error) {
	f.updateUserCalls++
	if f.next == nil {
		return UpdateUserResponse{}, nil
	}
	return f.next.UpdateUser(ctx, req)
}

func (f *recordingDatabase) DeleteUser(ctx context.Context, req DeleteUserRequest) (DeleteUserResponse, error) {
	f.deleteUserCalls++
	if f.next == nil {
		return DeleteUserResponse{}, nil
	}
	return f.next.DeleteUser(ctx, req)
}

func (f *recordingDatabase) Type() (string, error) {
	f.typeCalls++
	if f.next == nil {
		return "recordingDatabase", nil
	}
	return f.next.Type()
}

func (f *recordingDatabase) Close() error {
	f.closeCalls++
	if f.next == nil {
		return nil
	}
	return f.next.Close()
}

type fakeDatabaseWithVersion struct {
	version string
}

func (e fakeDatabaseWithVersion) PluginVersion() logical.PluginVersion {
	return logical.PluginVersion{Version: e.version}
}

func (e fakeDatabaseWithVersion) Initialize(_ context.Context, _ InitializeRequest) (InitializeResponse, error) {
	return InitializeResponse{}, nil
}

func (e fakeDatabaseWithVersion) NewUser(_ context.Context, _ NewUserRequest) (NewUserResponse, error) {
	return NewUserResponse{}, nil
}

func (e fakeDatabaseWithVersion) UpdateUser(_ context.Context, _ UpdateUserRequest) (UpdateUserResponse, error) {
	return UpdateUserResponse{}, nil
}

func (e fakeDatabaseWithVersion) DeleteUser(_ context.Context, _ DeleteUserRequest) (DeleteUserResponse, error) {
	return DeleteUserResponse{}, nil
}

func (e fakeDatabaseWithVersion) Type() (string, error) {
	return "", nil
}

func (e fakeDatabaseWithVersion) Close() error {
	return nil
}

var (
	_ Database                = (*fakeDatabaseWithVersion)(nil)
	_ logical.PluginVersioner = (*fakeDatabaseWithVersion)(nil)
)
