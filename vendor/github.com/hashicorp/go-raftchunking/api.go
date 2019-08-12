package raftchunking

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-raftchunking/types"
	"github.com/hashicorp/raft"
)

var (
	// ChunkSize is the threshold used for breaking a large value into chunks.
	// Defaults to the suggested max data size for the raft library.
	ChunkSize = raft.SuggestedMaxDataSize
)

// errorFuture is used to return a static error.
type errorFuture struct {
	err error
}

func (e errorFuture) Error() error {
	return e.err
}

func (e errorFuture) Response() interface{} {
	return nil
}

func (e errorFuture) Index() uint64 {
	return 0
}

// multiFuture is a future specialized for the chunking case. It contains some
// number of other futures in the order in which data was chunked and sent to
// apply.
type multiFuture []raft.ApplyFuture

// Error will return only when all Error functions in the contained futures
// return, in order.
func (m multiFuture) Error() error {
	for _, v := range m {
		if err := v.Error(); err != nil {
			return err
		}
	}

	return nil
}

// Index returns the index of the last chunk. Since required behavior is to not
// call this until Error is called, the last Index will correspond to the Apply
// of the final chunk.
func (m multiFuture) Index() uint64 {
	// This shouldn't happen but need an escape hatch
	if len(m) == 0 {
		return 0
	}

	return m[len(m)-1].Index()
}

// Response returns the response from underlying Apply of the last chunk.
func (m multiFuture) Response() interface{} {
	// This shouldn't happen but need an escape hatch
	if len(m) == 0 {
		return nil
	}

	return m[len(m)-1].Response()
}

type ApplyFunc func(raft.Log, time.Duration) raft.ApplyFuture

// ChunkingApply takes in a byte slice and chunks into ChunkSize (or less if
// EOF) chunks, calling Apply on each. It requires a corresponding wrapper
// around the FSM to handle reconstructing on the other end. Timeout will be the
// timeout for each individual operation, not total. The return value is a
// future whose Error() will return only when all underlying Apply futures have
// had Error() return. Note that any error indicates that the entire operation
// will not be applied, assuming the correct FSM wrapper is used. If extensions
// is passed in, it will be set as the Extensions value on the Apply once all
// chunks are received.
func ChunkingApply(cmd, extensions []byte, timeout time.Duration, applyFunc ApplyFunc) raft.ApplyFuture {
	// Generate a random op num via 64 random bits. These only have to be
	// unique across _in flight_ chunk operations until a Term changes so
	// should be fine.
	rb := make([]byte, 8)
	n, err := rand.Read(rb)
	if err != nil {
		return errorFuture{err: err}
	}
	if n != 8 {
		return errorFuture{err: fmt.Errorf("expected to read %d bytes for op num, read %d", 8, n)}
	}
	opNum := binary.BigEndian.Uint64(rb)

	var logs []raft.Log
	var byteChunks [][]byte
	var mf multiFuture

	// We break into chunks first so that we know how many chunks there will be
	// to put in NumChunks in the extensions info. This could probably be a bit
	// more efficient by just reslicing but doing it this way is a bit easier
	// for others to follow/track and in this kind of operation this won't be
	// the slow part anyways.
	reader := bytes.NewReader(cmd)
	remain := reader.Len()
	for {
		if remain <= 0 {
			break
		}

		if remain > ChunkSize {
			remain = ChunkSize
		}

		b := make([]byte, remain)
		n, err := reader.Read(b)
		if err != nil && err != io.EOF {
			return errorFuture{err: err}
		}
		if n != remain {
			return errorFuture{err: fmt.Errorf("expected to read %d bytes from buf, read %d", remain, n)}
		}

		byteChunks = append(byteChunks, b)
		remain = reader.Len()
	}

	// Create the underlying chunked logs
	for i, chunk := range byteChunks {
		chunkInfo := &types.ChunkInfo{
			OpNum:       opNum,
			SequenceNum: uint32(i),
			NumChunks:   uint32(len(byteChunks)),
		}

		// If extensions were passed in attach them to the last chunk so it
		// will go through Apply at the end.
		if i == len(byteChunks)-1 {
			chunkInfo.NextExtensions = extensions
		}

		chunkBytes, err := proto.Marshal(chunkInfo)
		if err != nil {
			return errorFuture{err: errwrap.Wrapf("error marshaling chunk info: {{err}}", err)}
		}
		logs = append(logs, raft.Log{
			Data:       chunk,
			Extensions: chunkBytes,
		})
	}

	for _, log := range logs {
		mf = append(mf, applyFunc(log, timeout))
	}

	return mf
}
