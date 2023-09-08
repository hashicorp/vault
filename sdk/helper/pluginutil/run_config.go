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
	labelVaultClusterID     = "com.hashicorp.vault.clusterID"
	labelVaultPluginName    = "com.hashicorp.vault.pluginName"
	labelVaultPluginVersion = "com.hashicorp.vault.pluginVersion"
	labelVaultPluginType    = "com.hashicorp.vault.pluginType"
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
}

func (rc runConfig) generateCmd(ctx context.Context) (cmd *exec.Cmd, clientTLSConfig *tls.Config, err error) {
	cmd = exec.Command(rc.command, rc.args...)
	cmd.Env = append(cmd.Env, rc.env...)

	// Add the mlock setting to the ENV of the plugin
	if rc.MLock || (rc.Wrapper != nil && rc.Wrapper.MlockEnabled()) {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginMlockEnabled, "true"))
	}
	version, err := rc.Wrapper.VaultVersion(ctx)
	if err != nil {
		return nil, nil, err
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginVaultVersionEnv, version))

	if rc.IsMetadataMode {
		rc.Logger = rc.Logger.With("metadata", "true")
	}
	metadataEnv := fmt.Sprintf("%s=%t", PluginMetadataModeEnv, rc.IsMetadataMode)
	cmd.Env = append(cmd.Env, metadataEnv)

	automtlsEnv := fmt.Sprintf("%s=%t", PluginAutoMTLSEnv, rc.AutoMTLS)
	cmd.Env = append(cmd.Env, automtlsEnv)

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
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginUnwrapTokenEnv, wrapToken))
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
		AutoMTLS: rc.AutoMTLS,
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
		clientConfig.SkipHostEnv = true
		clientConfig.RunnerFunc = containerCfg.NewContainerRunner
		clientConfig.UnixSocketConfig = &plugin.UnixSocketConfig{
			Group: strconv.Itoa(containerCfg.GroupAdd),
		}
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

		Env:      env,
		GroupAdd: os.Getgid(),
		Runtime:  consts.DefaultContainerPluginOCIRuntime,
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
