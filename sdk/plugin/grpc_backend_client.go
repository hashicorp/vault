// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"errors"
	"math"
	"sync/atomic"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrPluginShutdown       = errors.New("plugin is shut down")
	ErrClientInMetadataMode = errors.New("plugin client can not perform action while in metadata mode")
)

// Validate backendGRPCPluginClient satisfies the logical.Backend interface
var _ logical.Backend = (*backendGRPCPluginClient)(nil)

// backendPluginClient implements logical.Backend and is the
// go-plugin client.
type backendGRPCPluginClient struct {
	broker        *plugin.GRPCBroker
	client        pb.BackendClient
	versionClient logical.PluginVersionClient
	metadataMode  bool

	system logical.SystemView
	logger log.Logger

	// This is used to signal to the Cleanup function that it can proceed
	// because we have a defined server
	cleanupCh chan struct{}

	// server is the grpc server used for serving storage and sysview requests.
	server *atomic.Value

	doneCtx context.Context
}

func (b *backendGRPCPluginClient) Initialize(ctx context.Context, _ *logical.InitializationRequest) error {
	if b.metadataMode {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	reply, err := b.client.Initialize(ctx, &pb.InitializeArgs{}, largeMsgGRPCCallOpts...)
	if err != nil {
		if b.doneCtx.Err() != nil {
			return ErrPluginShutdown
		}

		// If the plugin doesn't have Initialize implemented we should not fail
		// the initialize call; otherwise this could halt startup of vault.
		grpcStatus, ok := status.FromError(err)
		if ok && grpcStatus.Code() == codes.Unimplemented {
			return nil
		}

		return err
	}
	if reply.Err != nil {
		return pb.ProtoErrToErr(reply.Err)
	}

	return nil
}

func (b *backendGRPCPluginClient) HandleRequest(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	if b.metadataMode {
		return nil, ErrClientInMetadataMode
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	protoReq, err := pb.LogicalRequestToProtoRequest(req)
	if err != nil {
		return nil, err
	}

	reply, err := b.client.HandleRequest(ctx, &pb.HandleRequestArgs{
		Request: protoReq,
	}, largeMsgGRPCCallOpts...)
	if err != nil {
		if b.doneCtx.Err() != nil {
			return nil, ErrPluginShutdown
		}

		return nil, err
	}
	resp, err := pb.ProtoResponseToLogicalResponse(reply.Response)
	if err != nil {
		return nil, err
	}
	if reply.Err != nil {
		return resp, pb.ProtoErrToErr(reply.Err)
	}

	return resp, nil
}

func (b *backendGRPCPluginClient) SpecialPaths() *logical.Paths {
	reply, err := b.client.SpecialPaths(b.doneCtx, &pb.Empty{})
	if err != nil {
		return nil
	}

	if reply.Paths == nil {
		return nil
	}

	return &logical.Paths{
		Root:                  reply.Paths.Root,
		Unauthenticated:       reply.Paths.Unauthenticated,
		LocalStorage:          reply.Paths.LocalStorage,
		SealWrapStorage:       reply.Paths.SealWrapStorage,
		WriteForwardedStorage: reply.Paths.WriteForwardedStorage,
	}
}

// System returns vault's system view. The backend client stores the view during
// Setup, so there is no need to shim the system just to get it back.
func (b *backendGRPCPluginClient) System() logical.SystemView {
	return b.system
}

// Logger returns vault's logger. The backend client stores the logger during
// Setup, so there is no need to shim the logger just to get it back.
func (b *backendGRPCPluginClient) Logger() log.Logger {
	return b.logger
}

func (b *backendGRPCPluginClient) HandleExistenceCheck(ctx context.Context, req *logical.Request) (bool, bool, error) {
	if b.metadataMode {
		return false, false, ErrClientInMetadataMode
	}

	protoReq, err := pb.LogicalRequestToProtoRequest(req)
	if err != nil {
		return false, false, err
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()
	reply, err := b.client.HandleExistenceCheck(ctx, &pb.HandleExistenceCheckArgs{
		Request: protoReq,
	}, largeMsgGRPCCallOpts...)
	if err != nil {
		if b.doneCtx.Err() != nil {
			return false, false, ErrPluginShutdown
		}
		return false, false, err
	}
	if reply.Err != nil {
		return false, false, pb.ProtoErrToErr(reply.Err)
	}

	return reply.CheckFound, reply.Exists, nil
}

func (b *backendGRPCPluginClient) Cleanup(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	// Only wait on graceful cleanup if we can establish communication with the
	// plugin, otherwise b.cleanupCh may never get closed.
	if _, err := b.client.Cleanup(ctx, &pb.Empty{}); status.Code(err) != codes.Unavailable {
		// This will block until Setup has run the function to create a new server
		// in b.server. If we stop here before it has a chance to actually start
		// listening, when it starts listening it will immediately error out and
		// exit, which is fine. Overall this ensures that we do not miss stopping
		// the server if it ends up being created after Cleanup is called.
		select {
		case <-b.cleanupCh:
		}
	}
	server := b.server.Load()
	if grpcServer, ok := server.(*grpc.Server); ok && grpcServer != nil {
		grpcServer.GracefulStop()
	}
}

func (b *backendGRPCPluginClient) InvalidateKey(ctx context.Context, key string) {
	if b.metadataMode {
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	b.client.InvalidateKey(ctx, &pb.InvalidateKeyArgs{
		Key: key,
	})
}

func (b *backendGRPCPluginClient) Setup(ctx context.Context, config *logical.BackendConfig) error {
	// Shim logical.Storage
	storageImpl := config.StorageView
	if b.metadataMode {
		storageImpl = &NOOPStorage{}
	}
	storage := &GRPCStorageServer{
		impl: storageImpl,
	}

	// Shim logical.SystemView
	sysViewImpl := config.System
	if b.metadataMode {
		sysViewImpl = &logical.StaticSystemView{}
	}

	events := &GRPCEventsServer{
		impl: config.EventsSender,
	}

	// Register the server in this closure.
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		opts = append(opts, grpc.MaxRecvMsgSize(math.MaxInt32))
		opts = append(opts, grpc.MaxSendMsgSize(math.MaxInt32))

		s := grpc.NewServer(opts...)
		registerSystemViewServer(s, sysViewImpl, config)
		pb.RegisterStorageServer(s, storage)
		pb.RegisterEventsServer(s, events)
		b.server.Store(s)
		close(b.cleanupCh)
		return s
	}
	brokerID := b.broker.NextId()
	go b.broker.AcceptAndServe(brokerID, serverFunc)

	args := &pb.SetupArgs{
		BrokerID:    brokerID,
		Config:      config.Config,
		BackendUUID: config.BackendUUID,
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	reply, err := b.client.Setup(ctx, args)
	if err != nil {
		return err
	}
	if reply.Err != "" {
		return errors.New(reply.Err)
	}

	// Set system and logger for getter methods
	b.system = config.System
	b.logger = config.Logger

	return nil
}

func (b *backendGRPCPluginClient) Type() logical.BackendType {
	reply, err := b.client.Type(b.doneCtx, &pb.Empty{})
	if err != nil {
		return logical.TypeUnknown
	}

	return logical.BackendType(reply.Type)
}

func (b *backendGRPCPluginClient) PluginVersion() logical.PluginVersion {
	reply, err := b.versionClient.Version(b.doneCtx, &logical.Empty{})
	if err != nil {
		if stErr, ok := status.FromError(err); ok {
			if stErr.Code() == codes.Unimplemented {
				return logical.EmptyPluginVersion
			}
		}
		b.Logger().Warn("Unknown error getting plugin version", "err", err)
		return logical.EmptyPluginVersion
	}
	return logical.PluginVersion{
		Version: reply.GetPluginVersion(),
	}
}

func (b *backendGRPCPluginClient) IsExternal() bool {
	return true
}
