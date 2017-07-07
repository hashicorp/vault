package plugin

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	bplugin "github.com/hashicorp/vault/logical/plugin"
)

// Factory returns a configured plugin logical.Backend.
func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	_, ok := conf.Config["plugin_name"]
	if !ok {
		return nil, fmt.Errorf("plugin_name not provided")
	}
	b, err := Backend(conf)
	if err != nil {
		return nil, err
	}

	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

// Backend returns an instance of the backend, either as a plugin if external
// or as a concrete implementation if builtin, casted as logical.Backend.
func Backend(conf *logical.BackendConfig) (logical.Backend, error) {
	name := conf.Config["plugin_name"]
	sys := conf.System

	b, err := bplugin.NewBackend(name, sys)
	if err != nil {
		return nil, err
	}

	return b, nil
}
