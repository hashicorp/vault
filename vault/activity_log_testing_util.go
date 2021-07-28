package vault

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
)

// InjectActivityLogDataThisMonth populates the in-memory client store
// with some entities and tokens, overriding what was already there
// It is currently used for API integration tests
func (c *Core) InjectActivityLogDataThisMonth(t *testing.T) (map[string]struct{}, map[string]uint64) {
	t.Helper()

	activeEntities := map[string]struct{}{
		"entity0": {},
		"entity1": {},
		"entity2": {},
	}
	tokens := map[string]uint64{
		"ns0": 5,
		"ns1": 1,
		"ns2": 10,
	}

	c.activityLog.l.Lock()
	defer c.activityLog.l.Unlock()
	c.activityLog.fragmentLock.Lock()
	defer c.activityLog.fragmentLock.Unlock()

	c.activityLog.activeEntities = activeEntities
	c.activityLog.currentSegment.tokenCount.CountByNamespaceID = tokens

	return activeEntities, tokens
}

// Return the in-memory activeEntities from an activity log
func (c *Core) GetActiveEntities() map[string]struct{} {
	out := make(map[string]struct{})

	c.stateLock.RLock()
	c.activityLog.fragmentLock.RLock()
	for k, v := range c.activityLog.activeEntities {
		out[k] = v
	}
	c.activityLog.fragmentLock.RUnlock()
	c.stateLock.RUnlock()

	return out
}

// GetCurrentEntities returns the current entity activity log
func (a *ActivityLog) GetCurrentEntities() *activity.EntityActivityLog {
	a.l.RLock()
	defer a.l.RUnlock()
	return a.currentSegment.currentEntities
}

// WriteToStorage is used to put entity data in storage
// `path` should be the complete path (not relative to the view)
func WriteToStorage(t *testing.T, c *Core, path string, data []byte) {
	t.Helper()
	err := c.barrier.Put(context.Background(), &logical.StorageEntry{
		Key:   path,
		Value: data,
	})
	if err != nil {
		t.Fatalf("Failed to write %s\nto %s\nerror: %v", data, path, err)
	}
}

// SetStandbyEnable sets enabled on a performance standby (using config)
func (a *ActivityLog) SetStandbyEnable(ctx context.Context, enabled bool) {
	var enableStr string
	if enabled {
		enableStr = "enable"
	} else {
		enableStr = "disable"
	}

	// TODO only patch enabled?
	a.SetConfigStandby(ctx, activityConfig{
		DefaultReportMonths: 12,
		RetentionMonths:     24,
		Enabled:             enableStr,
	})
}

// ExpectCurrentSegmentRefreshed verifies that the current segment has been refreshed
// non-nil empty components and updated with the `expectedStart` timestamp
// Note: if `verifyTimeNotZero` is true, ignore `expectedStart` and just make sure the timestamp isn't 0
func (a *ActivityLog) ExpectCurrentSegmentRefreshed(t *testing.T, expectedStart int64, verifyTimeNotZero bool) {
	t.Helper()

	a.l.RLock()
	defer a.l.RUnlock()
	a.fragmentLock.RLock()
	defer a.fragmentLock.RUnlock()
	if a.currentSegment.currentEntities == nil {
		t.Fatalf("expected non-nil currentSegment.currentEntities")
	}
	if a.currentSegment.currentEntities.Entities == nil {
		t.Errorf("expected non-nil currentSegment.currentEntities.Entities")
	}
	if a.activeEntities == nil {
		t.Errorf("expected non-nil activeEntities")
	}
	if a.currentSegment.tokenCount == nil {
		t.Fatalf("expected non-nil currentSegment.tokenCount")
	}
	if a.currentSegment.tokenCount.CountByNamespaceID == nil {
		t.Errorf("expected non-nil currentSegment.tokenCount.CountByNamespaceID")
	}

	if len(a.currentSegment.currentEntities.Entities) > 0 {
		t.Errorf("expected no current entity segment to be loaded. got: %v", a.currentSegment.currentEntities)
	}
	if len(a.activeEntities) > 0 {
		t.Errorf("expected no active entity segment to be loaded. got: %v", a.activeEntities)
	}
	if len(a.currentSegment.tokenCount.CountByNamespaceID) > 0 {
		t.Errorf("expected no token counts to be loaded. got: %v", a.currentSegment.tokenCount.CountByNamespaceID)
	}

	if verifyTimeNotZero {
		if a.currentSegment.startTimestamp == 0 {
			t.Error("bad start timestamp. expected no reset but timestamp was reset")
		}
	} else if a.currentSegment.startTimestamp != expectedStart {
		t.Errorf("bad start timestamp. expected: %v got: %v", expectedStart, a.currentSegment.startTimestamp)
	}
}

// ActiveEntitiesEqual checks that only the set of `test` exists in `active`
func ActiveEntitiesEqual(active map[string]struct{}, test []*activity.EntityRecord) bool {
	if len(active) != len(test) {
		return false
	}

	for _, ent := range test {
		if _, ok := active[ent.EntityID]; !ok {
			return false
		}
	}

	return true
}

// GetStartTimestamp returns the start timestamp on an activity log
func (a *ActivityLog) GetStartTimestamp() int64 {
	a.l.RLock()
	defer a.l.RUnlock()
	return a.currentSegment.startTimestamp
}

// SetStartTimestamp sets the start timestamp on an activity log
func (a *ActivityLog) SetStartTimestamp(timestamp int64) {
	a.l.Lock()
	defer a.l.Unlock()
	a.currentSegment.startTimestamp = timestamp
}

// SetTokenCount sets the tokenCount on an activity log
func (a *ActivityLog) SetTokenCount(tokenCount *activity.TokenCount) {
	a.l.Lock()
	defer a.l.Unlock()
	a.currentSegment.tokenCount = tokenCount
}

// GetCountByNamespaceID returns the count of tokens by namespace ID
func (a *ActivityLog) GetCountByNamespaceID() map[string]uint64 {
	a.l.RLock()
	defer a.l.RUnlock()
	return a.currentSegment.tokenCount.CountByNamespaceID
}

// GetEntitySequenceNumber returns the current entity sequence number
func (a *ActivityLog) GetEntitySequenceNumber() uint64 {
	a.l.RLock()
	defer a.l.RUnlock()
	return a.currentSegment.entitySequenceNumber
}

// SetEnable sets the enabled flag on the activity log
func (a *ActivityLog) SetEnable(enabled bool) {
	a.l.Lock()
	defer a.l.Unlock()
	a.fragmentLock.Lock()
	defer a.fragmentLock.Unlock()
	a.enabled = enabled
}

// GetEnabled returns the enabled flag on an activity log
func (a *ActivityLog) GetEnabled() bool {
	a.fragmentLock.RLock()
	defer a.fragmentLock.RUnlock()
	return a.enabled
}

// GetActivityLog returns a pointer to the (private) activity log on a core
// Note: you must do the usual locking scheme when modifying the ActivityLog
func (c *Core) GetActivityLog() *ActivityLog {
	return c.activityLog
}
