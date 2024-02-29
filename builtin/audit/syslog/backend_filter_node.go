// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package syslog

// configureFilterNode is used to configure a filter node and associated ID on the Backend.
func (b *Backend) configureFilterNode(_ string) error {
	return nil
}
