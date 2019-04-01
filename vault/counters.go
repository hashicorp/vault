package vault

import (
	"context"
	"sort"
	"sync/atomic"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
)

const (
	requestCounterDatePathFormat = "2006/01"
	countersPath                 = systemBarrierPrefix + "counters"
	requestCountersRelPath       = "counters/requests/"
)

type counters struct {
	// requests counts requests seen by Vault this month; does not include requests
	// excluded by design, e.g. health checks and UI asset requests.
	requests *uint64
	// activePath is set at startup to the path we primed the requests counter from,
	// or empty string if there wasn't a relevant path - either because this is the first
	// time Vault starts with the feature enabled, or because Vault hadn't written
	// out the request counter this month yet.
	// Whenever we write out the counters, we update activePath if it's no longer
	// accurate.  This coincides with a reset of the counters.
	// There's no lock because the only reader/writer of activePath is the goroutine
	// doing background syncs.
	activePath string
	// syncInterval determines how often the counters get written to storage (on primary)
	// or synced to primary.
	syncInterval time.Duration
}

// RequestCounter stores the state of request counters for a single unspecified period.
type RequestCounter struct {
	// Total is the number of requests seen during a given period.
	Total *uint64 `json:"total"`
}

// DatedRequestCounter holds request counters from a single period of time.
type DatedRequestCounter struct {
	// StartTime is when the period starts.
	StartTime time.Time `json:"start_time"`
	// RequestCounter counts requests.
	RequestCounter
}

// loadAllRequestCounters returns all request counters found in storage,
// ordered by time (oldest first.)
func (c *Core) loadAllRequestCounters(ctx context.Context, now time.Time) ([]DatedRequestCounter, error) {
	view := c.systemBarrierView.SubView(requestCountersRelPath)

	datepaths, err := view.List(ctx, "")
	if err != nil {
		return nil, errwrap.Wrapf("failed to read request counters: {{err}}", err)
	}

	var all []DatedRequestCounter
	sort.Strings(datepaths)
	for _, datepath := range datepaths {
		datesubpaths, err := view.List(ctx, datepath)
		if err != nil {
			return nil, errwrap.Wrapf("failed to read request counters: {{err}}", err)
		}
		sort.Strings(datesubpaths)
		for _, datesubpath := range datesubpaths {
			fullpath := datepath + datesubpath
			counter, err := c.loadRequestCounters(ctx, fullpath)
			if err != nil {
				return nil, err
			}

			t, err := time.Parse(requestCounterDatePathFormat, fullpath)
			if err != nil {
				return nil, err
			}

			all = append(all, DatedRequestCounter{StartTime: t, RequestCounter: *counter})
		}
	}

	start, _ := time.Parse(requestCounterDatePathFormat, now.Format(requestCounterDatePathFormat))
	idx := sort.Search(len(all), func(i int) bool {
		return !all[i].StartTime.Before(start)
	})
	cur := atomic.LoadUint64(c.counters.requests)
	if idx < len(all) {
		all[idx].RequestCounter.Total = &cur
	} else {
		all = append(all, DatedRequestCounter{StartTime: start, RequestCounter: RequestCounter{Total: &cur}})
	}

	return all, nil
}

// loadCurrentRequestCounters reads the current RequestCounter out of storage.
// The in-memory current request counter is populated with the value read, if any.
// now should be the current time; it is a parameter to facilitate testing.
func (c *Core) loadCurrentRequestCounters(ctx context.Context, now time.Time) error {
	datepath := now.Format(requestCounterDatePathFormat)
	counter, err := c.loadRequestCounters(ctx, datepath)
	if err != nil {
		return err
	}
	if counter != nil {
		c.counters.activePath = datepath
		atomic.StoreUint64(c.counters.requests, *counter.Total)
	}
	return nil
}

// loadRequestCounters reads a RequestCounter out of storage at location datepath.
// If nothing is found at that path, that isn't an error: a reference to a zero
// RequestCounter is returned.
func (c *Core) loadRequestCounters(ctx context.Context, datepath string) (*RequestCounter, error) {
	view := c.systemBarrierView.SubView(requestCountersRelPath)

	out, err := view.Get(ctx, datepath)
	if err != nil {
		return nil, errwrap.Wrapf("failed to read request counters: {{err}}", err)
	}
	if out == nil {
		return nil, nil
	}

	newCounters := &RequestCounter{}
	err = out.DecodeJSON(newCounters)
	if err != nil {
		return nil, err
	}

	return newCounters, nil
}

// saveCurrentRequestCounters writes the current RequestCounter to storage.
// The in-memory current request counter is reset to zero after writing if
// we've entered a new month.
// now should be the current time; it is a parameter to facilitate testing.
func (c *Core) saveCurrentRequestCounters(ctx context.Context, now time.Time) error {
	view := c.systemBarrierView.SubView(requestCountersRelPath)
	requests := atomic.LoadUint64(c.counters.requests)
	curDatePath := now.Format(requestCounterDatePathFormat)

	// If activePath is empty string, we were started with nothing in storage
	// for the current month, so we should not reset the in-mem counter.
	// But if activePath is nonempty and not curDatePath, we should reset.
	shouldReset, writeDatePath := false, curDatePath
	if c.counters.activePath != "" && c.counters.activePath != curDatePath {
		shouldReset, writeDatePath = true, c.counters.activePath
	}

	localCounters := &RequestCounter{
		Total: &requests,
	}
	entry, err := logical.StorageEntryJSON(writeDatePath, localCounters)
	if err != nil {
		return errwrap.Wrapf("failed to create request counters entry: {{err}}", err)
	}

	if err := view.Put(ctx, entry); err != nil {
		return errwrap.Wrapf("failed to save request counters: {{err}}", err)
	}

	if shouldReset {
		atomic.StoreUint64(c.counters.requests, 0)
	}
	if c.counters.activePath != curDatePath {
		c.counters.activePath = curDatePath
	}

	return nil
}
