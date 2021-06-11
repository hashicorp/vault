// +build !windows

package diagnose

import (
	"context"
	"fmt"

	"golang.org/x/sys/unix"
)

func OSChecks(ctx context.Context) {
	ctx, span := StartSpan(ctx, "operating system")
	defer span.End()

	var limit unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &limit); err != nil {
		SpotError(ctx, "open file limits", fmt.Errorf("could not determine open file limit: %w", err))
	} else {
		min := limit.Cur
		if limit.Max < min {
			min = limit.Max
		}
		if min <= 1024 {
			SpotWarn(ctx, "open file limits", fmt.Sprintf("set to %d, which may be insufficient.", min))
		} else {
			SpotOk(ctx, "open file limits", fmt.Sprintf("set to %d", min))
		}
	}

	diskUsage(ctx)
}
