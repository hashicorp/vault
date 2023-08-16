//go:build !enterprise

package vault

import (
	"context"
	"time"
)

// sendCurrentFragment is a no-op on OSS
func (a *ActivityLog) sendCurrentFragment(ctx context.Context) error {
	return nil
}

// CensusReport is a no-op on OSS
func (a *ActivityLog) CensusReport(context.Context, CensusReporter, time.Time) {}
