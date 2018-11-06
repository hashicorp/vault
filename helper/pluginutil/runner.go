package pluginutil

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"flag"
	"fmt"
	"os/exec"
	"time"

	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/wrapping"
	"github.com/hashicorp/vault/version"
)

// Looker defines the plugin Lookup function that looks into the plugin catalog
// for available plugins and returns a PluginRunner
type Looker interface {
	LookupPlugin(context.Context, string, consts.PluginType) (*PluginRunner, error)
}

// RunnerUtil interface defines the functions needed by the runner to wrap the
// metadata needed to run a plugin process. This includes looking up Mlock
// configuration and wrapping data in a response wrapped token.
// logical.SystemView implementations satisfy this interface.
type RunnerUtil interface {
	ResponseWrapData(ctx context.Context, data map[string]interface{}, ttl time.Duration, jwt bool) (*wrapping.ResponseWrapInfo, error)
	MlockEnabled() bool
}

// LookRunnerUtil defines the functions for both Looker and Wrapper
type LookRunnerUtil interface {
	Looker
	RunnerUtil
}

// PluginRunner defines the metadata needed to run a plugin securely with
// go-plugin.
type PluginRunner struct {
	Name           string                      `json:"name" structs:"name"`
	Type           consts.PluginType           `json:"type" structs:"type"`
	Command        string                      `json:"command" structs:"command"`
	Args           []string                    `json:"args" structs:"args"`
	Env            []string                    `json:"env" structs:"env"`
	Sha256         []byte                      `json:"sha256" structs:"sha256"`
	Builtin        bool                        `json:"builtin" structs:"builtin"`
	BuiltinFactory func() (interface{}, error) `json:"-" structs:"-"`
}

// Run takes a wrapper RunnerUtil instance along with the go-plugin parameters and
// returns a configured plugin.Client with TLS Configured and a wrapping token set
// on PluginUnwrapTokenEnv for plugin process consumption.
func (r *PluginRunner) Run(ctx context.Context, wrapper RunnerUtil, pluginSets map[int]plugin.PluginSet, hs plugin.HandshakeConfig, env []string, logger log.Logger) (*plugin.Client, error) {
	return r.runCommon(ctx, wrapper, pluginSets, hs, env, logger, false)
}

// RunMetadataMode returns a configured plugin.Client that will dispense a plugin
// in metadata mode. The PluginMetadataModeEnv is passed in as part of the Cmd to
// plugin.Client, and consumed by the plugin process on pluginutil.VaultPluginTLSProvider.
func (r *PluginRunner) RunMetadataMode(ctx context.Context, wrapper RunnerUtil, pluginSets map[int]plugin.PluginSet, hs plugin.HandshakeConfig, env []string, logger log.Logger) (*plugin.Client, error) {
	return r.runCommon(ctx, wrapper, pluginSets, hs, env, logger, true)

}

func (r *PluginRunner) runCommon(ctx context.Context, wrapper RunnerUtil, pluginSets map[int]plugin.PluginSet, hs plugin.HandshakeConfig, env []string, logger log.Logger, isMetadataMode bool) (*plugin.Client, error) {
	cmd := exec.Command(r.Command, r.Args...)

	// `env` should always go last to avoid overwriting internal values that might
	// have been provided externally.
	cmd.Env = append(cmd.Env, r.Env...)
	cmd.Env = append(cmd.Env, env...)

	// Add the mlock setting to the ENV of the plugin
	if wrapper != nil && wrapper.MlockEnabled() {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginMlockEnabled, "true"))
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginVaultVersionEnv, version.GetVersion().Version))

	var clientTLSConfig *tls.Config
	if !isMetadataMode {
		// Add the metadata mode ENV and set it to false
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginMetadataModeEnv, "false"))

		// Get a CA TLS Certificate
		certBytes, key, err := generateCert()
		if err != nil {
			return nil, err
		}

		// Use CA to sign a client cert and return a configured TLS config
		clientTLSConfig, err = createClientTLSConfig(certBytes, key)
		if err != nil {
			return nil, err
		}

		// Use CA to sign a server cert and wrap the values in a response wrapped
		// token.
		wrapToken, err := wrapServerConfig(ctx, wrapper, certBytes, key)
		if err != nil {
			return nil, err
		}

		// Add the response wrap token to the ENV of the plugin
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginUnwrapTokenEnv, wrapToken))
	} else {
		logger = logger.With("metadata", "true")
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginMetadataModeEnv, "true"))
	}

	secureConfig := &plugin.SecureConfig{
		Checksum: r.Sha256,
		Hash:     sha256.New(),
	}

	clientConfig := &plugin.ClientConfig{
		HandshakeConfig:  hs,
		VersionedPlugins: pluginSets,
		Cmd:              cmd,
		SecureConfig:     secureConfig,
		TLSConfig:        clientTLSConfig,
		Logger:           logger,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
			plugin.ProtocolGRPC,
		},
	}

	client := plugin.NewClient(clientConfig)

	return client, nil
}

// APIClientMeta is a helper that plugins can use to configure TLS connections
// back to Vault.
type APIClientMeta struct {
	// These are set by the command line flags.
	flagCACert     string
	flagCAPath     string
	flagClientCert string
	flagClientKey  string
	flagInsecure   bool
}

// FlagSet returns the flag set for configuring the TLS connection
func (f *APIClientMeta) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("vault plugin settings", flag.ContinueOnError)

	fs.StringVar(&f.flagCACert, "ca-cert", "", "")
	fs.StringVar(&f.flagCAPath, "ca-path", "", "")
	fs.StringVar(&f.flagClientCert, "client-cert", "", "")
	fs.StringVar(&f.flagClientKey, "client-key", "", "")
	fs.BoolVar(&f.flagInsecure, "tls-skip-verify", false, "")

	return fs
}

// GetTLSConfig will return a TLSConfig based off the values from the flags
func (f *APIClientMeta) GetTLSConfig() *api.TLSConfig {
	// If we need custom TLS configuration, then set it
	if f.flagCACert != "" || f.flagCAPath != "" || f.flagClientCert != "" || f.flagClientKey != "" || f.flagInsecure {
		t := &api.TLSConfig{
			CACert:        f.flagCACert,
			CAPath:        f.flagCAPath,
			ClientCert:    f.flagClientCert,
			ClientKey:     f.flagClientKey,
			TLSServerName: "",
			Insecure:      f.flagInsecure,
		}

		return t
	}

	return nil
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
