// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package server

import "fmt"

func (c *Config) IsMultisealEnabled() bool {
	return false
}

func validateFeatureFlags(featureFlags []string) error {
	if len(featureFlags) > 0 {
		return fmt.Errorf("feature_flags are enterprise-only")
	}
	return nil
}
