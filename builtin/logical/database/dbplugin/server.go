package dbplugin

import (
	"context"
	"crypto/tls"
	"encoding/json"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin/pb"
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

// ---- gRPC Server domain ----

type gRPCServer struct {
	impl Database
}

func (s *gRPCServer) Type(context.Context, *pb.Empty) (*pb.TypeResponse, error) {
	t, err := s.impl.Type()
	if err != nil {
		return nil, err
	}

	return &pb.TypeResponse{
		Type: t,
	}, nil
}

func (s *gRPCServer) CreateUser(_ context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	statements := Statements{
		CreationStatements:   req.Statements.CreationStatements,
		RevocationStatements: req.Statements.RevocationStatements,
		RollbackStatements:   req.Statements.RollbackStatements,
		RenewStatements:      req.Statements.RenewStatements,
	}

	usernameConfig := UsernameConfig{
		DisplayName: req.UsernameConfig.DisplayName,
		RoleName:    req.UsernameConfig.RoleName,
	}

	e, err := ptypes.Timestamp(req.Expiration)
	if err != nil {
		return nil, err
	}

	u, p, err := s.impl.CreateUser(statements, usernameConfig, e)

	return &pb.CreateUserResponse{
		Username: u,
		Password: p,
	}, err
}

func (s *gRPCServer) RenewUser(_ context.Context, req *pb.RenewUserRequest) (*pb.Empty, error) {
	statements := Statements{
		CreationStatements:   req.Statements.CreationStatements,
		RevocationStatements: req.Statements.RevocationStatements,
		RollbackStatements:   req.Statements.RollbackStatements,
		RenewStatements:      req.Statements.RenewStatements,
	}

	e, err := ptypes.Timestamp(req.Expiration)
	if err != nil {
		return nil, err
	}
	err = s.impl.RenewUser(statements, req.Username, e)
	return &pb.Empty{}, err
}

func (s *gRPCServer) RevokeUser(_ context.Context, req *pb.RevokeUserRequest) (*pb.Empty, error) {
	statements := Statements{
		CreationStatements:   req.Statements.CreationStatements,
		RevocationStatements: req.Statements.RevocationStatements,
		RollbackStatements:   req.Statements.RollbackStatements,
		RenewStatements:      req.Statements.RenewStatements,
	}

	err := s.impl.RevokeUser(statements, req.Username)
	return &pb.Empty{}, err
}

func (s *gRPCServer) Initialize(ctx context.Context, req *pb.InitializeRequest) (*pb.Empty, error) {
	config := map[string]interface{}{}

	err := json.Unmarshal([]byte(req.Config), config)
	if err != nil {
		return nil, err
	}

	err = s.impl.Initialize(config, req.VerifyConnection)
	return &pb.Empty{}, err
}

func (s *gRPCServer) Close(_ context.Context, _ *pb.Empty) (*pb.Empty, error) {
	s.impl.Close()
	return &pb.Empty{}, nil
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
