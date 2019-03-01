package vault

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-test/deep"
)

//noinspection SpellCheckingInspection
func testParseTime(t *testing.T, format, timeval string) time.Time {
	t.Helper()
	tm, err := time.Parse(format, timeval)
	if err != nil {
		t.Fatalf("Error parsing time '%s': %v", timeval, err)
	}
	return tm
}

// TestRequestCounterLoadCurrent exercises the code that primes the in-mem
// request counters from persistent storage.
func TestRequestCounterLoadCurrent(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	december2018 := testParseTime(t, time.RFC3339, "2018-12-05T09:44:12-05:00")
	decemberRequests := uint64(555)

	// It's December, and we got some requests.  Persist the counter.
	atomic.StoreUint64(c.counters.requests, decemberRequests)
	err := c.saveCurrentRequestCounters(context.Background(), december2018)
	if err != nil {
		t.Fatal(err)
	}

	// It's still December, simulate being restarted.  At startup the counter is
	// zero initially, until we read the counter from storage post-unseal via
	// loadCurrentRequestCounters.
	atomic.StoreUint64(c.counters.requests, 0)
	err = c.loadCurrentRequestCounters(context.Background(), december2018)
	if err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadUint64(c.counters.requests); got != decemberRequests {
		t.Fatalf("expected=%d, got=%d", decemberRequests, got)
	}

	// Now simulate being restarted in January. We never wrote anything out during
	// January, so the in-mem counter should remain zero.
	january2019 := testParseTime(t, time.RFC3339, "2019-01-02T08:21:11-05:00")
	atomic.StoreUint64(c.counters.requests, 0)
	err = c.loadCurrentRequestCounters(context.Background(), january2019)
	if err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadUint64(c.counters.requests); got != 0 {
		t.Fatalf("expected=%d, got=%d", 0, got)
	}
}

// TestRequestCounterSaveCurrent exercises the code that saves the in-mem
// request counters to persistent storage.
func TestRequestCounterSaveCurrent(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	// storeSaveLoad stores newValue in the in-mem counter, saves it to storage,
	// then verifies in-mem counter has value expectedPostLoad.
	storeSaveLoad := func(newValue, expectedPostLoad uint64, now time.Time) {
		t.Helper()
		atomic.StoreUint64(c.counters.requests, newValue)
		err := c.saveCurrentRequestCounters(context.Background(), now)
		if err != nil {
			t.Fatal(err)
		}
		if got := atomic.LoadUint64(c.counters.requests); got != expectedPostLoad {
			t.Fatalf("expected=%d, got=%d", expectedPostLoad, got)
		}
	}

	// Start in December. The first write ever should persist the current in-mem value.
	december2018 := testParseTime(t, time.RFC3339, "2018-12-05T09:44:12-05:00")
	decemberRequests := uint64(555)
	storeSaveLoad(decemberRequests, decemberRequests, december2018)

	// Update request count.
	decemberRequests++
	storeSaveLoad(decemberRequests, decemberRequests, december2018)

	decemberStartTime := testParseTime(t, requestCounterDatePathFormat, december2018.Format(requestCounterDatePathFormat))
	expected2018 := []DatedRequestCounter{
		{StartTime: decemberStartTime, RequestCounter: RequestCounter{Total: &decemberRequests}},
	}
	all, err := c.loadAllRequestCounters(context.Background(), december2018)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(all, expected2018); len(diff) != 0 {
		t.Errorf("Expected=%v, got=%v, diff=%v", expected2018, all, diff)
	}

	// Now it's January. Saving after transition to new month should reset in-mem
	// counter to zero, and also write zero to storage for the new month.
	january2019 := testParseTime(t, time.RFC3339, "2019-01-02T08:21:11-05:00")
	decemberRequests += 5
	storeSaveLoad(decemberRequests, 0, january2019)

	januaryRequests := uint64(333)
	storeSaveLoad(januaryRequests, januaryRequests, january2019)

	all, err = c.loadAllRequestCounters(context.Background(), january2019)
	if err != nil {
		t.Fatal(err)
	}

	januaryStartTime := testParseTime(t, requestCounterDatePathFormat, january2019.Format(requestCounterDatePathFormat))
	expected2019 := expected2018
	expected2019 = append(expected2019,
		DatedRequestCounter{januaryStartTime, RequestCounter{&januaryRequests}})
	if diff := deep.Equal(all, expected2019); len(diff) != 0 {
		t.Errorf("Expected=%v, got=%v, diff=%v", expected2019, all, diff)
	}
}
