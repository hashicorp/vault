// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package dbplugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

// PluginFactoryConfig carries the inputs needed to build a plugin database
// type. It exists so that new optional inputs can be added without breaking
// existing callers of PluginFactory and PluginFactoryVersion.
type PluginFactoryConfig struct {
	// PluginName is the name of the database plugin in the plugin catalog.
	PluginName string

	// PluginVersion optionally pins a specific plugin version. An empty string
	// selects the default, which is the builtin version when one exists.
	PluginVersion string

	Sys    pluginutil.LookRunnerUtil
	Logger log.Logger

	// Namespace, MountPoint, and ConnectionName identify the namespace, the
	// specific database secrets engine mount, and the configured connection
	// that this plugin instance serves. When set, the metrics middleware
	// attaches each as a label to the telemetry it emits. Leave them empty to
	// emit unlabeled metrics.
	Namespace      string
	MountPoint     string
	ConnectionName string
}

// validate reports whether the config carries the collaborators required to
// build a plugin. It exists because callers construct PluginFactoryConfig as a
// struct literal, so an omitted field would otherwise surface as a nil
// dereference deep inside plugin startup rather than as an error the caller can
// handle.
//
// Only the fields that would panic are checked. PluginName is deliberately not
// validated: an unknown or empty name is already reported by the plugin catalog
// lookup, which some callers rely on to resolve a default plugin.
func (cfg PluginFactoryConfig) validate() error {
	if cfg.Sys == nil {
		return errors.New("plugin system view is required")
	}
	if cfg.Logger == nil {
		return errors.New("logger is required")
	}
	return nil
}

// PluginFactory is used to build plugin database types. It wraps the database
// object in a logging and metrics middleware.
func PluginFactory(ctx context.Context, pluginName string, sys pluginutil.LookRunnerUtil, logger log.Logger) (Database, error) {
	return PluginFactoryVersion(ctx, pluginName, "", sys, logger)
}

// PluginFactoryVersion is used to build plugin database types with a version specified.
// It wraps the database object in a logging and metrics middleware.
func PluginFactoryVersion(ctx context.Context, pluginName string, pluginVersion string, sys pluginutil.LookRunnerUtil, logger log.Logger) (Database, error) {
	return PluginFactoryWithConfig(ctx, PluginFactoryConfig{
		PluginName:    pluginName,
		PluginVersion: pluginVersion,
		Sys:           sys,
		Logger:        logger,
	})
}

// PluginFactoryWithConfig is used to build plugin database types from a config
// struct. It wraps the database object in a logging and metrics middleware.
func PluginFactoryWithConfig(ctx context.Context, cfg PluginFactoryConfig) (Database, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	// Look for plugin in the plugin catalog
	pluginRunner, err := cfg.Sys.LookupPluginVersion(ctx, cfg.PluginName, consts.PluginTypeDatabase, cfg.PluginVersion)
	if err != nil {
		return nil, err
	}

	namedLogger := cfg.Logger.Named(cfg.PluginName)

	var transport string
	var db Database
	if pluginRunner.Builtin {
		// Plugin is builtin so we can retrieve an instance of the interface
		// from the pluginRunner. Then cast it to a Database.
		dbRaw, err := pluginRunner.BuiltinFactory()
		if err != nil {
			return nil, errwrap.Wrapf("error initializing plugin: {{err}}", err)
		}

		var ok bool
		db, ok = dbRaw.(Database)
		if !ok {
			return nil, fmt.Errorf("unsupported database type: %q", cfg.PluginName)
		}

		transport = "builtin"

	} else {
		if pluginRunner.Download {
			if err = cfg.Sys.DownloadExtractVerifyPlugin(ctx, pluginRunner); err != nil {
				return nil, fmt.Errorf("failed to extract and verify plugin=%q version=%q: %w",
					pluginRunner.Name, pluginRunner.Version, err)
			}
		}

		config := pluginutil.PluginClientConfig{
			Name:            cfg.PluginName,
			PluginType:      consts.PluginTypeDatabase,
			Version:         cfg.PluginVersion,
			PluginSets:      PluginSets,
			HandshakeConfig: HandshakeConfig,
			Logger:          namedLogger,
			IsMetadataMode:  false,
			AutoMTLS:        true,
			Wrapper:         cfg.Sys,
			Tier:            pluginRunner.Tier,
		}

		// create a DatabasePluginClient instance
		db, err = NewPluginClient(ctx, cfg.Sys, config)
		if err != nil {
			return nil, err
		}

		// Switch on the underlying database client type to get the transport
		// method.
		switch db.(*DatabasePluginClient).Database.(type) {
		case *gRPCClient:
			transport = "gRPC"
		}

	}

	typeStr, err := db.Type()
	if err != nil {
		return nil, errwrap.Wrapf("error getting plugin type: {{err}}", err)
	}
	cfg.Logger.Debug("got database plugin instance", "type", typeStr)

	// Wrap with metrics middleware
	db = &databaseMetricsMiddleware{
		next:           db,
		typeStr:        typeStr,
		namespace:      cfg.Namespace,
		mountPoint:     cfg.MountPoint,
		connectionName: cfg.ConnectionName,
	}

	// Wrap with tracing middleware
	if namedLogger.IsTrace() {
		db = &databaseTracingMiddleware{
			next:   db,
			logger: namedLogger.With("transport", transport),
		}
	}

	return db, nil
}
