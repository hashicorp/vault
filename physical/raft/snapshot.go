package raft

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/raft"
)

const (
	// boltSnapshotID is the stable ID for any boltDB snapshot. Keeping the ID
	// stable means there is only ever one bolt snapshot in the system
	boltSnapshotID = "bolt-snapshot"
)

// BoltSnapshotStore implements the SnapshotStore interface and allows
// snapshots to be made on the local disk. The main difference between this
// store and the file store is we make the distinction between snapshots that
// have been written by the FSM and by internal Raft operations. The former are
// treated as noop snapshots on Persist and are read in full from the FSM on
// Open. The latter are treated like normal file snapshots and are able to be
// opened and applied as usual.
type BoltSnapshotStore struct {
	// path is the directory in which to store file based snapshots
	path string
	// retain is the number of file based snapshots to keep
	retain int

	// We hold a copy of the FSM so we can stream snapshots straight out of the
	// database.
	fsm *FSM

	// fileSnapStore is used to fall back to file snapshots when the data is
	// being written from the raft library. This currently only happens on a
	// follower during a snapshot install RPC.
	fileSnapStore *raft.FileSnapshotStore
	logger        log.Logger
}

// BoltSnapshotSink implements SnapshotSink optionally choosing to write to a
// file.
type BoltSnapshotSink struct {
	store  *BoltSnapshotStore
	logger log.Logger
	meta   raft.SnapshotMeta
	trans  raft.Transport

	fileSink raft.SnapshotSink
	l        sync.Mutex
	closed   bool
}

// NewBoltSnapshotStore creates a new BoltSnapshotStore based
// on a base directory. The `retain` parameter controls how many
// snapshots are retained. Must be at least 1.
func NewBoltSnapshotStore(base string, retain int, logger log.Logger, fsm *FSM) (*BoltSnapshotStore, error) {
	if retain < 1 {
		return nil, fmt.Errorf("must retain at least one snapshot")
	}
	if logger == nil {
		return nil, fmt.Errorf("no logger provided")
	}

	fileStore, err := raft.NewFileSnapshotStore(base, retain, nil)
	if err != nil {
		return nil, err
	}

	// Setup the store
	store := &BoltSnapshotStore{
		logger:        logger,
		fsm:           fsm,
		fileSnapStore: fileStore,
	}

	{
		// TODO: I think this needs to be done before every NewRaft and
		// RecoverCluster call. Not just on Factory method.

		// Here we delete all the existing file based snapshots. This is necessary
		// because we do not issue a restore on NewRaft. If a previous file snapshot
		// had failed to apply we will be incorrectly setting the indexes. It's
		// safer to simply delete all file snapshots on startup and rely on Raft to
		// reconcile the FSM state.
		if err := store.ReapSnapshots(); err != nil {
			return nil, err
		}
	}

	return store, nil
}

// Create is used to start a new snapshot
func (f *BoltSnapshotStore) Create(version raft.SnapshotVersion, index, term uint64,
	configuration raft.Configuration, configurationIndex uint64, trans raft.Transport) (raft.SnapshotSink, error) {
	// We only support version 1 snapshots at this time.
	if version != 1 {
		return nil, fmt.Errorf("unsupported snapshot version %d", version)
	}

	// We are processing a snapshot, fastforward the index, term, and
	// configuration to the latest seen by the raft system. This could include
	// log indexes for operation types that are never sent to the FSM.
	if err := f.fsm.witnessSnapshot(index, term, configurationIndex, configuration); err != nil {
		return nil, err
	}

	// Create the sink
	sink := &BoltSnapshotSink{
		store:  f,
		logger: f.logger,
		meta: raft.SnapshotMeta{
			Version:            version,
			ID:                 boltSnapshotID,
			Index:              index,
			Term:               term,
			Configuration:      configuration,
			ConfigurationIndex: configurationIndex,
		},
		trans: trans,
	}

	// Done
	return sink, nil
}

// List returns available snapshots in the store. It only returns bolt
// snapshots. No snapshot will be returned if there are no indexes in the
// FSM.
func (f *BoltSnapshotStore) List() ([]*raft.SnapshotMeta, error) {
	meta, err := f.getBoltSnapshotMeta()
	if err != nil {
		return nil, err
	}

	// If we haven't seen any data yet do not return a snapshot
	if meta.Index == 0 {
		return nil, nil
	}

	return []*raft.SnapshotMeta{meta}, nil
}

// getBoltSnapshotMeta returns the fsm's latest state and configuration.
func (f *BoltSnapshotStore) getBoltSnapshotMeta() (*raft.SnapshotMeta, error) {
	latestIndex, latestConfig := f.fsm.LatestState()
	meta := &raft.SnapshotMeta{
		Version: 1,
		ID:      boltSnapshotID,
		Index:   latestIndex.Index,
		Term:    latestIndex.Term,
	}

	if latestConfig != nil {
		index, configuration := protoConfigurationToRaftConfiguration(latestConfig)
		meta.Configuration = configuration
		meta.ConfigurationIndex = index
	}

	return meta, nil
}

// Open takes a snapshot ID and returns a ReadCloser for that snapshot.
func (f *BoltSnapshotStore) Open(id string) (*raft.SnapshotMeta, io.ReadCloser, error) {
	var readCloser io.ReadCloser
	var meta *raft.SnapshotMeta
	switch id {
	case boltSnapshotID:

		var err error
		meta, err = f.getBoltSnapshotMeta()
		if err != nil {
			return nil, nil, err
		}
		// If we don't have any data return an error
		if meta.Index == 0 {
			return nil, nil, errors.New("no snapshot data")
		}

		// Stream data out of the FSM to calculate the size
		var writeCloser *io.PipeWriter
		readCloser, writeCloser = io.Pipe()
		metaReadCloser, metaWriteCloser := io.Pipe()
		go func() {
			f.fsm.writeTo(context.Background(), metaWriteCloser, writeCloser)
		}()

		// Compute the size
		n, err := io.Copy(ioutil.Discard, metaReadCloser)
		if err != nil {
			f.logger.Error("failed to read state file", "error", err)
			metaReadCloser.Close()
			readCloser.Close()
			return nil, nil, err
		}

		meta.Size = n

	default:
		var err error
		meta, readCloser, err = f.fileSnapStore.Open(id)
		if err != nil {
			return nil, nil, err
		}
	}

	return meta, readCloser, nil
}

// ReapSnapshots reaps any snapshots beyond the retain count.
func (f *BoltSnapshotStore) ReapSnapshots() error {
	return f.fileSnapStore.ReapSnapshots()
}

// ID returns the ID of the snapshot, can be used with Open()
// after the snapshot is finalized.
func (s *BoltSnapshotSink) ID() string {
	s.l.Lock()
	defer s.l.Unlock()

	if s.fileSink != nil {
		return s.fileSink.ID()
	}

	return s.meta.ID
}

// Write is used to append to the state file. We write to the
// buffered IO object to reduce the amount of context switches.
func (s *BoltSnapshotSink) Write(b []byte) (int, error) {
	s.l.Lock()
	defer s.l.Unlock()

	// If someone is writing to this sink then we need to create a file sink to
	// capture the data. This currently only happens when a follower is being
	// sent a snapshot.
	if s.fileSink == nil {
		fileSink, err := s.store.fileSnapStore.Create(s.meta.Version, s.meta.Index, s.meta.Term, s.meta.Configuration, s.meta.ConfigurationIndex, s.trans)
		if err != nil {
			return 0, err
		}
		s.fileSink = fileSink
	}

	return s.fileSink.Write(b)
}

// Close is used to indicate a successful end.
func (s *BoltSnapshotSink) Close() error {
	s.l.Lock()
	defer s.l.Unlock()

	// Make sure close is idempotent
	if s.closed {
		return nil
	}
	s.closed = true

	if s.fileSink != nil {
		return s.fileSink.Close()
	}

	return nil
}

// Cancel is used to indicate an unsuccessful end.
func (s *BoltSnapshotSink) Cancel() error {
	s.l.Lock()
	defer s.l.Unlock()

	// Make sure close is idempotent
	if s.closed {
		return nil
	}
	s.closed = true

	if s.fileSink != nil {
		return s.fileSink.Cancel()
	}

	return nil
}
