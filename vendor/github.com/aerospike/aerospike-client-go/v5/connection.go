// Copyright 2014-2021 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	"compress/zlib"
	"crypto/tls"
	"io"
	"net"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aerospike/aerospike-client-go/v5/logger"
	"github.com/aerospike/aerospike-client-go/v5/types"
)

// DefaultBufferSize specifies the initial size of the connection buffer when it is created.
// If not big enough (as big as the average record), it will be reallocated to size again
// which will be more expensive.
var DefaultBufferSize = 64 * 1024 // 64 KiB

// bufPool reuses the data buffers to remove pressure from
// the allocator and the GC during connection churns.
var bufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, DefaultBufferSize)
	},
}

// Connection represents a connection with a timeout.
type Connection struct {
	node *Node

	// timeouts
	socketTimeout time.Duration
	deadline      time.Time

	// duration after which connection is considered idle
	idleTimeout  time.Duration
	idleDeadline time.Time

	// connection object
	conn net.Conn

	// to avoid having a buffer pool and contention
	dataBuffer []byte

	compressed bool
	inflater   io.ReadCloser
	// inflater may consume more bytes than required.
	// LimitReader is used to avoid that problem.
	limitReader *io.LimitedReader

	closer sync.Once
}

// makes sure that the connection is closed eventually, even if it is not consumed
func connectionFinalizer(c *Connection) {
	c.Close()
}

// errToAerospikeErr will convert golang's net and io errors into *AerospikeError
// If the errors is nil, nil be returned. If conn is not nil, its node value
// will be set for the error.
func errToAerospikeErr(conn *Connection, err error) (aerr Error) {
	if err == nil {
		return nil
	}

	if terr, ok := err.(net.Error); ok {
		if terr.Timeout() {
			aerr = newErrorAndWrap(err, types.TIMEOUT)
		} else {
			aerr = newErrorAndWrap(err, types.NETWORK_ERROR)
		}
	} else {
		aerr = newErrorAndWrap(err, types.NETWORK_ERROR)
	}

	// set node if exists
	if conn != nil {
		aerr.setNode(conn.node)
	}

	return aerr
}

// newConnection creates a connection on the network and returns the pointer
// A minimum timeout of 2 seconds will always be applied.
// If the connection is not established in the specified timeout,
// an error will be returned
func newConnection(address string, timeout time.Duration) (*Connection, Error) {
	newConn := &Connection{dataBuffer: bufPool.Get().([]byte)}
	runtime.SetFinalizer(newConn, connectionFinalizer)

	// don't wait indefinitely
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		logger.Logger.Debug("Connection to address `%s` failed to establish with error: %s", address, err.Error())
		return nil, errToAerospikeErr(nil, err)
	}
	newConn.conn = conn
	newConn.limitReader = &io.LimitedReader{R: conn, N: 0}

	// set timeout at the last possible moment
	if err := newConn.SetTimeout(time.Now().Add(timeout), timeout); err != nil {
		newConn.Close()
		return nil, err
	}

	return newConn, nil
}

// NewConnection creates a TLS connection on the network and returns the pointer.
// A minimum timeout of 2 seconds will always be applied.
// If the connection is not established in the specified timeout,
// an error will be returned
func NewConnection(policy *ClientPolicy, host *Host) (*Connection, Error) {
	address := net.JoinHostPort(host.Name, strconv.Itoa(host.Port))
	conn, err := newConnection(address, policy.Timeout)
	if err != nil {
		return nil, err
	}

	if policy.TlsConfig == nil {
		return conn, nil
	}

	// Use version dependent clone function to clone the config
	tlsConfig := policy.TlsConfig.Clone()
	tlsConfig.ServerName = host.TLSName

	sconn := tls.Client(conn.conn, tlsConfig)
	if err := sconn.Handshake(); err != nil {
		nerr := newWrapNetworkError(err)
		if cerr := sconn.Close(); cerr != nil {
			logger.Logger.Debug("Closing connection after handshake error failed: %s", cerr.Error())
			nerr = chainErrors(newWrapNetworkError(cerr), nerr)
		}
		return nil, nerr
	}

	if host.TLSName != "" && !tlsConfig.InsecureSkipVerify {
		if err := sconn.VerifyHostname(host.TLSName); err != nil {
			nerr := newWrapNetworkError(err)
			if cerr := sconn.Close(); cerr != nil {
				logger.Logger.Debug("Closing connection after VerifyHostName error failed: %s", cerr.Error())
				nerr = chainErrors(newWrapNetworkError(cerr), nerr)
			}
			logger.Logger.Error("Connection to address `%s` failed to establish with error: %s", address, err.Error())
			return nil, nerr
		}
	}

	conn.conn = sconn
	return conn, nil
}

// Write writes the slice to the connection buffer.
func (ctn *Connection) Write(buf []byte) (total int, aerr Error) {
	var err error

	// make sure all bytes are written
	// Don't worry about the loop, timeout has been set elsewhere
	if err = ctn.updateDeadline(); err == nil {
		if total, err = ctn.conn.Write(buf); err == nil {
			return total, nil
		}

		// If all bytes are written, ignore any potential error
		// The error will bubble up on the next network io if it matters.
		if total == len(buf) {
			return total, nil
		}
	}

	aerr = chainErrors(errToAerospikeErr(ctn, err), aerr)

	if ctn.node != nil {
		ctn.node.incrErrorCount()
		atomic.AddInt64(&ctn.node.stats.ConnectionsFailed, 1)
	}

	ctn.Close()

	return total, aerr
}

// Read reads from connection buffer to the provided slice.
func (ctn *Connection) Read(buf []byte, length int) (total int, aerr Error) {
	var err error

	// if all bytes are not read, retry until successful
	// Don't worry about the loop; we've already set the timeout elsewhere
	for total < length {
		var r int
		if err = ctn.updateDeadline(); err != nil {
			break
		}

		if !ctn.compressed {
			r, err = ctn.conn.Read(buf[total:length])
		} else {
			r, err = ctn.inflater.Read(buf[total:length])
			if err == io.EOF && total+r == length {
				ctn.compressed = false
				err = ctn.inflater.Close()
			}
		}
		total += r
		if err != nil {
			break
		}
	}

	if total == length {
		// If all required bytes are read, ignore any potential error.
		// The error will bubble up on the next network io if it matters.
		return total, nil
	}

	aerr = chainErrors(errToAerospikeErr(ctn, err), aerr)

	if ctn.node != nil {
		ctn.node.incrErrorCount()
		atomic.AddInt64(&ctn.node.stats.ConnectionsFailed, 1)
	}

	ctn.Close()

	return total, aerr
}

// IsConnected returns true if the connection is not closed yet.
func (ctn *Connection) IsConnected() bool {
	return ctn.conn != nil
}

// updateDeadline sets connection timeout for both read and write operations.
// this function is called before each read and write operation. If deadline has passed,
// the function will return a TIMEOUT error.
func (ctn *Connection) updateDeadline() Error {
	now := time.Now()
	var socketDeadline time.Time
	if ctn.deadline.IsZero() {
		if ctn.socketTimeout > 0 {
			socketDeadline = now.Add(ctn.socketTimeout)
		}
	} else {
		if now.After(ctn.deadline) {
			return newError(types.TIMEOUT)
		}
		if ctn.socketTimeout == 0 {
			socketDeadline = ctn.deadline
		} else {
			tDeadline := now.Add(ctn.socketTimeout)
			if tDeadline.After(ctn.deadline) {
				socketDeadline = ctn.deadline
			} else {
				socketDeadline = tDeadline
			}
		}

		// floor to a millisecond to avoid too short timeouts
		if socketDeadline.Sub(now) < time.Millisecond {
			socketDeadline = now.Add(time.Millisecond)
		}
	}

	if err := ctn.conn.SetDeadline(socketDeadline); err != nil {
		if ctn.node != nil {
			atomic.AddInt64(&ctn.node.stats.ConnectionsFailed, 1)
		}
		return errToAerospikeErr(ctn, err)
	}

	return nil
}

// SetTimeout sets connection timeout for both read and write operations.
func (ctn *Connection) SetTimeout(deadline time.Time, socketTimeout time.Duration) Error {
	ctn.deadline = deadline
	ctn.socketTimeout = socketTimeout

	return nil
}

// Close closes the connection
func (ctn *Connection) Close() {
	ctn.closer.Do(func() {
		if ctn != nil && ctn.conn != nil {
			// deregister
			if ctn.node != nil {
				ctn.node.connectionCount.DecrementAndGet()
				atomic.AddInt64(&ctn.node.stats.ConnectionsClosed, 1)
			}

			if err := ctn.conn.Close(); err != nil {
				logger.Logger.Warn(err.Error())
			}
			ctn.conn = nil

			// put the data buffer back in the pool in case it gets used again
			if len(ctn.dataBuffer) >= DefaultBufferSize && len(ctn.dataBuffer) <= MaxBufferSize {
				bufPool.Put(ctn.dataBuffer)
			}

			ctn.dataBuffer = nil
			ctn.node = nil
		}
	})
}

// Login will send authentication information to the server.
func (ctn *Connection) login(policy *ClientPolicy, hashedPassword []byte, sessionInfo *sessionInfo) Error {
	// need to authenticate
	if policy.RequiresAuthentication() {
		var err Error
		command := newLoginCommand(ctn.dataBuffer)

		if !sessionInfo.isValid() {
			err = command.login(policy, ctn, hashedPassword)
		} else {
			err = command.authenticateViaToken(policy, ctn, sessionInfo.token)
			if err != nil && err.Matches(types.INVALID_CREDENTIAL, types.EXPIRED_SESSION) {
				// invalidate the token
				if ctn.node != nil {
					ctn.node.resetSessionInfo()
				}

				// retry via user/pass
				if hashedPassword != nil {
					command = newLoginCommand(ctn.dataBuffer)
					err = command.login(policy, ctn, hashedPassword)
				}
			}
		}

		if err != nil {
			if ctn.node != nil {
				atomic.AddInt64(&ctn.node.stats.ConnectionsFailed, 1)
			}
			// Socket not authenticated. Do not put back into pool.
			ctn.Close()
			return err
		}

		si := command.sessionInfo()
		if ctn.node != nil && si.isValid() {
			ctn.node.sessionInfo.Store(si)
		}
	}

	return nil
}

// Login will send authentication information to the server.
// This function is provided for using the connection in conjunction with external libraries.
// The password will be hashed everytime, which is a slow operation.
func (ctn *Connection) Login(policy *ClientPolicy) Error {
	if !policy.RequiresAuthentication() {
		return nil
	}

	hashedPassword, err := hashPassword(policy.Password)
	if err != nil {
		return err
	}

	return ctn.login(policy, hashedPassword, nil)
}

// RequestInfo gets info values by name from the specified connection.
// Timeout should already be set on the connection.
func (ctn *Connection) RequestInfo(names ...string) (map[string]string, Error) {
	info, err := newInfo(ctn, names...)
	if err != nil {
		return nil, err
	}

	return info.parseMultiResponse()
}

// setIdleTimeout sets the idle timeout for the connection.
func (ctn *Connection) setIdleTimeout(timeout time.Duration) {
	ctn.idleTimeout = timeout
}

// isIdle returns true if the connection has reached the idle deadline.
func (ctn *Connection) isIdle() bool {
	return ctn.idleTimeout > 0 && time.Now().After(ctn.idleDeadline)
}

// refresh extends the idle deadline of the connection.
func (ctn *Connection) refresh() {
	ctn.idleDeadline = time.Now().Add(ctn.idleTimeout)
	if ctn.inflater != nil {
		ctn.inflater.Close()
	}
	ctn.compressed = false
	ctn.inflater = nil
}

// initInflater sets up the zlib inflater to read compressed data from the connection
func (ctn *Connection) initInflater(enabled bool, length int) Error {
	ctn.compressed = enabled
	ctn.inflater = nil
	if ctn.compressed {
		ctn.limitReader.N = int64(length)
		r, err := zlib.NewReader(ctn.limitReader)
		if err != nil {
			return newCommonError(err)
		}
		ctn.inflater = r
	}
	return nil
}

// KeepConnection decides if a connection should be kept
// based on the error type.
func KeepConnection(err Error) bool {
	// Do not keep connection on client errors.
	if err.resultCode() < 0 {
		return false
	}

	return !err.Matches(types.QUERY_TERMINATED,
		types.SCAN_ABORT,
		types.QUERY_ABORTED,
		types.TIMEOUT)
}
