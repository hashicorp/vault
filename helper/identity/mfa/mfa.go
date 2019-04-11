package mfa

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
)

func (c *Config) Clone() (*Config, error) {
	if c == nil {
		return nil, fmt.Errorf("nil config")
	}

	marshaledConfig, err := proto.Marshal(c)
	if err != nil {
		return nil, errwrap.Wrapf("failed to marshal config: {{err}}", err)
	}

	var clonedConfig Config
	err = proto.Unmarshal(marshaledConfig, &clonedConfig)
	if err != nil {
		return nil, errwrap.Wrapf("failed to unmarshal config: {{err}}", err)
	}

	return &clonedConfig, nil
}
