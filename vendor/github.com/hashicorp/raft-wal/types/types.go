// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"errors"

	"github.com/hashicorp/raft"
)

var (
	// ErrNotFound is our own version of raft's not found error. It's important
	// it's exactly the same because the raft lib checks for equality with it's
	// own type as a crucial part of replication processing (detecting end of logs
	// and that a snapshot is needed for a follower).
	ErrNotFound = raft.ErrLogNotFound
	ErrCorrupt  = errors.New("WAL is corrupt")
	ErrSealed   = errors.New("segment is sealed")
	ErrClosed   = errors.New("closed")
)

// LogEntry represents an entry that has already been encoded.
type LogEntry struct {
	Index uint64
	Data  []byte
}
