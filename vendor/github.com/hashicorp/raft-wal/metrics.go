// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package wal

import (
	"github.com/hashicorp/raft-wal/metrics"
)

var (
	// MetricDefinitions describe the metrics emitted by this library via the
	// provided metrics.Collector implementation. It's public so that these can be
	// registered during init with metrics clients that support pre-defining
	// metrics.
	MetricDefinitions = metrics.Definitions{
		Counters: []metrics.Descriptor{
			{
				Name: "log_entry_bytes_written",
				Desc: "log_entry_bytes_written counts the bytes of log entry after encoding" +
					" with Codec. Actual bytes written to disk might be slightly higher as it" +
					" includes headers and index entries.",
			},
			{
				Name: "log_entries_written",
				Desc: "log_entries_written counts the number of entries written.",
			},
			{
				Name: "log_appends",
				Desc: "log_appends counts the number of calls to StoreLog(s) i.e." +
					" number of batches of entries appended.",
			},
			{
				Name: "log_entry_bytes_read",
				Desc: "log_entry_bytes_read counts the bytes of log entry read from" +
					" segments before decoding. actual bytes read from disk might be higher" +
					" as it includes headers and index entries and possible secondary reads" +
					" for large entries that don't fit in buffers.",
			},
			{
				Name: "log_entries_read",
				Desc: "log_entries_read counts the number of calls to get_log.",
			},
			{
				Name: "segment_rotations",
				Desc: "segment_rotations counts how many times we move to a new segment file.",
			},
			{
				Name: "head_truncations",
				Desc: "head_truncations counts how many log entries have been truncated" +
					" from the head - i.e. the oldest entries. by graphing the rate of" +
					" change over time you can see individual truncate calls as spikes.",
			},
			{
				Name: "tail_truncations",
				Desc: "tail_truncations counts how many log entries have been truncated" +
					" from the head - i.e. the newest entries. by graphing the rate of" +
					" change over time you can see individual truncate calls as spikes.",
			},
			{
				Name: "stable_gets",
				Desc: "stable_gets counts how many calls to StableStore.Get or GetUint64.",
			},
			{
				Name: "stable_sets",
				Desc: "stable_sets counts how many calls to StableStore.Set or SetUint64.",
			},
		},
		Gauges: []metrics.Descriptor{
			{
				Name: "last_segment_age_seconds",
				Desc: "last_segment_age_seconds is a gauge that is set each time we" +
					" rotate a segment and describes the number of seconds between when" +
					" that segment file was first created and when it was sealed. this" +
					" gives a rough estimate how quickly writes are filling the disk.",
			},
		},
	}
)
