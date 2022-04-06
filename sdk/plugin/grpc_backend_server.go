package plugin

import (
	"context"
	"errors"
	"fmt"
	"sync"

	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
	"google.golang.org/grpc"
)

var ErrServerInMetadataMode = errors.New("plugin server can not perform action while in metadata mode")

type backendInstance struct {
	brokeredClient *grpc.ClientConn
	backend        logical.Backend
}

type backendGRPCPluginServer struct {
	pb.UnimplementedBackendServer

	broker *plugin.GRPCBroker

	instances map[string]backendInstance
	sync.RWMutex

	factory logical.Factory

	logger log.Logger
}

// getBackendInternal returns the backend but does not hold a lock
func (b *backendGRPCPluginServer) getBackendInternal(ctx context.Context) (logical.Backend, error) {
	if singleImpl, ok := b.instances["single"]; ok {
		return singleImpl.backend, nil
	}

	id, err := pluginutil.GetMultiplexIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if inst, ok := b.instances[id]; ok {
		return inst.backend, nil
	}

	return nil, fmt.Errorf("no backend instance found")
}

// getBackend holds a read lock and returns the backend
func (b *backendGRPCPluginServer) getBackend(ctx context.Context) (logical.Backend, error) {
	b.RLock()
	defer b.RUnlock()
	return b.getBackendInternal(ctx)
}

// getBrokeredClientInternal returns the brokeredClient but does not hold a lock
func (b *backendGRPCPluginServer) getBrokeredClientInternal(ctx context.Context) (*grpc.ClientConn, error) {
	if singleImpl, ok := b.instances["single"]; ok {
		return singleImpl.brokeredClient, nil
	}

	id, err := pluginutil.GetMultiplexIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if inst, ok := b.instances[id]; ok {
		return inst.brokeredClient, nil
	}

	return nil, fmt.Errorf("no backend instance found")
}

// getBrokeredClient holds a read lock and returns the backend
func (b *backendGRPCPluginServer) getBrokeredClient(ctx context.Context) (*grpc.ClientConn, error) {
	b.RLock()
	defer b.RUnlock()
	return b.getBrokeredClientInternal(ctx)
}

// Setup dials into the plugin's broker to get a shimmed storage, logger, and
// system view of the backend. This method also instantiates the underlying
// backend through its factory func for the server side of the plugin.
func (b *backendGRPCPluginServer) Setup(ctx context.Context, args *pb.SetupArgs) (*pb.SetupReply, error) {
	var err error
	id := "single"

	if _, ok := b.instances[id]; !ok {
		id, err = pluginutil.GetMultiplexIDFromContext(ctx)
		if err != nil {
			return &pb.SetupReply{}, err
		}
	}

	// Dial for storage
	brokeredClient, err := b.broker.Dial(args.BrokerID)
	if err != nil {
		return &pb.SetupReply{}, err
	}

	storage := newGRPCStorageClient(brokeredClient)
	sysView := newGRPCSystemView(brokeredClient)

	config := &logical.BackendConfig{
		StorageView: storage,
		Logger:      b.logger,
		System:      sysView,
		Config:      args.Config,
		BackendUUID: args.BackendUUID,
	}

	// Call the underlying backend factory after shims have been created
	// to set b.backend
	backend, err := b.factory(ctx, config)
	if err != nil {
		return &pb.SetupReply{
			Err: pb.ErrToString(err),
		}, nil
	}
	b.instances[id] = backendInstance{
		brokeredClient: brokeredClient,
		backend:        backend,
	}

	return &pb.SetupReply{}, nil
}

func (b *backendGRPCPluginServer) HandleRequest(ctx context.Context, args *pb.HandleRequestArgs) (*pb.HandleRequestReply, error) {
	backend, err := b.getBackend(ctx)
	if err != nil {
		return &pb.HandleRequestReply{}, err
	}

	if pluginutil.InMetadataMode() {
		return &pb.HandleRequestReply{}, ErrServerInMetadataMode
	}

	logicalReq, err := pb.ProtoRequestToLogicalRequest(args.Request)
	if err != nil {
		return &pb.HandleRequestReply{}, err
	}

	brokeredClient, err := b.getBrokeredClient(ctx)
	if err != nil {
		return &pb.HandleRequestReply{}, err
	}

	logicalReq.Storage = newGRPCStorageClient(brokeredClient)

	resp, respErr := backend.HandleRequest(ctx, logicalReq)

	pbResp, err := pb.LogicalResponseToProtoResponse(resp)
	if err != nil {
		return &pb.HandleRequestReply{}, err
	}

	return &pb.HandleRequestReply{
		Response: pbResp,
		Err:      pb.ErrToProtoErr(respErr),
	}, nil
}

func (b *backendGRPCPluginServer) Initialize(ctx context.Context, _ *pb.InitializeArgs) (*pb.InitializeReply, error) {
	backend, err := b.getBackend(ctx)
	if err != nil {
		return &pb.InitializeReply{}, err
	}

	if pluginutil.InMetadataMode() {
		return &pb.InitializeReply{}, ErrServerInMetadataMode
	}

	brokeredClient, err := b.getBrokeredClient(ctx)
	if err != nil {
		return &pb.InitializeReply{}, err
	}

	req := &logical.InitializationRequest{
		Storage: newGRPCStorageClient(brokeredClient),
	}

	respErr := backend.Initialize(ctx, req)

	return &pb.InitializeReply{
		Err: pb.ErrToProtoErr(respErr),
	}, nil
}

func (b *backendGRPCPluginServer) SpecialPaths(ctx context.Context, args *pb.Empty) (*pb.SpecialPathsReply, error) {
	backend, err := b.getBackend(ctx)
	if err != nil {
		return &pb.SpecialPathsReply{}, err
	}

	paths := backend.SpecialPaths()
	if paths == nil {
		return &pb.SpecialPathsReply{
			Paths: nil,
		}, nil
	}

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
	backend, err := b.getBackend(ctx)
	if err != nil {
		return &pb.HandleExistenceCheckReply{}, err
	}

	if pluginutil.InMetadataMode() {
		return &pb.HandleExistenceCheckReply{}, ErrServerInMetadataMode
	}

	logicalReq, err := pb.ProtoRequestToLogicalRequest(args.Request)
	if err != nil {
		return &pb.HandleExistenceCheckReply{}, err
	}

	brokeredClient, err := b.getBrokeredClient(ctx)
	if err != nil {
		return &pb.HandleExistenceCheckReply{}, err
	}

	logicalReq.Storage = newGRPCStorageClient(brokeredClient)

	checkFound, exists, err := backend.HandleExistenceCheck(ctx, logicalReq)
	return &pb.HandleExistenceCheckReply{
		CheckFound: checkFound,
		Exists:     exists,
		Err:        pb.ErrToProtoErr(err),
	}, nil
}

func (b *backendGRPCPluginServer) Cleanup(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	b.Lock()
	defer b.Unlock()

	backend, err := b.getBackendInternal(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}

	backend.Cleanup(ctx)

	brokeredClient, err := b.getBrokeredClientInternal(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}

	// Close rpc clients
	brokeredClient.Close()

	if _, ok := b.instances["single"]; ok {
		delete(b.instances, "single")
	} else {
		id, err := pluginutil.GetMultiplexIDFromContext(ctx)
		if err != nil {
			return nil, err
		}
		delete(b.instances, id)
	}

	return &pb.Empty{}, nil
}

func (b *backendGRPCPluginServer) InvalidateKey(ctx context.Context, args *pb.InvalidateKeyArgs) (*pb.Empty, error) {
	backend, err := b.getBackend(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}

	if pluginutil.InMetadataMode() {
		return &pb.Empty{}, ErrServerInMetadataMode
	}

	backend.InvalidateKey(ctx, args.Key)
	return &pb.Empty{}, nil
}

func (b *backendGRPCPluginServer) Type(ctx context.Context, _ *pb.Empty) (*pb.TypeReply, error) {
	backend, err := b.getBackend(ctx)
	if err != nil {
		return &pb.TypeReply{}, err
	}

	return &pb.TypeReply{
		Type: uint32(backend.Type()),
	}, nil
}
