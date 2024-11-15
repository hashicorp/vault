// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// heartbeatInterval is the amount of time to wait between sending heartbeats
	// during an exec streaming operation
	heartbeatInterval = 10 * time.Second
)

type execSession struct {
	client  *Client
	alloc   *Allocation
	job     string
	task    string
	tty     bool
	command []string
	action  string

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	terminalSizeCh <-chan TerminalSize

	q *QueryOptions
}

func (s *execSession) run(ctx context.Context) (exitCode int, err error) {
	ctx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	conn, err := s.startConnection()
	if err != nil {
		return -2, err
	}
	defer conn.Close()

	sendErrCh := s.startTransmit(ctx, conn)
	exitCh, recvErrCh := s.startReceiving(ctx, conn)

	for {
		select {
		case <-ctx.Done():
			return -2, ctx.Err()
		case exitCode := <-exitCh:
			return exitCode, nil
		case recvErr := <-recvErrCh:
			// drop websocket code, not relevant to user
			if wsErr, ok := recvErr.(*websocket.CloseError); ok && wsErr.Text != "" {
				return -2, errors.New(wsErr.Text)
			}

			return -2, recvErr
		case sendErr := <-sendErrCh:
			return -2, fmt.Errorf("failed to send input: %w", sendErr)
		}
	}
}

func (s *execSession) startConnection() (*websocket.Conn, error) {
	// First, attempt to connect to the node directly, but may fail due to network isolation
	// and network errors.  Fallback to using server-side forwarding instead.
	nodeClient, err := s.client.GetNodeClientWithTimeout(s.alloc.NodeID, ClientConnTimeout, s.q)
	if err == NodeDownErr {
		return nil, NodeDownErr
	}

	q := s.q
	if q == nil {
		q = &QueryOptions{}
	}
	if q.Params == nil {
		q.Params = make(map[string]string)
	}

	commandBytes, err := json.Marshal(s.command)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal command: %W", err)
	}

	q.Params["tty"] = strconv.FormatBool(s.tty)
	q.Params["task"] = s.task
	q.Params["command"] = string(commandBytes)
	reqPath := fmt.Sprintf("/v1/client/allocation/%s/exec", s.alloc.ID)

	if s.action != "" {
		q.Params["action"] = s.action
		q.Params["allocID"] = s.alloc.ID
		q.Params["group"] = s.alloc.TaskGroup
		reqPath = fmt.Sprintf("/v1/job/%s/action", url.PathEscape(s.job))
	}

	var conn *websocket.Conn

	if nodeClient != nil {
		conn, _, _ = nodeClient.websocket(reqPath, q) //nolint:bodyclose // gorilla/websocket Dialer.DialContext() does not require the body to be closed.
	}

	if conn == nil {
		conn, _, err = s.client.websocket(reqPath, q) //nolint:bodyclose // gorilla/websocket Dialer.DialContext() does not require the body to be closed.
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

func (s *execSession) startTransmit(ctx context.Context, conn *websocket.Conn) <-chan error {

	// FIXME: Handle websocket send errors.
	// Currently, websocket write failures are dropped. As sending and
	// receiving are running concurrently, it's expected that some send
	// requests may fail with connection errors when connection closes.
	// Connection errors should surface in the receive paths already,
	// but I'm unsure about one-sided communication errors.
	var sendLock sync.Mutex
	send := func(v *ExecStreamingInput) {
		sendLock.Lock()
		defer sendLock.Unlock()

		conn.WriteJSON(v)
	}

	errCh := make(chan error, 4)

	// propagate stdin
	go func() {

		bytes := make([]byte, 2048)
		for {
			if ctx.Err() != nil {
				return
			}

			input := ExecStreamingInput{Stdin: &ExecStreamingIOOperation{}}

			n, err := s.stdin.Read(bytes)

			// always send data if we read some
			if n != 0 {
				input.Stdin.Data = bytes[:n]
				send(&input)
			}

			// then handle error
			if err == io.EOF {
				// if n != 0, send data and we'll get n = 0 on next read
				if n == 0 {
					input.Stdin.Close = true
					send(&input)
					return
				}
			} else if err != nil {
				errCh <- err
				return
			}
		}
	}()

	// propagate terminal sizing updates
	go func() {
		for {
			resizeInput := ExecStreamingInput{}

			select {
			case <-ctx.Done():
				return
			case size, ok := <-s.terminalSizeCh:
				if !ok {
					return
				}
				resizeInput.TTYSize = &size
				send(&resizeInput)
			}

		}
	}()

	// send a heartbeat every 10 seconds
	go func() {
		t := time.NewTimer(heartbeatInterval)
		defer t.Stop()

		for {
			t.Reset(heartbeatInterval)

			select {
			case <-ctx.Done():
				return
			case <-t.C:
				// heartbeat message
				send(&execStreamingInputHeartbeat)
			}
		}
	}()

	return errCh
}

func (s *execSession) startReceiving(ctx context.Context, conn *websocket.Conn) (<-chan int, <-chan error) {
	exitCodeCh := make(chan int, 1)
	errCh := make(chan error, 1)

	go func() {
		for ctx.Err() == nil {

			// Decode the next frame
			var frame ExecStreamingOutput
			err := conn.ReadJSON(&frame)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				errCh <- fmt.Errorf("websocket closed before receiving exit code: %w", err)
				return
			} else if err != nil {
				errCh <- err
				return
			}

			switch {
			case frame.Stdout != nil:
				if len(frame.Stdout.Data) != 0 {
					s.stdout.Write(frame.Stdout.Data)
				}
				// don't really do anything if stdout is closing
			case frame.Stderr != nil:
				if len(frame.Stderr.Data) != 0 {
					s.stderr.Write(frame.Stderr.Data)
				}
				// don't really do anything if stderr is closing
			case frame.Exited && frame.Result != nil:
				exitCodeCh <- frame.Result.ExitCode
				return
			default:
				// noop - heartbeat
			}

		}

	}()

	return exitCodeCh, errCh
}
