package dbplugin

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

var (
	ErrPluginShutdown          = errors.New("plugin shutdown")
	ErrPluginStaticUnsupported = errors.New("database plugin does not support Static Accounts")
)

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

func (s *gRPCServer) RotateRootCredentials(ctx context.Context, req *RotateRootCredentialsRequest) (*RotateRootCredentialsResponse, error) {
	resp, err := s.impl.RotateRootCredentials(ctx, req.Statements)
	if err != nil {
		return nil, err
	}

	respConfig, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return &RotateRootCredentialsResponse{
		Config: respConfig,
	}, err
}

func (s *gRPCServer) Initialize(ctx context.Context, req *InitializeRequest) (*Empty, error) {
	_, err := s.Init(ctx, &InitRequest{
		Config:           req.Config,
		VerifyConnection: req.VerifyConnection,
	})
	return &Empty{}, err
}

func (s *gRPCServer) Init(ctx context.Context, req *InitRequest) (*InitResponse, error) {
	config := map[string]interface{}{}
	err := json.Unmarshal(req.Config, &config)
	if err != nil {
		return nil, err
	}

	resp, err := s.impl.Init(ctx, config, req.VerifyConnection)
	if err != nil {
		return nil, err
	}

	respConfig, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return &InitResponse{
		Config: respConfig,
	}, err
}

func (s *gRPCServer) Close(_ context.Context, _ *Empty) (*Empty, error) {
	s.impl.Close()
	return &Empty{}, nil
}

func (s *gRPCServer) GenerateCredentials(ctx context.Context, _ *Empty) (*GenerateCredentialsResponse, error) {
	p, err := s.impl.GenerateCredentials(ctx)
	if err != nil {
		return nil, err
	}

	return &GenerateCredentialsResponse{
		Password: p,
	}, nil
}

func (s *gRPCServer) SetCredentials(ctx context.Context, req *SetCredentialsRequest) (*SetCredentialsResponse, error) {
	username, password, err := s.impl.SetCredentials(ctx, *req.Statements, *req.StaticUserConfig)
	if err != nil {
		return nil, err
	}

	return &SetCredentialsResponse{
		Username: username,
		Password: password,
	}, err
}

// ---- gRPC client domain ----

type gRPCClient struct {
	client     DatabaseClient
	clientConn *grpc.ClientConn

	doneCtx context.Context
}

func (c *gRPCClient) Type() (string, error) {
	resp, err := c.client.Type(c.doneCtx, &Empty{})
	if err != nil {
		return "", err
	}

	return resp.Type, err
}

func (c *gRPCClient) CreateUser(ctx context.Context, statements Statements, usernameConfig UsernameConfig, expiration time.Time) (username string, password string, err error) {
	t, err := ptypes.TimestampProto(expiration)
	if err != nil {
		return "", "", err
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, c.doneCtx)
	defer close(quitCh)
	defer cancel()

	resp, err := c.client.CreateUser(ctx, &CreateUserRequest{
		Statements:     &statements,
		UsernameConfig: &usernameConfig,
		Expiration:     t,
	})
	if err != nil {
		if c.doneCtx.Err() != nil {
			return "", "", ErrPluginShutdown
		}

		return "", "", err
	}

	return resp.Username, resp.Password, err
}

func (c *gRPCClient) RenewUser(ctx context.Context, statements Statements, username string, expiration time.Time) error {
	t, err := ptypes.TimestampProto(expiration)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, c.doneCtx)
	defer close(quitCh)
	defer cancel()

	_, err = c.client.RenewUser(ctx, &RenewUserRequest{
		Statements: &statements,
		Username:   username,
		Expiration: t,
	})
	if err != nil {
		if c.doneCtx.Err() != nil {
			return ErrPluginShutdown
		}

		return err
	}

	return nil
}

func (c *gRPCClient) RevokeUser(ctx context.Context, statements Statements, username string) error {
	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, c.doneCtx)
	defer close(quitCh)
	defer cancel()

	_, err := c.client.RevokeUser(ctx, &RevokeUserRequest{
		Statements: &statements,
		Username:   username,
	})
	if err != nil {
		if c.doneCtx.Err() != nil {
			return ErrPluginShutdown
		}

		return err
	}

	return nil
}

func (c *gRPCClient) RotateRootCredentials(ctx context.Context, statements []string) (conf map[string]interface{}, err error) {
	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, c.doneCtx)
	defer close(quitCh)
	defer cancel()

	resp, err := c.client.RotateRootCredentials(ctx, &RotateRootCredentialsRequest{
		Statements: statements,
	})
	if err != nil {
		if c.doneCtx.Err() != nil {
			return nil, ErrPluginShutdown
		}

		return nil, err
	}

	if err := json.Unmarshal(resp.Config, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *gRPCClient) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := c.Init(ctx, conf, verifyConnection)
	return err
}

func (c *gRPCClient) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (map[string]interface{}, error) {
	configRaw, err := json.Marshal(conf)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, c.doneCtx)
	defer close(quitCh)
	defer cancel()

	resp, err := c.client.Init(ctx, &InitRequest{
		Config:           configRaw,
		VerifyConnection: verifyConnection,
	})
	if err != nil {
		// Fall back to old call if not implemented
		grpcStatus, ok := status.FromError(err)
		if ok && grpcStatus.Code() == codes.Unimplemented {
			_, err = c.client.Initialize(ctx, &InitializeRequest{
				Config:           configRaw,
				VerifyConnection: verifyConnection,
			})
			if err == nil {
				return conf, nil
			}
		}

		if c.doneCtx.Err() != nil {
			return nil, ErrPluginShutdown
		}
		return nil, err
	}

	if err := json.Unmarshal(resp.Config, &conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func (c *gRPCClient) Close() error {
	_, err := c.client.Close(c.doneCtx, &Empty{})
	return err
}

func (c *gRPCClient) GenerateCredentials(ctx context.Context) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, c.doneCtx)
	defer close(quitCh)
	defer cancel()

	resp, err := c.client.GenerateCredentials(ctx, &Empty{})
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if ok && grpcStatus.Code() == codes.Unimplemented {
			return "", ErrPluginStaticUnsupported
		}

		if c.doneCtx.Err() != nil {
			return "", ErrPluginShutdown
		}
		return "", err
	}

	return resp.Password, nil
}

func (c *gRPCClient) SetCredentials(ctx context.Context, statements Statements, staticUser StaticUserConfig) (username, password string, err error) {
	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, c.doneCtx)
	defer close(quitCh)
	defer cancel()

	resp, err := c.client.SetCredentials(ctx, &SetCredentialsRequest{
		StaticUserConfig: &staticUser,
		Statements:       &statements,
	})
	if err != nil {
		// Fall back to old call if not implemented
		grpcStatus, ok := status.FromError(err)
		if ok && grpcStatus.Code() == codes.Unimplemented {
			return "", "", ErrPluginStaticUnsupported
		}

		if c.doneCtx.Err() != nil {
			return "", "", ErrPluginShutdown
		}
		return "", "", err
	}

	return resp.Username, resp.Password, err
}
