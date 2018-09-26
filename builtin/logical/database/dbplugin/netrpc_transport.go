package dbplugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/rpc"
	"strings"
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

	return dr.client.Call("Plugin.RenewUser", req, &struct{}{})
}

func (dr *databasePluginRPCClient) RevokeUser(_ context.Context, statements Statements, username string) error {
	req := RevokeUserRequestRPC{
		Statements: statements,
		Username:   username,
	}

	return dr.client.Call("Plugin.RevokeUser", req, &struct{}{})
}

func (dr *databasePluginRPCClient) RotateRootCredentials(_ context.Context, statements []string) (saveConf map[string]interface{}, err error) {
	req := RotateRootCredentialsRequestRPC{
		Statements: statements,
	}

	var resp RotateRootCredentialsResponse
	err = dr.client.Call("Plugin.RotateRootCredentials", req, &resp)

	err = json.Unmarshal(resp.Config, &saveConf)
	return saveConf, err
}

func (dr *databasePluginRPCClient) Initialize(_ context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := dr.Init(nil, conf, verifyConnection)
	return err
}

func (dr *databasePluginRPCClient) Init(_ context.Context, conf map[string]interface{}, verifyConnection bool) (saveConf map[string]interface{}, err error) {
	req := InitRequestRPC{
		Config:           conf,
		VerifyConnection: verifyConnection,
	}

	var resp InitResponse
	err = dr.client.Call("Plugin.Init", req, &resp)
	if err != nil {
		if strings.Contains(err.Error(), "can't find method Plugin.Init") {
			req := InitializeRequestRPC{
				Config:           conf,
				VerifyConnection: verifyConnection,
			}

			err = dr.client.Call("Plugin.Initialize", req, &struct{}{})
			if err == nil {
				return conf, nil
			}
		}
		return nil, err
	}

	err = json.Unmarshal(resp.Config, &saveConf)
	return saveConf, err
}

func (dr *databasePluginRPCClient) Close() error {
	return dr.client.Call("Plugin.Close", struct{}{}, &struct{}{})
}

// ---- RPC Request Args Domain ----

type InitializeRequestRPC struct {
	Config           map[string]interface{}
	VerifyConnection bool
}

type InitRequestRPC struct {
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

type RotateRootCredentialsRequestRPC struct {
	Statements []string
}
