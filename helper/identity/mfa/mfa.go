// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

func (c *MFAEnforcementConfig) Clone() (*MFAEnforcementConfig, error) {
	if c == nil {
		return nil, fmt.Errorf("nil config")
	}

	marshaledConfig, err := proto.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	var clonedConfig MFAEnforcementConfig
	err = proto.Unmarshal(marshaledConfig, &clonedConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &clonedConfig, nil
}
