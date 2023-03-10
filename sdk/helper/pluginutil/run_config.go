package pluginutil

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/consts"
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

// parse information on reattaching to unmanaged plugins out of a
// JSON-encoded environment variable.
func parseReattachPlugins(in string) (map[string]*plugin.ReattachConfig, error) {
	unmanagedPlugins := map[string]*plugin.ReattachConfig{}
	if in != "" {
		type reattachConfig struct {
			Protocol        string
			ProtocolVersion int
			Addr            struct {
				Network string
				String  string
			}
			Pid  int
			Test bool
		}
		var m map[string]reattachConfig
		err := json.Unmarshal([]byte(in), &m)
		if err != nil {
			return unmanagedPlugins, fmt.Errorf("Invalid format for VAULT_REATTACH_PLUGINS: %w", err)
		}
		for p, c := range m {
			// TODO(JM): parse plugin name??
			var addr net.Addr
			switch c.Addr.Network {
			case "unix":
				addr, err = net.ResolveUnixAddr("unix", c.Addr.String)
				if err != nil {
					return unmanagedPlugins, fmt.Errorf("Invalid unix socket path %q for %q: %w", c.Addr.String, p, err)
				}
			case "tcp":
				addr, err = net.ResolveTCPAddr("tcp", c.Addr.String)
				if err != nil {
					return unmanagedPlugins, fmt.Errorf("Invalid TCP address %q for %q: %w", c.Addr.String, p, err)
				}
			default:
				return unmanagedPlugins, fmt.Errorf("Unknown address type %q for %q", c.Addr.Network, p)
			}
			unmanagedPlugins["backend"] = &plugin.ReattachConfig{
				Protocol:        plugin.Protocol(c.Protocol),
				ProtocolVersion: c.ProtocolVersion,
				Pid:             c.Pid,
				Test:            c.Test,
				Addr:            addr,
			}
		}
	}
	return unmanagedPlugins, nil
}

func (rc runConfig) makeConfig(ctx context.Context) (*plugin.ClientConfig, error) {
	// The user can declare that certain plugins are being managed on
	// Vault's behalf using this environment variable. This is used
	// primarily for plugin debugging purposes.
	unmanagedPlugins, err := parseReattachPlugins(os.Getenv("VAULT_REATTACH_PLUGINS"))
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n\n%#+v\n\n", unmanagedPlugins)

	cmd := exec.Command(rc.command, rc.args...)
	cmd.Env = append(cmd.Env, rc.env...)

	// Add the mlock setting to the ENV of the plugin
	if rc.MLock || (rc.Wrapper != nil && rc.Wrapper.MlockEnabled()) {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginMlockEnabled, "true"))
	}
	version, err := rc.Wrapper.VaultVersion(ctx)
	if err != nil {
		return nil, err
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginVaultVersionEnv, version))

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
	_ = secureConfig

	clientConfig := &plugin.ClientConfig{
		HandshakeConfig:  rc.HandshakeConfig,
		VersionedPlugins: rc.PluginSets,
		// Can't set these with Reattach
		// Cmd:              cmd,
		// SecureConfig: secureConfig,
		TLSConfig: clientTLSConfig,
		Logger:    rc.Logger,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
			plugin.ProtocolGRPC,
		},
		AutoMTLS:   rc.AutoMTLS,
		Reattach:   unmanagedPlugins["backend"],
		SyncStdout: rc.Logger.StandardWriter(&log.StandardLoggerOptions{}),
		SyncStderr: rc.Logger.StandardWriter(&log.StandardLoggerOptions{}),
		// SyncStdout: logging.PluginOutputMonitor(fmt.Sprintf("%s:stdout", "debug-plugin")),
		// SyncStderr: logging.PluginOutputMonitor(fmt.Sprintf("%s:stderr", "debug-plugin")),
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
