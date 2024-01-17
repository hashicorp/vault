// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
	"google.golang.org/protobuf/testing/protocmp"
)

// InjectActivityLogDataThisMonth populates the in-memory client store
// with some entities and tokens, overriding what was already there
// It is currently used for API integration tests
func (c *Core) InjectActivityLogDataThisMonth(t *testing.T) map[string]*activity.EntityRecord {
	t.Helper()

	c.activityLog.l.Lock()
	defer c.activityLog.l.Unlock()
	c.activityLog.fragmentLock.Lock()
	defer c.activityLog.fragmentLock.Unlock()

	for i := 0; i < 3; i++ {
		er := &activity.EntityRecord{
			ClientID:      fmt.Sprintf("testclientid-%d", i),
			NamespaceID:   "root",
			MountAccessor: fmt.Sprintf("testmountaccessor-%d", i),
			Timestamp:     c.activityLog.clock.Now().Unix(),
			NonEntity:     i%2 == 0,
		}
		c.activityLog.partialMonthClientTracker[er.ClientID] = er
	}

	if constants.IsEnterprise {
		for j := 0; j < 2; j++ {
			for i := 0; i < 2; i++ {
				er := &activity.EntityRecord{
					ClientID:      fmt.Sprintf("ns-%d-testclientid-%d", j, i),
					NamespaceID:   fmt.Sprintf("ns-%d", j),
					MountAccessor: fmt.Sprintf("ns-%d-testmountaccessor-%d", j, i),
					Timestamp:     c.activityLog.clock.Now().Unix(),
					NonEntity:     i%2 == 0,
				}
				c.activityLog.partialMonthClientTracker[er.ClientID] = er
			}
		}
	}

	return c.activityLog.partialMonthClientTracker
}

// GetActiveClients returns the in-memory partialMonthClientTracker from an
// activity log.
func (c *Core) GetActiveClients() map[string]*activity.EntityRecord {
	out := make(map[string]*activity.EntityRecord)

	c.stateLock.RLock()
	c.activityLog.fragmentLock.RLock()
	for k, v := range c.activityLog.partialMonthClientTracker {
		out[k] = v
	}
	c.activityLog.fragmentLock.RUnlock()
	c.stateLock.RUnlock()

	return out
}

func (c *Core) GetActiveClientsList() []*activity.EntityRecord {
	out := []*activity.EntityRecord{}

	for _, v := range c.GetActiveClients() {
		out = append(out, v)
	}

	return out
}

// GetCurrentEntities returns the current entity activity log
func (a *ActivityLog) GetCurrentEntities() *activity.EntityActivityLog {
	a.l.RLock()
	defer a.l.RUnlock()
	return a.currentSegment.currentClients
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

// NOTE: AddTokenToFragment is deprecated and can no longer be used, except for
// testing backward compatibility. Please use AddClientToFragment instead.
func (a *ActivityLog) AddTokenToFragment(namespaceID string) {
	a.fragmentLock.Lock()
	defer a.fragmentLock.Unlock()

	if !a.enabled {
		return
	}

	a.createCurrentFragment()

	a.fragment.NonEntityTokens[namespaceID] += 1
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rand.Intn(26)) + 'a'
	}
	return string(b)
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
	if a.currentSegment.currentClients == nil {
		t.Fatalf("expected non-nil currentSegment.currentClients")
	}
	if a.currentSegment.currentClients.Clients == nil {
		t.Errorf("expected non-nil currentSegment.currentClients.Entities")
	}
	if a.currentSegment.tokenCount == nil {
		t.Fatalf("expected non-nil currentSegment.tokenCount")
	}
	if a.currentSegment.tokenCount.CountByNamespaceID == nil {
		t.Errorf("expected non-nil currentSegment.tokenCount.CountByNamespaceID")
	}
	if a.partialMonthClientTracker == nil {
		t.Errorf("expected non-nil partialMonthClientTracker")
	}
	if len(a.currentSegment.currentClients.Clients) > 0 {
		t.Errorf("expected no current entity segment to be loaded. got: %v", a.currentSegment.currentClients)
	}
	if len(a.currentSegment.tokenCount.CountByNamespaceID) > 0 {
		t.Errorf("expected no token counts to be loaded. got: %v", a.currentSegment.tokenCount.CountByNamespaceID)
	}
	if len(a.partialMonthClientTracker) > 0 {
		t.Errorf("expected no active entity segment to be loaded. got: %v", a.partialMonthClientTracker)
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
func ActiveEntitiesEqual(active []*activity.EntityRecord, test []*activity.EntityRecord) error {
	opts := []cmp.Option{protocmp.Transform(), cmpopts.SortSlices(func(x, y *activity.EntityRecord) bool {
		return x.ClientID < y.ClientID
	})}
	if diff := cmp.Diff(active, test, opts...); len(diff) > 0 {
		return fmt.Errorf("entity record mismatch: %v", diff)
	}
	return nil
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

// GetStoredTokenCountByNamespaceID returns the count of tokens by namespace ID
func (a *ActivityLog) GetStoredTokenCountByNamespaceID() map[string]uint64 {
	a.l.RLock()
	defer a.l.RUnlock()
	return a.currentSegment.tokenCount.CountByNamespaceID
}

// GetEntitySequenceNumber returns the current entity sequence number
func (a *ActivityLog) GetEntitySequenceNumber() uint64 {
	a.l.RLock()
	defer a.l.RUnlock()
	return a.currentSegment.clientSequenceNumber
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
