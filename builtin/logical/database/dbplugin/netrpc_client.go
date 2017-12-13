package dbplugin

import (
	"context"
	"fmt"
	"net/rpc"
	"time"
)

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
