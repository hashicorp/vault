// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package pki

func setupEntSpecificBackend(_ *backend) {
	// ENT hook is not used by OSS.
}
