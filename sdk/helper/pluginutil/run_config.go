package pluginutil

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"os/exec"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/version"
)

type PluginClientConfig struct {
	Name            string
	PluginType      consts.PluginType
	Version         string
	PluginSets      map[int]plugin.PluginSet
	HandshakeConfig plugin.HandshakeConfig
	Logger          log.Logger
	IsMetadataMode  bool
	AutoMTLS        bool
	MLock           bool
	Wrapper         RunnerUtil
}

type runConfig struct {
	// Provided by PluginRunner
	command string
	args    []string
	sha256  []byte

	// Initialized with what's in PluginRunner.Env, but can be added to
	env []string

	PluginClientConfig
}

func (rc runConfig) makeConfig(ctx context.Context) (*plugin.ClientConfig, error) {
	cmd := exec.Command(rc.command, rc.args...)
	cmd.Env = append(cmd.Env, rc.env...)

	// Add the mlock setting to the ENV of the plugin
	if rc.MLock || (rc.Wrapper != nil && rc.Wrapper.MlockEnabled()) {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginMlockEnabled, "true"))
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginVaultVersionEnv, version.GetVersion().Version))

	if rc.IsMetadataMode {
		rc.Logger = rc.Logger.With("metadata", "true")
	}
	metadataEnv := fmt.Sprintf("%s=%t", PluginMetadataModeEnv, rc.IsMetadataMode)
	cmd.Env = append(cmd.Env, metadataEnv)

	automtlsEnv := fmt.Sprintf("%s=%t", PluginAutoMTLSEnv, rc.AutoMTLS)
	cmd.Env = append(cmd.Env, automtlsEnv)

	var clientTLSConfig *tls.Config
	if !rc.AutoMTLS && !rc.IsMetadataMode {
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
		wrapToken, err := wrapServerConfig(ctx, rc.Wrapper, certBytes, key)
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
		HandshakeConfig:  rc.HandshakeConfig,
		VersionedPlugins: rc.PluginSets,
		Cmd:              cmd,
		SecureConfig:     secureConfig,
		TLSConfig:        clientTLSConfig,
		Logger:           rc.Logger,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
			plugin.ProtocolGRPC,
		},
		AutoMTLS: rc.AutoMTLS,
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
		rc.Wrapper = wrapper
	}
}

func PluginSets(pluginSets map[int]plugin.PluginSet) RunOpt {
	return func(rc *runConfig) {
		rc.PluginSets = pluginSets
	}
}

func HandshakeConfig(hs plugin.HandshakeConfig) RunOpt {
	return func(rc *runConfig) {
		rc.HandshakeConfig = hs
	}
}

func Logger(logger log.Logger) RunOpt {
	return func(rc *runConfig) {
		rc.Logger = logger
	}
}

func MetadataMode(isMetadataMode bool) RunOpt {
	return func(rc *runConfig) {
		rc.IsMetadataMode = isMetadataMode
	}
}

func AutoMTLS(autoMTLS bool) RunOpt {
	return func(rc *runConfig) {
		rc.AutoMTLS = autoMTLS
	}
}

func MLock(mlock bool) RunOpt {
	return func(rc *runConfig) {
		rc.MLock = mlock
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
