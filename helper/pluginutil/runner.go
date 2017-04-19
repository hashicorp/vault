package pluginutil

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"time"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/mlock"
)

var (
	// PluginUnwrapTokenEnv is the ENV name used to pass unwrap tokens to the
	// plugin.
	PluginMlockEnabled = "VAULT_PLUGIN_MLOCK_ENABLED"
)

// Looker defines the plugin Lookup function that looks into the plugin catalog
// for availible plugins and returns a PluginRunner
type Looker interface {
	LookupPlugin(string) (*PluginRunner, error)
}

// Wrapper interface defines the functions needed by the runner to wrap the
// metadata needed to run a plugin process. This includes looking up Mlock
// configuration and wrapping data in a respose wrapped token.
type Wrapper interface {
	ResponseWrapData(data map[string]interface{}, ttl time.Duration, jwt bool) (string, error)
	MlockDisabled() bool
}

// LookWrapper defines the functions for both Looker and Wrapper
type LookWrapper interface {
	Looker
	Wrapper
}

// PluginRunner defines the metadata needed to run a plugin securely with
// go-plugin.
type PluginRunner struct {
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Sha256  []byte   `json:"sha256"`
	Builtin bool     `json:"builtin"`
}

// Run takes a wrapper instance, and the go-plugin paramaters and executes a
// plugin.
func (r *PluginRunner) Run(wrapper Wrapper, pluginMap map[string]plugin.Plugin, hs plugin.HandshakeConfig, env []string) (*plugin.Client, error) {
	// Get a CA TLS Certificate
	certBytes, key, err := GenerateCert()
	if err != nil {
		return nil, err
	}

	// Use CA to sign a client cert and return a configured TLS config
	clientTLSConfig, err := CreateClientTLSConfig(certBytes, key)
	if err != nil {
		return nil, err
	}

	// Use CA to sign a server cert and wrap the values in a response wrapped
	// token.
	wrapToken, err := WrapServerConfig(wrapper, certBytes, key)
	if err != nil {
		return nil, err
	}

	mlock := "true"
	if wrapper.MlockDisabled() {
		mlock = "false"
	}

	cmd := exec.Command(r.Command, r.Args...)
	cmd.Env = append(cmd.Env, env...)
	// Add the response wrap token to the ENV of the plugin
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginUnwrapTokenEnv, wrapToken))
	// Add the mlock setting to the ENV of the plugin
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginMlockEnabled, mlock))

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

// OptionallyEnableMlock determines if mlock should be called, and if so enables
// mlock.
func OptionallyEnableMlock() error {
	if os.Getenv(PluginMlockEnabled) == "true" {
		return mlock.LockMemory()
	}

	return nil
}
