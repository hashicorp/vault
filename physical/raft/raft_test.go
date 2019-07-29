package raft

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	fmt "fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/golang/protobuf/proto"
	hclog "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/physical"
	bolt "go.etcd.io/bbolt"
)

func getRaft(t testing.TB, bootstrap bool, noStoreState bool) (*RaftBackend, string) {
	raftDir, err := ioutil.TempDir("", "vault-raft-")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("raft dir: %s", raftDir)

	return getRaftWithDir(t, bootstrap, noStoreState, raftDir)
}

func getRaftWithDir(t testing.TB, bootstrap bool, noStoreState bool, raftDir string) (*RaftBackend, string) {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "raft",
		Level: hclog.Trace,
	})
	logger.Info("raft dir", "dir", raftDir)

	conf := map[string]string{
		"path":          raftDir,
		"trailing_logs": "100",
	}

	if noStoreState {
		conf["doNotStoreLatestState"] = ""
	}

	backendRaw, err := NewRaftBackend(conf, logger)
	if err != nil {
		t.Fatal(err)
	}
	backend := backendRaw.(*RaftBackend)

	if bootstrap {
		err = backend.Bootstrap(context.Background(), []Peer{Peer{ID: backend.NodeID(), Address: backend.NodeID()}})
		if err != nil {
			t.Fatal(err)
		}

		err = backend.SetupCluster(context.Background(), nil, nil)
		if err != nil {
			t.Fatal(err)
		}

	}

	return backend, raftDir
}

func compareFSMs(t *testing.T, fsm1, fsm2 *FSM) {
	t.Helper()
	index1, config1 := fsm1.LatestState()
	index2, config2 := fsm2.LatestState()

	if !proto.Equal(index1, index2) {
		t.Fatalf("indexes did not match: %+v != %+v", index1, index2)
	}
	if !proto.Equal(config1, config2) {
		t.Fatalf("configs did not match: %+v != %+v", config1, config2)
	}

	compareDBs(t, fsm1.db, fsm2.db)
}

func compareDBs(t *testing.T, boltDB1, boltDB2 *bolt.DB) {
	db1 := make(map[string]string)
	db2 := make(map[string]string)

	err := boltDB1.View(func(tx *bolt.Tx) error {

		c := tx.Cursor()
		for bucketName, _ := c.First(); bucketName != nil; bucketName, _ = c.Next() {
			b := tx.Bucket(bucketName)

			cBucket := b.Cursor()

			for k, v := cBucket.First(); k != nil; k, v = cBucket.Next() {
				db1[string(k)] = base64.StdEncoding.EncodeToString(v)
			}
		}

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	err = boltDB2.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()
		for bucketName, _ := c.First(); bucketName != nil; bucketName, _ = c.Next() {
			b := tx.Bucket(bucketName)

			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				db2[string(k)] = base64.StdEncoding.EncodeToString(v)
			}
		}

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(db1, db2); diff != nil {
		t.Fatal(diff)
	}
}

func TestRaft_Backend(t *testing.T) {
	b, dir := getRaft(t, true, true)
	defer os.RemoveAll(dir)

	physical.ExerciseBackend(t, b)
}

func TestRaft_Backend_ListPrefix(t *testing.T) {
	b, dir := getRaft(t, true, true)
	defer os.RemoveAll(dir)

	physical.ExerciseBackend_ListPrefix(t, b)
}

func TestRaft_TransactionalBackend(t *testing.T) {
	b, dir := getRaft(t, true, true)
	defer os.RemoveAll(dir)

	physical.ExerciseTransactionalBackend(t, b)
}

func TestRaft_HABackend(t *testing.T) {
	t.Skip()
	raft, dir := getRaft(t, true, true)
	defer os.RemoveAll(dir)
	raft2, dir2 := getRaft(t, false, true)
	defer os.RemoveAll(dir2)

	// Add raft2 to the cluster
	addPeer(t, raft, raft2)

	physical.ExerciseHABackend(t, raft, raft2)
}

func TestRaft_Backend_ThreeNode(t *testing.T) {
	raft1, dir := getRaft(t, true, true)
	raft2, dir2 := getRaft(t, false, true)
	raft3, dir3 := getRaft(t, false, true)
	defer os.RemoveAll(dir)
	defer os.RemoveAll(dir2)
	defer os.RemoveAll(dir3)

	// Add raft2 to the cluster
	addPeer(t, raft1, raft2)

	// Add raft3 to the cluster
	addPeer(t, raft1, raft3)

	physical.ExerciseBackend(t, raft1)

	time.Sleep(10 * time.Second)
	// Make sure all stores are the same
	compareFSMs(t, raft1.fsm, raft2.fsm)
	compareFSMs(t, raft1.fsm, raft3.fsm)
}

func TestRaft_Recovery(t *testing.T) {
	// Create 4 raft nodes
	raft1, dir1 := getRaft(t, true, true)
	raft2, dir2 := getRaft(t, false, true)
	raft3, dir3 := getRaft(t, false, true)
	raft4, dir4 := getRaft(t, false, true)
	defer os.RemoveAll(dir1)
	defer os.RemoveAll(dir2)
	defer os.RemoveAll(dir3)
	defer os.RemoveAll(dir4)

	// Add them all to the cluster
	addPeer(t, raft1, raft2)
	addPeer(t, raft1, raft3)
	addPeer(t, raft1, raft4)

	// Add some data into the FSM
	physical.ExerciseBackend(t, raft1)

	time.Sleep(10 * time.Second)

	// Bring down all nodes
	raft1.TeardownCluster(nil)
	raft2.TeardownCluster(nil)
	raft3.TeardownCluster(nil)
	raft4.TeardownCluster(nil)

	// Prepare peers.json
	type RecoveryPeer struct {
		ID       string `json:"id"`
		Address  string `json:"address"`
		NonVoter bool   `json: non_voter`
	}

	// Leave out node 1 during recovery
	peersList := make([]*RecoveryPeer, 0, 3)
	peersList = append(peersList, &RecoveryPeer{
		ID:       raft1.NodeID(),
		Address:  raft1.NodeID(),
		NonVoter: false,
	})
	peersList = append(peersList, &RecoveryPeer{
		ID:       raft2.NodeID(),
		Address:  raft2.NodeID(),
		NonVoter: false,
	})
	peersList = append(peersList, &RecoveryPeer{
		ID:       raft4.NodeID(),
		Address:  raft4.NodeID(),
		NonVoter: false,
	})

	peersJSONBytes, err := jsonutil.EncodeJSON(peersList)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(filepath.Join(dir1, raftState), "peers.json"), peersJSONBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(filepath.Join(dir2, raftState), "peers.json"), peersJSONBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(filepath.Join(dir4, raftState), "peers.json"), peersJSONBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Bring up the nodes again
	raft1.SetupCluster(context.Background(), nil, nil)
	raft2.SetupCluster(context.Background(), nil, nil)
	raft4.SetupCluster(context.Background(), nil, nil)

	peers, err := raft1.Peers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(peers) != 3 {
		t.Fatalf("failed to recover the cluster")
	}

	time.Sleep(10 * time.Second)

	compareFSMs(t, raft1.fsm, raft2.fsm)
	compareFSMs(t, raft1.fsm, raft4.fsm)
}

func TestRaft_TransactionalBackend_ThreeNode(t *testing.T) {
	raft1, dir := getRaft(t, true, true)
	raft2, dir2 := getRaft(t, false, true)
	raft3, dir3 := getRaft(t, false, true)
	defer os.RemoveAll(dir)
	defer os.RemoveAll(dir2)
	defer os.RemoveAll(dir3)

	// Add raft2 to the cluster
	addPeer(t, raft1, raft2)

	// Add raft3 to the cluster
	addPeer(t, raft1, raft3)

	physical.ExerciseTransactionalBackend(t, raft1)

	time.Sleep(10 * time.Second)
	// Make sure all stores are the same
	compareFSMs(t, raft1.fsm, raft2.fsm)
	compareFSMs(t, raft1.fsm, raft3.fsm)
}

func TestRaft_Backend_Performance(t *testing.T) {
	b, dir := getRaft(t, true, false)
	defer os.RemoveAll(dir)

	defaultConfig := raft.DefaultConfig()

	localConfig := raft.DefaultConfig()
	b.applyConfigSettings(localConfig)

	if localConfig.ElectionTimeout != defaultConfig.ElectionTimeout*5 {
		t.Fatalf("bad config: %v", localConfig)
	}
	if localConfig.HeartbeatTimeout != defaultConfig.HeartbeatTimeout*5 {
		t.Fatalf("bad config: %v", localConfig)
	}
	if localConfig.LeaderLeaseTimeout != defaultConfig.LeaderLeaseTimeout*5 {
		t.Fatalf("bad config: %v", localConfig)
	}

	b.conf = map[string]string{
		"path":                   dir,
		"performance_multiplier": "5",
	}

	localConfig = raft.DefaultConfig()
	b.applyConfigSettings(localConfig)

	if localConfig.ElectionTimeout != defaultConfig.ElectionTimeout*5 {
		t.Fatalf("bad config: %v", localConfig)
	}
	if localConfig.HeartbeatTimeout != defaultConfig.HeartbeatTimeout*5 {
		t.Fatalf("bad config: %v", localConfig)
	}
	if localConfig.LeaderLeaseTimeout != defaultConfig.LeaderLeaseTimeout*5 {
		t.Fatalf("bad config: %v", localConfig)
	}

	b.conf = map[string]string{
		"path":                   dir,
		"performance_multiplier": "1",
	}

	localConfig = raft.DefaultConfig()
	b.applyConfigSettings(localConfig)

	if localConfig.ElectionTimeout != defaultConfig.ElectionTimeout {
		t.Fatalf("bad config: %v", localConfig)
	}
	if localConfig.HeartbeatTimeout != defaultConfig.HeartbeatTimeout {
		t.Fatalf("bad config: %v", localConfig)
	}
	if localConfig.LeaderLeaseTimeout != defaultConfig.LeaderLeaseTimeout {
		t.Fatalf("bad config: %v", localConfig)
	}

}

func BenchmarkDB_Puts(b *testing.B) {
	raft, dir := getRaft(b, true, false)
	defer os.RemoveAll(dir)
	raft2, dir2 := getRaft(b, true, false)
	defer os.RemoveAll(dir2)

	bench := func(b *testing.B, s physical.Backend, dataSize int) {
		data, err := uuid.GenerateRandomBytes(dataSize)
		if err != nil {
			b.Fatal(err)
		}

		ctx := context.Background()
		pe := &physical.Entry{
			Value: data,
		}
		testName := b.Name()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pe.Key = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%d", testName, i))))
			err := s.Put(ctx, pe)
			if err != nil {
				b.Fatal(err)
			}
		}
	}

	b.Run("256b", func(b *testing.B) { bench(b, raft, 256) })
	b.Run("256kb", func(b *testing.B) { bench(b, raft2, 256*1024) })
}

func BenchmarkDB_Snapshot(b *testing.B) {
	raft, dir := getRaft(b, true, false)
	defer os.RemoveAll(dir)

	data, err := uuid.GenerateRandomBytes(256 * 1024)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	pe := &physical.Entry{
		Value: data,
	}
	testName := b.Name()

	for i := 0; i < 100; i++ {
		pe.Key = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%d", testName, i))))
		err = raft.Put(ctx, pe)
		if err != nil {
			b.Fatal(err)
		}
	}

	bench := func(b *testing.B, s *FSM) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pe.Key = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%d", testName, i))))
			s.writeTo(ctx, discardCloser{Writer: ioutil.Discard}, discardCloser{Writer: ioutil.Discard})
		}
	}

	b.Run("256kb", func(b *testing.B) { bench(b, raft.fsm) })
}

type discardCloser struct {
	io.Writer
}

func (d discardCloser) Close() error               { return nil }
func (d discardCloser) CloseWithError(error) error { return nil }
