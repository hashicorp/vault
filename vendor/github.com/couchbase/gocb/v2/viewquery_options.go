package gocb

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

// ViewScanConsistency specifies the consistency required for a view query.
type ViewScanConsistency uint

const (
	// ViewScanConsistencyNotBounded indicates that no special behaviour should be used.
	ViewScanConsistencyNotBounded ViewScanConsistency = iota + 1
	// ViewScanConsistencyRequestPlus indicates to update the index before querying it.
	ViewScanConsistencyRequestPlus
	// ViewScanConsistencyUpdateAfter indicates to update the index asynchronously after querying.
	ViewScanConsistencyUpdateAfter
)

// ViewOrdering specifies the ordering for the view queries results.
type ViewOrdering uint

const (
	// ViewOrderingAscending indicates the query results should be sorted from lowest to highest.
	ViewOrderingAscending ViewOrdering = iota + 1
	// ViewOrderingDescending indicates the query results should be sorted from highest to lowest.
	ViewOrderingDescending
)

// ViewErrorMode specifies the behaviour of the query engine should an error occur during the gathering of
// view index results which would result in only partial results being available.
type ViewErrorMode uint

const (
	// ViewErrorModeContinue indicates to continue gathering results on error.
	ViewErrorModeContinue ViewErrorMode = iota + 1

	// ViewErrorModeStop indicates to stop gathering results on error
	ViewErrorModeStop
)

// ViewOptions represents the options available when executing view query.
type ViewOptions struct {
	ScanConsistency ViewScanConsistency
	Skip            uint32
	Limit           uint32
	Order           ViewOrdering
	Reduce          bool
	Group           bool
	GroupLevel      uint32
	Key             interface{}
	Keys            []interface{}
	StartKey        interface{}
	EndKey          interface{}
	InclusiveEnd    bool
	StartKeyDocID   string
	EndKeyDocID     string
	OnError         ViewErrorMode
	Debug           bool
	ParentSpan      RequestSpan

	// Raw provides a way to provide extra parameters in the request body for the query.
	Raw map[string]string

	Namespace DesignDocumentNamespace

	Timeout       time.Duration
	RetryStrategy RetryStrategy

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

func (opts *ViewOptions) toURLValues() (*url.Values, error) {
	options := &url.Values{}

	if opts.ScanConsistency != 0 {
		if opts.ScanConsistency == ViewScanConsistencyRequestPlus {
			options.Set("stale", "false")
		} else if opts.ScanConsistency == ViewScanConsistencyNotBounded {
			options.Set("stale", "ok")
		} else if opts.ScanConsistency == ViewScanConsistencyUpdateAfter {
			options.Set("stale", "update_after")
		} else {
			return nil, makeInvalidArgumentsError("unexpected stale option")
		}
	}

	if opts.Skip != 0 {
		options.Set("skip", strconv.FormatUint(uint64(opts.Skip), 10))
	}

	if opts.Limit != 0 {
		options.Set("limit", strconv.FormatUint(uint64(opts.Limit), 10))
	}

	if opts.Order != 0 {
		if opts.Order == ViewOrderingAscending {
			options.Set("descending", "false")
		} else if opts.Order == ViewOrderingDescending {
			options.Set("descending", "true")
		} else {
			return nil, makeInvalidArgumentsError("unexpected order option")
		}
	}

	options.Set("reduce", "false") // is this line necessary?
	if opts.Reduce {
		options.Set("reduce", "true")

		// Only set group if a reduce view
		if opts.Group {
			options.Set("group", "true")
		}

		if opts.GroupLevel != 0 {
			options.Set("group_level", strconv.FormatUint(uint64(opts.GroupLevel), 10))
		}
	}

	if opts.Key != nil {
		jsonKey, err := opts.marshalJSON(opts.Key)
		if err != nil {
			return nil, err
		}
		options.Set("key", string(jsonKey))
	}

	if len(opts.Keys) > 0 {
		jsonKeys, err := opts.marshalJSON(opts.Keys)
		if err != nil {
			return nil, err
		}
		options.Set("keys", string(jsonKeys))
	}

	if opts.StartKey != nil {
		jsonStartKey, err := opts.marshalJSON(opts.StartKey)
		if err != nil {
			return nil, err
		}
		options.Set("startkey", string(jsonStartKey))
	} else {
		options.Del("startkey")
	}

	if opts.EndKey != nil {
		jsonEndKey, err := opts.marshalJSON(opts.EndKey)
		if err != nil {
			return nil, err
		}
		options.Set("endkey", string(jsonEndKey))
	} else {
		options.Del("endkey")
	}

	if opts.StartKey != nil || opts.EndKey != nil {
		if opts.InclusiveEnd {
			options.Set("inclusive_end", "true")
		} else {
			options.Set("inclusive_end", "false")
		}
	}

	if opts.StartKeyDocID == "" {
		options.Del("startkey_docid")
	} else {
		options.Set("startkey_docid", opts.StartKeyDocID)
	}

	if opts.EndKeyDocID == "" {
		options.Del("endkey_docid")
	} else {
		options.Set("endkey_docid", opts.EndKeyDocID)
	}

	if opts.OnError > 0 {
		if opts.OnError == ViewErrorModeContinue {
			options.Set("on_error", "continue")
		} else if opts.OnError == ViewErrorModeStop {
			options.Set("on_error", "stop")
		} else {
			return nil, makeInvalidArgumentsError("unexpected onerror option")
		}
	}

	if opts.Debug {
		options.Set("debug", "true")
	}

	if opts.Raw != nil {
		for k, v := range opts.Raw {
			options.Set(k, v)
		}
	}

	return options, nil
}

func (opts *ViewOptions) marshalJSON(value interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
