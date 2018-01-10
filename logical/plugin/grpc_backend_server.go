package plugin

import (
	"context"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin/pb"
	log "github.com/mgutz/logxi/v1"
	"google.golang.org/grpc"
)

type backendGRPCPluginServer struct {
	broker  *plugin.GRPCBroker
	backend logical.Backend

	factory func(*logical.BackendConfig) (logical.Backend, error)

	brokeredClient *grpc.ClientConn

	logger log.Logger
}

// Setup dials into the plugin's broker to get a shimmed storage, logger, and
// system view of the backend. This method also instantiates the underlying
// backend through its factory func for the server side of the plugin.
func (b *backendGRPCPluginServer) Setup(ctx context.Context, args *pb.SetupArgs) (*pb.SetupReply, error) {
	// Dial for storage
	brokeredClient, err := b.broker.Dial(args.BrokerId)
	if err != nil {
		return &pb.SetupReply{}, err
	}
	b.brokeredClient = brokeredClient
	storage := newGRPCStorageClient(brokeredClient)
	sysView := newGRPCSystemView(brokeredClient)

	config := &logical.BackendConfig{
		StorageView: storage,
		Logger:      b.logger,
		System:      sysView,
		Config:      args.Config,
	}

	// Call the underlying backend factory after shims have been created
	// to set b.backend
	backend, err := b.factory(config)
	if err != nil {
		return &pb.SetupReply{
			Err: pb.ErrToString(err),
		}, nil
	}
	b.backend = backend

	return &pb.SetupReply{}, nil
}

func (b *backendGRPCPluginServer) HandleRequest(ctx context.Context, args *pb.HandleRequestArgs) (*pb.HandleRequestReply, error) {
	if inMetadataMode() {
		return &pb.HandleRequestReply{}, ErrServerInMetadataMode
	}

	logicalReq, err := pb.ProtoRequestToLogicalRequest(args.Request)
	if err != nil {
		return &pb.HandleRequestReply{}, err
	}

	logicalReq.Storage = newGRPCStorageClient(b.brokeredClient)

	resp, respErr := b.backend.HandleRequest(ctx, logicalReq)

	pbResp, err := pb.LogicalResponseToProtoResp(resp)
	if err != nil {
		return &pb.HandleRequestReply{}, err
	}

	return &pb.HandleRequestReply{
		Response: pbResp,
		Err:      pb.ErrToString(respErr),
	}, nil
}

func (b *backendGRPCPluginServer) SpecialPaths(ctx context.Context, args *pb.Empty) (*pb.SpecialPathsReply, error) {
	paths := b.backend.SpecialPaths()

	return &pb.SpecialPathsReply{
		Paths: &pb.Paths{
			Root:            paths.Root,
			Unauthenticated: paths.Unauthenticated,
			LocalStorage:    paths.LocalStorage,
			SealWrapStorage: paths.SealWrapStorage,
		},
	}, nil
}

func (b *backendGRPCPluginServer) HandleExistenceCheck(ctx context.Context, args *pb.HandleExistenceCheckArgs) (*pb.HandleExistenceCheckReply, error) {
	if inMetadataMode() {
		return &pb.HandleExistenceCheckReply{}, ErrServerInMetadataMode
	}

	logicalReq, err := pb.ProtoRequestToLogicalRequest(args.Request)
	if err != nil {
		return &pb.HandleExistenceCheckReply{}, err
	}
	logicalReq.Storage = newGRPCStorageClient(b.brokeredClient)

	checkFound, exists, err := b.backend.HandleExistenceCheck(ctx, logicalReq)
	return &pb.HandleExistenceCheckReply{
		CheckFound: checkFound,
		Exists:     exists,
		Err:        pb.ErrToString(err),
	}, nil
}

func (b *backendGRPCPluginServer) Cleanup(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	b.backend.Cleanup()

	// Close rpc clients
	b.brokeredClient.Close()
	return &pb.Empty{}, nil
}

func (b *backendGRPCPluginServer) Initialize(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	if inMetadataMode() {
		return &pb.Empty{}, ErrServerInMetadataMode
	}

	err := b.backend.Initialize()
	return &pb.Empty{}, err
}

func (b *backendGRPCPluginServer) InvalidateKey(ctx context.Context, args *pb.InvalidateKeyArgs) (*pb.Empty, error) {
	if inMetadataMode() {
		return &pb.Empty{}, ErrServerInMetadataMode
	}

	b.backend.InvalidateKey(args.Key)
	return &pb.Empty{}, nil
}

func (b *backendGRPCPluginServer) Type(ctx context.Context, _ *pb.Empty) (*pb.TypeReply, error) {
	return &pb.TypeReply{
		Type: uint32(b.backend.Type()),
	}, nil
}

func (b *backendGRPCPluginServer) RegisterLicense(ctx context.Context, _ *pb.RegisterLicenseArgs) (*pb.RegisterLicenseReply, error) {
	if inMetadataMode() {
		return &pb.RegisterLicenseReply{}, ErrServerInMetadataMode
	}

	err := b.backend.RegisterLicense(struct{}{})
	return &pb.RegisterLicenseReply{
		Err: pb.ErrToString(err),
	}, nil
}
