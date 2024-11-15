// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build batchtest
// +build batchtest

package raft

func init() {
	userSnapshotErrorsOnNoData = false
}

// ApplyBatch enables MockFSM to satisfy the BatchingFSM interface. This
// function is gated by the batchtest build flag.
//
// NOTE: This is exposed for middleware testing purposes and is not a stable API
func (m *MockFSM) ApplyBatch(logs []*Log) []interface{} {
	m.Lock()
	defer m.Unlock()

	ret := make([]interface{}, len(logs))
	for i, log := range logs {
		switch log.Type {
		case LogCommand:
			m.logs = append(m.logs, log.Data)
			ret[i] = len(m.logs)
		default:
			ret[i] = nil
		}
	}

	return ret
}
