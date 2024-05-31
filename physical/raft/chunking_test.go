// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package raft

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-raftchunking"
	raftchunkingtypes "github.com/hashicorp/go-raftchunking/types"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This chunks encoded data and then performing out-of-order applies of half
// the logs. It then snapshots, restores to a new FSM, and applies the rest.
// The goal is to verify that chunking snapshotting works as expected.
func TestRaft_Chunking_Lifecycle(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	assert := assert.New(t)

	b, dir := GetRaft(t, true, false)
	defer os.RemoveAll(dir)

	t.Log("applying configuration")

	b.applyConfigSettings(raft.DefaultConfig())

	t.Log("chunking")

	buf := []byte("let's see how this goes, shall we?")
	logData := &LogData{
		Operations: []*LogOperation{
			{
				OpType: putOp,
				Key:    "foobar",
				Value:  buf,
			},
		},
	}
	cmdBytes, err := proto.Marshal(logData)
	require.NoError(err)

	var logs []*raft.Log
	for i, b := range cmdBytes {
		// Stage multiple operations so we can test restoring across multiple opnums
		for j := 0; j < 10; j++ {
			chunkInfo := &raftchunkingtypes.ChunkInfo{
				OpNum:       uint64(32 + j),
				SequenceNum: uint32(i),
				NumChunks:   uint32(len(cmdBytes)),
			}
			chunkBytes, err := proto.Marshal(chunkInfo)
			require.NoError(err)

			logs = append(logs, &raft.Log{
				Data:       []byte{b},
				Extensions: chunkBytes,
			})
		}
	}

	t.Log("applying half of the logs")

	// The reason for the skipping is to test out-of-order applies which are
	// theoretically possible. Some of these will actually finish though!
	for i := 0; i < len(logs); i += 2 {
		resp := b.fsm.chunker.Apply(logs[i])
		if resp != nil {
			_, ok := resp.(raftchunking.ChunkingSuccess)
			assert.True(ok)
		}
	}

	t.Log("tearing down cluster")
	require.NoError(b.TeardownCluster(nil))
	require.NoError(b.fsm.getDB().Close())
	require.NoError(b.stableStore.(*raftboltdb.BoltStore).Close())

	t.Log("starting new backend")
	backendRaw, err := NewRaftBackend(b.conf, b.logger)
	require.NoError(err)
	b = backendRaw.(*RaftBackend)

	t.Log("applying rest of the logs")

	// Apply the rest of the logs
	var resp interface{}
	for i := 1; i < len(logs); i += 2 {
		resp = b.fsm.chunker.Apply(logs[i])
		if resp != nil {
			_, ok := resp.(raftchunking.ChunkingSuccess)
			assert.True(ok)
		}
	}

	assert.NotNil(resp)
	_, ok := resp.(raftchunking.ChunkingSuccess)
	assert.True(ok)
}

func TestFSM_Chunking_TermChange(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	assert := assert.New(t)

	b, dir := GetRaft(t, true, false)
	defer os.RemoveAll(dir)

	t.Log("applying configuration")

	b.applyConfigSettings(raft.DefaultConfig())

	t.Log("chunking")

	buf := []byte("let's see how this goes, shall we?")
	logData := &LogData{
		Operations: []*LogOperation{
			{
				OpType: putOp,
				Key:    "foobar",
				Value:  buf,
			},
		},
	}
	cmdBytes, err := proto.Marshal(logData)
	require.NoError(err)

	// Only need two chunks to test this
	chunks := [][]byte{
		cmdBytes[0:2],
		cmdBytes[2:],
	}
	var logs []*raft.Log
	for i, b := range chunks {
		chunkInfo := &raftchunkingtypes.ChunkInfo{
			OpNum:       uint64(32),
			SequenceNum: uint32(i),
			NumChunks:   uint32(len(chunks)),
		}
		chunkBytes, err := proto.Marshal(chunkInfo)
		if err != nil {
			t.Fatal(err)
		}
		logs = append(logs, &raft.Log{
			Term:       uint64(i),
			Data:       b,
			Extensions: chunkBytes,
		})
	}

	// We should see nil for both
	for _, log := range logs {
		resp := b.fsm.chunker.Apply(log)
		assert.Nil(resp)
	}

	// Now verify the other baseline, that when the term doesn't change we see
	// non-nil. First make the chunker have a clean state, then set the terms
	// to be the same.
	b.fsm.chunker.RestoreState(nil)
	logs[1].Term = uint64(0)

	// We should see nil only for the first one
	for i, log := range logs {
		resp := b.fsm.chunker.Apply(log)
		if i == 0 {
			assert.Nil(resp)
		}
		if i == 1 {
			assert.NotNil(resp)
			_, ok := resp.(raftchunking.ChunkingSuccess)
			assert.True(ok)
		}
	}
}

func TestRaft_Chunking_AppliedIndex(t *testing.T) {
	raft, dir := GetRaft(t, true, false)
	defer os.RemoveAll(dir)

	// Lower the size for tests
	raftchunking.ChunkSize = 1024
	val, err := uuid.GenerateRandomBytes(3 * raftchunking.ChunkSize)
	if err != nil {
		t.Fatal(err)
	}

	// Write a value to fastforward the index
	err = raft.Put(context.Background(), &physical.Entry{
		Key:   "key",
		Value: []byte("test"),
	})
	if err != nil {
		t.Fatal(err)
	}

	currentIndex := raft.AppliedIndex()
	// Write some data
	for i := 0; i < 10; i++ {
		err := raft.Put(context.Background(), &physical.Entry{
			Key:   fmt.Sprintf("key-%d", i),
			Value: val,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	newIndex := raft.AppliedIndex()

	// Each put should generate 4 chunks
	if newIndex-currentIndex != 10*4 {
		t.Fatalf("Did not apply chunks as expected, applied index = %d - %d = %d", newIndex, currentIndex, newIndex-currentIndex)
	}

	for i := 0; i < 10; i++ {
		entry, err := raft.Get(context.Background(), fmt.Sprintf("key-%d", i))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(entry.Value, val) {
			t.Fatal("value is corrupt")
		}
	}
}
