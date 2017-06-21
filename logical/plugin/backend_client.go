package plugin

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/logical"
	log "github.com/mgutz/logxi/v1"
)

// backendPluginClient implements logical.Backend and is the
// go-plugin client.
type backendPluginClient struct {
	broker       *plugin.MuxBroker
	client       *rpc.Client
	pluginClient *plugin.Client

	system logical.SystemView
	logger log.Logger
}

// HandleRequestArgs is the args for HandleRequest method.
type HandleRequestArgs struct {
	StorageID uint32
	Request   *logical.Request
}

// HandleRequestReply is the reply for HandleRequest method.
type HandleRequestReply struct {
	Response *logical.Response
	Error    *plugin.BasicError
}

// SpecialPathsReply is the reply for SpecialPaths method.
type SpecialPathsReply struct {
	Paths *logical.Paths
}

// SystemReply is the reply for System method.
type SystemReply struct {
	SystemView logical.SystemView
	Error      *plugin.BasicError
}

// HandleExistenceCheckArgs is the args for HandleExistenceCheck method.
type HandleExistenceCheckArgs struct {
	StorageID uint32
	Request   *logical.Request
}

// HandleExistenceCheckReply is the reply for HandleExistenceCheck method.
type HandleExistenceCheckReply struct {
	CheckFound bool
	Exists     bool
	Error      *plugin.BasicError
}

// ConfigureArgs is the args for Configure method.
type ConfigureArgs struct {
	StorageID uint32
	LoggerID  uint32
	SysViewID uint32
	Config    map[string]string
}

// ConfigureReply is the reply for Configure method.
type ConfigureReply struct {
	Error *plugin.BasicError
}

func (b *backendPluginClient) HandleRequest(req *logical.Request) (*logical.Response, error) {
	// Shim logical.Storage
	id := b.broker.NextId()
	go b.broker.AcceptAndServe(id, &StorageServer{
		impl: req.Storage,
	})

	argReq := new(logical.Request)
	*argReq = *req
	argReq.Storage = nil

	args := &HandleRequestArgs{
		StorageID: id,
		Request:   argReq,
	}
	var reply HandleRequestReply

	err := b.client.Call("Plugin.HandleRequest", args, &reply)
	if err != nil {
		return nil, err
	}
	if reply.Error != nil {
		if reply.Error.Error() == logical.ErrUnsupportedOperation.Error() {
			return nil, logical.ErrUnsupportedOperation
		}
		return nil, reply.Error
	}

	return reply.Response, nil
}

func (b *backendPluginClient) SpecialPaths() *logical.Paths {
	var reply SpecialPathsReply
	// var paths logical.Paths
	err := b.client.Call("Plugin.SpecialPaths", new(interface{}), &reply)
	if err != nil {
		return nil
	}

	return reply.Paths
}

// System returns vault's system view. The backend client stores the view during
// Configure, so there is no need to shim the system just to get it back.
func (b *backendPluginClient) System() logical.SystemView {
	return b.system
}

// Logger returns vault's logger. The backend client stores the logger during
// Configure, so there is no need to shim the logger just to get it back.
func (b *backendPluginClient) Logger() log.Logger {
	return b.logger
}

func (b *backendPluginClient) HandleExistenceCheck(req *logical.Request) (bool, bool, error) {
	// Shim logical.Storage
	id := b.broker.NextId()
	go b.broker.AcceptAndServe(id, &StorageServer{
		impl: req.Storage,
	})

	argReq := new(logical.Request)
	*argReq = *req
	argReq.Storage = nil

	args := &HandleExistenceCheckArgs{
		StorageID: id,
		Request:   argReq,
	}
	var reply HandleExistenceCheckReply

	err := b.client.Call("Plugin.HandleExistenceCheck", args, &reply)
	if err != nil {
		return false, false, err
	}
	if reply.Error != nil {
		// THINKING: Should be be a switch on all error types?
		if reply.Error.Error() == logical.ErrUnsupportedPath.Error() {
			return false, false, logical.ErrUnsupportedPath
		}
		return false, false, reply.Error
	}

	return reply.CheckFound, reply.Exists, nil
}

func (b *backendPluginClient) Cleanup() {
	b.client.Call("Plugin.Cleanup", new(interface{}), &struct{}{})
}

func (b *backendPluginClient) Initialize() error {
	err := b.client.Call("Plugin.Initialize", new(interface{}), &struct{}{})
	return err
}

func (b *backendPluginClient) InvalidateKey(key string) {
	b.client.Call("Plugin.InvalidateKey", key, &struct{}{})
}

func (b *backendPluginClient) Configure(config *logical.BackendConfig) error {
	// Shim logical.Storage
	storageID := b.broker.NextId()
	go b.broker.AcceptAndServe(storageID, &StorageServer{
		impl: config.StorageView,
	})

	// Shim log.Logger
	loggerID := b.broker.NextId()
	go b.broker.AcceptAndServe(loggerID, &LoggerServer{
		logger: config.Logger,
	})

	// Shim logical.SystemView
	sysViewID := b.broker.NextId()
	go b.broker.AcceptAndServe(sysViewID, &SystemViewServer{
		impl: config.System,
	})

	args := &ConfigureArgs{
		StorageID: storageID,
		LoggerID:  loggerID,
		SysViewID: sysViewID,
		Config:    config.Config,
	}
	var reply ConfigureReply

	err := b.client.Call("Plugin.Configure", args, &reply)
	if err != nil {
		return err
	}
	if reply.Error != nil {
		return reply.Error
	}

	// Set system and logger for getter methods
	b.system = config.System
	b.logger = config.Logger

	return nil
}
