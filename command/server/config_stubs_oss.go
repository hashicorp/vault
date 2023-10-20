// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package server

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func IsSealHABetaEnabled() (bool, error) {
	return false, nil
}
