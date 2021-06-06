package raftchunking

import "github.com/mitchellh/copystructure"

type ChunkStorage interface {
	// StoreChunk stores Data from ChunkInfo according to the other metadata
	// (OpNum, SeqNum). The bool returns whether or not all chunks have been
	// received, as in, the number of non-nil chunks is the same as NumChunks.
	StoreChunk(*ChunkInfo) (bool, error)

	// FinalizeOp gets all chunks for an op number and then removes the chunk
	// info for that op from the store. It should only be called when
	// StoreChunk for a given op number returns true but should be safe to call
	// at any time; clearing an op can be accomplished by calling this function
	// and ignoring the non-error result.
	FinalizeOp(uint64) ([]*ChunkInfo, error)

	// GetState gets all currently tracked ops, for snapshotting
	GetChunks() (ChunkMap, error)

	// RestoreChunks restores the current FSM state from a map
	RestoreChunks(ChunkMap) error
}

type State struct {
	ChunkMap ChunkMap
}

// ChunkInfo holds chunk information
type ChunkInfo struct {
	OpNum       uint64
	SequenceNum uint32
	NumChunks   uint32
	Term        uint64
	Data        []byte
}

// ChunkMap represents a set of data chunks. We use ChunkInfo with Data instead
// of bare []byte in case there is a need to extend this info later.
type ChunkMap map[uint64][]*ChunkInfo

// InmemChunkStorage satisfies ChunkStorage using an in-memory-only tracking
// method.
type InmemChunkStorage struct {
	chunks ChunkMap
}

func NewInmemChunkStorage() *InmemChunkStorage {
	return &InmemChunkStorage{
		chunks: make(ChunkMap),
	}
}

func (i *InmemChunkStorage) StoreChunk(chunk *ChunkInfo) (bool, error) {
	chunks, ok := i.chunks[chunk.OpNum]
	if !ok {
		chunks = make([]*ChunkInfo, chunk.NumChunks)
		i.chunks[chunk.OpNum] = chunks
	}

	chunks[chunk.SequenceNum] = chunk

	for _, c := range chunks {
		// Check for nil, but also check data length in case it ends up
		// unmarshaling weirdly for some reason where it makes a new struct
		// instead of keeping the pointer nil
		if c == nil || len(c.Data) == 0 {
			// Not done yet, so return
			return false, nil
		}
	}

	return true, nil
}

func (i *InmemChunkStorage) FinalizeOp(opNum uint64) ([]*ChunkInfo, error) {
	ret := i.chunks[opNum]
	delete(i.chunks, opNum)
	return ret, nil
}

func (i *InmemChunkStorage) GetChunks() (ChunkMap, error) {
	ret, err := copystructure.Copy(i.chunks)
	if err != nil {
		return nil, err
	}
	return ret.(ChunkMap), nil
}

func (i *InmemChunkStorage) RestoreChunks(chunks ChunkMap) error {
	// If passed in explicit emptiness, set state to empty
	if chunks == nil || len(chunks) == 0 {
		i.chunks = make(ChunkMap)
		return nil
	}

	chunksCopy, err := copystructure.Copy(chunks)
	if err != nil {
		return err
	}
	i.chunks = chunksCopy.(ChunkMap)
	return nil
}
