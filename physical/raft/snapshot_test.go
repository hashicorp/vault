package raft

import (
	"bytes"
	"context"
	fmt "fmt"
	"hash/crc64"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/vault/sdk/physical"
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

	err = follower.Bootstrap(context.Background(), peers)
	if err != nil {
		t.Fatal(err)
	}

	err = follower.SetupCluster(context.Background(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	leader.raftTransport.(*raft.InmemTransport).Connect(raft.ServerAddress(follower.NodeID()), follower.raftTransport)
	follower.raftTransport.(*raft.InmemTransport).Connect(raft.ServerAddress(leader.NodeID()), leader.raftTransport)
}

func TestRaft_Snapshot_Loading(t *testing.T) {
	raft, dir := getRaft(t, true, false)
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
	raft, dir := getRaft(t, true, false)
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
	if index.Term != 1 {
		t.Fatalf("unexpected term, got %d expected 1", index.Term)
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
	if index.Term != 1 {
		t.Fatalf("unexpected term, got %d expected 1", index.Term)
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
	if meta.Term != 1 {
		t.Fatalf("unexpected snapshot term %d", meta.Term)
	}
}

func TestRaft_Snapshot_Peers(t *testing.T) {
	raft1, dir := getRaft(t, true, false)
	raft2, dir2 := getRaft(t, false, false)
	raft3, dir3 := getRaft(t, false, false)
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

	// Add raft2 to the cluster
	addPeer(t, raft1, raft2)

	// TODO: remove sleeps from these tests
	time.Sleep(10 * time.Second)

	// Make sure the snapshot was applied correctly on the follower
	compareDBs(t, raft1.fsm.db, raft2.fsm.db)

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

	// Add raft3 to the cluster
	addPeer(t, raft1, raft3)

	// TODO: remove sleeps from these tests
	time.Sleep(10 * time.Second)

	// Make sure all stores are the same
	compareFSMs(t, raft1.fsm, raft2.fsm)
	compareFSMs(t, raft1.fsm, raft3.fsm)
}

func TestRaft_Snapshot_Restart(t *testing.T) {
	raft1, dir := getRaft(t, true, false)
	defer os.RemoveAll(dir)
	raft2, dir2 := getRaft(t, false, false)
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

	// Shutdown raft1
	if err := raft1.TeardownCluster(nil); err != nil {
		t.Fatal(err)
	}

	// Start Raft
	err = raft1.SetupCluster(context.Background(), nil, nil)
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

func TestRaft_Snapshot_Take_Restore(t *testing.T) {
	raft1, dir := getRaft(t, true, false)
	defer os.RemoveAll(dir)
	raft2, dir2 := getRaft(t, false, false)
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

	snap := &bytes.Buffer{}

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

	snapFile, cleanup, metadata, err := raft1.WriteSnapshotToTemp(ioutil.NopCloser(snap), nil)
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
