// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !fips_140_3

package cliconfig

func IsValidListener(tlsDisable bool) error {
	return nil
}
