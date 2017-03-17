package dbs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/rpc"
	"os/exec"
	"sync"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "VAULT_DATABASE_PLUGIN",
	MagicCookieValue: "926a0820-aea2-be28-51d6-83cdf00e8edb",
}

type DatabasePlugin struct {
	impl DatabaseType
}

func (d DatabasePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &databasePluginRPCServer{impl: d.impl}, nil
}

func (DatabasePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &databasePluginRPCClient{client: c}, nil
}

type DatabasePluginClient struct {
	client *plugin.Client
	sync.Mutex

	*databasePluginRPCClient
}

func (dc *DatabasePluginClient) Close() error {
	err := dc.databasePluginRPCClient.Close()
	dc.client.Kill()

	return err
}

func newPluginClient(sys logical.SystemView, command, checksum string) (DatabaseType, error) {
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"database": new(DatabasePlugin),
	}

	// Get a CA TLS Certificate
	CACertBytes, CACert, CAKey, err := pluginutil.GenerateCACert()
	if err != nil {
		return nil, err
	}

	// Use CA to sign a client cert and return a configured TLS config
	clientTLSConfig, err := pluginutil.CreateClientTLSConfig(CACert, CAKey)
	if err != nil {
		return nil, err
	}

	// Use CA to sign a server cert and wrap the values in a response wrapped
	// token.
	wrapToken, err := pluginutil.WrapServerConfig(sys, CACertBytes, CACert, CAKey)
	if err != nil {
		return nil, err
	}

	// Add the response wrap token to the ENV of the plugin
	cmd := exec.Command(command)
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", pluginutil.PluginUnwrapTokenEnv, wrapToken))

	checksumDecoded, err := hex.DecodeString(checksum)
	if err != nil {
		return nil, err
	}

	secureConfig := &plugin.SecureConfig{
		Checksum: checksumDecoded,
		Hash:     sha256.New(),
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             cmd,
		TLSConfig:       clientTLSConfig,
		SecureConfig:    secureConfig,
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("database")
	if err != nil {
		return nil, err
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	databaseRPC := raw.(*databasePluginRPCClient)

	return &DatabasePluginClient{
		client:                  client,
		databasePluginRPCClient: databaseRPC,
	}, nil
}

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

// ---- RPC client domain ----

type databasePluginRPCClient struct {
	client *rpc.Client
}

func (dr *databasePluginRPCClient) Type() string {
	return "plugin"
}

func (dr *databasePluginRPCClient) CreateUser(statements Statements, username, password, expiration string) error {
	req := CreateUserRequest{
		Statements: statements,
		Username:   username,
		Password:   password,
		Expiration: expiration,
	}

	err := dr.client.Call("Plugin.CreateUser", req, &struct{}{})

	return err
}

func (dr *databasePluginRPCClient) RenewUser(statements Statements, username, expiration string) error {
	req := RenewUserRequest{
		Statements: statements,
		Username:   username,
		Expiration: expiration,
	}

	err := dr.client.Call("Plugin.RenewUser", req, &struct{}{})

	return err
}

func (dr *databasePluginRPCClient) RevokeUser(statements Statements, username string) error {
	req := RevokeUserRequest{
		Statements: statements,
		Username:   username,
	}

	err := dr.client.Call("Plugin.RevokeUser", req, &struct{}{})

	return err
}

func (dr *databasePluginRPCClient) Initialize(conf map[string]interface{}) error {
	err := dr.client.Call("Plugin.Initialize", conf, &struct{}{})

	return err
}

func (dr *databasePluginRPCClient) Close() error {
	err := dr.client.Call("Plugin.Close", struct{}{}, &struct{}{})

	return err
}

func (dr *databasePluginRPCClient) GenerateUsername(displayName string) (string, error) {
	var username string
	err := dr.client.Call("Plugin.GenerateUsername", displayName, &username)

	return username, err
}

func (dr *databasePluginRPCClient) GeneratePassword() (string, error) {
	var password string
	err := dr.client.Call("Plugin.GeneratePassword", struct{}{}, &password)

	return password, err
}

func (dr *databasePluginRPCClient) GenerateExpiration(duration time.Duration) (string, error) {
	var expiration string
	err := dr.client.Call("Plugin.GenerateExpiration", duration, &expiration)

	return expiration, err
}

// ---- RPC server domain ----
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

func (ds *databasePluginRPCServer) Close(_ interface{}, _ *struct{}) error {
	ds.impl.Close()
	return nil
}

func (ds *databasePluginRPCServer) GenerateUsername(args string, resp *string) error {
	var err error
	*resp, err = ds.impl.GenerateUsername(args)

	return err
}

func (ds *databasePluginRPCServer) GeneratePassword(_ struct{}, resp *string) error {
	var err error
	*resp, err = ds.impl.GeneratePassword()

	return err
}

func (ds *databasePluginRPCServer) GenerateExpiration(args time.Duration, resp *string) error {
	var err error
	*resp, err = ds.impl.GenerateExpiration(args)

	return err
}

// ---- Request Args domain ----

type CreateUserRequest struct {
	Statements Statements
	Username   string
	Password   string
	Expiration string
}

type RenewUserRequest struct {
	Statements Statements
	Username   string
	Expiration string
}

type RevokeUserRequest struct {
	Statements Statements
	Username   string
}
