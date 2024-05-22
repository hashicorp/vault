// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package vault

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestActivityLog_doPrecomputedQueryCreation creates segments for the last 4
// months and then calls doPrecomputedQueryCreation, in order of oldest to most
// recent month. The test verifies that the count of clients in the generated
// precomputed query is equal to the number of deduplicated clients.
func TestActivityLog_doPrecomputedQueryCreation(t *testing.T) {
	core, _, token := TestCoreUnsealed(t)
	a := core.activityLog
	a.SetEnable(true)

	j, err := clientcountutil.NewActivityLogData(nil).
		// 8 new clients
		// across two segments
		NewPreviousMonthData(4).
		Segment().NewClientsSeen(5).
		Segment().NewClientsSeen(3).

		// 2 repeated clients
		// 10 new clients
		// across 3 segments
		NewPreviousMonthData(3).
		Segment().RepeatedClientsSeen(2).
		NewClientsSeen(3).
		Segment().NewClientsSeen(2).
		Segment().NewClientsSeen(5).

		// 7 new clients
		// single segment
		NewPreviousMonthData(2).
		NewClientsSeen(7).

		// 6 repeated clients
		// 5 new clients
		// across 2 segments
		NewPreviousMonthData(1).
		Segment().NewClientsSeen(5).
		Segment().RepeatedClientsSeen(6).
		ToJSON(generation.WriteOptions_WRITE_ENTITIES)
	require.NoError(t, err)

	r := logical.TestRequest(t, logical.UpdateOperation, "sys/internal/counters/activity/write")
	r.Data["input"] = string(j)
	r.ClientToken = token
	_, err = core.HandleRequest(namespace.RootContext(context.Background()), r)
	require.NoError(t, err)

	now := time.Now().UTC()
	times := map[int]time.Time{}
	for i := 1; i < 5; i++ {
		times[i] = timeutil.StartOfMonth(timeutil.MonthsPreviousTo(i, now))
	}

	testCases := []struct {
		name              string
		generateUpToMonth int
		strictEnforcement bool
		wantClients       int
	}{
		{
			name:              "only 4 months ago",
			generateUpToMonth: 4,
			wantClients:       8, // 8 clients from month 4
		},
		{
			name:              "3 months ago",
			generateUpToMonth: 3,
			// 8 clients (month 4) + 10 new clients (month 3)
			wantClients: 18,
		},
		{
			name:              "2 months ago",
			generateUpToMonth: 2,
			// 8 clients (month 4) + 10 new clients (month 3) + 7 new clients
			// (month 2)
			wantClients: 25,
		},
		{
			name:              "1 month ago",
			generateUpToMonth: 1,
			// 8 clients (month 4) + 10 new clients (month 3) + 7 new clients
			// (month 2) + 5 new clients (month 1)
			wantClients: 30,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			generateUpTo := times[tc.generateUpToMonth]
			nextMonth := timeutil.StartOfNextMonth(generateUpTo)
			err = a.precomputedQueryWorker(context.Background(), &ActivityIntentLog{PreviousMonth: generateUpTo.Unix(), NextMonth: nextMonth.Unix()})
			require.NoError(t, err)

			// get precomputed queries spanning the whole time period
			pq, err := a.queryStore.Get(context.Background(), times[4], now)
			require.NoError(t, err)
			require.Equal(t, tc.wantClients, int(pq.Namespaces[0].Entities))
		})
	}
}
