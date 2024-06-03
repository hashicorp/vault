// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package activity

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
)

func NewTestQueryStore(t *testing.T) *PrecomputedQueryStore {
	t.Helper()

	logger := logging.NewVaultLogger(log.Trace)
	view := &logical.InmemStorage{}
	return NewPrecomputedQueryStore(logger, view, 12)
}

func TestQueryStore_Inventory(t *testing.T) {
	startTimes := []time.Time{
		time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
	}

	endTimes := []time.Time{
		timeutil.EndOfMonth(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
		timeutil.EndOfMonth(time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)),
		timeutil.EndOfMonth(time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)),
		timeutil.EndOfMonth(time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)),
		timeutil.EndOfMonth(time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC)),
	}

	qs := NewTestQueryStore(t)
	ctx := context.Background()

	for _, s := range startTimes {
		for _, e := range endTimes {
			if e.Before(s) {
				continue
			}
			qs.Put(ctx, &PrecomputedQuery{
				StartTime:  s,
				EndTime:    e,
				Namespaces: []*NamespaceRecord{},
			})
		}
	}

	storedStartTimes, err := qs.listStartTimes(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(storedStartTimes) != len(startTimes) {
		t.Fatalf("bad length, expected %v got %v", len(startTimes), storedStartTimes)
	}
	sort.Slice(storedStartTimes, func(i, j int) bool {
		return storedStartTimes[i].Before(storedStartTimes[j])
	})
	if !reflect.DeepEqual(storedStartTimes, startTimes) {
		t.Fatalf("start time mismatch, expected %v got %v", startTimes, storedStartTimes)
	}

	storedEndTimes, err := qs.listEndTimes(ctx, startTimes[1])
	if err != nil {
		t.Fatal(err)
	}
	expected := endTimes[1:]
	if len(storedEndTimes) != len(expected) {
		t.Fatalf("bad length, expected %v got %v", len(expected), storedEndTimes)
	}
	sort.Slice(storedEndTimes, func(i, j int) bool {
		return storedEndTimes[i].Before(storedEndTimes[j])
	})
	if !reflect.DeepEqual(storedEndTimes, expected) {
		t.Fatalf("end time mismatch, expected %v got %v", expected, storedEndTimes)
	}
}

func TestQueryStore_MarshalDemarshal(t *testing.T) {
	tsStart := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	tsEnd := timeutil.EndOfMonth(tsStart)

	p := &PrecomputedQuery{
		StartTime: tsStart,
		EndTime:   tsEnd,
		Namespaces: []*NamespaceRecord{
			{
				NamespaceID:     "root",
				Entities:        20,
				NonEntityTokens: 42,
			},
			{
				NamespaceID:     "yzABC",
				Entities:        15,
				NonEntityTokens: 31,
			},
		},
	}

	qs := NewTestQueryStore(t)
	ctx := context.Background()
	qs.Put(ctx, p)
	result, err := qs.Get(ctx, tsStart, tsEnd)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("nil response from Get")
	}
	if !reflect.DeepEqual(result, p) {
		t.Fatalf("unequal query objects, expected %v got %v", p, result)
	}
}

func TestQueryStore_TimeRanges(t *testing.T) {
	qs := NewTestQueryStore(t)
	ctx := context.Background()

	// Scenario ranges: Jan 15 - Jan 31 (one month)
	//                  Feb 2 - Mar 31 (two months, but not contiguous)
	//                  April and May are skipped
	//                  June 1 - September 30 (4 months)
	periods := []struct {
		Begin time.Time
		Ends  []time.Time
	}{
		{
			time.Date(2020, 1, 15, 12, 45, 53, 0, time.UTC),
			[]time.Time{
				timeutil.EndOfMonth(time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC)),
			},
		},
		{
			time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				timeutil.EndOfMonth(time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)),
				timeutil.EndOfMonth(time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC),
			[]time.Time{
				timeutil.EndOfMonth(time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)),
				timeutil.EndOfMonth(time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)),
				timeutil.EndOfMonth(time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC)),
				timeutil.EndOfMonth(time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
	}

	for _, period := range periods {
		for _, e := range period.Ends {
			qs.Put(ctx, &PrecomputedQuery{
				StartTime: period.Begin,
				EndTime:   e,
				Namespaces: []*NamespaceRecord{
					{
						NamespaceID:     "root",
						Entities:        17,
						NonEntityTokens: 31,
					},
				},
			})
		}
	}

	testCases := []struct {
		Name          string
		StartTime     time.Time
		EndTime       time.Time
		Empty         bool
		ExpectedStart time.Time
		ExpectedEnd   time.Time
	}{
		{
			"year query in October",
			time.Date(2019, 10, 12, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC),
			false,
			// June - Sept
			periods[2].Begin,
			periods[2].Ends[3],
		},
		{
			"one day in January",
			time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
			false,
			// January, even though this is outside the range specified
			periods[0].Begin,
			periods[0].Ends[0],
		},
		{
			"one day in February",
			time.Date(2020, 2, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 2, 5, 0, 0, 0, 0, time.UTC),
			false,
			// February only
			periods[1].Begin,
			periods[1].Ends[0],
		},
		{
			"January through March",
			time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 3, 5, 0, 0, 0, 0, time.UTC),
			false,
			// February and March only
			// Fails due to bug in library function, TODO
			periods[1].Begin,
			periods[1].Ends[1],
		},
		{
			"the month of May",
			time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 5, 31, 0, 0, 0, 0, time.UTC),
			true, // no data
			time.Time{},
			time.Time{},
		},
		{
			"May through June",
			time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC),
			false,
			// June only
			periods[2].Begin,
			periods[2].Ends[0],
		},
		{
			"September",
			time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC),
			true, // We have June through September,
			// but not anything starting in September
			// (which does not match a real scenario)
			time.Time{},
			time.Time{},
		},
		{
			"December",
			time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC),
			true, // no data
			time.Time{},
			time.Time{},
		},
		{
			"June through December",
			time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC),
			time.Date(2020, 12, 31, 12, 0, 0, 0, time.UTC),
			false,
			// June through September
			periods[2].Begin,
			periods[2].Ends[3],
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			result, err := qs.Get(ctx, tc.StartTime, tc.EndTime)
			if err != nil {
				t.Fatal(err)
			}
			if result == nil {
				if tc.Empty {
					return
				} else {
					t.Fatal("unexpected empty result")
				}
			} else {
				if tc.Empty {
					t.Fatal("expected empty result")
				}
			}
			if !result.StartTime.Equal(tc.ExpectedStart) {
				t.Errorf("start time mismatch: %v, expected %v", result.StartTime, tc.ExpectedStart)
			}
			if !result.EndTime.Equal(tc.ExpectedEnd) {
				t.Errorf("end time mismatch: %v, expected %v", result.EndTime, tc.ExpectedEnd)
			}
		})
	}
}
