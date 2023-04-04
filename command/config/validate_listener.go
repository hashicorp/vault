// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !fips_140_3

package config

import "github.com/hashicorp/vault/internalshared/configutil"

func IsValidListener(listener *configutil.Listener) error {
	return nil
}
