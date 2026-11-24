// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package raft

import "github.com/hashicorp/go-hclog"

func (b *RaftBackend) entrySizeLimitForPath(path string) uint64 {
	return b.maxEntrySize
}

func emitEntWarning(logger hclog.Logger, field string) {
	logger.Warn("configuration for a Vault Enterprise feature has been ignored", "field", field)
}
