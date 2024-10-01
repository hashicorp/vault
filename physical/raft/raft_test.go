// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package raft

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/golang/protobuf/proto"
	bolt "github.com/hashicorp-forge/bbolt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/base62"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/stretchr/testify/require"
)

func testBothRaftBackends(t *testing.T, f func(t *testing.T, raftWALValue string)) {
	t.Helper()

	testCases := []struct {
		name   string
		useWAL string
	}{
		{
			name:   "use wal",
			useWAL: "true",
		},
		{
			name:   "use boltdb",
			useWAL: "false",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// we can't use t.Parallel() here because some raft tests manipulate package level variables
			f(t, tc.useWAL)
		})
	}
}

func testBothRaftBackendsBenchmark(b *testing.B, f func(raftWALValue string)) {
	testCases := []struct {
		name   string
		useWAL string
	}{
		{
			name:   "use wal",
			useWAL: "true",
		},
		{
			name:   "use boltdb",
			useWAL: "false",
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			f(tc.useWAL)
		})
	}
}

func connectPeers(nodes ...*RaftBackend) {
	for _, node := range nodes {
		for _, peer := range nodes {
			if node == peer {
				continue
			}

			node.raftTransport.(*raft.InmemTransport).Connect(raft.ServerAddress(peer.NodeID()), peer.raftTransport)
			peer.raftTransport.(*raft.InmemTransport).Connect(raft.ServerAddress(node.NodeID()), node.raftTransport)
		}
	}
}

func stepDownLeader(t *testing.T, node *RaftBackend) {
	t.Helper()

	if err := node.raft.LeadershipTransfer().Error(); err != nil {
		t.Fatal(err)
	}

	timeout := time.Now().Add(time.Second * 10)
	for !time.Now().After(timeout) {
		if err := node.raft.VerifyLeader().Error(); err != nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatal("still leader")
}

func waitForLeader(t *testing.T, nodes ...*RaftBackend) *RaftBackend {
	t.Helper()
	timeout := time.Now().Add(time.Second * 10)
	for !time.Now().After(timeout) {
		for _, node := range nodes {
			if node.raft.Leader() == raft.ServerAddress(node.NodeID()) {
				return node
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatal("no leader")
	return nil
}

func compareFSMs(t *testing.T, fsm1, fsm2 *FSM) {
	t.Helper()
	if err := compareFSMsWithErr(t, fsm1, fsm2); err != nil {
		t.Fatal(err)
	}
}

func compareFSMsWithErr(t *testing.T, fsm1, fsm2 *FSM) error {
	t.Helper()
	index1, config1 := fsm1.LatestState()
	index2, config2 := fsm2.LatestState()

	if !proto.Equal(index1, index2) {
		return fmt.Errorf("indexes did not match: %+v != %+v", index1, index2)
	}
	if !proto.Equal(config1, config2) {
		return fmt.Errorf("configs did not match: %+v != %+v", config1, config2)
	}

	return compareDBs(t, fsm1.getDB(), fsm2.getDB(), false)
}

func compareDBs(t *testing.T, boltDB1, boltDB2 *bolt.DB, dataOnly bool) error {
	t.Helper()
	db1 := make(map[string]string)
	db2 := make(map[string]string)

	err := boltDB1.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()
		for bucketName, _ := c.First(); bucketName != nil; bucketName, _ = c.Next() {
			if dataOnly && !bytes.Equal(bucketName, dataBucketName) {
				continue
			}

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
			if dataOnly && !bytes.Equal(bucketName, dataBucketName) {
				continue
			}
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
		return fmt.Errorf("%+v", diff)
	}

	return nil
}

func TestRaft_Backend(t *testing.T) {
	t.Parallel()

	testBothRaftBackends(t, func(t *testing.T, useRaftWal string) {
		conf := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		b, _ := GetRaftWithConfig(t, true, true, conf)
		physical.ExerciseBackend(t, b)
	})
}

// TestRaft_SwitchFromBoltDBToRaftWal is testing that we don't use raft-wal, even if configured to do so,
// if there is an existing raft.db file on disk (meaning BoltDB was previously in use).
func TestRaft_SwitchFromBoltDBToRaftWal(t *testing.T) {
	tmpDir := t.TempDir()

	// configured to use raft-wal
	conf := map[string]string{
		"path":          tmpDir,
		"trailing_logs": "100",
		"raft_wal":      "true",
	}

	// raftBaseDir will end up looking like $tmpDir/raft
	raftBaseDir := filepath.Join(tmpDir, raftState)
	err := os.MkdirAll(raftBaseDir, 0o777)
	if err != nil {
		t.Fatal(err)
	}

	// create a bogus $tmpDir/raft/raft.db file
	db, err := bolt.Open(filepath.Join(raftBaseDir, "raft.db"), 0o777, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewRaftBackend(conf, hclog.NewNullLogger())
	if err != nil {
		t.Fatal(err)
	}

	// Check to see if $tmpDir/raft/raft-wal exists. It should not, because we only create that if raft-wal is in use.
	// And since raft.db already existed, we should've skipped all the raft-wal setup code.
	raftWalExists, err := fileExists(filepath.Join(raftBaseDir, raftWalDir))
	if err != nil {
		t.Fatal(err)
	}

	if raftWalExists {
		t.Fatal("expected raft-wal dir to not exist, but it does")
	}
}

// TestRaft_VerifierEnabled is not checking to ensure that the verifier works correctly - the verifier has
// its own unit tests for that. What we're checking for here is that we've plumbed everything through correctly,
// i.e. we can stand up a raft cluster with the verifier enabled, do a bunch of raft things, let the verifier
// do its thing, and nothing blows up.
func TestRaft_VerifierEnabled(t *testing.T) {
	testBothRaftBackends(t, func(t *testing.T, useRaftWal string) {
		conf := map[string]string{
			"trailing_logs":             "100",
			"raft_wal":                  useRaftWal,
			"raft_log_verifier_enabled": "true",
		}

		b, _ := GetRaftWithConfig(t, true, true, conf)
		physical.ExerciseBackend(t, b)

		err := b.applyVerifierCheckpoint()
		if err != nil {
			t.Fatal(err)
		}
		physical.ExerciseBackend(t, b)
	})
}

func TestRaft_Backend_LargeKey(t *testing.T) {
	t.Parallel()

	testBothRaftBackends(t, func(t *testing.T, useRaftWal string) {
		conf := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		b, _ := GetRaftWithConfig(t, true, true, conf)
		key, err := base62.Random(bolt.MaxKeySize + 1)
		if err != nil {
			t.Fatal(err)
		}
		entry := &physical.Entry{Key: key, Value: []byte(key)}

		err = b.Put(context.Background(), entry)
		if err == nil {
			t.Fatal("expected error for put entry")
		}

		if !strings.Contains(err.Error(), physical.ErrKeyTooLarge) {
			t.Fatalf("expected %q, got %v", physical.ErrKeyTooLarge, err)
		}

		out, err := b.Get(context.Background(), entry.Key)
		if err != nil {
			t.Fatalf("unexpected error after failed put: %v", err)
		}
		if out != nil {
			t.Fatal("expected response entry to be nil after a failed put")
		}
	})
}

func TestRaft_Backend_LargeValue(t *testing.T) {
	t.Parallel()

	testBothRaftBackends(t, func(t *testing.T, useRaftWal string) {
		conf := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		b, _ := GetRaftWithConfig(t, true, true, conf)
		value := make([]byte, defaultMaxEntrySize+1)
		rand.Read(value)
		entry := &physical.Entry{Key: "foo", Value: value}

		err := b.Put(context.Background(), entry)
		if err == nil {
			t.Fatal("expected error for put entry")
		}

		if !strings.Contains(err.Error(), physical.ErrValueTooLarge) {
			t.Fatalf("expected %q, got %v", physical.ErrValueTooLarge, err)
		}

		out, err := b.Get(context.Background(), entry.Key)
		if err != nil {
			t.Fatalf("unexpected error after failed put: %v", err)
		}
		if out != nil {
			t.Fatal("expected response entry to be nil after a failed put")
		}
	})
}

// TestRaft_TransactionalBackend_GetTransactions tests that passing a slice of transactions to the
// raft backend will populate values for any transactions that are Get operations.
func TestRaft_TransactionalBackend_GetTransactions(t *testing.T) {
	t.Parallel()
	testBothRaftBackends(t, func(t *testing.T, useRaftWal string) {
		conf := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		b, _ := GetRaftWithConfig(t, true, true, conf)
		ctx := context.Background()
		txns := make([]*physical.TxnEntry, 0)

		// Add some seed values to our FSM, and prepare our slice of transactions at the same time
		for i := 0; i < 5; i++ {
			key := fmt.Sprintf("foo/%d", i)
			err := b.fsm.Put(ctx, &physical.Entry{Key: key, Value: []byte(fmt.Sprintf("value-%d", i))})
			if err != nil {
				t.Fatal(err)
			}

			txns = append(txns, &physical.TxnEntry{
				Operation: physical.GetOperation,
				Entry: &physical.Entry{
					Key: key,
				},
			})
		}

		// Add some additional transactions, so we have a mix of operations
		for i := 0; i < 10; i++ {
			txnEntry := &physical.TxnEntry{
				Entry: &physical.Entry{
					Key: fmt.Sprintf("lol-%d", i),
				},
			}

			if i%2 == 0 {
				txnEntry.Operation = physical.PutOperation
				txnEntry.Entry.Value = []byte("lol")
			} else {
				txnEntry.Operation = physical.DeleteOperation
			}

			txns = append(txns, txnEntry)
		}

		err := b.Transaction(ctx, txns)
		if err != nil {
			t.Fatal(err)
		}

		// Check that our Get operations were populated with their values
		for i, txn := range txns {
			if txn.Operation == physical.GetOperation {
				val := []byte(fmt.Sprintf("value-%d", i))
				if !bytes.Equal(val, txn.Entry.Value) {
					t.Fatalf("expected %s to equal %s but it didn't", hex.EncodeToString(val), hex.EncodeToString(txn.Entry.Value))
				}
			}
		}
	})
}

func TestRaft_TransactionalBackend_LargeKey(t *testing.T) {
	t.Parallel()
	testBothRaftBackends(t, func(t *testing.T, useRaftWal string) {
		conf := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		b, _ := GetRaftWithConfig(t, true, true, conf)
		value := make([]byte, defaultMaxEntrySize+1)
		rand.Read(value)

		key, err := base62.Random(bolt.MaxKeySize + 1)
		if err != nil {
			t.Fatal(err)
		}
		txns := []*physical.TxnEntry{
			{
				Operation: physical.PutOperation,
				Entry: &physical.Entry{
					Key:   key,
					Value: []byte(key),
				},
			},
		}

		err = b.Transaction(context.Background(), txns)
		if err == nil {
			t.Fatal("expected error for transactions")
		}

		if !strings.Contains(err.Error(), physical.ErrKeyTooLarge) {
			t.Fatalf("expected %q, got %v", physical.ErrValueTooLarge, err)
		}

		out, err := b.Get(context.Background(), txns[0].Entry.Key)
		if err != nil {
			t.Fatalf("unexpected error after failed put: %v", err)
		}
		if out != nil {
			t.Fatal("expected response entry to be nil after a failed put")
		}
	})
}

func TestRaft_TransactionalBackend_LargeValue(t *testing.T) {
	t.Parallel()
	testBothRaftBackends(t, func(t *testing.T, useRaftWal string) {
		conf := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		b, _ := GetRaftWithConfig(t, true, true, conf)
		value := make([]byte, defaultMaxEntrySize+1)
		rand.Read(value)

		txns := []*physical.TxnEntry{
			{
				Operation: physical.PutOperation,
				Entry: &physical.Entry{
					Key:   "foo",
					Value: value,
				},
			},
		}

		err := b.Transaction(context.Background(), txns)
		if err == nil {
			t.Fatal("expected error for transactions")
		}

		if !strings.Contains(err.Error(), physical.ErrValueTooLarge) {
			t.Fatalf("expected %q, got %v", physical.ErrValueTooLarge, err)
		}

		out, err := b.Get(context.Background(), txns[0].Entry.Key)
		if err != nil {
			t.Fatalf("unexpected error after failed put: %v", err)
		}
		if out != nil {
			t.Fatal("expected response entry to be nil after a failed put")
		}
	})
}

func TestRaft_Backend_ListPrefix(t *testing.T) {
	t.Parallel()
	testBothRaftBackends(t, func(t *testing.T, useRaftWal string) {
		conf := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		b, _ := GetRaftWithConfig(t, true, true, conf)
		physical.ExerciseBackend_ListPrefix(t, b)
	})
}

func TestRaft_TransactionalBackend(t *testing.T) {
	t.Parallel()
	testBothRaftBackends(t, func(t *testing.T, useRaftWal string) {
		conf := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		b, _ := GetRaftWithConfig(t, true, true, conf)
		physical.ExerciseTransactionalBackend(t, b)
	})
}

func TestRaft_HABackend(t *testing.T) {
	t.Skip()
	raft1, _ := GetRaft(t, true, true)
	raft2, _ := GetRaft(t, false, true)

	// Add raft2 to the cluster
	addPeer(t, raft1, raft2)
	physical.ExerciseHABackend(t, raft1, raft2)
}

func TestRaft_Backend_ThreeNode(t *testing.T) {
	t.Parallel()
	testBothRaftBackends(t, func(t *testing.T, useRaftWal string) {
		conf := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		raft1, _ := GetRaftWithConfig(t, true, true, conf)
		raft2, _ := GetRaftWithConfig(t, false, true, conf)
		raft3, _ := GetRaftWithConfig(t, false, true, conf)

		// Add raft2 to the cluster
		addPeer(t, raft1, raft2)

		// Add raft3 to the cluster
		addPeer(t, raft1, raft3)

		physical.ExerciseBackend(t, raft1)

		time.Sleep(10 * time.Second)
		// Make sure all stores are the same
		compareFSMs(t, raft1.fsm, raft2.fsm)
		compareFSMs(t, raft1.fsm, raft3.fsm)
	})
}

func TestRaft_GetOfflineConfig(t *testing.T) {
	t.Parallel()
	testBothRaftBackends(t, func(t *testing.T, useRaftWal string) {
		config := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		// Create 3 raft nodes
		raft1, _ := GetRaftWithConfig(t, true, true, config)
		raft2, _ := GetRaftWithConfig(t, false, true, config)
		raft3, _ := GetRaftWithConfig(t, false, true, config)

		// Add them all to the cluster
		addPeer(t, raft1, raft2)
		addPeer(t, raft1, raft3)

		// Add some data into the FSM
		physical.ExerciseBackend(t, raft1)

		time.Sleep(10 * time.Second)

		// Spin down the raft cluster and check that GetConfigurationOffline
		// returns 3 voters
		err := raft3.TeardownCluster(nil)
		if err != nil {
			t.Fatal(err)
		}
		err = raft2.TeardownCluster(nil)
		if err != nil {
			t.Fatal(err)
		}
		err = raft1.TeardownCluster(nil)
		if err != nil {
			t.Fatal(err)
		}

		conf, err := raft1.GetConfigurationOffline()
		if err != nil {
			t.Fatal(err)
		}
		if len(conf.Servers) != 3 {
			t.Fatalf("three raft nodes existed but we only see %d", len(conf.Servers))
		}
		for _, s := range conf.Servers {
			if s.Voter != true {
				t.Fatalf("one of the nodes is not a voter")
			}
		}
	})
}

func TestRaft_Recovery(t *testing.T) {
	t.Parallel()
	testBothRaftBackends(t, func(t *testing.T, useRaftWal string) {
		conf := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		// Create 4 raft nodes
		raft1, dir1 := GetRaftWithConfig(t, true, true, conf)
		raft2, dir2 := GetRaftWithConfig(t, false, true, conf)
		raft3, _ := GetRaftWithConfig(t, false, true, conf)
		raft4, dir4 := GetRaftWithConfig(t, false, true, conf)

		// Add them all to the cluster
		addPeer(t, raft1, raft2)
		addPeer(t, raft1, raft3)
		addPeer(t, raft1, raft4)

		// Add some data into the FSM
		physical.ExerciseBackend(t, raft1)

		time.Sleep(10 * time.Second)

		// Bring down all nodes
		err := raft1.TeardownCluster(nil)
		if err != nil {
			t.Fatal(err)
		}
		err = raft2.TeardownCluster(nil)
		if err != nil {
			t.Fatal(err)
		}
		err = raft3.TeardownCluster(nil)
		if err != nil {
			t.Fatal(err)
		}
		err = raft4.TeardownCluster(nil)
		if err != nil {
			t.Fatal(err)
		}

		// Prepare peers.json
		type RecoveryPeer struct {
			ID       string `json:"id"`
			Address  string `json:"address"`
			NonVoter bool   `json:"non_voter"`
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
		err = os.WriteFile(filepath.Join(filepath.Join(dir1, raftState), "peers.json"), peersJSONBytes, 0o644)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(filepath.Join(filepath.Join(dir2, raftState), "peers.json"), peersJSONBytes, 0o644)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(filepath.Join(filepath.Join(dir4, raftState), "peers.json"), peersJSONBytes, 0o644)
		if err != nil {
			t.Fatal(err)
		}

		// Bring up the nodes again
		err = raft1.SetupCluster(context.Background(), SetupOpts{})
		if err != nil {
			t.Fatal(err)
		}
		err = raft2.SetupCluster(context.Background(), SetupOpts{})
		if err != nil {
			t.Fatal(err)
		}
		err = raft4.SetupCluster(context.Background(), SetupOpts{})
		if err != nil {
			t.Fatal(err)
		}

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
	})
}

func TestRaft_TransactionalBackend_ThreeNode(t *testing.T) {
	t.Parallel()
	testBothRaftBackends(t, func(t *testing.T, useRaftWal string) {
		conf := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		raft1, _ := GetRaftWithConfig(t, true, true, conf)
		raft2, _ := GetRaftWithConfig(t, false, true, conf)
		raft3, _ := GetRaftWithConfig(t, false, true, conf)

		// Add raft2 to the cluster
		addPeer(t, raft1, raft2)

		// Add raft3 to the cluster
		addPeer(t, raft1, raft3)

		physical.ExerciseTransactionalBackend(t, raft1)

		time.Sleep(10 * time.Second)
		// Make sure all stores are the same
		compareFSMs(t, raft1.fsm, raft2.fsm)
		compareFSMs(t, raft1.fsm, raft3.fsm)
	})
}

// TestRaft_TransactionalLimitsEnvOverride ensures the ENV var overrides for
// transaction size limits are plumbed through as expected.
func TestRaft_TransactionalLimitsEnvOverride(t *testing.T) {
	tc := []struct {
		name        string
		envEntries  string
		envSize     string
		wantEntries int
		wantSize    int
		wantLog     string
	}{
		{
			name:        "defaults",
			wantEntries: defaultMaxBatchEntries,
			wantSize:    defaultMaxBatchSize,
		},
		{
			name:        "valid env",
			envEntries:  "123",
			envSize:     "456",
			wantEntries: 123,
			wantSize:    456,
		},
		{
			name:        "invalid entries",
			envEntries:  "not-a-number",
			envSize:     "100",
			wantEntries: defaultMaxBatchEntries,
			wantSize:    100,
			wantLog:     "failed to parse VAULT_RAFT_MAX_BATCH_ENTRIES",
		},
		{
			name:        "invalid entries",
			envEntries:  "100",
			envSize:     "asdasdsasd",
			wantEntries: 100,
			wantSize:    defaultMaxBatchSize,
			wantLog:     "failed to parse VAULT_RAFT_MAX_BATCH_SIZE_BYTES",
		},
		{
			name:        "zero entries",
			envEntries:  "0",
			envSize:     "100",
			wantEntries: defaultMaxBatchEntries,
			wantSize:    100,
			wantLog:     "failed to parse VAULT_RAFT_MAX_BATCH_ENTRIES as an integer > 0",
		},
		{
			name:        "zero size",
			envEntries:  "100",
			envSize:     "0",
			wantEntries: 100,
			wantSize:    defaultMaxBatchSize,
			wantLog:     "failed to parse VAULT_RAFT_MAX_BATCH_SIZE_BYTES as an integer > 0",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			// Set the env vars within this test
			if tt.envEntries != "" {
				t.Setenv(EnvVaultRaftMaxBatchEntries, tt.envEntries)
			}
			if tt.envSize != "" {
				t.Setenv(EnvVaultRaftMaxBatchSizeBytes, tt.envSize)
			}

			var logBuf bytes.Buffer
			raft1, dir := GetRaftWithLogOutput(t, false, true, &logBuf)
			defer os.RemoveAll(dir)

			e, s := raft1.TransactionLimits()

			require.Equal(t, tt.wantEntries, e)
			require.Equal(t, tt.wantSize, s)
			if tt.wantLog != "" {
				require.Contains(t, logBuf.String(), tt.wantLog)
			}
		})
	}
}

func TestRaft_Backend_Performance(t *testing.T) {
	t.Parallel()
	testBothRaftBackends(t, func(t *testing.T, useRaftWal string) {
		conf := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		b, dir := GetRaftWithConfig(t, true, true, conf)

		defaultConfig := raft.DefaultConfig()
		localConfig := raft.DefaultConfig()
		err := b.applyConfigSettings(localConfig)
		if err != nil {
			t.Fatal(err)
		}

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
		err = b.applyConfigSettings(localConfig)
		if err != nil {
			t.Fatal(err)
		}

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
		err = b.applyConfigSettings(localConfig)
		if err != nil {
			t.Fatal(err)
		}

		if localConfig.ElectionTimeout != defaultConfig.ElectionTimeout {
			t.Fatalf("bad config: %v", localConfig)
		}
		if localConfig.HeartbeatTimeout != defaultConfig.HeartbeatTimeout {
			t.Fatalf("bad config: %v", localConfig)
		}
		if localConfig.LeaderLeaseTimeout != defaultConfig.LeaderLeaseTimeout {
			t.Fatalf("bad config: %v", localConfig)
		}
	})
}

func BenchmarkDB_Puts(b *testing.B) {
	testBothRaftBackendsBenchmark(b, func(useRaftWal string) {
		conf := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		raft1, _ := GetRaftWithConfig(b, true, false, conf)
		raft2, _ := GetRaftWithConfig(b, true, false, conf)

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

		b.Run("256b", func(b *testing.B) { bench(b, raft1, 256) })
		b.Run("256kb", func(b *testing.B) { bench(b, raft2, 256*1024) })
	})
}

func BenchmarkDB_Snapshot(b *testing.B) {
	testBothRaftBackendsBenchmark(b, func(useRaftWal string) {
		conf := map[string]string{
			"trailing_logs": "100",
			"raft_wal":      useRaftWal,
		}

		raft1, _ := GetRaftWithConfig(b, true, false, conf)
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
			err = raft1.Put(ctx, pe)
			if err != nil {
				b.Fatal(err)
			}
		}

		bench := func(b *testing.B, s *FSM) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				pe.Key = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%d", testName, i))))
				s.writeTo(ctx, discardCloser{Writer: io.Discard}, discardCloser{Writer: io.Discard})
			}
		}

		b.Run("256kb", func(b *testing.B) { bench(b, raft1.fsm) })
	})
}

type discardCloser struct {
	io.Writer
}

func (d discardCloser) Close() error               { return nil }
func (d discardCloser) CloseWithError(error) error { return nil }
