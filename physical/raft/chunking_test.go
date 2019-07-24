package raft

import (
	"os"
	"testing"

	proto "github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-raftchunking"
	raftchunkingtypes "github.com/hashicorp/go-raftchunking/types"
	"github.com/hashicorp/raft"
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

	b, dir := getRaft(t, true, false)
	defer os.RemoveAll(dir)

	t.Log("applying configuration")

	b.applyConfigSettings(raft.DefaultConfig())

	t.Log("chunking")

	buf := []byte("let's see how this goes, shall we?")

	var logs []*raft.Log
	for i, b := range buf {
		chunkInfo := &raftchunkingtypes.ChunkInfo{
			OpNum:       uint64(32),
			SequenceNum: uint32(i),
			NumChunks:   uint32(len(buf)),
		}
		chunkBytes, err := proto.Marshal(chunkInfo)
		require.NoError(err)

		logs = append(logs, &raft.Log{
			Data:       []byte{b},
			Extensions: chunkBytes,
		})
	}

	t.Log("applying half of the logs")

	// The reason for the skipping is to test out-of-order applies which are
	// theoretically possible
	for i := 0; i < len(logs); i += 2 {
		resp := b.fsm.chunker.Apply(logs[i])
		assert.Nil(resp)
	}

	t.Log("tearing down cluster")
	require.NoError(b.TeardownCluster(nil))

	t.Log("starting new backend")
	backendRaw, err := newRaftBackend(b.conf, b.logger, b.fsm.db)
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

/*
func TestFSM_Chunking_TermChange(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	require := require.New(t)

	fsm, err := New(nil, os.Stderr)
	require.NoError(err)

	req := structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "foo",
		Address:    "127.0.0.1",
		Service: &structs.NodeService{
			ID:      "db",
			Service: "db",
			Tags:    []string{"master"},
			Port:    8000,
		},
		Check: &structs.HealthCheck{
			Node:      "foo",
			CheckID:   "db",
			Name:      "db connectivity",
			Status:    api.HealthPassing,
			ServiceID: "db",
		},
	}
	buf, err := structs.Encode(structs.RegisterRequestType, req)
	require.NoError(err)

	// Only need two chunks to test this
	chunks := [][]byte{
		buf[0:2],
		buf[2:],
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
		resp := fsm.chunker.Apply(log)
		assert.Nil(resp)
	}

	// Now verify the other baseline, that when the term doesn't change we see
	// non-nil. First make the chunker have a clean state, then set the terms
	// to be the same.
	fsm.chunker.RestoreState(nil)
	logs[1].Term = uint64(0)

	// We should see nil only for the first one
	for i, log := range logs {
		resp := fsm.chunker.Apply(log)
		if i == 0 {
			assert.Nil(resp)
		}
		if i == 1 {
			assert.NotNil(resp)
		}
	}
}
*/
