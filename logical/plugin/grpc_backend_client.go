package plugin

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin/pb"
	log "github.com/mgutz/logxi/v1"
)

// backendPluginClient implements logical.Backend and is the
// go-plugin client.
type backendGRPCPluginClient struct {
	broker       *plugin.GRPCBroker
	client       pb.BackendClient
	metadataMode bool

	system logical.SystemView
	logger log.Logger

	server *grpc.Server
}

func (b *backendGRPCPluginClient) HandleRequest(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	if b.metadataMode {
		return nil, ErrClientInMetadataMode
	}

	protoReq, err := pb.LogicalRequestToProtoRequest(req)
	if err != nil {
		return nil, err
	}

	reply, err := b.client.HandleRequest(ctx, &pb.HandleRequestArgs{
		Request: protoReq,
	})
	if err != nil {
		return nil, err
	}
	resp, err := pb.ProtoResponseToLogicalResponse(reply.Response)
	if err != nil {
		return nil, err
	}
	if reply.Err != "" {
		if reply.Err == logical.ErrUnsupportedOperation.Error() {
			return nil, logical.ErrUnsupportedOperation
		}

		return resp, errors.New(reply.Err)
	}

	return resp, nil
}

func (b *backendGRPCPluginClient) SpecialPaths() *logical.Paths {
	reply, err := b.client.SpecialPaths(context.Background(), &pb.Empty{})
	if err != nil {
		return nil
	}

	return &logical.Paths{
		Root:            reply.Paths.Root,
		Unauthenticated: reply.Paths.Unauthenticated,
		LocalStorage:    reply.Paths.LocalStorage,
		SealWrapStorage: reply.Paths.SealWrapStorage,
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

	reply, err := b.client.HandleExistenceCheck(ctx, &pb.HandleExistenceCheckArgs{
		Request: protoReq,
	})
	if err != nil {
		return false, false, err
	}
	if reply.Err != "" {
		if reply.Err == logical.ErrUnsupportedPath.Error() {
			return false, false, logical.ErrUnsupportedPath
		}
		return false, false, errors.New(reply.Err)
	}

	return reply.CheckFound, reply.Exists, nil
}

func (b *backendGRPCPluginClient) Cleanup() {
	b.client.Cleanup(context.Background(), &pb.Empty{})
}

func (b *backendGRPCPluginClient) Initialize() error {
	if b.metadataMode {
		return ErrClientInMetadataMode
	}
	_, err := b.client.Initialize(context.Background(), &pb.Empty{})
	return err
}

func (b *backendGRPCPluginClient) InvalidateKey(key string) {
	if b.metadataMode {
		return
	}
	b.client.InvalidateKey(context.Background(), &pb.InvalidateKeyArgs{
		Key: key,
	})
}

func (b *backendGRPCPluginClient) Setup(config *logical.BackendConfig) error {
	// Shim logical.Storage
	storageImpl := config.StorageView
	if b.metadataMode {
		storageImpl = &NOOPStorage{}
	}
	storage := &GRPCStorageServer{
		impl: storageImpl,
	}

	// Shim log.Logger
	/*	loggerImpl := config.Logger
		if b.metadataMode {
			loggerImpl = log.NullLog
		}
	*/
	// Shim logical.SystemView
	sysViewImpl := config.System
	if b.metadataMode {
		sysViewImpl = &logical.StaticSystemView{}
	}
	sysView := &gRPCSystemViewServer{
		impl: sysViewImpl,
	}

	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s := grpc.NewServer(opts...)
		pb.RegisterSystemViewServer(s, sysView)
		pb.RegisterStorageServer(s, storage)
		b.server = s
		return s
	}
	brokerID := b.broker.NextId()
	go b.broker.AcceptAndServe(brokerID, serverFunc)

	args := &pb.SetupArgs{
		BrokerId: brokerID,
		Config:   config.Config,
	}

	reply, err := b.client.Setup(context.Background(), args)
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
	reply, err := b.client.Type(context.Background(), &pb.Empty{})
	if err != nil {
		return logical.TypeUnknown
	}

	return logical.BackendType(reply.Type)
}

func (b *backendGRPCPluginClient) RegisterLicense(license interface{}) error {
	if b.metadataMode {
		return ErrClientInMetadataMode
	}

	args := &pb.RegisterLicenseArgs{
	//		License: license,
	}
	reply, err := b.client.RegisterLicense(context.Background(), args)
	if err != nil {
		return err
	}
	if reply.Err != "" {
		return errors.New(reply.Err)
	}

	return nil
}
