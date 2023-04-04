// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package server

import (
	"errors"
	"github.com/hashicorp/hcl/hcl/ast"
)

type entConfig struct{}

func (ec *entConfig) parseConfig(list *ast.ObjectList) error {
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

	disabledSeals := 0
	for _, seal := range c.Seals {
		if seal.Disabled {
			disabledSeals++
		}
	}

	if disabledSeals == len(c.Seals) {
		return errors.New("seals: multiple seals provided but all are disabled")
	}

	if disabledSeals < len(c.Seals)-1 {
		return errors.New("seals: only one seal can be enabled")
	}

	return nil
}
