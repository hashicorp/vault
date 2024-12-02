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
		c.activityLog.globalPartialMonthClientTracker[er.ClientID] = er
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
				c.activityLog.globalPartialMonthClientTracker[er.ClientID] = er
			}
		}
	}

	return c.activityLog.globalPartialMonthClientTracker
}

// GetActiveClients returns the in-memory globalPartialMonthClientTracker and  partialMonthLocalClientTracker from an
// activity log.
func (c *Core) GetActiveClients() map[string]*activity.EntityRecord {
	out := make(map[string]*activity.EntityRecord)

	c.stateLock.RLock()
	c.activityLog.globalFragmentLock.RLock()
	c.activityLog.localFragmentLock.RLock()

	// add active global clients
	for k, v := range c.activityLog.globalPartialMonthClientTracker {
		out[k] = v
	}

	// add active local clients
	for k, v := range c.activityLog.partialMonthLocalClientTracker {
		out[k] = v
	}

	c.activityLog.globalFragmentLock.RUnlock()
	c.activityLog.localFragmentLock.RUnlock()
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

// GetActiveLocalClientsList returns the active clients from globalPartialMonthClientTracker in activity log
func (c *Core) GetActiveGlobalClientsList() []*activity.EntityRecord {
	out := []*activity.EntityRecord{}
	c.activityLog.globalFragmentLock.RLock()
	// add active global clients
	for _, v := range c.activityLog.globalPartialMonthClientTracker {
		out = append(out, v)
	}
	c.activityLog.globalFragmentLock.RUnlock()
	return out
}

// GetActiveLocalClientsList returns the active clients from partialMonthLocalClientTracker in activity log
func (c *Core) GetActiveLocalClientsList() []*activity.EntityRecord {
	out := []*activity.EntityRecord{}
	c.activityLog.localFragmentLock.RLock()
	// add active global clients
	for _, v := range c.activityLog.partialMonthLocalClientTracker {
		out = append(out, v)
	}
	c.activityLog.localFragmentLock.RUnlock()
	return out
}

// GetCurrentGlobalEntities returns the current clients from currentGlobalSegment in activity log
func (a *ActivityLog) GetCurrentGlobalEntities() *activity.EntityActivityLog {
	a.l.RLock()
	defer a.l.RUnlock()
	return a.currentGlobalSegment.currentClients
}

// GetCurrentLocalEntities returns the current clients from currentLocalSegment in activity log
func (a *ActivityLog) GetCurrentLocalEntities() *activity.EntityActivityLog {
	a.l.RLock()
	defer a.l.RUnlock()
	return a.currentLocalSegment.currentClients
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
		RetentionMonths:     ActivityLogMinimumRetentionMonths,
		Enabled:             enableStr,
	})
}

// NOTE: AddTokenToFragment is deprecated and can no longer be used, except for
// testing backward compatibility. Please use AddClientToFragment instead.
func (a *ActivityLog) AddTokenToFragment(namespaceID string) {
	a.globalFragmentLock.Lock()
	defer a.globalFragmentLock.Unlock()

	a.localFragmentLock.Lock()
	defer a.localFragmentLock.Unlock()

	if !a.enabled {
		return
	}

	a.createCurrentFragment()

	a.localFragment.NonEntityTokens[namespaceID] += 1
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
	if a.currentGlobalSegment.currentClients == nil {
		t.Fatalf("expected non-nil currentSegment.currentClients")
	}
	if a.currentGlobalSegment.currentClients.Clients == nil {
		t.Errorf("expected non-nil currentSegment.currentClients.Entities")
	}
	if a.currentGlobalSegment.tokenCount == nil {
		t.Fatalf("expected non-nil currentSegment.tokenCount")
	}
	if a.currentGlobalSegment.tokenCount.CountByNamespaceID == nil {
		t.Errorf("expected non-nil currentSegment.tokenCount.CountByNamespaceID")
	}
	if a.currentLocalSegment.currentClients == nil {
		t.Fatalf("expected non-nil currentSegment.currentClients")
	}
	if a.currentLocalSegment.currentClients.Clients == nil {
		t.Errorf("expected non-nil currentSegment.currentClients.Entities")
	}
	if a.currentLocalSegment.tokenCount == nil {
		t.Fatalf("expected non-nil currentSegment.tokenCount")
	}
	if a.currentLocalSegment.tokenCount.CountByNamespaceID == nil {
		t.Errorf("expected non-nil currentSegment.tokenCount.CountByNamespaceID")
	}
	if a.partialMonthLocalClientTracker == nil {
		t.Errorf("expected non-nil partialMonthLocalClientTracker")
	}
	if a.globalPartialMonthClientTracker == nil {
		t.Errorf("expected non-nil globalPartialMonthClientTracker")
	}
	if len(a.currentGlobalSegment.currentClients.Clients) > 0 {
		t.Errorf("expected no current entity segment to be loaded. got: %v", a.currentGlobalSegment.currentClients)
	}
	if len(a.currentLocalSegment.currentClients.Clients) > 0 {
		t.Errorf("expected no current entity segment to be loaded. got: %v", a.currentLocalSegment.currentClients)
	}
	if len(a.currentLocalSegment.tokenCount.CountByNamespaceID) > 0 {
		t.Errorf("expected no token counts to be loaded. got: %v", a.currentLocalSegment.tokenCount.CountByNamespaceID)
	}
	if len(a.partialMonthLocalClientTracker) > 0 {
		t.Errorf("expected no active entity segment to be loaded. got: %v", a.partialMonthLocalClientTracker)
	}
	if len(a.globalPartialMonthClientTracker) > 0 {
		t.Errorf("expected no active entity segment to be loaded. got: %v", a.globalPartialMonthClientTracker)
	}

	if verifyTimeNotZero {
		if a.currentGlobalSegment.startTimestamp == 0 {
			t.Error("bad start timestamp. expected no reset but timestamp was reset")
		}
		if a.currentLocalSegment.startTimestamp == 0 {
			t.Error("bad start timestamp. expected no reset but timestamp was reset")
		}
	} else if a.currentGlobalSegment.startTimestamp != expectedStart {
		t.Errorf("bad start timestamp. expected: %v got: %v", expectedStart, a.currentGlobalSegment.startTimestamp)
	} else if a.currentLocalSegment.startTimestamp != expectedStart {
		t.Errorf("bad start timestamp. expected: %v got: %v", expectedStart, a.currentLocalSegment.startTimestamp)
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
	if a.currentGlobalSegment.startTimestamp != a.currentLocalSegment.startTimestamp {
		return -1
	}
	return a.currentGlobalSegment.startTimestamp
}

// SetStartTimestamp sets the start timestamp on an activity log
func (a *ActivityLog) SetStartTimestamp(timestamp int64) {
	a.l.Lock()
	defer a.l.Unlock()
	a.currentGlobalSegment.startTimestamp = timestamp
	a.currentLocalSegment.startTimestamp = timestamp
}

// GetStoredTokenCountByNamespaceID returns the count of tokens by namespace ID
func (a *ActivityLog) GetStoredTokenCountByNamespaceID() map[string]uint64 {
	a.l.RLock()
	defer a.l.RUnlock()
	return a.currentLocalSegment.tokenCount.CountByNamespaceID
}

// GetGlobalEntitySequenceNumber returns the current entity sequence number
func (a *ActivityLog) GetGlobalEntitySequenceNumber() uint64 {
	a.l.RLock()
	defer a.l.RUnlock()
	return a.currentGlobalSegment.clientSequenceNumber
}

// GetLocalEntitySequenceNumber returns the current entity sequence number
func (a *ActivityLog) GetLocalEntitySequenceNumber() uint64 {
	a.l.RLock()
	defer a.l.RUnlock()
	return a.currentLocalSegment.clientSequenceNumber
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

func (c *Core) GetActiveGlobalFragment() *activity.LogFragment {
	c.activityLog.globalFragmentLock.RLock()
	defer c.activityLog.globalFragmentLock.RUnlock()
	return c.activityLog.currentGlobalFragment
}

func (c *Core) GetSecondaryGlobalFragments() []*activity.LogFragment {
	c.activityLog.globalFragmentLock.RLock()
	defer c.activityLog.globalFragmentLock.RUnlock()
	return c.activityLog.secondaryGlobalClientFragments
}

func (c *Core) GetActiveLocalFragment() *activity.LogFragment {
	c.activityLog.localFragmentLock.RLock()
	defer c.activityLog.localFragmentLock.RUnlock()
	return c.activityLog.localFragment
}

// StoreCurrentSegment is a test only method to create and store
// segments from fragments. This allows createCurrentSegmentFromFragments to remain
// private
func (c *Core) StoreCurrentSegment(ctx context.Context, fragments []*activity.LogFragment, currentSegment *segmentInfo, force bool, storagePathPrefix string) error {
	return c.activityLog.createCurrentSegmentFromFragments(ctx, fragments, currentSegment, force, storagePathPrefix)
}

// DeleteLogsAtPath is test helper function deletes all logs at the given path
func (c *Core) DeleteLogsAtPath(ctx context.Context, t *testing.T, storagePath string, startTime int64) {
	basePath := storagePath + fmt.Sprint(startTime) + "/"
	a := c.activityLog
	segments, err := a.view.List(ctx, basePath)
	if err != nil {
		t.Fatalf("could not list path %v", err)
		return
	}
	for _, p := range segments {
		err = a.view.Delete(ctx, basePath+p)
		if err != nil {
			t.Fatalf("could not delete path %v", err)
		}
	}
}
