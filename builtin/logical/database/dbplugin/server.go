package dbplugin

import (
	"crypto/tls"

	"github.com/hashicorp/go-plugin"
)

// Serve is called from within a plugin and wraps the provided
// Database implementation in a databasePluginRPCServer object and starts a
// RPC server.
func Serve(db Database, tlsProvider func() (*tls.Config, error)) {
	dbPlugin := &DatabasePlugin{
		impl: db,
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"database": dbPlugin,
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		TLSProvider:     tlsProvider,
	})
}

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

func (ds *databasePluginRPCServer) CreateUser(args *CreateUserRequest, resp *CreateUserResponse) error {
	var err error
	resp.Username, resp.Password, err = ds.impl.CreateUser(args.Statements, args.UsernameConfig, args.Expiration)

	return err
}

func (ds *databasePluginRPCServer) RenewUser(args *RenewUserRequest, _ *struct{}) error {
	err := ds.impl.RenewUser(args.Statements, args.Username, args.Expiration)

	return err
}

func (ds *databasePluginRPCServer) RevokeUser(args *RevokeUserRequest, _ *struct{}) error {
	err := ds.impl.RevokeUser(args.Statements, args.Username)

	return err
}

func (ds *databasePluginRPCServer) Initialize(args *InitializeRequest, _ *struct{}) error {
	err := ds.impl.Initialize(args.Config, args.VerifyConnection)

	return err
}

func (ds *databasePluginRPCServer) Close(_ struct{}, _ *struct{}) error {
	ds.impl.Close()
	return nil
}
