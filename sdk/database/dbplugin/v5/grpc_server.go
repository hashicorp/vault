// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dbplugin

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"github.com/hashicorp/vault/sdk/helper/base62"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ proto.DatabaseServer = &gRPCServer{}

type gRPCServer struct {
	proto.UnimplementedDatabaseServer
	logical.UnimplementedPluginVersionServer

	// holds the non-multiplexed Database
	// when this is set the plugin does not support multiplexing
	singleImpl Database

	// instances holds the multiplexed Databases
	instances   map[string]Database
	factoryFunc func() (interface{}, error)

	sync.RWMutex
}

func (g *gRPCServer) getOrCreateDatabase(ctx context.Context) (Database, error) {
	g.Lock()
	defer g.Unlock()

	if g.singleImpl != nil {
		return g.singleImpl, nil
	}

	id, err := pluginutil.GetMultiplexIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if db, ok := g.instances[id]; ok {
		return db, nil
	}
	return g.createDatabase(id)
}

// must hold the g.Lock() to call this function
func (g *gRPCServer) createDatabase(id string) (Database, error) {
	db, err := g.factoryFunc()
	if err != nil {
		return nil, err
	}

	database := db.(Database)
	g.instances[id] = database

	return database, nil
}

// getDatabaseInternal returns the database but does not hold a lock
func (g *gRPCServer) getDatabaseInternal(ctx context.Context) (Database, error) {
	if g.singleImpl != nil {
		return g.singleImpl, nil
	}

	id, err := pluginutil.GetMultiplexIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if db, ok := g.instances[id]; ok {
		return db, nil
	}

	return nil, fmt.Errorf("no database instance found")
}

// getDatabase holds a read lock and returns the database
func (g *gRPCServer) getDatabase(ctx context.Context) (Database, error) {
	g.RLock()
	impl, err := g.getDatabaseInternal(ctx)
	g.RUnlock()
	return impl, err
}

// Initialize the database plugin
func (g *gRPCServer) Initialize(ctx context.Context, request *proto.InitializeRequest) (*proto.InitializeResponse, error) {
	impl, err := g.getOrCreateDatabase(ctx)
	if err != nil {
		return nil, err
	}

	rawConfig := structToMap(request.ConfigData)

	dbReq := InitializeRequest{
		Config:           rawConfig,
		VerifyConnection: request.VerifyConnection,
	}

	dbResp, err := impl.Initialize(ctx, dbReq)
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

	impl, err := g.getDatabase(ctx)
	if err != nil {
		return nil, err
	}

	dbReq := NewUserRequest{
		UsernameConfig: UsernameMetadata{
			DisplayName: req.GetUsernameConfig().GetDisplayName(),
			RoleName:    req.GetUsernameConfig().GetRoleName(),
		},
		CredentialType:     CredentialType(req.GetCredentialType()),
		Password:           req.GetPassword(),
		PublicKey:          req.GetPublicKey(),
		Subject:            req.GetSubject(),
		Expiration:         expiration,
		Statements:         getStatementsFromProto(req.GetStatements()),
		RollbackStatements: getStatementsFromProto(req.GetRollbackStatements()),
	}

	dbResp, err := impl.NewUser(ctx, dbReq)
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
		return &proto.UpdateUserResponse{}, status.Error(codes.InvalidArgument, "no username provided")
	}

	dbReq, err := getUpdateUserRequest(req)
	if err != nil {
		return &proto.UpdateUserResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	impl, err := g.getDatabase(ctx)
	if err != nil {
		return nil, err
	}

	_, err = impl.UpdateUser(ctx, dbReq)
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

	var publicKey *ChangePublicKey
	if req.GetPublicKey() != nil && len(req.GetPublicKey().GetNewPublicKey()) > 0 {
		publicKey = &ChangePublicKey{
			NewPublicKey: req.GetPublicKey().GetNewPublicKey(),
			Statements:   getStatementsFromProto(req.GetPublicKey().GetStatements()),
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
		Username:            req.GetUsername(),
		CredentialType:      CredentialType(req.GetCredentialType()),
		Password:            password,
		PublicKey:           publicKey,
		Expiration:          expiration,
		SelfManagedPassword: req.SelfManagedPassword,
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
	if dbReq.PublicKey != nil && len(dbReq.PublicKey.NewPublicKey) > 0 {
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

	impl, err := g.getDatabase(ctx)
	if err != nil {
		return nil, err
	}

	_, err = impl.DeleteUser(ctx, dbReq)
	if err != nil {
		return &proto.DeleteUserResponse{}, status.Errorf(codes.Internal, "unable to delete user: %s", err)
	}
	return &proto.DeleteUserResponse{}, nil
}

func (g *gRPCServer) Type(ctx context.Context, _ *proto.Empty) (*proto.TypeResponse, error) {
	impl, err := g.getOrCreateDatabase(ctx)
	if err != nil {
		return nil, err
	}

	t, err := impl.Type()
	if err != nil {
		return &proto.TypeResponse{}, status.Errorf(codes.Internal, "unable to retrieve type: %s", err)
	}

	resp := &proto.TypeResponse{
		Type: t,
	}
	return resp, nil
}

func (g *gRPCServer) Close(ctx context.Context, _ *proto.Empty) (*proto.Empty, error) {
	g.Lock()

	impl, err := g.getDatabaseInternal(ctx)
	if err != nil {
		g.Unlock()
		return nil, err
	}

	var id string
	if g.singleImpl == nil {
		// only cleanup instances map when multiplexing is supported
		id, err = pluginutil.GetMultiplexIDFromContext(ctx)
		if err != nil {
			g.Unlock()
			return nil, err
		}
		delete(g.instances, id)
	}

	// unlock here so that the subsequent call to Close() does not hold the
	// lock in case the DB is slow to respond
	g.Unlock()

	err = impl.Close()
	if err != nil {
		// The call to Close failed, so we will put the DB instance back in the
		// map. This might not be necessary, but we do this in case anything
		// relies on being able to retry Close.
		g.Lock()
		defer g.Unlock()
		if g.singleImpl == nil {
			// There is a chance that while we were calling Close another DB
			// config was created for the old ID. So we only put it back if
			// it's not set.
			if _, ok := g.instances[id]; !ok {
				g.instances[id] = impl
			}
		}
		return &proto.Empty{}, status.Errorf(codes.Internal, "unable to close database plugin: %s", err)
	}

	return &proto.Empty{}, nil
}

// getOrForceCreateDatabase will create a database even if the multiplexing ID is not present
func (g *gRPCServer) getOrForceCreateDatabase(ctx context.Context) (Database, error) {
	impl, err := g.getOrCreateDatabase(ctx)
	if errors.Is(err, pluginutil.ErrNoMultiplexingIDFound) {
		// if this is called without a multiplexing context, like from the plugin catalog directly,
		// then we won't have a database ID, so let's generate a new database instance
		id, err := base62.Random(10)
		if err != nil {
			return nil, err
		}

		g.Lock()
		defer g.Unlock()
		impl, err = g.createDatabase(id)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return impl, nil
}

// Version forwards the version request to the underlying Database implementation.
func (g *gRPCServer) Version(ctx context.Context, _ *logical.Empty) (*logical.VersionReply, error) {
	impl, err := g.getOrForceCreateDatabase(ctx)
	if err != nil {
		return nil, err
	}

	if versioner, ok := impl.(logical.PluginVersioner); ok {
		return &logical.VersionReply{PluginVersion: versioner.PluginVersion().Version}, nil
	}
	return &logical.VersionReply{}, nil
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
