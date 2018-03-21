package dbplugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/rpc"
	"strings"
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

func (ds *databasePluginRPCServer) RotateRootCredentials(args *RotateRootCredentialsRequestRPC, resp *RotateRootCredentialsResponse) error {
	config, err := ds.impl.RotateRootCredentials(context.Background(), args.Statements)
	if err != nil {
		return err
	}
	resp.Config, err = json.Marshal(config)
	return err
}

func (ds *databasePluginRPCServer) Initialize(args *InitializeRequestRPC, _ *struct{}) error {
	return ds.Init(&InitRequestRPC{
		Config:           args.Config,
		VerifyConnection: args.VerifyConnection,
	}, &InitResponse{})
}

func (ds *databasePluginRPCServer) Init(args *InitRequestRPC, resp *InitResponse) error {
	config, err := ds.impl.Init(context.Background(), args.Config, args.VerifyConnection)
	if err != nil {
		return err
	}
	resp.Config, err = json.Marshal(config)
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
