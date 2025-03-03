// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-secure-stdlib/base62"
	semver "github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/helper/versions"
	v4 "github.com/hashicorp/vault/sdk/database/dbplugin"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	backendplugin "github.com/hashicorp/vault/sdk/plugin"
	"github.com/hashicorp/vault/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	ErrDirectoryNotConfigured   = errors.New("could not set plugin, plugin directory is not configured")
	ErrPluginNotFound           = errors.New("plugin not found in the catalog")
	ErrPluginConnectionNotFound = errors.New("plugin connection not found for client")
	ErrPluginBadType            = errors.New("unable to determine plugin type")
	ErrPinnedVersion            = errors.New("cannot delete a pinned version")
	ErrPluginVersionMismatch    = errors.New("plugin version mismatch")
	ErrPluginUnableToRun        = errors.New("unable to run plugin")
)

// PluginCatalog keeps a record of plugins known to vault. External plugins need
// to be registered to the catalog before they can be used in backends. Builtin
// plugins are automatically detected and included in the catalog.
type PluginCatalog struct {
	builtinRegistry BuiltinRegistry
	catalogView     logical.Storage
	directory       string
	tmpdir          string
	logger          log.Logger

	// externalPlugins holds plugin process connections by a key which is
	// generated from the plugin runner config.
	//
	// This allows plugins that suppport multiplexing to use a single grpc
	// connection to communicate with multiple "backends". Each backend
	// configuration using the same plugin will be routed to the existing
	// plugin process.
	externalPlugins map[externalPluginsKey]*externalPlugin
	mlockPlugins    bool

	lock    sync.RWMutex
	wrapper pluginutil.RunnerUtil

	runtimeCatalog *PluginRuntimeCatalog
}

// Only plugins running with identical PluginRunner config can be multiplexed,
// so we use the PluginRunner input as the key for the external plugins map.
//
// However, to be a map key, it must be comparable:
// https://go.dev/ref/spec#Comparison_operators.
// In particular, the PluginRunner struct has slices and a function which are not
// comparable, so we need to transform it into a struct which is.
type externalPluginsKey struct {
	name     string
	typ      consts.PluginType
	version  string
	command  string
	ociImage string
	runtime  string
	args     string
	env      string
	sha256   string
	builtin  bool
}

func makeExternalPluginsKey(p *pluginutil.PluginRunner) (externalPluginsKey, error) {
	args, err := json.Marshal(p.Args)
	if err != nil {
		return externalPluginsKey{}, err
	}

	env, err := json.Marshal(p.Env)
	if err != nil {
		return externalPluginsKey{}, err
	}

	return externalPluginsKey{
		name:     p.Name,
		typ:      p.Type,
		version:  p.Version,
		command:  p.Command,
		ociImage: p.OCIImage,
		runtime:  p.Runtime,
		args:     string(args),
		env:      string(env),
		sha256:   hex.EncodeToString(p.Sha256),
		builtin:  p.Builtin,
	}, nil
}

// externalPlugin holds client connections for multiplexed and
// non-multiplexed plugin processes
type externalPlugin struct {
	// connections holds client connections by ID
	connections map[string]*pluginClient

	multiplexingSupport bool
}

// pluginClient represents a connection to a plugin process
type pluginClient struct {
	logger log.Logger

	// id is the connection ID
	id       string
	pluginID string

	// client handles the lifecycle of a plugin process
	// multiplexed plugins share the same client
	client      *plugin.Client
	clientConn  grpc.ClientConnInterface
	cleanupFunc func() error
	reloadFunc  func() error

	plugin.ClientProtocol
}

type PluginCatalogInput struct {
	Logger               log.Logger
	BuiltinRegistry      BuiltinRegistry
	CatalogView          logical.Storage
	PluginDirectory      string
	Tmpdir               string
	EnableMlock          bool
	PluginRuntimeCatalog *PluginRuntimeCatalog
}

func SetupPluginCatalog(ctx context.Context, in *PluginCatalogInput) (*PluginCatalog, error) {
	logger := in.Logger
	catalog := &PluginCatalog{
		builtinRegistry: in.BuiltinRegistry,
		catalogView:     in.CatalogView,
		directory:       in.PluginDirectory,
		tmpdir:          in.Tmpdir,
		logger:          logger,
		mlockPlugins:    in.EnableMlock,
		wrapper:         logical.StaticSystemView{VersionString: version.GetVersion().Version},
		runtimeCatalog:  in.PluginRuntimeCatalog,
	}

	// Run upgrade if untyped plugins exist
	err := catalog.upgradePlugins(ctx, logger)
	if err != nil {
		logger.Error("error while upgrading plugin storage", "error", err)
		return nil, err
	}

	// Sanitize the plugin catalog
	err = catalog.entValidate(ctx)
	if err != nil {
		logger.Error("error while sanitizing plugin storage", "error", err)
		return nil, err
	}

	if legacy, _ := strconv.ParseBool(os.Getenv(pluginutil.PluginUseLegacyEnvLayering)); legacy {
		conflicts := false
		osKeys := envKeys(os.Environ())
		plugins, err := catalog.collectAllPlugins(ctx)
		if err != nil {
			return nil, err
		}
		for _, plugin := range plugins {
			pluginKeys := envKeys(plugin.Env)
			for k := range pluginKeys {
				if _, ok := osKeys[k]; ok {
					conflicts = true
					logger.Warn("conflict between system and plugin environment variable", "type", plugin.Type, "name", plugin.Name, "version", plugin.Version, "variable", k)
				}
			}
		}
		if conflicts {
			logger.Warn("conflicts found between system and plugin environment variables, "+
				"system environment variables will take precedence until flag is disabled",
				pluginutil.PluginUseLegacyEnvLayering, os.Getenv(pluginutil.PluginUseLegacyEnvLayering))
		} else {
			logger.Info("no conflicts found between system and plugin environment variables")
		}
	}

	logger.Info("successfully setup plugin catalog", "plugin-directory", catalog.directory)
	if catalog.tmpdir != "" {
		logger.Debug("plugin temporary directory configured", "tmpdir", catalog.tmpdir)
	}

	return catalog, nil
}

func envKeys(env []string) map[string]struct{} {
	keys := make(map[string]struct{}, len(env))
	for _, env := range env {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 0 {
			continue
		}
		keys[parts[0]] = struct{}{}
	}
	return keys
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

func (p *pluginClient) Reload() error {
	p.logger.Debug("reload external plugin process")
	return p.reloadFunc()
}

func (c *PluginCatalog) Processes() int {
	return len(c.externalPlugins)
}

// reloadExternalPlugin
// This should be called with the write lock held.
func (c *PluginCatalog) reloadExternalPlugin(key externalPluginsKey, id, pluginBinaryRef string) error {
	extPlugin, ok := c.externalPlugins[key]
	if !ok {
		return fmt.Errorf("plugin client not found")
	}
	if !extPlugin.multiplexingSupport {
		err := c.cleanupExternalPlugin(key, id, pluginBinaryRef)
		if err != nil {
			return err
		}
		return nil
	}

	pc, ok := extPlugin.connections[id]
	if !ok {
		return fmt.Errorf("%w id: %s", ErrPluginConnectionNotFound, id)
	}

	delete(c.externalPlugins, key)
	pc.client.Kill()
	c.logger.Debug("killed external plugin process for reload", "plugin", pluginBinaryRef, "pluginID", pc.pluginID)

	return nil
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
func (c *PluginCatalog) cleanupExternalPlugin(key externalPluginsKey, id, pluginBinaryRef string) error {
	extPlugin, ok := c.externalPlugins[key]
	if !ok {
		return fmt.Errorf("plugin client not found")
	}

	pc, ok := extPlugin.connections[id]
	if !ok {
		// this can happen if the backend is reloaded due to a plugin process
		// being killed out of band
		c.logger.Warn(ErrPluginConnectionNotFound.Error(), "id", id)
		return fmt.Errorf("%w id: %s", ErrPluginConnectionNotFound, id)
	}

	delete(extPlugin.connections, id)
	c.logger.Debug("removed plugin client connection", "id", id)

	if !extPlugin.multiplexingSupport {
		pc.client.Kill()

		if len(extPlugin.connections) == 0 {
			delete(c.externalPlugins, key)
		}
		c.logger.Debug("killed external plugin process", "plugin", pluginBinaryRef, "pluginID", pc.pluginID)
	} else if len(extPlugin.connections) == 0 || pc.client.Exited() {
		pc.client.Kill()
		delete(c.externalPlugins, key)
		c.logger.Debug("killed external multiplexed plugin process", "plugin", pluginBinaryRef, "pluginID", pc.pluginID)
	}

	return nil
}

func (c *PluginCatalog) getExternalPlugin(key externalPluginsKey) *externalPlugin {
	if extPlugin, ok := c.externalPlugins[key]; ok {
		return extPlugin
	}

	return c.newExternalPlugin(key)
}

func (c *PluginCatalog) newExternalPlugin(key externalPluginsKey) *externalPlugin {
	if c.externalPlugins == nil {
		c.externalPlugins = make(map[externalPluginsKey]*externalPlugin)
	}

	extPlugin := &externalPlugin{
		connections: make(map[string]*pluginClient),
	}

	c.externalPlugins[key] = extPlugin
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

	pluginRunner, err := c.get(ctx, config.Name, config.PluginType, config.Version)
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

	key, err := makeExternalPluginsKey(pluginRunner)
	if err != nil {
		return nil, err
	}

	extPlugin := c.getExternalPlugin(key)
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
			return c.cleanupExternalPlugin(key, id, pluginRunner.BinaryReference())
		},
		reloadFunc: func() error {
			c.lock.Lock()
			defer c.lock.Unlock()
			return c.reloadExternalPlugin(key, id, pluginRunner.BinaryReference())
		},
	}

	// Multiplexing support will always be false initially, but will be
	// adjusted once we query from the plugin whether it can multiplex or not
	var spawnedPlugin bool
	if !extPlugin.multiplexingSupport || len(extPlugin.connections) == 0 {
		c.logger.Debug("spawning a new plugin process", "plugin_name", pluginRunner.Name, "id", id)
		client, err := pluginRunner.RunConfig(ctx,
			pluginutil.PluginSets(config.PluginSets),
			pluginutil.HandshakeConfig(config.HandshakeConfig),
			pluginutil.Logger(config.Logger),
			pluginutil.MetadataMode(config.IsMetadataMode),
			pluginutil.MLock(c.mlockPlugins),
			pluginutil.AutoMTLS(config.AutoMTLS),
			pluginutil.Runner(config.Wrapper),
		)
		if err != nil {
			return nil, err
		}

		spawnedPlugin = true
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
		// Make sure we kill any spawned plugins that didn't make it into our
		// map of connections.
		if spawnedPlugin {
			pc.client.Kill()
		}
		return nil, err
	}

	// get the external plugin id
	pc.pluginID = pc.client.ID()

	clientConn := rpcClient.(*plugin.GRPCClient).Conn

	muxed, err := pluginutil.MultiplexingSupported(ctx, clientConn, config.Name)
	if err != nil {
		// Make sure we kill any spawned plugins that didn't make it into our
		// map of connections.
		if spawnedPlugin {
			pc.client.Kill()
		}
		return nil, err
	}

	pc.clientConn = &pluginClientConn{
		ClientConn: clientConn,
		id:         id,
	}

	pc.ClientProtocol = rpcClient

	extPlugin.connections[id] = pc
	extPlugin.multiplexingSupport = muxed

	return extPlugin.connections[id], nil
}

// getPluginTypeFromUnknown will attempt to run the plugin to determine the
// type. It will first attempt to run as a database plugin then a backend
// plugin.
func (c *PluginCatalog) getPluginTypeFromUnknown(ctx context.Context, plugin *pluginutil.PluginRunner) (consts.PluginType, error) {
	merr := &multierror.Error{}
	err := c.isDatabasePlugin(ctx, plugin)
	if err == nil {
		return consts.PluginTypeDatabase, nil
	}
	merr = multierror.Append(merr, err)

	pluginType, err := c.getBackendPluginType(ctx, plugin)
	if err == nil {
		return pluginType, nil
	}
	merr = multierror.Append(merr, err)

	return consts.PluginTypeUnknown, merr
}

// getBackendPluginType returns an error if the plugin is not a backend plugin.
func (c *PluginCatalog) getBackendPluginType(ctx context.Context, pluginRunner *pluginutil.PluginRunner) (consts.PluginType, error) {
	merr := &multierror.Error{}
	// Attempt to run as backend plugin
	config := pluginutil.PluginClientConfig{
		Name:            pluginRunner.Name,
		PluginSets:      backendplugin.PluginSet,
		HandshakeConfig: backendplugin.HandshakeConfig,
		Logger:          log.NewNullLogger(),
		IsMetadataMode:  false,
		AutoMTLS:        true,
		Wrapper:         c.wrapper,
	}

	var client logical.Backend
	var attemptV4 bool
	// First, attempt to run as backend V5 plugin
	c.logger.Debug("attempting to load backend plugin", "name", pluginRunner.Name)
	pc, err := c.newPluginClient(ctx, pluginRunner, config)
	if err == nil {
		// we spawned a subprocess, so make sure to clean it up
		key, err := makeExternalPluginsKey(pluginRunner)
		if err != nil {
			return consts.PluginTypeUnknown, err
		}
		defer func() {
			// Close the client and cleanup the plugin process
			err = c.cleanupExternalPlugin(key, pc.id, pluginRunner.BinaryReference())
			if err != nil {
				c.logger.Error("error closing plugin client", "error", err)
			}
		}()

		// dispense the plugin so we can get its type
		client, err = backendplugin.Dispense(pc.ClientProtocol, pc)
		if err != nil {
			merr = multierror.Append(merr, fmt.Errorf("failed to dispense plugin as backend v5: %w", err))
			c.logger.Debug("failed to dispense v5 backend plugin", "name", pluginRunner.Name)
			attemptV4 = true
		} else {
			c.logger.Debug("successfully dispensed v5 backend plugin", "name", pluginRunner.Name)
		}
	} else {
		attemptV4 = true
	}

	if attemptV4 {
		c.logger.Debug("failed to dispense v5 backend plugin", "name", pluginRunner.Name)
		config.AutoMTLS = false
		config.IsMetadataMode = true
		// attempt to run as a v4 backend plugin
		client, err = backendplugin.NewPluginClient(ctx, c.wrapper, pluginRunner, log.NewNullLogger(), true)
		if err != nil {
			merr = multierror.Append(merr, fmt.Errorf("failed to dispense v4 backend plugin: %w", err))
			c.logger.Debug("failed to dispense v4 backend plugin", "name", pluginRunner.Name, "error", merr)
			return consts.PluginTypeUnknown, merr.ErrorOrNil()
		}
		c.logger.Debug("successfully dispensed v4 backend plugin", "name", pluginRunner.Name)
		defer client.Cleanup(ctx)
	}

	err = client.Setup(ctx, &logical.BackendConfig{})
	if err != nil {
		return consts.PluginTypeUnknown, err
	}
	backendType := client.Type()

	switch backendType {
	case logical.TypeCredential:
		return consts.PluginTypeCredential, nil
	case logical.TypeLogical:
		return consts.PluginTypeSecrets, nil
	}

	if client == nil || client.Type() == logical.TypeUnknown {
		c.logger.Warn("unknown plugin type",
			"plugin name", pluginRunner.Name,
			"error", merr.Error())
	} else {
		c.logger.Warn("unsupported plugin type",
			"plugin name", pluginRunner.Name,
			"plugin type", client.Type().String(),
			"error", merr.Error())
	}

	merr = multierror.Append(merr, fmt.Errorf("failed to load plugin as backend plugin: %w", err))

	return consts.PluginTypeUnknown, merr.ErrorOrNil()
}

// getBackendRunningVersion attempts to get the plugin version
func (c *PluginCatalog) getBackendRunningVersion(ctx context.Context, pluginRunner *pluginutil.PluginRunner) (logical.PluginVersion, error) {
	merr := &multierror.Error{}
	// Attempt to run as backend plugin
	config := pluginutil.PluginClientConfig{
		Name:            pluginRunner.Name,
		PluginSets:      backendplugin.PluginSet,
		HandshakeConfig: backendplugin.HandshakeConfig,
		Logger:          log.NewNullLogger(),
		IsMetadataMode:  false,
		AutoMTLS:        true,
		Wrapper:         c.wrapper,
	}

	var client logical.Backend
	// First, attempt to run as backend V5 plugin
	c.logger.Debug("attempting to load backend plugin", "name", pluginRunner.Name)
	pc, err := c.newPluginClient(ctx, pluginRunner, config)
	if err == nil {
		// we spawned a subprocess, so make sure to clean it up
		key, err := makeExternalPluginsKey(pluginRunner)
		if err != nil {
			return logical.EmptyPluginVersion, err
		}
		defer func() {
			// Close the client and cleanup the plugin process
			err = c.cleanupExternalPlugin(key, pc.id, pluginRunner.BinaryReference())
			if err != nil {
				c.logger.Error("error closing plugin client", "error", err)
			}
		}()

		// dispense the plugin so we can get its version
		client, err = backendplugin.Dispense(pc.ClientProtocol, pc)
		if err == nil {
			c.logger.Debug("successfully dispensed v5 backend plugin", "name", pluginRunner.Name)

			err = client.Setup(ctx, &logical.BackendConfig{})
			if err != nil {
				return logical.EmptyPluginVersion, nil
			}
			if versioner, ok := client.(logical.PluginVersioner); ok {
				return versioner.PluginVersion(), nil
			}
			return logical.EmptyPluginVersion, nil
		}
		merr = multierror.Append(merr, err)
	}
	c.logger.Debug("failed to dispense v5 backend plugin", "name", pluginRunner.Name, "error", err)
	config.AutoMTLS = false
	config.IsMetadataMode = true
	// attempt to run as a v4 backend plugin
	client, err = backendplugin.NewPluginClient(ctx, c.wrapper, pluginRunner, log.NewNullLogger(), true)
	if err != nil {
		merr = multierror.Append(merr, err)
		c.logger.Debug("failed to dispense v4 backend plugin", "name", pluginRunner.Name, "error", err)
		if pluginRunner.OCIImage != "" {
			return logical.EmptyPluginVersion, fmt.Errorf("%w: %w", ErrPluginUnableToRun, merr.ErrorOrNil())
		}
		return logical.EmptyPluginVersion, merr.ErrorOrNil()
	}
	c.logger.Debug("successfully dispensed v4 backend plugin", "name", pluginRunner.Name)
	defer client.Cleanup(ctx)

	err = client.Setup(ctx, &logical.BackendConfig{})
	if err != nil {
		return logical.EmptyPluginVersion, err
	}
	if versioner, ok := client.(logical.PluginVersioner); ok {
		return versioner.PluginVersion(), nil
	}
	return logical.EmptyPluginVersion, nil
}

// getDatabaseRunningVersion returns the version reported by a database plugin
func (c *PluginCatalog) getDatabaseRunningVersion(ctx context.Context, pluginRunner *pluginutil.PluginRunner) (logical.PluginVersion, error) {
	merr := &multierror.Error{}
	config := pluginutil.PluginClientConfig{
		Name:            pluginRunner.Name,
		PluginSets:      v5.PluginSets,
		PluginType:      consts.PluginTypeDatabase,
		Version:         pluginRunner.Version,
		HandshakeConfig: v5.HandshakeConfig,
		Logger:          log.Default(),
		IsMetadataMode:  true,
		AutoMTLS:        true,
		Wrapper:         c.wrapper,
	}

	// Attempt to run as database V5+ multiplexed plugin
	c.logger.Debug("attempting to load database plugin as v5", "name", pluginRunner.Name)
	v5Client, err := c.newPluginClient(ctx, pluginRunner, config)
	if err == nil {
		key, err := makeExternalPluginsKey(pluginRunner)
		if err != nil {
			return logical.EmptyPluginVersion, err
		}
		defer func() {
			// Close the client and cleanup the plugin process
			err = c.cleanupExternalPlugin(key, v5Client.id, pluginRunner.BinaryReference())
			if err != nil {
				c.logger.Error("error closing plugin client", "error", err)
			}
		}()

		raw, err := v5Client.Dispense("database")
		if err != nil {
			return logical.EmptyPluginVersion, err
		}
		if versioner, ok := raw.(logical.PluginVersioner); ok {
			return versioner.PluginVersion(), nil
		}
		return logical.EmptyPluginVersion, nil
	}
	merr = multierror.Append(merr, err)

	c.logger.Debug("attempting to load database plugin as v4", "name", pluginRunner.Name)
	v4Client, err := v4.NewPluginClient(ctx, c.wrapper, pluginRunner, log.NewNullLogger(), true)
	if err == nil {
		// Close the client and cleanup the plugin process
		defer func() {
			err = v4Client.Close()
			if err != nil {
				c.logger.Error("error closing plugin client", "error", err)
			}
		}()

		if versioner, ok := v4Client.(logical.PluginVersioner); ok {
			return versioner.PluginVersion(), nil
		}

		return logical.EmptyPluginVersion, nil
	}

	merr = multierror.Append(merr, err)
	if pluginRunner.OCIImage != "" {
		return logical.EmptyPluginVersion, fmt.Errorf("%w: %w", ErrPluginUnableToRun, merr.ErrorOrNil())
	}

	return logical.EmptyPluginVersion, merr.ErrorOrNil()
}

// isDatabasePlugin returns an error if the plugin is not a database plugin.
func (c *PluginCatalog) isDatabasePlugin(ctx context.Context, pluginRunner *pluginutil.PluginRunner) error {
	merr := &multierror.Error{}
	config := pluginutil.PluginClientConfig{
		Name:            pluginRunner.Name,
		PluginSets:      v5.PluginSets,
		PluginType:      consts.PluginTypeDatabase,
		Version:         pluginRunner.Version,
		HandshakeConfig: v5.HandshakeConfig,
		Logger:          log.NewNullLogger(),
		IsMetadataMode:  true,
		AutoMTLS:        true,
		Wrapper:         c.wrapper,
	}

	// Attempt to run as database V5+ multiplexed plugin
	c.logger.Debug("attempting to load database plugin as v5", "name", pluginRunner.Name)
	v5Client, err := c.newPluginClient(ctx, pluginRunner, config)
	if err == nil {
		// Close the client and cleanup the plugin process
		key, err := makeExternalPluginsKey(pluginRunner)
		if err != nil {
			return err
		}
		err = c.cleanupExternalPlugin(key, v5Client.id, pluginRunner.BinaryReference())
		if err != nil {
			c.logger.Error("error closing plugin client", "error", err)
		}

		return nil
	}
	merr = multierror.Append(merr, fmt.Errorf("failed to load plugin as database v5: %w", err))

	c.logger.Debug("attempting to load database plugin as v4", "name", pluginRunner.Name)
	v4Client, err := v4.NewPluginClient(ctx, c.wrapper, pluginRunner, log.NewNullLogger(), true)
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

// upgradePlugins will loop over all the plugins of unknown type and attempt to
// upgrade them to typed plugins
func (c *PluginCatalog) upgradePlugins(ctx context.Context, logger log.Logger) error {
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
	if len(pluginsRaw) == 0 {
		return nil
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

		if pluginRaw == nil {
			retErr = multierror.Append(fmt.Errorf("%q plugin entry was nil", pluginName))
			continue
		}

		plugin := new(pluginutil.PluginRunner)
		if err := jsonutil.DecodeJSON(pluginRaw.Value, plugin); err != nil {
			retErr = multierror.Append(fmt.Errorf("failed to decode plugin entry: %w", err))
			continue
		}

		// Upgrade the storage. At this point we don't know what type of plugin this is so pass in the unknown type.
		runner, err := c.setInternal(ctx, pluginutil.SetPluginInput{
			Name:    pluginName,
			Type:    consts.PluginTypeUnknown,
			Version: plugin.Version,
			Command: plugin.Command,
			Args:    plugin.Args,
			Env:     plugin.Env,
			Sha256:  plugin.Sha256,
		})
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
func (c *PluginCatalog) Get(ctx context.Context, name string, pluginType consts.PluginType, version string) (*pluginutil.PluginRunner, error) {
	c.lock.RLock()
	runner, err := c.get(ctx, name, pluginType, version)
	c.lock.RUnlock()
	return runner, err
}

func (c *PluginCatalog) get(ctx context.Context, name string, pluginType consts.PluginType, version string) (*pluginutil.PluginRunner, error) {
	// Look for external plugins in the barrier
	storageKey := path.Join(pluginType.String(), name)
	if version != "" {
		storageKey = path.Join(storageKey, version)
	}
	out, err := c.catalogView.Get(ctx, storageKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve plugin %q: %w", name, err)
	}
	if out == nil && version == "" {
		// Also look for external plugins under what their name would have been if they
		// were registered before plugin types existed.
		out, err = c.catalogView.Get(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve plugin %q: %w", name, err)
		}
	}
	entry := new(pluginutil.PluginRunner)
	if out != nil {
		if err := jsonutil.DecodeJSON(out.Value, entry); err != nil {
			return nil, fmt.Errorf("failed to decode plugin entry: %w", err)
		}
		if entry.Type != pluginType && entry.Type != consts.PluginTypeUnknown {
			return nil, nil
		}

		// If none of the cases are satisfied, we'll search for a builtin plugin below.
		switch {
		case entry.OCIImage != "":
			entry.Tmpdir = c.tmpdir
			if entry.Runtime != "" {
				entry.RuntimeConfig, err = c.runtimeCatalog.Get(ctx, entry.Runtime, consts.PluginRuntimeTypeContainer)
				if err != nil {
					return nil, fmt.Errorf("failed to get configured runtime for plugin %q: %w", name, err)
				}
			}
			return entry, nil
		case c.directory != "":
			// Only allow returning non-container external plugins if we have a plugin directory.
			// Make the command path fully rooted.
			entry.Command = filepath.Join(c.directory, entry.Command)
			return entry, nil
		}
	}

	builtinVersion := versions.GetBuiltinVersion(pluginType, name)
	if version == "" || version == builtinVersion {
		if version == builtinVersion {
			// Don't return the builtin if it's shadowed by an unversioned plugin.
			unversioned, err := c.get(ctx, name, pluginType, "")
			if err == nil && unversioned != nil && !unversioned.Builtin {
				return nil, nil
			}
		}

		// Look for builtin plugins
		if factory, ok := c.builtinRegistry.Get(name, pluginType); ok {
			return &pluginutil.PluginRunner{
				Name:           name,
				Type:           pluginType,
				Builtin:        true,
				BuiltinFactory: factory,
				Version:        builtinVersion,
			}, nil
		}
	}

	return nil, nil
}

// Set registers a new external plugin with the catalog, or updates an existing
// external plugin. It takes the name, command and SHA256 of the plugin.
func (c *PluginCatalog) Set(ctx context.Context, plugin pluginutil.SetPluginInput) error {
	if c.directory == "" && plugin.OCIImage == "" {
		return ErrDirectoryNotConfigured
	}

	switch {
	case strings.Contains(plugin.Name, ".."):
		fallthrough
	case strings.Contains(plugin.Command, ".."):
		return consts.ErrPathContainsParentReferences
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.setInternal(ctx, plugin)
	return err
}

// Delete is used to remove an external plugin from the catalog. Builtin plugins
// can not be deleted.
func (c *PluginCatalog) Delete(ctx context.Context, name string, pluginType consts.PluginType, pluginVersion string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	pin, err := c.getPinnedVersionInternal(ctx, pinnedVersionStorageKey(pluginType, name))
	if err != nil && !errors.Is(err, pluginutil.ErrPinnedVersionNotFound) {
		return err
	}
	if pin != nil && pin.Version == pluginVersion {
		return ErrPinnedVersion
	}

	// Check the name under which the plugin exists, but if it's unfound, don't return any error.
	pluginKey := path.Join(pluginType.String(), name)
	if pluginVersion != "" {
		pluginKey = path.Join(pluginKey, pluginVersion)
	}
	out, err := c.catalogView.Get(ctx, pluginKey)
	if err != nil || out == nil {
		pluginKey = name
	}

	return c.catalogView.Delete(ctx, pluginKey)
}

// List returns a list of all the known plugin names. If an external and builtin
// plugin share the same name, only one instance of the name will be returned.
func (c *PluginCatalog) List(ctx context.Context, pluginType consts.PluginType) ([]string, error) {
	plugins, err := c.listInternal(ctx, pluginType, false)
	if err != nil {
		return nil, err
	}

	// Use a set to de-dupe between builtin and unversioned external plugins.
	// External plugins with the same name as a builtin override the builtin.
	uniquePluginNames := make(map[string]struct{})
	for _, plugin := range plugins {
		uniquePluginNames[plugin.Name] = struct{}{}
	}

	retList := make([]string, 0, len(uniquePluginNames))
	for plugin := range uniquePluginNames {
		retList = append(retList, plugin)
	}

	return retList, nil
}

// ListPluginsWithRuntime lists the plugins that are registered with a given runtime
func (c *PluginCatalog) ListPluginsWithRuntime(ctx context.Context, runtime string) ([]string, error) {
	plugins, err := c.collectAllPlugins(ctx)
	if err != nil {
		return nil, err
	}

	var ret []string
	for _, plugin := range plugins {
		if plugin.Runtime == runtime {
			ret = append(ret, plugin.Name)
		}
		if plugin.OCIImage != "" {
			plugin.Tmpdir = c.tmpdir
		}
	}
	return ret, nil
}

func (c *PluginCatalog) ListVersionedPlugins(ctx context.Context, pluginType consts.PluginType) ([]pluginutil.VersionedPlugin, error) {
	return c.listInternal(ctx, pluginType, true)
}

func (c *PluginCatalog) listInternal(ctx context.Context, pluginType consts.PluginType, includeVersioned bool) ([]pluginutil.VersionedPlugin, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var result []pluginutil.VersionedPlugin

	plugins, err := c.collectAllPlugins(ctx)
	if err != nil {
		return nil, err
	}

	unversionedPlugins := make(map[string]struct{})
	for _, plugin := range plugins {
		var semanticVersion *semver.Version
		if plugin.Version == "" {
			semanticVersion, err = semver.NewVersion("0.0.0")
			if err != nil {
				return nil, err
			}
		} else {
			if !includeVersioned {
				continue
			}

			semanticVersion, err = semver.NewVersion(plugin.Version)
			if err != nil {
				return nil, fmt.Errorf("unexpected error parsing %s %s plugin version: %w", plugin.Type.String(), plugin.Name, err)
			}
		}

		// Only list user-added plugins if they're of the given type.
		if plugin.Type != consts.PluginTypeUnknown && plugin.Type != pluginType {
			continue
		}

		if plugin.OCIImage != "" {
			plugin.Tmpdir = c.tmpdir
		}

		result = append(result, pluginutil.VersionedPlugin{
			Name:            plugin.Name,
			Type:            plugin.Type.String(),
			OCIImage:        plugin.OCIImage,
			Runtime:         plugin.Runtime,
			Version:         plugin.Version,
			SHA256:          hex.EncodeToString(plugin.Sha256),
			SemanticVersion: semanticVersion,
		})

		if plugin.Version == "" {
			unversionedPlugins[plugin.Name] = struct{}{}
		}
	}

	// Get the builtin plugins.
	builtinPlugins := c.builtinRegistry.Keys(pluginType)
	for _, plugin := range builtinPlugins {
		// Unversioned plugins fully replace builtins of the same name.
		if _, ok := unversionedPlugins[plugin]; ok {
			continue
		}

		version := versions.GetBuiltinVersion(pluginType, plugin)
		semanticVersion, err := semver.NewVersion(version)
		deprecationStatus, _ := c.builtinRegistry.DeprecationStatus(plugin, pluginType)
		if err != nil {
			return nil, err
		}
		result = append(result, pluginutil.VersionedPlugin{
			Name:              plugin,
			Type:              pluginType.String(),
			Version:           version,
			Builtin:           true,
			SemanticVersion:   semanticVersion,
			DeprecationStatus: deprecationStatus.String(),
		})
	}

	// Sort the result for consistent ordering.
	sortVersionedPlugins(result)

	return result, nil
}

func (c *PluginCatalog) collectAllPlugins(ctx context.Context) ([]*pluginutil.PluginRunner, error) {
	// Collect keys for external plugins in the barrier.
	keys, err := logical.CollectKeys(ctx, c.catalogView)
	if err != nil {
		return nil, err
	}

	var plugins []*pluginutil.PluginRunner
	for _, key := range keys {
		// Skip: pinned version entry.
		if strings.HasPrefix(key, pinnedVersionStoragePrefix) {
			continue
		}

		entry, err := c.catalogView.Get(ctx, key)
		if err != nil || entry == nil {
			continue
		}

		plugin := new(pluginutil.PluginRunner)
		if err := jsonutil.DecodeJSON(entry.Value, plugin); err != nil {
			return nil, fmt.Errorf("failed to decode plugin entry: %w", err)
		}

		plugins = append(plugins, plugin)
	}

	return plugins, nil
}

func sortVersionedPlugins(versionedPlugins []pluginutil.VersionedPlugin) {
	sort.SliceStable(versionedPlugins, func(i, j int) bool {
		left, right := versionedPlugins[i], versionedPlugins[j]
		if left.Type != right.Type {
			return left.Type < right.Type
		}
		if left.Name != right.Name {
			return left.Name < right.Name
		}
		if left.Version != right.Version {
			return right.SemanticVersion.GreaterThan(left.SemanticVersion)
		}

		return false
	})
}
