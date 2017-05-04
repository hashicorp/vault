package pluginutil

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"os/exec"
	"time"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/wrapping"
)

// Looker defines the plugin Lookup function that looks into the plugin catalog
// for availible plugins and returns a PluginRunner
type Looker interface {
	LookupPlugin(string) (*PluginRunner, error)
}

// Wrapper interface defines the functions needed by the runner to wrap the
// metadata needed to run a plugin process. This includes looking up Mlock
// configuration and wrapping data in a respose wrapped token.
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
	Name           string                      `json:"name"`
	Command        string                      `json:"command"`
	Args           []string                    `json:"args"`
	Sha256         []byte                      `json:"sha256"`
	Builtin        bool                        `json:"builtin"`
	BuiltinFactory func() (interface{}, error) `json:"-"`
}

// Run takes a wrapper instance, and the go-plugin paramaters and executes a
// plugin.
func (r *PluginRunner) Run(wrapper RunnerUtil, pluginMap map[string]plugin.Plugin, hs plugin.HandshakeConfig, env []string) (*plugin.Client, error) {
	// Get a CA TLS Certificate
	certBytes, key, err := generateCert()
	if err != nil {
		return nil, err
	}

	// Use CA to sign a client cert and return a configured TLS config
	clientTLSConfig, err := createClientTLSConfig(certBytes, key)
	if err != nil {
		return nil, err
	}

	// Use CA to sign a server cert and wrap the values in a response wrapped
	// token.
	wrapToken, err := wrapServerConfig(wrapper, certBytes, key)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(r.Command, r.Args...)
	cmd.Env = append(cmd.Env, env...)
	// Add the response wrap token to the ENV of the plugin
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginUnwrapTokenEnv, wrapToken))
	// Add the mlock setting to the ENV of the plugin
	if wrapper.MlockEnabled() {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginMlockEnabled, "true"))
	}

	secureConfig := &plugin.SecureConfig{
		Checksum: r.Sha256,
		Hash:     sha256.New(),
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: hs,
		Plugins:         pluginMap,
		Cmd:             cmd,
		TLSConfig:       clientTLSConfig,
		SecureConfig:    secureConfig,
	})

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
	fs := flag.NewFlagSet("tls settings", flag.ContinueOnError)

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
