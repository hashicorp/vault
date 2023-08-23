// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package vault

import (
	"context"
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// TestSystemBackend_handleActivityWriteData calls the activity log write endpoint and confirms that the inputs are
// correctly validated
func TestSystemBackend_handleActivityWriteData(t *testing.T) {
	testCases := []struct {
		name      string
		operation logical.Operation
		input     map[string]interface{}
		wantError error
		wantPaths int
	}{
		{
			name:      "read fails",
			operation: logical.ReadOperation,
			wantError: logical.ErrUnsupportedOperation,
		},
		{
			name:      "empty write fails",
			operation: logical.UpdateOperation,
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "wrong key fails",
			operation: logical.UpdateOperation,
			input:     map[string]interface{}{"other": "data"},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "incorrectly formatted data fails",
			operation: logical.UpdateOperation,
			input:     map[string]interface{}{"input": "data"},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "incorrect json data fails",
			operation: logical.UpdateOperation,
			input:     map[string]interface{}{"input": `{"other":"json"}`},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "empty write value fails",
			operation: logical.UpdateOperation,
			input:     map[string]interface{}{"input": `{"write":[],"data":[]}`},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "empty data value fails",
			operation: logical.UpdateOperation,
			input:     map[string]interface{}{"input": `{"write":["WRITE_PRECOMPUTED_QUERIES"],"data":[]}`},
			wantError: logical.ErrInvalidRequest,
		},
		{
			name:      "correctly formatted data succeeds",
			operation: logical.UpdateOperation,
			input:     map[string]interface{}{"input": `{"write":["WRITE_PRECOMPUTED_QUERIES"],"data":[{"current_month":true,"all":{"clients":[{"count":5}]}}]}`},
		},
		{
			name:      "entities with multiple segments",
			operation: logical.UpdateOperation,
			input:     map[string]interface{}{"input": `{"write":["WRITE_ENTITIES"],"data":[{"current_month":true,"num_segments":3,"all":{"clients":[{"count":5}]}}]}`},
			wantPaths: 3,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := testSystemBackend(t)
			req := logical.TestRequest(t, tc.operation, "internal/counters/activity/write")
			req.Data = tc.input
			resp, err := b.HandleRequest(namespace.RootContext(nil), req)
			if tc.wantError != nil {
				require.Equal(t, tc.wantError, err, resp.Error())
			} else {
				require.NoError(t, err)
				paths := resp.Data["paths"].([]string)
				require.Len(t, paths, tc.wantPaths)
			}
		})
	}
}

// Test_singleMonthActivityClients_addNewClients verifies that new clients are
// created correctly, adhering to the requested parameters. The clients should
// use the inputted mount and a generated ID if one is not supplied. The new
// client should be added to the month's `clients` slice and segment map, if
// a segment index is supplied
func Test_singleMonthActivityClients_addNewClients(t *testing.T) {
	segmentIndex := 0
	tests := []struct {
		name          string
		mount         string
		clients       *generation.Client
		wantNamespace string
		wantMount     string
		wantID        string
		segmentIndex  *int
	}{
		{
			name:      "default mount is used",
			mount:     "default_mount",
			wantMount: "default_mount",
			clients:   &generation.Client{},
		},
		{
			name:          "record namespace is used, default mount is used",
			mount:         "default_mount",
			wantNamespace: "ns",
			wantMount:     "default_mount",
			clients: &generation.Client{
				Namespace: "ns",
				Mount:     "mount",
			},
		},
		{
			name: "predefined ID is used",
			clients: &generation.Client{
				Id: "client_id",
			},
			wantID: "client_id",
		},
		{
			name: "non zero count",
			clients: &generation.Client{
				Count: 5,
			},
		},
		{
			name: "non entity client",
			clients: &generation.Client{
				ClientType: "non-entity",
			},
		},
		{
			name: "acme client",
			clients: &generation.Client{
				ClientType: "acme",
			},
		},
		{
			name:         "added to segment",
			clients:      &generation.Client{},
			segmentIndex: &segmentIndex,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &singleMonthActivityClients{
				predefinedSegments: make(map[int][]int),
			}
			err := m.addNewClients(tt.clients, tt.mount, tt.segmentIndex)
			require.NoError(t, err)
			numNew := tt.clients.Count
			if numNew == 0 {
				numNew = 1
			}
			require.Len(t, m.clients, int(numNew))
			for i, rec := range m.clients {
				require.NotNil(t, rec)
				require.Equal(t, tt.wantNamespace, rec.NamespaceID)
				require.Equal(t, tt.wantMount, rec.MountAccessor)
				require.Equal(t, tt.clients.ClientType, rec.ClientType)
				if tt.wantID != "" {
					require.Equal(t, tt.wantID, rec.ClientID)
				} else {
					require.NotEqual(t, "", rec.ClientID)
				}
				if tt.segmentIndex != nil {
					require.Contains(t, m.predefinedSegments[*tt.segmentIndex], i)
				}
			}
		})
	}
}

// Test_multipleMonthsActivityClients_processMonth verifies that a month of data
// is added correctly. The test checks that default values are handled correctly
// for mounts and namespaces.
func Test_multipleMonthsActivityClients_processMonth(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	tests := []struct {
		name      string
		clients   *generation.Data
		wantError bool
		numMonths int
	}{
		{
			name: "specified namespace and mount exist",
			clients: &generation.Data{
				Clients: &generation.Data_All{All: &generation.Clients{Clients: []*generation.Client{{
					Namespace: namespace.RootNamespaceID,
					Mount:     "identity/",
				}}}},
			},
			numMonths: 1,
		},
		{
			name: "specified namespace exists, mount empty",
			clients: &generation.Data{
				Clients: &generation.Data_All{All: &generation.Clients{Clients: []*generation.Client{{
					Namespace: namespace.RootNamespaceID,
				}}}},
			},
			numMonths: 1,
		},
		{
			name: "empty namespace and mount",
			clients: &generation.Data{
				Clients: &generation.Data_All{All: &generation.Clients{Clients: []*generation.Client{{}}}},
			},
			numMonths: 1,
		},
		{
			name: "namespace doesn't exist",
			clients: &generation.Data{
				Clients: &generation.Data_All{All: &generation.Clients{Clients: []*generation.Client{{
					Namespace: "abcd",
				}}}},
			},
			wantError: true,
			numMonths: 1,
		},
		{
			name: "namespace exists, mount doesn't exist",
			clients: &generation.Data{
				Clients: &generation.Data_All{All: &generation.Clients{Clients: []*generation.Client{{
					Namespace: namespace.RootNamespaceID,
					Mount:     "mount",
				}}}},
			},
			wantError: true,
			numMonths: 1,
		},
		{
			name: "older month",
			clients: &generation.Data{
				Month:   &generation.Data_MonthsAgo{MonthsAgo: 4},
				Clients: &generation.Data_All{All: &generation.Clients{Clients: []*generation.Client{{}}}},
			},
			numMonths: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMultipleMonthsActivityClients(tt.numMonths)
			err := m.processMonth(context.Background(), core, tt.clients)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, m.months[tt.clients.GetMonthsAgo()].clients, len(tt.clients.GetAll().Clients))
				for _, month := range m.months {
					for _, c := range month.clients {
						require.NotEmpty(t, c.NamespaceID)
						require.NotEmpty(t, c.MountAccessor)
					}
				}
			}
		})
	}
}

// Test_multipleMonthsActivityClients_processMonth_segmented verifies that segments
// are filled correctly when a month is processed with segmented data. The clients
// should be in the clients array, and should also be in the predefinedSegments map
// at the correct segment index
func Test_multipleMonthsActivityClients_processMonth_segmented(t *testing.T) {
	index7 := int32(7)
	data := &generation.Data{
		Clients: &generation.Data_Segments{
			Segments: &generation.Segments{
				Segments: []*generation.Segment{
					{
						Clients: &generation.Clients{Clients: []*generation.Client{
							{},
						}},
					},
					{
						Clients: &generation.Clients{Clients: []*generation.Client{{}}},
					},
					{
						SegmentIndex: &index7,
						Clients:      &generation.Clients{Clients: []*generation.Client{{}}},
					},
				},
			},
		},
	}
	m := newMultipleMonthsActivityClients(1)
	core, _, _ := TestCoreUnsealed(t)
	require.NoError(t, m.processMonth(context.Background(), core, data))
	require.Len(t, m.months[0].predefinedSegments, 3)
	require.Len(t, m.months[0].clients, 3)

	// segment indexes are correct
	require.Contains(t, m.months[0].predefinedSegments, 0)
	require.Contains(t, m.months[0].predefinedSegments, 1)
	require.Contains(t, m.months[0].predefinedSegments, 7)

	// the data in each segment is correct
	require.Contains(t, m.months[0].predefinedSegments[0], 0)
	require.Contains(t, m.months[0].predefinedSegments[1], 1)
	require.Contains(t, m.months[0].predefinedSegments[7], 2)
}

// Test_multipleMonthsActivityClients_addRepeatedClients adds repeated clients
// from 1 month ago and 2 months ago, and verifies that the correct clients are
// added based on namespace, mount, and non-entity attributes
func Test_multipleMonthsActivityClients_addRepeatedClients(t *testing.T) {
	m := newMultipleMonthsActivityClients(3)
	defaultMount := "default"

	require.NoError(t, m.addClientToMonth(2, &generation.Client{Count: 2}, "identity", nil))
	require.NoError(t, m.addClientToMonth(2, &generation.Client{Count: 2, Namespace: "other_ns"}, defaultMount, nil))
	require.NoError(t, m.addClientToMonth(1, &generation.Client{Count: 2}, defaultMount, nil))
	require.NoError(t, m.addClientToMonth(1, &generation.Client{Count: 2, ClientType: "non-entity"}, defaultMount, nil))

	month2Clients := m.months[2].clients
	month1Clients := m.months[1].clients

	thisMonth := m.months[0]
	// this will match the first client in month 1
	require.NoError(t, m.addRepeatedClients(0, &generation.Client{Count: 1, Repeated: true}, defaultMount, nil))
	require.Contains(t, month1Clients, thisMonth.clients[0])

	// this will match the 3rd client in month 1
	require.NoError(t, m.addRepeatedClients(0, &generation.Client{Count: 1, Repeated: true, ClientType: "non-entity"}, defaultMount, nil))
	require.Equal(t, month1Clients[2], thisMonth.clients[1])

	// this will match the first two clients in month 1
	require.NoError(t, m.addRepeatedClients(0, &generation.Client{Count: 2, Repeated: true}, defaultMount, nil))
	require.Equal(t, month1Clients[0:2], thisMonth.clients[2:4])

	// this will match the first client in month 2
	require.NoError(t, m.addRepeatedClients(0, &generation.Client{Count: 1, RepeatedFromMonth: 2}, "identity", nil))
	require.Equal(t, month2Clients[0], thisMonth.clients[4])

	// this will match the 3rd client in month 2
	require.NoError(t, m.addRepeatedClients(0, &generation.Client{Count: 1, RepeatedFromMonth: 2, Namespace: "other_ns"}, defaultMount, nil))
	require.Equal(t, month2Clients[2], thisMonth.clients[5])

	require.Error(t, m.addRepeatedClients(0, &generation.Client{Count: 1, RepeatedFromMonth: 2, Namespace: "other_ns"}, "other_mount", nil))
}

// Test_singleMonthActivityClients_populateSegments calls populateSegments for a
// collection of 5 clients, segmented in various ways. The test ensures that the
// resulting map has the correct clients for each segment index
func Test_singleMonthActivityClients_populateSegments(t *testing.T) {
	clients := []*activity.EntityRecord{
		{ClientID: "a"},
		{ClientID: "b"},
		{ClientID: "c"},
		{ClientID: "d"},
		{ClientID: "e"},
	}
	cases := []struct {
		name         string
		segments     map[int][]int
		numSegments  int
		emptyIndexes []int32
		skipIndexes  []int32
		wantSegments map[int][]*activity.EntityRecord
	}{
		{
			name: "segmented",
			segments: map[int][]int{
				0: {0, 1},
				1: {2, 3},
				2: {4},
			},
			wantSegments: map[int][]*activity.EntityRecord{
				0: {{ClientID: "a"}, {ClientID: "b"}},
				1: {{ClientID: "c"}, {ClientID: "d"}},
				2: {{ClientID: "e"}},
			},
		},
		{
			name: "segmented with skip and empty",
			segments: map[int][]int{
				0: {0, 1},
				2: {0, 1},
			},
			emptyIndexes: []int32{1, 4},
			skipIndexes:  []int32{3},
			wantSegments: map[int][]*activity.EntityRecord{
				0: {{ClientID: "a"}, {ClientID: "b"}},
				1: {},
				2: {{ClientID: "a"}, {ClientID: "b"}},
				3: nil,
				4: {},
			},
		},
		{
			name:        "all clients",
			numSegments: 0,
			wantSegments: map[int][]*activity.EntityRecord{
				0: {{ClientID: "a"}, {ClientID: "b"}, {ClientID: "c"}, {ClientID: "d"}, {ClientID: "e"}},
			},
		},
		{
			name:        "all clients split",
			numSegments: 2,
			wantSegments: map[int][]*activity.EntityRecord{
				0: {{ClientID: "a"}, {ClientID: "b"}, {ClientID: "c"}},
				1: {{ClientID: "d"}, {ClientID: "e"}},
			},
		},
		{
			name:         "all clients with skip and empty",
			numSegments:  5,
			skipIndexes:  []int32{0, 3},
			emptyIndexes: []int32{2},
			wantSegments: map[int][]*activity.EntityRecord{
				0: nil,
				1: {{ClientID: "a"}, {ClientID: "b"}, {ClientID: "c"}},
				2: {},
				3: nil,
				4: {{ClientID: "d"}, {ClientID: "e"}},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := singleMonthActivityClients{predefinedSegments: tc.segments, clients: clients, generationParameters: &generation.Data{EmptySegmentIndexes: tc.emptyIndexes, SkipSegmentIndexes: tc.skipIndexes, NumSegments: int32(tc.numSegments)}}
			gotSegments, err := s.populateSegments()
			require.NoError(t, err)
			require.Equal(t, tc.wantSegments, gotSegments)
		})
	}
}

// Test_handleActivityWriteData writes 4 months of data splitting some months
// across segments and using empty segments and skipped segments. Entities and
// precomputed queries are written. written and then storage is queried. The
// test verifies that the correct timestamps are present in the activity log and
// that the correct segment numbers for each month contain the correct number of
// clients
func Test_handleActivityWriteData(t *testing.T) {
	index5 := int32(5)
	index4 := int32(4)
	data := []*generation.Data{
		{
			// segments: 0:[x,y], 1:[z]
			Month:       &generation.Data_MonthsAgo{MonthsAgo: 3},
			Clients:     &generation.Data_All{All: &generation.Clients{Clients: []*generation.Client{{Count: 3}}}},
			NumSegments: 2,
		},
		{
			// segments: 1:[a,b,c], 2:[d,e]
			Month:              &generation.Data_MonthsAgo{MonthsAgo: 2},
			Clients:            &generation.Data_All{All: &generation.Clients{Clients: []*generation.Client{{Count: 5}}}},
			NumSegments:        3,
			SkipSegmentIndexes: []int32{0},
		},
		{
			// segments: 5:[f,g]
			Month: &generation.Data_MonthsAgo{MonthsAgo: 1},
			Clients: &generation.Data_Segments{
				Segments: &generation.Segments{Segments: []*generation.Segment{{
					SegmentIndex: &index5,
					Clients:      &generation.Clients{Clients: []*generation.Client{{Count: 2}}},
				}}},
			},
		},
		{
			// segments: 1:[], 2:[], 4:[n], 5:[o]
			Month:               &generation.Data_CurrentMonth{},
			EmptySegmentIndexes: []int32{1, 2},
			Clients: &generation.Data_Segments{
				Segments: &generation.Segments{Segments: []*generation.Segment{
					{
						SegmentIndex: &index5,
						Clients:      &generation.Clients{Clients: []*generation.Client{{Count: 1}}},
					},
					{
						SegmentIndex: &index4,
						Clients:      &generation.Clients{Clients: []*generation.Client{{Count: 1}}},
					},
				}},
			},
		},
	}

	t.Run("write entitites", func(t *testing.T) {
		core, _, _ := TestCoreUnsealed(t)
		marshaled, err := protojson.Marshal(&generation.ActivityLogMockInput{
			Data:  data,
			Write: []generation.WriteOptions{generation.WriteOptions_WRITE_ENTITIES},
		})
		require.NoError(t, err)
		req := logical.TestRequest(t, logical.UpdateOperation, "internal/counters/activity/write")
		req.Data = map[string]interface{}{"input": string(marshaled)}
		resp, err := core.systemBackend.HandleRequest(namespace.RootContext(nil), req)
		require.NoError(t, err)
		paths := resp.Data["paths"].([]string)
		require.Len(t, paths, 9)

		times, err := core.activityLog.availableLogs(context.Background())
		require.NoError(t, err)
		require.Len(t, times, 4)

		sortPaths := func(monthPaths []string) {
			sort.Slice(monthPaths, func(i, j int) bool {
				iVal, _ := parseSegmentNumberFromPath(monthPaths[i])
				jVal, _ := parseSegmentNumberFromPath(monthPaths[j])
				return iVal < jVal
			})
		}

		month0Paths := paths[0:4]
		month1Paths := paths[4:5]
		month2Paths := paths[5:7]
		month3Paths := paths[7:9]
		sortPaths(month0Paths)
		sortPaths(month1Paths)
		sortPaths(month2Paths)
		sortPaths(month3Paths)
		entities := func(paths []string) map[int][]*activity.EntityRecord {
			segments := make(map[int][]*activity.EntityRecord)
			for _, path := range paths {
				segmentNum, _ := parseSegmentNumberFromPath(path)
				entry, err := core.activityLog.view.Get(context.Background(), path)
				require.NoError(t, err)
				if entry == nil {
					segments[segmentNum] = []*activity.EntityRecord{}
					continue
				}
				activities := &activity.EntityActivityLog{}
				err = proto.Unmarshal(entry.Value, activities)
				require.NoError(t, err)
				segments[segmentNum] = activities.Clients
			}
			return segments
		}
		month0Entities := entities(month0Paths)
		require.Len(t, month0Entities, 4)
		require.Contains(t, month0Entities, 1)
		require.Contains(t, month0Entities, 2)
		require.Contains(t, month0Entities, 4)
		require.Contains(t, month0Entities, 5)
		require.Len(t, month0Entities[1], 0)
		require.Len(t, month0Entities[2], 0)
		require.Len(t, month0Entities[4], 1)
		require.Len(t, month0Entities[5], 1)

		month1Entities := entities(month1Paths)
		require.Len(t, month1Entities, 1)
		require.Contains(t, month1Entities, 5)
		require.Len(t, month1Entities[5], 2)

		month2Entities := entities(month2Paths)
		require.Len(t, month2Entities, 2)
		require.Contains(t, month2Entities, 1)
		require.Contains(t, month2Entities, 2)
		require.Len(t, month2Entities[1], 3)
		require.Len(t, month2Entities[2], 2)

		month3Entities := entities(month3Paths)
		require.Len(t, month3Entities, 2)
		require.Contains(t, month3Entities, 0)
		require.Contains(t, month3Entities, 1)
		require.Len(t, month3Entities[0], 2)
		require.Len(t, month3Entities[1], 1)
	})
	t.Run("write precomputed queries", func(t *testing.T) {
		core, _, _ := TestCoreUnsealed(t)
		marshaled, err := protojson.Marshal(&generation.ActivityLogMockInput{
			Data:  data,
			Write: []generation.WriteOptions{generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES},
		})
		require.NoError(t, err)
		req := logical.TestRequest(t, logical.UpdateOperation, "internal/counters/activity/write")
		req.Data = map[string]interface{}{"input": string(marshaled)}
		_, err = core.systemBackend.HandleRequest(namespace.RootContext(nil), req)
		require.NoError(t, err)

		queries, err := core.activityLog.queryStore.QueriesAvailable(context.Background())
		require.NoError(t, err)
		require.True(t, queries)

		now := time.Now().UTC()
		start := timeutil.StartOfMonth(timeutil.MonthsPreviousTo(3, now))
		end := timeutil.EndOfMonth(now)
		pq, err := core.activityLog.queryStore.Get(context.Background(), start, end)
		require.NoError(t, err)
		require.NotNil(t, pq)
		require.Equal(t, end, pq.EndTime)
		require.Equal(t, start, pq.StartTime)
		require.Len(t, pq.Namespaces, 1)
		require.Equal(t, uint64(12), pq.Namespaces[0].Entities)
		require.Len(t, pq.Months, 4)
	})
	t.Run("write intent logs", func(t *testing.T) {
		core, _, _ := TestCoreUnsealed(t)
		marshaled, err := protojson.Marshal(&generation.ActivityLogMockInput{
			Data:  data,
			Write: []generation.WriteOptions{generation.WriteOptions_WRITE_INTENT_LOGS},
		})
		require.NoError(t, err)
		now := time.Now().UTC()
		req := logical.TestRequest(t, logical.UpdateOperation, "internal/counters/activity/write")
		req.Data = map[string]interface{}{"input": string(marshaled)}
		_, err = core.systemBackend.HandleRequest(namespace.RootContext(nil), req)
		require.NoError(t, err)

		entry, err := core.activityLog.view.Get(context.Background(), activityIntentLogKey)
		require.NoError(t, err)
		var intent ActivityIntentLog
		err = json.Unmarshal(entry.Value, &intent)
		require.NoError(t, err)
		prev := time.Unix(intent.PreviousMonth, 0)
		next := time.Unix(intent.NextMonth, 0)

		require.Equal(t, timeutil.StartOfMonth(now), next.UTC())
		require.Equal(t, timeutil.StartOfMonth(timeutil.MonthsPreviousTo(3, now)), prev.UTC())

		times, err := core.activityLog.availableLogs(context.Background())
		require.NoError(t, err)
		require.Len(t, times, 4)
	})
}
