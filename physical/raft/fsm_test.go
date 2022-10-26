package raft

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"testing"

	"github.com/go-test/deep"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/vault/sdk/physical"
)

func getFSM(t testing.TB) (*FSM, string) {
	raftDir, err := ioutil.TempDir("", "vault-raft-")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("raft dir: %s", raftDir)

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "raft",
		Level: hclog.Trace,
	})

	fsm, err := NewFSM(raftDir, "", logger)
	if err != nil {
		t.Fatal(err)
	}

	return fsm, raftDir
}

func TestFSM_Batching(t *testing.T) {
	fsm, dir := getFSM(t)
	defer func() { _ = os.RemoveAll(dir) }()

	var index uint64
	var term uint64 = 1

	getLog := func(i uint64) (int, *raft.Log) {
		if rand.Intn(10) >= 8 {
			term += 1
			return 0, &raft.Log{
				Index: i,
				Term:  term,
				Type:  raft.LogConfiguration,
				Data: raft.EncodeConfiguration(raft.Configuration{
					Servers: []raft.Server{
						{
							Address: "test",
							ID:      "test",
						},
					},
				}),
			}
		}

		command := &LogData{
			Operations: make([]*LogOperation, rand.Intn(10)),
		}

		for j := range command.Operations {
			command.Operations[j] = &LogOperation{
				OpType: putOp,
				Key:    fmt.Sprintf("key-%d-%d", i, j),
				Value:  []byte(fmt.Sprintf("value-%d-%d", i, j)),
			}
		}
		commandBytes, err := proto.Marshal(command)
		if err != nil {
			t.Fatal(err)
		}
		return len(command.Operations), &raft.Log{
			Index: i,
			Term:  term,
			Type:  raft.LogCommand,
			Data:  commandBytes,
		}
	}

	totalKeys := 0
	for i := 0; i < 100; i++ {
		batchSize := rand.Intn(64)
		batch := make([]*raft.Log, batchSize)
		for j := 0; j < batchSize; j++ {
			var keys int
			index++
			keys, batch[j] = getLog(index)
			totalKeys += keys
		}

		resp := fsm.ApplyBatch(batch)
		if len(resp) != batchSize {
			t.Fatalf("incorrect response length: got %d expected %d", len(resp), batchSize)
		}

		for _, r := range resp {
			if _, ok := r.(*FSMApplyResponse); !ok {
				t.Fatal("bad response type")
			}
		}
	}

	keys, err := fsm.List(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}

	if len(keys) != totalKeys {
		t.Fatalf("incorrect number of keys: got %d expected %d", len(keys), totalKeys)
	}

	latestIndex, latestConfig := fsm.LatestState()
	if latestIndex.Index != index {
		t.Fatalf("bad latest index: got %d expected %d", latestIndex.Index, index)
	}
	if latestIndex.Term != term {
		t.Fatalf("bad latest term: got %d expected %d", latestIndex.Term, term)
	}

	if latestConfig == nil && term > 1 {
		t.Fatal("config wasn't updated")
	}
}

func TestFSM_Batching_BadCommand(t *testing.T) {
	fsm, dir := getFSM(t)
	defer func() { _ = os.RemoveAll(dir) }()

	var index uint64
	var term uint64 = 1

	getLog := func(i uint64) (int, *raft.Log) {
		command := &LogData{
			Operations: make([]*LogOperation, 2),
		}

		for j := range command.Operations {
			command.Operations[j] = &LogOperation{
				OpType: putOp,
				Key:    fmt.Sprintf("key-%d-%d", i, j),
				Value:  []byte(fmt.Sprintf("value-%d-%d", i, j)),
			}

			// we have a batch size of 2 and each batch has 2 operations for a total of 4 operations.
			// 1 out of the 4 operations will be bad
			if i%2 == 1 && j%2 == 1 {
				command.Operations[j].OpType = 9000 // not a valid op
			}
		}
		commandBytes, err := proto.Marshal(command)
		if err != nil {
			t.Fatal(err)
		}
		return len(command.Operations), &raft.Log{
			Index: i,
			Term:  term,
			Type:  raft.LogCommand,
			Data:  commandBytes,
		}
	}

	batchSize := 2
	batch := make([]*raft.Log, batchSize)
	for j := 0; j < batchSize; j++ {
		index++
		_, batch[j] = getLog(index)
	}

	resp := fsm.ApplyBatch(batch)
	if len(resp) != batchSize {
		t.Fatalf("incorrect response length: got %d expected %d", len(resp), batchSize)
	}

	for i, r := range resp {
		fsmar, ok := r.(*FSMApplyResponse)
		if !ok {
			t.Fatal("bad response type")
		}

		if i == 0 && fsmar.Error == nil {
			t.Fatal("error should not be nil but it was")
		}

		if i == 1 && fsmar.Error != nil {
			t.Fatalf("error should've been nil but got: %v", fsmar.Error)
		}
	}
}

func TestFSM_List(t *testing.T) {
	fsm, dir := getFSM(t)
	defer func() { _ = os.RemoveAll(dir) }()

	ctx := context.Background()
	count := 100
	keys := rand.Perm(count)
	var sorted []string
	for _, k := range keys {
		err := fsm.Put(ctx, &physical.Entry{Key: fmt.Sprintf("foo/%d/bar", k)})
		if err != nil {
			t.Fatal(err)
		}
		err = fsm.Put(ctx, &physical.Entry{Key: fmt.Sprintf("foo/%d/baz", k)})
		if err != nil {
			t.Fatal(err)
		}
		sorted = append(sorted, fmt.Sprintf("%d/", k))
	}
	sort.Strings(sorted)

	got, err := fsm.List(ctx, "foo/")
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)
	if diff := deep.Equal(sorted, got); len(diff) > 0 {
		t.Fatal(diff)
	}
}
