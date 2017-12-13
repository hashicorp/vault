package dbplugin

import (
	"context"
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
