package vault

import (
	"context"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// activitySubPath is the directory under the system view where
	// the log will be stored.
	activitySubPath = "activity/"
)

// ActivityLog tracks unique entity counts and non-entity token counts.
// It handles assembling log fragments (and sending them to the active
// node), writing log segments, and precomputing queries.
type ActivityLog struct {
	// log destination
	logger log.Logger

	// view is the storage location used by ActivityLog,
	// defaults to sys/activity.
	view logical.Storage
}

// NewActivityLog creates an activity log.
func NewActivityLog(_ context.Context, logger log.Logger, view logical.Storage) (*ActivityLog, error) {
	return &ActivityLog{
		logger: logger,
		view:   view,
	}, nil
}

// setupActivityLog hooks up the singleton ActivityLog into Core.
func (c *Core) setupActivityLog(ctx context.Context) error {
	view := c.systemBarrierView.SubView(activitySubPath)
	logger := c.baseLogger.Named("activity")
	c.AddLogger(logger)

	manager, err := NewActivityLog(ctx, logger, view)
	if err != nil {
		return err
	}
	c.activityLog = manager
	return nil
}
