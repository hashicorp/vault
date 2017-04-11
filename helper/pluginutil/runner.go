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

type Looker interface {
	LookupPlugin(string) (*PluginRunner, error)
}

type Wrapper interface {
	ResponseWrapData(data map[string]interface{}, ttl time.Duration, jwt bool) (string, error)
	MlockDisabled() bool
}

type LookWrapper interface {
	Looker
	Wrapper
}

type PluginRunner struct {
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Sha256  []byte   `json:"sha256"`
	Builtin bool     `json:"builtin"`
}

func (r *PluginRunner) Run(wrapper Wrapper, pluginMap map[string]plugin.Plugin, hs plugin.HandshakeConfig, env []string) (*plugin.Client, error) {
	// Get a CA TLS Certificate
	CACertBytes, CACert, CAKey, err := GenerateCACert()
	if err != nil {
		return nil, err
	}

	// Use CA to sign a client cert and return a configured TLS config
	clientTLSConfig, err := CreateClientTLSConfig(CACert, CAKey)
	if err != nil {
		return nil, err
	}

	// Use CA to sign a server cert and wrap the values in a response wrapped
	// token.
	wrapToken, err := WrapServerConfig(wrapper, CACertBytes, CACert, CAKey)
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

func OptionallyEnableMlock() error {
	if os.Getenv(PluginMlockEnabled) == "true" {
		return mlock.LockMemory()
	}

	return nil
}
