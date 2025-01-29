// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dbplugin

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	_ Database                = gRPCClient{}
	_ logical.PluginVersioner = gRPCClient{}

	ErrPluginShutdown = errors.New("plugin shutdown")
)

type gRPCClient struct {
	client        proto.DatabaseClient
	versionClient logical.PluginVersionClient
	doneCtx       context.Context
}

func (c gRPCClient) PluginVersion() logical.PluginVersion {
	version, _ := c.versionClient.Version(context.Background(), &logical.Empty{})
	if version != nil {
		return logical.PluginVersion{Version: version.PluginVersion}
	}
	return logical.EmptyPluginVersion
}

func (c gRPCClient) Initialize(ctx context.Context, req InitializeRequest) (InitializeResponse, error) {
	rpcReq, err := initReqToProto(req)
	if err != nil {
		return InitializeResponse{}, err
	}

	rpcResp, err := c.client.Initialize(ctx, rpcReq)
	if err != nil {
		return InitializeResponse{}, fmt.Errorf("unable to initialize: %s", err.Error())
	}

	return initRespFromProto(rpcResp)
}

func initReqToProto(req InitializeRequest) (*proto.InitializeRequest, error) {
	config, err := mapToStruct(req.Config)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal config: %w", err)
	}

	rpcReq := &proto.InitializeRequest{
		ConfigData:       config,
		VerifyConnection: req.VerifyConnection,
	}
	return rpcReq, nil
}

func initRespFromProto(rpcResp *proto.InitializeResponse) (InitializeResponse, error) {
	newConfig := structToMap(rpcResp.GetConfigData())

	resp := InitializeResponse{
		Config: newConfig,
	}
	return resp, nil
}

func (c gRPCClient) NewUser(ctx context.Context, req NewUserRequest) (NewUserResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, c.doneCtx)
	defer close(quitCh)
	defer cancel()

	rpcReq, err := newUserReqToProto(req)
	if err != nil {
		return NewUserResponse{}, err
	}

	rpcResp, err := c.client.NewUser(ctx, rpcReq)
	if err != nil {
		if c.doneCtx.Err() != nil {
			return NewUserResponse{}, ErrPluginShutdown
		}
		return NewUserResponse{}, fmt.Errorf("unable to create new user: %w", err)
	}

	return newUserRespFromProto(rpcResp)
}

func newUserReqToProto(req NewUserRequest) (*proto.NewUserRequest, error) {
	switch req.CredentialType {
	case CredentialTypePassword:
		if req.Password == "" {
			return nil, fmt.Errorf("missing password credential")
		}
	case CredentialTypeRSAPrivateKey:
		if len(req.PublicKey) == 0 {
			return nil, fmt.Errorf("missing public key credential")
		}
	case CredentialTypeClientCertificate:
		if req.Subject == "" {
			return nil, fmt.Errorf("missing certificate subject")
		}
	default:
		return nil, fmt.Errorf("unknown credential type")
	}

	expiration, err := ptypes.TimestampProto(req.Expiration)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal expiration date: %w", err)
	}

	rpcReq := &proto.NewUserRequest{
		UsernameConfig: &proto.UsernameConfig{
			DisplayName: req.UsernameConfig.DisplayName,
			RoleName:    req.UsernameConfig.RoleName,
		},
		CredentialType: int32(req.CredentialType),
		Password:       req.Password,
		PublicKey:      req.PublicKey,
		Subject:        req.Subject,
		Expiration:     expiration,
		Statements: &proto.Statements{
			Commands: req.Statements.Commands,
		},
		RollbackStatements: &proto.Statements{
			Commands: req.RollbackStatements.Commands,
		},
	}
	return rpcReq, nil
}

func newUserRespFromProto(rpcResp *proto.NewUserResponse) (NewUserResponse, error) {
	resp := NewUserResponse{
		Username: rpcResp.GetUsername(),
	}
	return resp, nil
}

func (c gRPCClient) UpdateUser(ctx context.Context, req UpdateUserRequest) (UpdateUserResponse, error) {
	rpcReq, err := updateUserReqToProto(req)
	if err != nil {
		return UpdateUserResponse{}, err
	}

	rpcResp, err := c.client.UpdateUser(ctx, rpcReq)
	if err != nil {
		if c.doneCtx.Err() != nil {
			return UpdateUserResponse{}, ErrPluginShutdown
		}

		return UpdateUserResponse{}, fmt.Errorf("unable to update user: %w", err)
	}

	return updateUserRespFromProto(rpcResp)
}

func updateUserReqToProto(req UpdateUserRequest) (*proto.UpdateUserRequest, error) {
	if req.Username == "" {
		return nil, fmt.Errorf("missing username")
	}

	if (req.Password == nil || req.Password.NewPassword == "") &&
		(req.PublicKey == nil || len(req.PublicKey.NewPublicKey) == 0) &&
		(req.Expiration == nil || req.Expiration.NewExpiration.IsZero()) {
		return nil, fmt.Errorf("missing changes")
	}

	expiration, err := expirationToProto(req.Expiration)
	if err != nil {
		return nil, fmt.Errorf("unable to parse new expiration date: %w", err)
	}

	var password *proto.ChangePassword
	if req.Password != nil && req.Password.NewPassword != "" {
		password = &proto.ChangePassword{
			NewPassword: req.Password.NewPassword,
			Statements: &proto.Statements{
				Commands: req.Password.Statements.Commands,
			},
		}
	}

	var publicKey *proto.ChangePublicKey
	if req.PublicKey != nil && len(req.PublicKey.NewPublicKey) > 0 {
		publicKey = &proto.ChangePublicKey{
			NewPublicKey: req.PublicKey.NewPublicKey,
			Statements: &proto.Statements{
				Commands: req.PublicKey.Statements.Commands,
			},
		}
	}

	rpcReq := &proto.UpdateUserRequest{
		Username:            req.Username,
		CredentialType:      int32(req.CredentialType),
		Password:            password,
		PublicKey:           publicKey,
		Expiration:          expiration,
		SelfManagedPassword: req.SelfManagedPassword,
	}
	return rpcReq, nil
}

func updateUserRespFromProto(rpcResp *proto.UpdateUserResponse) (UpdateUserResponse, error) {
	// Placeholder for future conversion if data is returned
	return UpdateUserResponse{}, nil
}

func expirationToProto(exp *ChangeExpiration) (*proto.ChangeExpiration, error) {
	if exp == nil {
		return nil, nil
	}

	expiration, err := ptypes.TimestampProto(exp.NewExpiration)
	if err != nil {
		return nil, err
	}

	changeExp := &proto.ChangeExpiration{
		NewExpiration: expiration,
		Statements: &proto.Statements{
			Commands: exp.Statements.Commands,
		},
	}
	return changeExp, nil
}

func (c gRPCClient) DeleteUser(ctx context.Context, req DeleteUserRequest) (DeleteUserResponse, error) {
	rpcReq, err := deleteUserReqToProto(req)
	if err != nil {
		return DeleteUserResponse{}, err
	}

	rpcResp, err := c.client.DeleteUser(ctx, rpcReq)
	if err != nil {
		if c.doneCtx.Err() != nil {
			return DeleteUserResponse{}, ErrPluginShutdown
		}
		return DeleteUserResponse{}, fmt.Errorf("unable to delete user: %w", err)
	}

	return deleteUserRespFromProto(rpcResp)
}

func deleteUserReqToProto(req DeleteUserRequest) (*proto.DeleteUserRequest, error) {
	if req.Username == "" {
		return nil, fmt.Errorf("missing username")
	}

	rpcReq := &proto.DeleteUserRequest{
		Username: req.Username,
		Statements: &proto.Statements{
			Commands: req.Statements.Commands,
		},
	}
	return rpcReq, nil
}

func deleteUserRespFromProto(rpcResp *proto.DeleteUserResponse) (DeleteUserResponse, error) {
	// Placeholder for future conversion if data is returned
	return DeleteUserResponse{}, nil
}

func (c gRPCClient) Type() (string, error) {
	ctx, cancel := getContextWithTimeout(pluginutil.PluginGRPCTimeoutType)
	defer cancel()

	typeResp, err := c.client.Type(ctx, &proto.Empty{})
	if err != nil {
		if c.doneCtx.Err() != nil {
			return "", ErrPluginShutdown
		}
		return "", fmt.Errorf("unable to get database plugin type: %w", err)
	}
	return typeResp.GetType(), nil
}

func (c gRPCClient) Close() error {
	ctx, cancel := getContextWithTimeout(pluginutil.PluginGRPCTimeoutClose)
	defer cancel()

	_, err := c.client.Close(ctx, &proto.Empty{})
	if err != nil {
		if c.doneCtx.Err() != nil {
			return ErrPluginShutdown
		}
		return err
	}
	return nil
}

func getContextWithTimeout(env string) (context.Context, context.CancelFunc) {
	timeout := 1 // default timeout
	if envTimeout, err := strconv.Atoi(os.Getenv(env)); err == nil && envTimeout > 0 {
		timeout = envTimeout
	}
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}
