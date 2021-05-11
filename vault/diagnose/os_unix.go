// +build !windows

package diagnose

import (
	"context"
	"fmt"
	"golang.org/x/sys/unix"
)

func OSChecks(ctx context.Context) {
	var limit unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &limit); err != nil {
		SpotError(ctx, "open file limits", fmt.Errorf("could not determine open file limit: %w", err))
	} else {
		if limit.Cur <= 1024 || limit.Max <= 1024 {
			SpotWarn(ctx, "open file limits", fmt.Sprintf("open file limits are set to %d, which may be insufficient.", limit.Max))
		} else {
			SpotOk(ctx, "open file limits", fmt.Sprintf("open file limits are set to %d", limit.Max))
		}
	}
}
