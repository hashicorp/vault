// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build cgo

package version

func init() {
	CgoEnabled = true
}
