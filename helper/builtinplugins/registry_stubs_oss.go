// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package builtinplugins

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func entAddExtPlugins(r *registry) {
}
