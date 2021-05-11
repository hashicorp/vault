package mfa

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

func (c *Config) Clone() (*Config, error) {
	if c == nil {
		return nil, fmt.Errorf("nil config")
	}

	marshaledConfig, err := proto.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	var clonedConfig Config
	err = proto.Unmarshal(marshaledConfig, &clonedConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &clonedConfig, nil
}
