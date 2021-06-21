// +build !windows

package diagnose

import (
	"context"
	"fmt"

	"golang.org/x/sys/unix"
)

func OSChecks(ctx context.Context) {
	ctx, span := StartSpan(ctx, "Operating System")
	defer span.End()

	var limit unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &limit); err != nil {
		SpotError(ctx, "open file limits", fmt.Errorf("Could not determine open file limit: %w", err))
	} else {
		min := limit.Cur
		if limit.Max < min {
			min = limit.Max
		}
		if min <= 1024 {
			SpotWarn(ctx, "open file limits", fmt.Sprintf("Open file limits are set to %d", min))
			Advise(ctx, "These limits may be insufficient. We recommend raising the soft and hard limits to 1024768.")
		} else {
			SpotOk(ctx, "open file limits", fmt.Sprintf("Open file limits are set to %d", min))
		}
	}

	diskUsage(ctx)
}
