package dbplugin

import (
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
)

// NewPluginServer is called from within a plugin and wraps the provided
// DatabaseType implimentation in a databasePluginRPCServer object and starts a
// RPC server.
func NewPluginServer(db DatabaseType) {
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
		TLSProvider:     pluginutil.VaultPluginTLSProvider,
	})
}

// ---- RPC server domain ----

// databasePluginRPCServer impliments DatabaseType and is run inside a plugin
type databasePluginRPCServer struct {
	impl DatabaseType
}

func (ds *databasePluginRPCServer) Type(_ struct{}, resp *string) error {
	*resp = ds.impl.Type()
	return nil
}

func (ds *databasePluginRPCServer) CreateUser(args *CreateUserRequest, _ *struct{}) error {
	err := ds.impl.CreateUser(args.Statements, args.Username, args.Password, args.Expiration)

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

func (ds *databasePluginRPCServer) Initialize(args map[string]interface{}, _ *struct{}) error {
	err := ds.impl.Initialize(args)

	return err
}

func (ds *databasePluginRPCServer) Close(_ struct{}, _ *struct{}) error {
	ds.impl.Close()
	return nil
}

func (ds *databasePluginRPCServer) GenerateUsername(args string, resp *GenerateUsernameResponse) error {
	var err error
	resp.Username, err = ds.impl.GenerateUsername(args)

	return err
}

func (ds *databasePluginRPCServer) GeneratePassword(_ struct{}, resp *GeneratePasswordResponse) error {
	var err error
	resp.Password, err = ds.impl.GeneratePassword()

	return err
}

func (ds *databasePluginRPCServer) GenerateExpiration(args time.Duration, resp *GenerateExpirationResponse) error {
	var err error
	resp.Expiration, err = ds.impl.GenerateExpiration(args)

	return err
}
