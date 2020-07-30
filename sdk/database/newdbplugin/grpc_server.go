package newdbplugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/vault/sdk/database/newdbplugin/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ proto.DatabaseServer = gRPCServer{}

type gRPCServer struct {
	impl Database
}

// Initialize the database plugin
func (g gRPCServer) Initialize(ctx context.Context, request *proto.InitializeRequest) (*proto.InitializeResponse, error) {
	// Parse the config back from JSON to a map[string]interface{}
	rawConfig, err := parseConfigData(request.ConfigData)
	if err != nil {
		return &proto.InitializeResponse{}, status.Errorf(codes.InvalidArgument, "unable to parse config data: %s", err)
	}

	dbReq := InitializeRequest{
		Config:           rawConfig,
		VerifyConnection: request.VerifyConnection,
	}

	dbResp, err := g.impl.Initialize(ctx, dbReq)
	if err != nil {
		return &proto.InitializeResponse{}, status.Errorf(codes.Internal, "failed to initialize: %s", err)
	}

	newConfig, err := json.Marshal(dbResp.Config)
	if err != nil {
		return &proto.InitializeResponse{}, status.Errorf(codes.Internal, "failed to marshal new config to JSON: %s", err)
	}

	resp := &proto.InitializeResponse{
		ConfigData: newConfig,
	}

	return resp, nil
}

func parseConfigData(b []byte) (map[string]interface{}, error) {
	config := map[string]interface{}{}
	if len(b) == 0 {
		return config, nil
	}
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	err := decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	cleanNumbers(config)
	return config, nil
}

func cleanNumbers(data map[string]interface{}) {
	for k, v := range data {
		switch val := v.(type) {
		case json.Number:
			newNum, err := coerceToScalarNumber(val)
			if err != nil {
				continue
			}
			data[k] = newNum
		case map[string]interface{}:
			cleanNumbers(val)
		}
	}
}

func coerceToScalarNumber(num json.Number) (newVal interface{}, err error) {
	intNum, err := num.Int64()
	if err == nil {
		return intNum, nil
	}

	floatNum, err := num.Float64()
	if err == nil {
		return floatNum, nil
	}

	return nil, fmt.Errorf("unrecognized number: %w", err)
}

func (g gRPCServer) NewUser(ctx context.Context, req *proto.NewUserRequest) (*proto.NewUserResponse, error) {
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
		Statements: Statements{
			Commands: req.GetStatements().GetCommands(),
		},
		Password:   req.GetPassword(),
		Expiration: expiration,
	}

	dbResp, err := g.impl.NewUser(ctx, dbReq)
	if err != nil {
		return &proto.NewUserResponse{}, status.Errorf(codes.Internal, "unable to create new user: %s", err)
	}

	resp := &proto.NewUserResponse{
		Username: dbResp.Username,
	}
	return resp, nil
}

func (g gRPCServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	if req.GetUsername() == "" {
		return &proto.UpdateUserResponse{}, status.Errorf(codes.InvalidArgument, "no username provided")
	}

	dbReq, err := getUpdateUserRequest(req)
	if err != nil {
		return &proto.UpdateUserResponse{}, status.Errorf(codes.InvalidArgument, err.Error())
	}

	_, err = g.impl.UpdateUser(ctx, dbReq)
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
			Statements: Statements{
				Commands: req.GetExpiration().GetStatements().GetCommands(),
			},
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

func (g gRPCServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	if req.GetUsername() == "" {
		return &proto.DeleteUserResponse{}, status.Errorf(codes.InvalidArgument, "no username provided")
	}
	dbReq := DeleteUserRequest{
		Username: req.GetUsername(),
	}

	_, err := g.impl.DeleteUser(ctx, dbReq)
	if err != nil {
		return &proto.DeleteUserResponse{}, status.Errorf(codes.Internal, "unable to delete user: %s", err)
	}
	return &proto.DeleteUserResponse{}, nil
}

func (g gRPCServer) Type(ctx context.Context, _ *proto.Empty) (*proto.TypeResponse, error) {
	t, err := g.impl.Type()
	if err != nil {
		return &proto.TypeResponse{}, status.Errorf(codes.Internal, "unable to retrieve type: %s", err)
	}

	resp := &proto.TypeResponse{
		Type: t,
	}
	return resp, nil
}

func (g gRPCServer) Close(ctx context.Context, _ *proto.Empty) (*proto.Empty, error) {
	err := g.impl.Close()
	if err != nil {
		return &proto.Empty{}, status.Errorf(codes.Internal, "unable to close database plugin: %s", err)
	}
	return &proto.Empty{}, nil
}
