package dbplugin

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin/pb"
	"github.com/hashicorp/vault/helper/pluginutil"
	log "github.com/mgutz/logxi/v1"
)

// DatabasePluginClient embeds a databasePluginRPCClient and wraps it's Close
// method to also call Kill() on the plugin.Client.
type DatabasePluginClient struct {
	client *plugin.Client
	sync.Mutex

	*gRPCClient
}

func (dc *DatabasePluginClient) Close() error {
	err := dc.gRPCClient.Close()
	dc.client.Kill()

	return err
}

// newPluginClient returns a databaseRPCClient with a connection to a running
// plugin. The client is wrapped in a DatabasePluginClient object to ensure the
// plugin is killed on call of Close().
func newPluginClient(sys pluginutil.RunnerUtil, pluginRunner *pluginutil.PluginRunner, logger log.Logger) (Database, error) {
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"database": new(DatabasePlugin),
	}

	client, err := pluginRunner.Run(sys, pluginMap, handshakeConfig, []string{}, logger)
	if err != nil {
		return nil, err
	}

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

	// We should have a database type now. This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	databaseRPC := raw.(*gRPCClient)

	// Wrap RPC implimentation in DatabasePluginClient
	return &DatabasePluginClient{
		client:     client,
		gRPCClient: databaseRPC,
	}, nil
}

// ---- gRPC client domain ----

type gRPCClient struct {
	client pb.DatabaseClient
}

func (c gRPCClient) Type() (string, error) {
	resp, err := c.client.Type(context.Background(), &pb.Empty{})
	if err != nil {
		return "", err
	}

	return resp.Type, err
}

func (c gRPCClient) CreateUser(ctx context.Context, statements Statements, usernameConfig UsernameConfig, expiration time.Time) (username string, password string, err error) {
	s := &pb.Statements{
		CreationStatements:   statements.CreationStatements,
		RevocationStatements: statements.RevocationStatements,
		RollbackStatements:   statements.RollbackStatements,
		RenewStatements:      statements.RenewStatements,
	}

	u := &pb.UsernameConfig{
		DisplayName: usernameConfig.DisplayName,
		RoleName:    usernameConfig.RoleName,
	}

	t, err := ptypes.TimestampProto(expiration)
	if err != nil {
		return "", "", err
	}

	resp, err := c.client.CreateUser(ctx, &pb.CreateUserRequest{
		Statements:     s,
		UsernameConfig: u,
		Expiration:     t,
	})
	if err != nil {
		return "", "", err
	}

	return resp.Username, resp.Password, err
}

func (c *gRPCClient) RenewUser(ctx context.Context, statements Statements, username string, expiration time.Time) error {
	s := &pb.Statements{
		CreationStatements:   statements.CreationStatements,
		RevocationStatements: statements.RevocationStatements,
		RollbackStatements:   statements.RollbackStatements,
		RenewStatements:      statements.RenewStatements,
	}

	t, err := ptypes.TimestampProto(expiration)
	if err != nil {
		return err
	}

	_, err = c.client.RenewUser(ctx, &pb.RenewUserRequest{
		Statements: s,
		Username:   username,
		Expiration: t,
	})

	return err
}

func (c *gRPCClient) RevokeUser(ctx context.Context, statements Statements, username string) error {
	s := &pb.Statements{
		CreationStatements:   statements.CreationStatements,
		RevocationStatements: statements.RevocationStatements,
		RollbackStatements:   statements.RollbackStatements,
		RenewStatements:      statements.RenewStatements,
	}

	_, err := c.client.RevokeUser(ctx, &pb.RevokeUserRequest{
		Statements: s,
		Username:   username,
	})

	return err
}

func (c *gRPCClient) Initialize(ctx context.Context, config map[string]interface{}, verifyConnection bool) error {
	configRaw, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = c.client.Initialize(ctx, &pb.InitializeRequest{
		Config:           configRaw,
		VerifyConnection: verifyConnection,
	})

	return err
}

func (c *gRPCClient) Close() error {
	_, err := c.client.Close(context.Background(), &pb.Empty{})
	return err
}
