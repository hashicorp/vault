package dbplugin

import (
	"fmt"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
)

// NewPluginServer is called from within a plugin and wraps the provided
// DatabaseType implementation in a databasePluginRPCServer object and starts a
// RPC server.
func NewPluginServer(db DatabaseType) {
	dbPlugin := &DatabasePlugin{
		impl: db,
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"database": dbPlugin,
	}

	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		fmt.Println(err)
		return
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		TLSProvider:     pluginutil.VaultPluginTLSProvider,
	})
}

// ---- RPC server domain ----

// databasePluginRPCServer implements an RPC version of DatabaseType and is run
// inside a plugin. It wraps an underlying implementation of DatabaseType.
type databasePluginRPCServer struct {
	impl DatabaseType
}

func (ds *databasePluginRPCServer) Type(_ struct{}, resp *string) error {
	*resp = ds.impl.Type()
	return nil
}

func (ds *databasePluginRPCServer) CreateUser(args *CreateUserRequest, resp *CreateUserResponse) error {
	var err error
	resp.Username, resp.Password, err = ds.impl.CreateUser(args.Statements, args.UsernamePrefix, args.Expiration)

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
