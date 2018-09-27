package plugin

import (
	"context"
	"errors"
	"net/rpc"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
)

var (
	ErrServerInMetadataMode = errors.New("plugin server can not perform action while in metadata mode")
)

// backendPluginServer is the RPC server that backendPluginClient talks to,
// it methods conforming to requirements by net/rpc
type backendPluginServer struct {
	broker  *plugin.MuxBroker
	backend logical.Backend
	factory logical.Factory

	logger        hclog.Logger
	sysViewClient *rpc.Client
	storageClient *rpc.Client
}

func (b *backendPluginServer) HandleRequest(args *HandleRequestArgs, reply *HandleRequestReply) error {
	if pluginutil.InMetadataMode() {
		return ErrServerInMetadataMode
	}

	storage := &StorageClient{client: b.storageClient}
	args.Request.Storage = storage

	resp, err := b.backend.HandleRequest(context.Background(), args.Request)
	*reply = HandleRequestReply{
		Response: resp,
		Error:    wrapError(err),
	}

	return nil
}

func (b *backendPluginServer) SpecialPaths(_ interface{}, reply *SpecialPathsReply) error {
	*reply = SpecialPathsReply{
		Paths: b.backend.SpecialPaths(),
	}
	return nil
}

func (b *backendPluginServer) HandleExistenceCheck(args *HandleExistenceCheckArgs, reply *HandleExistenceCheckReply) error {
	if pluginutil.InMetadataMode() {
		return ErrServerInMetadataMode
	}

	storage := &StorageClient{client: b.storageClient}
	args.Request.Storage = storage

	checkFound, exists, err := b.backend.HandleExistenceCheck(context.TODO(), args.Request)
	*reply = HandleExistenceCheckReply{
		CheckFound: checkFound,
		Exists:     exists,
		Error:      wrapError(err),
	}

	return nil
}

func (b *backendPluginServer) Cleanup(_ interface{}, _ *struct{}) error {
	b.backend.Cleanup(context.Background())

	// Close rpc clients
	b.sysViewClient.Close()
	b.storageClient.Close()
	return nil
}

func (b *backendPluginServer) InvalidateKey(args string, _ *struct{}) error {
	if pluginutil.InMetadataMode() {
		return ErrServerInMetadataMode
	}

	b.backend.InvalidateKey(context.Background(), args)
	return nil
}

// Setup dials into the plugin's broker to get a shimmed storage, logger, and
// system view of the backend. This method also instantiates the underlying
// backend through its factory func for the server side of the plugin.
func (b *backendPluginServer) Setup(args *SetupArgs, reply *SetupReply) error {
	// Dial for storage
	storageConn, err := b.broker.Dial(args.StorageID)
	if err != nil {
		*reply = SetupReply{
			Error: wrapError(err),
		}
		return nil
	}
	rawStorageClient := rpc.NewClient(storageConn)
	b.storageClient = rawStorageClient

	storage := &StorageClient{client: rawStorageClient}

	// Dial for sys view
	sysViewConn, err := b.broker.Dial(args.SysViewID)
	if err != nil {
		*reply = SetupReply{
			Error: wrapError(err),
		}
		return nil
	}
	rawSysViewClient := rpc.NewClient(sysViewConn)
	b.sysViewClient = rawSysViewClient

	sysView := &SystemViewClient{client: rawSysViewClient}

	config := &logical.BackendConfig{
		StorageView: storage,
		Logger:      b.logger,
		System:      sysView,
		Config:      args.Config,
		BackendUUID: args.BackendUUID,
	}

	// Call the underlying backend factory after shims have been created
	// to set b.backend
	backend, err := b.factory(context.Background(), config)
	if err != nil {
		*reply = SetupReply{
			Error: wrapError(err),
		}
	}
	b.backend = backend

	return nil
}

func (b *backendPluginServer) Type(_ interface{}, reply *TypeReply) error {
	*reply = TypeReply{
		Type: b.backend.Type(),
	}

	return nil
}
