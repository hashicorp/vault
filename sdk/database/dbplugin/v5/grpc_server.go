package dbplugin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ proto.DatabaseServer = (*gRPCServer)(nil)

type gRPCServer struct {
	mu *sync.RWMutex

	dbFactory func() (Database, error)
	dbType    string
	dbs       map[string]Database
}

func (g *gRPCServer) getDatabase(id string) (Database, error) {
	g.mu.RLock()
	db, exists := g.dbs[id]
	if exists {
		g.mu.RUnlock()
		return db, nil
	}

	// Upgrade to a write lock
	g.mu.RUnlock()
	g.mu.Lock()
	defer g.mu.Unlock()

	// Check if another goroutine has added it to the cache
	db, exists = g.dbs[id]
	if exists {
		return db, nil
	}

	db, err := g.dbFactory()
	if err != nil {
		return nil, fmt.Errorf("failed to create database instance: %w", err)
	}
	g.dbs[id] = db

	if g.dbType != "" {
		return db, nil
	}

	dbType, err := db.Type()
	if err != nil {
		// Failed to get the DB type, close it and return - this is a developer error on the plugin side
		db.Close()
		delete(g.dbs, id)
		return nil, fmt.Errorf("failed to get database type: %w", err)
	}

	g.dbType = dbType
	return db, nil
}

// Initialize the database plugin
func (g *gRPCServer) Initialize(ctx context.Context, req *proto.InitializeRequest) (*proto.InitializeResponse, error) {
	rawConfig := structToMap(req.ConfigData)

	dbReq := InitializeRequest{
		Config:           rawConfig,
		VerifyConnection: req.VerifyConnection,
	}

	db, err := g.getDatabase(req.ID)
	if err != nil {
		return &proto.InitializeResponse{}, status.Errorf(codes.Internal, "failed to get Database instance: %w", err)
	}

	dbResp, err := db.Initialize(ctx, dbReq)
	if err != nil {
		return &proto.InitializeResponse{}, status.Errorf(codes.Internal, "failed to initialize: %s", err)
	}

	newConfig, err := mapToStruct(dbResp.Config)
	if err != nil {
		return &proto.InitializeResponse{}, status.Errorf(codes.Internal, "failed to marshal new config to JSON: %s", err)
	}

	resp := &proto.InitializeResponse{
		ConfigData: newConfig,
	}

	return resp, nil
}

func (g *gRPCServer) NewUser(ctx context.Context, req *proto.NewUserRequest) (*proto.NewUserResponse, error) {
	if req.GetUsernameConfig() == nil {
		return &proto.NewUserResponse{}, status.Errorf(codes.InvalidArgument, "missing username config")
	}

	var expiration time.Time

	if req.GetExpiration() != nil {
		exp, err := ptypes.Timestamp(req.GetExpiration())
		if err != nil {
			return &proto.NewUserResponse{}, status.Errorf(codes.InvalidArgument, "unable to parse expiration date: %s", err)
		}
		expiration = exp
	}

	dbReq := NewUserRequest{
		UsernameConfig: UsernameMetadata{
			DisplayName: req.GetUsernameConfig().GetDisplayName(),
			RoleName:    req.GetUsernameConfig().GetRoleName(),
		},
		Password:           req.GetPassword(),
		Expiration:         expiration,
		Statements:         getStatementsFromProto(req.GetStatements()),
		RollbackStatements: getStatementsFromProto(req.GetRollbackStatements()),
	}

	db, err := g.getDatabase(req.ID)
	if err != nil {
		return &proto.NewUserResponse{}, status.Errorf(codes.Internal, "failed to get Database instance: %w", err)
	}

	dbResp, err := db.NewUser(ctx, dbReq)
	if err != nil {
		return &proto.NewUserResponse{}, status.Errorf(codes.Internal, "unable to create new user: %s", err)
	}

	resp := &proto.NewUserResponse{
		Username: dbResp.Username,
	}
	return resp, nil
}

func (g *gRPCServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	if req.GetUsername() == "" {
		return &proto.UpdateUserResponse{}, status.Errorf(codes.InvalidArgument, "no username provided")
	}

	dbReq, err := getUpdateUserRequest(req)
	if err != nil {
		return &proto.UpdateUserResponse{}, status.Errorf(codes.InvalidArgument, err.Error())
	}

	db, err := g.getDatabase(req.ID)
	if err != nil {
		return &proto.UpdateUserResponse{}, status.Errorf(codes.Internal, "failed to get Database instance: %w", err)
	}

	_, err = db.UpdateUser(ctx, dbReq)
	if err != nil {
		return &proto.UpdateUserResponse{}, status.Errorf(codes.Internal, "unable to update user: %s", err)
	}
	return &proto.UpdateUserResponse{}, nil
}

func getUpdateUserRequest(req *proto.UpdateUserRequest) (UpdateUserRequest, error) {
	var password *ChangePassword
	if req.GetPassword() != nil && req.GetPassword().GetNewPassword() != "" {
		password = &ChangePassword{
			NewPassword: req.GetPassword().GetNewPassword(),
			Statements:  getStatementsFromProto(req.GetPassword().GetStatements()),
		}
	}

	var expiration *ChangeExpiration
	if req.GetExpiration() != nil && req.GetExpiration().GetNewExpiration() != nil {
		newExpiration, err := ptypes.Timestamp(req.GetExpiration().GetNewExpiration())
		if err != nil {
			return UpdateUserRequest{}, fmt.Errorf("unable to parse new expiration: %w", err)
		}

		expiration = &ChangeExpiration{
			NewExpiration: newExpiration,
			Statements:    getStatementsFromProto(req.GetExpiration().GetStatements()),
		}
	}

	dbReq := UpdateUserRequest{
		Username:   req.GetUsername(),
		Password:   password,
		Expiration: expiration,
	}

	if !hasChange(dbReq) {
		return UpdateUserRequest{}, fmt.Errorf("update user request has no changes")
	}

	return dbReq, nil
}

func hasChange(dbReq UpdateUserRequest) bool {
	if dbReq.Password != nil && dbReq.Password.NewPassword != "" {
		return true
	}
	if dbReq.Expiration != nil && !dbReq.Expiration.NewExpiration.IsZero() {
		return true
	}
	return false
}

func (g *gRPCServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	if req.GetUsername() == "" {
		return &proto.DeleteUserResponse{}, status.Errorf(codes.InvalidArgument, "no username provided")
	}
	dbReq := DeleteUserRequest{
		Username:   req.GetUsername(),
		Statements: getStatementsFromProto(req.GetStatements()),
	}

	db, err := g.getDatabase(req.ID)
	if err != nil {
		return &proto.DeleteUserResponse{}, status.Errorf(codes.Internal, "failed to get Database instance: %w", err)
	}

	_, err = db.DeleteUser(ctx, dbReq)
	if err != nil {
		return &proto.DeleteUserResponse{}, status.Errorf(codes.Internal, "unable to delete user: %s", err)
	}
	return &proto.DeleteUserResponse{}, nil
}

func (g *gRPCServer) Type(_ context.Context, _ *proto.Empty) (*proto.TypeResponse, error) {
	// Attempt to load a previously set DB type
	g.mu.RLock()
	dbType := g.dbType
	if dbType != "" {
		g.mu.RUnlock()
		resp := &proto.TypeResponse{
			Type: dbType,
		}
		return resp, nil
	}

	// Upgrade to a write lock
	g.mu.RUnlock()
	g.mu.Lock()
	defer g.mu.Unlock()

	// Check if a different goroutine has updated the type
	dbType = g.dbType
	if dbType != "" {
		g.mu.RUnlock()
		resp := &proto.TypeResponse{
			Type: dbType,
		}
		return resp, nil
	}

	// Create a database instance so we can get the type from it.
	// Then immediately close it since we only need the type and don't
	// want to pollute the cache
	db, err := g.dbFactory()
	if err != nil {
		return &proto.TypeResponse{}, status.Errorf(codes.Internal, "failed to create database instance: %s", err)
	}

	dbType, err = db.Type()
	db.Close() // Cleanup before error checking
	if err != nil {
		return &proto.TypeResponse{}, status.Errorf(codes.Internal, "failed to get database type: %s", err)
	}

	g.dbType = dbType
	resp := &proto.TypeResponse{
		Type: dbType,
	}
	return resp, nil
}

func (g *gRPCServer) Close(_ context.Context, req *proto.CloseRequest) (*proto.CloseResponse, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	db, exists := g.dbs[req.ID]
	if !exists {
		// Database has already been removed, just return immediately
		return &proto.CloseResponse{}, nil
	}

	// Delete from the cache before trying to close so the cache is kept clean in the event of an error
	delete(g.dbs, req.ID)

	err := db.Close()
	if err != nil {
		return &proto.CloseResponse{}, status.Errorf(codes.Internal, "unable to close database plugin: %s", err)
	}

	return &proto.CloseResponse{}, nil
}

func getStatementsFromProto(protoStmts *proto.Statements) (statements Statements) {
	if protoStmts == nil {
		return statements
	}
	cmds := protoStmts.GetCommands()
	statements = Statements{
		Commands: cmds,
	}
	return statements
}
