package vault

import (
	"context"
	"sort"
	"sync/atomic"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
)

const requestCounterDatePathFormat = "2006/01"

// RequestCounter stores the state of request counters for a single unspecified period.
type RequestCounter struct {
	// Total holds the sum total of all requests seen during the period.
	// "All" does not include requests excluded by design, e.g. health checks and UI
	// asset requests.
	Total *uint64 `json:"total"`
}

// DatedRequestCounter holds request counters from a single period of time.
type DatedRequestCounter struct {
	// StartTime is when the period starts.
	StartTime time.Time `json:"start_time"`
	// RequestCounter counts requests.
	RequestCounter
}

// AllRequestCounters contains all request counters from the dawn of time.
type AllRequestCounters struct {
	// Dated holds the request counters dating back to when the feature was first
	// introduced in this instance, ordered by time (oldest first).
	Dated []DatedRequestCounter
}

// loadAllRequestCounters returns all request counters found in storage.
func (c *Core) loadAllRequestCounters(ctx context.Context) (*AllRequestCounters, error) {
	view := c.systemBarrierView.SubView("counters/requests/")

	datepaths, err := view.List(ctx, "")
	if err != nil {
		return nil, errwrap.Wrapf("failed to read request counters: {{err}}", err)
	}
	if datepaths == nil {
		return nil, nil
	}

	var all AllRequestCounters
	sort.Strings(datepaths)
	for _, datepath := range datepaths {
		datesubpaths, err := view.List(ctx, datepath)
		if err != nil {
			return nil, errwrap.Wrapf("failed to read request counters: {{err}}", err)
		}
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

			all.Dated = append(all.Dated, DatedRequestCounter{StartTime: t, RequestCounter: *counter})
		}
	}

	return &all, nil
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
		atomic.StoreUint64(&c.counters.requests, *counter.Total)
	}
	return nil
}

// loadRequestCounters reads a RequestCounter out of storage at location datepath.
// If nothing is found at that path, that isn't an error: a reference to a zero
// RequestCounter is returned.
func (c *Core) loadRequestCounters(ctx context.Context, datepath string) (*RequestCounter, error) {
	view := c.systemBarrierView.SubView("counters/requests/")

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
	view := c.systemBarrierView.SubView("counters/requests/")
	datepath := now.Format(requestCounterDatePathFormat)

	var requests uint64
	if c.counters.activePath == "" || c.counters.activePath == datepath {
		// We leave requests at 0 in the case where the month has just changed; we don't
		// want to write the current count of requests (from last month) to the new month
		// datepath.  This means we discard any requests counted since the last time
		// we wrote, but since we write frequently that's okay.
		requests = atomic.LoadUint64(&c.counters.requests)
	}

	localCounters := &RequestCounter{
		Total: &requests,
	}
	entry, err := logical.StorageEntryJSON(datepath, localCounters)
	if err != nil {
		return errwrap.Wrapf("failed to create request counters entry: {{err}}", err)
	}

	if err := view.Put(ctx, entry); err != nil {
		return errwrap.Wrapf("failed to save request counters: {{err}}", err)
	}

	if datepath != c.counters.activePath {
		if c.counters.activePath != "" {
			atomic.StoreUint64(&c.counters.requests, 0)
		}
		c.counters.activePath = datepath
	}

	return nil
}
