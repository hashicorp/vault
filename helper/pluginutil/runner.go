package pluginutil

import (
	"crypto/sha256"
	"crypto/tls"
	"flag"
	"fmt"
	"os/exec"
	"time"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/wrapping"
	log "github.com/mgutz/logxi/v1"
)

// Looker defines the plugin Lookup function that looks into the plugin catalog
// for availible plugins and returns a PluginRunner
type Looker interface {
	LookupPlugin(string) (*PluginRunner, error)
}

// Wrapper interface defines the functions needed by the runner to wrap the
// metadata needed to run a plugin process. This includes looking up Mlock
// configuration and wrapping data in a respose wrapped token.
// logical.SystemView implementataions satisfy this interface.
type RunnerUtil interface {
	ResponseWrapData(data map[string]interface{}, ttl time.Duration, jwt bool) (*wrapping.ResponseWrapInfo, error)
	MlockEnabled() bool
}

// LookWrapper defines the functions for both Looker and Wrapper
type LookRunnerUtil interface {
	Looker
	RunnerUtil
}

// PluginRunner defines the metadata needed to run a plugin securely with
// go-plugin.
type PluginRunner struct {
	Name           string                      `json:"name" structs:"name"`
	Command        string                      `json:"command" structs:"command"`
	Args           []string                    `json:"args" structs:"args"`
	Sha256         []byte                      `json:"sha256" structs:"sha256"`
	Builtin        bool                        `json:"builtin" structs:"builtin"`
	BuiltinFactory func() (interface{}, error) `json:"-" structs:"-"`
}

// Run takes a wrapper RunnerUtil instance along with the go-plugin paramaters and
// returns a configured plugin.Client with TLS Configured and a wrapping token set
// on PluginUnwrapTokenEnv for plugin process consumption.
func (r *PluginRunner) Run(wrapper RunnerUtil, pluginMap map[string]plugin.Plugin, hs plugin.HandshakeConfig, env []string, logger log.Logger) (*plugin.Client, error) {
	return r.runCommon(wrapper, pluginMap, hs, env, logger, false)
}

// RunMetadataMode returns a configured plugin.Client that will dispense a plugin
// in metadata mode. The PluginMetadaModeEnv is passed in as part of the Cmd to
// plugin.Client, and consumed by the plugin process on pluginutil.VaultPluginTLSProvider.
func (r *PluginRunner) RunMetadataMode(wrapper RunnerUtil, pluginMap map[string]plugin.Plugin, hs plugin.HandshakeConfig, env []string, logger log.Logger) (*plugin.Client, error) {
	return r.runCommon(wrapper, pluginMap, hs, env, logger, true)

}

func (r *PluginRunner) runCommon(wrapper RunnerUtil, pluginMap map[string]plugin.Plugin, hs plugin.HandshakeConfig, env []string, logger log.Logger, isMetadataMode bool) (*plugin.Client, error) {
	cmd := exec.Command(r.Command, r.Args...)
	cmd.Env = append(cmd.Env, env...)

	// Add the mlock setting to the ENV of the plugin
	if wrapper.MlockEnabled() {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginMlockEnabled, "true"))
	}

	// Create logger for the plugin client
	clogger := &hclogFaker{
		logger: logger,
	}
	namedLogger := clogger.ResetNamed("plugin")

	var clientTLSConfig *tls.Config
	if !isMetadataMode {
		// Add the metadata mode ENV and set it to false
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginMetadaModeEnv, "false"))

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
		wrapToken, err := wrapServerConfig(wrapper, certBytes, key)
		if err != nil {
			return nil, err
		}

		// Add the response wrap token to the ENV of the plugin
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginUnwrapTokenEnv, wrapToken))
	} else {
		namedLogger = clogger.ResetNamed("plugin.metadata")
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginMetadaModeEnv, "true"))
	}

	secureConfig := &plugin.SecureConfig{
		Checksum: r.Sha256,
		Hash:     sha256.New(),
	}

	clientConfig := &plugin.ClientConfig{
		HandshakeConfig: hs,
		Plugins:         pluginMap,
		Cmd:             cmd,
		SecureConfig:    secureConfig,
		TLSConfig:       clientTLSConfig,
		Logger:          namedLogger,
	}

	client := plugin.NewClient(clientConfig)

	return client, nil
}

type APIClientMeta struct {
	// These are set by the command line flags.
	flagCACert     string
	flagCAPath     string
	flagClientCert string
	flagClientKey  string
	flagInsecure   bool
}

func (f *APIClientMeta) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("vault plugin settings", flag.ContinueOnError)

	fs.StringVar(&f.flagCACert, "ca-cert", "", "")
	fs.StringVar(&f.flagCAPath, "ca-path", "", "")
	fs.StringVar(&f.flagClientCert, "client-cert", "", "")
	fs.StringVar(&f.flagClientKey, "client-key", "", "")
	fs.BoolVar(&f.flagInsecure, "tls-skip-verify", false, "")

	return fs
}

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
