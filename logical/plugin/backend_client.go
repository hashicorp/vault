package plugin

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/logical"
	log "github.com/mgutz/logxi/v1"
)

type backendPluginClient struct {
	broker *plugin.MuxBroker
	client *rpc.Client

	system logical.SystemView
	logger log.Logger
}

type HandleRequestArgs struct {
	StorageID uint32
	Request   *logical.Request
}

type HandleRequestReply struct {
	Response *logical.Response
	Error    *plugin.BasicError
}

type SpecialPathsReply struct {
	Paths *logical.Paths
}

type SystemReply struct {
	SystemView logical.SystemView
	Error      *plugin.BasicError
}

type HandleExistenceCheckArgs struct {
	StorageID uint32
	Request   *logical.Request
}

type BackendLoggerReply struct {
	Logger log.Logger
}

type HandleExistenceCheckReply struct {
	CheckFound bool
	Exists     bool
	Error      *plugin.BasicError
}

type ConfigureArgs struct {
	StorageID uint32
	LoggerID  uint32
	SysViewID uint32
	Config    map[string]string
}

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

func (b *backendPluginClient) System() logical.SystemView {
	return b.system
}

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
	// This should no-op for plugins, since Configure() gets called instead of Initialize()
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
