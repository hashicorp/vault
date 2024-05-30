// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/axiomhq/hyperloglog"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/vault/activity"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// Test_ActivityLog_ComputeCurrentMonthForBillingPeriodInternal creates 3 months
// of hyperloglogs and fills them with overlapping clients. The test calls
// computeCurrentMonthForBillingPeriodInternal with the current month map having
// some overlap with the previous months. The test then verifies that the
// results have the correct number of entity, non-entity, and secret sync
// association clients. The test also calls
// computeCurrentMonthForBillingPeriodInternal with an empty current month map,
// and verifies that the results are all 0.
func Test_ActivityLog_ComputeCurrentMonthForBillingPeriodInternal(t *testing.T) {
	// populate the first month with clients 1-20
	monthOneHLL := hyperloglog.New()
	// populate the second month with clients 10-30
	monthTwoHLL := hyperloglog.New()
	// populate the third month with clients 20-40
	monthThreeHLL := hyperloglog.New()

	for i := 0; i < 40; i++ {
		clientID := []byte(fmt.Sprintf("client_%d", i))
		if i < 20 {
			monthOneHLL.Insert(clientID)
		}
		if 10 <= i && i < 20 {
			monthTwoHLL.Insert(clientID)
		}
		if 20 <= i && i < 40 {
			monthThreeHLL.Insert(clientID)
		}
	}
	mockHLLGetFunc := func(ctx context.Context, startTime time.Time) (*hyperloglog.Sketch, error) {
		currMonthStart := timeutil.StartOfMonth(time.Now())
		if startTime.Equal(timeutil.MonthsPreviousTo(3, currMonthStart)) {
			return monthThreeHLL, nil
		}
		if startTime.Equal(timeutil.MonthsPreviousTo(2, currMonthStart)) {
			return monthTwoHLL, nil
		}
		if startTime.Equal(timeutil.MonthsPreviousTo(1, currMonthStart)) {
			return monthOneHLL, nil
		}
		return nil, fmt.Errorf("bad start time")
	}

	// Below we register the entity, non-entity, and secret sync clients that
	// are seen in the current month

	// Let's add 2 entities exclusive to month 1 (clients 0,1),
	// 2 entities shared by month 1 and 2 (clients 10,11),
	// 2 entities shared by month 2 and 3 (clients 20,21), and
	// 2 entities exclusive to month 3 (30,31). Furthermore, we can add
	// 3 new entities (clients 40,41,42).
	entitiesStruct := map[string]struct{}{
		"client_0":  {},
		"client_1":  {},
		"client_10": {},
		"client_11": {},
		"client_20": {},
		"client_21": {},
		"client_30": {},
		"client_31": {},
		"client_40": {},
		"client_41": {},
		"client_42": {},
	}

	// We will add 3 nonentity clients from month 1 (clients 2,3,4),
	// 3 shared by months 1 and 2 (12,13,14),
	// 3 shared by months 2 and 3 (22,23,24), and
	// 3 exclusive to month 3 (32,33,34). We will also
	// add 4 new nonentity clients (43,44,45,46)
	nonEntitiesStruct := map[string]struct{}{
		"client_2":  {},
		"client_3":  {},
		"client_4":  {},
		"client_12": {},
		"client_13": {},
		"client_14": {},
		"client_22": {},
		"client_23": {},
		"client_24": {},
		"client_32": {},
		"client_33": {},
		"client_34": {},
		"client_43": {},
		"client_44": {},
		"client_45": {},
		"client_46": {},
	}

	// secret syncs have 1 client from month 1 (5)
	// 1 shared by months 1 and 2 (15)
	// 1 shared by months 2 and 3 (25)
	// 2 exclusive to month 3 (35,36)
	// and 2 new clients (47,48)
	secretSyncStruct := map[string]struct{}{
		"client_5":  {},
		"client_15": {},
		"client_25": {},
		"client_35": {},
		"client_36": {},
		"client_47": {},
		"client_48": {},
	}

	counts := &processCounts{
		ClientsByType: map[string]clientIDSet{
			entityActivityType:         entitiesStruct,
			nonEntityTokenActivityType: nonEntitiesStruct,
			secretSyncActivityType:     secretSyncStruct,
		},
	}

	currentMonthClientsMap := make(map[int64]*processMonth, 1)
	currentMonthClients := &processMonth{
		Counts:     counts,
		NewClients: &processNewClients{Counts: counts},
	}
	// Technially I think currentMonthClientsMap should have the keys as
	// unix timestamps, but for the purposes of the unit test it doesn't
	// matter what the values actually are.
	currentMonthClientsMap[0] = currentMonthClients

	core, _, _ := TestCoreUnsealed(t)
	a := core.activityLog

	endTime := timeutil.StartOfMonth(time.Now())
	startTime := timeutil.MonthsPreviousTo(3, endTime)

	monthRecord, err := a.computeCurrentMonthForBillingPeriodInternal(context.Background(), currentMonthClientsMap, mockHLLGetFunc, startTime, endTime)
	require.NoError(t, err)

	require.Equal(t, &activity.CountsRecord{
		EntityClients:    11,
		NonEntityClients: 16,
		SecretSyncs:      7,
	}, monthRecord.Counts)

	require.Equal(t, &activity.CountsRecord{
		EntityClients:    3,
		NonEntityClients: 4,
		SecretSyncs:      2,
	}, monthRecord.NewClients.Counts)

	// Attempt to compute current month when no records exist
	endTime = time.Now().UTC()
	startTime = timeutil.StartOfMonth(endTime)
	emptyClientsMap := make(map[int64]*processMonth, 0)
	monthRecord, err = a.computeCurrentMonthForBillingPeriodInternal(context.Background(), emptyClientsMap, mockHLLGetFunc, startTime, endTime)
	require.NoError(t, err)

	require.Equal(t, &activity.CountsRecord{}, monthRecord.Counts)
	require.Equal(t, &activity.CountsRecord{}, monthRecord.NewClients.Counts)
}

// writeEntitySegment writes a single segment file with the given time and index for an entity
func writeEntitySegment(t *testing.T, core *Core, ts time.Time, index int, item *activity.EntityActivityLog) {
	t.Helper()
	protoItem, err := proto.Marshal(item)
	require.NoError(t, err)
	WriteToStorage(t, core, makeSegmentPath(t, activityEntityBasePath, ts, index), protoItem)
}

// writeTokenSegment writes a single segment file with the given time and index for a token
func writeTokenSegment(t *testing.T, core *Core, ts time.Time, index int, item *activity.TokenCount) {
	t.Helper()
	protoItem, err := proto.Marshal(item)
	require.NoError(t, err)
	WriteToStorage(t, core, makeSegmentPath(t, activityTokenBasePath, ts, index), protoItem)
}

// makeSegmentPath formats the path for a segment at a particular time and index
func makeSegmentPath(t *testing.T, typ string, ts time.Time, index int) string {
	t.Helper()
	return fmt.Sprintf("%s%s%d/%d", ActivityPrefix, typ, ts.Unix(), index)
}

// TestSegmentFileReader_BadData verifies that the reader returns errors when the data is unable to be parsed
// However, the next time that Read*() is called, the reader should still progress and be able to then return any
// valid data without errors
func TestSegmentFileReader_BadData(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	now := time.Now()

	// write bad data that won't be able to be unmarshaled at index 0
	WriteToStorage(t, core, makeSegmentPath(t, activityTokenBasePath, now, 0), []byte("fake data"))
	WriteToStorage(t, core, makeSegmentPath(t, activityEntityBasePath, now, 0), []byte("fake data"))

	// write entity at index 1
	entity := &activity.EntityActivityLog{Clients: []*activity.EntityRecord{
		{
			ClientID: "id",
		},
	}}
	writeEntitySegment(t, core, now, 1, entity)

	// write token at index 1
	token := &activity.TokenCount{CountByNamespaceID: map[string]uint64{
		"ns": 1,
	}}
	writeTokenSegment(t, core, now, 1, token)
	reader, err := core.activityLog.NewSegmentFileReader(context.Background(), now)
	require.NoError(t, err)

	// first the bad entity is read, which returns an error
	_, err = reader.ReadEntity(context.Background())
	require.Error(t, err)
	// then, the reader can read the good entity at index 1
	gotEntity, err := reader.ReadEntity(context.Background())
	require.True(t, proto.Equal(gotEntity, entity))
	require.Nil(t, err)

	// the bad token causes an error
	_, err = reader.ReadToken(context.Background())
	require.Error(t, err)
	// but the good token is able to be read
	gotToken, err := reader.ReadToken(context.Background())
	require.True(t, proto.Equal(gotToken, token))
	require.Nil(t, err)
}

// TestSegmentFileReader_MissingData verifies that the segment file reader will skip over missing segment paths without
// errorring until it is able to find a valid segment path
func TestSegmentFileReader_MissingData(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	now := time.Now()
	// write entities and tokens at indexes 0, 1, 2
	for i := 0; i < 3; i++ {
		WriteToStorage(t, core, makeSegmentPath(t, activityTokenBasePath, now, i), []byte("fake data"))
		WriteToStorage(t, core, makeSegmentPath(t, activityEntityBasePath, now, i), []byte("fake data"))

	}
	// write entity at index 3
	entity := &activity.EntityActivityLog{Clients: []*activity.EntityRecord{
		{
			ClientID: "id",
		},
	}}
	writeEntitySegment(t, core, now, 3, entity)
	// write token at index 3
	token := &activity.TokenCount{CountByNamespaceID: map[string]uint64{
		"ns": 1,
	}}
	writeTokenSegment(t, core, now, 3, token)
	reader, err := core.activityLog.NewSegmentFileReader(context.Background(), now)
	require.NoError(t, err)

	// delete the indexes 0, 1, 2
	for i := 0; i < 3; i++ {
		require.NoError(t, core.barrier.Delete(context.Background(), makeSegmentPath(t, activityTokenBasePath, now, i)))
		require.NoError(t, core.barrier.Delete(context.Background(), makeSegmentPath(t, activityEntityBasePath, now, i)))
	}

	// we expect the reader to only return the data at index 3, and then be done
	gotEntity, err := reader.ReadEntity(context.Background())
	require.NoError(t, err)
	require.True(t, proto.Equal(gotEntity, entity))
	_, err = reader.ReadEntity(context.Background())
	require.Equal(t, err, io.EOF)

	gotToken, err := reader.ReadToken(context.Background())
	require.NoError(t, err)
	require.True(t, proto.Equal(gotToken, token))
	_, err = reader.ReadToken(context.Background())
	require.Equal(t, err, io.EOF)
}

// TestSegmentFileReader_NoData verifies that the reader return io.EOF when there is no data
func TestSegmentFileReader_NoData(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	now := time.Now()
	reader, err := core.activityLog.NewSegmentFileReader(context.Background(), now)
	require.NoError(t, err)
	entity, err := reader.ReadEntity(context.Background())
	require.Nil(t, entity)
	require.Equal(t, err, io.EOF)
	token, err := reader.ReadToken(context.Background())
	require.Nil(t, token)
	require.Equal(t, err, io.EOF)
}

// TestSegmentFileReader verifies that the reader iterates through all segments paths in ascending order and returns
// io.EOF when it's done
func TestSegmentFileReader(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	now := time.Now()
	entities := make([]*activity.EntityActivityLog, 0, 3)
	tokens := make([]*activity.TokenCount, 0, 3)

	// write 3 entity segment pieces and 3 token segment pieces
	for i := 0; i < 3; i++ {
		entity := &activity.EntityActivityLog{Clients: []*activity.EntityRecord{
			{
				ClientID: fmt.Sprintf("id-%d", i),
			},
		}}
		token := &activity.TokenCount{CountByNamespaceID: map[string]uint64{
			fmt.Sprintf("ns-%d", i): uint64(i),
		}}
		writeEntitySegment(t, core, now, i, entity)
		writeTokenSegment(t, core, now, i, token)
		entities = append(entities, entity)
		tokens = append(tokens, token)
	}

	reader, err := core.activityLog.NewSegmentFileReader(context.Background(), now)
	require.NoError(t, err)

	gotEntities := make([]*activity.EntityActivityLog, 0, 3)
	gotTokens := make([]*activity.TokenCount, 0, 3)

	// read the entities from the reader
	for entity, err := reader.ReadEntity(context.Background()); !errors.Is(err, io.EOF); entity, err = reader.ReadEntity(context.Background()) {
		require.NoError(t, err)
		gotEntities = append(gotEntities, entity)
	}

	// read the tokens from the reader
	for token, err := reader.ReadToken(context.Background()); !errors.Is(err, io.EOF); token, err = reader.ReadToken(context.Background()) {
		require.NoError(t, err)
		gotTokens = append(gotTokens, token)
	}
	require.Len(t, gotEntities, 3)
	require.Len(t, gotTokens, 3)

	// verify that the entities and tokens we got from the reader are correct
	// we can't use require.Equals() here because there are protobuf differences in unexported fields
	for i := 0; i < 3; i++ {
		require.True(t, proto.Equal(gotEntities[i], entities[i]))
		require.True(t, proto.Equal(gotTokens[i], tokens[i]))
	}
}
