// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clientcountutil

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
	"github.com/stretchr/testify/require"
)

// TestNewCurrentMonthData verifies that current month is set correctly and that
// there are no open segments
func TestNewCurrentMonthData(t *testing.T) {
	generator := NewActivityLogData(nil).NewCurrentMonthData()
	require.True(t, generator.data.Data[0].GetCurrentMonth())
	require.True(t, generator.addingToMonth.GetCurrentMonth())
	require.Nil(t, generator.addingToSegment)
}

// TestNewMonthDataMonthsAgo verifies that months ago is set correctly and that
// there are no open segments
func TestNewMonthDataMonthsAgo(t *testing.T) {
	generator := NewActivityLogData(nil).NewPreviousMonthData(3)
	require.Equal(t, int32(3), generator.data.Data[0].GetMonthsAgo())
	require.Equal(t, int32(3), generator.addingToMonth.GetMonthsAgo())
	require.Nil(t, generator.addingToSegment)
}

// TestNewMonthData_MultipleMonths opens a month 3 months ago then 2 months ago.
// The test verifies that the generator is set to add to the correct month. We
// then open a current month, and verify that the generator will add to the
// current month.
func TestNewMonthData_MultipleMonths(t *testing.T) {
	generator := NewActivityLogData(nil).NewPreviousMonthData(3).NewPreviousMonthData(2)
	require.Equal(t, int32(2), generator.data.Data[1].GetMonthsAgo())
	require.Equal(t, int32(2), generator.addingToMonth.GetMonthsAgo())
	generator = generator.NewCurrentMonthData()
	require.True(t, generator.data.Data[2].GetCurrentMonth())
	require.True(t, generator.addingToMonth.GetCurrentMonth())
}

// TestNewCurrentMonthData_ClientsSeen calls ClientsSeen with 3 clients, and
// verifies that they are added to the input data
func TestNewCurrentMonthData_ClientsSeen(t *testing.T) {
	wantClients := []*generation.Client{
		{
			Id:         "1",
			Namespace:  "ns",
			Mount:      "mount",
			ClientType: "non-entity",
		},
		{
			Id: "2",
		},
		{
			Id:    "3",
			Count: int32(3),
		},
	}
	generator := NewActivityLogData(nil).NewCurrentMonthData().ClientsSeen(wantClients...)
	require.Equal(t, generator.data.Data[0].GetAll().Clients, wantClients)
	require.True(t, generator.data.Data[0].GetCurrentMonth())
}

// TestSegment_AddClients adds clients in a variety of ways to an open segment
// and verifies that the clients are present in the segment with the correct
// options
func TestSegment_AddClients(t *testing.T) {
	testAddClients(t, func() *ActivityLogDataGenerator {
		return NewActivityLogData(nil).NewCurrentMonthData().Segment()
	}, func(g *ActivityLogDataGenerator) *generation.Client {
		return g.data.Data[0].GetSegments().Segments[0].Clients.Clients[0]
	})
}

// TestSegment_MultipleSegments opens a current month and adds a client to an
// un-indexed segment, then opens an indexed segment and adds a client. The test
// verifies that clients are present in both segments, and that the segment
// index is correctly recorded
func TestSegment_MultipleSegments(t *testing.T) {
	generator := NewActivityLogData(nil).NewCurrentMonthData().Segment().NewClientSeen().Segment(WithSegmentIndex(2)).NewClientSeen()
	require.Len(t, generator.data.Data[0].GetSegments().Segments[0].Clients.Clients, 1)
	require.Len(t, generator.data.Data[0].GetSegments().Segments[1].Clients.Clients, 1)
	require.Equal(t, int32(2), *generator.data.Data[0].GetSegments().Segments[1].SegmentIndex)
	require.Equal(t, int32(2), *generator.addingToSegment.SegmentIndex)
}

// TestSegment_NewMonth adds a client to a segment, then starts a new month. The
// test verifies that there are no open segments
func TestSegment_NewMonth(t *testing.T) {
	generator := NewActivityLogData(nil).NewCurrentMonthData().Segment().NewClientSeen().NewPreviousMonthData(1)
	require.Nil(t, generator.addingToSegment)
}

// TestNewCurrentMonthData_AddClients adds clients in a variety of ways to an
// the current month and verifies that the clients are present in the month with
// the correct options
func TestNewCurrentMonthData_AddClients(t *testing.T) {
	testAddClients(t, func() *ActivityLogDataGenerator {
		return NewActivityLogData(nil).NewCurrentMonthData()
	}, func(g *ActivityLogDataGenerator) *generation.Client {
		return g.data.Data[0].GetAll().Clients[0]
	})
}

// TestWrite creates a mock http server and writes generated data to it. The
// test verifies that the returned paths are parsed correctly, and that the JSON
// sent to the server is correct.
func TestWrite(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `{"data":{"paths":["path1","path2"]}}`)
		require.NoError(t, err)
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		raw := map[string]string{}
		err = json.Unmarshal(body, &raw)
		require.NoError(t, err)
		require.JSONEq(t, `{"write":["WRITE_ENTITIES"],"data":[{"monthsAgo":3,"all":{"clients":[{"count":1}]}},{"monthsAgo":2,"segments":{"segments":[{"segmentIndex":2,"clients":{"clients":[{"count":1,"repeated":true}]}}]}},{"currentMonth":true}]}`, raw["input"])
	}))
	defer ts.Close()

	client, err := api.NewClient(&api.Config{
		Address: ts.URL,
	})
	require.NoError(t, err)
	paths, err := NewActivityLogData(client).
		NewPreviousMonthData(3).
		NewClientSeen().
		NewPreviousMonthData(2).
		Segment(WithSegmentIndex(2)).
		RepeatedClientSeen().
		NewCurrentMonthData().Write(context.Background(), generation.WriteOptions_WRITE_ENTITIES)

	require.NoError(t, err)
	require.Equal(t, []string{"path1", "path2"}, paths)
}

func testAddClients(t *testing.T, makeGenerator func() *ActivityLogDataGenerator, getClient func(data *ActivityLogDataGenerator) *generation.Client) {
	t.Helper()
	clientOptions := []ClientOption{
		WithClientNamespace("ns"), WithClientMount("mount"), WithClientIsNonEntity(), WithClientID("1"),
	}
	generator := makeGenerator().NewClientSeen(clientOptions...)
	require.Equal(t, getClient(generator), &generation.Client{
		Id:         "1",
		Count:      1,
		Namespace:  "ns",
		Mount:      "mount",
		ClientType: "non-entity",
	})

	generator = makeGenerator().NewClientsSeen(4, clientOptions...)
	require.Equal(t, getClient(generator), &generation.Client{
		Id:         "1",
		Count:      4,
		Namespace:  "ns",
		Mount:      "mount",
		ClientType: "non-entity",
	})

	generator = makeGenerator().RepeatedClientSeen(clientOptions...)
	require.Equal(t, getClient(generator), &generation.Client{
		Id:         "1",
		Count:      1,
		Repeated:   true,
		Namespace:  "ns",
		Mount:      "mount",
		ClientType: "non-entity",
	})

	generator = makeGenerator().RepeatedClientsSeen(4, clientOptions...)
	require.Equal(t, getClient(generator), &generation.Client{
		Id:         "1",
		Count:      4,
		Repeated:   true,
		Namespace:  "ns",
		Mount:      "mount",
		ClientType: "non-entity",
	})

	generator = makeGenerator().RepeatedClientSeenFromMonthsAgo(3, clientOptions...)
	require.Equal(t, getClient(generator), &generation.Client{
		Id:                "1",
		Count:             1,
		RepeatedFromMonth: 3,
		Namespace:         "ns",
		Mount:             "mount",
		ClientType:        "non-entity",
	})

	generator = makeGenerator().RepeatedClientsSeenFromMonthsAgo(4, 3, clientOptions...)
	require.Equal(t, getClient(generator), &generation.Client{
		Id:                "1",
		Count:             4,
		RepeatedFromMonth: 3,
		Namespace:         "ns",
		Mount:             "mount",
		ClientType:        "non-entity",
	})
}

// TestSetMonthOptions sets month options and verifies that they are saved
func TestSetMonthOptions(t *testing.T) {
	generator := NewActivityLogData(nil).NewCurrentMonthData().SetMonthOptions(WithEmptySegmentIndexes(3, 4),
		WithMaximumSegmentIndex(7), WithSkipSegmentIndexes(1, 2))
	require.Equal(t, int32(7), generator.data.Data[0].NumSegments)
	require.Equal(t, []int32{3, 4}, generator.data.Data[0].EmptySegmentIndexes)
	require.Equal(t, []int32{1, 2}, generator.data.Data[0].SkipSegmentIndexes)
}

// TestVerifyInput constructs invalid inputs and ensures that VerifyInput
// returns an error
func TestVerifyInput(t *testing.T) {
	cases := []struct {
		name      string
		generator *ActivityLogDataGenerator
	}{
		{
			name: "repeated client with only 1 month",
			generator: NewActivityLogData(nil).
				NewCurrentMonthData().
				RepeatedClientSeen(),
		},
		{
			name: "repeated client with segment",
			generator: NewActivityLogData(nil).
				NewCurrentMonthData().
				Segment().
				RepeatedClientSeen(),
		},
		{
			name: "repeated client with earliest month",
			generator: NewActivityLogData(nil).
				NewCurrentMonthData().
				NewClientSeen().
				NewPreviousMonthData(2).
				RepeatedClientSeen(),
		},
		{
			name: "repeated month",
			generator: NewActivityLogData(nil).
				NewPreviousMonthData(1).
				NewPreviousMonthData(1),
		},
		{
			name: "repeated current month",
			generator: NewActivityLogData(nil).
				NewCurrentMonthData().
				NewCurrentMonthData(),
		},
		{
			name: "repeated segment index",
			generator: NewActivityLogData(nil).
				NewCurrentMonthData().
				Segment(WithSegmentIndex(1)).
				Segment(WithSegmentIndex(1)),
		},
		{
			name: "segment with num segments",
			generator: NewActivityLogData(nil).
				NewCurrentMonthData().
				Segment().
				SetMonthOptions(WithMaximumSegmentIndex(1)),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Error(t, VerifyInput(tc.generator.data))
		})
	}
}
