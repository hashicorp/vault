// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package builtinplugins

import (
	"context"

	credJWT "github.com/hashicorp/vault-plugin-auth-jwt"
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	credCert "github.com/hashicorp/vault/builtin/credential/cert"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	logicalPki "github.com/hashicorp/vault/builtin/logical/pki"
	logicalSsh "github.com/hashicorp/vault/builtin/logical/ssh"
	logicalTransit "github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

// Registry is inherently thread-safe because it's immutable.
// Thus, rather than creating multiple instances of it, we only need one.
var Registry = newRegistry()

// BuiltinFactory is the func signature that should be returned by
// the plugin's New() func.
type BuiltinFactory func() (interface{}, error)

// There are three forms of Backends which exist in the BuiltinRegistry.
type credentialBackend struct {
	logical.Factory
	consts.DeprecationStatus
}

type databasePlugin struct {
	Factory BuiltinFactory
	consts.DeprecationStatus
}

type logicalBackend struct {
	logical.Factory
	consts.DeprecationStatus
}

type removedBackend struct {
	*framework.Backend
}

func removedFactory(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
	removedBackend := &removedBackend{}
	removedBackend.Backend = &framework.Backend{}
	return removedBackend, nil
}

func newMinimalRegistry() *registry {
	return &registry{
		credentialBackends: map[string]credentialBackend{
			"approle":  {Factory: credAppRole.Factory},
			"cert":     {Factory: credCert.Factory},
			"jwt":      {Factory: credJWT.Factory},
			"oidc":     {Factory: credJWT.Factory},
			"userpass": {Factory: credUserpass.Factory},
		},
		databasePlugins: map[string]databasePlugin{},
		logicalBackends: map[string]logicalBackend{
			"kv":      {Factory: logicalKv.Factory},
			"pki":     {Factory: logicalPki.Factory},
			"ssh":     {Factory: logicalSsh.Factory},
			"transit": {Factory: logicalTransit.Factory},
		},
	}
}

func newRegistry() *registry {
	reg := newMinimalRegistry()

	extendAddonPlugins(reg)

	entAddExtPlugins(reg)

	return reg
}

func addExtPluginsImpl(r *registry) {}

type registry struct {
	credentialBackends map[string]credentialBackend
	databasePlugins    map[string]databasePlugin
	logicalBackends    map[string]logicalBackend
}

// Get returns the Factory func for a particular backend plugin from the
// plugins map.
func (r *registry) Get(name string, pluginType consts.PluginType) (func() (interface{}, error), bool) {
	switch pluginType {
	case consts.PluginTypeCredential:
		if f, ok := r.credentialBackends[name]; ok {
			return toFunc(f.Factory), ok
		}
	case consts.PluginTypeSecrets:
		if f, ok := r.logicalBackends[name]; ok {
			return toFunc(f.Factory), ok
		}
	case consts.PluginTypeDatabase:
		if f, ok := r.databasePlugins[name]; ok {
			return f.Factory, ok
		}
	default:
		return nil, false
	}

	return nil, false
}

// Keys returns the list of plugin names that are considered builtin plugins.
func (r *registry) Keys(pluginType consts.PluginType) []string {
	var keys []string
	switch pluginType {
	case consts.PluginTypeDatabase:
		for key, backend := range r.databasePlugins {
			keys = appendIfNotRemoved(keys, key, backend.DeprecationStatus)
		}
	case consts.PluginTypeCredential:
		for key, backend := range r.credentialBackends {
			keys = appendIfNotRemoved(keys, key, backend.DeprecationStatus)
		}
	case consts.PluginTypeSecrets:
		for key, backend := range r.logicalBackends {
			keys = appendIfNotRemoved(keys, key, backend.DeprecationStatus)
		}
	}
	return keys
}

func (r *registry) Contains(name string, pluginType consts.PluginType) bool {
	for _, key := range r.Keys(pluginType) {
		if key == name {
			return true
		}
	}
	return false
}

// DeprecationStatus returns the Deprecation status for a builtin with type `pluginType`
func (r *registry) DeprecationStatus(name string, pluginType consts.PluginType) (consts.DeprecationStatus, bool) {
	switch pluginType {
	case consts.PluginTypeCredential:
		if f, ok := r.credentialBackends[name]; ok {
			return f.DeprecationStatus, ok
		}
	case consts.PluginTypeSecrets:
		if f, ok := r.logicalBackends[name]; ok {
			return f.DeprecationStatus, ok
		}
	case consts.PluginTypeDatabase:
		if f, ok := r.databasePlugins[name]; ok {
			return f.DeprecationStatus, ok
		}
	default:
		return consts.Unknown, false
	}

	return consts.Unknown, false
}

func toFunc(ifc interface{}) func() (interface{}, error) {
	return func() (interface{}, error) {
		return ifc, nil
	}
}

func appendIfNotRemoved(keys []string, name string, status consts.DeprecationStatus) []string {
	if status != consts.Removed {
		return append(keys, name)
	}
	return keys
}
