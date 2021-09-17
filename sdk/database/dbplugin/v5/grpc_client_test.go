package dbplugin

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"google.golang.org/grpc"
)

func TestGRPCClient_Initialize(t *testing.T) {
	type testCase struct {
		client       proto.DatabaseClient
		req          InitializeRequest
		expectedResp InitializeResponse
		assertErr    errorAssertion
	}

	tests := map[string]testCase{
		"bad config": {
			client: fakeClient{},
			req: InitializeRequest{
				Config: map[string]interface{}{
					"foo": badJSONValue{},
				},
			},
			assertErr: assertErrNotNil,
		},
		"database error": {
			client: fakeClient{
				initErr: errors.New("initialize error"),
			},
			req: InitializeRequest{
				Config: map[string]interface{}{
					"foo": "bar",
				},
			},
			assertErr: assertErrNotNil,
		},
		"happy path": {
			client: fakeClient{
				initResp: &proto.InitializeResponse{
					ConfigData: marshal(t, map[string]interface{}{
						"foo": "bar",
						"baz": "biz",
					}),
				},
			},
			req: InitializeRequest{
				Config: map[string]interface{}{
					"foo": "bar",
				},
			},
			expectedResp: InitializeResponse{
				Config: map[string]interface{}{
					"foo": "bar",
					"baz": "biz",
				},
			},
			assertErr: assertErrNil,
		},
		"JSON number type in initialize request": {
			client: fakeClient{
				initResp: &proto.InitializeResponse{
					ConfigData: marshal(t, map[string]interface{}{
						"foo": "bar",
						"max": "10",
					}),
				},
			},
			req: InitializeRequest{
				Config: map[string]interface{}{
					"foo": "bar",
					"max": json.Number("10"),
				},
			},
			expectedResp: InitializeResponse{
				Config: map[string]interface{}{
					"foo": "bar",
					"max": "10",
				},
			},
			assertErr: assertErrNil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := gRPCClient{
				client:  test.client,
				doneCtx: nil,
			}

			// Context doesn't need to timeout since this is just passed through
			ctx := context.Background()

			resp, err := c.Initialize(ctx, test.req)
			test.assertErr(t, err)

			if !reflect.DeepEqual(resp, test.expectedResp) {
				t.Fatalf("Actual response: %#v\nExpected response: %#v", resp, test.expectedResp)
			}
		})
	}
}

func TestGRPCClient_NewUser(t *testing.T) {
	runningCtx := context.Background()
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	type testCase struct {
		client       proto.DatabaseClient
		req          NewUserRequest
		doneCtx      context.Context
		expectedResp NewUserResponse
		assertErr    errorAssertion
	}

	tests := map[string]testCase{
		"missing password": {
			client: fakeClient{},
			req: NewUserRequest{
				Password:   "",
				Expiration: time.Now(),
			},
			doneCtx:   runningCtx,
			assertErr: assertErrNotNil,
		},
		"bad expiration": {
			client: fakeClient{},
			req: NewUserRequest{
				Password:   "njkvcb8y934u90grsnkjl",
				Expiration: invalidExpiration,
			},
			doneCtx:   runningCtx,
			assertErr: assertErrNotNil,
		},
		"database error": {
			client: fakeClient{
				newUserErr: errors.New("new user error"),
			},
			req: NewUserRequest{
				Password:   "njkvcb8y934u90grsnkjl",
				Expiration: time.Now(),
			},
			doneCtx:   runningCtx,
			assertErr: assertErrNotNil,
		},
		"plugin shut down": {
			client: fakeClient{
				newUserErr: errors.New("new user error"),
			},
			req: NewUserRequest{
				Password:   "njkvcb8y934u90grsnkjl",
				Expiration: time.Now(),
			},
			doneCtx:   cancelledCtx,
			assertErr: assertErrEquals(ErrPluginShutdown),
		},
		"happy path": {
			client: fakeClient{
				newUserResp: &proto.NewUserResponse{
					Username: "new_user",
				},
			},
			req: NewUserRequest{
				Password:   "njkvcb8y934u90grsnkjl",
				Expiration: time.Now(),
			},
			doneCtx: runningCtx,
			expectedResp: NewUserResponse{
				Username: "new_user",
			},
			assertErr: assertErrNil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := gRPCClient{
				client:  test.client,
				doneCtx: test.doneCtx,
			}

			ctx := context.Background()

			resp, err := c.NewUser(ctx, test.req)
			test.assertErr(t, err)

			if !reflect.DeepEqual(resp, test.expectedResp) {
				t.Fatalf("Actual response: %#v\nExpected response: %#v", resp, test.expectedResp)
			}
		})
	}
}

func TestGRPCClient_UpdateUser(t *testing.T) {
	runningCtx := context.Background()
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	type testCase struct {
		client    proto.DatabaseClient
		req       UpdateUserRequest
		doneCtx   context.Context
		assertErr errorAssertion
	}

	tests := map[string]testCase{
		"missing username": {
			client:    fakeClient{},
			req:       UpdateUserRequest{},
			doneCtx:   runningCtx,
			assertErr: assertErrNotNil,
		},
		"missing changes": {
			client: fakeClient{},
			req: UpdateUserRequest{
				Username: "user",
			},
			doneCtx:   runningCtx,
			assertErr: assertErrNotNil,
		},
		"empty password": {
			client: fakeClient{},
			req: UpdateUserRequest{
				Username: "user",
				Password: &ChangePassword{
					NewPassword: "",
				},
			},
			doneCtx:   runningCtx,
			assertErr: assertErrNotNil,
		},
		"zero expiration": {
			client: fakeClient{},
			req: UpdateUserRequest{
				Username: "user",
				Expiration: &ChangeExpiration{
					NewExpiration: time.Time{},
				},
			},
			doneCtx:   runningCtx,
			assertErr: assertErrNotNil,
		},
		"bad expiration": {
			client: fakeClient{},
			req: UpdateUserRequest{
				Username: "user",
				Expiration: &ChangeExpiration{
					NewExpiration: invalidExpiration,
				},
			},
			doneCtx:   runningCtx,
			assertErr: assertErrNotNil,
		},
		"database error": {
			client: fakeClient{
				updateUserErr: errors.New("update user error"),
			},
			req: UpdateUserRequest{
				Username: "user",
				Password: &ChangePassword{
					NewPassword: "asdf",
				},
			},
			doneCtx:   runningCtx,
			assertErr: assertErrNotNil,
		},
		"plugin shut down": {
			client: fakeClient{
				updateUserErr: errors.New("update user error"),
			},
			req: UpdateUserRequest{
				Username: "user",
				Password: &ChangePassword{
					NewPassword: "asdf",
				},
			},
			doneCtx:   cancelledCtx,
			assertErr: assertErrEquals(ErrPluginShutdown),
		},
		"happy path - change password": {
			client: fakeClient{},
			req: UpdateUserRequest{
				Username: "user",
				Password: &ChangePassword{
					NewPassword: "asdf",
				},
			},
			doneCtx:   runningCtx,
			assertErr: assertErrNil,
		},
		"happy path - change expiration": {
			client: fakeClient{},
			req: UpdateUserRequest{
				Username: "user",
				Expiration: &ChangeExpiration{
					NewExpiration: time.Now(),
				},
			},
			doneCtx:   runningCtx,
			assertErr: assertErrNil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := gRPCClient{
				client:  test.client,
				doneCtx: test.doneCtx,
			}

			ctx := context.Background()

			_, err := c.UpdateUser(ctx, test.req)
			test.assertErr(t, err)
		})
	}
}

func TestGRPCClient_DeleteUser(t *testing.T) {
	runningCtx := context.Background()
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	type testCase struct {
		client    proto.DatabaseClient
		req       DeleteUserRequest
		doneCtx   context.Context
		assertErr errorAssertion
	}

	tests := map[string]testCase{
		"missing username": {
			client:    fakeClient{},
			req:       DeleteUserRequest{},
			doneCtx:   runningCtx,
			assertErr: assertErrNotNil,
		},
		"database error": {
			client: fakeClient{
				deleteUserErr: errors.New("delete user error'"),
			},
			req: DeleteUserRequest{
				Username: "user",
			},
			doneCtx:   runningCtx,
			assertErr: assertErrNotNil,
		},
		"plugin shut down": {
			client: fakeClient{
				deleteUserErr: errors.New("delete user error'"),
			},
			req: DeleteUserRequest{
				Username: "user",
			},
			doneCtx:   cancelledCtx,
			assertErr: assertErrEquals(ErrPluginShutdown),
		},
		"happy path": {
			client: fakeClient{},
			req: DeleteUserRequest{
				Username: "user",
			},
			doneCtx:   runningCtx,
			assertErr: assertErrNil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := gRPCClient{
				client:  test.client,
				doneCtx: test.doneCtx,
			}

			ctx := context.Background()

			_, err := c.DeleteUser(ctx, test.req)
			test.assertErr(t, err)
		})
	}
}

func TestGRPCClient_Type(t *testing.T) {
	runningCtx := context.Background()
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	type testCase struct {
		client       proto.DatabaseClient
		doneCtx      context.Context
		expectedType string
		assertErr    errorAssertion
	}

	tests := map[string]testCase{
		"database error": {
			client: fakeClient{
				typeErr: errors.New("type error"),
			},
			doneCtx:   runningCtx,
			assertErr: assertErrNotNil,
		},
		"plugin shut down": {
			client: fakeClient{
				typeErr: errors.New("type error"),
			},
			doneCtx:   cancelledCtx,
			assertErr: assertErrEquals(ErrPluginShutdown),
		},
		"happy path": {
			client: fakeClient{
				typeResp: &proto.TypeResponse{
					Type: "test type",
				},
			},
			doneCtx:      runningCtx,
			expectedType: "test type",
			assertErr:    assertErrNil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := gRPCClient{
				client:  test.client,
				doneCtx: test.doneCtx,
			}

			dbType, err := c.Type()
			test.assertErr(t, err)

			if dbType != test.expectedType {
				t.Fatalf("Actual type: %s Expected type: %s", dbType, test.expectedType)
			}
		})
	}
}

func TestGRPCClient_Close(t *testing.T) {
	runningCtx := context.Background()
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	type testCase struct {
		client    proto.DatabaseClient
		doneCtx   context.Context
		assertErr errorAssertion
	}

	tests := map[string]testCase{
		"database error": {
			client: fakeClient{
				typeErr: errors.New("type error"),
			},
			doneCtx:   runningCtx,
			assertErr: assertErrNotNil,
		},
		"plugin shut down": {
			client: fakeClient{
				typeErr: errors.New("type error"),
			},
			doneCtx:   cancelledCtx,
			assertErr: assertErrEquals(ErrPluginShutdown),
		},
		"happy path": {
			client:    fakeClient{},
			doneCtx:   runningCtx,
			assertErr: assertErrNil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := gRPCClient{
				client:  test.client,
				doneCtx: test.doneCtx,
			}

			err := c.Close()
			test.assertErr(t, err)
		})
	}
}

type errorAssertion func(*testing.T, error)

func assertErrNotNil(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("err expected, got nil")
	}
}

func assertErrNil(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("no error expected, got: %s", err)
	}
}

func assertErrEquals(expectedErr error) errorAssertion {
	return func(t *testing.T, err error) {
		t.Helper()
		if err != expectedErr {
			t.Fatalf("Actual err: %#v Expected err: %#v", err, expectedErr)
		}
	}
}

var _ proto.DatabaseClient = fakeClient{}

type fakeClient struct {
	initResp *proto.InitializeResponse
	initErr  error

	newUserResp *proto.NewUserResponse
	newUserErr  error

	updateUserResp *proto.UpdateUserResponse
	updateUserErr  error

	deleteUserResp *proto.DeleteUserResponse
	deleteUserErr  error

	typeResp *proto.TypeResponse
	typeErr  error

	closeErr error
}

func (f fakeClient) Initialize(context.Context, *proto.InitializeRequest, ...grpc.CallOption) (*proto.InitializeResponse, error) {
	return f.initResp, f.initErr
}

func (f fakeClient) NewUser(context.Context, *proto.NewUserRequest, ...grpc.CallOption) (*proto.NewUserResponse, error) {
	return f.newUserResp, f.newUserErr
}

func (f fakeClient) UpdateUser(context.Context, *proto.UpdateUserRequest, ...grpc.CallOption) (*proto.UpdateUserResponse, error) {
	return f.updateUserResp, f.updateUserErr
}

func (f fakeClient) DeleteUser(context.Context, *proto.DeleteUserRequest, ...grpc.CallOption) (*proto.DeleteUserResponse, error) {
	return f.deleteUserResp, f.deleteUserErr
}

func (f fakeClient) Type(context.Context, *proto.Empty, ...grpc.CallOption) (*proto.TypeResponse, error) {
	return f.typeResp, f.typeErr
}

func (f fakeClient) Close(context.Context, *proto.Empty, ...grpc.CallOption) (*proto.Empty, error) {
	return &proto.Empty{}, f.typeErr
}
