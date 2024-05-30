// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package audit

import "fmt"

// validate ensures that this if we're not running Vault Enterprise, we cannot
// supply Enterprise-only audit configuration options.
func (c *BackendConfig) validate() error {
	if HasInvalidOptions(c.Config) {
		return fmt.Errorf("enterprise-only options supplied: %w", ErrExternalOptions)
	}

	return nil
}
