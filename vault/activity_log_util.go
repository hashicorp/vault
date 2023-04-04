// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package vault

import "context"

// sendCurrentFragment is a no-op on OSS
func (a *ActivityLog) sendCurrentFragment(ctx context.Context) error {
	return nil
}

// CensusReport is a no-op on OSS
func (a *ActivityLog) CensusReport(_ctx context.Context, _ca *CensusAgent) {}
