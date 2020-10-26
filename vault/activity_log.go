package vault

import (
	"context"
	"os"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
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

	// nodeID is the ID to use for all fragments that
	// are generated.
	// TODO: use secondary ID when available?
	nodeID string

	// current log fragment (may be nil) and a mutex to protect it
	fragmentLock     sync.RWMutex
	fragment         *activity.LogFragment
	fragmentCreation time.Time
}

// NewActivityLog creates an activity log.
func NewActivityLog(_ context.Context, logger log.Logger, view logical.Storage) (*ActivityLog, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &ActivityLog{
		logger: logger,
		view:   view,
		nodeID: hostname,
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

func (a *ActivityLog) AddEntityToFragment(entityID string, namespaceID string, timestamp time.Time) {
	a.fragmentLock.Lock()
	defer a.fragmentLock.Unlock()

	// TODO: check whether entity ID already recorded

	a.createCurrentFragment()

	a.fragment.Entities = append(a.fragment.Entities,
		&activity.EntityRecord{
			EntityID:    entityID,
			NamespaceID: namespaceID,
			Timestamp:   timestamp.UnixNano(),
		})
}

func (a *ActivityLog) AddTokenToFragment(namespaceID string) {
	a.fragmentLock.Lock()
	defer a.fragmentLock.Unlock()

	a.createCurrentFragment()

	a.fragment.NonEntityTokens[namespaceID] += 1
}

// Create the current fragment if it doesn't already exist.
// Must be called with the lock held.
func (a *ActivityLog) createCurrentFragment() {
	if a.fragment == nil {
		a.fragment = &activity.LogFragment{
			OriginatingNode: a.nodeID,
			Entities:        make([]*activity.EntityRecord, 0, 120),
			NonEntityTokens: make(map[string]uint64),
		}
		a.fragmentCreation = time.Now()

		// TODO: start a timer to send it, if we're a performance standby
	}
}
