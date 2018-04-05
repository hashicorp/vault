package config

import (
	"errors"

	"github.com/hashicorp/vault/helper/activedirectory"
)

// newUnsetEngineConf provides an engineConf that's essentially empty,
// but won't cause a panic.
// If `engineConf.IsSet == false`, that reflects that the user doesn't
// currently have any configuration stored.
func newUnsetEngineConf() *EngineConf {
	return &EngineConf{
		PasswordConf: &PasswordConf{},
		ADConf:       &activedirectory.Configuration{},
		IsSet:        false,
	}
}

func newEngineConf(passwordConf *PasswordConf, adConf *activedirectory.Configuration) (*EngineConf, error) {

	// usernames and passwords aren't required for the AD client in general,
	// but for this backend, we need them for password rotation.
	if adConf.Username == "" || adConf.Password == "" {
		return nil, errors.New("a username and password must be provided to perform password rotation")
	}

	return &EngineConf{
		PasswordConf: passwordConf,
		ADConf:       adConf,
		IsSet:        true,
	}, nil
}

type EngineConf struct {
	PasswordConf *PasswordConf
	ADConf       *activedirectory.Configuration

	// IsSet reflects whether the *user* has set the configuration.
	IsSet bool
}

func (c *EngineConf) Map() map[string]interface{} {
	combined := make(map[string]interface{})
	if !c.IsSet {
		return combined
	}
	for k, v := range c.PasswordConf.Map() {
		combined[k] = v
	}
	for k, v := range c.ADConf.Map() {
		combined[k] = v
	}
	return combined
}
