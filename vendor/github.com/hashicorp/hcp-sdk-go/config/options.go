// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

// HCPConfigOption are functions that modify the hcpConfig struct. They can be
// passed to NewHCPConfig.
type HCPConfigOption = func(config *hcpConfig) error

// apply can be used to directly apply an option. This can be helpful to call
// an option from another one or during tests.
//
// E.g. apply(config, WithClientCredentials(clientID, clientSecret))
func apply(config *hcpConfig, option HCPConfigOption) error {
	return option(config)
}
