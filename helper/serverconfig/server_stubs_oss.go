// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package server

import "github.com/hashicorp/vault/internalshared/configutil"

func entValidateConfig(_ *Config, _ string) []configutil.ConfigError {
	return nil
}
