package plugin

import (
	"errors"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/logical"
	log "github.com/mgutz/logxi/v1"
)

var (
	ErrClientInMetadataMode = errors.New("plugin client can not perform action while in metadata mode")
)

// backendPluginClient implements logical.Backend and is the
// go-plugin client.
type backendPluginClient struct {
	broker       *plugin.MuxBroker
	client       *rpc.Client
	metadataMode bool

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
	Error    error
}

// SpecialPathsReply is the reply for SpecialPaths method.
type SpecialPathsReply struct {
	Paths *logical.Paths
}

// SystemReply is the reply for System method.
type SystemReply struct {
	SystemView logical.SystemView
	Error      error
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
	Error      error
}

// SetupArgs is the args for Setup method.
type SetupArgs struct {
	StorageID uint32
	LoggerID  uint32
	SysViewID uint32
	Config    map[string]string
}

// SetupReply is the reply for Setup method.
type SetupReply struct {
	Error error
}

// TypeReply is the reply for the Type method.
type TypeReply struct {
	Type logical.BackendType
}

// RegisterLicenseArgs is the args for the RegisterLicense method.
type RegisterLicenseArgs struct {
	License interface{}
}

// RegisterLicenseReply is the reply for the RegisterLicense method.
type RegisterLicenseReply struct {
	Error error
}

func (b *backendPluginClient) HandleRequest(req *logical.Request) (*logical.Response, error) {
	if b.metadataMode {
		return nil, ErrClientInMetadataMode
	}

	// Do not send the storage, since go-plugin cannot serialize
	// interfaces. The server will pick up the storage from the shim.
	req.Storage = nil
	args := &HandleRequestArgs{
		Request: req,
	}
	var reply HandleRequestReply

	if req.Connection != nil {
		oldConnState := req.Connection.ConnState
		req.Connection.ConnState = nil
		defer func() {
			req.Connection.ConnState = oldConnState
		}()
	}

	err := b.client.Call("Plugin.HandleRequest", args, &reply)
	if err != nil {
		return nil, err
	}
	if reply.Error != nil {
		if reply.Error.Error() == logical.ErrUnsupportedOperation.Error() {
			return nil, logical.ErrUnsupportedOperation
		}

		return reply.Response, reply.Error
	}

	return reply.Response, nil
}

func (b *backendPluginClient) SpecialPaths() *logical.Paths {
	var reply SpecialPathsReply
	err := b.client.Call("Plugin.SpecialPaths", new(interface{}), &reply)
	if err != nil {
		return nil
	}

	return reply.Paths
}

// System returns vault's system view. The backend client stores the view during
// Setup, so there is no need to shim the system just to get it back.
func (b *backendPluginClient) System() logical.SystemView {
	return b.system
}

// Logger returns vault's logger. The backend client stores the logger during
// Setup, so there is no need to shim the logger just to get it back.
func (b *backendPluginClient) Logger() log.Logger {
	return b.logger
}

func (b *backendPluginClient) HandleExistenceCheck(req *logical.Request) (bool, bool, error) {
	if b.metadataMode {
		return false, false, ErrClientInMetadataMode
	}

	// Do not send the storage, since go-plugin cannot serialize
	// interfaces. The server will pick up the storage from the shim.
	req.Storage = nil
	args := &HandleExistenceCheckArgs{
		Request: req,
	}
	var reply HandleExistenceCheckReply

	if req.Connection != nil {
		oldConnState := req.Connection.ConnState
		req.Connection.ConnState = nil
		defer func() {
			req.Connection.ConnState = oldConnState
		}()
	}

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
	if b.metadataMode {
		return ErrClientInMetadataMode
	}
	err := b.client.Call("Plugin.Initialize", new(interface{}), &struct{}{})
	return err
}

func (b *backendPluginClient) InvalidateKey(key string) {
	if b.metadataMode {
		return
	}
	b.client.Call("Plugin.InvalidateKey", key, &struct{}{})
}

func (b *backendPluginClient) Setup(config *logical.BackendConfig) error {
	// Shim logical.Storage
	storageImpl := config.StorageView
	if b.metadataMode {
		storageImpl = &NOOPStorage{}
	}
	storageID := b.broker.NextId()
	go b.broker.AcceptAndServe(storageID, &StorageServer{
		impl: storageImpl,
	})

	// Shim log.Logger
	loggerImpl := config.Logger
	if b.metadataMode {
		loggerImpl = log.NullLog
	}
	loggerID := b.broker.NextId()
	go b.broker.AcceptAndServe(loggerID, &LoggerServer{
		logger: loggerImpl,
	})

	// Shim logical.SystemView
	sysViewImpl := config.System
	if b.metadataMode {
		sysViewImpl = &logical.StaticSystemView{}
	}
	sysViewID := b.broker.NextId()
	go b.broker.AcceptAndServe(sysViewID, &SystemViewServer{
		impl: sysViewImpl,
	})

	args := &SetupArgs{
		StorageID: storageID,
		LoggerID:  loggerID,
		SysViewID: sysViewID,
		Config:    config.Config,
	}
	var reply SetupReply

	err := b.client.Call("Plugin.Setup", args, &reply)
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

func (b *backendPluginClient) Type() logical.BackendType {
	var reply TypeReply
	err := b.client.Call("Plugin.Type", new(interface{}), &reply)
	if err != nil {
		return logical.TypeUnknown
	}

	return logical.BackendType(reply.Type)
}

func (b *backendPluginClient) RegisterLicense(license interface{}) error {
	if b.metadataMode {
		return ErrClientInMetadataMode
	}

	var reply RegisterLicenseReply
	args := RegisterLicenseArgs{
		License: license,
	}
	err := b.client.Call("Plugin.RegisterLicense", args, &reply)
	if err != nil {
		return err
	}
	if reply.Error != nil {
		return reply.Error
	}

	return nil
}
