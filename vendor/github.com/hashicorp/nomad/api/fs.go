// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
)

const (
	// OriginStart and OriginEnd are the available parameters for the origin
	// argument when streaming a file. They respectively offset from the start
	// and end of a file.
	OriginStart = "start"
	OriginEnd   = "end"

	// FSLogNameStdout is the name given to the stdout log stream of a task. It
	// can be used when calling AllocFS.Logs as the logType parameter.
	FSLogNameStdout = "stdout"

	// FSLogNameStderr is the name given to the stderr log stream of a task. It
	// can be used when calling AllocFS.Logs as the logType parameter.
	FSLogNameStderr = "stderr"
)

// AllocFileInfo holds information about a file inside the AllocDir
type AllocFileInfo struct {
	Name        string
	IsDir       bool
	Size        int64
	FileMode    string
	ModTime     time.Time
	ContentType string
}

// StreamFrame is used to frame data of a file when streaming
type StreamFrame struct {
	Offset    int64  `json:",omitempty"`
	Data      []byte `json:",omitempty"`
	File      string `json:",omitempty"`
	FileEvent string `json:",omitempty"`
}

// IsHeartbeat returns if the frame is a heartbeat frame
func (s *StreamFrame) IsHeartbeat() bool {
	return len(s.Data) == 0 && s.FileEvent == "" && s.File == "" && s.Offset == 0
}

// AllocFS is used to introspect an allocation directory on a Nomad client
type AllocFS struct {
	client *Client
}

// AllocFS returns an handle to the AllocFS endpoints
func (c *Client) AllocFS() *AllocFS {
	return &AllocFS{client: c}
}

// List is used to list the files at a given path of an allocation directory.
// Note: for cluster topologies where API consumers don't have network access to
// Nomad clients, set api.ClientConnTimeout to a small value (ex 1ms) to avoid
// long pauses on this API call.
func (a *AllocFS) List(alloc *Allocation, path string, q *QueryOptions) ([]*AllocFileInfo, *QueryMeta, error) {
	if q == nil {
		q = &QueryOptions{}
	}
	if q.Params == nil {
		q.Params = make(map[string]string)
	}
	q.Params["path"] = path

	var resp []*AllocFileInfo
	qm, err := a.client.query(fmt.Sprintf("/v1/client/fs/ls/%s", alloc.ID), &resp, q)
	if err != nil {
		return nil, nil, err
	}

	return resp, qm, nil
}

// Stat is used to stat a file at a given path of an allocation directory.
// Note: for cluster topologies where API consumers don't have network access to
// Nomad clients, set api.ClientConnTimeout to a small value (ex 1ms) to avoid
// long pauses on this API call.
func (a *AllocFS) Stat(alloc *Allocation, path string, q *QueryOptions) (*AllocFileInfo, *QueryMeta, error) {
	if q == nil {
		q = &QueryOptions{}
	}
	if q.Params == nil {
		q.Params = make(map[string]string)
	}

	q.Params["path"] = path

	var resp AllocFileInfo
	qm, err := a.client.query(fmt.Sprintf("/v1/client/fs/stat/%s", alloc.ID), &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// ReadAt is used to read bytes at a given offset until limit at the given path
// in an allocation directory. If limit is <= 0, there is no limit.
// Note: for cluster topologies where API consumers don't have network access to
// Nomad clients, set api.ClientConnTimeout to a small value (ex 1ms) to avoid
// long pauses on this API call.
func (a *AllocFS) ReadAt(alloc *Allocation, path string, offset int64, limit int64, q *QueryOptions) (io.ReadCloser, error) {
	reqPath := fmt.Sprintf("/v1/client/fs/readat/%s", alloc.ID)

	return queryClientNode(a.client, alloc, reqPath, q,
		func(q *QueryOptions) {
			q.Params["path"] = path
			q.Params["offset"] = strconv.FormatInt(offset, 10)
			q.Params["limit"] = strconv.FormatInt(limit, 10)
		})
}

// Cat is used to read contents of a file at the given path in an allocation
// directory.
// Note: for cluster topologies where API consumers don't have network access to
// Nomad clients, set api.ClientConnTimeout to a small value (ex 1ms) to avoid
// long pauses on this API call.
func (a *AllocFS) Cat(alloc *Allocation, path string, q *QueryOptions) (io.ReadCloser, error) {
	reqPath := fmt.Sprintf("/v1/client/fs/cat/%s", alloc.ID)
	return queryClientNode(a.client, alloc, reqPath, q,
		func(q *QueryOptions) {
			q.Params["path"] = path
		})
}

// Stream streams the content of a file blocking on EOF.
// The parameters are:
// * path: path to file to stream.
// * offset: The offset to start streaming data at.
// * origin: Either "start" or "end" and defines from where the offset is applied.
// * cancel: A channel that when closed, streaming will end.
//
// The return value is a channel that will emit StreamFrames as they are read.
//
// Note: for cluster topologies where API consumers don't have network access to
// Nomad clients, set api.ClientConnTimeout to a small value (ex 1ms) to avoid
// long pauses on this API call.
func (a *AllocFS) Stream(alloc *Allocation, path, origin string, offset int64,
	cancel <-chan struct{}, q *QueryOptions) (<-chan *StreamFrame, <-chan error) {

	errCh := make(chan error, 1)

	reqPath := fmt.Sprintf("/v1/client/fs/stream/%s", alloc.ID)
	r, err := queryClientNode(a.client, alloc, reqPath, q,
		func(q *QueryOptions) {
			q.Params["path"] = path
			q.Params["offset"] = strconv.FormatInt(offset, 10)
			q.Params["origin"] = origin
		})
	if err != nil {
		errCh <- err
		return nil, errCh
	}

	// Create the output channel
	frames := make(chan *StreamFrame, 10)

	go func() {
		// Close the body
		defer r.Close()

		// Create a decoder
		dec := json.NewDecoder(r)

		for {
			// Check if we have been cancelled
			select {
			case <-cancel:
				return
			default:
			}

			// Decode the next frame
			var frame StreamFrame
			if err := dec.Decode(&frame); err != nil {
				errCh <- err
				close(frames)
				return
			}

			// Discard heartbeat frames
			if frame.IsHeartbeat() {
				continue
			}

			frames <- &frame
		}
	}()

	return frames, errCh
}

func queryClientNode(c *Client, alloc *Allocation, reqPath string, q *QueryOptions, customizeQ func(*QueryOptions)) (io.ReadCloser, error) {
	nodeClient, _ := c.GetNodeClientWithTimeout(alloc.NodeID, ClientConnTimeout, q)

	if q == nil {
		q = &QueryOptions{}
	}
	if q.Params == nil {
		q.Params = make(map[string]string)
	}
	if customizeQ != nil {
		customizeQ(q)
	}

	var r io.ReadCloser
	var err error

	if nodeClient != nil {
		r, err = nodeClient.rawQuery(reqPath, q)
		if _, ok := err.(net.Error); err != nil && !ok {
			// found a non networking error talking to client directly
			return nil, err
		}

	}

	// failed to query node, access through server directly
	// or network error when talking to the client directly
	if r == nil {
		return c.rawQuery(reqPath, q)
	}

	return r, err
}

// Logs streams the content of a tasks logs blocking on EOF.
// The parameters are:
// * allocation: the allocation to stream from.
// * follow: Whether the logs should be followed.
// * task: the tasks name to stream logs for.
// * logType: Either "stdout" or "stderr"
// * origin: Either "start" or "end" and defines from where the offset is applied.
// * offset: The offset to start streaming data at.
// * cancel: A channel that when closed, streaming will end.
//
// The return value is a channel that will emit StreamFrames as they are read.
// The chan will be closed when follow=false and the end of the file is
// reached.
//
// Unexpected (non-EOF) errors will be sent on the error chan.
//
// Note: for cluster topologies where API consumers don't have network access to
// Nomad clients, set api.ClientConnTimeout to a small value (ex 1ms) to avoid
// long pauses on this API call.
func (a *AllocFS) Logs(alloc *Allocation, follow bool, task, logType, origin string,
	offset int64, cancel <-chan struct{}, q *QueryOptions) (<-chan *StreamFrame, <-chan error) {

	errCh := make(chan error, 1)

	reqPath := fmt.Sprintf("/v1/client/fs/logs/%s", alloc.ID)
	r, err := queryClientNode(a.client, alloc, reqPath, q,
		func(q *QueryOptions) {
			q.Params["follow"] = strconv.FormatBool(follow)
			q.Params["task"] = task
			q.Params["type"] = logType
			q.Params["origin"] = origin
			q.Params["offset"] = strconv.FormatInt(offset, 10)
		})
	if err != nil {
		errCh <- err
		return nil, errCh
	}

	// Create the output channel
	frames := make(chan *StreamFrame, 10)

	go func() {
		// Close the body
		defer r.Close()

		// Create a decoder
		dec := json.NewDecoder(r)

		for {
			// Check if we have been cancelled
			select {
			case <-cancel:
				close(frames)
				return
			default:
			}

			// Decode the next frame
			var frame StreamFrame
			if err := dec.Decode(&frame); err != nil {
				if err == io.EOF || err == io.ErrClosedPipe {
					close(frames)
				} else {
					buf, err2 := io.ReadAll(dec.Buffered())
					if err2 != nil {
						errCh <- fmt.Errorf("failed to decode and failed to read buffered data: %w", multierror.Append(err, err2))
					} else {
						errCh <- fmt.Errorf("failed to decode log endpoint response as JSON: %q", buf)
					}
				}
				return
			}

			// Discard heartbeat frames
			if frame.IsHeartbeat() {
				continue
			}

			frames <- &frame
		}
	}()

	return frames, errCh
}

// FrameReader is used to convert a stream of frames into a read closer.
type FrameReader struct {
	frames   <-chan *StreamFrame
	errCh    <-chan error
	cancelCh chan struct{}

	closedLock sync.Mutex
	closed     bool

	unblockTime time.Duration

	frame       *StreamFrame
	frameOffset int

	byteOffset int
}

// NewFrameReader takes a channel of frames and returns a FrameReader which
// implements io.ReadCloser
func NewFrameReader(frames <-chan *StreamFrame, errCh <-chan error, cancelCh chan struct{}) *FrameReader {
	return &FrameReader{
		frames:   frames,
		errCh:    errCh,
		cancelCh: cancelCh,
	}
}

// SetUnblockTime sets the time to unblock and return zero bytes read. If the
// duration is unset or is zero or less, the read will block until data is read.
func (f *FrameReader) SetUnblockTime(d time.Duration) {
	f.unblockTime = d
}

// Offset returns the offset into the stream.
func (f *FrameReader) Offset() int {
	return f.byteOffset
}

// Read reads the data of the incoming frames into the bytes buffer. Returns EOF
// when there are no more frames.
func (f *FrameReader) Read(p []byte) (n int, err error) {
	f.closedLock.Lock()
	closed := f.closed
	f.closedLock.Unlock()
	if closed {
		return 0, io.EOF
	}

	if f.frame == nil {
		var unblock <-chan time.Time
		if f.unblockTime.Nanoseconds() > 0 {
			unblock = time.After(f.unblockTime)
		}

		select {
		case frame, ok := <-f.frames:
			if !ok {
				return 0, io.EOF
			}
			f.frame = frame

			// Store the total offset into the file
			f.byteOffset = int(f.frame.Offset)
		case <-unblock:
			return 0, nil
		case err := <-f.errCh:
			return 0, err
		case <-f.cancelCh:
			return 0, io.EOF
		}
	}

	// Copy the data out of the frame and update our offset
	n = copy(p, f.frame.Data[f.frameOffset:])
	f.frameOffset += n

	// Clear the frame and its offset once we have read everything
	if len(f.frame.Data) == f.frameOffset {
		f.frame = nil
		f.frameOffset = 0
	}

	return n, nil
}

// Close cancels the stream of frames
func (f *FrameReader) Close() error {
	f.closedLock.Lock()
	defer f.closedLock.Unlock()
	if f.closed {
		return nil
	}

	close(f.cancelCh)
	f.closed = true
	return nil
}
