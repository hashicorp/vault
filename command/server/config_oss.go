// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package server

func (c *Config) IsMultisealEnabled() bool {
	return false
}
