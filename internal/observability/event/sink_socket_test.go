// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/eventlogger"
	"github.com/stretchr/testify/require"
)

// TestNewSocketSink ensures that we validate the input arguments and can create
// the SocketSink if everything goes to plan.
func TestNewSocketSink(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		address        string
		format         string
		opts           []Option
		want           *SocketSink
		wantErr        bool
		expectedErrMsg string
	}{
		"address-empty": {
			address:        "",
			wantErr:        true,
			expectedErrMsg: "address is required: invalid parameter",
		},
		"address-whitespace": {
			address:        "    ",
			wantErr:        true,
			expectedErrMsg: "address is required: invalid parameter",
		},
		"format-empty": {
			address:        "addr",
			format:         "",
			wantErr:        true,
			expectedErrMsg: "format is required: invalid parameter",
		},
		"format-whitespace": {
			address:        "addr",
			format:         "   ",
			wantErr:        true,
			expectedErrMsg: "format is required: invalid parameter",
		},
		"bad-max-duration": {
			address:        "addr",
			format:         "json",
			opts:           []Option{WithMaxDuration("bar")},
			wantErr:        true,
			expectedErrMsg: "unable to parse max duration: invalid parameter: time: invalid duration \"bar\"",
		},
		"happy": {
			address: "wss://foo",
			format:  "json",
			want: &SocketSink{
				requiredFormat: "json",
				address:        "wss://foo",
				socketType:     "tcp",           // defaults to tcp
				maxDuration:    2 * time.Second, // defaults to 2 secs
			},
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := NewSocketSink(tc.address, tc.format, tc.opts...)

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrMsg)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}

// TestSocketSink_Process_cancelledContextWhileWaitingForPermit verifies
// Process returns promptly when the context is canceled while waiting
// to acquire the sink serialization permit.
func TestSocketSink_Process_cancelledContextWhileWaitingForPermit(t *testing.T) {
	sink, err := NewSocketSink("127.0.0.1:0", "json")
	require.NoError(t, err)

	require.NoError(t, sink.acquirePermit(context.Background()))
	t.Cleanup(func() { sink.releasePermit() })

	e := &eventlogger.Event{Formatted: make(map[string][]byte)}
	e.FormattedAs("json", []byte("{\"foo\":\"bar\"}"))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	_, err = sink.Process(ctx, e)
	require.Less(t, time.Since(start), 200*time.Millisecond)
	require.True(t, errors.Is(err, context.Canceled))
}

// TestSocketSink_Process_timeoutWhileWaitingForPermit verifies
// Process returns DeadlineExceeded when the context expires while
// waiting to acquire the sink serialization permit.
func TestSocketSink_Process_timeoutWhileWaitingForPermit(t *testing.T) {
	sink, err := NewSocketSink("127.0.0.1:0", "json")
	require.NoError(t, err)

	require.NoError(t, sink.acquirePermit(context.Background()))

	errCh := make(chan error, 1)
	doneCh := make(chan struct{})

	// Wait cleanup is registered before release cleanup so release happens first.
	t.Cleanup(func() {
		select {
		case <-doneCh:
		case <-time.After(2 * time.Second):
			t.Error("timed out waiting for Process goroutine to exit")
		}
	})
	t.Cleanup(func() { sink.releasePermit() })

	e := &eventlogger.Event{Formatted: make(map[string][]byte)}
	e.FormattedAs("json", []byte("{\"foo\":\"bar\"}"))

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	t.Cleanup(cancel)

	start := time.Now()
	go func() {
		defer close(doneCh)
		_, processErr := sink.Process(ctx, e)
		errCh <- processErr
	}()

	select {
	case err = <-errCh:
		require.Less(t, time.Since(start), 500*time.Millisecond)
		require.True(t, errors.Is(err, context.DeadlineExceeded))
	case <-time.After(500 * time.Millisecond):
		t.Fatal("SocketSink.Process did not return promptly while waiting on permit")
	}
}

// TestSocketSink_ProcessAndReopenAreSerialized verifies Process and
// Reopen remain serialized and do not operate on the socket connection
// concurrently.
func TestSocketSink_ProcessAndReopenAreSerialized(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = listener.Close() })

	acceptDone := make(chan error, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			acceptDone <- acceptErr
			return
		}
		_ = conn.Close()
		acceptDone <- nil
	}()

	sink, err := NewSocketSink(listener.Addr().String(), "json")
	require.NoError(t, err)

	blockConn := &blockingWriteConn{
		writeStarted: make(chan struct{}),
		releaseWrite: make(chan struct{}),
		closed:       make(chan struct{}),
		localAddr:    staticAddr("local"),
		remoteAddr:   staticAddr("remote"),
	}
	sink.connection = blockConn

	e := &eventlogger.Event{Formatted: make(map[string][]byte)}
	e.FormattedAs("json", []byte("{\"foo\":\"bar\"}"))

	processDone := make(chan struct{})
	processErrCh := make(chan error, 1)
	reopenStarted := make(chan struct{})
	reopenDone := make(chan struct{})
	reopenErrCh := make(chan error, 1)
	go func() {
		_, processErr := sink.Process(context.Background(), e)
		processErrCh <- processErr
		close(processDone)
	}()

	var releaseWriteOnce sync.Once
	releaseWrite := func() {
		releaseWriteOnce.Do(func() { close(blockConn.releaseWrite) })
	}
	t.Cleanup(func() {
		releaseWrite()
		select {
		case <-processDone:
		case <-time.After(2 * time.Second):
			t.Error("timed out waiting for Process goroutine cleanup")
		}
		select {
		case <-reopenDone:
		case <-time.After(2 * time.Second):
			t.Error("timed out waiting for Reopen goroutine cleanup")
		}
	})

	<-blockConn.writeStarted

	go func() {
		close(reopenStarted)
		reopenErrCh <- sink.Reopen()
		close(reopenDone)
	}()

	<-reopenStarted
	select {
	case <-reopenDone:
		err = <-reopenErrCh
		t.Fatalf("Reopen returned while Process was blocked in Write: %v", err)
	default:
	}

	releaseWrite()

	select {
	case <-processDone:
		err = <-processErrCh
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Process did not complete after write was released")
	}

	select {
	case <-reopenDone:
		err = <-reopenErrCh
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Reopen did not complete after Process released serialization permit")
	}

	select {
	case acceptErr := <-acceptDone:
		require.NoError(t, acceptErr)
	case <-time.After(2 * time.Second):
		t.Fatal("listener accept goroutine did not exit")
	}
}

type blockingWriteConn struct {
	writeStarted chan struct{}
	releaseWrite chan struct{}
	closed       chan struct{}
	closeOnce    sync.Once
	writeOnce    sync.Once
	localAddr    net.Addr
	remoteAddr   net.Addr
}

func (c *blockingWriteConn) Read(_ []byte) (int, error) {
	return 0, io.EOF
}

func (c *blockingWriteConn) Write(b []byte) (int, error) {
	c.writeOnce.Do(func() { close(c.writeStarted) })
	<-c.releaseWrite
	return len(b), nil
}

func (c *blockingWriteConn) Close() error {
	c.closeOnce.Do(func() { close(c.closed) })
	return nil
}

func (c *blockingWriteConn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *blockingWriteConn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *blockingWriteConn) SetDeadline(_ time.Time) error {
	return nil
}

func (c *blockingWriteConn) SetReadDeadline(_ time.Time) error {
	return nil
}

func (c *blockingWriteConn) SetWriteDeadline(_ time.Time) error {
	return nil
}

type staticAddr string

func (a staticAddr) Network() string { return "tcp" }
func (a staticAddr) String() string  { return string(a) }
