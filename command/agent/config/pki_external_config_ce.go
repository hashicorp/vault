// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package config

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/hcl/ast"
)

// PKIExternalCA is a CE stub; pki_external_ca is an enterprise-only feature.
type PKIExternalCA struct{}

// validatePKIExternalCAConfig is a no-op in CE builds.
func (c *Config) validatePKIExternalCAConfig(_ hclog.Logger) error {
	return nil
}

// parsePKIExternalCA is a no-op in CE builds.
func parsePKIExternalCA(_ *Config, _ *ast.ObjectList) error {
	return nil
}
