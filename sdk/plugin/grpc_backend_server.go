// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"errors"
	"fmt"
	"sync"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
	"google.golang.org/grpc"
)

var ErrServerInMetadataMode = errors.New("plugin server can not perform action while in metadata mode")

// singleImplementationID is the string used to define the instance ID of a
// non-multiplexed plugin
const singleImplementationID string = "single"

type backendInstance struct {
	brokeredClient *grpc.ClientConn
	backend        logical.Backend
}

type backendGRPCPluginServer struct {
	pb.UnimplementedBackendServer
	logical.UnimplementedPluginVersionServer

	broker *plugin.GRPCBroker

	instances           map[string]backendInstance
	instancesLock       sync.RWMutex
	multiplexingSupport bool

	factory logical.Factory

	logger log.Logger
}

// getBackendAndBrokeredClientInternal returns the backend and client
// connection but does not hold a lock
func (b *backendGRPCPluginServer) getBackendAndBrokeredClientInternal(ctx context.Context) (logical.Backend, *grpc.ClientConn, error) {
	if b.multiplexingSupport {
		id, err := pluginutil.GetMultiplexIDFromContext(ctx)
		if err != nil {
			return nil, nil, err
		}

		if inst, ok := b.instances[id]; ok {
			return inst.backend, inst.brokeredClient, nil
		}

	}

	if singleImpl, ok := b.instances[singleImplementationID]; ok {
		return singleImpl.backend, singleImpl.brokeredClient, nil
	}

	return nil, nil, fmt.Errorf("no backend instance found")
}

// getBackendAndBrokeredClient holds a read lock and returns the backend and
// client connection
func (b *backendGRPCPluginServer) getBackendAndBrokeredClient(ctx context.Context) (logical.Backend, *grpc.ClientConn, error) {
	b.instancesLock.RLock()
	defer b.instancesLock.RUnlock()
	return b.getBackendAndBrokeredClientInternal(ctx)
}

// Setup dials into the plugin's broker to get a shimmed storage, logger, and
// system view of the backend. This method also instantiates the underlying
// backend through its factory func for the server side of the plugin.
func (b *backendGRPCPluginServer) Setup(ctx context.Context, args *pb.SetupArgs) (*pb.SetupReply, error) {
	var err error
	id := singleImplementationID

	if b.multiplexingSupport {
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
	events := newGRPCEventsClient(brokeredClient)

	config := &logical.BackendConfig{
		StorageView:  storage,
		Logger:       b.logger,
		System:       newGRPCSystemViewFromSetupArgs(brokeredClient, args),
		Config:       args.Config,
		BackendUUID:  args.BackendUUID,
		EventsSender: events,
	}

	// Call the underlying backend factory after shims have been created
	// to set b.backend
	backend, err := b.factory(ctx, config)
	if err != nil {
		return &pb.SetupReply{
			Err: pb.ErrToString(err),
		}, nil
	}

	b.instancesLock.Lock()
	defer b.instancesLock.Unlock()
	b.instances[id] = backendInstance{
		brokeredClient: brokeredClient,
		backend:        backend,
	}

	return &pb.SetupReply{}, nil
}

func (b *backendGRPCPluginServer) HandleRequest(ctx context.Context, args *pb.HandleRequestArgs) (*pb.HandleRequestReply, error) {
	backend, brokeredClient, err := b.getBackendAndBrokeredClient(ctx)
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
	backend, brokeredClient, err := b.getBackendAndBrokeredClient(ctx)
	if err != nil {
		return &pb.InitializeReply{}, err
	}

	if pluginutil.InMetadataMode() {
		return &pb.InitializeReply{}, ErrServerInMetadataMode
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
	backend, _, err := b.getBackendAndBrokeredClient(ctx)
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
			Root:                  paths.Root,
			Unauthenticated:       paths.Unauthenticated,
			LocalStorage:          paths.LocalStorage,
			SealWrapStorage:       paths.SealWrapStorage,
			WriteForwardedStorage: paths.WriteForwardedStorage,
			Binary:                paths.Binary,
			Limited:               paths.Limited,
		},
	}, nil
}

func (b *backendGRPCPluginServer) HandleExistenceCheck(ctx context.Context, args *pb.HandleExistenceCheckArgs) (*pb.HandleExistenceCheckReply, error) {
	backend, brokeredClient, err := b.getBackendAndBrokeredClient(ctx)
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

	logicalReq.Storage = newGRPCStorageClient(brokeredClient)

	checkFound, exists, err := backend.HandleExistenceCheck(ctx, logicalReq)
	return &pb.HandleExistenceCheckReply{
		CheckFound: checkFound,
		Exists:     exists,
		Err:        pb.ErrToProtoErr(err),
	}, nil
}

func (b *backendGRPCPluginServer) Cleanup(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	b.instancesLock.Lock()
	defer b.instancesLock.Unlock()

	backend, brokeredClient, err := b.getBackendAndBrokeredClientInternal(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}

	backend.Cleanup(ctx)

	// Close rpc clients
	brokeredClient.Close()

	if b.multiplexingSupport {
		id, err := pluginutil.GetMultiplexIDFromContext(ctx)
		if err != nil {
			return nil, err
		}
		delete(b.instances, id)
	} else if _, ok := b.instances[singleImplementationID]; ok {
		delete(b.instances, singleImplementationID)
	}

	return &pb.Empty{}, nil
}

func (b *backendGRPCPluginServer) InvalidateKey(ctx context.Context, args *pb.InvalidateKeyArgs) (*pb.Empty, error) {
	backend, _, err := b.getBackendAndBrokeredClient(ctx)
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
	backend, _, err := b.getBackendAndBrokeredClient(ctx)
	if err != nil {
		return &pb.TypeReply{}, err
	}

	return &pb.TypeReply{
		Type: uint32(backend.Type()),
	}, nil
}

func (b *backendGRPCPluginServer) Version(ctx context.Context, _ *logical.Empty) (*logical.VersionReply, error) {
	backend, _, err := b.getBackendAndBrokeredClient(ctx)
	if err != nil {
		return &logical.VersionReply{}, err
	}

	if versioner, ok := backend.(logical.PluginVersioner); ok {
		return &logical.VersionReply{
			PluginVersion: versioner.PluginVersion().Version,
		}, nil
	}
	return &logical.VersionReply{
		PluginVersion: "",
	}, nil
}
