// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build openbsd && arm

package diagnose

import "context"

func diskUsage(ctx context.Context) error {
	SpotSkipped(ctx, "Check Disk Usage", "Disk Usage diagnostics are unsupported on this platform.")
	return nil
}
