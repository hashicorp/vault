package raftchunking

import (
	"errors"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-raftchunking/types"
	"github.com/hashicorp/raft"
)

var (
	ErrTermMismatch = errors.New("term mismatch during reconstruction of chunks, please resubmit")

	ErrInvalidOpNum = errors.New("no op num found when reconstructing chunks")

	ErrNoExistingChunks = errors.New("no existing chunks but non-zero sequence num")

	ErrSequenceNumberMismatch = errors.New("sequence number skipped")

	ErrMissingChunk = errors.New("missing sequence number during reconstruction")
)

type chunkInfo struct {
	term   uint64
	seqNum uint32
	data   []byte
}

var _ raft.FSM = (*ChunkingFSM)(nil)
var _ raft.ConfigurationStore = (*ChunkingConfigurationStore)(nil)

type ChunkingFSM struct {
	underlying raft.FSM
	opMap      map[uint64][]chunkInfo
}

type ChunkingConfigurationStore struct {
	*ChunkingFSM
	underlyingConfigurationStore raft.ConfigurationStore
}

func NewChunkingFSM(underlying raft.FSM) raft.FSM {
	return &ChunkingFSM{
		underlying: underlying,
		opMap:      make(map[uint64][]chunkInfo),
	}
}

func NewChunkingConfigurationStore(underlying raft.ConfigurationStore) raft.ConfigurationStore {
	return &ChunkingConfigurationStore{
		ChunkingFSM: &ChunkingFSM{
			underlying: underlying,
			opMap:      make(map[uint64][]chunkInfo),
		},
		underlyingConfigurationStore: underlying,
	}
}

// Apply applies the log, handling chunking as needed. The return value will
// either be an error or whatever is returned from the underlying Apply.
func (c *ChunkingFSM) Apply(l *raft.Log) interface{} {
	// Not chunking or wrong type, pass through
	if l.Type != raft.LogCommand || l.Extensions == nil {
		return c.underlying.Apply(l)
	}

	var ci types.ChunkInfo
	if err := proto.Unmarshal(l.Extensions, &ci); err != nil {
		return errwrap.Wrapf("error unmarshaling chunk info: {{err}}", err)
	}

	opNum := ci.OpNum
	seqNum := ci.SequenceNum

	if opNum == 0 {
		return ErrInvalidOpNum
	}

	// Look up existing chunks
	chunks, ok := c.opMap[opNum]
	if !ok {
		if seqNum != 0 {
			return ErrNoExistingChunks
		}
	}

	// Do early detection of a loss or other problem
	if int(seqNum) != len(chunks) {
		delete(c.opMap, opNum)
		return ErrSequenceNumberMismatch
	}

	chunks = append(chunks, chunkInfo{
		term:   l.Term,
		seqNum: seqNum,
		data:   l.Data,
	})

	if ci.SequenceNum == ci.NumChunks-1 {
		// Run through and reconstruct the data
		finalData := make([]byte, 0, len(chunks)*raft.SuggestedMaxDataSize)

		var term uint64
		for i, chunk := range chunks {
			if i == 0 {
				term = chunk.term
			}
			if chunk.seqNum != uint32(i) {
				delete(c.opMap, opNum)
				return ErrMissingChunk
			}
			if chunk.term != term {
				delete(c.opMap, opNum)
				return ErrTermMismatch
			}
			finalData = append(finalData, chunk.data...)
		}

		// Use the latest log's values with the final data
		logToApply := &raft.Log{
			Index:      l.Index,
			Term:       l.Term,
			Type:       l.Type,
			Data:       finalData,
			Extensions: ci.NextExtensions,
		}

		delete(c.opMap, opNum)
		return c.Apply(logToApply)
	}

	// Otherwise, re-add to map and return
	c.opMap[opNum] = chunks
	return nil
}

func (c *ChunkingFSM) Snapshot() (raft.FSMSnapshot, error) {
	return c.underlying.Snapshot()
}

func (c *ChunkingFSM) Restore(rc io.ReadCloser) error {
	return c.underlying.Restore(rc)
}

// Note: this is used in tests via the Raft package test helper functions, even
// if it's not used in client code
func (c *ChunkingFSM) Underlying() raft.FSM {
	return c.underlying
}

func (c *ChunkingConfigurationStore) StoreConfiguration(index uint64, configuration raft.Configuration) {
	c.underlyingConfigurationStore.StoreConfiguration(index, configuration)
}
