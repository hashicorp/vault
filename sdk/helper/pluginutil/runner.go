// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginutil

import (
	"context"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"google.golang.org/grpc"
)

// Looker defines the plugin Lookup function that looks into the plugin catalog
// for available plugins and returns a PluginRunner
type Looker interface {
	LookupPlugin(ctx context.Context, pluginName string, pluginType consts.PluginType) (*PluginRunner, error)
	LookupPluginVersion(ctx context.Context, pluginName string, pluginType consts.PluginType, version string) (*PluginRunner, error)
}

// RunnerUtil interface defines the functions needed by the runner to wrap the
// metadata needed to run a plugin process. This includes looking up Mlock
// configuration and wrapping data in a response wrapped token.
// logical.SystemView implementations satisfy this interface.
type RunnerUtil interface {
	NewPluginClient(ctx context.Context, config PluginClientConfig) (PluginClient, error)
	ResponseWrapData(ctx context.Context, data map[string]interface{}, ttl time.Duration, jwt bool) (*wrapping.ResponseWrapInfo, error)
	MlockEnabled() bool
	VaultVersion(ctx context.Context) (string, error)
}

// LookRunnerUtil defines the functions for both Looker and Wrapper
type LookRunnerUtil interface {
	Looker
	RunnerUtil
}

type PluginClient interface {
	Conn() grpc.ClientConnInterface
	Reload() error
	plugin.ClientProtocol
}

const MultiplexingCtxKey string = "multiplex_id"

// PluginRunner defines the metadata needed to run a plugin securely with
// go-plugin.
type PluginRunner struct {
	Name           string                      `json:"name" structs:"name"`
	Type           consts.PluginType           `json:"type" structs:"type"`
	Version        string                      `json:"version" structs:"version"`
	Command        string                      `json:"command" structs:"command"`
	Args           []string                    `json:"args" structs:"args"`
	Env            []string                    `json:"env" structs:"env"`
	Sha256         []byte                      `json:"sha256" structs:"sha256"`
	Builtin        bool                        `json:"builtin" structs:"builtin"`
	BuiltinFactory func() (interface{}, error) `json:"-" structs:"-"`
}

// SetPluginInput is only used as input for the plugin catalog's set methods.
// We don't use the very similar PluginRunner struct to avoid confusion about
// what's settable, which does not include the builtin fields.
type SetPluginInput struct {
	Name     string
	Type     consts.PluginType
	Version  string
	Command  string
	OCIImage string
	Args     []string
	Env      []string
	Sha256   []byte
}

// Run takes a wrapper RunnerUtil instance along with the go-plugin parameters and
// returns a configured plugin.Client with TLS Configured and a wrapping token set
// on PluginUnwrapTokenEnv for plugin process consumption.
func (r *PluginRunner) Run(ctx context.Context, wrapper RunnerUtil, pluginSets map[int]plugin.PluginSet, hs plugin.HandshakeConfig, env []string, logger log.Logger) (*plugin.Client, error) {
	return r.RunConfig(ctx,
		Runner(wrapper),
		PluginSets(pluginSets),
		HandshakeConfig(hs),
		Env(env...),
		Logger(logger),
		MetadataMode(false),
	)
}

// RunMetadataMode returns a configured plugin.Client that will dispense a plugin
// in metadata mode. The PluginMetadataModeEnv is passed in as part of the Cmd to
// plugin.Client, and consumed by the plugin process on api.VaultPluginTLSProvider.
func (r *PluginRunner) RunMetadataMode(ctx context.Context, wrapper RunnerUtil, pluginSets map[int]plugin.PluginSet, hs plugin.HandshakeConfig, env []string, logger log.Logger) (*plugin.Client, error) {
	return r.RunConfig(ctx,
		Runner(wrapper),
		PluginSets(pluginSets),
		HandshakeConfig(hs),
		Env(env...),
		Logger(logger),
		MetadataMode(true),
	)
}

// VersionedPlugin holds any versioning information stored about a plugin in the
// plugin catalog.
type VersionedPlugin struct {
	Type              string `json:"type"` // string instead of consts.PluginType so that we get the string form in API responses.
	Name              string `json:"name"`
	Version           string `json:"version"`
	SHA256            string `json:"sha256,omitempty"`
	Builtin           bool   `json:"builtin"`
	DeprecationStatus string `json:"deprecation_status,omitempty"`

	// Pre-parsed semver struct of the Version field
	SemanticVersion *version.Version `json:"-"`
}

// CtxCancelIfCanceled takes a context cancel func and a context. If the context is
// shutdown the cancelfunc is called. This is useful for merging two cancel
// functions.
func CtxCancelIfCanceled(f context.CancelFunc, ctxCanceler context.Context) chan struct{} {
	quitCh := make(chan struct{})
	go func() {
		select {
		case <-quitCh:
		case <-ctxCanceler.Done():
			f()
		}
	}()
	return quitCh
}
