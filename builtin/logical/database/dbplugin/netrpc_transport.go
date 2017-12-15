package dbplugin

import (
	"context"
	"fmt"
	"net/rpc"
	"time"
)

// ---- RPC server domain ----

// databasePluginRPCServer implements an RPC version of Database and is run
// inside a plugin. It wraps an underlying implementation of Database.
type databasePluginRPCServer struct {
	impl Database
}

func (ds *databasePluginRPCServer) Type(_ struct{}, resp *string) error {
	var err error
	*resp, err = ds.impl.Type()
	return err
}

func (ds *databasePluginRPCServer) CreateUser(args *CreateUserRequestRPC, resp *CreateUserResponse) error {
	var err error
	resp.Username, resp.Password, err = ds.impl.CreateUser(context.Background(), args.Statements, args.UsernameConfig, args.Expiration)
	return err
}

func (ds *databasePluginRPCServer) RenewUser(args *RenewUserRequestRPC, _ *struct{}) error {
	err := ds.impl.RenewUser(context.Background(), args.Statements, args.Username, args.Expiration)
	return err
}

func (ds *databasePluginRPCServer) RevokeUser(args *RevokeUserRequestRPC, _ *struct{}) error {
	err := ds.impl.RevokeUser(context.Background(), args.Statements, args.Username)
	return err
}

func (ds *databasePluginRPCServer) Initialize(args *InitializeRequestRPC, _ *struct{}) error {
	err := ds.impl.Initialize(context.Background(), args.Config, args.VerifyConnection)
	return err
}

func (ds *databasePluginRPCServer) Close(_ struct{}, _ *struct{}) error {
	ds.impl.Close()
	return nil
}

// ---- RPC client domain ----
// databasePluginRPCClient implements Database and is used on the client to
// make RPC calls to a plugin.
type databasePluginRPCClient struct {
	client *rpc.Client
}

func (dr *databasePluginRPCClient) Type() (string, error) {
	var dbType string
	err := dr.client.Call("Plugin.Type", struct{}{}, &dbType)

	return fmt.Sprintf("plugin-%s", dbType), err
}

func (dr *databasePluginRPCClient) CreateUser(_ context.Context, statements Statements, usernameConfig UsernameConfig, expiration time.Time) (username string, password string, err error) {
	req := CreateUserRequestRPC{
		Statements:     statements,
		UsernameConfig: usernameConfig,
		Expiration:     expiration,
	}

	var resp CreateUserResponse
	err = dr.client.Call("Plugin.CreateUser", req, &resp)

	return resp.Username, resp.Password, err
}

func (dr *databasePluginRPCClient) RenewUser(_ context.Context, statements Statements, username string, expiration time.Time) error {
	req := RenewUserRequestRPC{
		Statements: statements,
		Username:   username,
		Expiration: expiration,
	}

	err := dr.client.Call("Plugin.RenewUser", req, &struct{}{})

	return err
}

func (dr *databasePluginRPCClient) RevokeUser(_ context.Context, statements Statements, username string) error {
	req := RevokeUserRequestRPC{
		Statements: statements,
		Username:   username,
	}

	err := dr.client.Call("Plugin.RevokeUser", req, &struct{}{})

	return err
}

func (dr *databasePluginRPCClient) Initialize(_ context.Context, conf map[string]interface{}, verifyConnection bool) error {
	req := InitializeRequestRPC{
		Config:           conf,
		VerifyConnection: verifyConnection,
	}

	err := dr.client.Call("Plugin.Initialize", req, &struct{}{})

	return err
}

func (dr *databasePluginRPCClient) Close() error {
	err := dr.client.Call("Plugin.Close", struct{}{}, &struct{}{})

	return err
}

// ---- RPC Request Args Domain ----

type InitializeRequestRPC struct {
	Config           map[string]interface{}
	VerifyConnection bool
}

type CreateUserRequestRPC struct {
	Statements     Statements
	UsernameConfig UsernameConfig
	Expiration     time.Time
}

type RenewUserRequestRPC struct {
	Statements Statements
	Username   string
	Expiration time.Time
}

type RevokeUserRequestRPC struct {
	Statements Statements
	Username   string
}
