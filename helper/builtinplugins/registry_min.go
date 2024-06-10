// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build minimal

package builtinplugins

func newRegistry() *registry {
	return newCommonRegistry()
}
