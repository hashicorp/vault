package pluginutil

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"os/exec"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/version"
)

type runConfig struct {
	// Provided by PluginRunner
	command string
	args    []string
	sha256  []byte

	// Initialized with what's in PluginRunner.Env, but can be added to
	env []string

	wrapper        RunnerUtil
	pluginSets     map[int]plugin.PluginSet
	hs             plugin.HandshakeConfig
	logger         log.Logger
	isMetadataMode bool
	autoMTLS       bool
}

func (rc runConfig) makeConfig(ctx context.Context) (*plugin.ClientConfig, error) {
	cmd := exec.Command(rc.command, rc.args...)
	cmd.Env = append(cmd.Env, rc.env...)

	// Add the mlock setting to the ENV of the plugin
	if rc.wrapper != nil && rc.wrapper.MlockEnabled() {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginMlockEnabled, "true"))
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginVaultVersionEnv, version.GetVersion().Version))

	if rc.isMetadataMode {
		rc.logger = rc.logger.With("metadata", "true")
	}
	metadataEnv := fmt.Sprintf("%s=%t", PluginMetadataModeEnv, rc.isMetadataMode)
	cmd.Env = append(cmd.Env, metadataEnv)

	var clientTLSConfig *tls.Config
	if !rc.autoMTLS && !rc.isMetadataMode {
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
		wrapToken, err := wrapServerConfig(ctx, rc.wrapper, certBytes, key)
		if err != nil {
			return nil, err
		}

		// Add the response wrap token to the ENV of the plugin
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginUnwrapTokenEnv, wrapToken))
	}

	secureConfig := &plugin.SecureConfig{
		Checksum: rc.sha256,
		Hash:     sha256.New(),
	}

	clientConfig := &plugin.ClientConfig{
		HandshakeConfig:  rc.hs,
		VersionedPlugins: rc.pluginSets,
		Cmd:              cmd,
		SecureConfig:     secureConfig,
		TLSConfig:        clientTLSConfig,
		Logger:           rc.logger,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
			plugin.ProtocolGRPC,
		},
		AutoMTLS: rc.autoMTLS,
	}
	return clientConfig, nil
}

func (rc runConfig) run(ctx context.Context) (*plugin.Client, error) {
	clientConfig, err := rc.makeConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := plugin.NewClient(clientConfig)
	return client, nil
}

type RunOpt func(*runConfig)

func Env(env ...string) RunOpt {
	return func(rc *runConfig) {
		rc.env = append(rc.env, env...)
	}
}

func Runner(wrapper RunnerUtil) RunOpt {
	return func(rc *runConfig) {
		rc.wrapper = wrapper
	}
}

func PluginSets(pluginSets map[int]plugin.PluginSet) RunOpt {
	return func(rc *runConfig) {
		rc.pluginSets = pluginSets
	}
}

func HandshakeConfig(hs plugin.HandshakeConfig) RunOpt {
	return func(rc *runConfig) {
		rc.hs = hs
	}
}

func Logger(logger log.Logger) RunOpt {
	return func(rc *runConfig) {
		rc.logger = logger
	}
}

func MetadataMode(isMetadataMode bool) RunOpt {
	return func(rc *runConfig) {
		rc.isMetadataMode = isMetadataMode
	}
}

func AutoMTLS(autoMTLS bool) RunOpt {
	return func(rc *runConfig) {
		rc.autoMTLS = autoMTLS
	}
}

func (r *PluginRunner) RunConfig(ctx context.Context, opts ...RunOpt) (*plugin.Client, error) {
	rc := runConfig{
		command: r.Command,
		args:    r.Args,
		sha256:  r.Sha256,
		env:     r.Env,
	}

	for _, opt := range opts {
		opt(&rc)
	}

	return rc.run(ctx)
}
