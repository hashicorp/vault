// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package server

import "github.com/hashicorp/vault/internalshared/configutil"

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func entValidateConfig(_ *Config, _ string) []configutil.ConfigError {
	return nil
}
