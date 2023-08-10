// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build cgo

package version

func init() {
	CgoEnabled = true
}
