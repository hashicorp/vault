package config

import (
	"github.com/hashicorp/vault/helper/activedirectory"
)

type EngineConf struct {
	PasswordConf *PasswordConf
	ADConf       *activedirectory.Configuration
}

func (c *EngineConf) Map() map[string]interface{} {
	combined := make(map[string]interface{})
	for k, v := range c.PasswordConf.Map() {
		combined[k] = v
	}
	for k, v := range c.ADConf.Map() {
		combined[k] = v
	}
	return combined
}
