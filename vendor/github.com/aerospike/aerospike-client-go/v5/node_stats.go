// Copyright 2014-2021 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import "sync/atomic"

// nodeStats keeps track of client's internal node statistics
// These statistics are aggregated once per tend in the cluster object
// and then are served to the end-user.
type nodeStats struct {
	ConnectionsAttempts   int64 `json:"connections-attempts"`
	ConnectionsSuccessful int64 `json:"connections-successful"`
	ConnectionsFailed     int64 `json:"connections-failed"`
	ConnectionsPoolEmpty  int64 `json:"connections-pool-empty"`
	ConnectionsOpen       int64 `json:"open-connections"`
	ConnectionsClosed     int64 `json:"closed-connections"`
	TendsTotal            int64 `json:"tends-total"`
	TendsSuccessful       int64 `json:"tends-successful"`
	TendsFailed           int64 `json:"tends-failed"`
	PartitionMapUpdates   int64 `json:"partition-map-updates"`
	NodeAdded             int64 `json:"node-added-count"`
	NodeRemoved           int64 `json:"node-removed-count"`
}

// latest returns the latest values to be used in aggregation and then resets the values
func (ns *nodeStats) getAndReset() *nodeStats {
	return &nodeStats{
		ConnectionsAttempts:   atomic.SwapInt64(&ns.ConnectionsAttempts, 0),
		ConnectionsSuccessful: atomic.SwapInt64(&ns.ConnectionsSuccessful, 0),
		ConnectionsFailed:     atomic.SwapInt64(&ns.ConnectionsFailed, 0),
		ConnectionsPoolEmpty:  atomic.SwapInt64(&ns.ConnectionsPoolEmpty, 0),
		ConnectionsOpen:       atomic.SwapInt64(&ns.ConnectionsOpen, 0),
		ConnectionsClosed:     atomic.SwapInt64(&ns.ConnectionsClosed, 0),
		TendsTotal:            atomic.SwapInt64(&ns.TendsTotal, 0),
		TendsSuccessful:       atomic.SwapInt64(&ns.TendsSuccessful, 0),
		TendsFailed:           atomic.SwapInt64(&ns.TendsFailed, 0),
		PartitionMapUpdates:   atomic.SwapInt64(&ns.PartitionMapUpdates, 0),
		NodeAdded:             atomic.SwapInt64(&ns.NodeAdded, 0),
		NodeRemoved:           atomic.SwapInt64(&ns.NodeRemoved, 0),
	}
}

// latest returns the latest values to be used in aggregation and then resets the values
func (ns *nodeStats) clone() nodeStats {
	return nodeStats{
		ConnectionsAttempts:   atomic.LoadInt64(&ns.ConnectionsAttempts),
		ConnectionsSuccessful: atomic.LoadInt64(&ns.ConnectionsSuccessful),
		ConnectionsFailed:     atomic.LoadInt64(&ns.ConnectionsFailed),
		ConnectionsPoolEmpty:  atomic.LoadInt64(&ns.ConnectionsPoolEmpty),
		ConnectionsOpen:       atomic.LoadInt64(&ns.ConnectionsOpen),
		ConnectionsClosed:     atomic.LoadInt64(&ns.ConnectionsClosed),
		TendsTotal:            atomic.LoadInt64(&ns.TendsTotal),
		TendsSuccessful:       atomic.LoadInt64(&ns.TendsSuccessful),
		TendsFailed:           atomic.LoadInt64(&ns.TendsFailed),
		PartitionMapUpdates:   atomic.LoadInt64(&ns.PartitionMapUpdates),
		NodeAdded:             atomic.LoadInt64(&ns.NodeAdded),
		NodeRemoved:           atomic.LoadInt64(&ns.NodeRemoved),
	}
}

func (ns *nodeStats) aggregate(newStats *nodeStats) {
	atomic.AddInt64(&ns.ConnectionsAttempts, newStats.ConnectionsAttempts)
	atomic.AddInt64(&ns.ConnectionsSuccessful, newStats.ConnectionsSuccessful)
	atomic.AddInt64(&ns.ConnectionsFailed, newStats.ConnectionsFailed)
	atomic.AddInt64(&ns.ConnectionsPoolEmpty, newStats.ConnectionsPoolEmpty)
	atomic.AddInt64(&ns.ConnectionsOpen, newStats.ConnectionsOpen)
	atomic.AddInt64(&ns.ConnectionsClosed, newStats.ConnectionsClosed)
	atomic.AddInt64(&ns.TendsTotal, newStats.TendsTotal)
	atomic.AddInt64(&ns.TendsSuccessful, newStats.TendsSuccessful)
	atomic.AddInt64(&ns.TendsFailed, newStats.TendsFailed)
	atomic.AddInt64(&ns.PartitionMapUpdates, newStats.PartitionMapUpdates)
	atomic.AddInt64(&ns.NodeAdded, newStats.NodeAdded)
	atomic.AddInt64(&ns.NodeRemoved, newStats.NodeRemoved)
}
