// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/axiomhq/hyperloglog"
	"github.com/go-test/deep"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

// TestActivityLog_Creation calls AddEntityToFragment and verifies that it appears correctly in a.fragment.
func TestActivityLog_Creation(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)

	a := core.activityLog
	a.SetEnable(true)

	if a == nil {
		t.Fatal("no activity log found")
	}
	if a.logger == nil || a.view == nil {
		t.Fatal("activity log not initialized")
	}
	if a.fragment != nil {
		t.Fatal("activity log already has fragment")
	}

	const entity_id = "entity_id_75432"
	const namespace_id = "ns123"
	ts := time.Now()

	a.AddEntityToFragment(entity_id, namespace_id, ts.Unix())
	if a.fragment == nil {
		t.Fatal("no fragment created")
	}

	if a.fragment.OriginatingNode != a.nodeID {
		t.Errorf("mismatched node ID, %q vs %q", a.fragment.OriginatingNode, a.nodeID)
	}

	if a.fragment.Clients == nil {
		t.Fatal("no fragment entity slice")
	}

	if a.fragment.NonEntityTokens == nil {
		t.Fatal("no fragment token map")
	}

	if len(a.fragment.Clients) != 1 {
		t.Fatalf("wrong number of entities %v", len(a.fragment.Clients))
	}

	er := a.fragment.Clients[0]
	if er.ClientID != entity_id {
		t.Errorf("mimatched entity ID, %q vs %q", er.ClientID, entity_id)
	}
	if er.NamespaceID != namespace_id {
		t.Errorf("mimatched namespace ID, %q vs %q", er.NamespaceID, namespace_id)
	}
	if er.Timestamp != ts.Unix() {
		t.Errorf("mimatched timestamp, %v vs %v", er.Timestamp, ts.Unix())
	}

	// Reset and test the other code path
	a.fragment = nil
	a.AddTokenToFragment(namespace_id)

	if a.fragment == nil {
		t.Fatal("no fragment created")
	}

	if a.fragment.NonEntityTokens == nil {
		t.Fatal("no fragment token map")
	}

	actual := a.fragment.NonEntityTokens[namespace_id]
	if actual != 1 {
		t.Errorf("mismatched number of tokens, %v vs %v", actual, 1)
	}
}

// TestActivityLog_Creation_WrappingTokens calls HandleTokenUsage for two wrapping tokens, and verifies that this
// doesn't create a fragment.
func TestActivityLog_Creation_WrappingTokens(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)

	a := core.activityLog
	a.SetEnable(true)

	if a == nil {
		t.Fatal("no activity log found")
	}
	if a.logger == nil || a.view == nil {
		t.Fatal("activity log not initialized")
	}
	a.fragmentLock.Lock()
	if a.fragment != nil {
		t.Fatal("activity log already has fragment")
	}
	a.fragmentLock.Unlock()
	const namespace_id = "ns123"

	te := &logical.TokenEntry{
		Path:         "test",
		Policies:     []string{responseWrappingPolicyName},
		CreationTime: time.Now().Unix(),
		TTL:          3600,
		NamespaceID:  namespace_id,
	}

	id, isTWE := te.CreateClientID()
	err := a.HandleTokenUsage(context.Background(), te, id, isTWE)
	if err != nil {
		t.Fatal(err)
	}

	a.fragmentLock.Lock()
	if a.fragment != nil {
		t.Fatal("fragment created")
	}
	a.fragmentLock.Unlock()

	teNew := &logical.TokenEntry{
		Path:         "test",
		Policies:     []string{controlGroupPolicyName},
		CreationTime: time.Now().Unix(),
		TTL:          3600,
		NamespaceID:  namespace_id,
	}

	id, isTWE = teNew.CreateClientID()
	err = a.HandleTokenUsage(context.Background(), teNew, id, isTWE)
	if err != nil {
		t.Fatal(err)
	}

	a.fragmentLock.Lock()
	if a.fragment != nil {
		t.Fatal("fragment created")
	}
	a.fragmentLock.Unlock()
}

func checkExpectedEntitiesInMap(t *testing.T, a *ActivityLog, entityIDs []string) {
	t.Helper()

	activeClients := a.core.GetActiveClients()
	if len(activeClients) != len(entityIDs) {
		t.Fatalf("mismatched number of entities, expected %v got %v", len(entityIDs), activeClients)
	}
	for _, e := range entityIDs {
		if _, present := activeClients[e]; !present {
			t.Errorf("entity ID %q is missing", e)
		}
	}
}

// TestActivityLog_UniqueEntities calls AddEntityToFragment 4 times with 2 different clients, then verifies that there
// are only 2 clients in the fragment and that they have the earlier timestamps.
func TestActivityLog_UniqueEntities(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	a.SetEnable(true)

	id1 := "11111111-1111-1111-1111-111111111111"
	t1 := time.Now()

	id2 := "22222222-2222-2222-2222-222222222222"
	t2 := time.Now()
	t3 := t2.Add(60 * time.Second)

	a.AddEntityToFragment(id1, "root", t1.Unix())
	a.AddEntityToFragment(id2, "root", t2.Unix())
	a.AddEntityToFragment(id2, "root", t3.Unix())
	a.AddEntityToFragment(id1, "root", t3.Unix())

	if a.fragment == nil {
		t.Fatal("no current fragment")
	}

	if len(a.fragment.Clients) != 2 {
		t.Fatalf("number of entities is %v", len(a.fragment.Clients))
	}

	for i, e := range a.fragment.Clients {
		expectedID := id1
		expectedTime := t1.Unix()
		expectedNS := "root"
		if i == 1 {
			expectedID = id2
			expectedTime = t2.Unix()
		}
		if e.ClientID != expectedID {
			t.Errorf("%v: expected %q, got %q", i, expectedID, e.ClientID)
		}
		if e.NamespaceID != expectedNS {
			t.Errorf("%v: expected %q, got %q", i, expectedNS, e.NamespaceID)
		}
		if e.Timestamp != expectedTime {
			t.Errorf("%v: expected %v, got %v", i, expectedTime, e.Timestamp)
		}
	}

	checkExpectedEntitiesInMap(t, a, []string{id1, id2})
}

func readSegmentFromStorageNil(t *testing.T, c *Core, path string) {
	t.Helper()
	logSegment, err := c.barrier.Get(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if logSegment != nil {
		t.Fatalf("expected non-nil log segment at %q", path)
	}
}

func readSegmentFromStorage(t *testing.T, c *Core, path string) *logical.StorageEntry {
	t.Helper()
	logSegment, err := c.barrier.Get(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if logSegment == nil {
		t.Fatalf("expected non-nil log segment at %q", path)
	}

	return logSegment
}

func expectMissingSegment(t *testing.T, c *Core, path string) {
	t.Helper()
	logSegment, err := c.barrier.Get(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if logSegment != nil {
		t.Fatalf("expected nil log segment at %q", path)
	}
}

func expectedEntityIDs(t *testing.T, out *activity.EntityActivityLog, ids []string) {
	t.Helper()

	if len(out.Clients) != len(ids) {
		t.Fatalf("entity log expected length %v, actual %v", len(ids), len(out.Clients))
	}

	// Double loop, OK for small cases
	for _, id := range ids {
		found := false
		for _, e := range out.Clients {
			if e.ClientID == id {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("did not find entity ID %v", id)
		}
	}
}

// TestActivityLog_SaveTokensToStorage calls AddTokenToFragment with duplicate namespaces and then saves the segment to
// storage. The test then reads and unmarshals the segment, and verifies that the results have the correct counts by
// namespace.
func TestActivityLog_SaveTokensToStorage(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ctx := context.Background()

	a := core.activityLog
	a.SetStandbyEnable(ctx, true)
	a.SetStartTimestamp(time.Now().Unix()) // set a nonzero segment

	nsIDs := [...]string{"ns1_id", "ns2_id", "ns3_id"}
	path := fmt.Sprintf("%sdirecttokens/%d/0", ActivityLogPrefix, a.GetStartTimestamp())

	for i := 0; i < 3; i++ {
		a.AddTokenToFragment(nsIDs[0])
	}
	a.AddTokenToFragment(nsIDs[1])
	err := a.saveCurrentSegmentToStorage(ctx, false)
	if err != nil {
		t.Fatalf("got error writing tokens to storage: %v", err)
	}
	if a.fragment != nil {
		t.Errorf("fragment was not reset after write to storage")
	}

	out := &activity.TokenCount{}
	protoSegment := readSegmentFromStorage(t, core, path)
	err = proto.Unmarshal(protoSegment.Value, out)
	if err != nil {
		t.Fatalf("could not unmarshal protobuf: %v", err)
	}

	if len(out.CountByNamespaceID) != 2 {
		t.Fatalf("unexpected token length. Expected %d, got %d", 2, len(out.CountByNamespaceID))
	}
	for i := 0; i < 2; i++ {
		if _, ok := out.CountByNamespaceID[nsIDs[i]]; !ok {
			t.Fatalf("namespace ID %s missing from token counts", nsIDs[i])
		}
	}
	if out.CountByNamespaceID[nsIDs[0]] != 3 {
		t.Errorf("namespace ID %s has %d count, expected %d", nsIDs[0], out.CountByNamespaceID[nsIDs[0]], 3)
	}
	if out.CountByNamespaceID[nsIDs[1]] != 1 {
		t.Errorf("namespace ID %s has %d count, expected %d", nsIDs[1], out.CountByNamespaceID[nsIDs[1]], 1)
	}

	a.AddTokenToFragment(nsIDs[0])
	a.AddTokenToFragment(nsIDs[2])
	err = a.saveCurrentSegmentToStorage(ctx, false)
	if err != nil {
		t.Fatalf("got error writing tokens to storage: %v", err)
	}
	if a.fragment != nil {
		t.Errorf("fragment was not reset after write to storage")
	}

	protoSegment = readSegmentFromStorage(t, core, path)
	out = &activity.TokenCount{}
	err = proto.Unmarshal(protoSegment.Value, out)
	if err != nil {
		t.Fatalf("could not unmarshal protobuf: %v", err)
	}

	if len(out.CountByNamespaceID) != 3 {
		t.Fatalf("unexpected token length. Expected %d, got %d", 3, len(out.CountByNamespaceID))
	}
	for i := 0; i < 3; i++ {
		if _, ok := out.CountByNamespaceID[nsIDs[i]]; !ok {
			t.Fatalf("namespace ID %s missing from token counts", nsIDs[i])
		}
	}
	if out.CountByNamespaceID[nsIDs[0]] != 4 {
		t.Errorf("namespace ID %s has %d count, expected %d", nsIDs[0], out.CountByNamespaceID[nsIDs[0]], 4)
	}
	if out.CountByNamespaceID[nsIDs[1]] != 1 {
		t.Errorf("namespace ID %s has %d count, expected %d", nsIDs[1], out.CountByNamespaceID[nsIDs[1]], 1)
	}
	if out.CountByNamespaceID[nsIDs[2]] != 1 {
		t.Errorf("namespace ID %s has %d count, expected %d", nsIDs[2], out.CountByNamespaceID[nsIDs[2]], 1)
	}
}

// TestActivityLog_SaveTokensToStorageDoesNotUpdateTokenCount ensures that
// a new fragment with nonEntityTokens will not update the currentSegment's
// tokenCount, as this field will not be used going forward.
func TestActivityLog_SaveTokensToStorageDoesNotUpdateTokenCount(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ctx := context.Background()

	a := core.activityLog
	a.SetStandbyEnable(ctx, true)
	a.SetStartTimestamp(time.Now().Unix()) // set a nonzero segment

	tokenPath := fmt.Sprintf("%sdirecttokens/%d/0", ActivityLogPrefix, a.GetStartTimestamp())
	clientPath := fmt.Sprintf("sys/counters/activity/log/entity/%d/0", a.GetStartTimestamp())
	// Create some entries without entityIDs
	tokenEntryOne := logical.TokenEntry{NamespaceID: namespace.RootNamespaceID, Policies: []string{"hi"}}
	entityEntry := logical.TokenEntry{EntityID: "foo", NamespaceID: namespace.RootNamespaceID, Policies: []string{"hi"}}

	idNonEntity, isTWE := tokenEntryOne.CreateClientID()

	for i := 0; i < 3; i++ {
		err := a.HandleTokenUsage(ctx, &tokenEntryOne, idNonEntity, isTWE)
		if err != nil {
			t.Fatal(err)
		}
	}

	idEntity, isTWE := entityEntry.CreateClientID()
	for i := 0; i < 2; i++ {
		err := a.HandleTokenUsage(ctx, &entityEntry, idEntity, isTWE)
		if err != nil {
			t.Fatal(err)
		}
	}
	err := a.saveCurrentSegmentToStorage(ctx, false)
	if err != nil {
		t.Fatalf("got error writing TWEs to storage: %v", err)
	}

	// Assert that new elements have been written to the fragment
	if a.fragment != nil {
		t.Errorf("fragment was not reset after write to storage")
	}

	// Assert that no tokens have been written to the fragment
	readSegmentFromStorageNil(t, core, tokenPath)

	e := readSegmentFromStorage(t, core, clientPath)
	out := &activity.EntityActivityLog{}
	err = proto.Unmarshal(e.Value, out)
	if err != nil {
		t.Fatalf("could not unmarshal protobuf: %v", err)
	}
	if len(out.Clients) != 2 {
		t.Fatalf("added 3 distinct TWEs and 2 distinct entity tokens that should all result in the same ID, got: %d", len(out.Clients))
	}
	nonEntityTokenFlag := false
	entityTokenFlag := false
	for _, client := range out.Clients {
		if client.NonEntity == true {
			nonEntityTokenFlag = true
			if client.ClientID != idNonEntity {
				t.Fatalf("expected a client ID of %s, but saved instead %s", idNonEntity, client.ClientID)
			}
		}
		if client.NonEntity == false {
			entityTokenFlag = true
			if client.ClientID != idEntity {
				t.Fatalf("expected a client ID of %s, but saved instead %s", idEntity, client.ClientID)
			}
		}
	}

	if !nonEntityTokenFlag || !entityTokenFlag {
		t.Fatalf("Saved clients missing TWE: %v; saved clients missing entity token: %v", nonEntityTokenFlag, entityTokenFlag)
	}
}

// TestActivityLog_SaveEntitiesToStorage calls AddEntityToFragment with clients with different namespaces and then
// writes the segment to storage. Read back from storage, and verify that client IDs exist in storage.
func TestActivityLog_SaveEntitiesToStorage(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ctx := context.Background()

	a := core.activityLog
	a.SetStandbyEnable(ctx, true)
	a.SetStartTimestamp(time.Now().Unix()) // set a nonzero segment

	now := time.Now()
	ids := []string{"11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222", "33333333-2222-2222-2222-222222222222"}
	times := [...]int64{
		now.Unix(),
		now.Add(1 * time.Second).Unix(),
		now.Add(2 * time.Second).Unix(),
	}
	path := fmt.Sprintf("%sentity/%d/0", ActivityLogPrefix, a.GetStartTimestamp())

	a.AddEntityToFragment(ids[0], "root", times[0])
	a.AddEntityToFragment(ids[1], "root2", times[1])
	err := a.saveCurrentSegmentToStorage(ctx, false)
	if err != nil {
		t.Fatalf("got error writing entities to storage: %v", err)
	}
	if a.fragment != nil {
		t.Errorf("fragment was not reset after write to storage")
	}

	protoSegment := readSegmentFromStorage(t, core, path)
	out := &activity.EntityActivityLog{}
	err = proto.Unmarshal(protoSegment.Value, out)
	if err != nil {
		t.Fatalf("could not unmarshal protobuf: %v", err)
	}
	expectedEntityIDs(t, out, ids[:2])

	a.AddEntityToFragment(ids[0], "root", times[2])
	a.AddEntityToFragment(ids[2], "root", times[2])
	err = a.saveCurrentSegmentToStorage(ctx, false)
	if err != nil {
		t.Fatalf("got error writing segments to storage: %v", err)
	}

	protoSegment = readSegmentFromStorage(t, core, path)
	out = &activity.EntityActivityLog{}
	err = proto.Unmarshal(protoSegment.Value, out)
	if err != nil {
		t.Fatalf("could not unmarshal protobuf: %v", err)
	}
	expectedEntityIDs(t, out, ids)
}

// TestActivityLog_StoreAndReadHyperloglog inserts into a hyperloglog, stores it and then reads it back. The test
// verifies the estimate count is correct.
func TestActivityLog_StoreAndReadHyperloglog(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ctx := context.Background()

	a := core.activityLog
	a.SetStandbyEnable(ctx, true)
	a.SetStartTimestamp(time.Now().Unix()) // set a nonzero segment
	currentMonth := timeutil.StartOfMonth(time.Now())
	currentMonthHll := hyperloglog.New()
	currentMonthHll.Insert([]byte("a"))
	currentMonthHll.Insert([]byte("a"))
	currentMonthHll.Insert([]byte("b"))
	currentMonthHll.Insert([]byte("c"))
	currentMonthHll.Insert([]byte("d"))
	currentMonthHll.Insert([]byte("d"))

	err := a.StoreHyperlogLog(ctx, currentMonth, currentMonthHll)
	if err != nil {
		t.Fatalf("error storing hyperloglog in storage: %v", err)
	}
	fetchedHll, err := a.CreateOrFetchHyperlogLog(ctx, currentMonth)
	// check the distinct count stored from hll
	if fetchedHll.Estimate() != 4 {
		t.Fatalf("wrong number of distinct elements: expected: 5 actual: %v", fetchedHll.Estimate())
	}
}

// TestModifyResponseMonthsNilAppend calls modifyResponseMonths for a range of 5 months ago to now. It verifies that the
// 5 months in the range are correct.
func TestModifyResponseMonthsNilAppend(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	end := time.Now().UTC()
	start := timeutil.StartOfMonth(end).AddDate(0, -5, 0)
	responseMonthTimestamp := timeutil.StartOfMonth(end).AddDate(0, -3, 0).Format(time.RFC3339)
	responseMonths := []*ResponseMonth{{Timestamp: responseMonthTimestamp}}
	months := a.modifyResponseMonths(responseMonths, start, end)
	if len(months) != 5 {
		t.Fatal("wrong number of months padded")
	}
	for _, m := range months {
		ts, err := time.Parse(time.RFC3339, m.Timestamp)
		if err != nil {
			t.Fatal(err)
		}
		if !ts.Equal(start) {
			t.Fatalf("incorrect time in month sequence timestamps: expected %+v, got %+v", start, ts)
		}
		start = timeutil.StartOfMonth(start).AddDate(0, 1, 0)
	}
	// The following is a redundant check, but for posterity and readability I've
	// made it explicit.
	lastMonth, err := time.Parse(time.RFC3339, months[4].Timestamp)
	if err != nil {
		t.Fatal(err)
	}
	if timeutil.IsCurrentMonth(lastMonth, time.Now().UTC()) {
		t.Fatalf("do not include current month timestamp in nil padding for months")
	}
}

// TestActivityLog_ReceivedFragment calls receivedFragment with a fragment and verifies it gets added to
// standbyFragmentsReceived. Send the same fragment again and then verify that it doesn't change the entity map but does
// get added to standbyFragmentsReceived.
func TestActivityLog_ReceivedFragment(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	a.SetEnable(true)

	ids := []string{
		"11111111-1111-1111-1111-111111111111",
		"22222222-2222-2222-2222-222222222222",
	}

	entityRecords := []*activity.EntityRecord{
		{
			ClientID:    ids[0],
			NamespaceID: "root",
			Timestamp:   time.Now().Unix(),
		},
		{
			ClientID:    ids[1],
			NamespaceID: "root",
			Timestamp:   time.Now().Unix(),
		},
	}

	fragment := &activity.LogFragment{
		OriginatingNode: "test-123",
		Clients:         entityRecords,
		NonEntityTokens: make(map[string]uint64),
	}

	if len(a.standbyFragmentsReceived) != 0 {
		t.Fatalf("fragment already received")
	}

	a.receivedFragment(fragment)

	checkExpectedEntitiesInMap(t, a, ids)

	if len(a.standbyFragmentsReceived) != 1 {
		t.Fatalf("fragment count is %v, expected 1", len(a.standbyFragmentsReceived))
	}

	// Send a duplicate, should be stored but not change entity map
	a.receivedFragment(fragment)

	checkExpectedEntitiesInMap(t, a, ids)

	if len(a.standbyFragmentsReceived) != 2 {
		t.Fatalf("fragment count is %v, expected 2", len(a.standbyFragmentsReceived))
	}
}

// TestActivityLog_availableLogsEmptyDirectory verifies that availableLogs returns an empty slice when the log directory
// is empty.
func TestActivityLog_availableLogsEmptyDirectory(t *testing.T) {
	// verify that directory is empty, and nothing goes wrong
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	times, err := a.availableLogs(context.Background(), time.Now())
	if err != nil {
		t.Fatalf("error getting start_time(s) for empty activity log")
	}
	if len(times) != 0 {
		t.Fatalf("invalid number of start_times returned. expected 0, got %d", len(times))
	}
}

// TestActivityLog_availableLogs writes to the direct token paths and entity paths and verifies that the correct start
// times are returned.
func TestActivityLog_availableLogs(t *testing.T) {
	// set up a few files in storage
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	paths := [...]string{"entity/1111/1", "directtokens/1111/1", "directtokens/1000000/1", "entity/992/3", "directtokens/992/1"}
	expectedTimes := [...]time.Time{time.Unix(1000000, 0), time.Unix(1111, 0), time.Unix(992, 0)}

	for _, path := range paths {
		WriteToStorage(t, core, ActivityLogPrefix+path, []byte("test"))
	}

	// verify above files are there, and dates in correct order
	times, err := a.availableLogs(context.Background(), time.Now())
	if err != nil {
		t.Fatalf("error getting start_time(s) for activity log")
	}

	if len(times) != len(expectedTimes) {
		t.Fatalf("invalid number of start_times returned. expected %d, got %d", len(expectedTimes), len(times))
	}
	for i := range times {
		if !times[i].Equal(expectedTimes[i]) {
			t.Errorf("invalid time. expected %v, got %v", expectedTimes[i], times[i])
		}
	}
}

// TestActivityLog_createRegenerationIntentLog tests that we can correctly create a regeneration intent log given the segments in storage
func TestActivityLog_createRegenerationIntentLog(t *testing.T) {
	testCases := []struct {
		name          string
		times         []time.Time
		expectedLog   *ActivityIntentLog
		expectedError bool
	}{
		{
			"no segments",
			[]time.Time{},
			nil,
			true,
		},
		{
			"one segment",
			[]time.Time{
				time.Date(2024, 4, 4, 10, 54, 12, 0, time.UTC),
			},
			nil,
			true,
		},
		{
			"most recent segment is 3 months ago",
			[]time.Time{
				time.Date(2024, 1, 4, 10, 54, 12, 0, time.UTC),
				time.Date(2024, 1, 3, 10, 54, 12, 0, time.UTC),
			},
			&ActivityIntentLog{NextMonth: 0, PreviousMonth: 1704365652},
			false,
		},
		{
			"lots of segments",
			[]time.Time{
				// two this month
				time.Date(2024, 4, 4, 10, 54, 12, 0, time.UTC),
				time.Date(2024, 4, 6, 10, 54, 12, 0, time.UTC),
				// three last month
				time.Date(2024, 3, 3, 10, 54, 12, 0, time.UTC),
				time.Date(2024, 3, 6, 10, 54, 12, 0, time.UTC),
				time.Date(2024, 3, 14, 10, 54, 12, 0, time.UTC),
				// two the month before that
				time.Date(2024, 2, 10, 10, 54, 12, 0, time.UTC),
				time.Date(2024, 2, 17, 10, 54, 12, 0, time.UTC),
			},
			&ActivityIntentLog{NextMonth: 1712228052, PreviousMonth: 1710413652},
			false,
		},
	}

	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	now := time.Date(2024, 4, 10, 10, 54, 12, 0, time.UTC)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deletePaths := make([]string, 0)

			// insert the times we're given
			paths := make([]string, 0, len(tc.times))
			for _, tm := range tc.times {
				paths = append(paths, fmt.Sprintf("entity/%d/1", tm.Unix()))
			}

			for _, subPath := range paths {
				fullPath := ActivityLogPrefix + subPath
				WriteToStorage(t, core, fullPath, []byte("test"))
				deletePaths = append(deletePaths, fullPath)
			}

			// regenerate the log
			intentLog, err := a.createRegenerationIntentLog(context.Background(), now)
			if tc.expectedError && err == nil {
				t.Fatal("expected an error and got none")
			}
			if !tc.expectedError && err != nil {
				t.Fatal(err)
			}

			// verify it's what we expect
			if diff := deep.Equal(intentLog, tc.expectedLog); len(diff) != 0 {
				t.Errorf("got=%v, expected=%v, diff=%v", intentLog, tc.expectedLog, diff)
			}

			// delete everything we wrote so the next test starts fresh
			for _, p := range deletePaths {
				err := core.barrier.Delete(ctx, p)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestActivityLog_MultipleFragmentsAndSegments adds 4000 clients to a fragment
// and saves it and reads it. The test then adds 4000 more clients and calls
// receivedFragment with 200 more entities. The current segment is saved to
// storage and read back. The test verifies that there are ActivitySegmentClientCapacity clients in the
// first and second segment index, then the rest in the third index.
func TestActivityLog_MultipleFragmentsAndSegments(t *testing.T) {
	core, _, _ := TestCoreUnsealedWithConfig(t, &CoreConfig{
		ActivityLogConfig: ActivityLogCoreConfig{
			DisableFragmentWorker: true,
			DisableTimers:         true,
		},
	})
	a := core.activityLog

	// enabled check is now inside AddClientToFragment
	a.SetEnable(true)
	a.SetStartTimestamp(time.Now().Unix()) // set a nonzero segment

	startTimestamp := a.GetStartTimestamp()
	path0 := fmt.Sprintf("sys/counters/activity/log/entity/%d/0", startTimestamp)
	path1 := fmt.Sprintf("sys/counters/activity/log/entity/%d/1", startTimestamp)
	path2 := fmt.Sprintf("sys/counters/activity/log/entity/%d/2", startTimestamp)
	tokenPath := fmt.Sprintf("sys/counters/activity/log/directtokens/%d/0", startTimestamp)

	genID := func(i int) string {
		return fmt.Sprintf("11111111-1111-1111-1111-%012d", i)
	}
	ts := time.Now().Unix()

	// First ActivitySegmentClientCapacity should fit in one segment
	for i := 0; i < 4000; i++ {
		a.AddEntityToFragment(genID(i), "root", ts)
	}

	// Consume new fragment notification.
	// The worker may have gotten it first, before processing
	// the close!
	select {
	case <-a.newFragmentCh:
	default:
	}

	// Save segment
	err := a.saveCurrentSegmentToStorage(context.Background(), false)
	if err != nil {
		t.Fatalf("got error writing entities to storage: %v", err)
	}

	protoSegment0 := readSegmentFromStorage(t, core, path0)
	entityLog0 := activity.EntityActivityLog{}
	err = proto.Unmarshal(protoSegment0.Value, &entityLog0)
	if err != nil {
		t.Fatalf("could not unmarshal protobuf: %v", err)
	}
	if len(entityLog0.Clients) != ActivitySegmentClientCapacity {
		t.Fatalf("unexpected entity length. Expected %d, got %d", ActivitySegmentClientCapacity, len(entityLog0.Clients))
	}

	// 4000 more local entities
	for i := 4000; i < 8000; i++ {
		a.AddEntityToFragment(genID(i), "root", ts)
	}

	// Simulated remote fragment with 100 duplicate entities
	tokens1 := map[string]uint64{
		"root":  3,
		"aaaaa": 4,
		"bbbbb": 5,
	}
	fragment1 := &activity.LogFragment{
		OriginatingNode: "test-123",
		Clients:         make([]*activity.EntityRecord, 0, 100),
		NonEntityTokens: tokens1,
	}
	for i := 4000; i < 4100; i++ {
		fragment1.Clients = append(fragment1.Clients, &activity.EntityRecord{
			ClientID:    genID(i),
			NamespaceID: "root",
			Timestamp:   ts,
		})
	}

	// Simulated remote fragment with 100 new entities
	tokens2 := map[string]uint64{
		"root":  6,
		"aaaaa": 7,
		"bbbbb": 8,
	}
	fragment2 := &activity.LogFragment{
		OriginatingNode: "test-123",
		Clients:         make([]*activity.EntityRecord, 0, 100),
		NonEntityTokens: tokens2,
	}
	for i := 8000; i < 8100; i++ {
		fragment2.Clients = append(fragment2.Clients, &activity.EntityRecord{
			ClientID:    genID(i),
			NamespaceID: "root",
			Timestamp:   ts,
		})
	}
	a.receivedFragment(fragment1)
	a.receivedFragment(fragment2)

	select {
	case <-a.newFragmentCh:
	case <-time.After(time.Minute):
		t.Fatal("timed out waiting for new fragment")
	}

	err = a.saveCurrentSegmentToStorage(context.Background(), false)
	if err != nil {
		t.Fatalf("got error writing entities to storage: %v", err)
	}

	seqNum := a.GetEntitySequenceNumber()
	if seqNum != 2 {
		t.Fatalf("expected sequence number 2, got %v", seqNum)
	}

	protoSegment0 = readSegmentFromStorage(t, core, path0)
	err = proto.Unmarshal(protoSegment0.Value, &entityLog0)
	if err != nil {
		t.Fatalf("could not unmarshal protobuf: %v", err)
	}
	if len(entityLog0.Clients) != ActivitySegmentClientCapacity {
		t.Fatalf("unexpected client length. Expected %d, got %d", ActivitySegmentClientCapacity,
			len(entityLog0.Clients))
	}

	protoSegment1 := readSegmentFromStorage(t, core, path1)
	entityLog1 := activity.EntityActivityLog{}
	err = proto.Unmarshal(protoSegment1.Value, &entityLog1)
	if err != nil {
		t.Fatalf("could not unmarshal protobuf: %v", err)
	}
	if len(entityLog1.Clients) != ActivitySegmentClientCapacity {
		t.Fatalf("unexpected entity length. Expected %d, got %d", ActivitySegmentClientCapacity,
			len(entityLog1.Clients))
	}

	protoSegment2 := readSegmentFromStorage(t, core, path2)
	entityLog2 := activity.EntityActivityLog{}
	err = proto.Unmarshal(protoSegment2.Value, &entityLog2)
	if err != nil {
		t.Fatalf("could not unmarshal protobuf: %v", err)
	}
	expectedCount := 8100 - (ActivitySegmentClientCapacity * 2)
	if len(entityLog2.Clients) != expectedCount {
		t.Fatalf("unexpected entity length. Expected %d, got %d", expectedCount,
			len(entityLog1.Clients))
	}

	entityPresent := make(map[string]struct{})
	for _, e := range entityLog0.Clients {
		entityPresent[e.ClientID] = struct{}{}
	}
	for _, e := range entityLog1.Clients {
		entityPresent[e.ClientID] = struct{}{}
	}
	for _, e := range entityLog2.Clients {
		entityPresent[e.ClientID] = struct{}{}
	}
	for i := 0; i < 8100; i++ {
		expectedID := genID(i)
		if _, present := entityPresent[expectedID]; !present {
			t.Fatalf("entity ID %v = %v not present", i, expectedID)
		}
	}
	expectedNSCounts := map[string]uint64{
		"root":  9,
		"aaaaa": 11,
		"bbbbb": 13,
	}
	tokenSegment := readSegmentFromStorage(t, core, tokenPath)
	tokenCount := activity.TokenCount{}
	err = proto.Unmarshal(tokenSegment.Value, &tokenCount)
	if err != nil {
		t.Fatalf("could not unmarshal protobuf: %v", err)
	}

	if !reflect.DeepEqual(expectedNSCounts, tokenCount.CountByNamespaceID) {
		t.Fatalf("token counts are not equal, expected %v got %v", expectedNSCounts, tokenCount.CountByNamespaceID)
	}
}

// TestActivityLog_API_ConfigCRUD_Census performs various CRUD operations on internal/counters/config
// depending on license reporting
func TestActivityLog_API_ConfigCRUD_Census(t *testing.T) {
	core, b, _ := testCoreSystemBackend(t)
	view := core.systemBarrierView

	req := logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
	req.Storage = view
	req.Data["retention_months"] = 2
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if core.ManualLicenseReportingEnabled() {
		if err == nil {
			t.Fatal("expected error")
		}
		if resp.Data["error"] != `retention_months must be at least 48 while Reporting is enabled` {
			t.Fatalf("bad: %v", resp)
		}
	} else {
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
	req.Storage = view
	req.Data["retention_months"] = 56
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
	req.Storage = view
	req.Data["enabled"] = "disable"
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if core.ManualLicenseReportingEnabled() {
		if err == nil {
			t.Fatal("expected error")
		}
		if resp.Data["error"] != `cannot disable the activity log while Reporting is enabled` {
			t.Fatalf("bad: %v", resp)
		}
	} else {
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if resp != nil {
			t.Fatalf("bad: %#v", resp)
		}
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
	req.Storage = view
	req.Data["enabled"] = "enable"
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "internal/counters/config")
	req.Storage = view
	resp, err = b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expected := map[string]interface{}{
		"retention_months":         56,
		"enabled":                  "enable",
		"queries_available":        false,
		"reporting_enabled":        core.AutomatedLicenseReportingEnabled(),
		"billing_start_timestamp":  core.BillingStart(),
		"minimum_retention_months": core.activityLog.configOverrides.MinimumRetentionMonths,
	}

	if diff := deep.Equal(resp.Data, expected); len(diff) > 0 {
		t.Fatalf("diff: %v", diff)
	}
}

// TestActivityLog_parseSegmentNumberFromPath verifies that the segment number is extracted correctly from a path.
func TestActivityLog_parseSegmentNumberFromPath(t *testing.T) {
	testCases := []struct {
		input        string
		expected     int
		expectExists bool
	}{
		{
			input:        "path/to/log/5",
			expected:     5,
			expectExists: true,
		},
		{
			input:        "/path/to/log/5",
			expected:     5,
			expectExists: true,
		},
		{
			input:        "path/to/log/",
			expected:     0,
			expectExists: false,
		},
		{
			input:        "path/to/log/foo",
			expected:     0,
			expectExists: false,
		},
		{
			input:        "",
			expected:     0,
			expectExists: false,
		},
		{
			input:        "5",
			expected:     5,
			expectExists: true,
		},
	}

	for _, tc := range testCases {
		result, ok := parseSegmentNumberFromPath(tc.input)
		if result != tc.expected {
			t.Errorf("expected: %d, got: %d for input %q", tc.expected, result, tc.input)
		}
		if ok != tc.expectExists {
			t.Errorf("unexpected value presence. expected exists: %t, got: %t for input %q", tc.expectExists, ok, tc.input)
		}
	}
}

// TestActivityLog_getLastEntitySegmentNumber verifies that the last segment number is correctly returned.
func TestActivityLog_getLastEntitySegmentNumber(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	paths := [...]string{"entity/992/0", "entity/1000/-1", "entity/1001/foo", "entity/1111/0", "entity/1111/1"}
	for _, path := range paths {
		WriteToStorage(t, core, ActivityLogPrefix+path, []byte("test"))
	}

	testCases := []struct {
		input        int64
		expectedVal  uint64
		expectExists bool
	}{
		{
			input:        992,
			expectedVal:  0,
			expectExists: true,
		},
		{
			input:        1000,
			expectedVal:  0,
			expectExists: false,
		},
		{
			input:        1001,
			expectedVal:  0,
			expectExists: false,
		},
		{
			input:        1111,
			expectedVal:  1,
			expectExists: true,
		},
		{
			input:        2222,
			expectedVal:  0,
			expectExists: false,
		},
	}

	ctx := context.Background()
	for _, tc := range testCases {
		result, exists, err := a.getLastEntitySegmentNumber(ctx, time.Unix(tc.input, 0))
		if err != nil {
			t.Fatalf("unexpected error for input %d: %v", tc.input, err)
		}
		if exists != tc.expectExists {
			t.Errorf("expected result exists: %t, got: %t for input: %d", tc.expectExists, exists, tc.input)
		}
		if result != tc.expectedVal {
			t.Errorf("expected: %d got: %d for input: %d", tc.expectedVal, result, tc.input)
		}
	}
}

// TestActivityLog_tokenCountExists writes to the direct tokens segment path and verifies that segment count exists
// returns true for the segments at these paths.
func TestActivityLog_tokenCountExists(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	paths := [...]string{"directtokens/992/0", "directtokens/1001/foo", "directtokens/1111/0", "directtokens/2222/1"}
	for _, path := range paths {
		WriteToStorage(t, core, ActivityLogPrefix+path, []byte("test"))
	}

	testCases := []struct {
		input        int64
		expectExists bool
	}{
		{
			input:        992,
			expectExists: true,
		},
		{
			input:        1001,
			expectExists: false,
		},
		{
			input:        1111,
			expectExists: true,
		},
		{
			input:        2222,
			expectExists: false,
		},
	}

	ctx := context.Background()
	for _, tc := range testCases {
		exists, err := a.tokenCountExists(ctx, time.Unix(tc.input, 0))
		if err != nil {
			t.Fatalf("unexpected error for input %d: %v", tc.input, err)
		}
		if exists != tc.expectExists {
			t.Errorf("expected segment to exist: %t but got: %t for input: %d", tc.expectExists, exists, tc.input)
		}
	}
}

// entityRecordsEqual compares the parts we care about from two activity entity record slices
// note: this makes a copy of the []*activity.EntityRecord so that misordered slices won't fail the comparison,
// but the function won't modify the order of the slices to compare
func entityRecordsEqual(t *testing.T, record1, record2 []*activity.EntityRecord) bool {
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

func (a *ActivityLog) resetEntitiesInMemory(t *testing.T) {
	t.Helper()

	a.l.Lock()
	defer a.l.Unlock()
	a.fragmentLock.Lock()
	defer a.fragmentLock.Unlock()
	a.currentSegment = segmentInfo{
		startTimestamp: time.Time{}.Unix(),
		currentClients: &activity.EntityActivityLog{
			Clients: make([]*activity.EntityRecord, 0),
		},
		tokenCount:           a.currentSegment.tokenCount,
		clientSequenceNumber: 0,
	}

	a.partialMonthClientTracker = make(map[string]*activity.EntityRecord)
}

// TestActivityLog_loadCurrentClientSegment writes entity segments and calls loadCurrentClientSegment, then verifies
// that the correct values are returned when querying the current segment.
func TestActivityLog_loadCurrentClientSegment(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	// we must verify that loadCurrentClientSegment doesn't overwrite the in-memory token count
	tokenRecords := make(map[string]uint64)
	tokenRecords["test"] = 1
	tokenCount := &activity.TokenCount{
		CountByNamespaceID: tokenRecords,
	}
	a.l.Lock()
	a.currentSegment.tokenCount = tokenCount
	a.l.Unlock()

	// setup in-storage data to load for testing
	entityRecords := []*activity.EntityRecord{
		{
			ClientID:    "11111111-1111-1111-1111-111111111111",
			NamespaceID: "root",
			Timestamp:   time.Now().Unix(),
		},
		{
			ClientID:    "22222222-2222-2222-2222-222222222222",
			NamespaceID: "root",
			Timestamp:   time.Now().Unix(),
		},
	}
	testEntities1 := &activity.EntityActivityLog{
		Clients: entityRecords[:1],
	}
	testEntities2 := &activity.EntityActivityLog{
		Clients: entityRecords[1:2],
	}
	testEntities3 := &activity.EntityActivityLog{
		Clients: entityRecords[:2],
	}

	time1 := time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC).Unix()
	time2 := time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC).Unix()
	testCases := []struct {
		time     int64
		seqNum   uint64
		path     string
		entities *activity.EntityActivityLog
	}{
		{
			time:     time1,
			seqNum:   0,
			path:     "entity/" + fmt.Sprint(time1) + "/0",
			entities: testEntities1,
		},
		{
			// we want to verify that data from segment 0 hasn't been loaded
			time:     time1,
			seqNum:   1,
			path:     "entity/" + fmt.Sprint(time1) + "/1",
			entities: testEntities2,
		},
		{
			time:     time2,
			seqNum:   0,
			path:     "entity/" + fmt.Sprint(time2) + "/0",
			entities: testEntities3,
		},
	}

	for _, tc := range testCases {
		data, err := proto.Marshal(tc.entities)
		if err != nil {
			t.Fatal(err.Error())
		}
		WriteToStorage(t, core, ActivityLogPrefix+tc.path, data)
	}

	ctx := context.Background()
	for _, tc := range testCases {
		a.l.Lock()
		a.fragmentLock.Lock()
		// loadCurrentClientSegment requires us to grab the fragment lock and the
		// activityLog lock, as per the comment in the loadCurrentClientSegment
		// function
		err := a.loadCurrentClientSegment(ctx, time.Unix(tc.time, 0), tc.seqNum)
		a.fragmentLock.Unlock()
		a.l.Unlock()

		if err != nil {
			t.Fatalf("got error loading data for %q: %v", tc.path, err)
		}
		if !reflect.DeepEqual(a.GetStoredTokenCountByNamespaceID(), tokenCount.CountByNamespaceID) {
			t.Errorf("this function should not wipe out the in-memory token count")
		}

		// verify accurate data in in-memory current segment
		startTimestamp := a.GetStartTimestamp()
		if startTimestamp != tc.time {
			t.Errorf("bad timestamp loaded. expected: %v, got: %v for path %q", tc.time, startTimestamp, tc.path)
		}

		seqNum := a.GetEntitySequenceNumber()
		if seqNum != tc.seqNum {
			t.Errorf("bad sequence number loaded. expected: %v, got: %v for path %q", tc.seqNum, seqNum, tc.path)
		}

		currentEntities := a.GetCurrentEntities()
		if !entityRecordsEqual(t, currentEntities.Clients, tc.entities.Clients) {
			t.Errorf("bad data loaded. expected: %v, got: %v for path %q", tc.entities.Clients, currentEntities, tc.path)
		}

		activeClients := core.GetActiveClientsList()
		if err := ActiveEntitiesEqual(activeClients, tc.entities.Clients); err != nil {
			t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v for path %q: %v", tc.entities.Clients, activeClients, tc.path, err)
		}

		a.resetEntitiesInMemory(t)
	}
}

// TestActivityLog_loadPriorEntitySegment writes entities to two months and calls loadPriorEntitySegment for each month,
// verifying that the active clients are correct.
func TestActivityLog_loadPriorEntitySegment(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	a.SetEnable(true)

	// setup in-storage data to load for testing
	entityRecords := []*activity.EntityRecord{
		{
			ClientID:    "11111111-1111-1111-1111-111111111111",
			NamespaceID: "root",
			Timestamp:   time.Now().Unix(),
		},
		{
			ClientID:    "22222222-2222-2222-2222-222222222222",
			NamespaceID: "root",
			Timestamp:   time.Now().Unix(),
		},
	}
	testEntities1 := &activity.EntityActivityLog{
		Clients: entityRecords[:1],
	}
	testEntities2 := &activity.EntityActivityLog{
		Clients: entityRecords[:2],
	}

	time1 := time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC).Unix()
	time2 := time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC).Unix()
	testCases := []struct {
		time     int64
		seqNum   uint64
		path     string
		entities *activity.EntityActivityLog
		// set true if the in-memory active entities should be wiped because the next test case is a new month
		// this also means that currentSegment.startTimestamp must be updated with :time:
		refresh bool
	}{
		{
			time:     time1,
			seqNum:   0,
			path:     "entity/" + fmt.Sprint(time1) + "/0",
			entities: testEntities1,
			refresh:  true,
		},
		{
			// verify that we don't have a duplicate (shouldn't be possible with the current implementation)
			time:     time1,
			seqNum:   1,
			path:     "entity/" + fmt.Sprint(time1) + "/1",
			entities: testEntities2,
			refresh:  true,
		},
		{
			time:     time2,
			seqNum:   0,
			path:     "entity/" + fmt.Sprint(time2) + "/0",
			entities: testEntities2,
			refresh:  true,
		},
	}

	for _, tc := range testCases {
		data, err := proto.Marshal(tc.entities)
		if err != nil {
			t.Fatal(err.Error())
		}
		WriteToStorage(t, core, ActivityLogPrefix+tc.path, data)
	}

	ctx := context.Background()
	for _, tc := range testCases {
		if tc.refresh {
			a.l.Lock()
			a.fragmentLock.Lock()
			a.partialMonthClientTracker = make(map[string]*activity.EntityRecord)
			a.currentSegment.startTimestamp = tc.time
			a.fragmentLock.Unlock()
			a.l.Unlock()
		}

		err := a.loadPriorEntitySegment(ctx, time.Unix(tc.time, 0), tc.seqNum)
		if err != nil {
			t.Fatalf("got error loading data for %q: %v", tc.path, err)
		}

		activeClients := core.GetActiveClientsList()
		if err := ActiveEntitiesEqual(activeClients, tc.entities.Clients); err != nil {
			t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v for path %q: %v", tc.entities.Clients, activeClients, tc.path, err)
		}
	}
}

// TestActivityLog_loadTokenCount ensures that previous segments with tokenCounts
// can still be read from storage, even when TWE's have distinct, tracked clientIDs.
func TestActivityLog_loadTokenCount(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog

	// setup in-storage data to load for testing
	tokenRecords := make(map[string]uint64)
	for i := 1; i < 4; i++ {
		nsID := "ns" + strconv.Itoa(i)
		tokenRecords[nsID] = uint64(i)
	}
	tokenCount := &activity.TokenCount{
		CountByNamespaceID: tokenRecords,
	}

	data, err := proto.Marshal(tokenCount)
	if err != nil {
		t.Fatal(err.Error())
	}

	testCases := []struct {
		time int64
		path string
	}{
		{
			time: 1111,
			path: "directtokens/1111/0",
		},
		{
			time: 2222,
			path: "directtokens/2222/0",
		},
	}

	ctx := context.Background()
	for _, tc := range testCases {
		WriteToStorage(t, core, ActivityLogPrefix+tc.path, data)
	}

	for _, tc := range testCases {
		err := a.loadTokenCount(ctx, time.Unix(tc.time, 0))
		if err != nil {
			t.Fatalf("got error loading data for %q: %v", tc.path, err)
		}

		nsCount := a.GetStoredTokenCountByNamespaceID()
		if !reflect.DeepEqual(nsCount, tokenRecords) {
			t.Errorf("bad token count loaded. expected: %v got: %v for path %q", tokenRecords, nsCount, tc.path)
		}
	}
}

// TestActivityLog_StopAndRestart disables the activity log, waits for deletes to complete, and then enables the
// activity log. The activity log is then stopped and started again, to simulate a seal and unseal. The test then
// verifies that there's no error adding an entity, direct token, and when writing a segment to storage.
func TestActivityLog_StopAndRestart(t *testing.T) {
	core, b, _ := testCoreSystemBackend(t)
	sysView := core.systemBarrierView

	a := core.activityLog
	ctx := namespace.RootContext(nil)

	// Disable, then enable, to exercise newly-enabled code
	a.SetConfig(ctx, activityConfig{
		Enabled:             "disable",
		RetentionMonths:     12,
		DefaultReportMonths: 12,
	})

	// On enterprise, a segment will be created, and
	// disabling it will trigger deletion, so wait
	// for that deletion to finish.
	// (Alternatively, we could ensure that the next segment
	// uses a different timestamp by waiting 1 second.)
	a.WaitForDeletion()

	// Go through request to ensure config is persisted
	req := logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
	req.Storage = sysView
	req.Data["enabled"] = "enable"
	resp, err := b.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Simulate seal/unseal cycle
	core.stopActivityLog()
	var wg sync.WaitGroup
	core.setupActivityLog(ctx, &wg, false)
	wg.Wait()

	a = core.activityLog
	if a.GetStoredTokenCountByNamespaceID() == nil {
		t.Fatalf("nil token count map")
	}

	a.AddEntityToFragment("1111-1111", "root", time.Now().Unix())
	a.AddTokenToFragment("root")

	err = a.saveCurrentSegmentToStorage(ctx, false)
	if err != nil {
		t.Fatal(err)
	}
}

// :base: is the timestamp to start from for the setup logic (use to simulate newest log from past or future)
// entity records returned include [0] data from a previous month and [1:] data from the current month
// token counts returned are from the current month
func setupActivityRecordsInStorage(t *testing.T, base time.Time, includeEntities, includeTokens bool) (*ActivityLog, []*activity.EntityRecord, map[string]uint64) {
	t.Helper()

	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	monthsAgo := base.AddDate(0, -3, 0)

	var entityRecords []*activity.EntityRecord
	if includeEntities {
		entityRecords = []*activity.EntityRecord{
			{
				ClientID:    "11111111-1111-1111-1111-111111111111",
				NamespaceID: namespace.RootNamespaceID,
				Timestamp:   time.Now().Unix(),
			},
			{
				ClientID:    "22222222-2222-2222-2222-222222222222",
				NamespaceID: namespace.RootNamespaceID,
				Timestamp:   time.Now().Unix(),
			},
			{
				ClientID:    "33333333-2222-2222-2222-222222222222",
				NamespaceID: namespace.RootNamespaceID,
				Timestamp:   time.Now().Unix(),
			},
		}
		if constants.IsEnterprise {
			entityRecords = append(entityRecords, []*activity.EntityRecord{
				{
					ClientID:    "44444444-1111-1111-1111-111111111111",
					NamespaceID: "ns1",
					Timestamp:   time.Now().Unix(),
				},
			}...)
		}
		for i, entityRecord := range entityRecords {
			entityData, err := proto.Marshal(&activity.EntityActivityLog{
				Clients: []*activity.EntityRecord{entityRecord},
			})
			if err != nil {
				t.Fatal(err.Error())
			}
			if i == 0 {
				WriteToStorage(t, core, ActivityLogPrefix+"entity/"+fmt.Sprint(monthsAgo.Unix())+"/0", entityData)
			} else {
				WriteToStorage(t, core, ActivityLogPrefix+"entity/"+fmt.Sprint(base.Unix())+"/"+strconv.Itoa(i-1), entityData)
			}
		}
	}

	var tokenRecords map[string]uint64
	if includeTokens {
		tokenRecords = make(map[string]uint64)
		tokenRecords[namespace.RootNamespaceID] = uint64(1)
		if constants.IsEnterprise {
			for i := 1; i < 4; i++ {
				nsID := "ns" + strconv.Itoa(i)
				tokenRecords[nsID] = uint64(i)
			}
		}
		tokenCount := &activity.TokenCount{
			CountByNamespaceID: tokenRecords,
		}

		tokenData, err := proto.Marshal(tokenCount)
		if err != nil {
			t.Fatal(err.Error())
		}

		WriteToStorage(t, core, ActivityLogPrefix+"directtokens/"+fmt.Sprint(base.Unix())+"/0", tokenData)
	}

	return a, entityRecords, tokenRecords
}

// TestActivityLog_refreshFromStoredLog writes records for 3 months ago and this month, then calls refreshFromStoredLog.
// The test verifies that current entities and current tokens are correct.
func TestActivityLog_refreshFromStoredLog(t *testing.T) {
	a, expectedClientRecords, expectedTokenCounts := setupActivityRecordsInStorage(t, time.Now().UTC(), true, true)
	a.SetEnable(true)

	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(context.Background(), &wg, time.Now().UTC())
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	expectedActive := &activity.EntityActivityLog{
		Clients: expectedClientRecords[1:],
	}
	expectedCurrent := &activity.EntityActivityLog{
		Clients: expectedClientRecords[len(expectedClientRecords)-1:],
	}

	currentEntities := a.GetCurrentEntities()
	if !entityRecordsEqual(t, currentEntities.Clients, expectedCurrent.Clients) {
		// we only expect the newest entity segment to be loaded (for the current month)
		t.Errorf("bad activity entity logs loaded. expected: %v got: %v", expectedCurrent, currentEntities)
	}

	nsCount := a.GetStoredTokenCountByNamespaceID()
	if !reflect.DeepEqual(nsCount, expectedTokenCounts) {
		// we expect all token counts to be loaded
		t.Errorf("bad activity token counts loaded. expected: %v got: %v", expectedTokenCounts, nsCount)
	}

	activeClients := a.core.GetActiveClientsList()
	if err := ActiveEntitiesEqual(activeClients, expectedActive.Clients); err != nil {
		// we expect activeClients to be loaded for the entire month
		t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v: %v", expectedActive.Clients, activeClients, err)
	}
}

// TestActivityLog_refreshFromStoredLogWithBackgroundLoadingCancelled writes data from 3 months ago to this month. The
// test closes a.doneCh and calls refreshFromStoredLog, which will not do any processing because the doneCh is closed.
// The test verifies that the current data is not loaded.
func TestActivityLog_refreshFromStoredLogWithBackgroundLoadingCancelled(t *testing.T) {
	a, expectedClientRecords, expectedTokenCounts := setupActivityRecordsInStorage(t, time.Now().UTC(), true, true)
	a.SetEnable(true)

	var wg sync.WaitGroup
	close(a.doneCh)
	defer func() {
		a.l.Lock()
		a.doneCh = make(chan struct{}, 1)
		a.l.Unlock()
	}()

	err := a.refreshFromStoredLog(context.Background(), &wg, time.Now().UTC())
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	expected := &activity.EntityActivityLog{
		Clients: expectedClientRecords[len(expectedClientRecords)-1:],
	}

	currentEntities := a.GetCurrentEntities()
	if !entityRecordsEqual(t, currentEntities.Clients, expected.Clients) {
		// we only expect the newest entity segment to be loaded (for the current month)
		t.Errorf("bad activity entity logs loaded. expected: %v got: %v", expected, currentEntities)
	}

	nsCount := a.GetStoredTokenCountByNamespaceID()
	if !reflect.DeepEqual(nsCount, expectedTokenCounts) {
		// we expect all token counts to be loaded
		t.Errorf("bad activity token counts loaded. expected: %v got: %v", expectedTokenCounts, nsCount)
	}

	activeClients := a.core.GetActiveClientsList()
	if err := ActiveEntitiesEqual(activeClients, expected.Clients); err != nil {
		// we only expect activeClients to be loaded for the newest segment (for the current month)
		t.Error(err)
	}
}

// TestActivityLog_refreshFromStoredLogContextCancelled writes data from 3 months ago to this month and calls
// refreshFromStoredLog with a canceled context, verifying that the function errors because of the canceled context.
func TestActivityLog_refreshFromStoredLogContextCancelled(t *testing.T) {
	a, _, _ := setupActivityRecordsInStorage(t, time.Now().UTC(), true, true)

	var wg sync.WaitGroup
	ctx, cancelFn := context.WithCancel(context.Background())
	cancelFn()

	err := a.refreshFromStoredLog(ctx, &wg, time.Now().UTC())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancelled error, got: %v", err)
	}
}

// TestActivityLog_refreshFromStoredLogNoTokens writes only entities from 3 months ago to today, then calls
// refreshFromStoredLog. It verifies that there are no tokens loaded.
func TestActivityLog_refreshFromStoredLogNoTokens(t *testing.T) {
	a, expectedClientRecords, _ := setupActivityRecordsInStorage(t, time.Now().UTC(), true, false)
	a.SetEnable(true)

	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(context.Background(), &wg, time.Now().UTC())
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	expectedActive := &activity.EntityActivityLog{
		Clients: expectedClientRecords[1:],
	}
	expectedCurrent := &activity.EntityActivityLog{
		Clients: expectedClientRecords[len(expectedClientRecords)-1:],
	}

	currentEntities := a.GetCurrentEntities()
	if !entityRecordsEqual(t, currentEntities.Clients, expectedCurrent.Clients) {
		// we expect all segments for the current month to be loaded
		t.Errorf("bad activity entity logs loaded. expected: %v got: %v", expectedCurrent, currentEntities)
	}
	activeClients := a.core.GetActiveClientsList()
	if err := ActiveEntitiesEqual(activeClients, expectedActive.Clients); err != nil {
		t.Error(err)
	}

	// we expect no tokens
	nsCount := a.GetStoredTokenCountByNamespaceID()
	if len(nsCount) > 0 {
		t.Errorf("expected no token counts to be loaded. got: %v", nsCount)
	}
}

// TestActivityLog_refreshFromStoredLogNoEntities writes only direct tokens from 3 months ago to today, and runs
// refreshFromStoredLog. It verifies that there are no entities or clients loaded.
func TestActivityLog_refreshFromStoredLogNoEntities(t *testing.T) {
	a, _, expectedTokenCounts := setupActivityRecordsInStorage(t, time.Now().UTC(), false, true)
	a.SetEnable(true)

	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(context.Background(), &wg, time.Now().UTC())
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	nsCount := a.GetStoredTokenCountByNamespaceID()
	if !reflect.DeepEqual(nsCount, expectedTokenCounts) {
		// we expect all token counts to be loaded
		t.Errorf("bad activity token counts loaded. expected: %v got: %v", expectedTokenCounts, nsCount)
	}

	currentEntities := a.GetCurrentEntities()
	if len(currentEntities.Clients) > 0 {
		t.Errorf("expected no current entity segment to be loaded. got: %v", currentEntities)
	}
	activeClients := a.core.GetActiveClientsList()
	if len(activeClients) > 0 {
		t.Errorf("expected no active entity segment to be loaded. got: %v", activeClients)
	}
}

// TestActivityLog_refreshFromStoredLogNoData writes nothing and calls refreshFromStoredLog, and verifies that the
// current segment counts are zero.
func TestActivityLog_refreshFromStoredLogNoData(t *testing.T) {
	now := time.Now().UTC()
	a, _, _ := setupActivityRecordsInStorage(t, now, false, false)
	a.SetEnable(true)

	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(context.Background(), &wg, now)
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	a.ExpectCurrentSegmentRefreshed(t, now.Unix(), false)
}

// TestActivityLog_refreshFromStoredLogTwoMonthsPrevious creates segment data from 5 months ago to 2 months ago and
// calls refreshFromStoredLog, then verifies that the current segment counts are zero.
func TestActivityLog_refreshFromStoredLogTwoMonthsPrevious(t *testing.T) {
	// test what happens when the most recent data is from month M-2 (or earlier - same effect)
	now := time.Now().UTC()
	twoMonthsAgoStart := timeutil.StartOfPreviousMonth(timeutil.StartOfPreviousMonth(now))
	a, _, _ := setupActivityRecordsInStorage(t, twoMonthsAgoStart, true, true)
	a.SetEnable(true)

	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(context.Background(), &wg, now)
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	a.ExpectCurrentSegmentRefreshed(t, now.Unix(), false)
}

// TestActivityLog_refreshFromStoredLogPreviousMonth creates segment data from 4 months ago to 1 month ago, then calls
// refreshFromStoredLog, then verifies that these clients are included in the current segment.
func TestActivityLog_refreshFromStoredLogPreviousMonth(t *testing.T) {
	// test what happens when most recent data is from month M-1
	// we expect to load the data from the previous month so that the activeFragmentWorker
	// can handle end of month rotations
	monthStart := timeutil.StartOfMonth(time.Now().UTC())
	oneMonthAgoStart := timeutil.StartOfPreviousMonth(monthStart)
	a, expectedClientRecords, expectedTokenCounts := setupActivityRecordsInStorage(t, oneMonthAgoStart, true, true)
	a.SetEnable(true)

	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(context.Background(), &wg, time.Now().UTC())
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	expectedActive := &activity.EntityActivityLog{
		Clients: expectedClientRecords[1:],
	}
	expectedCurrent := &activity.EntityActivityLog{
		Clients: expectedClientRecords[len(expectedClientRecords)-1:],
	}

	currentEntities := a.GetCurrentEntities()
	if !entityRecordsEqual(t, currentEntities.Clients, expectedCurrent.Clients) {
		// we only expect the newest entity segment to be loaded (for the current month)
		t.Errorf("bad activity entity logs loaded. expected: %v got: %v", expectedCurrent, currentEntities)
	}

	nsCount := a.GetStoredTokenCountByNamespaceID()
	if !reflect.DeepEqual(nsCount, expectedTokenCounts) {
		// we expect all token counts to be loaded
		t.Errorf("bad activity token counts loaded. expected: %v got: %v", expectedTokenCounts, nsCount)
	}

	activeClients := a.core.GetActiveClientsList()
	if err := ActiveEntitiesEqual(activeClients, expectedActive.Clients); err != nil {
		// we expect activeClients to be loaded for the entire month
		t.Error(err)
	}
}

type fakeResponseWriter struct {
	buffer  *bytes.Buffer
	headers http.Header
}

func (f *fakeResponseWriter) Write(b []byte) (int, error) {
	return f.buffer.Write(b)
}

func (f *fakeResponseWriter) Header() http.Header {
	return f.headers
}

func (f *fakeResponseWriter) WriteHeader(statusCode int) {
	panic("unimplmeneted")
}

// TestActivityLog_IncludeNamespace verifies that includeInResponse returns true for namespaces that are children of
// their parents.
func TestActivityLog_IncludeNamespace(t *testing.T) {
	root := namespace.RootNamespace
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog

	nsA := &namespace.Namespace{
		ID:   "aaaaa",
		Path: "a/",
	}
	nsC := &namespace.Namespace{
		ID:   "ccccc",
		Path: "c/",
	}
	nsAB := &namespace.Namespace{
		ID:   "bbbbb",
		Path: "a/b/",
	}

	testCases := []struct {
		QueryNS  *namespace.Namespace
		RecordNS *namespace.Namespace
		Expected bool
	}{
		{root, nil, true},
		{root, root, true},
		{root, nsA, true},
		{root, nsAB, true},
		{nsA, nsA, true},
		{nsA, nsAB, true},
		{nsAB, nsAB, true},

		{nsA, root, false},
		{nsA, nil, false},
		{nsAB, root, false},
		{nsAB, nil, false},
		{nsAB, nsA, false},
		{nsC, nsA, false},
		{nsC, nsAB, false},
	}

	for _, tc := range testCases {
		if a.includeInResponse(tc.QueryNS, tc.RecordNS) != tc.Expected {
			t.Errorf("bad response for query %v record %v, expected %v",
				tc.QueryNS, tc.RecordNS, tc.Expected)
		}
	}
}

// TestActivityLog_DeleteWorker writes segments for entities and direct tokens for 2 different timestamps, then runs the
// deleteLogWorker for one of the timestamps. The test verifies that the correct segment is deleted, and the other remains.
func TestActivityLog_DeleteWorker(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog

	paths := []string{
		"entity/1111/1",
		"entity/1111/2",
		"entity/1111/3",
		"entity/1112/1",
		"directtokens/1111/1",
		"directtokens/1112/1",
	}
	for _, path := range paths {
		WriteToStorage(t, core, ActivityLogPrefix+path, []byte("test"))
	}

	doneCh := make(chan struct{})
	timeout := time.After(20 * time.Second)

	go a.deleteLogWorker(namespace.RootContext(nil), 1111, doneCh)
	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatalf("timed out")
	}

	// Check segments still present
	readSegmentFromStorage(t, core, ActivityLogPrefix+"entity/1112/1")
	readSegmentFromStorage(t, core, ActivityLogPrefix+"directtokens/1112/1")

	// Check other segments not present
	expectMissingSegment(t, core, ActivityLogPrefix+"entity/1111/1")
	expectMissingSegment(t, core, ActivityLogPrefix+"entity/1111/2")
	expectMissingSegment(t, core, ActivityLogPrefix+"entity/1111/3")
	expectMissingSegment(t, core, ActivityLogPrefix+"directtokens/1111/1")
}

// checkAPIWarnings ensures there is a warning if switching from enabled -> disabled,
// and no response otherwise
func checkAPIWarnings(t *testing.T, originalEnabled, newEnabled bool, resp *logical.Response) {
	t.Helper()

	expectWarning := originalEnabled == true && newEnabled == false

	switch {
	case !expectWarning && resp != nil:
		t.Fatalf("got unexpected response: %#v", resp)
	case expectWarning && resp == nil:
		t.Fatal("expected response (containing warning) when switching from enabled to disabled")
	case expectWarning && len(resp.Warnings) == 0:
		t.Fatal("expected warning when switching from enabled to disabled")
	}
}

// TestActivityLog_EnableDisable writes a segment, adds an entity to the in-memory fragment, then disables the activity
// log. The test verifies that activity log cannot be disabled if manual reporting is enabled and no segment data is lost.
// If manual reporting is not enabled(OSS), The test verifies that the segment doesn't exist. The activity log is enabled, then verified that an empty
// segment is written and new clients can be added and written to segments.
func TestActivityLog_EnableDisable(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)

	core, b, _ := testCoreSystemBackend(t)
	a := core.activityLog
	view := core.systemBarrierView
	ctx := namespace.RootContext(nil)

	enableRequest := func() {
		t.Helper()
		originalEnabled := a.GetEnabled()

		req := logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
		req.Storage = view
		req.Data["enabled"] = "enable"
		resp, err := b.HandleRequest(ctx, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		// don't really need originalEnabled, but might as well be correct
		checkAPIWarnings(t, originalEnabled, true, resp)
	}
	disableRequest := func() {
		t.Helper()
		originalEnabled := a.GetEnabled()

		req := logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
		req.Storage = view
		req.Data["enabled"] = "disable"
		resp, err := b.HandleRequest(ctx, req)
		if a.core.ManualLicenseReportingEnabled() {
			if err == nil {
				t.Fatal("expected error")
			}
			if resp.Data["error"] != `cannot disable the activity log while Reporting is enabled` {
				t.Fatalf("bad: %v", resp)
			}
		} else {
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			checkAPIWarnings(t, originalEnabled, false, resp)
		}
	}

	// enable (if not already) and write a segment
	enableRequest()

	id1 := "11111111-1111-1111-1111-111111111111"
	id2 := "22222222-2222-2222-2222-222222222222"
	id3 := "33333333-3333-3333-3333-333333333333"
	a.AddEntityToFragment(id1, "root", time.Now().Unix())
	a.AddEntityToFragment(id2, "root", time.Now().Unix())

	a.SetStartTimestamp(a.GetStartTimestamp() - 10)
	seg1 := a.GetStartTimestamp()
	err := a.saveCurrentSegmentToStorage(ctx, false)
	if err != nil {
		t.Fatal(err)
	}

	// verify segment exists
	path := fmt.Sprintf("%ventity/%v/0", ActivityLogPrefix, seg1)
	readSegmentFromStorage(t, core, path)

	// Add in-memory fragment
	a.AddEntityToFragment(id3, "root", time.Now().Unix())

	// disable and verify segment exists
	disableRequest()

	if !a.core.ManualLicenseReportingEnabled() {
		timeout := time.After(20 * time.Second)
		select {
		case <-a.deleteDone:
			break
		case <-timeout:
			t.Fatalf("timed out")
		}

		expectMissingSegment(t, core, path)
		a.ExpectCurrentSegmentRefreshed(t, 0, false)

		// enable (if not already) which force-writes an empty segment
		enableRequest()

		seg2 := a.GetStartTimestamp()
		if seg1 >= seg2 {
			t.Errorf("bad second segment timestamp, %v >= %v", seg1, seg2)
		}

		// Verify empty segments are present
		path = fmt.Sprintf("%ventity/%v/0", ActivityLogPrefix, seg2)
		readSegmentFromStorage(t, core, path)

		path = fmt.Sprintf("%vdirecttokens/%v/0", ActivityLogPrefix, seg2)
	}
	readSegmentFromStorage(t, core, path)
}

func TestActivityLog_EndOfMonth(t *testing.T) {
	// We only want *fake* end of months, *real* ones are too scary.
	timeutil.SkipAtEndOfMonth(t)

	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	ctx := namespace.RootContext(nil)

	// Make sure we're enabled.
	a.SetConfig(ctx, activityConfig{
		Enabled:             "enable",
		RetentionMonths:     12,
		DefaultReportMonths: 12,
	})

	id1 := "11111111-1111-1111-1111-111111111111"
	id2 := "22222222-2222-2222-2222-222222222222"
	id3 := "33333333-3333-3333-3333-333333333333"
	a.AddEntityToFragment(id1, "root", time.Now().Unix())

	month0 := time.Now().UTC()
	segment0 := a.GetStartTimestamp()
	month1 := timeutil.StartOfNextMonth(month0)
	month2 := timeutil.StartOfNextMonth(month1)

	// Trigger end-of-month
	a.HandleEndOfMonth(ctx, month1)

	// Check segment is present, with 1 entity
	path := fmt.Sprintf("%ventity/%v/0", ActivityLogPrefix, segment0)
	protoSegment := readSegmentFromStorage(t, core, path)
	out := &activity.EntityActivityLog{}
	err := proto.Unmarshal(protoSegment.Value, out)
	if err != nil {
		t.Fatal(err)
	}

	segment1 := a.GetStartTimestamp()
	expectedTimestamp := timeutil.StartOfMonth(month1).Unix()
	if segment1 != expectedTimestamp {
		t.Errorf("expected segment timestamp %v got %v", expectedTimestamp, segment1)
	}

	// Check intent log is present
	intentRaw, err := core.barrier.Get(ctx, "sys/counters/activity/endofmonth")
	if err != nil {
		t.Fatal(err)
	}
	if intentRaw == nil {
		t.Fatal("no intent log present")
	}

	var intent ActivityIntentLog
	err = intentRaw.DecodeJSON(&intent)
	if err != nil {
		t.Fatal(err)
	}

	if intent.PreviousMonth != segment0 {
		t.Errorf("expected previous month %v got %v", segment0, intent.PreviousMonth)
	}

	if intent.NextMonth != segment1 {
		t.Errorf("expected previous month %v got %v", segment1, intent.NextMonth)
	}

	a.AddEntityToFragment(id2, "root", time.Now().Unix())

	a.HandleEndOfMonth(ctx, month2)
	segment2 := a.GetStartTimestamp()

	a.AddEntityToFragment(id3, "root", time.Now().Unix())

	err = a.saveCurrentSegmentToStorage(ctx, false)
	if err != nil {
		t.Fatal(err)
	}

	// Check all three segments still present, with correct entities
	testCases := []struct {
		SegmentTimestamp  int64
		ExpectedEntityIDs []string
	}{
		{segment0, []string{id1}},
		{segment1, []string{id2}},
		{segment2, []string{id3}},
	}

	for i, tc := range testCases {
		t.Logf("checking segment %v timestamp %v", i, tc.SegmentTimestamp)
		path := fmt.Sprintf("%ventity/%v/0", ActivityLogPrefix, tc.SegmentTimestamp)
		protoSegment := readSegmentFromStorage(t, core, path)
		out := &activity.EntityActivityLog{}
		err = proto.Unmarshal(protoSegment.Value, out)
		if err != nil {
			t.Fatalf("could not unmarshal protobuf: %v", err)
		}
		expectedEntityIDs(t, out, tc.ExpectedEntityIDs)
	}
}

// TestActivityLog_CalculatePrecomputedQueriesWithMixedTWEs tests that precomputed
// queries work when new months have tokens without entities saved in the TokenCount,
// as clients, or both.
func TestActivityLog_CalculatePrecomputedQueriesWithMixedTWEs(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)

	// root namespace will have TWEs with clientIDs and untracked TWEs
	// ns1 namespace will only have TWEs with clientIDs
	// aaaa, bbbb, and cccc namespace will only have untracked TWEs
	// 1. January tests clientIDs from a segment don't roll over into another month's
	// client counts in same segment.
	// 2. August tests that client counts work when split across segment.
	// 3. September tests that an entire segment in a month yields correct cc.
	// 4. October tests that a month with only a segment rolled over from previous
	// month yields correct client count.

	january := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	august := time.Date(2020, 8, 15, 12, 0, 0, 0, time.UTC)
	september := timeutil.StartOfMonth(time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC))
	october := timeutil.StartOfMonth(time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC))
	november := timeutil.StartOfMonth(time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC))

	conf := &CoreConfig{
		ActivityLogConfig: ActivityLogCoreConfig{
			ForceEnable:   true,
			DisableTimers: true,
		},
	}
	sink := SetupMetrics(conf)
	core, _, _ := TestCoreUnsealedWithConfig(t, conf)
	a := core.activityLog
	<-a.computationWorkerDone
	ctx := namespace.RootContext(nil)

	// Generate overlapping sets of entity IDs from this list.

	clientRecords := make([]*activity.EntityRecord, 45)
	clientNamespaces := []string{"root", "aaaaa", "bbbbb", "root", "root", "ccccc", "root", "bbbbb", "rrrrr"}

	for i := range clientRecords {
		clientRecords[i] = &activity.EntityRecord{
			ClientID:    fmt.Sprintf("111122222-3333-4444-5555-%012v", i),
			NamespaceID: clientNamespaces[i/5],
			Timestamp:   time.Now().Unix(),
			NonEntity:   true,
		}
	}

	toInsert := []struct {
		StartTime int64
		Segment   uint64
		Clients   []*activity.EntityRecord
	}{
		// January, should not be included
		{
			january.Unix(),
			0,
			clientRecords[40:45],
		},
		{
			august.Unix(),
			0,
			clientRecords[:13],
		},
		{
			august.Unix(),
			1,
			clientRecords[13:20],
		},
		{
			september.Unix(),
			1,
			clientRecords[10:30],
		},
		{
			september.Unix(),
			2,
			clientRecords[15:40],
		},
		{
			september.Unix(),
			3,
			clientRecords[15:40],
		},
		{
			october.Unix(),
			3,
			clientRecords[17:23],
		},
	}

	// Insert token counts for all 3 segments
	toInsertTokenCount := []struct {
		StartTime          int64
		Segment            uint64
		CountByNamespaceID map[string]uint64
	}{
		{
			january.Unix(),
			0,
			map[string]uint64{"root": 3, "ns1": 5},
		},
		{
			august.Unix(),
			0,
			map[string]uint64{"root": 40, "ns1": 50},
		},
		{
			august.Unix(),
			1,
			map[string]uint64{"root": 60, "ns1": 70},
		},
		{
			september.Unix(),
			1,
			map[string]uint64{"root": 400, "ns1": 500},
		},
		{
			september.Unix(),
			2,
			map[string]uint64{"root": 700, "ns1": 800},
		},
		{
			september.Unix(),
			3,
			map[string]uint64{"root": 0, "ns1": 0},
		},
		{
			october.Unix(),
			3,
			map[string]uint64{"root": 0, "ns1": 0},
		},
	}
	doInsertTokens := func(i int) {
		segment := toInsertTokenCount[i]
		tc := &activity.TokenCount{
			CountByNamespaceID: segment.CountByNamespaceID,
		}
		data, err := proto.Marshal(tc)
		if err != nil {
			t.Fatal(err)
		}
		tokenPath := fmt.Sprintf("%vdirecttokens/%v/%v", ActivityLogPrefix, segment.StartTime, segment.Segment)
		WriteToStorage(t, core, tokenPath, data)
	}

	// Note that precomputedQuery worker doesn't filter
	// for times <= the one it was asked to do. Is that a problem?
	// Here, it means that we can't insert everything *first* and do multiple
	// test cases, we have to write logs incrementally.
	doInsert := func(i int) {
		segment := toInsert[i]
		eal := &activity.EntityActivityLog{
			Clients: segment.Clients,
		}
		data, err := proto.Marshal(eal)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("%ventity/%v/%v", ActivityLogPrefix, segment.StartTime, segment.Segment)
		WriteToStorage(t, core, path, data)
	}
	expectedCounts := []struct {
		StartTime   time.Time
		EndTime     time.Time
		ByNamespace map[string]int
	}{
		// First test case
		{
			august,
			timeutil.EndOfMonth(august),
			map[string]int{
				"root":  110, // 10 TWEs + 50 + 60 direct tokens
				"ns1":   120, // 60 + 70 direct tokens
				"aaaaa": 5,
				"bbbbb": 5,
			},
		},
		// Second test case
		{
			august,
			timeutil.EndOfMonth(september),
			map[string]int{
				"root":  1220, // 110 from august + 10 non-overlapping TWEs in September, + 400 + 700 direct tokens in september
				"ns1":   1420, // 120 from August + 500 + 800 direct tokens in september
				"aaaaa": 5,
				"bbbbb": 10,
				"ccccc": 5,
			},
		},
		{
			september,
			timeutil.EndOfMonth(september),
			map[string]int{
				"root":  1115, // 15 TWEs in September, + 400 + 700 direct tokens
				"ns1":   1300, // 500 direct tokens in september
				"bbbbb": 10,
				"ccccc": 5,
			},
		},
		// Third test case
		{
			august,
			timeutil.EndOfMonth(october),
			map[string]int{
				"root":  1220, // 1220 from Aug to Sept
				"ns1":   1420, // 1420 from Aug to Sept
				"aaaaa": 5,
				"bbbbb": 10,
				"ccccc": 5,
			},
		},
		{
			september,
			timeutil.EndOfMonth(october),
			map[string]int{
				"root":  1115, // 1115 in Sept
				"ns1":   1300, // 1300 in Sept
				"bbbbb": 10,
				"ccccc": 5,
			},
		},
		{
			october,
			timeutil.EndOfMonth(october),
			map[string]int{
				"root": 6, // 6 overlapping TWEs in October
				"ns1":  0, // No new direct tokens in october
			},
		},
	}

	checkPrecomputedQuery := func(i int) {
		t.Helper()
		pq, err := a.queryStore.Get(ctx, expectedCounts[i].StartTime, expectedCounts[i].EndTime)
		if err != nil {
			t.Fatal(err)
		}
		if pq == nil {
			t.Errorf("empty result for %v -- %v", expectedCounts[i].StartTime, expectedCounts[i].EndTime)
		}
		if len(pq.Namespaces) != len(expectedCounts[i].ByNamespace) {
			t.Errorf("mismatched number of namespaces, expected %v got %v",
				len(expectedCounts[i].ByNamespace), len(pq.Namespaces))
		}
		for _, nsRecord := range pq.Namespaces {
			val, ok := expectedCounts[i].ByNamespace[nsRecord.NamespaceID]
			if !ok {
				t.Errorf("unexpected namespace %v", nsRecord.NamespaceID)
				continue
			}
			if uint64(val) != nsRecord.NonEntityTokens {
				t.Errorf("wrong number of entities in %v: expected %v, got %v",
					nsRecord.NamespaceID, val, nsRecord.NonEntityTokens)
			}
		}
		if !pq.StartTime.Equal(expectedCounts[i].StartTime) {
			t.Errorf("mismatched start time: expected %v got %v",
				expectedCounts[i].StartTime, pq.StartTime)
		}
		if !pq.EndTime.Equal(expectedCounts[i].EndTime) {
			t.Errorf("mismatched end time: expected %v got %v",
				expectedCounts[i].EndTime, pq.EndTime)
		}
	}

	testCases := []struct {
		InsertUpTo   int // index in the toInsert array
		PrevMonth    int64
		NextMonth    int64
		ExpectedUpTo int // index in the expectedCounts array
	}{
		{
			2, // jan-august
			august.Unix(),
			september.Unix(),
			0, // august-august
		},
		{
			5, // jan-sept
			september.Unix(),
			october.Unix(),
			2, // august-september
		},
		{
			6, // jan-oct
			october.Unix(),
			november.Unix(),
			5, // august-october
		},
	}

	inserted := -1
	for _, tc := range testCases {
		t.Logf("tc %+v", tc)

		// Persists across loops
		for inserted < tc.InsertUpTo {
			inserted += 1
			t.Logf("inserting segment %v", inserted)
			doInsert(inserted)
			doInsertTokens(inserted)
		}

		intent := &ActivityIntentLog{
			PreviousMonth: tc.PrevMonth,
			NextMonth:     tc.NextMonth,
		}
		data, err := json.Marshal(intent)
		if err != nil {
			t.Fatal(err)
		}
		WriteToStorage(t, core, "sys/counters/activity/endofmonth", data)

		// Pretend we've successfully rolled over to the following month
		a.SetStartTimestamp(tc.NextMonth)

		err = a.precomputedQueryWorker(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}

		expectMissingSegment(t, core, "sys/counters/activity/endofmonth")

		for i := 0; i <= tc.ExpectedUpTo; i++ {
			checkPrecomputedQuery(i)
		}
	}

	// Check metrics on the last precomputed query
	// (otherwise we need a way to reset the in-memory metrics between test cases.)

	intervals := sink.Data()
	// Test crossed an interval boundary, don't try to deal with it.
	if len(intervals) > 1 {
		t.Skip("Detected interval crossing.")
	}
	expectedGauges := []struct {
		Name           string
		NamespaceLabel string
		Value          float32
	}{
		// october values
		{
			"identity.nonentity.active.monthly",
			"root",
			6.0,
		},
		{
			"identity.nonentity.active.monthly",
			"deleted-bbbbb", // No namespace entry for this fake ID
			10.0,
		},
		{
			"identity.nonentity.active.monthly",
			"deleted-ccccc",
			5.0,
		},
		// january-september values
		{
			"identity.nonentity.active.reporting_period",
			"root",
			1223.0,
		},
		{
			"identity.nonentity.active.reporting_period",
			"deleted-aaaaa",
			5.0,
		},
		{
			"identity.nonentity.active.reporting_period",
			"deleted-bbbbb",
			10.0,
		},
		{
			"identity.nonentity.active.reporting_period",
			"deleted-ccccc",
			5.0,
		},
	}
	for _, g := range expectedGauges {
		found := false
		for _, actual := range intervals[0].Gauges {
			actualNamespaceLabel := ""
			for _, l := range actual.Labels {
				if l.Name == "namespace" {
					actualNamespaceLabel = l.Value
					break
				}
			}
			if actual.Name == g.Name && actualNamespaceLabel == g.NamespaceLabel {
				found = true
				if actual.Value != g.Value {
					t.Errorf("Mismatched value for %v %v %v != %v",
						g.Name, g.NamespaceLabel, actual.Value, g.Value)
				}
				break
			}
		}
		if !found {
			t.Errorf("No gauge found for %v %v",
				g.Name, g.NamespaceLabel)
		}
	}
}

func TestActivityLog_SaveAfterDisable(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	a := core.activityLog
	a.SetConfig(ctx, activityConfig{
		Enabled:             "enable",
		RetentionMonths:     12,
		DefaultReportMonths: 12,
	})

	a.AddEntityToFragment("1111-1111-11111111", "root", time.Now().Unix())
	startTimestamp := a.GetStartTimestamp()

	// This kicks off an asynchronous delete
	a.SetConfig(ctx, activityConfig{
		Enabled:             "disable",
		RetentionMonths:     12,
		DefaultReportMonths: 12,
	})

	timer := time.After(10 * time.Second)
	select {
	case <-timer:
		t.Fatal("timeout waiting for delete to finish")
	case <-a.deleteDone:
		break
	}

	// Segment should not be written even with force
	err := a.saveCurrentSegmentToStorage(context.Background(), true)
	if err != nil {
		t.Fatal(err)
	}

	path := ActivityLogPrefix + "entity/0/0"
	expectMissingSegment(t, core, path)

	path = fmt.Sprintf("%ventity/%v/0", ActivityLogPrefix, startTimestamp)
	expectMissingSegment(t, core, path)
}

// TestActivityLog_Precompute creates segments over a range of 11 months, with overlapping clients and namespaces.
// Create intent logs and run precomputedQueryWorker for various month ranges. Verify that the precomputed queries have
// the correct counts, including per namespace.
func TestActivityLog_Precompute(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)

	january := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	august := time.Date(2020, 8, 15, 12, 0, 0, 0, time.UTC)
	september := timeutil.StartOfMonth(time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC))
	october := timeutil.StartOfMonth(time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC))
	november := timeutil.StartOfMonth(time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC))

	conf := &CoreConfig{
		ActivityLogConfig: ActivityLogCoreConfig{
			ForceEnable:   true,
			DisableTimers: true,
		},
	}
	sink := SetupMetrics(conf)
	core, _, _ := TestCoreUnsealedWithConfig(t, conf)
	a := core.activityLog
	ctx := namespace.RootContext(nil)

	// Generate overlapping sets of entity IDs from this list.
	//   january:      40-44                                          RRRRR
	//   first month:   0-19  RRRRRAAAAABBBBBRRRRR
	//   second month: 10-29            BBBBBRRRRRRRRRRCCCCC
	//   third month:  15-39                 RRRRRRRRRRCCCCCRRRRRBBBBB

	entityRecords := make([]*activity.EntityRecord, 45)
	entityNamespaces := []string{"root", "aaaaa", "bbbbb", "root", "root", "ccccc", "root", "bbbbb", "rrrrr"}

	for i := range entityRecords {
		entityRecords[i] = &activity.EntityRecord{
			ClientID:    fmt.Sprintf("111122222-3333-4444-5555-%012v", i),
			NamespaceID: entityNamespaces[i/5],
			Timestamp:   time.Now().Unix(),
		}
	}

	toInsert := []struct {
		StartTime int64
		Segment   uint64
		Clients   []*activity.EntityRecord
	}{
		// January, should not be included
		{
			january.Unix(),
			0,
			entityRecords[40:45],
		},
		// Artifically split August and October
		{ // 1
			august.Unix(),
			0,
			entityRecords[:13],
		},
		{ // 2
			august.Unix(),
			1,
			entityRecords[13:20],
		},
		{ // 3
			september.Unix(),
			0,
			entityRecords[10:30],
		},
		{ // 4
			october.Unix(),
			0,
			entityRecords[15:40],
		},
		{
			october.Unix(),
			1,
			entityRecords[15:40],
		},
		{
			october.Unix(),
			2,
			entityRecords[17:23],
		},
	}

	// Note that precomputedQuery worker doesn't filter
	// for times <= the one it was asked to do. Is that a problem?
	// Here, it means that we can't insert everything *first* and do multiple
	// test cases, we have to write logs incrementally.
	doInsert := func(i int) {
		segment := toInsert[i]
		eal := &activity.EntityActivityLog{
			Clients: segment.Clients,
		}
		data, err := proto.Marshal(eal)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("%ventity/%v/%v", ActivityLogPrefix, segment.StartTime, segment.Segment)
		WriteToStorage(t, core, path, data)
	}

	expectedCounts := []struct {
		StartTime   time.Time
		EndTime     time.Time
		ByNamespace map[string]int
	}{
		// First test case
		{
			august,
			timeutil.EndOfMonth(august),
			map[string]int{
				"root":  10,
				"aaaaa": 5,
				"bbbbb": 5,
			},
		},
		// Second test case
		{
			august,
			timeutil.EndOfMonth(september),
			map[string]int{
				"root":  15,
				"aaaaa": 5,
				"bbbbb": 5,
				"ccccc": 5,
			},
		},
		{
			september,
			timeutil.EndOfMonth(september),
			map[string]int{
				"root":  10,
				"bbbbb": 5,
				"ccccc": 5,
			},
		},
		// Third test case
		{
			august,
			timeutil.EndOfMonth(october),
			map[string]int{
				"root":  20,
				"aaaaa": 5,
				"bbbbb": 10,
				"ccccc": 5,
			},
		},
		{
			september,
			timeutil.EndOfMonth(october),
			map[string]int{
				"root":  15,
				"bbbbb": 10,
				"ccccc": 5,
			},
		},
		{
			october,
			timeutil.EndOfMonth(october),
			map[string]int{
				"root":  15,
				"bbbbb": 5,
				"ccccc": 5,
			},
		},
	}

	checkPrecomputedQuery := func(i int) {
		t.Helper()
		pq, err := a.queryStore.Get(ctx, expectedCounts[i].StartTime, expectedCounts[i].EndTime)
		if err != nil {
			t.Fatal(err)
		}
		if pq == nil {
			t.Errorf("empty result for %v -- %v", expectedCounts[i].StartTime, expectedCounts[i].EndTime)
		}
		if len(pq.Namespaces) != len(expectedCounts[i].ByNamespace) {
			t.Errorf("mismatched number of namespaces, expected %v got %v",
				len(expectedCounts[i].ByNamespace), len(pq.Namespaces))
		}
		for _, nsRecord := range pq.Namespaces {
			val, ok := expectedCounts[i].ByNamespace[nsRecord.NamespaceID]
			if !ok {
				t.Errorf("unexpected namespace %v", nsRecord.NamespaceID)
				continue
			}
			if uint64(val) != nsRecord.Entities {
				t.Errorf("wrong number of entities in %v: expected %v, got %v",
					nsRecord.NamespaceID, val, nsRecord.Entities)
			}
		}
		if !pq.StartTime.Equal(expectedCounts[i].StartTime) {
			t.Errorf("mismatched start time: expected %v got %v",
				expectedCounts[i].StartTime, pq.StartTime)
		}
		if !pq.EndTime.Equal(expectedCounts[i].EndTime) {
			t.Errorf("mismatched end time: expected %v got %v",
				expectedCounts[i].EndTime, pq.EndTime)
		}
	}

	testCases := []struct {
		InsertUpTo   int // index in the toInsert array
		PrevMonth    int64
		NextMonth    int64
		ExpectedUpTo int // index in the expectedCounts array
	}{
		{
			2, // jan-august
			august.Unix(),
			september.Unix(),
			0, // august-august
		},
		{
			3, // jan-sept
			september.Unix(),
			october.Unix(),
			2, // august-september
		},
		{
			6, // jan-oct
			october.Unix(),
			november.Unix(),
			5, // august-september
		},
	}

	inserted := -1
	for _, tc := range testCases {
		t.Logf("tc %+v", tc)

		// Persists across loops
		for inserted < tc.InsertUpTo {
			inserted += 1
			t.Logf("inserting segment %v", inserted)
			doInsert(inserted)
		}

		intent := &ActivityIntentLog{
			PreviousMonth: tc.PrevMonth,
			NextMonth:     tc.NextMonth,
		}
		data, err := json.Marshal(intent)
		if err != nil {
			t.Fatal(err)
		}
		WriteToStorage(t, core, "sys/counters/activity/endofmonth", data)

		// Pretend we've successfully rolled over to the following month
		a.SetStartTimestamp(tc.NextMonth)

		err = a.precomputedQueryWorker(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}

		expectMissingSegment(t, core, "sys/counters/activity/endofmonth")

		for i := 0; i <= tc.ExpectedUpTo; i++ {
			checkPrecomputedQuery(i)
		}
	}

	// Check metrics on the last precomputed query
	// (otherwise we need a way to reset the in-memory metrics between test cases.)

	intervals := sink.Data()
	// Test crossed an interval boundary, don't try to deal with it.
	if len(intervals) > 1 {
		t.Skip("Detected interval crossing.")
	}
	expectedGauges := []struct {
		Name           string
		NamespaceLabel string
		Value          float32
	}{
		// october values
		{
			"identity.entity.active.monthly",
			"root",
			15.0,
		},
		{
			"identity.entity.active.monthly",
			"deleted-bbbbb", // No namespace entry for this fake ID
			5.0,
		},
		{
			"identity.entity.active.monthly",
			"deleted-ccccc",
			5.0,
		},
		// august-september values
		{
			"identity.entity.active.reporting_period",
			"root",
			20.0,
		},
		{
			"identity.entity.active.reporting_period",
			"deleted-aaaaa",
			5.0,
		},
		{
			"identity.entity.active.reporting_period",
			"deleted-bbbbb",
			10.0,
		},
		{
			"identity.entity.active.reporting_period",
			"deleted-ccccc",
			5.0,
		},
	}
	for _, g := range expectedGauges {
		found := false
		for _, actual := range intervals[0].Gauges {
			actualNamespaceLabel := ""
			for _, l := range actual.Labels {
				if l.Name == "namespace" {
					actualNamespaceLabel = l.Value
					break
				}
			}
			if actual.Name == g.Name && actualNamespaceLabel == g.NamespaceLabel {
				found = true
				if actual.Value != g.Value {
					t.Errorf("Mismatched value for %v %v %v != %v",
						g.Name, g.NamespaceLabel, actual.Value, g.Value)
				}
				break
			}
		}
		if !found {
			t.Errorf("No guage found for %v %v",
				g.Name, g.NamespaceLabel)
		}
	}
}

// TestActivityLog_Precompute_SkipMonth will put two non-contiguous chunks of
// data in the activity log, and then run precomputedQueryWorker. Finally it
// will perform a query get over the skip month and expect a query for the entire
// time segment (non-contiguous)
func TestActivityLog_Precompute_SkipMonth(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)

	august := time.Date(2020, 8, 15, 12, 0, 0, 0, time.UTC)
	september := timeutil.StartOfMonth(time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC))
	october := timeutil.StartOfMonth(time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC))
	november := timeutil.StartOfMonth(time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC))
	december := timeutil.StartOfMonth(time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC))

	core, _, _ := TestCoreUnsealedWithConfig(t, &CoreConfig{
		ActivityLogConfig: ActivityLogCoreConfig{
			ForceEnable:   true,
			DisableTimers: true,
		},
	})
	a := core.activityLog
	ctx := namespace.RootContext(nil)

	entityRecords := make([]*activity.EntityRecord, 45)

	for i := range entityRecords {
		entityRecords[i] = &activity.EntityRecord{
			ClientID:    fmt.Sprintf("111122222-3333-4444-5555-%012v", i),
			NamespaceID: "root",
			Timestamp:   time.Now().Unix(),
		}
	}

	toInsert := []struct {
		StartTime int64
		Segment   uint64
		Clients   []*activity.EntityRecord
	}{
		{
			august.Unix(),
			0,
			entityRecords[:20],
		},
		{
			september.Unix(),
			0,
			entityRecords[20:30],
		},
		{
			november.Unix(),
			0,
			entityRecords[30:45],
		},
	}

	// Note that precomputedQuery worker doesn't filter
	// for times <= the one it was asked to do. Is that a problem?
	// Here, it means that we can't insert everything *first* and do multiple
	// test cases, we have to write logs incrementally.
	doInsert := func(i int) {
		t.Helper()
		segment := toInsert[i]
		eal := &activity.EntityActivityLog{
			Clients: segment.Clients,
		}
		data, err := proto.Marshal(eal)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("%ventity/%v/%v", ActivityLogPrefix, segment.StartTime, segment.Segment)
		WriteToStorage(t, core, path, data)
	}

	expectedCounts := []struct {
		StartTime   time.Time
		EndTime     time.Time
		ByNamespace map[string]int
	}{
		// First test case
		{
			august,
			timeutil.EndOfMonth(september),
			map[string]int{
				"root": 30,
			},
		},
		// Second test case
		{
			august,
			timeutil.EndOfMonth(november),
			map[string]int{
				"root": 45,
			},
		},
	}

	checkPrecomputedQuery := func(i int) {
		t.Helper()
		pq, err := a.queryStore.Get(ctx, expectedCounts[i].StartTime, expectedCounts[i].EndTime)
		if err != nil {
			t.Fatal(err)
		}
		if pq == nil {
			t.Errorf("empty result for %v -- %v", expectedCounts[i].StartTime, expectedCounts[i].EndTime)
		}
		if len(pq.Namespaces) != len(expectedCounts[i].ByNamespace) {
			t.Errorf("mismatched number of namespaces, expected %v got %v",
				len(expectedCounts[i].ByNamespace), len(pq.Namespaces))
		}
		for _, nsRecord := range pq.Namespaces {
			val, ok := expectedCounts[i].ByNamespace[nsRecord.NamespaceID]
			if !ok {
				t.Errorf("unexpected namespace %v", nsRecord.NamespaceID)
				continue
			}
			if uint64(val) != nsRecord.Entities {
				t.Errorf("wrong number of entities in %v: expected %v, got %v",
					nsRecord.NamespaceID, val, nsRecord.Entities)
			}
		}
		if !pq.StartTime.Equal(expectedCounts[i].StartTime) {
			t.Errorf("mismatched start time: expected %v got %v",
				expectedCounts[i].StartTime, pq.StartTime)
		}
		if !pq.EndTime.Equal(expectedCounts[i].EndTime) {
			t.Errorf("mismatched end time: expected %v got %v",
				expectedCounts[i].EndTime, pq.EndTime)
		}
	}

	testCases := []struct {
		InsertUpTo   int // index in the toInsert array
		PrevMonth    int64
		NextMonth    int64
		ExpectedUpTo int // index in the expectedCounts array
	}{
		{
			1,
			september.Unix(),
			october.Unix(),
			0,
		},
		{
			2,
			november.Unix(),
			december.Unix(),
			1,
		},
	}

	inserted := -1
	for _, tc := range testCases {
		t.Logf("tc %+v", tc)

		// Persists across loops
		for inserted < tc.InsertUpTo {
			inserted += 1
			t.Logf("inserting segment %v", inserted)
			doInsert(inserted)
		}

		intent := &ActivityIntentLog{
			PreviousMonth: tc.PrevMonth,
			NextMonth:     tc.NextMonth,
		}
		data, err := json.Marshal(intent)
		if err != nil {
			t.Fatal(err)
		}
		WriteToStorage(t, core, "sys/counters/activity/endofmonth", data)

		// Pretend we've successfully rolled over to the following month
		a.SetStartTimestamp(tc.NextMonth)

		err = a.precomputedQueryWorker(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}

		expectMissingSegment(t, core, "sys/counters/activity/endofmonth")

		for i := 0; i <= tc.ExpectedUpTo; i++ {
			checkPrecomputedQuery(i)
		}
	}
}

// TestActivityLog_PrecomputeNonEntityTokensWithID is the same test as
// TestActivityLog_Precompute, except all the clients are tokens without
// entities. This ensures the deduplication logic and separation logic between
// entities and TWE clients is correct.
func TestActivityLog_PrecomputeNonEntityTokensWithID(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)

	january := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	august := time.Date(2020, 8, 15, 12, 0, 0, 0, time.UTC)
	september := timeutil.StartOfMonth(time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC))
	october := timeutil.StartOfMonth(time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC))
	november := timeutil.StartOfMonth(time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC))

	conf := &CoreConfig{
		ActivityLogConfig: ActivityLogCoreConfig{
			ForceEnable:   true,
			DisableTimers: true,
		},
	}
	sink := SetupMetrics(conf)
	core, _, _ := TestCoreUnsealedWithConfig(t, conf)
	a := core.activityLog
	ctx := namespace.RootContext(nil)

	// Generate overlapping sets of entity IDs from this list.
	//   january:      40-44                                          RRRRR
	//   first month:   0-19  RRRRRAAAAABBBBBRRRRR
	//   second month: 10-29            BBBBBRRRRRRRRRRCCCCC
	//   third month:  15-39                 RRRRRRRRRRCCCCCRRRRRBBBBB

	clientRecords := make([]*activity.EntityRecord, 45)
	clientNamespaces := []string{"root", "aaaaa", "bbbbb", "root", "root", "ccccc", "root", "bbbbb", "rrrrr"}

	for i := range clientRecords {
		clientRecords[i] = &activity.EntityRecord{
			ClientID:    fmt.Sprintf("111122222-3333-4444-5555-%012v", i),
			NamespaceID: clientNamespaces[i/5],
			Timestamp:   time.Now().Unix(),
			NonEntity:   true,
		}
	}

	toInsert := []struct {
		StartTime int64
		Segment   uint64
		Clients   []*activity.EntityRecord
	}{
		// January, should not be included
		{
			january.Unix(),
			0,
			clientRecords[40:45],
		},
		// Artifically split August and October
		{ // 1
			august.Unix(),
			0,
			clientRecords[:13],
		},
		{ // 2
			august.Unix(),
			1,
			clientRecords[13:20],
		},
		{ // 3
			september.Unix(),
			0,
			clientRecords[10:30],
		},
		{ // 4
			october.Unix(),
			0,
			clientRecords[15:40],
		},
		{
			october.Unix(),
			1,
			clientRecords[15:40],
		},
		{
			october.Unix(),
			2,
			clientRecords[17:23],
		},
	}

	// Note that precomputedQuery worker doesn't filter
	// for times <= the one it was asked to do. Is that a problem?
	// Here, it means that we can't insert everything *first* and do multiple
	// test cases, we have to write logs incrementally.
	doInsert := func(i int) {
		segment := toInsert[i]
		eal := &activity.EntityActivityLog{
			Clients: segment.Clients,
		}
		data, err := proto.Marshal(eal)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("%ventity/%v/%v", ActivityLogPrefix, segment.StartTime, segment.Segment)
		WriteToStorage(t, core, path, data)
	}

	expectedCounts := []struct {
		StartTime   time.Time
		EndTime     time.Time
		ByNamespace map[string]int
	}{
		// First test case
		{
			august,
			timeutil.EndOfMonth(august),
			map[string]int{
				"root":  10,
				"aaaaa": 5,
				"bbbbb": 5,
			},
		},
		// Second test case
		{
			august,
			timeutil.EndOfMonth(september),
			map[string]int{
				"root":  15,
				"aaaaa": 5,
				"bbbbb": 5,
				"ccccc": 5,
			},
		},
		{
			september,
			timeutil.EndOfMonth(september),
			map[string]int{
				"root":  10,
				"bbbbb": 5,
				"ccccc": 5,
			},
		},
		// Third test case
		{
			august,
			timeutil.EndOfMonth(october),
			map[string]int{
				"root":  20,
				"aaaaa": 5,
				"bbbbb": 10,
				"ccccc": 5,
			},
		},
		{
			september,
			timeutil.EndOfMonth(october),
			map[string]int{
				"root":  15,
				"bbbbb": 10,
				"ccccc": 5,
			},
		},
		{
			october,
			timeutil.EndOfMonth(october),
			map[string]int{
				"root":  15,
				"bbbbb": 5,
				"ccccc": 5,
			},
		},
	}

	checkPrecomputedQuery := func(i int) {
		t.Helper()
		pq, err := a.queryStore.Get(ctx, expectedCounts[i].StartTime, expectedCounts[i].EndTime)
		if err != nil {
			t.Fatal(err)
		}
		if pq == nil {
			t.Errorf("empty result for %v -- %v", expectedCounts[i].StartTime, expectedCounts[i].EndTime)
		}
		if len(pq.Namespaces) != len(expectedCounts[i].ByNamespace) {
			t.Errorf("mismatched number of namespaces, expected %v got %v",
				len(expectedCounts[i].ByNamespace), len(pq.Namespaces))
		}
		for _, nsRecord := range pq.Namespaces {
			val, ok := expectedCounts[i].ByNamespace[nsRecord.NamespaceID]
			if !ok {
				t.Errorf("unexpected namespace %v", nsRecord.NamespaceID)
				continue
			}
			if uint64(val) != nsRecord.NonEntityTokens {
				t.Errorf("wrong number of entities in %v: expected %v, got %v",
					nsRecord.NamespaceID, val, nsRecord.NonEntityTokens)
			}
		}
		if !pq.StartTime.Equal(expectedCounts[i].StartTime) {
			t.Errorf("mismatched start time: expected %v got %v",
				expectedCounts[i].StartTime, pq.StartTime)
		}
		if !pq.EndTime.Equal(expectedCounts[i].EndTime) {
			t.Errorf("mismatched end time: expected %v got %v",
				expectedCounts[i].EndTime, pq.EndTime)
		}
	}

	testCases := []struct {
		InsertUpTo   int // index in the toInsert array
		PrevMonth    int64
		NextMonth    int64
		ExpectedUpTo int // index in the expectedCounts array
	}{
		{
			2, // jan-august
			august.Unix(),
			september.Unix(),
			0, // august-august
		},
		{
			3, // jan-sept
			september.Unix(),
			october.Unix(),
			2, // august-september
		},
		{
			6, // jan-oct
			october.Unix(),
			november.Unix(),
			5, // august-september
		},
	}

	inserted := -1
	for _, tc := range testCases {
		t.Logf("tc %+v", tc)

		// Persists across loops
		for inserted < tc.InsertUpTo {
			inserted += 1
			t.Logf("inserting segment %v", inserted)
			doInsert(inserted)
		}

		intent := &ActivityIntentLog{
			PreviousMonth: tc.PrevMonth,
			NextMonth:     tc.NextMonth,
		}
		data, err := json.Marshal(intent)
		if err != nil {
			t.Fatal(err)
		}
		WriteToStorage(t, core, "sys/counters/activity/endofmonth", data)

		// Pretend we've successfully rolled over to the following month
		a.SetStartTimestamp(tc.NextMonth)

		err = a.precomputedQueryWorker(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}

		expectMissingSegment(t, core, "sys/counters/activity/endofmonth")

		for i := 0; i <= tc.ExpectedUpTo; i++ {
			checkPrecomputedQuery(i)
		}
	}

	// Check metrics on the last precomputed query
	// (otherwise we need a way to reset the in-memory metrics between test cases.)

	intervals := sink.Data()
	// Test crossed an interval boundary, don't try to deal with it.
	if len(intervals) > 1 {
		t.Skip("Detected interval crossing.")
	}
	expectedGauges := []struct {
		Name           string
		NamespaceLabel string
		Value          float32
	}{
		// october values
		{
			"identity.nonentity.active.monthly",
			"root",
			15.0,
		},
		{
			"identity.nonentity.active.monthly",
			"deleted-bbbbb", // No namespace entry for this fake ID
			5.0,
		},
		{
			"identity.nonentity.active.monthly",
			"deleted-ccccc",
			5.0,
		},
		// august-september values
		{
			"identity.nonentity.active.reporting_period",
			"root",
			20.0,
		},
		{
			"identity.nonentity.active.reporting_period",
			"deleted-aaaaa",
			5.0,
		},
		{
			"identity.nonentity.active.reporting_period",
			"deleted-bbbbb",
			10.0,
		},
		{
			"identity.nonentity.active.reporting_period",
			"deleted-ccccc",
			5.0,
		},
	}
	for _, g := range expectedGauges {
		found := false
		for _, actual := range intervals[0].Gauges {
			actualNamespaceLabel := ""
			for _, l := range actual.Labels {
				if l.Name == "namespace" {
					actualNamespaceLabel = l.Value
					break
				}
			}
			if actual.Name == g.Name && actualNamespaceLabel == g.NamespaceLabel {
				found = true
				if actual.Value != g.Value {
					t.Errorf("Mismatched value for %v %v %v != %v",
						g.Name, g.NamespaceLabel, actual.Value, g.Value)
				}
				break
			}
		}
		if !found {
			t.Errorf("No guage found for %v %v",
				g.Name, g.NamespaceLabel)
		}
	}
}

type BlockingInmemStorage struct{}

func (b *BlockingInmemStorage) List(ctx context.Context, prefix string) ([]string, error) {
	<-ctx.Done()
	return nil, errors.New("fake implementation")
}

func (b *BlockingInmemStorage) Get(ctx context.Context, key string) (*logical.StorageEntry, error) {
	<-ctx.Done()
	return nil, errors.New("fake implementation")
}

func (b *BlockingInmemStorage) Put(ctx context.Context, entry *logical.StorageEntry) error {
	<-ctx.Done()
	return errors.New("fake implementation")
}

func (b *BlockingInmemStorage) Delete(ctx context.Context, key string) error {
	<-ctx.Done()
	return errors.New("fake implementation")
}

// TestActivityLog_PrecomputeCancel stops the activity log before running the precomputedQueryWorker, and verifies that
// the context used to query storage has been canceled.
func TestActivityLog_PrecomputeCancel(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog

	// Substitute in a new view
	a.view = NewBarrierView(&BlockingInmemStorage{}, "test")

	core.stopActivityLog()

	done := make(chan struct{})

	// This will block if the shutdown didn't work.
	go func() {
		// We expect this to error because of BlockingInmemStorage
		_ = a.precomputedQueryWorker(namespace.RootContext(nil), nil)
		close(done)
	}()

	timeout := time.After(5 * time.Second)

	select {
	case <-done:
		break
	case <-timeout:
		t.Fatalf("timeout waiting for worker to finish")
	}
}

// TestActivityLog_NextMonthStart sets the activity log start timestamp, then verifies that StartOfNextMonth returns the
// correct value.
func TestActivityLog_NextMonthStart(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)

	now := time.Now().UTC()
	year, month, _ := now.Date()
	computedStart := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)

	testCases := []struct {
		SegmentStart int64
		ExpectedTime time.Time
	}{
		{
			0,
			computedStart,
		},
		{
			time.Date(2021, 2, 12, 13, 14, 15, 0, time.UTC).Unix(),
			time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC).Unix(),
			time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog

	for _, tc := range testCases {
		t.Logf("segmentStart=%v", tc.SegmentStart)
		a.SetStartTimestamp(tc.SegmentStart)

		actual := a.StartOfNextMonth()
		if !actual.Equal(tc.ExpectedTime) {
			t.Errorf("expected %v, got %v", tc.ExpectedTime, actual)
		}
	}
}

// The retention worker is called on unseal; wait for it to finish before
// proceeding with the test.
func waitForRetentionWorkerToFinish(t *testing.T, a *ActivityLog) {
	t.Helper()
	timeout := time.After(30 * time.Second)
	select {
	case <-a.retentionDone:
		return
	case <-timeout:
		t.Fatal("timeout waiting for retention worker to finish")
	}
}

// TestActivityLog_Deletion writes entity, direct tokens, and queries for dates ranging over 20 months. Then the test
// calls the retentionWorker with decreasing retention values, and verifies that the correct paths are being deleted.
func TestActivityLog_Deletion(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)

	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	waitForRetentionWorkerToFinish(t, a)

	times := []time.Time{
		time.Date(2019, 1, 15, 1, 2, 3, 0, time.UTC), // 0
		time.Date(2019, 3, 15, 1, 2, 3, 0, time.UTC),
		time.Date(2019, 4, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 7, 1, 0, 0, 0, 0, time.UTC), // 5
		time.Date(2019, 8, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 9, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 10, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 11, 1, 0, 0, 0, 0, time.UTC), // <-- 12 months starts here
		time.Date(2019, 12, 1, 0, 0, 0, 0, time.UTC), // 10
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC), // 15
		time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC), // 20
		time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC),
	}

	novIndex := len(times) - 1
	paths := make([][]string, len(times))

	for i, start := range times {
		// no entities in some months, just for fun
		for j := 0; j < (i+3)%5; j++ {
			entityPath := fmt.Sprintf("%ventity/%v/%v", ActivityLogPrefix, start.Unix(), j)
			paths[i] = append(paths[i], entityPath)
			WriteToStorage(t, core, entityPath, []byte("test"))
		}
		tokenPath := fmt.Sprintf("%vdirecttokens/%v/0", ActivityLogPrefix, start.Unix())
		paths[i] = append(paths[i], tokenPath)
		WriteToStorage(t, core, tokenPath, []byte("test"))

		// No queries for November yet
		if i < novIndex {
			for _, endTime := range times[i+1 : novIndex] {
				queryPath := fmt.Sprintf("sys/counters/activity/queries/%v/%v", start.Unix(), endTime.Unix())
				paths[i] = append(paths[i], queryPath)
				WriteToStorage(t, core, queryPath, []byte("test"))
			}
		}
	}

	checkPresent := func(i int) {
		t.Helper()
		for _, p := range paths[i] {
			readSegmentFromStorage(t, core, p)
		}
	}

	checkAbsent := func(i int) {
		t.Helper()
		for _, p := range paths[i] {
			expectMissingSegment(t, core, p)
		}
	}

	ctx := namespace.RootContext(nil)
	t.Log("24 months")
	now := times[len(times)-1]
	err := a.retentionWorker(ctx, now, 24)
	if err != nil {
		t.Fatal(err)
	}
	for i := range times {
		checkPresent(i)
	}

	t.Log("12 months")
	err = a.retentionWorker(ctx, now, 12)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i <= 8; i++ {
		checkAbsent(i)
	}
	for i := 9; i <= 21; i++ {
		checkPresent(i)
	}

	t.Log("1 month")
	err = a.retentionWorker(ctx, now, 1)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i <= 19; i++ {
		checkAbsent(i)
	}
	checkPresent(20)
	checkPresent(21)

	t.Log("0 months")
	err = a.retentionWorker(ctx, now, 0)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i <= 20; i++ {
		checkAbsent(i)
	}
	checkPresent(21)
}

// TestActivityLog_partialMonthClientCount writes segment data for the curren month and runs refreshFromStoredLog and
// then partialMonthClientCount. The test verifies that the values returned by partialMonthClientCount are correct.
func TestActivityLog_partialMonthClientCount(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)

	ctx := namespace.RootContext(nil)
	now := time.Now().UTC()
	a, clients, _ := setupActivityRecordsInStorage(t, timeutil.StartOfMonth(now), true, true)

	// clients[0] belongs to previous month
	clients = clients[1:]

	clientCounts := make(map[string]uint64)
	for _, client := range clients {
		clientCounts[client.NamespaceID] += 1
	}

	a.SetEnable(true)
	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(ctx, &wg, now)
	if err != nil {
		t.Fatalf("error loading clients: %v", err)
	}
	wg.Wait()

	results, err := a.partialMonthClientCount(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if results == nil {
		t.Fatal("no results to test")
	}

	byNamespace, ok := results["by_namespace"]
	if !ok {
		t.Fatalf("malformed results. got %v", results)
	}

	clientCountResponse := make([]*ResponseNamespace, 0)
	err = mapstructure.Decode(byNamespace, &clientCountResponse)
	if err != nil {
		t.Fatal(err)
	}

	for _, clientCount := range clientCountResponse {
		if int(clientCounts[clientCount.NamespaceID]) != clientCount.Counts.EntityClients {
			t.Errorf("bad entity count for namespace %s . expected %d, got %d", clientCount.NamespaceID, int(clientCounts[clientCount.NamespaceID]), clientCount.Counts.EntityClients)
		}
		totalCount := int(clientCounts[clientCount.NamespaceID])
		if totalCount != clientCount.Counts.Clients {
			t.Errorf("bad client count for namespace %s . expected %d, got %d", clientCount.NamespaceID, totalCount, clientCount.Counts.Clients)
		}
	}

	entityClients, ok := results["entity_clients"]
	if !ok {
		t.Fatalf("malformed results. got %v", results)
	}
	if entityClients != len(clients) {
		t.Errorf("bad entity count. expected %d, got %d", len(clients), entityClients)
	}

	clientCount, ok := results["clients"]
	if !ok {
		t.Fatalf("malformed results. got %v", results)
	}
	if clientCount != len(clients) {
		t.Errorf("bad client count. expected %d, got %d", len(clients), clientCount)
	}
}

// TestActivityLog_partialMonthClientCountUsingHandleQuery writes segments for the current month and calls
// refreshFromStoredLog, then handleQuery. The test verifies that the results from handleQuery are correct.
func TestActivityLog_partialMonthClientCountUsingHandleQuery(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)

	ctx := namespace.RootContext(nil)
	now := time.Now().UTC()
	a, clients, _ := setupActivityRecordsInStorage(t, timeutil.StartOfMonth(now), true, true)

	// clients[0] belongs to previous month
	clients = clients[1:]

	clientCounts := make(map[string]uint64)
	for _, client := range clients {
		clientCounts[client.NamespaceID] += 1
	}

	a.SetEnable(true)
	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(ctx, &wg, now)
	if err != nil {
		t.Fatalf("error loading clients: %v", err)
	}
	wg.Wait()

	results, err := a.handleQuery(ctx, time.Now().UTC(), time.Now().UTC(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if results == nil {
		t.Fatal("no results to test")
	}
	if err != nil {
		t.Fatal(err)
	}
	if results == nil {
		t.Fatal("no results to test")
	}

	byNamespace, ok := results["by_namespace"]
	if !ok {
		t.Fatalf("malformed results. got %v", results)
	}

	clientCountResponse := make([]*ResponseNamespace, 0)
	err = mapstructure.Decode(byNamespace, &clientCountResponse)
	if err != nil {
		t.Fatal(err)
	}

	for _, clientCount := range clientCountResponse {
		if int(clientCounts[clientCount.NamespaceID]) != clientCount.Counts.EntityClients {
			t.Errorf("bad entity count for namespace %s . expected %d, got %d", clientCount.NamespaceID, int(clientCounts[clientCount.NamespaceID]), clientCount.Counts.EntityClients)
		}
		totalCount := int(clientCounts[clientCount.NamespaceID])
		if totalCount != clientCount.Counts.Clients {
			t.Errorf("bad client count for namespace %s . expected %d, got %d", clientCount.NamespaceID, totalCount, clientCount.Counts.Clients)
		}
	}

	totals, ok := results["total"]
	if !ok {
		t.Fatalf("malformed results. got %v", results)
	}
	totalCounts := ResponseCounts{}
	err = mapstructure.Decode(totals, &totalCounts)
	entityClients := totalCounts.EntityClients
	if entityClients != len(clients) {
		t.Errorf("bad entity count. expected %d, got %d", len(clients), entityClients)
	}

	clientCount := totalCounts.Clients
	if clientCount != len(clients) {
		t.Errorf("bad client count. expected %d, got %d", len(clients), clientCount)
	}
	// Ensure that the month response is the same as the totals, because all clients
	// are new clients and there will be no approximation in the single month partial
	// case
	monthsRaw, ok := results["months"]
	if !ok {
		t.Fatalf("malformed results. got %v", results)
	}
	monthsResponse := make([]ResponseMonth, 0)
	err = mapstructure.Decode(monthsRaw, &monthsResponse)
	if len(monthsResponse) != 1 {
		t.Fatalf("wrong number of months returned. got %v", monthsResponse)
	}
	if monthsResponse[0].Counts.Clients != totalCounts.Clients {
		t.Fatalf("wrong client count. got %v, expected %v", monthsResponse[0].Counts.Clients, totalCounts.Clients)
	}
	if monthsResponse[0].Counts.EntityClients != totalCounts.EntityClients {
		t.Fatalf("wrong entity client count. got %v, expected %v", monthsResponse[0].Counts.EntityClients, totalCounts.EntityClients)
	}
	if monthsResponse[0].Counts.NonEntityClients != totalCounts.NonEntityClients {
		t.Fatalf("wrong non-entity client count. got %v, expected %v", monthsResponse[0].Counts.NonEntityClients, totalCounts.NonEntityClients)
	}
	if monthsResponse[0].Counts.Clients != monthsResponse[0].NewClients.Counts.Clients {
		t.Fatalf("wrong client count. got %v, expected %v", monthsResponse[0].Counts.Clients, monthsResponse[0].NewClients.Counts.Clients)
	}
	if monthsResponse[0].Counts.EntityClients != monthsResponse[0].NewClients.Counts.EntityClients {
		t.Fatalf("wrong entity client count. got %v, expected %v", monthsResponse[0].Counts.EntityClients, monthsResponse[0].NewClients.Counts.EntityClients)
	}
	if monthsResponse[0].Counts.NonEntityClients != monthsResponse[0].NewClients.Counts.NonEntityClients {
		t.Fatalf("wrong non-entity client count. got %v, expected %v", monthsResponse[0].Counts.NonEntityClients, monthsResponse[0].NewClients.Counts.NonEntityClients)
	}
	namespaceResponseMonth := monthsResponse[0].Namespaces

	for _, clientCount := range namespaceResponseMonth {
		if int(clientCounts[clientCount.NamespaceID]) != clientCount.Counts.EntityClients {
			t.Errorf("bad entity count for namespace %s . expected %d, got %d", clientCount.NamespaceID, int(clientCounts[clientCount.NamespaceID]), clientCount.Counts.EntityClients)
		}
		totalCount := int(clientCounts[clientCount.NamespaceID])
		if totalCount != clientCount.Counts.Clients {
			t.Errorf("bad client count for namespace %s . expected %d, got %d", clientCount.NamespaceID, totalCount, clientCount.Counts.Clients)
		}
	}
}

// TestActivityLog_handleQuery_normalizedMountPaths ensures that the mount paths returned by the activity log always have a trailing slash and client accounting is done correctly when there's no trailing slash.
// Two clients that have the same mount path, but one has a trailing slash, should be considered part of the same mount path.
func TestActivityLog_handleQuery_normalizedMountPaths(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)

	core, _, _ := TestCoreUnsealed(t)
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "auth/")
	ctx := namespace.RootContext(nil)
	now := time.Now().UTC()
	a := core.activityLog
	a.SetEnable(true)

	uuid1, err := uuid.GenerateUUID()
	require.NoError(t, err)
	uuid2, err := uuid.GenerateUUID()
	require.NoError(t, err)
	accessor1 := "accessor1"
	accessor2 := "accessor2"
	pathWithSlash := "auth/foo/"
	pathWithoutSlash := "auth/foo"

	// create two mounts of the same name. One has a trailing slash, the other doesn't
	err = core.router.Mount(&NoopBackend{}, "auth/foo", &MountEntry{UUID: uuid1, Accessor: accessor1, NamespaceID: namespace.RootNamespaceID, namespace: namespace.RootNamespace, Path: pathWithSlash}, view)
	require.NoError(t, err)
	err = core.router.Mount(&NoopBackend{}, "auth/bar", &MountEntry{UUID: uuid2, Accessor: accessor2, NamespaceID: namespace.RootNamespaceID, namespace: namespace.RootNamespace, Path: pathWithoutSlash}, view)
	require.NoError(t, err)

	// handle token usage for each of the mount paths
	a.HandleTokenUsage(ctx, &logical.TokenEntry{Path: pathWithSlash, NamespaceID: namespace.RootNamespaceID}, "id1", false)
	a.HandleTokenUsage(ctx, &logical.TokenEntry{Path: pathWithoutSlash, NamespaceID: namespace.RootNamespaceID}, "id2", false)
	// and have client 2 use both mount paths
	a.HandleTokenUsage(ctx, &logical.TokenEntry{Path: pathWithSlash, NamespaceID: namespace.RootNamespaceID}, "id2", false)

	// query the data for the month
	results, err := a.handleQuery(ctx, timeutil.StartOfMonth(now), timeutil.EndOfMonth(now), 0)
	require.NoError(t, err)

	byNamespace := results["by_namespace"].([]*ResponseNamespace)
	require.Len(t, byNamespace, 1)
	byMount := byNamespace[0].Mounts
	require.Len(t, byMount, 1)
	mountPath := byMount[0].MountPath

	// verify that both clients are recorded for the mount path with the slash
	require.Equal(t, mountPath, pathWithSlash)
	require.Equal(t, byMount[0].Counts.Clients, 2)
}

// TestActivityLog_partialMonthClientCountWithMultipleMountPaths verifies that logic in refreshFromStoredLog includes all mount paths
// in its mount data. In this test we create 3 entity records with different mount accessors: one is empty, one is
// valid, one can't be found (so it's assumed the mount is deleted). These records are written to storage, then this data is
// refreshed in refreshFromStoredLog, and finally we verify the results returned with partialMonthClientCount.
func TestActivityLog_partialMonthClientCountWithMultipleMountPaths(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)

	core, _, _ := TestCoreUnsealed(t)
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "auth/")

	ctx := namespace.RootContext(nil)
	now := time.Now().UTC()
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	a := core.activityLog
	path := "auth/foo/bar/"
	accessor := "authfooaccessor"

	// we mount a path using the accessor 'authfooaccessor' which has mount path "auth/foo/bar"
	// when an entity record references this accessor, activity log will be able to find it on its mounts and translate the mount accessor
	// into a mount path
	err = core.router.Mount(&NoopBackend{}, "auth/foo/", &MountEntry{UUID: meUUID, Accessor: accessor, NamespaceID: namespace.RootNamespaceID, namespace: namespace.RootNamespace, Path: path}, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	entityRecords := []*activity.EntityRecord{
		{
			// this record has no mount accessor, so it'll get recorded as a pre-1.10 upgrade
			ClientID:    "11111111-1111-1111-1111-111111111111",
			NamespaceID: namespace.RootNamespaceID,
			Timestamp:   time.Now().Unix(),
		},
		{
			// this record's mount path won't be able to be found, because there's no mount with the accessor 'deleted'
			// the code in mountAccessorToMountPath assumes that if the mount accessor isn't empty but the mount path
			// can't be found, then the mount must have been deleted
			ClientID:      "22222222-2222-2222-2222-222222222222",
			NamespaceID:   namespace.RootNamespaceID,
			Timestamp:     time.Now().Unix(),
			MountAccessor: "deleted",
		},
		{
			// this record will have mount path 'auth/foo/bar', because we set up the mount above
			ClientID:      "33333333-2222-2222-2222-222222222222",
			NamespaceID:   namespace.RootNamespaceID,
			Timestamp:     time.Now().Unix(),
			MountAccessor: "authfooaccessor",
		},
	}
	for i, entityRecord := range entityRecords {
		entityData, err := proto.Marshal(&activity.EntityActivityLog{
			Clients: []*activity.EntityRecord{entityRecord},
		})
		if err != nil {
			t.Fatal(err.Error())
		}
		storagePath := fmt.Sprintf("%sentity/%d/%d", ActivityLogPrefix, timeutil.StartOfMonth(now).Unix(), i)
		WriteToStorage(t, core, storagePath, entityData)
	}

	a.SetEnable(true)
	var wg sync.WaitGroup
	err = a.refreshFromStoredLog(ctx, &wg, now)
	if err != nil {
		t.Fatalf("error loading clients: %v", err)
	}
	wg.Wait()

	results, err := a.partialMonthClientCount(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if results == nil {
		t.Fatal("no results to test")
	}

	byNamespace, ok := results["by_namespace"]
	if !ok {
		t.Fatalf("malformed results. got %v", results)
	}

	clientCountResponse := make([]*ResponseNamespace, 0)
	err = mapstructure.Decode(byNamespace, &clientCountResponse)
	if err != nil {
		t.Fatal(err)
	}
	if len(clientCountResponse) != 1 {
		t.Fatalf("incorrect client count responses, expected 1 but got %d", len(clientCountResponse))
	}
	if len(clientCountResponse[0].Mounts) != len(entityRecords) {
		t.Fatalf("incorrect client mounts, expected %d but got %d", len(entityRecords), len(clientCountResponse[0].Mounts))
	}
	byPath := make(map[string]int, len(clientCountResponse[0].Mounts))
	for _, mount := range clientCountResponse[0].Mounts {
		byPath[mount.MountPath] = byPath[mount.MountPath] + mount.Counts.Clients
	}

	// these are the paths that are expected and correspond with the entity records created above
	expectedPaths := []string{
		noMountAccessor,
		fmt.Sprintf(DeletedMountFmt, "deleted"),
		path,
	}
	for _, expectedPath := range expectedPaths {
		count, ok := byPath[expectedPath]
		if !ok {
			t.Fatalf("path %s not found", expectedPath)
		}
		if count != 1 {
			t.Fatalf("incorrect count value %d for path %s", count, expectedPath)
		}
	}
}

// TestActivityLog_processNewClients_delete ensures that the correct clients are deleted from a processNewClients struct
func TestActivityLog_processNewClients_delete(t *testing.T) {
	mount := "mount"
	ns := "namespace"
	clientID := "client-id"
	run := func(t *testing.T, clientType string) {
		t.Helper()
		isNonEntity := clientType == nonEntityTokenActivityType || clientType == ACMEActivityType
		record := &activity.EntityRecord{
			MountAccessor: mount,
			NamespaceID:   ns,
			ClientID:      clientID,
			NonEntity:     isNonEntity,
			ClientType:    clientType,
		}
		newClients := newProcessNewClients()
		newClients.add(record)

		require.True(t, newClients.Counts.contains(record))
		require.True(t, newClients.Namespaces[ns].Counts.contains(record))
		require.True(t, newClients.Namespaces[ns].Mounts[mount].Counts.contains(record))

		newClients.delete(record)

		byNS := newClients.Namespaces
		counts := newClients.Counts
		for _, typ := range []string{nonEntityTokenActivityType, secretSyncActivityType, entityActivityType, ACMEActivityType} {
			require.NotContains(t, counts.clientsByType(typ), clientID)
			require.NotContains(t, byNS[ns].Mounts[mount].Counts.clientsByType(typ), clientID)
			require.NotContains(t, byNS[ns].Counts.clientsByType(typ), clientID)
		}
	}
	t.Run("entity", func(t *testing.T) {
		run(t, entityActivityType)
	})
	t.Run("non entity", func(t *testing.T) {
		run(t, nonEntityTokenActivityType)
	})
	t.Run("secret sync", func(t *testing.T) {
		run(t, secretSyncActivityType)
	})
	t.Run("acme", func(t *testing.T) {
		run(t, ACMEActivityType)
	})
}

// TestActivityLog_processClientRecord calls processClientRecord for an entity and a non-entity record and verifies that
// the record is present in the namespace and month maps
func TestActivityLog_processClientRecord(t *testing.T) {
	startTime := time.Now()
	mount := "mount"
	ns := "namespace"
	clientID := "client-id"

	run := func(t *testing.T, clientType string) {
		t.Helper()
		isNonEntity := clientType == nonEntityTokenActivityType || clientType == ACMEActivityType
		record := &activity.EntityRecord{
			MountAccessor: mount,
			NamespaceID:   ns,
			ClientID:      clientID,
			NonEntity:     isNonEntity,
			ClientType:    clientType,
		}
		byNS := make(summaryByNamespace)
		byMonth := make(summaryByMonth)
		processClientRecord(record, byNS, byMonth, startTime)
		require.Contains(t, byNS, ns)
		require.Contains(t, byNS[ns].Mounts, mount)
		monthIndex := timeutil.StartOfMonth(startTime).UTC().Unix()
		require.Contains(t, byMonth, monthIndex)
		require.Equal(t, byMonth[monthIndex].Namespaces, byNS)
		require.Equal(t, byMonth[monthIndex].NewClients.Namespaces, byNS)

		for _, typ := range ActivityClientTypes {
			if clientType == typ {
				require.Contains(t, byMonth[monthIndex].Counts.clientsByType(typ), clientID)
				require.Contains(t, byMonth[monthIndex].NewClients.Counts.clientsByType(typ), clientID)
				require.Contains(t, byNS[ns].Mounts[mount].Counts.clientsByType(typ), clientID)
				require.Contains(t, byNS[ns].Counts.clientsByType(typ), clientID)
			} else {
				require.NotContains(t, byMonth[monthIndex].Counts.clientsByType(typ), clientID)
				require.NotContains(t, byMonth[monthIndex].NewClients.Counts.clientsByType(typ), clientID)
				require.NotContains(t, byNS[ns].Mounts[mount].Counts.clientsByType(typ), clientID)
				require.NotContains(t, byNS[ns].Counts.clientsByType(typ), clientID)
			}
		}
	}

	t.Run("non entity", func(t *testing.T) {
		run(t, nonEntityTokenActivityType)
	})
	t.Run("entity", func(t *testing.T) {
		run(t, entityActivityType)
	})
	t.Run("secret sync", func(t *testing.T) {
		run(t, secretSyncActivityType)
	})
	t.Run("acme", func(t *testing.T) {
		run(t, ACMEActivityType)
	})
}

func verifyByNamespaceContains(t *testing.T, s summaryByNamespace, clients ...*activity.EntityRecord) {
	t.Helper()
	for _, c := range clients {
		require.Contains(t, s, c.NamespaceID)
		counts := s[c.NamespaceID].Counts
		require.True(t, counts.contains(c))
		mounts := s[c.NamespaceID].Mounts
		require.Contains(t, mounts, c.MountAccessor)
		require.True(t, mounts[c.MountAccessor].Counts.contains(c))
	}
}

func (s summaryByMonth) firstSeen(t *testing.T, client *activity.EntityRecord) time.Time {
	t.Helper()
	var seen int64
	for month, data := range s {
		present := data.NewClients.Counts.contains(client)
		if present {
			if seen != 0 {
				require.Fail(t, "client seen more than once", client.ClientID, s)
			}
			seen = month
		}
	}
	return time.Unix(seen, 0).UTC()
}

// TestActivityLog_handleEntitySegment verifies that the by namespace and by month summaries are correctly filled in a
// variety of scenarios
func TestActivityLog_handleEntitySegment(t *testing.T) {
	finalTime := timeutil.StartOfMonth(time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC))
	addMonths := func(i int) time.Time {
		return timeutil.StartOfMonth(finalTime.AddDate(0, i, 0))
	}
	currentSegmentClients := make([]*activity.EntityRecord, 0, 3)
	for i := 0; i < 3; i++ {
		currentSegmentClients = append(currentSegmentClients, &activity.EntityRecord{
			ClientID:      fmt.Sprintf("id-%d", i),
			NamespaceID:   fmt.Sprintf("ns-%d", i),
			MountAccessor: fmt.Sprintf("mnt-%d", i),
			NonEntity:     i == 0,
		})
	}
	a := &ActivityLog{}
	t.Run("older segment empty", func(t *testing.T) {
		hll := hyperloglog.New()
		byNS := make(summaryByNamespace)
		byMonth := make(summaryByMonth)
		segmentTime := addMonths(-3)
		// our 3 clients were seen 3 months ago, with no other clients having been seen
		err := a.handleEntitySegment(&activity.EntityActivityLog{Clients: currentSegmentClients}, segmentTime, hll, pqOptions{
			byNamespace:       byNS,
			byMonth:           byMonth,
			endTime:           timeutil.EndOfMonth(segmentTime),
			activePeriodStart: addMonths(-12),
			activePeriodEnd:   addMonths(12),
		})
		require.NoError(t, err)
		require.Len(t, byNS, 3)
		verifyByNamespaceContains(t, byNS, currentSegmentClients...)
		require.Len(t, byMonth, 1)
		// they should all be registered as having first been seen 3 months ago
		require.Equal(t, byMonth.firstSeen(t, currentSegmentClients[0]), segmentTime)
		require.Equal(t, byMonth.firstSeen(t, currentSegmentClients[1]), segmentTime)
		require.Equal(t, byMonth.firstSeen(t, currentSegmentClients[2]), segmentTime)
		// and all 3 should be in the hyperloglog
		require.Equal(t, hll.Estimate(), uint64(3))
	})
	t.Run("older segment clients seen earlier", func(t *testing.T) {
		hll := hyperloglog.New()
		byNS := make(summaryByNamespace)
		byNS.add(currentSegmentClients[0])
		byNS.add(currentSegmentClients[1])
		byMonth := make(summaryByMonth)
		segmentTime := addMonths(-3)
		seenBefore2Months := addMonths(-2)
		seenBefore1Month := addMonths(-1)

		// client 0 was seen 2 months ago
		byMonth.add(currentSegmentClients[0], seenBefore2Months)
		// client 1 was seen 1 month ago
		byMonth.add(currentSegmentClients[1], seenBefore1Month)

		// handle clients 0, 1, and 2 as having been seen 3 months ago
		err := a.handleEntitySegment(&activity.EntityActivityLog{Clients: currentSegmentClients}, segmentTime, hll, pqOptions{
			byNamespace:       byNS,
			byMonth:           byMonth,
			endTime:           timeutil.EndOfMonth(segmentTime),
			activePeriodStart: addMonths(-12),
			activePeriodEnd:   addMonths(12),
		})
		require.NoError(t, err)
		require.Len(t, byNS, 3)
		verifyByNamespaceContains(t, byNS, currentSegmentClients...)
		// we expect that they will only be registered as new 3 months ago, because that's when they were first seen
		require.Equal(t, byMonth.firstSeen(t, currentSegmentClients[0]), segmentTime)
		require.Equal(t, byMonth.firstSeen(t, currentSegmentClients[1]), segmentTime)
		require.Equal(t, byMonth.firstSeen(t, currentSegmentClients[2]), segmentTime)

		require.Equal(t, hll.Estimate(), uint64(3))
	})
	t.Run("disjoint set of clients", func(t *testing.T) {
		hll := hyperloglog.New()
		byNS := make(summaryByNamespace)
		byNS.add(currentSegmentClients[0])
		byNS.add(currentSegmentClients[1])
		byMonth := make(summaryByMonth)
		segmentTime := addMonths(-3)
		seenBefore2Months := addMonths(-2)
		seenBefore1Month := addMonths(-1)

		// client 0 was seen 2 months ago
		byMonth.add(currentSegmentClients[0], seenBefore2Months)
		// client 1 was seen 1 month ago
		byMonth.add(currentSegmentClients[1], seenBefore1Month)

		// handle client 2 as having been seen 3 months ago
		err := a.handleEntitySegment(&activity.EntityActivityLog{Clients: currentSegmentClients[2:]}, segmentTime, hll, pqOptions{
			byNamespace:       byNS,
			byMonth:           byMonth,
			endTime:           timeutil.EndOfMonth(segmentTime),
			activePeriodStart: addMonths(-12),
			activePeriodEnd:   addMonths(12),
		})
		require.NoError(t, err)
		require.Len(t, byNS, 3)
		verifyByNamespaceContains(t, byNS, currentSegmentClients...)
		// client 2 should be added to the map, and the other clients should stay where they were
		require.Equal(t, byMonth.firstSeen(t, currentSegmentClients[0]), seenBefore2Months)
		require.Equal(t, byMonth.firstSeen(t, currentSegmentClients[1]), seenBefore1Month)
		require.Equal(t, byMonth.firstSeen(t, currentSegmentClients[2]), segmentTime)
		// the hyperloglog will have 1 element, because there was only 1 client in the segment
		require.Equal(t, hll.Estimate(), uint64(1))
	})
	t.Run("new clients same namespaces", func(t *testing.T) {
		hll := hyperloglog.New()
		byNS := make(summaryByNamespace)
		byNS.add(currentSegmentClients[0])
		byNS.add(currentSegmentClients[1])
		byNS.add(currentSegmentClients[2])
		byMonth := make(summaryByMonth)
		segmentTime := addMonths(-3)
		seenBefore2Months := addMonths(-2)
		seenBefore1Month := addMonths(-1)

		// client 0 and 2 were seen 2 months ago
		byMonth.add(currentSegmentClients[0], seenBefore2Months)
		byMonth.add(currentSegmentClients[2], seenBefore2Months)
		// client 1 was seen 1 month ago
		byMonth.add(currentSegmentClients[1], seenBefore1Month)

		// create 3 additional clients
		// these have ns-1, ns-2, ns-3 and mnt-1, mnt-2, mnt-3
		moreSegmentClients := make([]*activity.EntityRecord, 0, 3)
		for i := 0; i < 3; i++ {
			moreSegmentClients = append(moreSegmentClients, &activity.EntityRecord{
				ClientID:      fmt.Sprintf("id-%d", i+3),
				NamespaceID:   fmt.Sprintf("ns-%d", i),
				MountAccessor: fmt.Sprintf("ns-%d", i),
				NonEntity:     i == 1,
			})
		}
		// 3 new clients have been seen 3 months ago
		err := a.handleEntitySegment(&activity.EntityActivityLog{Clients: moreSegmentClients}, segmentTime, hll, pqOptions{
			byNamespace:       byNS,
			byMonth:           byMonth,
			endTime:           timeutil.EndOfMonth(segmentTime),
			activePeriodStart: addMonths(-12),
			activePeriodEnd:   addMonths(12),
		})
		require.NoError(t, err)
		// there are only 3 namespaces, since both currentSegmentClients and moreSegmentClients use the same namespaces
		require.Len(t, byNS, 3)
		verifyByNamespaceContains(t, byNS, currentSegmentClients...)
		verifyByNamespaceContains(t, byNS, moreSegmentClients...)
		// The segment clients that have already been seen have their same first seen dates
		require.Equal(t, byMonth.firstSeen(t, currentSegmentClients[0]), seenBefore2Months)
		require.Equal(t, byMonth.firstSeen(t, currentSegmentClients[1]), seenBefore1Month)
		require.Equal(t, byMonth.firstSeen(t, currentSegmentClients[2]), seenBefore2Months)
		// and the new clients should be first seen at segmentTime
		require.Equal(t, byMonth.firstSeen(t, moreSegmentClients[0]), segmentTime)
		require.Equal(t, byMonth.firstSeen(t, moreSegmentClients[1]), segmentTime)
		require.Equal(t, byMonth.firstSeen(t, moreSegmentClients[2]), segmentTime)
		// the hyperloglog will have 3 elements, because there were the 3 new elements in moreSegmentClients seen
		require.Equal(t, hll.Estimate(), uint64(3))
	})
}

// TestActivityLog_breakdownTokenSegment verifies that tokens are correctly added to a map that tracks counts per namespace
func TestActivityLog_breakdownTokenSegment(t *testing.T) {
	toAdd := map[string]uint64{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	a := &ActivityLog{}
	testCases := []struct {
		name                    string
		existingNamespaceCounts map[string]uint64
		wantCounts              map[string]uint64
	}{
		{
			name:       "empty",
			wantCounts: toAdd,
		},
		{
			name: "some overlap",
			existingNamespaceCounts: map[string]uint64{
				"a": 2,
				"z": 1,
			},
			wantCounts: map[string]uint64{
				"a": 3,
				"b": 2,
				"c": 3,
				"z": 1,
			},
		},
		{
			name: "disjoint sets",
			existingNamespaceCounts: map[string]uint64{
				"z": 5,
				"y": 3,
				"x": 2,
			},
			wantCounts: map[string]uint64{
				"a": 1,
				"b": 2,
				"c": 3,
				"z": 5,
				"y": 3,
				"x": 2,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			byNamespace := make(map[string]*processByNamespace)
			for k, v := range tc.existingNamespaceCounts {
				byNamespace[k] = newByNamespace()
				byNamespace[k].Counts.Tokens = v
			}
			a.breakdownTokenSegment(&activity.TokenCount{CountByNamespaceID: toAdd}, byNamespace)
			got := make(map[string]uint64)
			for k, v := range byNamespace {
				got[k] = v.Counts.Tokens
			}
			require.Equal(t, tc.wantCounts, got)
		})
	}
}

// TestActivityLog_writePrecomputedQuery calls writePrecomputedQuery for a
// segment with 1 non entity, 1 entity, and 1 secret sync assoc client,
// which have different namespaces and mounts. The precomputed query is then
// retrieved from storage and we verify that the data structure is filled
// correctly
func TestActivityLog_writePrecomputedQuery(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)

	a := core.activityLog
	a.SetEnable(true)

	byMonth := make(summaryByMonth)
	byNS := make(summaryByNamespace)
	clientEntity := &activity.EntityRecord{
		ClientID:      "id-1",
		NamespaceID:   "ns-1",
		MountAccessor: "mnt-1",
	}
	clientNonEntity := &activity.EntityRecord{
		ClientID:      "id-2",
		NamespaceID:   "ns-2",
		MountAccessor: "mnt-2",
		NonEntity:     true,
	}
	secretSync := &activity.EntityRecord{
		ClientID:      "id-3",
		NamespaceID:   "ns-3",
		MountAccessor: "mnt-3",
		ClientType:    secretSyncActivityType,
	}
	acmeClient := &activity.EntityRecord{
		ClientID:      "id-4",
		NamespaceID:   "ns-4",
		MountAccessor: "mnt-4",
		ClientType:    ACMEActivityType,
	}

	now := time.Now()

	// add the 2 clients to the namespace and month summaries
	processClientRecord(clientEntity, byNS, byMonth, now)
	processClientRecord(clientNonEntity, byNS, byMonth, now)
	processClientRecord(secretSync, byNS, byMonth, now)
	processClientRecord(acmeClient, byNS, byMonth, now)

	endTime := timeutil.EndOfMonth(now)
	opts := pqOptions{
		byNamespace: byNS,
		byMonth:     byMonth,
		endTime:     endTime,
	}

	err := a.writePrecomputedQuery(context.Background(), now, opts)
	require.NoError(t, err)

	// read the query back from storage
	val, err := a.queryStore.Get(context.Background(), now, endTime)
	require.NoError(t, err)
	require.Equal(t, now.UTC().Unix(), val.StartTime.UTC().Unix())
	require.Equal(t, endTime.UTC().Unix(), val.EndTime.UTC().Unix())

	// ns-1, ns-2, ns-3, ns-4 should all be present in the results
	require.Len(t, val.Namespaces, 4)
	require.Len(t, val.Months, 1)
	resultByNS := make(map[string]*activity.NamespaceRecord)
	for _, ns := range val.Namespaces {
		resultByNS[ns.NamespaceID] = ns
	}
	ns1 := resultByNS["ns-1"]
	ns2 := resultByNS["ns-2"]
	ns3 := resultByNS["ns-3"]
	ns4 := resultByNS["ns-4"]

	require.Equal(t, ns1.Entities, uint64(1))
	require.Equal(t, ns1.NonEntityTokens, uint64(0))
	require.Equal(t, ns1.SecretSyncs, uint64(0))
	require.Equal(t, ns1.ACMEClients, uint64(0))
	require.Equal(t, ns2.Entities, uint64(0))
	require.Equal(t, ns2.NonEntityTokens, uint64(1))
	require.Equal(t, ns2.SecretSyncs, uint64(0))
	require.Equal(t, ns2.ACMEClients, uint64(0))
	require.Equal(t, ns3.Entities, uint64(0))
	require.Equal(t, ns3.NonEntityTokens, uint64(0))
	require.Equal(t, ns3.SecretSyncs, uint64(1))
	require.Equal(t, ns3.ACMEClients, uint64(0))
	require.Equal(t, ns4.Entities, uint64(0))
	require.Equal(t, ns4.NonEntityTokens, uint64(0))
	require.Equal(t, ns4.SecretSyncs, uint64(0))
	require.Equal(t, ns4.ACMEClients, uint64(1))

	require.Len(t, ns1.Mounts, 1)
	require.Len(t, ns2.Mounts, 1)
	require.Len(t, ns3.Mounts, 1)
	require.Len(t, ns4.Mounts, 1)

	// ns-1 needs to have mnt-1
	require.Contains(t, ns1.Mounts[0].MountPath, "mnt-1")
	// ns-2 needs to have mnt-2
	require.Contains(t, ns2.Mounts[0].MountPath, "mnt-2")
	// ns-3 needs to have mnt-3
	require.Contains(t, ns3.Mounts[0].MountPath, "mnt-3")
	// ns-4 needs to have mnt-4
	require.Contains(t, ns4.Mounts[0].MountPath, "mnt-4")

	// ns1 only has an entity client
	require.Equal(t, 1, ns1.Mounts[0].Counts.EntityClients)
	require.Equal(t, 0, ns1.Mounts[0].Counts.NonEntityClients)
	require.Equal(t, 0, ns1.Mounts[0].Counts.SecretSyncs)
	require.Equal(t, 0, ns1.Mounts[0].Counts.ACMEClients)

	// ns2 only has a non entity client
	require.Equal(t, 0, ns2.Mounts[0].Counts.EntityClients)
	require.Equal(t, 1, ns2.Mounts[0].Counts.NonEntityClients)
	require.Equal(t, 0, ns2.Mounts[0].Counts.SecretSyncs)
	require.Equal(t, 0, ns2.Mounts[0].Counts.ACMEClients)

	// ns3 only has a secret sync association
	require.Equal(t, 0, ns3.Mounts[0].Counts.EntityClients)
	require.Equal(t, 0, ns3.Mounts[0].Counts.NonEntityClients)
	require.Equal(t, 1, ns3.Mounts[0].Counts.SecretSyncs)
	require.Equal(t, 0, ns3.Mounts[0].Counts.ACMEClients)

	// ns4 only has an ACME client
	require.Equal(t, 0, ns4.Mounts[0].Counts.EntityClients)
	require.Equal(t, 0, ns4.Mounts[0].Counts.NonEntityClients)
	require.Equal(t, 0, ns4.Mounts[0].Counts.SecretSyncs)
	require.Equal(t, 1, ns4.Mounts[0].Counts.ACMEClients)

	monthRecord := val.Months[0]
	// there should only be one month present, since the clients were added with the same timestamp
	require.Equal(t, monthRecord.Timestamp, timeutil.StartOfMonth(now).UTC().Unix())
	require.Equal(t, 1, monthRecord.Counts.NonEntityClients)
	require.Equal(t, 1, monthRecord.Counts.EntityClients)
	require.Equal(t, 1, monthRecord.Counts.SecretSyncs)
	require.Equal(t, 1, monthRecord.Counts.ACMEClients)
	require.Len(t, monthRecord.Namespaces, 4)
	require.Len(t, monthRecord.NewClients.Namespaces, 4)
	require.Equal(t, 1, monthRecord.NewClients.Counts.EntityClients)
	require.Equal(t, 1, monthRecord.NewClients.Counts.NonEntityClients)
	require.Equal(t, 1, monthRecord.NewClients.Counts.SecretSyncs)
	require.Equal(t, 1, monthRecord.NewClients.Counts.ACMEClients)
}

type mockTimeNowClock struct {
	timeutil.DefaultClock
	start   time.Time
	created time.Time
}

func newMockTimeNowClock(startAt time.Time) timeutil.Clock {
	return &mockTimeNowClock{start: startAt, created: time.Now()}
}

// NewTimer returns a timer with a channel that will return the correct time,
// relative to the starting time. This is used when testing the
// activeFragmentWorker, as that function uses the returned value from timer.C
// to perform additional functionality
func (m mockTimeNowClock) NewTimer(d time.Duration) *time.Timer {
	timerStarted := m.Now()
	t := time.NewTimer(d)
	readCh := t.C
	writeCh := make(chan time.Time, 1)
	go func() {
		<-readCh
		writeCh <- timerStarted.Add(d)
	}()
	t.C = writeCh
	return t
}

func (m mockTimeNowClock) Now() time.Time {
	return m.start.Add(time.Since(m.created))
}

// TestActivityLog_HandleEndOfMonth runs the activity log with a mock clock.
// The current time is set to be 3 seconds before the end of a month. The test
// verifies that the precomputedQueryWorker runs and writes precomputed queries
// with the proper start and end times when the end of the month is triggered
func TestActivityLog_HandleEndOfMonth(t *testing.T) {
	// 3 seconds until a new month
	now := time.Date(2021, 1, 31, 23, 59, 57, 0, time.UTC)
	core, _, _ := TestCoreUnsealedWithConfig(t, &CoreConfig{ActivityLogConfig: ActivityLogCoreConfig{Clock: newMockTimeNowClock(now)}})
	done := make(chan struct{})
	go func() {
		defer close(done)
		<-core.activityLog.precomputedQueryWritten
	}()
	core.activityLog.SetEnable(true)
	core.activityLog.SetStartTimestamp(now.Unix())
	core.activityLog.AddClientToFragment("id", "ns", now.Unix(), false, "mount")

	// wait for the end of month to be triggered
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for precomputed query")
	}

	// verify that a precomputed query was written
	exists, err := core.activityLog.queryStore.QueriesAvailable(context.Background())
	require.NoError(t, err)
	require.True(t, exists)

	// verify that the timestamp is correct
	pq, err := core.activityLog.queryStore.Get(context.Background(), now, now.Add(24*time.Hour))
	require.NoError(t, err)
	require.Equal(t, now, pq.StartTime)
	require.Equal(t, timeutil.EndOfMonth(now), pq.EndTime)
}

// TestAddActivityToFragment calls AddActivityToFragment for different types of
// clients and verifies that they are added correctly to the tracking data
// structures
func TestAddActivityToFragment(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	a.SetEnable(true)

	mount := "mount"
	ns := "root"
	id := "id1"
	a.AddActivityToFragment(id, ns, 0, entityActivityType, mount)

	testCases := []struct {
		name         string
		id           string
		activityType string
		isAdded      bool
		expectedID   string
		isNonEntity  bool
	}{
		{
			name:         "duplicate",
			id:           id,
			activityType: entityActivityType,
			isAdded:      false,
			expectedID:   id,
		},
		{
			name:         "new entity",
			id:           "new-id",
			activityType: entityActivityType,
			isAdded:      true,
			expectedID:   "new-id",
		},
		{
			name:         "new nonentity",
			id:           "new-nonentity",
			activityType: nonEntityTokenActivityType,
			isAdded:      true,
			expectedID:   "new-nonentity",
			isNonEntity:  true,
		},
		{
			name:         "new acme",
			id:           "new-acme",
			activityType: ACMEActivityType,
			isAdded:      true,
			expectedID:   "pki-acme.new-acme",
			isNonEntity:  true,
		},
		{
			name:         "new secret sync",
			id:           "new-secret-sync",
			activityType: secretSyncActivityType,
			isAdded:      true,
			expectedID:   "new-secret-sync",
			isNonEntity:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a.fragmentLock.RLock()
			numClientsBefore := len(a.fragment.Clients)
			a.fragmentLock.RUnlock()

			a.AddActivityToFragment(tc.id, ns, 0, tc.activityType, mount)
			a.fragmentLock.RLock()
			defer a.fragmentLock.RUnlock()
			numClientsAfter := len(a.fragment.Clients)

			if tc.isAdded {
				require.Equal(t, numClientsBefore+1, numClientsAfter)
			} else {
				require.Equal(t, numClientsBefore, numClientsAfter)
			}

			require.Contains(t, a.partialMonthClientTracker, tc.expectedID)
			require.True(t, proto.Equal(&activity.EntityRecord{
				ClientID:      tc.expectedID,
				NamespaceID:   ns,
				Timestamp:     0,
				NonEntity:     tc.isNonEntity,
				MountAccessor: mount,
				ClientType:    tc.activityType,
			}, a.partialMonthClientTracker[tc.expectedID]))
		})
	}
}

// TestActivityLog_reportPrecomputedQueryMetrics creates 3 clients per type and
// calls reportPrecomputedQueryMetrics. The test verifies that the metric sink
// gets metrics reported correctly, based on the segment time matching the
// active period start or end
func TestActivityLog_reportPrecomputedQueryMetrics(t *testing.T) {
	core, _, _, metricsSink := TestCoreUnsealedWithMetrics(t)
	a := core.activityLog
	byMonth := make(summaryByMonth)
	byNS := make(summaryByNamespace)
	segmentTime := time.Now()

	// for each client type, make 3 clients in their own namespaces
	for i := 0; i < 3; i++ {
		for _, clientType := range ActivityClientTypes {
			client := &activity.EntityRecord{
				ClientID:      fmt.Sprintf("%s-%d", clientType, i),
				NamespaceID:   fmt.Sprintf("ns-%d", i),
				MountAccessor: fmt.Sprintf("mnt-%d", i),
				ClientType:    clientType,
				NonEntity:     clientType == nonEntityTokenActivityType || clientType == ACMEActivityType,
			}
			processClientRecord(client, byNS, byMonth, segmentTime)
		}
	}

	endTime := timeutil.EndOfMonth(segmentTime)
	opts := pqOptions{
		byNamespace: byNS,
		byMonth:     byMonth,
		endTime:     endTime,
	}

	otherTime := segmentTime.Add(time.Hour)
	hasNoMetric := func(t *testing.T, intervals []*metrics.IntervalMetrics, name string) {
		t.Helper()
		gauges := intervals[len(intervals)-1].Gauges
		for _, metric := range gauges {
			if metric.Name == name {
				require.Fail(t, "metric found", name)
			}
		}
	}

	hasMetric := func(t *testing.T, intervals []*metrics.IntervalMetrics, name string, value float32, namespaceLabel *string) {
		t.Helper()
		fullMetric := fmt.Sprintf("%s;cluster=test-cluster", name)
		if namespaceLabel != nil {
			fullMetric = fmt.Sprintf("%s;namespace=%s;cluster=test-cluster", name, *namespaceLabel)
		}
		gauges := intervals[len(intervals)-1].Gauges
		require.Contains(t, gauges, fullMetric)
		metric := gauges[fullMetric]
		require.Equal(t, value, metric.Value)
	}

	t.Run("no metrics", func(t *testing.T) {
		// neither option is equal to the segment time, so no metrics should be
		// reported
		opts.activePeriodStart = otherTime
		opts.activePeriodEnd = otherTime
		a.reportPrecomputedQueryMetrics(context.Background(), segmentTime, opts)

		data := metricsSink.Data()
		hasNoMetric(t, data, "identity.entity.active.monthly")
		hasNoMetric(t, data, "identity.nonentity.active.monthly")
		hasNoMetric(t, data, "identity.secret_sync.active.monthly")
		hasNoMetric(t, data, "identity.entity.active.reporting_period")
		hasNoMetric(t, data, "identity.entity.active.reporting_period")
		hasNoMetric(t, data, "identity.secret_sync.active.reporting_period")
	})

	t.Run("monthly metric", func(t *testing.T) {
		// activePeriodEnd is equal to the segment time, indicating that monthly
		// metrics should be reported
		opts.activePeriodEnd = segmentTime
		opts.activePeriodStart = otherTime
		a.reportPrecomputedQueryMetrics(context.Background(), segmentTime, opts)

		data := metricsSink.Data()
		// expect the metrics ending with "monthly"
		// the namespace was never registered in core, so it'll be
		// reported with a "deleted-" prefix
		for i := 0; i < 3; i++ {
			ns := fmt.Sprintf("deleted-ns-%d", i)
			hasMetric(t, data, "identity.entity.active.monthly", 1, &ns)
			hasMetric(t, data, "identity.nonentity.active.monthly", 1, &ns)
		}
		// secret sync metrics should be the sum of clients across all
		// namespaces
		hasMetric(t, data, "identity.secret_sync.active.monthly", 3, nil)
		hasMetric(t, data, "identity.pki_acme.active.monthly", 3, nil)
	})

	t.Run("reporting period metric", func(t *testing.T) {
		// activePeriodEnd is not equal to the segment time but activePeriodStart
		// is, which indicates that metrics for the reporting period should be
		// reported
		opts.activePeriodEnd = otherTime
		opts.activePeriodStart = segmentTime
		a.reportPrecomputedQueryMetrics(context.Background(), segmentTime, opts)

		data := metricsSink.Data()
		// expect the metrics ending with "reporting_period"
		// the namespace was never registered in core, so it'll be
		// reported with a "deleted-" prefix
		for i := 0; i < 3; i++ {
			ns := fmt.Sprintf("deleted-ns-%d", i)
			hasMetric(t, data, "identity.entity.active.reporting_period", 1, &ns)
			hasMetric(t, data, "identity.nonentity.active.reporting_period", 1, &ns)
		}
		// secret sync metrics should be the sum of clients across all
		// namespaces
		hasMetric(t, data, "identity.secret_sync.active.reporting_period", 3, nil)
		hasMetric(t, data, "identity.pki_acme.active.reporting_period", 3, nil)
	})
}

// TestActivityLog_Export_CSV_Header verifies that the export API properly
// generates a CSV column index and header. Various ActivityLogExportRecords
// are used to mimic an export discovering new map and slice fields that are
// meant to be flattened as new columns.
func TestActivityLog_Export_CSV_Header(t *testing.T) {
	encoder, err := newCSVEncoder(nil)
	require.NoError(t, err)

	expectedColumnIndex := make(map[string]int)

	// set expected index as base columnIndex upon encoder initialization
	for k, v := range encoder.columnIndex {
		expectedColumnIndex[k] = v
	}

	err = encoder.accumulateHeaderFields(&ActivityLogExportRecord{
		Policies: []string{
			"foo",
		},
		EntityMetadata: map[string]string{
			"email_address": "jdoe@abc.com",
		},
	})
	require.NoError(t, err)

	expectedColumnIndex["policies.0"] = exportCSVFlatteningInitIndex
	expectedColumnIndex["entity_metadata.email_address"] = exportCSVFlatteningInitIndex

	require.Empty(t, deep.Equal(expectedColumnIndex, encoder.columnIndex))

	err = encoder.accumulateHeaderFields(&ActivityLogExportRecord{
		Policies: []string{
			"foo",
			"bar",
			"baz",
		},
		EntityAliasCustomMetadata: map[string]string{
			"region": "west",
			"group":  "san_francisco",
		},
	})
	require.NoError(t, err)

	expectedColumnIndex["policies.1"] = exportCSVFlatteningInitIndex
	expectedColumnIndex["policies.2"] = exportCSVFlatteningInitIndex
	expectedColumnIndex["entity_alias_custom_metadata.group"] = exportCSVFlatteningInitIndex
	expectedColumnIndex["entity_alias_custom_metadata.region"] = exportCSVFlatteningInitIndex

	require.Empty(t, deep.Equal(expectedColumnIndex, encoder.columnIndex))

	err = encoder.accumulateHeaderFields(&ActivityLogExportRecord{
		Policies: []string{
			"foo",
		},
		EntityGroupIDs: []string{
			"97798e02-51e5-4ef3-906e-82c76d1a396e",
		},
		EntityMetadata: map[string]string{
			"first_name": "John",
			"last_name":  "Doe",
		},
		EntityAliasMetadata: map[string]string{
			"contact_email": "foo@abc.com",
		},
	})
	require.NoError(t, err)

	expectedColumnIndex["entity_metadata.last_name"] = exportCSVFlatteningInitIndex
	expectedColumnIndex["entity_metadata.first_name"] = exportCSVFlatteningInitIndex
	expectedColumnIndex["entity_alias_metadata.contact_email"] = exportCSVFlatteningInitIndex
	expectedColumnIndex["entity_group_ids.0"] = exportCSVFlatteningInitIndex

	require.Empty(t, deep.Equal(expectedColumnIndex, encoder.columnIndex))

	// no change because all the fields have seen before
	err = encoder.accumulateHeaderFields(&ActivityLogExportRecord{
		EntityAliasCustomMetadata: map[string]string{
			"group":  "does-not-matter",
			"region": "does-not-matter",
		},
		EntityAliasMetadata: map[string]string{
			"contact_email": "does-not-matter",
		},
		EntityGroupIDs: []string{
			"does-not-matter",
		},
		EntityMetadata: map[string]string{
			"first_name": "does-not-matter",
			"last_name":  "does-not-matter",
		},
		Policies: []string{
			"does-not-matter",
			"does-not-matter",
			"does-not-matter",
		},
	})
	require.NoError(t, err)
	require.Empty(t, deep.Equal(expectedColumnIndex, encoder.columnIndex))

	// no change because there are no slice or map fields
	err = encoder.accumulateHeaderFields(&ActivityLogExportRecord{})
	require.NoError(t, err)
	require.Empty(t, deep.Equal(expectedColumnIndex, encoder.columnIndex))

	expectedHeader := append(baseActivityExportCSVHeader(),
		"entity_alias_custom_metadata.group",
		"entity_alias_custom_metadata.region",
		"entity_alias_metadata.contact_email",
		"entity_group_ids.0",
		"entity_metadata.email_address",
		"entity_metadata.first_name",
		"entity_metadata.last_name",
		"policies.0",
		"policies.1",
		"policies.2")

	header := encoder.generateHeader()
	require.Empty(t, deep.Equal(expectedHeader, header))

	expectedColumnIndex = make(map[string]int)

	for idx, col := range expectedHeader {
		expectedColumnIndex[col] = idx
	}

	require.Empty(t, deep.Equal(expectedColumnIndex, encoder.columnIndex))
}
