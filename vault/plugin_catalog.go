package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	log "github.com/hashicorp/go-hclog"
	multierror "github.com/hashicorp/go-multierror"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-secure-stdlib/base62"
	v4 "github.com/hashicorp/vault/sdk/database/dbplugin"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	backendplugin "github.com/hashicorp/vault/sdk/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	pluginCatalogPath         = "core/plugin-catalog/"
	ErrDirectoryNotConfigured = errors.New("could not set plugin, plugin directory is not configured")
	ErrPluginNotFound         = errors.New("plugin not found in the catalog")
	ErrPluginBadType          = errors.New("unable to determine plugin type")
)

// PluginCatalog keeps a record of plugins known to vault. External plugins need
// to be registered to the catalog before they can be used in backends. Builtin
// plugins are automatically detected and included in the catalog.
type PluginCatalog struct {
	builtinRegistry BuiltinRegistry
	catalogView     *BarrierView
	directory       string
	logger          log.Logger

	// externalPlugins holds plugin process connections by plugin name
	//
	// This allows plugins that suppport multiplexing to use a single grpc
	// connection to communicate with multiple "backends". Each backend
	// configuration using the same plugin will be routed to the existing
	// plugin process.
	externalPlugins map[string]*externalPlugin
	mlockPlugins    bool

	lock sync.RWMutex
}

// externalPlugin holds client connections for multiplexed and
// non-multiplexed plugin processes
type externalPlugin struct {
	// name is the plugin name
	name string

	// connections holds client connections by ID
	connections map[string]*pluginClient

	multiplexingSupport bool
}

// pluginClient represents a connection to a plugin process
type pluginClient struct {
	logger log.Logger

	// id is the connection ID
	id string

	// client handles the lifecycle of a plugin process
	// multiplexed plugins share the same client
	client      *plugin.Client
	clientConn  grpc.ClientConnInterface
	cleanupFunc func() error

	plugin.ClientProtocol
}

func wrapFactoryCheckPerms(core *Core, f logical.Factory) logical.Factory {
	return func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
		if err := core.CheckPluginPerms(conf.Config["plugin_name"]); err != nil {
			return nil, err
		}
		return f(ctx, conf)
	}
}

func (c *Core) setupPluginCatalog(ctx context.Context) error {
	c.pluginCatalog = &PluginCatalog{
		builtinRegistry: c.builtinRegistry,
		catalogView:     NewBarrierView(c.barrier, pluginCatalogPath),
		directory:       c.pluginDirectory,
		logger:          c.logger,
		mlockPlugins:    c.enableMlock,
	}

	// Run upgrade if untyped plugins exist
	err := c.pluginCatalog.UpgradePlugins(ctx, c.logger)
	if err != nil {
		c.logger.Error("error while upgrading plugin storage", "error", err)
		return err
	}

	if c.logger.IsInfo() {
		c.logger.Info("successfully setup plugin catalog", "plugin-directory", c.pluginDirectory)
	}

	return nil
}

type pluginClientConn struct {
	*grpc.ClientConn
	id string
}

var _ grpc.ClientConnInterface = &pluginClientConn{}

func (d *pluginClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	// Inject ID to the context
	md := metadata.Pairs(pluginutil.MultiplexingCtxKey, d.id)
	idCtx := metadata.NewOutgoingContext(ctx, md)

	return d.ClientConn.Invoke(idCtx, method, args, reply, opts...)
}

func (d *pluginClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// Inject ID to the context
	md := metadata.Pairs(pluginutil.MultiplexingCtxKey, d.id)
	idCtx := metadata.NewOutgoingContext(ctx, md)

	return d.ClientConn.NewStream(idCtx, desc, method, opts...)
}

func (p *pluginClient) Conn() grpc.ClientConnInterface {
	return p.clientConn
}

// Close calls the plugin client's cleanupFunc to do any necessary cleanup on
// the plugin client and the PluginCatalog. This implements the
// plugin.ClientProtocol interface.
func (p *pluginClient) Close() error {
	p.logger.Debug("cleaning up plugin client connection", "id", p.id)
	return p.cleanupFunc()
}

// cleanupExternalPlugin will kill plugin processes and perform any necessary
// cleanup on the externalPlugins map for multiplexed and non-multiplexed
// plugins. This should be called with the write lock held.
func (c *PluginCatalog) cleanupExternalPlugin(name, id string) error {
	extPlugin, ok := c.externalPlugins[name]
	if !ok {
		return fmt.Errorf("plugin client not found")
	}

	pc, ok := extPlugin.connections[id]
	if !ok {
		return fmt.Errorf("plugin connection not found")
	}

	delete(extPlugin.connections, id)
	if !extPlugin.multiplexingSupport {
		pc.client.Kill()

		if len(extPlugin.connections) == 0 {
			delete(c.externalPlugins, name)
		}
	} else if len(extPlugin.connections) == 0 || pc.client.Exited() {
		pc.client.Kill()
		delete(c.externalPlugins, name)
	}

	return nil
}

func (c *PluginCatalog) getExternalPlugin(pluginName string) *externalPlugin {
	if extPlugin, ok := c.externalPlugins[pluginName]; ok {
		return extPlugin
	}

	return c.newExternalPlugin(pluginName)
}

func (c *PluginCatalog) newExternalPlugin(pluginName string) *externalPlugin {
	if c.externalPlugins == nil {
		c.externalPlugins = make(map[string]*externalPlugin)
	}

	extPlugin := &externalPlugin{
		connections: make(map[string]*pluginClient),
		name:        pluginName,
	}

	c.externalPlugins[pluginName] = extPlugin
	return extPlugin
}

// NewPluginClient returns a client for managing the lifecycle of a plugin
// process
func (c *PluginCatalog) NewPluginClient(ctx context.Context, config pluginutil.PluginClientConfig) (*pluginClient, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if config.Name == "" {
		return nil, fmt.Errorf("no name provided for plugin")
	}
	if config.PluginType == consts.PluginTypeUnknown {
		return nil, fmt.Errorf("no plugin type provided")
	}

	pluginRunner, err := c.get(ctx, config.Name, config.PluginType)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup plugin: %w", err)
	}
	if pluginRunner == nil {
		return nil, fmt.Errorf("no plugin found")
	}
	pc, err := c.newPluginClient(ctx, pluginRunner, config)
	return pc, err
}

// newPluginClient returns a client for managing the lifecycle of a plugin
// process. Callers should have the write lock held.
func (c *PluginCatalog) newPluginClient(ctx context.Context, pluginRunner *pluginutil.PluginRunner, config pluginutil.PluginClientConfig) (*pluginClient, error) {
	if pluginRunner == nil {
		return nil, fmt.Errorf("no plugin found")
	}

	extPlugin := c.getExternalPlugin(pluginRunner.Name)
	id, err := base62.Random(10)
	if err != nil {
		return nil, err
	}

	pc := &pluginClient{
		id:     id,
		logger: c.logger.Named(pluginRunner.Name),
		cleanupFunc: func() error {
			c.lock.Lock()
			defer c.lock.Unlock()
			return c.cleanupExternalPlugin(pluginRunner.Name, id)
		},
	}

	// Multiplexing support will always be false initially, but will be
	// adjusted once we query from the plugin whether it can multiplex or not
	if !extPlugin.multiplexingSupport || len(extPlugin.connections) == 0 {
		c.logger.Debug("spawning a new plugin process", "plugin_name", pluginRunner.Name, "id", id)
		client, err := pluginRunner.RunConfig(ctx,
			pluginutil.PluginSets(config.PluginSets),
			pluginutil.HandshakeConfig(config.HandshakeConfig),
			pluginutil.Logger(config.Logger),
			pluginutil.MetadataMode(config.IsMetadataMode),
			pluginutil.MLock(c.mlockPlugins),

			// NewPluginClient only supports AutoMTLS today
			pluginutil.AutoMTLS(true),
		)
		if err != nil {
			return nil, err
		}

		pc.client = client
	} else {
		c.logger.Debug("returning existing plugin client for multiplexed plugin", "id", id)

		// get the first client, since they are all the same
		for k := range extPlugin.connections {
			pc.client = extPlugin.connections[k].client
			break
		}

		if pc.client == nil {
			return nil, fmt.Errorf("plugin client is nil")
		}
	}

	// Get the protocol client for this connection.
	// Subsequent calls to this will return the same client.
	rpcClient, err := pc.client.Client()
	if err != nil {
		return nil, err
	}

	clientConn := rpcClient.(*plugin.GRPCClient).Conn

	muxed, err := pluginutil.MultiplexingSupported(ctx, clientConn)
	if err != nil {
		return nil, err
	}

	if muxed {
		// Wrap rpcClient with our implementation so that we can inject the
		// ID into the context
		pc.clientConn = &pluginClientConn{
			ClientConn: clientConn,
			id:         id,
		}
	} else {
		pc.clientConn = clientConn
	}

	pc.ClientProtocol = rpcClient

	extPlugin.connections[id] = pc
	extPlugin.name = pluginRunner.Name
	extPlugin.multiplexingSupport = muxed

	return extPlugin.connections[id], nil
}

// getPluginTypeFromUnknown will attempt to run the plugin to determine the
// type and if it supports multiplexing. It will first attempt to run as a
// database plugin then a backend plugin. Both of these will be run in metadata
// mode.
func (c *PluginCatalog) getPluginTypeFromUnknown(ctx context.Context, logger log.Logger, plugin *pluginutil.PluginRunner) (consts.PluginType, error) {
	merr := &multierror.Error{}
	err := c.isDatabasePlugin(ctx, plugin)
	if err == nil {
		return consts.PluginTypeDatabase, nil
	}
	merr = multierror.Append(merr, err)

	// Attempt to run as backend plugin
	client, err := backendplugin.NewPluginClient(ctx, nil, plugin, log.NewNullLogger(), true)
	if err == nil {
		err := client.Setup(ctx, &logical.BackendConfig{})
		if err != nil {
			return consts.PluginTypeUnknown, err
		}

		backendType := client.Type()
		client.Cleanup(ctx)

		switch backendType {
		case logical.TypeCredential:
			return consts.PluginTypeCredential, nil
		case logical.TypeLogical:
			return consts.PluginTypeSecrets, nil
		}
	} else {
		merr = multierror.Append(merr, err)
	}

	if client == nil || client.Type() == logical.TypeUnknown {
		logger.Warn("unknown plugin type",
			"plugin name", plugin.Name,
			"error", merr.Error())
	} else {
		logger.Warn("unsupported plugin type",
			"plugin name", plugin.Name,
			"plugin type", client.Type().String(),
			"error", merr.Error())
	}

	return consts.PluginTypeUnknown, nil
}

// isDatabasePlugin returns true if the plugin supports multiplexing. An error
// is returned if the plugin is not a database plugin.
func (c *PluginCatalog) isDatabasePlugin(ctx context.Context, pluginRunner *pluginutil.PluginRunner) error {
	merr := &multierror.Error{}
	config := pluginutil.PluginClientConfig{
		Name:            pluginRunner.Name,
		PluginSets:      v5.PluginSets,
		PluginType:      consts.PluginTypeDatabase,
		HandshakeConfig: v5.HandshakeConfig,
		Logger:          log.NewNullLogger(),
		IsMetadataMode:  true,
		AutoMTLS:        true,
	}

	// Attempt to run as database V5 or V6 multiplexed plugin
	c.logger.Debug("attempting to load database plugin as v5", "name", pluginRunner.Name)
	v5Client, err := c.newPluginClient(ctx, pluginRunner, config)
	if err == nil {
		// Close the client and cleanup the plugin process
		err = c.cleanupExternalPlugin(pluginRunner.Name, v5Client.id)
		if err != nil {
			c.logger.Error("error closing plugin client", "error", err)
		}

		return nil
	}
	merr = multierror.Append(merr, fmt.Errorf("failed to load plugin as database v5: %w", err))

	c.logger.Debug("attempting to load database plugin as v4", "name", pluginRunner.Name)
	v4Client, err := v4.NewPluginClient(ctx, nil, pluginRunner, log.NewNullLogger(), true)
	if err == nil {
		// Close the client and cleanup the plugin process
		err = v4Client.Close()
		if err != nil {
			c.logger.Error("error closing plugin client", "error", err)
		}

		return nil
	}
	merr = multierror.Append(merr, fmt.Errorf("failed to load plugin as database v4: %w", err))

	return merr.ErrorOrNil()
}

// UpdatePlugins will loop over all the plugins of unknown type and attempt to
// upgrade them to typed plugins
func (c *PluginCatalog) UpgradePlugins(ctx context.Context, logger log.Logger) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// If the directory isn't set we can skip the upgrade attempt
	if c.directory == "" {
		return nil
	}

	// List plugins from old location
	pluginsRaw, err := c.catalogView.List(ctx, "")
	if err != nil {
		return err
	}
	plugins := make([]string, 0, len(pluginsRaw))
	for _, p := range pluginsRaw {
		if !strings.HasSuffix(p, "/") {
			plugins = append(plugins, p)
		}
	}

	logger.Info("upgrading plugin information", "plugins", plugins)

	var retErr error
	for _, pluginName := range plugins {
		pluginRaw, err := c.catalogView.Get(ctx, pluginName)
		if err != nil {
			retErr = multierror.Append(fmt.Errorf("failed to load plugin entry: %w", err))
			continue
		}

		plugin := new(pluginutil.PluginRunner)
		if err := jsonutil.DecodeJSON(pluginRaw.Value, plugin); err != nil {
			retErr = multierror.Append(fmt.Errorf("failed to decode plugin entry: %w", err))
			continue
		}

		// prepend the plugin directory to the command
		cmdOld := plugin.Command
		plugin.Command = filepath.Join(c.directory, plugin.Command)

		// Upgrade the storage. At this point we don't know what type of plugin this is so pass in the unkonwn type.
		runner, err := c.setInternal(ctx, pluginName, consts.PluginTypeUnknown, cmdOld, plugin.Args, plugin.Env, plugin.Sha256)
		if err != nil {
			if errors.Is(err, ErrPluginBadType) {
				retErr = multierror.Append(retErr, fmt.Errorf("could not upgrade plugin %s: plugin of unknown type", pluginName))
				continue
			}

			retErr = multierror.Append(retErr, fmt.Errorf("could not upgrade plugin %s: %s", pluginName, err))
			continue
		}

		err = c.catalogView.Delete(ctx, pluginName)
		if err != nil {
			logger.Error("could not remove plugin", "plugin", pluginName, "error", err)
		}

		logger.Info("upgraded plugin type", "plugin", pluginName, "type", runner.Type.String())
	}

	return retErr
}

// Get retrieves a plugin with the specified name from the catalog. It first
// looks for external plugins with this name and then looks for builtin plugins.
// It returns a PluginRunner or an error if no plugin was found.
func (c *PluginCatalog) Get(ctx context.Context, name string, pluginType consts.PluginType) (*pluginutil.PluginRunner, error) {
	c.lock.RLock()
	runner, err := c.get(ctx, name, pluginType)
	c.lock.RUnlock()
	return runner, err
}

func (c *PluginCatalog) get(ctx context.Context, name string, pluginType consts.PluginType) (*pluginutil.PluginRunner, error) {
	// If the directory isn't set only look for builtin plugins.
	if c.directory != "" {
		// Look for external plugins in the barrier
		out, err := c.catalogView.Get(ctx, pluginType.String()+"/"+name)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve plugin %q: %w", name, err)
		}
		if out == nil {
			// Also look for external plugins under what their name would have been if they
			// were registered before plugin types existed.
			out, err = c.catalogView.Get(ctx, name)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve plugin %q: %w", name, err)
			}
		}
		if out != nil {
			entry := new(pluginutil.PluginRunner)
			if err := jsonutil.DecodeJSON(out.Value, entry); err != nil {
				return nil, fmt.Errorf("failed to decode plugin entry: %w", err)
			}
			if entry.Type != pluginType && entry.Type != consts.PluginTypeUnknown {
				return nil, nil
			}

			// prepend the plugin directory to the command
			entry.Command = filepath.Join(c.directory, entry.Command)

			return entry, nil
		}
	}
	// Look for builtin plugins
	if factory, ok := c.builtinRegistry.Get(name, pluginType); ok {
		return &pluginutil.PluginRunner{
			Name:           name,
			Type:           pluginType,
			Builtin:        true,
			BuiltinFactory: factory,
		}, nil
	}

	return nil, nil
}

// Set registers a new external plugin with the catalog, or updates an existing
// external plugin. It takes the name, command and SHA256 of the plugin.
func (c *PluginCatalog) Set(ctx context.Context, name string, pluginType consts.PluginType, command string, args []string, env []string, sha256 []byte) error {
	if c.directory == "" {
		return ErrDirectoryNotConfigured
	}

	switch {
	case strings.Contains(name, ".."):
		fallthrough
	case strings.Contains(command, ".."):
		return consts.ErrPathContainsParentReferences
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.setInternal(ctx, name, pluginType, command, args, env, sha256)
	return err
}

func (c *PluginCatalog) setInternal(ctx context.Context, name string, pluginType consts.PluginType, command string, args []string, env []string, sha256 []byte) (*pluginutil.PluginRunner, error) {
	// Best effort check to make sure the command isn't breaking out of the
	// configured plugin directory.
	commandFull := filepath.Join(c.directory, command)
	sym, err := filepath.EvalSymlinks(commandFull)
	if err != nil {
		return nil, fmt.Errorf("error while validating the command path: %w", err)
	}
	symAbs, err := filepath.Abs(filepath.Dir(sym))
	if err != nil {
		return nil, fmt.Errorf("error while validating the command path: %w", err)
	}

	if symAbs != c.directory {
		return nil, errors.New("cannot execute files outside of configured plugin directory")
	}

	// If the plugin type is unknown, we want to attempt to determine the type
	if pluginType == consts.PluginTypeUnknown {
		// entryTmp should only be used for the below type check, it uses the
		// full command instead of the relative command.
		entryTmp := &pluginutil.PluginRunner{
			Name:    name,
			Command: commandFull,
			Args:    args,
			Env:     env,
			Sha256:  sha256,
			Builtin: false,
		}

		pluginType, err = c.getPluginTypeFromUnknown(ctx, log.Default(), entryTmp)
		if err != nil {
			return nil, err
		}
		if pluginType == consts.PluginTypeUnknown {
			return nil, ErrPluginBadType
		}
	}

	entry := &pluginutil.PluginRunner{
		Name:    name,
		Type:    pluginType,
		Command: command,
		Args:    args,
		Env:     env,
		Sha256:  sha256,
		Builtin: false,
	}

	buf, err := json.Marshal(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to encode plugin entry: %w", err)
	}

	logicalEntry := logical.StorageEntry{
		Key:   pluginType.String() + "/" + name,
		Value: buf,
	}
	if err := c.catalogView.Put(ctx, &logicalEntry); err != nil {
		return nil, fmt.Errorf("failed to persist plugin entry: %w", err)
	}
	return entry, nil
}

// Delete is used to remove an external plugin from the catalog. Builtin plugins
// can not be deleted.
func (c *PluginCatalog) Delete(ctx context.Context, name string, pluginType consts.PluginType) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check the name under which the plugin exists, but if it's unfound, don't return any error.
	pluginKey := pluginType.String() + "/" + name
	out, err := c.catalogView.Get(ctx, pluginKey)
	if err != nil || out == nil {
		pluginKey = name
	}

	return c.catalogView.Delete(ctx, pluginKey)
}

// List returns a list of all the known plugin names. If an external and builtin
// plugin share the same name, only one instance of the name will be returned.
func (c *PluginCatalog) List(ctx context.Context, pluginType consts.PluginType) ([]string, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	// Collect keys for external plugins in the barrier.
	keys, err := logical.CollectKeys(ctx, c.catalogView)
	if err != nil {
		return nil, err
	}

	// Get the builtin plugins.
	builtinKeys := c.builtinRegistry.Keys(pluginType)

	// Use a map to unique the two lists.
	mapKeys := make(map[string]bool)

	pluginTypePrefix := pluginType.String() + "/"

	for _, plugin := range keys {
		// Only list user-added plugins if they're of the given type.
		if entry, err := c.get(ctx, plugin, pluginType); err == nil && entry != nil {

			// Some keys will be prepended with the plugin type, but other ones won't.
			// Users don't expect to see the plugin type, so we need to strip that here.
			idx := strings.Index(plugin, pluginTypePrefix)
			if idx == 0 {
				plugin = plugin[len(pluginTypePrefix):]
			}
			mapKeys[plugin] = true
		}
	}

	for _, plugin := range builtinKeys {
		mapKeys[plugin] = true
	}

	retList := make([]string, len(mapKeys))
	i := 0
	for k := range mapKeys {
		retList[i] = k
		i++
	}
	// sort for consistent ordering of builtin plugins
	sort.Strings(retList)

	return retList, nil
}
