// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
	"github.com/stretchr/testify/require"
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

// ExpectOldSegmentRefreshed verifies that the old current segment structure has been refreshed
// non-nil empty components and updated with the `expectedStart` timestamp. This is expected when
// an upgrade has not yet completed.
// Note: if `verifyTimeNotZero` is true, ignore `expectedStart` and just make sure the timestamp isn't 0
func (a *ActivityLog) ExpectOldSegmentRefreshed(t *testing.T, expectedStart int64, verifyTimeNotZero bool, expectedEntities []*activity.EntityRecord, directTokens map[string]uint64) {
	t.Helper()

	a.l.RLock()
	defer a.l.RUnlock()
	a.fragmentLock.RLock()
	defer a.fragmentLock.RUnlock()
	require.NotNil(t, a.currentSegment.currentClients)
	require.NotNil(t, a.currentSegment.currentClients.Clients)
	require.NotNil(t, a.currentSegment.tokenCount)
	require.NotNil(t, a.currentSegment.tokenCount.CountByNamespaceID)
	if !EntityRecordsEqual(t, a.currentSegment.currentClients.Clients, expectedEntities) {
		// we only expect the newest entity segment to be loaded (for the current month)
		t.Errorf("bad activity entity logs loaded. expected: %v got: %v", a.currentSegment.currentClients.Clients, expectedEntities)
	}
	require.Equal(t, directTokens, a.currentSegment.tokenCount.CountByNamespaceID)
	if verifyTimeNotZero {
		require.NotEqual(t, a.currentSegment.startTimestamp, 0)
	} else {
		require.Equal(t, a.currentSegment.startTimestamp, expectedStart)
	}
}

// ExpectCurrentSegmentsRefreshed verifies that the current segment has been refreshed
// non-nil empty components and updated with the `expectedStart` timestamp
// Note: if `verifyTimeNotZero` is true, ignore `expectedStart` and just make sure the timestamp isn't 0
func (a *ActivityLog) ExpectCurrentSegmentsRefreshed(t *testing.T, expectedStart int64, verifyTimeNotZero bool) {
	t.Helper()

	a.l.RLock()
	defer a.l.RUnlock()
	a.fragmentLock.RLock()
	defer a.fragmentLock.RUnlock()
	require.NotNil(t, a.currentGlobalSegment.currentClients)
	require.NotNil(t, a.currentGlobalSegment.currentClients.Clients)
	require.NotNil(t, a.currentGlobalSegment.tokenCount)
	require.NotNil(t, a.currentGlobalSegment.tokenCount.CountByNamespaceID)

	require.NotNil(t, a.currentLocalSegment.currentClients)
	require.NotNil(t, a.currentLocalSegment.currentClients.Clients)
	require.NotNil(t, a.currentLocalSegment.tokenCount)
	require.NotNil(t, a.currentLocalSegment.tokenCount.CountByNamespaceID)

	require.NotNil(t, a.partialMonthLocalClientTracker)
	require.NotNil(t, a.globalPartialMonthClientTracker)

	require.Equal(t, 0, len(a.currentGlobalSegment.currentClients.Clients))
	require.Equal(t, 0, len(a.currentLocalSegment.currentClients.Clients))
	require.Equal(t, 0, len(a.currentLocalSegment.tokenCount.CountByNamespaceID))

	require.Equal(t, 0, len(a.partialMonthLocalClientTracker))
	require.Equal(t, 0, len(a.globalPartialMonthClientTracker))

	if verifyTimeNotZero {
		require.NotEqual(t, 0, a.currentGlobalSegment.startTimestamp)
		require.NotEqual(t, 0, a.currentLocalSegment.startTimestamp)
		require.NotEqual(t, 0, a.currentSegment.startTimestamp)
	} else {
		require.Equal(t, expectedStart, a.currentGlobalSegment.startTimestamp)
		require.Equal(t, expectedStart, a.currentLocalSegment.startTimestamp)
	}
}

// EntityRecordsEqual compares the parts we care about from two activity entity record slices
// note: this makes a copy of the []*activity.EntityRecord so that misordered slices won't fail the comparison,
// but the function won't modify the order of the slices to compare
func EntityRecordsEqual(t *testing.T, record1, record2 []*activity.EntityRecord) bool {
	t.Helper()

	if record1 == nil {
		return record2 == nil
	}
	if record2 == nil {
		return record1 == nil
	}

	if len(record1) != len(record2) {
		return false
	}

	// sort first on namespace, then on ID, then on timestamp
	entityLessFn := func(e []*activity.EntityRecord, i, j int) bool {
		ei := e[i]
		ej := e[j]

		nsComp := strings.Compare(ei.NamespaceID, ej.NamespaceID)
		if nsComp == -1 {
			return true
		}
		if nsComp == 1 {
			return false
		}

		idComp := strings.Compare(ei.ClientID, ej.ClientID)
		if idComp == -1 {
			return true
		}
		if idComp == 1 {
			return false
		}

		return ei.Timestamp < ej.Timestamp
	}

	entitiesCopy1 := make([]*activity.EntityRecord, len(record1))
	entitiesCopy2 := make([]*activity.EntityRecord, len(record2))
	copy(entitiesCopy1, record1)
	copy(entitiesCopy2, record2)

	sort.Slice(entitiesCopy1, func(i, j int) bool {
		return entityLessFn(entitiesCopy1, i, j)
	})
	sort.Slice(entitiesCopy2, func(i, j int) bool {
		return entityLessFn(entitiesCopy2, i, j)
	})

	for i, a := range entitiesCopy1 {
		b := entitiesCopy2[i]
		if a.ClientID != b.ClientID || a.NamespaceID != b.NamespaceID || a.Timestamp != b.Timestamp {
			return false
		}
	}

	return true
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
	a.currentSegment.startTimestamp = timestamp
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

// SaveEntitySegment is a test helper function to keep the savePreviousEntitySegments function internal
func (a *ActivityLog) SaveEntitySegment(ctx context.Context, startTime int64, pathPrefix string, fragments []*activity.LogFragment) error {
	return a.savePreviousEntitySegments(ctx, startTime, pathPrefix, fragments)
}

// LaunchMigrationWorker is a test only helper function that launches the migration workers.
// This allows us to keep the migration worker methods internal
func (a *ActivityLog) LaunchMigrationWorker(ctx context.Context, isSecondary bool) {
	if isSecondary {
		go a.core.secondaryDuplicateClientMigrationWorker(ctx)
	} else {
		go a.core.primaryDuplicateClientMigrationWorker(ctx)
	}
}

// DedupUpgradeComplete is a test helper function that indicates whether the
// all correct states have been set after completing upgrade processes to 1.19+
func (a *ActivityLog) DedupUpgradeComplete(ctx context.Context) bool {
	return a.hasDedupClientsUpgrade(ctx)
}

// ResetDedupUpgrade is a test helper function that resets the state to reflect
// how the system should look before running/completing any upgrade process to 1.19+
func (a *ActivityLog) ResetDedupUpgrade(ctx context.Context) {
	a.view.Delete(ctx, activityDeduplicationUpgradeKey)
	a.view.Delete(ctx, activitySecondaryDataRecCount)
}

// RefreshActivityLog is a test helper functions that refreshes the activity logs
// segments and current month data. This allows us to keep the refreshFromStoredLog
// function internal
func (a *ActivityLog) RefreshActivityLog(ctx context.Context) {
	a.refreshFromStoredLog(ctx, &sync.WaitGroup{}, time.Now().UTC())
}
