// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package verifier

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
				Name: "checkpoints_written",
				Desc: "checkpoints_written counts the number of checkpoint entries" +
					" written to the LogStore.",
			},
			{
				Name: "ranges_verified",
				Desc: "ranges_verified counts the number of log ranges for which a" +
					" verification report has been completed.",
			},
			{
				Name: "read_checksum_failures",
				Desc: "read_checksum_failures counts the number of times a range of" +
					" logs between two check points contained at least one corruption.",
			},
			{
				Name: "write_checksum_failures",
				Desc: "write_checksum_failures counts the number of times a follower" +
					" has a different checksum to the leader at the point where it" +
					" writes to the log. This could be caused by either a disk-corruption" +
					" on the leader (unlikely) or some other corruption of the log" +
					" entries in-flight.",
			},
			{
				Name: "dropped_reports",
				Desc: "dropped_reports counts how many times the verifier routine was" +
					" still busy when the next checksum came in and so verification for" +
					" a range was skipped. If you see this happen consider increasing" +
					" the interval between checkpoints.",
			},
		},
	}
)
