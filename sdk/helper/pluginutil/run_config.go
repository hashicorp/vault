// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginutil

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-secure-stdlib/plugincontainer"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginruntimeutil"
)

const (
	// Labels for plugin container ownership
	labelVaultPID           = "com.hashicorp.vault.pid"
	labelVaultClusterID     = "com.hashicorp.vault.cluster.id"
	labelVaultPluginName    = "com.hashicorp.vault.plugin.name"
	labelVaultPluginVersion = "com.hashicorp.vault.plugin.version"
	labelVaultPluginType    = "com.hashicorp.vault.plugin.type"
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
	command  string
	image    string
	imageTag string
	args     []string
	sha256   []byte

	// Initialized with what's in PluginRunner.Env, but can be added to
	env []string

	runtimeConfig *pluginruntimeutil.PluginRuntimeConfig

	PluginClientConfig
	tmpdir string
}

func (rc runConfig) mlockEnabled() bool {
	return rc.MLock || (rc.Wrapper != nil && rc.Wrapper.MlockEnabled())
}

func (rc runConfig) generateCmd(ctx context.Context) (cmd *exec.Cmd, clientTLSConfig *tls.Config, err error) {
	cmd = exec.Command(rc.command, rc.args...)
	env := rc.env

	// Add the mlock setting to the ENV of the plugin
	if rc.mlockEnabled() {
		env = append(env, fmt.Sprintf("%s=%s", PluginMlockEnabled, "true"))
	}
	version, err := rc.Wrapper.VaultVersion(ctx)
	if err != nil {
		return nil, nil, err
	}
	env = append(env, fmt.Sprintf("%s=%s", PluginVaultVersionEnv, version))

	if rc.IsMetadataMode {
		rc.Logger = rc.Logger.With("metadata", "true")
	}
	metadataEnv := fmt.Sprintf("%s=%t", PluginMetadataModeEnv, rc.IsMetadataMode)
	env = append(env, metadataEnv)

	automtlsEnv := fmt.Sprintf("%s=%t", PluginAutoMTLSEnv, rc.AutoMTLS)
	env = append(env, automtlsEnv)

	if !rc.AutoMTLS && !rc.IsMetadataMode {
		// Get a CA TLS Certificate
		certBytes, key, err := generateCert()
		if err != nil {
			return nil, nil, err
		}

		// Use CA to sign a client cert and return a configured TLS config
		clientTLSConfig, err = createClientTLSConfig(certBytes, key)
		if err != nil {
			return nil, nil, err
		}

		// Use CA to sign a server cert and wrap the values in a response wrapped
		// token.
		wrapToken, err := wrapServerConfig(ctx, rc.Wrapper, certBytes, key)
		if err != nil {
			return nil, nil, err
		}

		// Add the response wrap token to the ENV of the plugin
		env = append(env, fmt.Sprintf("%s=%s", PluginUnwrapTokenEnv, wrapToken))
	}

	if rc.image == "" {
		// go-plugin has always overridden user-provided env vars with the OS
		// (Vault process) env vars, but we want plugins to be able to override
		// the Vault process env. We don't want to make a breaking change in
		// go-plugin so always set SkipHostEnv and replicate the legacy behavior
		// ourselves if user opts in.
		if legacy, _ := strconv.ParseBool(os.Getenv(PluginUseLegacyEnvLayering)); legacy {
			// Env vars are layered as follows, with later entries overriding
			// earlier entries if there are duplicate keys:
			// 1. Env specified at plugin registration
			// 2. Env from Vault SDK
			// 3. Env from Vault process (OS)
			// 4. Env from go-plugin
			cmd.Env = append(env, os.Environ()...)
		} else {
			// Env vars are layered as follows, with later entries overriding
			// earlier entries if there are duplicate keys:
			// 1. Env from Vault process (OS)
			// 2. Env specified at plugin registration
			// 3. Env from Vault SDK
			// 4. Env from go-plugin
			cmd.Env = append(os.Environ(), env...)
		}
	} else {
		// Containerized plugins do not inherit any env vars from Vault.
		cmd.Env = env
	}

	return cmd, clientTLSConfig, nil
}

func (rc runConfig) makeConfig(ctx context.Context) (*plugin.ClientConfig, error) {
	cmd, clientTLSConfig, err := rc.generateCmd(ctx)
	if err != nil {
		return nil, err
	}

	clientConfig := &plugin.ClientConfig{
		HandshakeConfig:  rc.HandshakeConfig,
		VersionedPlugins: rc.PluginSets,
		TLSConfig:        clientTLSConfig,
		Logger:           rc.Logger,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
			plugin.ProtocolGRPC,
		},
		AutoMTLS:    rc.AutoMTLS,
		SkipHostEnv: true,
	}
	if rc.image == "" {
		clientConfig.Cmd = cmd
		clientConfig.SecureConfig = &plugin.SecureConfig{
			Checksum: rc.sha256,
			Hash:     sha256.New(),
		}
	} else {
		containerCfg, err := rc.containerConfig(ctx, cmd.Env)
		if err != nil {
			return nil, err
		}
		clientConfig.RunnerFunc = containerCfg.NewContainerRunner
		clientConfig.UnixSocketConfig = &plugin.UnixSocketConfig{
			Group:   strconv.Itoa(containerCfg.GroupAdd),
			TempDir: rc.tmpdir,
		}
		clientConfig.GRPCBrokerMultiplex = true
	}
	return clientConfig, nil
}

func (rc runConfig) containerConfig(ctx context.Context, env []string) (*plugincontainer.Config, error) {
	clusterID, err := rc.Wrapper.ClusterID(ctx)
	if err != nil {
		return nil, err
	}
	cfg := &plugincontainer.Config{
		Image:  rc.image,
		Tag:    rc.imageTag,
		SHA256: fmt.Sprintf("%x", rc.sha256),

		Env:        env,
		GroupAdd:   os.Getegid(),
		Runtime:    consts.DefaultContainerPluginOCIRuntime,
		CapIPCLock: rc.mlockEnabled(),
		Labels: map[string]string{
			labelVaultPID:           strconv.Itoa(os.Getpid()),
			labelVaultClusterID:     clusterID,
			labelVaultPluginName:    rc.PluginClientConfig.Name,
			labelVaultPluginType:    rc.PluginClientConfig.PluginType.String(),
			labelVaultPluginVersion: rc.PluginClientConfig.Version,
		},
	}

	// Use rc.command and rc.args directly instead of cmd.Path and cmd.Args, as
	// exec.Command may mutate the provided command.
	if rc.command != "" {
		cfg.Entrypoint = []string{rc.command}
	}
	if len(rc.args) > 0 {
		cfg.Args = rc.args
	}
	if rc.runtimeConfig != nil {
		cfg.CgroupParent = rc.runtimeConfig.CgroupParent
		cfg.NanoCpus = rc.runtimeConfig.CPU
		cfg.Memory = rc.runtimeConfig.Memory
		if rc.runtimeConfig.OCIRuntime != "" {
			cfg.Runtime = rc.runtimeConfig.OCIRuntime
		}
		if rc.runtimeConfig.Rootless {
			cfg.Rootless = true
		}
	}

	return cfg, nil
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
	var image, imageTag string
	if r.OCIImage != "" {
		image = r.OCIImage
		imageTag = strings.TrimPrefix(r.Version, "v")
	}
	rc := runConfig{
		command:       r.Command,
		image:         image,
		imageTag:      imageTag,
		args:          r.Args,
		sha256:        r.Sha256,
		env:           r.Env,
		runtimeConfig: r.RuntimeConfig,
		tmpdir:        r.Tmpdir,
		PluginClientConfig: PluginClientConfig{
			Name:       r.Name,
			PluginType: r.Type,
			Version:    r.Version,
		},
	}

	for _, opt := range opts {
		opt(&rc)
	}

	return rc.run(ctx)
}
