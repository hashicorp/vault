package plugin

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/logical"
)

// backendPluginServer is the RPC server that backendPluginClient talks to,
// it methods conforming to requirements by net/rpc
type backendPluginServer struct {
	broker  *plugin.MuxBroker
	backend logical.Backend

	loggerClient  *rpc.Client
	sysViewClient *rpc.Client
	storageClient *rpc.Client
}

func (b *backendPluginServer) HandleRequest(args *HandleRequestArgs, reply *HandleRequestReply) error {
	conn, err := b.broker.Dial(args.StorageID)
	if err != nil {
		*reply = HandleRequestReply{
			Error: plugin.NewBasicError(err),
		}
		return nil
	}
	c := rpc.NewClient(conn)
	defer c.Close()

	storage := &StorageClient{client: c}
	args.Request.Storage = storage

	resp, err := b.backend.HandleRequest(args.Request)
	*reply = HandleRequestReply{
		Response: resp,
		Error:    plugin.NewBasicError(err),
	}

	return nil
}

func (b *backendPluginServer) SpecialPaths(_ interface{}, reply *SpecialPathsReply) error {
	*reply = SpecialPathsReply{
		Paths: b.backend.SpecialPaths(),
	}
	return nil
}

func (b *backendPluginServer) Logger(_ interface{}, reply *BackendLoggerReply) error {
	*reply = BackendLoggerReply{
		Logger: b.backend.Logger(),
	}
	return nil
}

func (b *backendPluginServer) HandleExistenceCheck(args *HandleExistenceCheckArgs, reply *HandleExistenceCheckReply) error {
	conn, err := b.broker.Dial(args.StorageID)
	if err != nil {
		*reply = HandleExistenceCheckReply{
			Error: plugin.NewBasicError(err),
		}
		return nil
	}
	c := rpc.NewClient(conn)
	defer c.Close()

	storage := &StorageClient{client: c}
	args.Request.Storage = storage

	checkFound, exists, err := b.backend.HandleExistenceCheck(args.Request)
	*reply = HandleExistenceCheckReply{
		CheckFound: checkFound,
		Exists:     exists,
		Error:      plugin.NewBasicError(err),
	}

	return nil
}

func (b *backendPluginServer) Cleanup(_ interface{}, _ *struct{}) error {
	b.backend.Cleanup()

	// Close rpc clients
	b.loggerClient.Close()
	b.sysViewClient.Close()
	b.storageClient.Close()
	return nil
}

func (b *backendPluginServer) Initialize(_ interface{}, _ *struct{}) error {
	err := b.backend.Initialize()
	return err
}

func (b *backendPluginServer) InvalidateKey(args string, _ *struct{}) error {
	b.backend.InvalidateKey(args)
	return nil
}

func (b *backendPluginServer) Configure(args *ConfigureArgs, reply *ConfigureReply) error {
	// Dial for storage
	storageConn, err := b.broker.Dial(args.StorageID)
	if err != nil {
		*reply = ConfigureReply{
			Error: plugin.NewBasicError(err),
		}
		return nil
	}
	rawStorageClient := rpc.NewClient(storageConn)
	b.storageClient = rawStorageClient

	storage := &StorageClient{client: rawStorageClient}

	// Dial for logger
	loggerConn, err := b.broker.Dial(args.LoggerID)
	if err != nil {
		*reply = ConfigureReply{
			Error: plugin.NewBasicError(err),
		}
		return nil
	}
	rawLoggerClient := rpc.NewClient(loggerConn)
	b.loggerClient = rawLoggerClient

	logger := &LoggerClient{client: rawLoggerClient}

	// Dial for sys view
	sysViewConn, err := b.broker.Dial(args.SysViewID)
	if err != nil {
		*reply = ConfigureReply{
			Error: plugin.NewBasicError(err),
		}
		return nil
	}
	rawSysViewClient := rpc.NewClient(sysViewConn)
	b.sysViewClient = rawSysViewClient

	sysView := &SystemViewClient{client: rawSysViewClient}

	config := &logical.BackendConfig{
		StorageView: storage,
		Logger:      logger,
		System:      sysView,
		Config:      args.Config,
	}

	err = b.backend.Configure(config)
	if err != nil {
		*reply = ConfigureReply{
			Error: plugin.NewBasicError(err),
		}
	}

	return nil
}
