//go:build !windows

package diagnose

import (
	"context"
	"fmt"

	"golang.org/x/sys/unix"
)

func OSChecks(ctx context.Context) {
	ctx, span := StartSpan(ctx, "Check Operating System")
	defer span.End()

	fileLimitsName := "Check Open File Limits"

	var limit unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &limit); err != nil {
		SpotError(ctx, fileLimitsName, fmt.Errorf("Could not determine open file limit: %w.", err))
	} else {
		min := limit.Cur
		if limit.Max < min {
			min = limit.Max
		}
		if min <= 1024 {
			SpotWarn(ctx, fileLimitsName, fmt.Sprintf("Open file limits are set to %d", min),
				Advice("These limits may be insufficient. We recommend raising the soft and hard limits to 1024768."))
		} else {
			SpotOk(ctx, fileLimitsName, fmt.Sprintf("Open file limits are set to %d.", min))
		}
	}

	diskUsage(ctx)
}
