// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package raft

import (
	"bytes"
	"context"
	"fmt"
	"hash/crc64"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
)

type idAddr struct {
	id string
}

func (a *idAddr) Network() string { return "inmem" }
func (a *idAddr) String() string  { return a.id }

func addPeer(t *testing.T, leader, follower *RaftBackend) {
	t.Helper()
	if err := leader.AddPeer(context.Background(), follower.NodeID(), follower.NodeID()); err != nil {
		t.Fatal(err)
	}

	peers, err := leader.Peers(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	err = follower.Bootstrap(peers)
	if err != nil {
		t.Fatal(err)
	}

	err = follower.SetupCluster(context.Background(), SetupOpts{})
	if err != nil {
		t.Fatal(err)
	}

	leader.raftTransport.(*raft.InmemTransport).Connect(raft.ServerAddress(follower.NodeID()), follower.raftTransport)
	follower.raftTransport.(*raft.InmemTransport).Connect(raft.ServerAddress(leader.NodeID()), leader.raftTransport)
}

func TestRaft_Snapshot_Loading(t *testing.T) {
	raft, dir := GetRaft(t, true, false)
	defer os.RemoveAll(dir)

	// Write some data
	for i := 0; i < 1000; i++ {
		err := raft.Put(context.Background(), &physical.Entry{
			Key:   fmt.Sprintf("key-%d", i),
			Value: []byte(fmt.Sprintf("value-%d", i)),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	readCloser, writeCloser := io.Pipe()
	metaReadCloser, metaWriteCloser := io.Pipe()

	go func() {
		raft.fsm.writeTo(context.Background(), metaWriteCloser, writeCloser)
	}()

	// Create a CRC64 hash
	stateHash := crc64.New(crc64.MakeTable(crc64.ECMA))

	// Compute the hash
	size1, err := io.Copy(stateHash, metaReadCloser)
	if err != nil {
		t.Fatal(err)
	}

	computed1 := stateHash.Sum(nil)

	// Create a CRC64 hash
	stateHash = crc64.New(crc64.MakeTable(crc64.ECMA))

	// Compute the hash
	size2, err := io.Copy(stateHash, readCloser)
	if err != nil {
		t.Fatal(err)
	}

	computed2 := stateHash.Sum(nil)

	if size1 != size2 {
		t.Fatal("sizes did not match")
	}

	if !bytes.Equal(computed1, computed2) {
		t.Fatal("hashes did not match")
	}

	snapFuture := raft.raft.Snapshot()
	if err := snapFuture.Error(); err != nil {
		t.Fatal(err)
	}

	meta, reader, err := snapFuture.Open()
	if err != nil {
		t.Fatal(err)
	}
	if meta.Size != size1 {
		t.Fatal("meta size did not match expected")
	}

	// Create a CRC64 hash
	stateHash = crc64.New(crc64.MakeTable(crc64.ECMA))

	// Compute the hash
	size3, err := io.Copy(stateHash, reader)
	if err != nil {
		t.Fatal(err)
	}

	computed3 := stateHash.Sum(nil)
	if size1 != size3 {
		t.Fatal("sizes did not match")
	}

	if !bytes.Equal(computed1, computed3) {
		t.Fatal("hashes did not match")
	}
}

func TestRaft_Snapshot_Index(t *testing.T) {
	raft, dir := GetRaft(t, true, false)
	defer os.RemoveAll(dir)

	err := raft.Put(context.Background(), &physical.Entry{
		Key:   "key",
		Value: []byte("value"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get index
	index, _ := raft.fsm.LatestState()
	if index.Term != 2 {
		t.Fatalf("unexpected term, got %d expected 2", index.Term)
	}
	if index.Index != 3 {
		t.Fatalf("unexpected index, got %d expected 3", index.Term)
	}

	// Write some data
	for i := 0; i < 100; i++ {
		err := raft.Put(context.Background(), &physical.Entry{
			Key:   fmt.Sprintf("key-%d", i),
			Value: []byte(fmt.Sprintf("value-%d", i)),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Get index
	index, _ = raft.fsm.LatestState()
	if index.Term != 2 {
		t.Fatalf("unexpected term, got %d expected 2", index.Term)
	}
	if index.Index != 103 {
		t.Fatalf("unexpected index, got %d expected 103", index.Term)
	}

	// Take a snapshot
	snapFuture := raft.raft.Snapshot()
	if err := snapFuture.Error(); err != nil {
		t.Fatal(err)
	}

	meta, reader, err := snapFuture.Open()
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(ioutil.Discard, reader)

	if meta.Index != index.Index {
		t.Fatalf("indexes did not match, got %d expected %d", meta.Index, index.Index)
	}
	if meta.Term != index.Term {
		t.Fatalf("term did not match, got %d expected %d", meta.Term, index.Term)
	}

	// Write some more data
	for i := 0; i < 100; i++ {
		err := raft.Put(context.Background(), &physical.Entry{
			Key:   fmt.Sprintf("key-%d", i),
			Value: []byte(fmt.Sprintf("value-%d", i)),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Open the same snapshot again
	meta, reader, err = raft.snapStore.Open(meta.ID)
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(ioutil.Discard, reader)

	// Make sure the meta data has updated to the new values
	if meta.Index != 203 {
		t.Fatalf("unexpected snapshot index %d", meta.Index)
	}
	if meta.Term != 2 {
		t.Fatalf("unexpected snapshot term %d", meta.Term)
	}
}

func TestRaft_Snapshot_Peers(t *testing.T) {
	raft1, dir := GetRaft(t, true, false)
	raft2, dir2 := GetRaft(t, false, false)
	raft3, dir3 := GetRaft(t, false, false)
	defer os.RemoveAll(dir)
	defer os.RemoveAll(dir2)
	defer os.RemoveAll(dir3)

	// Write some data
	for i := 0; i < 1000; i++ {
		err := raft1.Put(context.Background(), &physical.Entry{
			Key:   fmt.Sprintf("key-%d", i),
			Value: []byte(fmt.Sprintf("value-%d", i)),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Force a snapshot
	snapFuture := raft1.raft.Snapshot()
	if err := snapFuture.Error(); err != nil {
		t.Fatal(err)
	}

	commitIdx := raft1.CommittedIndex()

	// Add raft2 to the cluster
	addPeer(t, raft1, raft2)

	ensureCommitApplied(t, commitIdx, raft2)

	// Make sure the snapshot was applied correctly on the follower
	if err := compareDBs(t, raft1.fsm.getDB(), raft2.fsm.getDB(), false); err != nil {
		t.Fatal(err)
	}

	// Write some more data
	for i := 1000; i < 2000; i++ {
		err := raft1.Put(context.Background(), &physical.Entry{
			Key:   fmt.Sprintf("key-%d", i),
			Value: []byte(fmt.Sprintf("value-%d", i)),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	snapFuture = raft1.raft.Snapshot()
	if err := snapFuture.Error(); err != nil {
		t.Fatal(err)
	}

	commitIdx = raft1.CommittedIndex()

	// Add raft3 to the cluster
	addPeer(t, raft1, raft3)

	ensureCommitApplied(t, commitIdx, raft2)
	ensureCommitApplied(t, commitIdx, raft3)

	// Make sure all stores are the same
	compareFSMs(t, raft1.fsm, raft2.fsm)
	compareFSMs(t, raft1.fsm, raft3.fsm)
}

func ensureCommitApplied(t *testing.T, leaderCommitIdx uint64, backend *RaftBackend) {
	t.Helper()

	timeout := time.Now().Add(10 * time.Second)
	for {
		if time.Now().After(timeout) {
			t.Fatal("timeout reached while verifying applied index on raft backend")
		}

		if backend.AppliedIndex() >= leaderCommitIdx {
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func TestRaft_Snapshot_Restart(t *testing.T) {
	raft1, dir := GetRaft(t, true, false)
	defer os.RemoveAll(dir)
	raft2, dir2 := GetRaft(t, false, false)
	defer os.RemoveAll(dir2)

	// Write some data
	for i := 0; i < 100; i++ {
		err := raft1.Put(context.Background(), &physical.Entry{
			Key:   fmt.Sprintf("key-%d", i),
			Value: []byte(fmt.Sprintf("value-%d", i)),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Take a snapshot
	snapFuture := raft1.raft.Snapshot()
	if err := snapFuture.Error(); err != nil {
		t.Fatal(err)
	}
	// Advance FSM's index past configuration change
	raft1.Put(context.Background(), &physical.Entry{
		Key:   "key",
		Value: []byte("value"),
	})

	// Add raft2 to the cluster
	addPeer(t, raft1, raft2)

	time.Sleep(2 * time.Second)

	peers, err := raft2.Peers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(peers) != 2 {
		t.Fatal(peers)
	}

	// Finalize raft1
	if err := raft1.TeardownCluster(nil); err != nil {
		t.Fatal(err)
	}

	// Start Raft
	err = raft1.SetupCluster(context.Background(), SetupOpts{})
	if err != nil {
		t.Fatal(err)
	}

	peers, err = raft1.Peers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(peers) != 2 {
		t.Fatal(peers)
	}

	compareFSMs(t, raft1.fsm, raft2.fsm)
}

/*
func TestRaft_Snapshot_ErrorRecovery(t *testing.T) {
	raft1, dir := GetRaft(t, true, false)
	raft2, dir2 := GetRaft(t, false, false)
	raft3, dir3 := GetRaft(t, false, false)
	defer os.RemoveAll(dir)
	defer os.RemoveAll(dir2)
	defer os.RemoveAll(dir3)

	// Add raft2 to the cluster
	addPeer(t, raft1, raft2)

	// Write some data
	for i := 0; i < 100; i++ {
		err := raft1.Put(context.Background(), &physical.Entry{
			Key:   fmt.Sprintf("key-%d", i),
			Value: []byte(fmt.Sprintf("value-%d", i)),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Take a snapshot on each node to ensure we no longer have older logs
	snapFuture := raft1.raft.Snapshot()
	if err := snapFuture.Error(); err != nil {
		t.Fatal(err)
	}

	stepDownLeader(t, raft1)
	leader := waitForLeader(t, raft1, raft2)

	snapFuture = leader.raft.Snapshot()
	if err := snapFuture.Error(); err != nil {
		t.Fatal(err)
	}

	// Advance FSM's index past snapshot index
	leader.Put(context.Background(), &physical.Entry{
		Key:   "key",
		Value: []byte("value"),
	})

	// Error on snapshot restore
	raft3.fsm.testSnapshotRestoreError = true

	// Add raft3 to the cluster
	addPeer(t, leader, raft3)

	time.Sleep(2 * time.Second)

	// Restart the failing node to make sure fresh state does not have invalid
	// values.
	if err := raft3.TeardownCluster(nil); err != nil {
		t.Fatal(err)
	}

	// Ensure the databases are not equal
	if err := compareFSMsWithErr(t, leader.fsm, raft3.fsm); err == nil {
		t.Fatal("nil error")
	}

	// Remove error and make sure we can reconcile state
	raft3.fsm.testSnapshotRestoreError = false

	// Step down leader node
	stepDownLeader(t, leader)
	leader = waitForLeader(t, raft1, raft2)

	// Start Raft3
	if err := raft3.SetupCluster(context.Background(), SetupOpts{}); err != nil {
		t.Fatal(err)
	}

	connectPeers(raft1, raft2, raft3)
	waitForLeader(t, raft1, raft2)

	time.Sleep(5 * time.Second)

	// Make sure state gets re-replicated.
	compareFSMs(t, raft1.fsm, raft3.fsm)
}*/

func TestRaft_Snapshot_Take_Restore(t *testing.T) {
	raft1, dir := GetRaft(t, true, false)
	defer os.RemoveAll(dir)
	raft2, dir2 := GetRaft(t, false, false)
	defer os.RemoveAll(dir2)

	addPeer(t, raft1, raft2)

	// Write some data
	for i := 0; i < 100; i++ {
		err := raft1.Put(context.Background(), &physical.Entry{
			Key:   fmt.Sprintf("key-%d", i),
			Value: []byte(fmt.Sprintf("value-%d", i)),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	recorder := httptest.NewRecorder()
	snap := logical.NewHTTPResponseWriter(recorder)

	err := raft1.Snapshot(snap, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Write some more data
	for i := 100; i < 200; i++ {
		err := raft1.Put(context.Background(), &physical.Entry{
			Key:   fmt.Sprintf("key-%d", i),
			Value: []byte(fmt.Sprintf("value-%d", i)),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	snapFile, cleanup, metadata, err := raft1.WriteSnapshotToTemp(ioutil.NopCloser(recorder.Body), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	err = raft1.RestoreSnapshot(context.Background(), metadata, snapFile)
	if err != nil {
		t.Fatal(err)
	}

	// make sure we don't have the second batch of writes
	for i := 100; i < 200; i++ {
		{
			value, err := raft1.Get(context.Background(), fmt.Sprintf("key-%d", i))
			if err != nil {
				t.Fatal(err)
			}
			if value != nil {
				t.Fatal("didn't remove data")
			}
		}
		{
			value, err := raft2.Get(context.Background(), fmt.Sprintf("key-%d", i))
			if err != nil {
				t.Fatal(err)
			}
			if value != nil {
				t.Fatal("didn't remove data")
			}
		}
	}

	time.Sleep(10 * time.Second)
	compareFSMs(t, raft1.fsm, raft2.fsm)
}

func TestBoltSnapshotStore_CreateSnapshotMissingParentDir(t *testing.T) {
	parent, err := ioutil.TempDir("", "raft")
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	defer os.RemoveAll(parent)

	dir, err := ioutil.TempDir(parent, "raft")
	if err != nil {
		t.Fatalf("err: %v ", err)
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "raft",
		Level: hclog.Trace,
	})

	snap, err := NewBoltSnapshotStore(dir, logger, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	os.RemoveAll(parent)
	_, trans := raft.NewInmemTransport(raft.NewInmemAddr())
	sink, err := snap.Create(raft.SnapshotVersionMax, 10, 3, raft.Configuration{}, 0, trans)
	if err != nil {
		t.Fatal(err)
	}
	defer sink.Cancel()

	_, err = sink.Write([]byte("test"))
	if err != nil {
		t.Fatalf("should not fail when using non existing parent: %s", err)
	}

	// Ensure the snapshot file exists
	_, err = os.Stat(filepath.Join(snap.path, sink.ID()+tmpSuffix, databaseFilename))
	if err != nil {
		t.Fatal(err)
	}
}

func TestBoltSnapshotStore_Listing(t *testing.T) {
	// Create a test dir
	parent, err := ioutil.TempDir("", "raft")
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	defer os.RemoveAll(parent)

	dir, err := ioutil.TempDir(parent, "raft")
	if err != nil {
		t.Fatalf("err: %v ", err)
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "raft",
		Level: hclog.Trace,
	})

	fsm, err := NewFSM(parent, "", logger)
	if err != nil {
		t.Fatal(err)
	}

	snap, err := NewBoltSnapshotStore(dir, logger, fsm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// FSM has no data, should have empty snapshot list
	snaps, err := snap.List()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(snaps) != 0 {
		t.Fatalf("expect 0 snapshots: %v", snaps)
	}

	// Move the fsm forward
	err = fsm.witnessSnapshot(&raft.SnapshotMeta{
		Index:              100,
		Term:               20,
		Configuration:      raft.Configuration{},
		ConfigurationIndex: 0,
	})
	if err != nil {
		t.Fatal(err)
	}

	snaps, err = snap.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(snaps) != 1 {
		t.Fatalf("expect 1 snapshots: %v", snaps)
	}

	if snaps[0].Index != 100 || snaps[0].Term != 20 {
		t.Fatalf("bad snapshot: %+v", snaps[0])
	}

	if snaps[0].ID != boltSnapshotID {
		t.Fatalf("bad snapshot: %+v", snaps[0])
	}
}

func TestBoltSnapshotStore_CreateInstallSnapshot(t *testing.T) {
	// Create a test dir
	parent, err := ioutil.TempDir("", "raft")
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	defer os.RemoveAll(parent)

	dir, err := ioutil.TempDir(parent, "raft")
	if err != nil {
		t.Fatalf("err: %v ", err)
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "raft",
		Level: hclog.Trace,
	})

	fsm, err := NewFSM(parent, "", logger)
	if err != nil {
		t.Fatal(err)
	}
	defer fsm.Close()

	snap, err := NewBoltSnapshotStore(dir, logger, fsm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check no snapshots
	snaps, err := snap.List()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(snaps) != 0 {
		t.Fatalf("did not expect any snapshots: %v", snaps)
	}

	// Create a new sink
	var configuration raft.Configuration
	configuration.Servers = append(configuration.Servers, raft.Server{
		Suffrage: raft.Voter,
		ID:       raft.ServerID("my id"),
		Address:  raft.ServerAddress("over here"),
	})
	_, trans := raft.NewInmemTransport(raft.NewInmemAddr())
	sink, err := snap.Create(raft.SnapshotVersionMax, 10, 3, configuration, 2, trans)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	protoWriter := NewDelimitedWriter(sink)

	err = fsm.Put(context.Background(), &physical.Entry{
		Key:   "test-key",
		Value: []byte("test-value"),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = fsm.Put(context.Background(), &physical.Entry{
		Key:   "test-key1",
		Value: []byte("test-value1"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Write to the sink
	err = protoWriter.WriteMsg(&pb.StorageEntry{
		Key:   "test-key",
		Value: []byte("test-value"),
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	err = protoWriter.WriteMsg(&pb.StorageEntry{
		Key:   "test-key1",
		Value: []byte("test-value1"),
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Done!
	err = sink.Close()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Read the snapshot
	meta, r, err := snap.Open(sink.ID())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the latest
	if meta.Index != 10 {
		t.Fatalf("bad snapshot: %+v", meta)
	}
	if meta.Term != 3 {
		t.Fatalf("bad snapshot: %+v", meta)
	}
	if !reflect.DeepEqual(meta.Configuration, configuration) {
		t.Fatalf("bad snapshot: %+v", meta)
	}
	if meta.ConfigurationIndex != 2 {
		t.Fatalf("bad snapshot: %+v", meta)
	}

	installer, ok := r.(*boltSnapshotInstaller)
	if !ok {
		t.Fatal("expected snapshot installer object")
	}

	newFSM, err := NewFSM(filepath.Dir(installer.Filename()), "", logger)
	if err != nil {
		t.Fatal(err)
	}

	err = compareDBs(t, fsm.getDB(), newFSM.getDB(), true)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure config data is different
	err = compareDBs(t, fsm.getDB(), newFSM.getDB(), false)
	if err == nil {
		t.Fatal("expected error")
	}

	if err := newFSM.Close(); err != nil {
		t.Fatal(err)
	}

	err = fsm.Restore(installer)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 2; i++ {
		latestIndex, latestConfigRaw := fsm.LatestState()
		latestConfigIndex, latestConfig := protoConfigurationToRaftConfiguration(latestConfigRaw)
		if latestIndex.Index != 10 {
			t.Fatalf("bad install: %+v", latestIndex)
		}
		if latestIndex.Term != 3 {
			t.Fatalf("bad install: %+v", latestIndex)
		}
		if !reflect.DeepEqual(latestConfig, configuration) {
			t.Fatalf("bad install: %+v", latestConfig)
		}
		if latestConfigIndex != 2 {
			t.Fatalf("bad install: %+v", latestConfigIndex)
		}

		v, err := fsm.Get(context.Background(), "test-key")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(v.Value, []byte("test-value")) {
			t.Fatalf("bad: %+v", v)
		}

		v, err = fsm.Get(context.Background(), "test-key1")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(v.Value, []byte("test-value1")) {
			t.Fatalf("bad: %+v", v)
		}

		// Close/Reopen the db and make sure we still match
		fsm.Close()
		fsm, err = NewFSM(parent, "", logger)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestBoltSnapshotStore_CancelSnapshot(t *testing.T) {
	// Create a test dir
	dir, err := ioutil.TempDir("", "raft")
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	defer os.RemoveAll(dir)

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "raft",
		Level: hclog.Trace,
	})

	snap, err := NewBoltSnapshotStore(dir, logger, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	_, trans := raft.NewInmemTransport(raft.NewInmemAddr())
	sink, err := snap.Create(raft.SnapshotVersionMax, 10, 3, raft.Configuration{}, 0, trans)
	if err != nil {
		t.Fatal(err)
	}
	_, err = sink.Write([]byte("test"))
	if err != nil {
		t.Fatalf("should not fail when using non existing parent: %s", err)
	}

	// Ensure the snapshot file exists
	_, err = os.Stat(filepath.Join(snap.path, sink.ID()+tmpSuffix, databaseFilename))
	if err != nil {
		t.Fatal(err)
	}

	// Cancel the snapshot! Should delete
	err = sink.Cancel()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Ensure the snapshot file does not exist
	_, err = os.Stat(filepath.Join(snap.path, sink.ID()+tmpSuffix, databaseFilename))
	if !os.IsNotExist(err) {
		t.Fatal(err)
	}

	// Make sure future writes fail
	_, err = sink.Write([]byte("test"))
	if err == nil {
		t.Fatal("expected write to fail")
	}
}

func TestBoltSnapshotStore_BadPerm(t *testing.T) {
	var err error
	if runtime.GOOS == "windows" {
		t.Skip("skipping file permission test on windows")
	}

	// Create a temp dir
	var dir1 string
	dir1, err = ioutil.TempDir("", "raft")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(dir1)

	// Create a sub dir and remove all permissions
	var dir2 string
	dir2, err = ioutil.TempDir(dir1, "badperm")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if err = os.Chmod(dir2, 0o00); err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Chmod(dir2, 777) // Set perms back for delete

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "raft",
		Level: hclog.Trace,
	})

	_, err = NewBoltSnapshotStore(dir2, logger, nil)
	if err == nil {
		t.Fatalf("should fail to use dir with bad perms")
	}
}

func TestBoltSnapshotStore_CloseFailure(t *testing.T) {
	// Create a test dir
	dir, err := ioutil.TempDir("", "raft")
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	defer os.RemoveAll(dir)

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "raft",
		Level: hclog.Trace,
	})

	snap, err := NewBoltSnapshotStore(dir, logger, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	_, trans := raft.NewInmemTransport(raft.NewInmemAddr())
	sink, err := snap.Create(raft.SnapshotVersionMax, 10, 3, raft.Configuration{}, 0, trans)
	if err != nil {
		t.Fatal(err)
	}

	// This should stash an error value
	_, err = sink.Write([]byte("test"))
	if err != nil {
		t.Fatalf("should not fail when using non existing parent: %s", err)
	}

	// Cancel the snapshot! Should delete
	err = sink.Close()
	if err == nil {
		t.Fatalf("expected error")
	}

	// Ensure the snapshot file does not exist
	_, err = os.Stat(filepath.Join(snap.path, sink.ID()+tmpSuffix, databaseFilename))
	if !os.IsNotExist(err) {
		t.Fatal(err)
	}

	// Make sure future writes fail
	_, err = sink.Write([]byte("test"))
	if err == nil {
		t.Fatal("expected write to fail")
	}
}
