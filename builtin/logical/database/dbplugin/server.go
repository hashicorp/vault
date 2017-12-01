package dbplugin

import (
	"context"
	"crypto/tls"
	"encoding/json"

	"github.com/golang/protobuf/ptypes"
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
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}

// ---- gRPC Server domain ----

type gRPCServer struct {
	impl Database
}

func (s *gRPCServer) Type(context.Context, *Empty) (*TypeResponse, error) {
	t, err := s.impl.Type()
	if err != nil {
		return nil, err
	}

	return &TypeResponse{
		Type: t,
	}, nil
}

func (s *gRPCServer) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	e, err := ptypes.Timestamp(req.Expiration)
	if err != nil {
		return nil, err
	}

	u, p, err := s.impl.CreateUser(ctx, *req.Statements, *req.UsernameConfig, e)

	return &CreateUserResponse{
		Username: u,
		Password: p,
	}, err
}

func (s *gRPCServer) RenewUser(ctx context.Context, req *RenewUserRequest) (*Empty, error) {
	e, err := ptypes.Timestamp(req.Expiration)
	if err != nil {
		return nil, err
	}
	err = s.impl.RenewUser(ctx, *req.Statements, req.Username, e)
	return &Empty{}, err
}

func (s *gRPCServer) RevokeUser(ctx context.Context, req *RevokeUserRequest) (*Empty, error) {
	err := s.impl.RevokeUser(ctx, *req.Statements, req.Username)
	return &Empty{}, err
}

func (s *gRPCServer) Initialize(ctx context.Context, req *InitializeRequest) (*Empty, error) {
	config := map[string]interface{}{}

	err := json.Unmarshal(req.Config, &config)
	if err != nil {
		return nil, err
	}

	err = s.impl.Initialize(ctx, config, req.VerifyConnection)
	return &Empty{}, err
}

func (s *gRPCServer) Close(_ context.Context, _ *Empty) (*Empty, error) {
	s.impl.Close()
	return &Empty{}, nil
}
