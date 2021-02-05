// Copyright 2013-2020 Aerospike, Inc.
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

	. "github.com/aerospike/aerospike-client-go/logger"
	. "github.com/aerospike/aerospike-client-go/types"
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

func errToTimeoutErr(conn *Connection, err error) error {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return ErrTimeout
	}
	return err
}

// newConnection creates a connection on the network and returns the pointer
// A minimum timeout of 2 seconds will always be applied.
// If the connection is not established in the specified timeout,
// an error will be returned
func newConnection(address string, timeout time.Duration) (*Connection, error) {
	newConn := &Connection{dataBuffer: bufPool.Get().([]byte)}
	runtime.SetFinalizer(newConn, connectionFinalizer)

	// don't wait indefinitely
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		Logger.Debug("Connection to address `" + address + "` failed to establish with error: " + err.Error())
		return nil, errToTimeoutErr(nil, err)
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
func NewConnection(policy *ClientPolicy, host *Host) (*Connection, error) {
	address := net.JoinHostPort(host.Name, strconv.Itoa(host.Port))
	conn, err := newConnection(address, policy.Timeout)
	if err != nil {
		return nil, err
	}

	if policy.TlsConfig == nil {
		return conn, nil
	}

	// Use version dependent clone function to clone the config
	tlsConfig := cloneTLSConfig(policy.TlsConfig)
	tlsConfig.ServerName = host.TLSName

	sconn := tls.Client(conn.conn, tlsConfig)
	if err := sconn.Handshake(); err != nil {
		sconn.Close()
		return nil, err
	}

	if host.TLSName != "" && !tlsConfig.InsecureSkipVerify {
		if err := sconn.VerifyHostname(host.TLSName); err != nil {
			sconn.Close()
			Logger.Error("Connection to address `" + address + "` failed to establish with error: " + err.Error())
			return nil, errToTimeoutErr(nil, err)
		}
	}

	conn.conn = sconn
	return conn, nil
}

// Write writes the slice to the connection buffer.
func (ctn *Connection) Write(buf []byte) (total int, err error) {
	// make sure all bytes are written
	// Don't worry about the loop, timeout has been set elsewhere
	length := len(buf)
	for total < length {
		var r int
		if err = ctn.updateDeadline(); err != nil {
			break
		}
		r, err = ctn.conn.Write(buf[total:])
		total += r
		if err != nil {
			break
		}
	}

	// If all bytes are written, ignore any potential error
	// The error will bubble up on the next network io if it matters.
	if total == len(buf) {
		return total, nil
	}

	if ctn.node != nil {
		ctn.node.incrErrorCount()
		atomic.AddInt64(&ctn.node.stats.ConnectionsFailed, 1)
	}

	ctn.Close()

	return total, errToTimeoutErr(ctn, err)
}

// Read reads from connection buffer to the provided slice.
func (ctn *Connection) Read(buf []byte, length int) (total int, err error) {
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
				ctn.inflater.Close()
				err = nil
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

	if ctn.node != nil {
		ctn.node.incrErrorCount()
		atomic.AddInt64(&ctn.node.stats.ConnectionsFailed, 1)
	}

	ctn.Close()

	return total, errToTimeoutErr(ctn, err)
}

// IsConnected returns true if the connection is not closed yet.
func (ctn *Connection) IsConnected() bool {
	return ctn.conn != nil
}

// updateDeadline sets connection timeout for both read and write operations.
// this function is called before each read and write operation. If deadline has passed,
// the function will return a TIMEOUT error.
func (ctn *Connection) updateDeadline() error {
	now := time.Now()
	var socketDeadline time.Time
	if ctn.deadline.IsZero() {
		if ctn.socketTimeout > 0 {
			socketDeadline = now.Add(ctn.socketTimeout)
		}
	} else {
		if now.After(ctn.deadline) {
			return NewAerospikeError(TIMEOUT)
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
		return err
	}

	return nil
}

// SetTimeout sets connection timeout for both read and write operations.
func (ctn *Connection) SetTimeout(deadline time.Time, socketTimeout time.Duration) error {
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
				Logger.Warn(err.Error())
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

// Authenticate will send authentication information to the server.
// Notice: This method does not support external authentication mechanisms like LDAP.
// This method is deprecated and will be removed in the future.
func (ctn *Connection) Authenticate(user string, password string) error {
	// need to authenticate
	if user != "" {
		hashedPass, err := hashPassword(password)
		if err != nil {
			return err
		}

		return ctn.authenticateFast(user, hashedPass)
	}
	return nil
}

// authenticateFast will send authentication information to the server.
func (ctn *Connection) authenticateFast(user string, hashedPass []byte) error {
	// need to authenticate
	if len(user) > 0 {
		command := newLoginCommand(ctn.dataBuffer)
		if err := command.authenticateInternal(ctn, user, hashedPass); err != nil {
			if ctn.node != nil {
				atomic.AddInt64(&ctn.node.stats.ConnectionsFailed, 1)
			}
			// Socket not authenticated. Do not put back into pool.
			ctn.Close()
			return err
		}
	}
	return nil
}

// Login will send authentication information to the server.
func (ctn *Connection) login(policy *ClientPolicy, hashedPassword []byte, sessionToken []byte) error {
	// need to authenticate
	if policy.RequiresAuthentication() {
		switch policy.AuthMode {
		case AuthModeExternal:
			var err error
			command := newLoginCommand(ctn.dataBuffer)
			if sessionToken == nil {
				err = command.login(policy, ctn, hashedPassword)
			} else {
				err = command.authenticateViaToken(policy, ctn, sessionToken)
			}

			if err != nil {
				if ctn.node != nil {
					atomic.AddInt64(&ctn.node.stats.ConnectionsFailed, 1)
				}
				// Socket not authenticated. Do not put back into pool.
				ctn.Close()
				return err
			}

			if ctn.node != nil && command.SessionToken != nil {
				ctn.node._sessionToken.Store(command.SessionToken)
				ctn.node._sessionExpiration.Store(command.SessionExpiration)
			}

			return nil

		case AuthModeInternal:
			return ctn.authenticateFast(policy.User, hashedPassword)
		}
	}

	return nil
}

// Login will send authentication information to the server.
// This function is provided for using the connection in conjunction with external libraries.
// The password will be hashed everytime, which is a slow operation.
func (ctn *Connection) Login(policy *ClientPolicy) error {
	if !policy.RequiresAuthentication() {
		return nil
	}

	hashedPassword, err := hashPassword(policy.Password)
	if err != nil {
		return err
	}

	return ctn.login(policy, hashedPassword, nil)
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
func (ctn *Connection) initInflater(enabled bool, length int) error {
	ctn.compressed = enabled
	ctn.inflater = nil
	if ctn.compressed {
		ctn.limitReader.N = int64(length)
		r, err := zlib.NewReader(ctn.limitReader)
		if err != nil {
			return err
		}
		ctn.inflater = r
	}
	return nil
}
