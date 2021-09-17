package raft

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/plugin/pb"
	"github.com/rboyer/safeio"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/atomic"

	"github.com/hashicorp/raft"
)

const (
	// boltSnapshotID is the stable ID for any boltDB snapshot. Keeping the ID
	// stable means there is only ever one bolt snapshot in the system
	boltSnapshotID = "bolt-snapshot"
	tmpSuffix      = ".tmp"
	snapPath       = "snapshots"
)

// BoltSnapshotStore implements the SnapshotStore interface and allows snapshots
// to be stored in BoltDB files on local disk. Since we always have an up to
// date FSM we use a special snapshot ID to indicate that the snapshot can be
// pulled from the BoltDB file that is currently backing the FSM. This allows us
// to provide just-in-time snapshots without doing incremental data dumps.
//
// When a snapshot is being installed on the node we will Create and Write data
// to it. This will cause the snapshot store to create a new BoltDB file and
// write the snapshot data to it. Then, we can simply rename the snapshot to the
// FSM's filename. This allows us to atomically install the snapshot and
// reduces the amount of disk i/o. Older snapshots are reaped on startup and
// before each subsequent snapshot write. This ensures we only have one snapshot
// on disk at a time.
type BoltSnapshotStore struct {
	// path is the directory in which to store file based snapshots
	path string

	// We hold a copy of the FSM so we can stream snapshots straight out of the
	// database.
	fsm *FSM

	logger log.Logger
}

// BoltSnapshotSink implements SnapshotSink optionally choosing to write to a
// file.
type BoltSnapshotSink struct {
	store  *BoltSnapshotStore
	logger log.Logger
	meta   raft.SnapshotMeta
	trans  raft.Transport

	// These fields will be used if we are writing a snapshot (vs. reading
	// one)
	written       atomic.Bool
	writer        io.WriteCloser
	writeError    error
	dir           string
	parentDir     string
	doneWritingCh chan struct{}

	l      sync.Mutex
	closed bool
}

// NewBoltSnapshotStore creates a new BoltSnapshotStore based
// on a base directory.
func NewBoltSnapshotStore(base string, logger log.Logger, fsm *FSM) (*BoltSnapshotStore, error) {
	if logger == nil {
		return nil, fmt.Errorf("no logger provided")
	}

	// Ensure our path exists
	path := filepath.Join(base, snapPath)
	if err := os.MkdirAll(path, 0o755); err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("snapshot path not accessible: %v", err)
	}

	// Setup the store
	store := &BoltSnapshotStore{
		logger: logger,
		fsm:    fsm,
		path:   path,
	}

	// Cleanup any old or failed snapshots on startup.
	if err := store.ReapSnapshots(); err != nil {
		return nil, err
	}

	return store, nil
}

// Create is used to start a new snapshot
func (f *BoltSnapshotStore) Create(version raft.SnapshotVersion, index, term uint64, configuration raft.Configuration, configurationIndex uint64, trans raft.Transport) (raft.SnapshotSink, error) {
	// We only support version 1 snapshots at this time.
	if version != 1 {
		return nil, fmt.Errorf("unsupported snapshot version %d", version)
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

	return sink, nil
}

// List returns available snapshots in the store. It only returns bolt
// snapshots. No snapshot will be returned if there are no indexes in the
// FSM.
func (f *BoltSnapshotStore) List() ([]*raft.SnapshotMeta, error) {
	meta, err := f.getMetaFromFSM()
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
func (f *BoltSnapshotStore) getMetaFromFSM() (*raft.SnapshotMeta, error) {
	latestIndex, latestConfig := f.fsm.LatestState()
	meta := &raft.SnapshotMeta{
		Version: 1,
		ID:      boltSnapshotID,
		Index:   latestIndex.Index,
		Term:    latestIndex.Term,
	}

	if latestConfig != nil {
		meta.ConfigurationIndex, meta.Configuration = protoConfigurationToRaftConfiguration(latestConfig)
	}

	return meta, nil
}

// Open takes a snapshot ID and returns a ReadCloser for that snapshot.
func (f *BoltSnapshotStore) Open(id string) (*raft.SnapshotMeta, io.ReadCloser, error) {
	if id == boltSnapshotID {
		return f.openFromFSM()
	}

	return f.openFromFile(id)
}

func (f *BoltSnapshotStore) openFromFSM() (*raft.SnapshotMeta, io.ReadCloser, error) {
	meta, err := f.getMetaFromFSM()
	if err != nil {
		return nil, nil, err
	}
	// If we don't have any data return an error
	if meta.Index == 0 {
		return nil, nil, errors.New("no snapshot data")
	}

	// Stream data out of the FSM to calculate the size
	readCloser, writeCloser := io.Pipe()
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
	metaReadCloser.Close()

	return meta, readCloser, nil
}

func (f *BoltSnapshotStore) getMetaFromDB(id string) (*raft.SnapshotMeta, error) {
	if len(id) == 0 {
		return nil, errors.New("can not open empty snapshot ID")
	}

	filename := filepath.Join(f.path, id, databaseFilename)
	boltDB, err := bolt.Open(filename, 0o600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	defer boltDB.Close()

	meta := &raft.SnapshotMeta{
		Version: 1,
		ID:      id,
	}

	err = boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(configBucketName)
		val := b.Get(latestIndexKey)
		if val != nil {
			var snapshotIndexes IndexValue
			err := proto.Unmarshal(val, &snapshotIndexes)
			if err != nil {
				return err
			}

			meta.Index = snapshotIndexes.Index
			meta.Term = snapshotIndexes.Term
		}

		// Read in our latest config and populate it inmemory
		val = b.Get(latestConfigKey)
		if val != nil {
			var config ConfigurationValue
			err := proto.Unmarshal(val, &config)
			if err != nil {
				return err
			}

			meta.ConfigurationIndex, meta.Configuration = protoConfigurationToRaftConfiguration(&config)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (f *BoltSnapshotStore) openFromFile(id string) (*raft.SnapshotMeta, io.ReadCloser, error) {
	meta, err := f.getMetaFromDB(id)
	if err != nil {
		return nil, nil, err
	}

	filename := filepath.Join(f.path, id, databaseFilename)
	installer := &boltSnapshotInstaller{
		meta:       meta,
		ReadCloser: ioutil.NopCloser(strings.NewReader(filename)),
		filename:   filename,
	}

	return meta, installer, nil
}

// ReapSnapshots reaps all snapshots.
func (f *BoltSnapshotStore) ReapSnapshots() error {
	snapshots, err := ioutil.ReadDir(f.path)
	switch {
	case err == nil:
	case os.IsNotExist(err):
		return nil
	default:
		f.logger.Error("failed to scan snapshot directory", "error", err)
		return err
	}

	for _, snap := range snapshots {
		// Ignore any files
		if !snap.IsDir() {
			continue
		}

		// Warn about temporary snapshots, this indicates a previously failed
		// snapshot attempt. We still want to clean these up.
		dirName := snap.Name()
		if strings.HasSuffix(dirName, tmpSuffix) {
			f.logger.Warn("found temporary snapshot", "name", dirName)
		}

		path := filepath.Join(f.path, dirName)
		f.logger.Info("reaping snapshot", "path", path)
		if err := os.RemoveAll(path); err != nil {
			f.logger.Error("failed to reap snapshot", "path", snap.Name(), "error", err)
			return err
		}
	}

	return nil
}

// ID returns the ID of the snapshot, can be used with Open()
// after the snapshot is finalized.
func (s *BoltSnapshotSink) ID() string {
	s.l.Lock()
	defer s.l.Unlock()

	return s.meta.ID
}

func (s *BoltSnapshotSink) writeBoltDBFile() error {
	// Create a new path
	name := snapshotName(s.meta.Term, s.meta.Index)
	path := filepath.Join(s.store.path, name+tmpSuffix)
	s.logger.Info("creating new snapshot", "path", path)

	// Make the directory
	if err := os.MkdirAll(path, 0o755); err != nil {
		s.logger.Error("failed to make snapshot directory", "error", err)
		return err
	}

	// Create the BoltDB file
	dbPath := filepath.Join(path, databaseFilename)
	boltDB, err := bolt.Open(dbPath, 0o600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	// Write the snapshot metadata
	if err := writeSnapshotMetaToDB(&s.meta, boltDB); err != nil {
		return err
	}

	// Set the snapshot ID to the generated name.
	s.meta.ID = name

	// Create the done channel
	s.doneWritingCh = make(chan struct{})

	// Store the directories so we can commit the changes on success or abort
	// them on failure.
	s.dir = path
	s.parentDir = s.store.path

	// Create a pipe so we pipe writes into the go routine below.
	reader, writer := io.Pipe()
	s.writer = writer

	// Start a go routine in charge of piping data from the snapshot's Write
	// call to the delimtedreader and the BoltDB file.
	go func() {
		defer close(s.doneWritingCh)
		defer boltDB.Close()

		// The delimted reader will parse full proto messages from the snapshot
		// data.
		protoReader := NewDelimitedReader(reader, math.MaxInt32)
		defer protoReader.Close()

		var done bool
		var keys int
		entry := new(pb.StorageEntry)
		for !done {
			err := boltDB.Update(func(tx *bolt.Tx) error {
				b, err := tx.CreateBucketIfNotExists(dataBucketName)
				if err != nil {
					return err
				}

				// Commit in batches of 50k. Bolt holds all the data in memory and
				// doesn't split the pages until commit so we do incremental writes.
				for i := 0; i < 50000; i++ {
					err := protoReader.ReadMsg(entry)
					if err != nil {
						if err == io.EOF {
							done = true
							return nil
						}
						return err
					}

					err = b.Put([]byte(entry.Key), entry.Value)
					if err != nil {
						return err
					}
					keys += 1
				}

				return nil
			})
			if err != nil {
				s.logger.Error("snapshot write: failed to write transaction", "error", err)
				s.writeError = err
				return
			}

			s.logger.Trace("snapshot write: writing keys", "num_written", keys)
		}
	}()

	return nil
}

// Write is used to append to the bolt file. The first call to write ensures we
// have the file created.
func (s *BoltSnapshotSink) Write(b []byte) (int, error) {
	s.l.Lock()
	defer s.l.Unlock()

	// If this is the first call to Write we need to setup the boltDB file and
	// kickoff the pipeline write
	if previouslyWritten := s.written.Swap(true); !previouslyWritten {
		// Reap any old snapshots
		if err := s.store.ReapSnapshots(); err != nil {
			return 0, err
		}

		if err := s.writeBoltDBFile(); err != nil {
			return 0, err
		}
	}

	return s.writer.Write(b)
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

	if s.writer != nil {
		s.writer.Close()
		<-s.doneWritingCh

		if s.writeError != nil {
			// If we encountered an error while writing then we should remove
			// the directory and return the error
			_ = os.RemoveAll(s.dir)
			return s.writeError
		}

		// Move the directory into place
		newPath := strings.TrimSuffix(s.dir, tmpSuffix)

		var err error
		if runtime.GOOS != "windows" {
			err = safeio.Rename(s.dir, newPath)
		} else {
			err = os.Rename(s.dir, newPath)
		}

		if err != nil {
			s.logger.Error("failed to move snapshot into place", "error", err)
			return err
		}
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

	if s.writer != nil {
		s.writer.Close()
		<-s.doneWritingCh

		// Attempt to remove all artifacts
		return os.RemoveAll(s.dir)
	}

	return nil
}

type boltSnapshotInstaller struct {
	io.ReadCloser
	meta     *raft.SnapshotMeta
	filename string
}

func (i *boltSnapshotInstaller) Filename() string {
	return i.filename
}

func (i *boltSnapshotInstaller) Metadata() *raft.SnapshotMeta {
	return i.meta
}

func (i *boltSnapshotInstaller) Install(filename string) error {
	if len(i.filename) == 0 {
		return errors.New("snapshot filename empty")
	}

	if len(filename) == 0 {
		return errors.New("fsm filename empty")
	}

	// Rename the snapshot to the FSM location
	if runtime.GOOS != "windows" {
		return safeio.Rename(i.filename, filename)
	} else {
		return os.Rename(i.filename, filename)
	}
}

// snapshotName generates a name for the snapshot.
func snapshotName(term, index uint64) string {
	now := time.Now()
	msec := now.UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%d-%d-%d", term, index, msec)
}
