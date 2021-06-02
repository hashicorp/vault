// +build windows

package diagnose

import (
	"context"
)

func OSChecks(ctx context.Context) {
	ctx, span := StartSpan(ctx, "operating system")
	defer span.End()
	diskUsage(ctx)
}
