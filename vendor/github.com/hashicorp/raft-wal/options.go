// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package wal

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft-wal/fs"
	"github.com/hashicorp/raft-wal/metadb"
	"github.com/hashicorp/raft-wal/metrics"
	"github.com/hashicorp/raft-wal/segment"
	"github.com/hashicorp/raft-wal/types"
)

// WithCodec is an option that allows a custom Codec to be provided to the WAL.
// If not used the default Codec is used.
func WithCodec(c Codec) walOpt {
	return func(w *WAL) {
		w.codec = c
	}
}

// WithMetaStore is an option that allows a custom MetaStore to be provided to
// the WAL. If not used the default MetaStore is used.
func WithMetaStore(db types.MetaStore) walOpt {
	return func(w *WAL) {
		w.metaDB = db
	}
}

// WithSegmentFiler is an option that allows a custom SegmentFiler (and hence
// Segment Reader/Writer implementation) to be provided to the WAL. If not used
// the default SegmentFiler is used.
func WithSegmentFiler(sf types.SegmentFiler) walOpt {
	return func(w *WAL) {
		w.sf = sf
	}
}

// WithLogger is an option that allows a custom logger to be used.
func WithLogger(logger hclog.Logger) walOpt {
	return func(w *WAL) {
		w.log = logger
	}
}

// WithSegmentSize is an option that allows a custom segmentSize to be set.
func WithSegmentSize(size int) walOpt {
	return func(w *WAL) {
		w.segmentSize = size
	}
}

// WithMetricsCollector is an option that allows a custom segmentSize to be set.
func WithMetricsCollector(c metrics.Collector) walOpt {
	return func(w *WAL) {
		w.metrics = c
	}
}

func (w *WAL) applyDefaultsAndValidate() error {
	// Check if an external codec has been used that it's not using a reserved ID.
	if w.codec != nil && w.codec.ID() < FirstExternalCodecID {
		return fmt.Errorf("codec is using a reserved ID (below %d)", FirstExternalCodecID)
	}

	// Defaults
	if w.log == nil {
		w.log = hclog.Default().Named("wal")
	}
	if w.codec == nil {
		w.codec = &BinaryCodec{}
	}
	if w.sf == nil {
		// These are not actually swappable via options right now but we override
		// them in tests. Only load the default implementations if they are not set.
		vfs := fs.New()
		w.sf = segment.NewFiler(w.dir, vfs)
	}
	if w.metrics == nil {
		w.metrics = &metrics.NoOpCollector{}
	}
	if w.metaDB == nil {
		w.metaDB = &metadb.BoltMetaDB{}
	}
	if w.segmentSize == 0 {
		w.segmentSize = DefaultSegmentSize
	}
	return nil
}
