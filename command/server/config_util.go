// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package server

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/hcl/ast"
)

type entConfig struct{}

func (ec *entConfig) parseConfig(list *ast.ObjectList, source string) error {
	return nil
}

func (ec entConfig) Merge(ec2 entConfig) entConfig {
	result := entConfig{}
	return result
}

func (ec entConfig) Sanitized() map[string]interface{} {
	return nil
}

func (c *Config) checkSealConfig() error {
	if len(c.Seals) == 0 {
		return nil
	}

	if len(c.Seals) > 2 {
		return fmt.Errorf("seals: at most 2 seals can be provided: received %d", len(c.Seals))
	}

	disabledSeals := 0
	for _, seal := range c.Seals {
		if seal.Disabled {
			disabledSeals++
		}
	}

	if len(c.Seals) > 1 && disabledSeals == len(c.Seals) {
		return errors.New("seals: seals provided but all are disabled")
	}

	if disabledSeals < len(c.Seals)-1 {
		return errors.New("seals: only one seal can be enabled")
	}

	return nil
}
