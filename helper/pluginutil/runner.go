package pluginutil

import (
	"crypto/sha256"
	"fmt"
	"os/exec"

	plugin "github.com/hashicorp/go-plugin"
)

type Looker interface {
	LookupPlugin(string) (*PluginRunner, error)
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

	// Add the response wrap token to the ENV of the plugin
	cmd := exec.Command(r.Command, r.Args...)
	cmd.Env = append(cmd.Env, env...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", PluginUnwrapTokenEnv, wrapToken))

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
