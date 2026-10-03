// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build windows

package diagnose

import (
	"context"
)

func OSChecks(ctx context.Context) {
	ctx, span := StartSpan(ctx, "Check Operating System")
	defer span.End()
	diskUsage(ctx)
}
