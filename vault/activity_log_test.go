package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
)

const (
	logPrefix          = "sys/counters/activity/log/"
	activityPrefix     = "sys/counters/activity/"
	activityConfigPath = "sys/counters/activity/config"
)

func TestActivityLog_Creation(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)

	a := core.activityLog
	a.enabled = true

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

	if a.fragment.Entities == nil {
		t.Fatal("no fragment entity slice")
	}

	if a.fragment.NonEntityTokens == nil {
		t.Fatal("no fragment token map")
	}

	if len(a.fragment.Entities) != 1 {
		t.Fatalf("wrong number of entities %v", len(a.fragment.Entities))
	}

	er := a.fragment.Entities[0]
	if er.EntityID != entity_id {
		t.Errorf("mimatched entity ID, %q vs %q", er.EntityID, entity_id)
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

func checkExpectedEntitiesInMap(t *testing.T, a *ActivityLog, entityIDs []string) {
	t.Helper()

	activeEntities := a.core.GetActiveEntities()
	if len(activeEntities) != len(entityIDs) {
		t.Fatalf("mismatched number of entities, expected %v got %v", len(entityIDs), activeEntities)
	}
	for _, e := range entityIDs {
		if _, present := activeEntities[e]; !present {
			t.Errorf("entity ID %q is missing", e)
		}
	}
}

func TestActivityLog_UniqueEntities(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	a.enabled = true

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

	if len(a.fragment.Entities) != 2 {
		t.Fatalf("number of entities is %v", len(a.fragment.Entities))
	}

	for i, e := range a.fragment.Entities {
		expectedID := id1
		expectedTime := t1.Unix()
		expectedNS := "root"
		if i == 1 {
			expectedID = id2
			expectedTime = t2.Unix()
		}
		if e.EntityID != expectedID {
			t.Errorf("%v: expected %q, got %q", i, expectedID, e.EntityID)
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

func writeToStorage(t *testing.T, c *Core, path string, data []byte) {
	t.Helper()
	err := c.barrier.Put(context.Background(), &logical.StorageEntry{
		Key:   path,
		Value: data,
	})

	if err != nil {
		t.Fatalf("Failed to write %s to %s", data, path)
	}
}

func expectedEntityIDs(t *testing.T, out *activity.EntityActivityLog, ids []string) {
	t.Helper()

	if len(out.Entities) != len(ids) {
		t.Fatalf("entity log expected length %v, actual %v", len(ids), len(out.Entities))
	}

	// Double loop, OK for small cases
	for _, id := range ids {
		found := false
		for _, e := range out.Entities {
			if e.EntityID == id {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("did not find entity ID %v", id)
		}
	}
}

func TestActivityLog_SaveTokensToStorage(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	a.enabled = true
	// set a nonzero segment
	a.currentSegment.startTimestamp = time.Now().Unix()

	nsIDs := [...]string{"ns1_id", "ns2_id", "ns3_id"}
	path := fmt.Sprintf("%sdirecttokens/%d/0", logPrefix, a.currentSegment.startTimestamp)

	for i := 0; i < 3; i++ {
		a.AddTokenToFragment(nsIDs[0])
	}
	a.AddTokenToFragment(nsIDs[1])
	err := a.saveCurrentSegmentToStorage(context.Background(), false)
	if err != nil {
		t.Fatalf("got error writing tokens to storage: %v", err)
	}
	if a.fragment != nil {
		t.Errorf("fragment was not reset after write to storage")
	}

	protoSegment := readSegmentFromStorage(t, core, path)
	out := &activity.TokenCount{}
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
	err = a.saveCurrentSegmentToStorage(context.Background(), false)
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

func TestActivityLog_SaveEntitiesToStorage(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	a.enabled = true
	// set a nonzero segment
	a.currentSegment.startTimestamp = time.Now().Unix()

	now := time.Now()
	ids := []string{"11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222", "33333333-2222-2222-2222-222222222222"}
	times := [...]int64{
		now.Unix(),
		now.Add(1 * time.Second).Unix(),
		now.Add(2 * time.Second).Unix(),
	}
	path := fmt.Sprintf("%sentity/%d/0", logPrefix, a.currentSegment.startTimestamp)

	a.AddEntityToFragment(ids[0], "root", times[0])
	a.AddEntityToFragment(ids[1], "root2", times[1])
	err := a.saveCurrentSegmentToStorage(context.Background(), false)
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
	err = a.saveCurrentSegmentToStorage(context.Background(), false)
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

func TestActivityLog_ReceivedFragment(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	a.enabled = true

	ids := []string{
		"11111111-1111-1111-1111-111111111111",
		"22222222-2222-2222-2222-222222222222",
	}

	entityRecords := []*activity.EntityRecord{
		&activity.EntityRecord{
			EntityID:    ids[0],
			NamespaceID: "root",
			Timestamp:   time.Now().Unix(),
		},
		&activity.EntityRecord{
			EntityID:    ids[1],
			NamespaceID: "root",
			Timestamp:   time.Now().Unix(),
		},
	}

	fragment := &activity.LogFragment{
		OriginatingNode: "test-123",
		Entities:        entityRecords,
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

func TestActivityLog_availableLogsEmptyDirectory(t *testing.T) {
	// verify that directory is empty, and nothing goes wrong
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	times, err := a.availableLogs(context.Background())

	if err != nil {
		t.Fatalf("error getting start_time(s) for empty activity log")
	}
	if len(times) != 0 {
		t.Fatalf("invalid number of start_times returned. expected 0, got %d", len(times))
	}
}

func TestActivityLog_availableLogs(t *testing.T) {
	// set up a few files in storage
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	paths := [...]string{"entity/1111/1", "directtokens/1111/1", "directtokens/1000000/1", "entity/992/3", "directtokens/992/1"}
	expectedTimes := [...]time.Time{time.Unix(1000000, 0), time.Unix(1111, 0), time.Unix(992, 0)}

	for _, path := range paths {
		writeToStorage(t, core, logPrefix+path, []byte("test"))
	}

	// verify above files are there, and dates in correct order
	times, err := a.availableLogs(context.Background())
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

func TestActivityLog_MultipleFragmentsAndSegments(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog

	// enabled check is now inside AddEntityToFragment
	a.enabled = true
	// set a nonzero segment
	a.currentSegment.startTimestamp = time.Now().Unix()

	// Stop timers for test purposes
	close(a.doneCh)

	path0 := fmt.Sprintf("sys/counters/activity/log/entity/%d/0", a.currentSegment.startTimestamp)
	path1 := fmt.Sprintf("sys/counters/activity/log/entity/%d/1", a.currentSegment.startTimestamp)
	tokenPath := fmt.Sprintf("sys/counters/activity/log/directtokens/%d/0", a.currentSegment.startTimestamp)

	genID := func(i int) string {
		return fmt.Sprintf("11111111-1111-1111-1111-%012d", i)
	}
	ts := time.Now().Unix()

	// First 7000 should fit in one segment
	for i := 0; i < 7000; i++ {
		a.AddEntityToFragment(genID(i), "root", ts)
	}

	// Consume new fragment notification.
	// The worker may have gotten it first, before processing
	// the close!
	select {
	case <-a.newFragmentCh:
	default:
	}

	// Save incomplete segment
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
	if len(entityLog0.Entities) != 7000 {
		t.Fatalf("unexpected entity length. Expected %d, got %d", 7000, len(entityLog0.Entities))
	}

	// 7000 more local entities
	for i := 7000; i < 14000; i++ {
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
		Entities:        make([]*activity.EntityRecord, 0, 100),
		NonEntityTokens: tokens1,
	}
	for i := 7000; i < 7100; i++ {
		fragment1.Entities = append(fragment1.Entities, &activity.EntityRecord{
			EntityID:    genID(i),
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
		Entities:        make([]*activity.EntityRecord, 0, 100),
		NonEntityTokens: tokens2,
	}
	for i := 14000; i < 14100; i++ {
		fragment2.Entities = append(fragment2.Entities, &activity.EntityRecord{
			EntityID:    genID(i),
			NamespaceID: "root",
			Timestamp:   ts,
		})
	}
	a.receivedFragment(fragment1)
	a.receivedFragment(fragment2)

	<-a.newFragmentCh

	err = a.saveCurrentSegmentToStorage(context.Background(), false)
	if err != nil {
		t.Fatalf("got error writing entities to storage: %v", err)
	}

	if a.currentSegment.entitySequenceNumber != 1 {
		t.Fatalf("expected sequence number 1, got %v", a.currentSegment.entitySequenceNumber)
	}

	protoSegment0 = readSegmentFromStorage(t, core, path0)
	err = proto.Unmarshal(protoSegment0.Value, &entityLog0)
	if err != nil {
		t.Fatalf("could not unmarshal protobuf: %v", err)
	}
	if len(entityLog0.Entities) != activitySegmentEntityCapacity {
		t.Fatalf("unexpected entity length. Expected %d, got %d", activitySegmentEntityCapacity,
			len(entityLog0.Entities))
	}

	protoSegment1 := readSegmentFromStorage(t, core, path1)
	entityLog1 := activity.EntityActivityLog{}
	err = proto.Unmarshal(protoSegment1.Value, &entityLog1)
	if err != nil {
		t.Fatalf("could not unmarshal protobuf: %v", err)
	}
	expectedCount := 14100 - activitySegmentEntityCapacity
	if len(entityLog1.Entities) != expectedCount {
		t.Fatalf("unexpected entity length. Expected %d, got %d", expectedCount,
			len(entityLog1.Entities))
	}

	entityPresent := make(map[string]struct{})
	for _, e := range entityLog0.Entities {
		entityPresent[e.EntityID] = struct{}{}
	}
	for _, e := range entityLog1.Entities {
		entityPresent[e.EntityID] = struct{}{}
	}
	for i := 0; i < 14100; i++ {
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

func TestActivityLog_API_ConfigCRUD(t *testing.T) {
	core, b, _ := testCoreSystemBackend(t)
	view := core.systemBarrierView

	// Test reading the defaults
	{
		req := logical.TestRequest(t, logical.ReadOperation, "internal/counters/config")
		req.Storage = view
		resp, err := b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		defaults := map[string]interface{}{
			"default_report_months": 12,
			"retention_months":      24,
			"enabled":               activityLogEnabledDefaultValue,
			"queries_available":     false,
		}

		if diff := deep.Equal(resp.Data, defaults); len(diff) > 0 {
			t.Fatalf("diff: %v", diff)
		}
	}

	// Check Error Cases
	{
		req := logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
		req.Storage = view
		req.Data["default_report_months"] = 0
		_, err := b.HandleRequest(namespace.RootContext(nil), req)
		if err == nil {
			t.Fatal("expected error")
		}

		req = logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
		req.Storage = view
		req.Data["enabled"] = "bad-value"
		_, err = b.HandleRequest(namespace.RootContext(nil), req)
		if err == nil {
			t.Fatal("expected error")
		}

		req = logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
		req.Storage = view
		req.Data["retention_months"] = 0
		req.Data["enabled"] = "enable"
		_, err = b.HandleRequest(namespace.RootContext(nil), req)
		if err == nil {
			t.Fatal("expected error")
		}
	}

	// Test single key updates
	{
		req := logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
		req.Storage = view
		req.Data["default_report_months"] = 1
		resp, err := b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if resp != nil {
			t.Fatalf("bad: %#v", resp)
		}

		req = logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
		req.Storage = view
		req.Data["retention_months"] = 2
		resp, err = b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if resp != nil {
			t.Fatalf("bad: %#v", resp)
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
			"default_report_months": 1,
			"retention_months":      2,
			"enabled":               "enable",
			"queries_available":     false,
		}

		if diff := deep.Equal(resp.Data, expected); len(diff) > 0 {
			t.Fatalf("diff: %v", diff)
		}
	}

	// Test updating all keys
	{
		req := logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
		req.Storage = view
		req.Data["enabled"] = "default"
		req.Data["retention_months"] = 24
		req.Data["default_report_months"] = 12
		resp, err := b.HandleRequest(namespace.RootContext(nil), req)
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

		defaults := map[string]interface{}{
			"default_report_months": 12,
			"retention_months":      24,
			"enabled":               activityLogEnabledDefaultValue,
			"queries_available":     false,
		}

		if diff := deep.Equal(resp.Data, defaults); len(diff) > 0 {
			t.Fatalf("diff: %v", diff)
		}
	}
}

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

func TestActivityLog_getLastEntitySegmentNumber(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	paths := [...]string{"entity/992/0", "entity/1000/-1", "entity/1001/foo", "entity/1111/0", "entity/1111/1"}
	for _, path := range paths {
		writeToStorage(t, core, logPrefix+path, []byte("test"))
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

func TestActivityLog_tokenCountExists(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	paths := [...]string{"directtokens/992/0", "directtokens/1001/foo", "directtokens/1111/0", "directtokens/2222/1"}
	for _, path := range paths {
		writeToStorage(t, core, logPrefix+path, []byte("test"))
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

		idComp := strings.Compare(ei.EntityID, ej.EntityID)
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
		if a.EntityID != b.EntityID || a.NamespaceID != b.NamespaceID || a.Timestamp != b.Timestamp {
			return false
		}
	}

	return true
}

func activeEntitiesEqual(t *testing.T, active map[string]struct{}, test []*activity.EntityRecord) bool {
	t.Helper()

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

func (a *ActivityLog) resetEntitiesInMemory(t *testing.T) {
	t.Helper()

	a.l.Lock()
	defer a.l.Unlock()
	a.fragmentLock.Lock()
	defer a.fragmentLock.Unlock()
	a.currentSegment = segmentInfo{
		startTimestamp: time.Time{}.Unix(),
		currentEntities: &activity.EntityActivityLog{
			Entities: make([]*activity.EntityRecord, 0),
		},
		tokenCount:           a.currentSegment.tokenCount,
		entitySequenceNumber: 0,
	}

	a.activeEntities = make(map[string]struct{})
}

func TestActivityLog_loadCurrentEntitySegment(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog

	// we must verify that loadCurrentEntitySegment doesn't overwrite the in-memory token count
	tokenRecords := make(map[string]uint64)
	tokenRecords["test"] = 1
	tokenCount := &activity.TokenCount{
		CountByNamespaceID: tokenRecords,
	}
	a.currentSegment.tokenCount = tokenCount

	// setup in-storage data to load for testing
	entityRecords := []*activity.EntityRecord{
		&activity.EntityRecord{
			EntityID:    "11111111-1111-1111-1111-111111111111",
			NamespaceID: "root",
			Timestamp:   time.Now().Unix(),
		},
		&activity.EntityRecord{
			EntityID:    "22222222-2222-2222-2222-222222222222",
			NamespaceID: "root",
			Timestamp:   time.Now().Unix(),
		},
	}
	testEntities1 := &activity.EntityActivityLog{
		Entities: entityRecords[:1],
	}
	testEntities2 := &activity.EntityActivityLog{
		Entities: entityRecords[1:2],
	}
	testEntities3 := &activity.EntityActivityLog{
		Entities: entityRecords[:2],
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
			t.Fatalf(err.Error())
		}
		writeToStorage(t, core, logPrefix+tc.path, data)
	}

	ctx := context.Background()
	for _, tc := range testCases {
		err := a.loadCurrentEntitySegment(ctx, time.Unix(tc.time, 0), tc.seqNum)
		if err != nil {
			t.Fatalf("got error loading data for %q: %v", tc.path, err)
		}
		if !reflect.DeepEqual(a.currentSegment.tokenCount.CountByNamespaceID, tokenCount.CountByNamespaceID) {
			t.Errorf("this function should not wipe out the in-memory token count")
		}

		// verify accurate data in in-memory current segment
		if a.currentSegment.startTimestamp != tc.time {
			t.Errorf("bad timestamp loaded. expected: %v, got: %v for path %q", tc.time, a.currentSegment.startTimestamp, tc.path)
		}
		if a.currentSegment.entitySequenceNumber != tc.seqNum {
			t.Errorf("bad sequence number loaded. expected: %v, got: %v for path %q", tc.seqNum, a.currentSegment.entitySequenceNumber, tc.path)
		}
		if !entityRecordsEqual(t, a.currentSegment.currentEntities.Entities, tc.entities.Entities) {
			t.Errorf("bad data loaded. expected: %v, got: %v for path %q", tc.entities.Entities, a.currentSegment.currentEntities, tc.path)
		}

		activeEntities := core.GetActiveEntities()
		if !activeEntitiesEqual(t, activeEntities, tc.entities.Entities) {
			t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v for path %q", tc.entities.Entities, activeEntities, tc.path)
		}

		a.resetEntitiesInMemory(t)
	}
}

func TestActivityLog_loadPriorEntitySegment(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog
	a.enabled = true

	// setup in-storage data to load for testing
	entityRecords := []*activity.EntityRecord{
		&activity.EntityRecord{
			EntityID:    "11111111-1111-1111-1111-111111111111",
			NamespaceID: "root",
			Timestamp:   time.Now().Unix(),
		},
		&activity.EntityRecord{
			EntityID:    "22222222-2222-2222-2222-222222222222",
			NamespaceID: "root",
			Timestamp:   time.Now().Unix(),
		},
	}
	testEntities1 := &activity.EntityActivityLog{
		Entities: entityRecords[:1],
	}
	testEntities2 := &activity.EntityActivityLog{
		Entities: entityRecords[:2],
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
			t.Fatalf(err.Error())
		}
		writeToStorage(t, core, logPrefix+tc.path, data)
	}

	ctx := context.Background()
	for _, tc := range testCases {
		if tc.refresh {
			a.l.Lock()
			a.fragmentLock.Lock()
			a.activeEntities = make(map[string]struct{})
			a.currentSegment.startTimestamp = tc.time
			a.fragmentLock.Unlock()
			a.l.Unlock()
		}

		err := a.loadPriorEntitySegment(ctx, time.Unix(tc.time, 0), tc.seqNum)
		if err != nil {
			t.Fatalf("got error loading data for %q: %v", tc.path, err)
		}

		activeEntities := core.GetActiveEntities()
		if !activeEntitiesEqual(t, activeEntities, tc.entities.Entities) {
			t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v for path %q", tc.entities.Entities, activeEntities, tc.path)
		}
	}
}

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
		t.Fatalf(err.Error())
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

	for _, tc := range testCases {
		writeToStorage(t, core, logPrefix+tc.path, data)
	}

	ctx := context.Background()
	for _, tc := range testCases {
		// a.currentSegment.tokenCount doesn't need to be wiped each iter since it happens in loadTokenSegment()
		err := a.loadTokenCount(ctx, time.Unix(tc.time, 0))
		if err != nil {
			t.Fatalf("got error loading data for %q: %v", tc.path, err)
		}
		if !reflect.DeepEqual(a.currentSegment.tokenCount.CountByNamespaceID, tokenRecords) {
			t.Errorf("bad token count loaded. expected: %v got: %v for path %q", tokenRecords, a.currentSegment.tokenCount.CountByNamespaceID, tc.path)
		}
	}
}

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
	core.setupActivityLog(ctx)

	a = core.activityLog
	if a.currentSegment.tokenCount.CountByNamespaceID == nil {
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
			&activity.EntityRecord{
				EntityID:    "11111111-1111-1111-1111-111111111111",
				NamespaceID: "root",
				Timestamp:   time.Now().Unix(),
			},
			&activity.EntityRecord{
				EntityID:    "22222222-2222-2222-2222-222222222222",
				NamespaceID: "root",
				Timestamp:   time.Now().Unix(),
			},
			&activity.EntityRecord{
				EntityID:    "33333333-2222-2222-2222-222222222222",
				NamespaceID: "root",
				Timestamp:   time.Now().Unix(),
			},
		}

		testEntities1 := &activity.EntityActivityLog{
			Entities: entityRecords[:1],
		}
		entityData1, err := proto.Marshal(testEntities1)
		if err != nil {
			t.Fatalf(err.Error())
		}
		testEntities2 := &activity.EntityActivityLog{
			Entities: entityRecords[1:2],
		}
		entityData2, err := proto.Marshal(testEntities2)
		if err != nil {
			t.Fatalf(err.Error())
		}
		testEntities3 := &activity.EntityActivityLog{
			Entities: entityRecords[2:],
		}
		entityData3, err := proto.Marshal(testEntities3)
		if err != nil {
			t.Fatalf(err.Error())
		}

		writeToStorage(t, core, logPrefix+"entity/"+fmt.Sprint(monthsAgo.Unix())+"/0", entityData1)
		writeToStorage(t, core, logPrefix+"entity/"+fmt.Sprint(base.Unix())+"/0", entityData2)
		writeToStorage(t, core, logPrefix+"entity/"+fmt.Sprint(base.Unix())+"/1", entityData3)
	}

	var tokenRecords map[string]uint64
	if includeTokens {
		tokenRecords = make(map[string]uint64)
		for i := 1; i < 4; i++ {
			nsID := "ns" + strconv.Itoa(i)
			tokenRecords[nsID] = uint64(i)
		}
		tokenCount := &activity.TokenCount{
			CountByNamespaceID: tokenRecords,
		}

		tokenData, err := proto.Marshal(tokenCount)
		if err != nil {
			t.Fatalf(err.Error())
		}

		writeToStorage(t, core, logPrefix+"directtokens/"+fmt.Sprint(base.Unix())+"/0", tokenData)
	}

	return a, entityRecords, tokenRecords
}

func TestActivityLog_refreshFromStoredLog(t *testing.T) {
	a, expectedEntityRecords, expectedTokenCounts := setupActivityRecordsInStorage(t, time.Now().UTC(), true, true)
	a.enabled = true

	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(context.Background(), &wg)
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	expectedActive := &activity.EntityActivityLog{
		Entities: expectedEntityRecords[1:],
	}
	expectedCurrent := &activity.EntityActivityLog{
		Entities: expectedEntityRecords[2:],
	}
	if !entityRecordsEqual(t, a.currentSegment.currentEntities.Entities, expectedCurrent.Entities) {
		// we only expect the newest entity segment to be loaded (for the current month)
		t.Errorf("bad activity entity logs loaded. expected: %v got: %v", expectedCurrent, a.currentSegment.currentEntities)
	}
	if !reflect.DeepEqual(a.currentSegment.tokenCount.CountByNamespaceID, expectedTokenCounts) {
		// we expect all token counts to be loaded
		t.Errorf("bad activity token counts loaded. expected: %v got: %v", expectedTokenCounts, a.currentSegment.tokenCount.CountByNamespaceID)
	}

	activeEntities := a.core.GetActiveEntities()
	if !activeEntitiesEqual(t, activeEntities, expectedActive.Entities) {
		// we expect activeEntities to be loaded for the entire month
		t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v", expectedActive.Entities, activeEntities)
	}
}

func TestActivityLog_refreshFromStoredLogWithBackgroundLoadingCancelled(t *testing.T) {
	a, expectedEntityRecords, expectedTokenCounts := setupActivityRecordsInStorage(t, time.Now().UTC(), true, true)
	a.enabled = true

	var wg sync.WaitGroup
	close(a.doneCh)

	err := a.refreshFromStoredLog(context.Background(), &wg)
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	expected := &activity.EntityActivityLog{
		Entities: expectedEntityRecords[2:],
	}
	if !entityRecordsEqual(t, a.currentSegment.currentEntities.Entities, expected.Entities) {
		// we only expect the newest entity segment to be loaded (for the current month)
		t.Errorf("bad activity entity logs loaded. expected: %v got: %v", expected, a.currentSegment.currentEntities)
	}
	if !reflect.DeepEqual(a.currentSegment.tokenCount.CountByNamespaceID, expectedTokenCounts) {
		// we expect all token counts to be loaded
		t.Errorf("bad activity token counts loaded. expected: %v got: %v", expectedTokenCounts, a.currentSegment.tokenCount.CountByNamespaceID)
	}

	activeEntities := a.core.GetActiveEntities()
	if !activeEntitiesEqual(t, activeEntities, expected.Entities) {
		// we only expect activeEntities to be loaded for the newest segment (for the current month)
		t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v", expected.Entities, activeEntities)
	}
}

func TestActivityLog_refreshFromStoredLogContextCancelled(t *testing.T) {
	a, _, _ := setupActivityRecordsInStorage(t, time.Now().UTC(), true, true)

	var wg sync.WaitGroup
	ctx, cancelFn := context.WithCancel(context.Background())
	cancelFn()

	err := a.refreshFromStoredLog(ctx, &wg)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancelled error, got: %v", err)
	}
}

func TestActivityLog_refreshFromStoredLogNoTokens(t *testing.T) {
	a, expectedEntityRecords, _ := setupActivityRecordsInStorage(t, time.Now().UTC(), true, false)
	a.enabled = true

	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(context.Background(), &wg)
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	expectedActive := &activity.EntityActivityLog{
		Entities: expectedEntityRecords[1:],
	}
	expectedCurrent := &activity.EntityActivityLog{
		Entities: expectedEntityRecords[2:],
	}
	if !entityRecordsEqual(t, a.currentSegment.currentEntities.Entities, expectedCurrent.Entities) {
		// we expect all segments for the current month to be loaded
		t.Errorf("bad activity entity logs loaded. expected: %v got: %v", expectedCurrent, a.currentSegment.currentEntities)
	}
	activeEntities := a.core.GetActiveEntities()
	if !activeEntitiesEqual(t, activeEntities, expectedActive.Entities) {
		t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v", expectedActive.Entities, activeEntities)
	}

	// we expect no tokens
	if len(a.currentSegment.tokenCount.CountByNamespaceID) > 0 {
		t.Errorf("expected no token counts to be loaded. got: %v", a.currentSegment.tokenCount.CountByNamespaceID)
	}
}

func TestActivityLog_refreshFromStoredLogNoEntities(t *testing.T) {
	a, _, expectedTokenCounts := setupActivityRecordsInStorage(t, time.Now().UTC(), false, true)
	a.enabled = true

	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(context.Background(), &wg)
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	if !reflect.DeepEqual(a.currentSegment.tokenCount.CountByNamespaceID, expectedTokenCounts) {
		// we expect all token counts to be loaded
		t.Errorf("bad activity token counts loaded. expected: %v got: %v", expectedTokenCounts, a.currentSegment.tokenCount.CountByNamespaceID)
	}

	if len(a.currentSegment.currentEntities.Entities) > 0 {
		t.Errorf("expected no current entity segment to be loaded. got: %v", a.currentSegment.currentEntities)
	}
	activeEntities := a.core.GetActiveEntities()
	if len(activeEntities) > 0 {
		t.Errorf("expected no active entity segment to be loaded. got: %v", activeEntities)
	}
}

// verify current segment refreshed with non-nil empty components and the :expectedStart: timestamp
// note: if :verifyTimeNotZero: is true, ignore :expectedStart: and just make sure the timestamp
// isn't 0
func expectCurrentSegmentRefreshed(t *testing.T, a *ActivityLog, expectedStart int64, verifyTimeNotZero bool) {
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

func TestActivityLog_refreshFromStoredLogNoData(t *testing.T) {
	now := time.Now().UTC()
	a, _, _ := setupActivityRecordsInStorage(t, now, false, false)
	a.enabled = true

	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(context.Background(), &wg)
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	expectCurrentSegmentRefreshed(t, a, now.Unix(), false)
}

func TestActivityLog_refreshFromStoredLogTwoMonthsPrevious(t *testing.T) {
	// test what happens when the most recent data is from month M-2 (or earlier - same effect)
	now := time.Now().UTC()
	twoMonthsAgoStart := timeutil.StartOfPreviousMonth(timeutil.StartOfPreviousMonth(now))
	a, _, _ := setupActivityRecordsInStorage(t, twoMonthsAgoStart, true, true)
	a.enabled = true

	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(context.Background(), &wg)
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	expectCurrentSegmentRefreshed(t, a, now.Unix(), false)
}

func TestActivityLog_refreshFromStoredLogPreviousMonth(t *testing.T) {
	// test what happens when most recent data is from month M-1
	// we expect to load the data from the previous month so that the activeFragmentWorker
	// can handle end of month rotations
	monthStart := timeutil.StartOfMonth(time.Now().UTC())
	oneMonthAgoStart := timeutil.StartOfPreviousMonth(monthStart)
	a, expectedEntityRecords, expectedTokenCounts := setupActivityRecordsInStorage(t, oneMonthAgoStart, true, true)
	a.enabled = true

	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(context.Background(), &wg)
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	expectedActive := &activity.EntityActivityLog{
		Entities: expectedEntityRecords[1:],
	}
	expectedCurrent := &activity.EntityActivityLog{
		Entities: expectedEntityRecords[2:],
	}
	if !entityRecordsEqual(t, a.currentSegment.currentEntities.Entities, expectedCurrent.Entities) {
		// we only expect the newest entity segment to be loaded (for the current month)
		t.Errorf("bad activity entity logs loaded. expected: %v got: %v", expectedCurrent, a.currentSegment.currentEntities)
	}
	if !reflect.DeepEqual(a.currentSegment.tokenCount.CountByNamespaceID, expectedTokenCounts) {
		// we expect all token counts to be loaded
		t.Errorf("bad activity token counts loaded. expected: %v got: %v", expectedTokenCounts, a.currentSegment.tokenCount.CountByNamespaceID)
	}

	activeEntities := a.core.GetActiveEntities()
	if !activeEntitiesEqual(t, activeEntities, expectedActive.Entities) {
		// we expect activeEntities to be loaded for the entire month
		t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v", expectedActive.Entities, activeEntities)
	}
}

func TestActivityLog_refreshFromStoredLogNextMonth(t *testing.T) {
	t.Skip("works on enterprise, fails on oss (oss boots with activity log disabled)")

	// test what happens when most recent data is from month M+1
	nextMonthStart := timeutil.StartOfNextMonth(time.Now().UTC())
	a, _, _ := setupActivityRecordsInStorage(t, nextMonthStart, true, true)
	a.enabled = true

	var wg sync.WaitGroup
	err := a.refreshFromStoredLog(context.Background(), &wg)
	if err != nil {
		t.Fatalf("got error loading stored activity logs: %v", err)
	}
	wg.Wait()

	// we can't know exactly what the timestamp should be set to, just that it shouldn't be zero
	expectCurrentSegmentRefreshed(t, a, time.Now().Unix(), true)
}

func TestActivityLog_IncludeNamespace(t *testing.T) {
	root := namespace.RootNamespace
	a := &ActivityLog{}

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
		writeToStorage(t, core, logPrefix+path, []byte("test"))
	}

	doneCh := make(chan struct{})
	timeout := time.After(20 * time.Second)

	go a.deleteLogWorker(1111, doneCh)
	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatalf("timed out")
	}

	// Check segments still present
	readSegmentFromStorage(t, core, logPrefix+"entity/1112/1")
	readSegmentFromStorage(t, core, logPrefix+"directtokens/1112/1")

	// Check other segments not present
	expectMissingSegment(t, core, logPrefix+"entity/1111/1")
	expectMissingSegment(t, core, logPrefix+"entity/1111/2")
	expectMissingSegment(t, core, logPrefix+"entity/1111/3")
	expectMissingSegment(t, core, logPrefix+"directtokens/1111/1")
}

// Skip this test if too close to the end of a month!
// TODO: move testhelper?
func SkipAtEndOfMonth(t *testing.T) {
	thisMonth := timeutil.StartOfMonth(time.Now().UTC())
	endOfMonth := timeutil.EndOfMonth(thisMonth)
	if endOfMonth.Sub(time.Now()) < 10*time.Minute {
		t.Skip("too close to end of month")
	}
}

func TestActivityLog_EnableDisable(t *testing.T) {
	SkipAtEndOfMonth(t)

	core, b, _ := testCoreSystemBackend(t)
	a := core.activityLog
	view := core.systemBarrierView
	ctx := namespace.RootContext(nil)

	enableRequest := func() {
		t.Helper()
		req := logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
		req.Storage = view
		req.Data["enabled"] = "enable"
		resp, err := b.HandleRequest(ctx, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if resp != nil {
			t.Fatalf("bad: %#v", resp)
		}
	}
	disableRequest := func() {
		t.Helper()
		req := logical.TestRequest(t, logical.UpdateOperation, "internal/counters/config")
		req.Storage = view
		req.Data["enabled"] = "disable"
		resp, err := b.HandleRequest(ctx, req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if resp != nil {
			t.Fatalf("bad: %#v", resp)
		}
	}

	// enable (if not already) and write a segment
	enableRequest()

	id1 := "11111111-1111-1111-1111-111111111111"
	id2 := "22222222-2222-2222-2222-222222222222"
	id3 := "33333333-3333-3333-3333-333333333333"
	a.AddEntityToFragment(id1, "root", time.Now().Unix())
	a.AddEntityToFragment(id2, "root", time.Now().Unix())

	a.currentSegment.startTimestamp -= 10
	seg1 := a.currentSegment.startTimestamp
	err := a.saveCurrentSegmentToStorage(ctx, false)
	if err != nil {
		t.Fatal(err)
	}

	// verify segment exists
	path := fmt.Sprintf("%ventity/%v/0", logPrefix, seg1)
	readSegmentFromStorage(t, core, path)

	// Add in-memory fragment
	a.AddEntityToFragment(id3, "root", time.Now().Unix())

	// disable and verify segment no longer exists
	disableRequest()

	timeout := time.After(20 * time.Second)
	select {
	case <-a.deleteDone:
		break
	case <-timeout:
		t.Fatalf("timed out")
	}

	expectMissingSegment(t, core, path)
	expectCurrentSegmentRefreshed(t, a, 0, false)

	// enable (if not already) which force-writes an empty segment
	enableRequest()

	seg2 := a.currentSegment.startTimestamp
	if seg1 >= seg2 {
		t.Errorf("bad second segment timestamp, %v >= %v", seg1, seg2)
	}

	// Verify empty segments are present
	path = fmt.Sprintf("%ventity/%v/0", logPrefix, seg2)
	readSegmentFromStorage(t, core, path)

	path = fmt.Sprintf("%vdirecttokens/%v/0", logPrefix, seg2)
	readSegmentFromStorage(t, core, path)
}

func TestActivityLog_EndOfMonth(t *testing.T) {
	// We only want *fake* end of months, *real* ones are too scary.
	SkipAtEndOfMonth(t)

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
	segment0 := a.currentSegment.startTimestamp
	month1 := month0.AddDate(0, 1, 0)
	month2 := month0.AddDate(0, 2, 0)

	// Trigger end-of-month
	a.HandleEndOfMonth(month1)

	// Check segment is present, with 1 entity
	path := fmt.Sprintf("%ventity/%v/0", logPrefix, segment0)
	protoSegment := readSegmentFromStorage(t, core, path)
	out := &activity.EntityActivityLog{}
	err := proto.Unmarshal(protoSegment.Value, out)
	if err != nil {
		t.Fatal(err)
	}

	segment1 := a.currentSegment.startTimestamp
	expectedTimestamp := timeutil.StartOfMonth(month1).Unix()
	if segment1 != expectedTimestamp {
		t.Errorf("expected segment timestamp %v got %v", expectedTimestamp, segment1)
	}

	// Check intent log is present
	intentRaw, err := core.barrier.Get(ctx, "sys/counters/activity/endofmonth")
	if err != nil {
		t.Fatal(err)
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

	a.HandleEndOfMonth(month2)
	segment2 := a.currentSegment.startTimestamp

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
		path := fmt.Sprintf("%ventity/%v/0", logPrefix, tc.SegmentTimestamp)
		protoSegment := readSegmentFromStorage(t, core, path)
		out := &activity.EntityActivityLog{}
		err = proto.Unmarshal(protoSegment.Value, out)
		if err != nil {
			t.Fatalf("could not unmarshal protobuf: %v", err)
		}
		expectedEntityIDs(t, out, tc.ExpectedEntityIDs)
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
	startTimestamp := a.currentSegment.startTimestamp

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

	path := logPrefix + "entity/0/0"
	expectMissingSegment(t, core, path)

	path = fmt.Sprintf("%ventity/%v/0", logPrefix, startTimestamp)
	expectMissingSegment(t, core, path)
}

func TestActivityLog_Precompute(t *testing.T) {
	SkipAtEndOfMonth(t)

	january := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	august := time.Date(2020, 8, 15, 12, 0, 0, 0, time.UTC)
	september := timeutil.StartOfMonth(time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC))
	october := timeutil.StartOfMonth(time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC))
	november := timeutil.StartOfMonth(time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC))

	core, _, _ := TestCoreUnsealed(t)
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
			EntityID:    fmt.Sprintf("111122222-3333-4444-5555-%012v", i),
			NamespaceID: entityNamespaces[i/5],
			Timestamp:   time.Now().Unix(),
		}
	}

	toInsert := []struct {
		StartTime int64
		Segment   uint64
		Entities  []*activity.EntityRecord
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
			Entities: segment.Entities,
		}
		data, err := proto.Marshal(eal)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("%ventity/%v/%v", logPrefix, segment.StartTime, segment.Segment)
		writeToStorage(t, core, path, data)
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
		writeToStorage(t, core, "sys/counters/activity/endofmonth", data)

		// Pretend we've successfully rolled over to the following month
		a.l.Lock()
		a.currentSegment.startTimestamp = tc.NextMonth
		a.l.Unlock()

		err = a.precomputedQueryWorker()
		if err != nil {
			t.Fatal(err)
		}

		expectMissingSegment(t, core, "sys/counters/activity/endofmonth")

		for i := 0; i <= tc.ExpectedUpTo; i++ {
			checkPrecomputedQuery(i)
		}

	}
}

type BlockingInmemStorage struct {
}

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

func TestActivityLog_PrecomputeCancel(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog

	// Substitute in a new view
	a.view = NewBarrierView(&BlockingInmemStorage{}, "test")

	core.stopActivityLog()

	done := make(chan struct{})

	// This will block if the shutdown didn't work.
	go func() {
		a.precomputedQueryWorker()
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

func TestActivityLog_NextMonthStart(t *testing.T) {
	SkipAtEndOfMonth(t)

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
		a.l.Lock()
		a.currentSegment.startTimestamp = tc.SegmentStart
		a.l.Unlock()

		actual := a.StartOfNextMonth()
		if !actual.Equal(tc.ExpectedTime) {
			t.Errorf("expected %v, got %v", tc.ExpectedTime, actual)
		}
	}
}

func TestActivityLog_Deletion(t *testing.T) {
	SkipAtEndOfMonth(t)

	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog

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
			entityPath := fmt.Sprintf("%ventity/%v/%v", logPrefix, start.Unix(), j)
			paths[i] = append(paths[i], entityPath)
			writeToStorage(t, core, entityPath, []byte("test"))
		}
		tokenPath := fmt.Sprintf("%vdirecttokens/%v/0", logPrefix, start.Unix())
		paths[i] = append(paths[i], tokenPath)
		writeToStorage(t, core, tokenPath, []byte("test"))

		// No queries for November yet
		if i < novIndex {
			for _, endTime := range times[i+1 : novIndex] {
				queryPath := fmt.Sprintf("sys/counters/activity/queries/%v/%v", start.Unix(), endTime.Unix())
				paths[i] = append(paths[i], queryPath)
				writeToStorage(t, core, queryPath, []byte("test"))
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

	t.Log("24 months")
	now := times[len(times)-1]
	err := a.retentionWorker(now, 24)
	if err != nil {
		t.Fatal(err)
	}
	for i := range times {
		checkPresent(i)
	}

	t.Log("12 months")
	err = a.retentionWorker(now, 12)
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
	err = a.retentionWorker(now, 1)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i <= 19; i++ {
		checkAbsent(i)
	}
	checkPresent(20)
	checkPresent(21)

	t.Log("0 months")
	err = a.retentionWorker(now, 0)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i <= 20; i++ {
		checkAbsent(i)
	}
	checkPresent(21)

}
