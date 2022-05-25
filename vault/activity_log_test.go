package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
	"github.com/mitchellh/mapstructure"
)

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
	a.HandleTokenUsage(context.Background(), te, id, isTWE)

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
	a.HandleTokenUsage(context.Background(), teNew, id, isTWE)

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
	tokenEntryOne := logical.TokenEntry{NamespaceID: "ns1_id", Policies: []string{"hi"}}
	entityEntry := logical.TokenEntry{EntityID: "foo", NamespaceID: "ns1_id", Policies: []string{"hi"}}

	idNonEntity, isTWE := tokenEntryOne.CreateClientID()

	for i := 0; i < 3; i++ {
		a.HandleTokenUsage(ctx, &tokenEntryOne, idNonEntity, isTWE)
	}

	idEntity, isTWE := entityEntry.CreateClientID()
	for i := 0; i < 2; i++ {
		a.HandleTokenUsage(ctx, &entityEntry, idEntity, isTWE)
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
		WriteToStorage(t, core, ActivityLogPrefix+path, []byte("test"))
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

	// enabled check is now inside AddClientToFragment
	a.SetEnable(true)
	a.SetStartTimestamp(time.Now().Unix()) // set a nonzero segment

	// Stop timers for test purposes
	close(a.doneCh)
	defer func() {
		a.l.Lock()
		a.doneCh = make(chan struct{}, 1)
		a.l.Unlock()
	}()

	startTimestamp := a.GetStartTimestamp()
	path0 := fmt.Sprintf("sys/counters/activity/log/entity/%d/0", startTimestamp)
	path1 := fmt.Sprintf("sys/counters/activity/log/entity/%d/1", startTimestamp)
	tokenPath := fmt.Sprintf("sys/counters/activity/log/directtokens/%d/0", startTimestamp)

	genID := func(i int) string {
		return fmt.Sprintf("11111111-1111-1111-1111-%012d", i)
	}
	ts := time.Now().Unix()

	// First 4000 should fit in one segment
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
	if len(entityLog0.Clients) != 4000 {
		t.Fatalf("unexpected entity length. Expected %d, got %d", 4000, len(entityLog0.Clients))
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
	if seqNum != 1 {
		t.Fatalf("expected sequence number 1, got %v", seqNum)
	}

	protoSegment0 = readSegmentFromStorage(t, core, path0)
	err = proto.Unmarshal(protoSegment0.Value, &entityLog0)
	if err != nil {
		t.Fatalf("could not unmarshal protobuf: %v", err)
	}
	if len(entityLog0.Clients) != activitySegmentClientCapacity {
		t.Fatalf("unexpected client length. Expected %d, got %d", activitySegmentClientCapacity,
			len(entityLog0.Clients))
	}

	protoSegment1 := readSegmentFromStorage(t, core, path1)
	entityLog1 := activity.EntityActivityLog{}
	err = proto.Unmarshal(protoSegment1.Value, &entityLog1)
	if err != nil {
		t.Fatalf("could not unmarshal protobuf: %v", err)
	}
	expectedCount := 8100 - activitySegmentClientCapacity
	if len(entityLog1.Clients) != expectedCount {
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

		originalEnabled := core.activityLog.GetEnabled()
		newEnabled := activityLogEnabledDefault

		resp, err := b.HandleRequest(namespace.RootContext(nil), req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		checkAPIWarnings(t, originalEnabled, newEnabled, resp)

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
			t.Fatalf(err.Error())
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

		activeClients := core.GetActiveClients()
		if !ActiveEntitiesEqual(activeClients, tc.entities.Clients) {
			t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v for path %q", tc.entities.Clients, activeClients, tc.path)
		}

		a.resetEntitiesInMemory(t)
	}
}

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
			t.Fatalf(err.Error())
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

		activeClients := core.GetActiveClients()
		if !ActiveEntitiesEqual(activeClients, tc.entities.Clients) {
			t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v for path %q", tc.entities.Clients, activeClients, tc.path)
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
	core.setupActivityLog(ctx, &wg)
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
				t.Fatalf(err.Error())
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
			t.Fatalf(err.Error())
		}

		WriteToStorage(t, core, ActivityLogPrefix+"directtokens/"+fmt.Sprint(base.Unix())+"/0", tokenData)
	}

	return a, entityRecords, tokenRecords
}

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

	activeClients := a.core.GetActiveClients()
	if !ActiveEntitiesEqual(activeClients, expectedActive.Clients) {
		// we expect activeClients to be loaded for the entire month
		t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v", expectedActive.Clients, activeClients)
	}
}

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

	activeClients := a.core.GetActiveClients()
	if !ActiveEntitiesEqual(activeClients, expected.Clients) {
		// we only expect activeClients to be loaded for the newest segment (for the current month)
		t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v", expected.Clients, activeClients)
	}
}

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
	activeClients := a.core.GetActiveClients()
	if !ActiveEntitiesEqual(activeClients, expectedActive.Clients) {
		t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v", expectedActive.Clients, activeClients)
	}

	// we expect no tokens
	nsCount := a.GetStoredTokenCountByNamespaceID()
	if len(nsCount) > 0 {
		t.Errorf("expected no token counts to be loaded. got: %v", nsCount)
	}
}

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
	activeClients := a.core.GetActiveClients()
	if len(activeClients) > 0 {
		t.Errorf("expected no active entity segment to be loaded. got: %v", activeClients)
	}
}

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

	activeClients := a.core.GetActiveClients()
	if !ActiveEntitiesEqual(activeClients, expectedActive.Clients) {
		// we expect activeClients to be loaded for the entire month
		t.Errorf("bad data loaded into active entities. expected only set of EntityID from %v in %v", expectedActive.Clients, activeClients)
	}
}

func TestActivityLog_Export(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)

	january := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	august := time.Date(2020, 8, 15, 12, 0, 0, 0, time.UTC)
	september := timeutil.StartOfMonth(time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC))
	october := timeutil.StartOfMonth(time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC))
	november := timeutil.StartOfMonth(time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC))

	core, _, _, _ := TestCoreUnsealedWithMetrics(t)
	a := core.activityLog
	ctx := namespace.RootContext(nil)

	// Generate overlapping sets of entity IDs from this list.
	//   january:      40-44                                          RRRRR
	//   first month:   0-19  RRRRRAAAAABBBBBRRRRR
	//   second month: 10-29            BBBBBRRRRRRRRRRCCCCC
	//   third month:  15-39                 RRRRRRRRRRCCCCCRRRRRBBBBB

	entityRecords := make([]*activity.EntityRecord, 45)
	entityNamespaces := []string{"root", "aaaaa", "bbbbb", "root", "root", "ccccc", "root", "bbbbb", "rrrrr"}
	authMethods := []string{"auth_1", "auth_2", "auth_3", "auth_4", "auth_5", "auth_6", "auth_7", "auth_8", "auth_9"}

	for i := range entityRecords {
		entityRecords[i] = &activity.EntityRecord{
			ClientID:      fmt.Sprintf("111122222-3333-4444-5555-%012v", i),
			NamespaceID:   entityNamespaces[i/5],
			MountAccessor: authMethods[i/5],
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

	for i, segment := range toInsert {
		eal := &activity.EntityActivityLog{
			Clients: segment.Clients,
		}

		// Mimic a lower time stamp for earlier clients
		for _, c := range eal.Clients {
			c.Timestamp = int64(i)
		}

		data, err := proto.Marshal(eal)
		if err != nil {
			t.Fatal(err)
		}
		path := fmt.Sprintf("%ventity/%v/%v", ActivityLogPrefix, segment.StartTime, segment.Segment)
		WriteToStorage(t, core, path, data)
	}

	tCases := []struct {
		format    string
		startTime time.Time
		endTime   time.Time
		expected  string
	}{
		{
			format:    "json",
			startTime: august,
			endTime:   timeutil.EndOfMonth(september),
			expected:  "aug_sep.json",
		},
		{
			format:    "csv",
			startTime: august,
			endTime:   timeutil.EndOfMonth(september),
			expected:  "aug_sep.csv",
		},
		{
			format:    "json",
			startTime: january,
			endTime:   timeutil.EndOfMonth(november),
			expected:  "full_history.json",
		},
		{
			format:    "csv",
			startTime: january,
			endTime:   timeutil.EndOfMonth(november),
			expected:  "full_history.csv",
		},
		{
			format:    "json",
			startTime: august,
			endTime:   timeutil.EndOfMonth(october),
			expected:  "aug_oct.json",
		},
		{
			format:    "csv",
			startTime: august,
			endTime:   timeutil.EndOfMonth(october),
			expected:  "aug_oct.csv",
		},
		{
			format:    "json",
			startTime: august,
			endTime:   timeutil.EndOfMonth(august),
			expected:  "aug.json",
		},
		{
			format:    "csv",
			startTime: august,
			endTime:   timeutil.EndOfMonth(august),
			expected:  "aug.csv",
		},
	}

	for _, tCase := range tCases {
		rw := &fakeResponseWriter{
			buffer:  &bytes.Buffer{},
			headers: http.Header{},
		}
		if err := a.writeExport(ctx, rw, tCase.format, tCase.startTime, tCase.endTime); err != nil {
			t.Fatal(err)
		}

		expected, err := os.ReadFile(filepath.Join("activity", "test_fixtures", tCase.expected))
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(rw.buffer.Bytes(), expected) {
			t.Fatal(rw.buffer.String())
		}
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
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		checkAPIWarnings(t, originalEnabled, false, resp)
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

	core, _, _, sink := TestCoreUnsealedWithMetrics(t)
	a := core.activityLog
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

		err = a.precomputedQueryWorker(ctx)
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

func TestActivityLog_Precompute(t *testing.T) {
	timeutil.SkipAtEndOfMonth(t)

	january := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	august := time.Date(2020, 8, 15, 12, 0, 0, 0, time.UTC)
	september := timeutil.StartOfMonth(time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC))
	october := timeutil.StartOfMonth(time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC))
	november := timeutil.StartOfMonth(time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC))

	core, _, _, sink := TestCoreUnsealedWithMetrics(t)
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

		err = a.precomputedQueryWorker(ctx)
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

	core, _, _, _ := TestCoreUnsealedWithMetrics(t)
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

		err = a.precomputedQueryWorker(ctx)
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

	core, _, _, sink := TestCoreUnsealedWithMetrics(t)
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

		err = a.precomputedQueryWorker(ctx)
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
		_ = a.precomputedQueryWorker(namespace.RootContext(nil))
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
		if int(clientCounts[clientCount.NamespaceID]) != clientCount.Counts.DistinctEntities {
			t.Errorf("bad entity count for namespace %s . expected %d, got %d", clientCount.NamespaceID, int(clientCounts[clientCount.NamespaceID]), clientCount.Counts.DistinctEntities)
		}
		totalCount := int(clientCounts[clientCount.NamespaceID])
		if totalCount != clientCount.Counts.Clients {
			t.Errorf("bad client count for namespace %s . expected %d, got %d", clientCount.NamespaceID, totalCount, clientCount.Counts.Clients)
		}
	}

	distinctEntities, ok := results["distinct_entities"]
	if !ok {
		t.Fatalf("malformed results. got %v", results)
	}
	if distinctEntities != len(clients) {
		t.Errorf("bad entity count. expected %d, got %d", len(clients), distinctEntities)
	}

	clientCount, ok := results["clients"]
	if !ok {
		t.Fatalf("malformed results. got %v", results)
	}
	if clientCount != len(clients) {
		t.Errorf("bad client count. expected %d, got %d", len(clients), clientCount)
	}
}
