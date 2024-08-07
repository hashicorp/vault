// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

// BackendConfig contains configuration parameters used in the factory func to
// instantiate audit backends
type BackendConfig struct {
	// The view to store the salt
	SaltView logical.Storage

	// The salt config that should be used for any secret obfuscation
	SaltConfig *salt.Config

	// Config is the opaque user configuration provided when mounting
	Config map[string]string

	// MountPath is the path where this Backend is mounted
	MountPath string

	// Logger is used to emit log messages usually captured in the server logs.
	Logger hclog.Logger
}

// Validate ensures that we have the required configuration to create audit backends.
func (c *BackendConfig) Validate() error {
	if c.SaltConfig == nil {
		return fmt.Errorf("nil salt config: %w", ErrInvalidParameter)
	}

	if c.SaltView == nil {
		return fmt.Errorf("nil salt view: %w", ErrInvalidParameter)
	}

	if c.Logger == nil || reflect.ValueOf(c.Logger).IsNil() {
		return fmt.Errorf("nil logger: %w", ErrInvalidParameter)
	}

	if c.Config == nil {
		return fmt.Errorf("config cannot be nil: %w", ErrInvalidParameter)
	}

	if strings.TrimSpace(c.MountPath) == "" {
		return fmt.Errorf("mount path cannot be empty: %w", ErrExternalOptions)
	}

	// Validate actual config specific to Vault version (Enterprise/CE).
	if err := c.validate(); err != nil {
		return err
	}

	return nil
}
