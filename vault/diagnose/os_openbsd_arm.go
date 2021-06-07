// +build openbsd,arm

package diagnose

import "context"

func diskUsage(ctx context.Context) error {
	return SpotSkipped(ctx, "disk usage", "unsupported on this platform")
}
