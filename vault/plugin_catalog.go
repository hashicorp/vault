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

	// multiplexedClients holds plugin process connections by plugin name
	// This allows a single grpc connection to communicate with multiple
	// databases. Each database configuration using the same plugin will be
	// routed to the existing plugin process.
	multiplexedClients map[string]*MultiplexedClient

	lock sync.RWMutex
}

type PluginClient struct {
	logger   log.Logger
	id       string
	protocol plugin.ClientProtocol
	// client handles the lifecycle of a plugin process
	client *plugin.Client
}

type MultiplexedClient struct {
	logger log.Logger

	// name is the plugin name
	name string

	client *plugin.Client

	connections         map[string]*PluginClient
	multiplexingSupport bool
}

func (p *PluginClient) Conn() *grpc.ClientConn {
	gc := p.protocol.(*plugin.GRPCClient)
	return gc.Conn
}

func (p *PluginClient) ID() string {
	return p.id
}

func (p *MultiplexedClient) ByID(id string) *PluginClient {
	return p.connections[id]
}

func (p *PluginClient) MultiplexingSupport() bool {
	if p.client == nil {
		return false
	}
	return pluginutil.MultiplexingSupport(p.client.NegotiatedVersion())
}

func (c *PluginCatalog) removeMultiplexedClient(name, id string) {
	mpc, ok := c.multiplexedClients[name]
	if !ok {
		return
	}

	mpc.connections[id].Close()

	delete(mpc.connections, id)
	if len(mpc.connections) == 0 {
		delete(c.multiplexedClients, name)
	}
}

func (p *PluginClient) Close() error {
	p.logger.Debug("attempting to kill plugin process", "id", p.id)
	p.client.Kill()
	err := p.protocol.Close()
	p.client = nil
	p.protocol = nil
	p.logger.Debug("killed plugin process", "id", p.id)
	return err
}

func (p *PluginClient) Dispense(name string) (interface{}, error) {
	pluginInstance, err := p.protocol.Dispense(name)
	if err != nil {
		return nil, err
	}
	return pluginInstance, nil
}

func (p *PluginClient) Ping() error {
	err := p.protocol.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) setupPluginCatalog(ctx context.Context) error {
	c.pluginCatalog = &PluginCatalog{
		builtinRegistry: c.builtinRegistry,
		catalogView:     NewBarrierView(c.barrier, pluginCatalogPath),
		directory:       c.pluginDirectory,
		logger:          c.logger,
	}

	// Run upgrade if untyped plugins exist
	err := c.pluginCatalog.UpgradePlugins(ctx, c.logger)
	if err != nil {
		c.logger.Error("error while upgrading plugin storage", "error", err)
	}

	if c.logger.IsInfo() {
		c.logger.Info("successfully setup plugin catalog", "plugin-directory", c.pluginDirectory)
	}

	return nil
}

func (c *PluginCatalog) getMultiplexedClient(pluginName string) *MultiplexedClient {
	if mpc, ok := c.multiplexedClients[pluginName]; ok {
		c.logger.Debug("MultiplexedClient exists", "pluginName", pluginName)

		return mpc
	}

	c.logger.Debug("MultiplexedClient does not exist", "pluginName", pluginName)

	return c.newMultiplexedClient(pluginName)
}

func (c *PluginCatalog) newMultiplexedClient(pluginName string) *MultiplexedClient {
	if c.multiplexedClients == nil {
		c.multiplexedClients = make(map[string]*MultiplexedClient)
		c.logger.Debug("created multiplexedClients map")
	}

	mpc := &MultiplexedClient{connections: make(map[string]*PluginClient), logger: c.logger}

	// set the MultiplexedClient for the given plugin name
	c.multiplexedClients[pluginName] = mpc
	c.logger.Debug("set the MultiplexedClient for", "pluginName", pluginName)

	return mpc
}

// GetPluginClient returns a client for managing the lifecycle of a plugin
// process
func (c *PluginCatalog) GetPluginClient(ctx context.Context, sys pluginutil.RunnerUtil, pluginRunner *pluginutil.PluginRunner, config pluginutil.PluginClientConfig) (*PluginClient, error) {
	c.lock.Lock()
	pc, err := c.getPluginClient(ctx, sys, pluginRunner, config)
	c.lock.Unlock()
	return pc, err
}

// getPluginClient returns a client for managing the lifecycle of a plugin
// process
func (c *PluginCatalog) getPluginClient(ctx context.Context, sys pluginutil.RunnerUtil, pluginRunner *pluginutil.PluginRunner, config pluginutil.PluginClientConfig) (*PluginClient, error) {
	mpc := c.getMultiplexedClient(pluginRunner.Name)
	c.logger.Debug("getPluginClient", "pluginRunner.MultiplexingSupport", pluginRunner.MultiplexingSupport)

	id, err := base62.Random(10)
	if err != nil {
		return nil, err
	}

	pc := &PluginClient{
		id:     id,
		logger: c.logger,
	}
	if !pluginRunner.MultiplexingSupport || mpc.client == nil {
		// get a new client
		c.logger.Debug("spawning a new plugin process")
		client, err := pluginRunner.RunConfig(ctx,
			pluginutil.Runner(sys),
			pluginutil.PluginSets(config.PluginSets),
			pluginutil.HandshakeConfig(config.HandshakeConfig),
			pluginutil.Logger(config.Logger),
			pluginutil.MetadataMode(config.IsMetadataMode),
			pluginutil.AutoMTLS(config.AutoMTLS),
		)
		if err != nil {
			return nil, err
		}

		mpc.client = client
		pc.client = client

	}

	if pluginRunner.MultiplexingSupport && mpc.client != nil {
		// return existing client
		c.logger.Debug("return existing client")

		pc.client = mpc.client
	}

	// Get the protocol client for this connection.
	// Subsequent calls to this will return the same client.
	rpcClient, err := pc.client.Client()
	if err != nil {
		return nil, err
	}

	pc.protocol = rpcClient
	mpc.connections[id] = pc
	mpc.multiplexingSupport = pc.MultiplexingSupport()

	mpc.name = pluginRunner.Name

	return mpc.connections[id], nil
}

// getPluginTypeFromUnknown will attempt to run the plugin to determine the
// type. It will first attempt to run as a database plugin then a backend
// plugin. Both of these will be run in metadata mode.
func (c *PluginCatalog) getPluginTypeFromUnknown(ctx context.Context, logger log.Logger, plugin *pluginutil.PluginRunner) (consts.PluginType, bool, error) {
	merr := &multierror.Error{}
	multiplexingSupport, err := c.isDatabasePlugin(ctx, plugin)
	if err == nil {
		return consts.PluginTypeDatabase, multiplexingSupport, nil
	}
	merr = multierror.Append(merr, err)

	// Attempt to run as backend plugin
	client, err := backendplugin.NewPluginClient(ctx, nil, plugin, log.NewNullLogger(), true)
	if err == nil {
		err := client.Setup(ctx, &logical.BackendConfig{})
		if err != nil {
			return consts.PluginTypeUnknown, false, err
		}

		backendType := client.Type()
		client.Cleanup(ctx)

		switch backendType {
		case logical.TypeCredential:
			return consts.PluginTypeCredential, false, nil
		case logical.TypeLogical:
			return consts.PluginTypeSecrets, false, nil
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

	return consts.PluginTypeUnknown, false, nil
}

func (c *PluginCatalog) isDatabasePlugin(ctx context.Context, pluginRunner *pluginutil.PluginRunner) (bool, error) {
	merr := &multierror.Error{}
	config := pluginutil.PluginClientConfig{
		Name:            pluginRunner.Name,
		PluginSets:      v5.PluginSets,
		HandshakeConfig: v5.HandshakeConfig,
		Logger:          log.NewNullLogger(),
		IsMetadataMode:  true,
		AutoMTLS:        true,
	}
	// Attempt to run as database V5 or V6 plugin
	v5Client, err := c.getPluginClient(ctx, nil, pluginRunner, config)
	if err == nil {
		multiplexingSupport := v5Client.MultiplexingSupport()
		// Close the client and cleanup the plugin process
		c.removeMultiplexedClient(pluginRunner.Name, v5Client.ID())
		return multiplexingSupport, nil
	}
	merr = multierror.Append(merr, fmt.Errorf("failed to load plugin as database v5: %w", err))

	v4Client, err := v4.NewPluginClient(ctx, nil, pluginRunner, log.NewNullLogger(), true)
	if err == nil {
		// Close the client and cleanup the plugin process
		v4Client.Close()
		return false, nil
	}
	merr = multierror.Append(merr, fmt.Errorf("failed to load plugin as database v4: %w", err))

	return false, merr.ErrorOrNil()
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

		pluginType, multiplexingSupport, err := c.getPluginTypeFromUnknown(ctx, logger, plugin)
		if err != nil {
			retErr = multierror.Append(retErr, fmt.Errorf("could not upgrade plugin %s: %s", pluginName, err))
			continue
		}
		if pluginType == consts.PluginTypeUnknown {
			retErr = multierror.Append(retErr, fmt.Errorf("could not upgrade plugin %s: plugin of unknown type", pluginName))
			continue
		}

		// Upgrade the storage
		err = c.setInternal(ctx, pluginName, pluginType, multiplexingSupport, cmdOld, plugin.Args, plugin.Env, plugin.Sha256)
		if err != nil {
			retErr = multierror.Append(retErr, fmt.Errorf("could not upgrade plugin %s: %s", pluginName, err))
			continue
		}

		err = c.catalogView.Delete(ctx, pluginName)
		if err != nil {
			logger.Error("could not remove plugin", "plugin", pluginName, "error", err)
		}

		logger.Info("upgraded plugin type", "plugin", pluginName, "type", pluginType.String())
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

	return c.setInternal(ctx, name, pluginType, false, command, args, env, sha256)
}

func (c *PluginCatalog) setInternal(ctx context.Context, name string, pluginType consts.PluginType, multiplexingSupport bool, command string, args []string, env []string, sha256 []byte) error {
	// Best effort check to make sure the command isn't breaking out of the
	// configured plugin directory.
	commandFull := filepath.Join(c.directory, command)
	sym, err := filepath.EvalSymlinks(commandFull)
	if err != nil {
		return fmt.Errorf("error while validating the command path: %w", err)
	}
	symAbs, err := filepath.Abs(filepath.Dir(sym))
	if err != nil {
		return fmt.Errorf("error while validating the command path: %w", err)
	}

	if symAbs != c.directory {
		return errors.New("cannot execute files outside of configured plugin directory")
	}

	// If the plugin type is unknown, we want to attempt to determine the type
	if pluginType == consts.PluginTypeUnknown {
		// entryTmp should only be used for the below type check, it uses the
		// full command instead of the relative command.
		entryTmp := &pluginutil.PluginRunner{
			Name:                name,
			Command:             commandFull,
			Args:                args,
			Env:                 env,
			Sha256:              sha256,
			Builtin:             false,
			MultiplexingSupport: multiplexingSupport,
		}

		pluginType, multiplexingSupport, err = c.getPluginTypeFromUnknown(ctx, log.Default(), entryTmp)
		if err != nil {
			return err
		}
		if pluginType == consts.PluginTypeUnknown {
			return ErrPluginBadType
		}
	}

	entry := &pluginutil.PluginRunner{
		Name:                name,
		Type:                pluginType,
		Command:             command,
		Args:                args,
		Env:                 env,
		Sha256:              sha256,
		Builtin:             false,
		MultiplexingSupport: multiplexingSupport,
	}

	buf, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to encode plugin entry: %w", err)
	}

	logicalEntry := logical.StorageEntry{
		Key:   pluginType.String() + "/" + name,
		Value: buf,
	}
	if err := c.catalogView.Put(ctx, &logicalEntry); err != nil {
		return fmt.Errorf("failed to persist plugin entry: %w", err)
	}
	return nil
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
