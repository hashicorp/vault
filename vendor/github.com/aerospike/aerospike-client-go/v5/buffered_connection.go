// Copyright 2014-2019 Aerospike, Inbc.
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
	"fmt"

	"github.com/aerospike/aerospike-client-go/v5/logger"
	"github.com/aerospike/aerospike-client-go/v5/types"
)

type bufferedConn struct {
	conn       *Connection
	remaining  int
	head, tail int
}

func newBufferedConn(conn *Connection, total int) bufferedConn {
	return bufferedConn{
		conn:      conn,
		remaining: total,
	}
}

// emptyCap returns the empty capacity in the buffer.
func (bc *bufferedConn) emptyCap() int {
	return len(bc.buf()) - bc.len()
}

// len returns the number of unread bytes in the buffer.
func (bc *bufferedConn) len() int {
	return bc.tail - bc.head
}

// buf returns the connection's byte buffer.
func (bc *bufferedConn) buf() []byte {
	return bc.conn.dataBuffer
}

// shiftContentToHead will move the unread bytes to the head of the buffer.
// It will also resize the buffer if there is not enough empty capacity to read
// the minimum number of bytes that are requested.
// If the buffer is empty, head and tail will be reset to the beginning of the buffer.
func (bc *bufferedConn) shiftContentToHead(length int) {
	// shift data to the head of the byte slice
	if length > bc.emptyCap() {
		buf := make([]byte, bc.len()+length)
		copy(buf, bc.buf()[bc.head:bc.tail])
		bc.conn.dataBuffer = buf

		bc.tail -= bc.head
		bc.head = 0
	} else if bc.len() > 0 {
		copy(bc.buf(), bc.buf()[bc.head:bc.tail])

		bc.tail -= bc.head
		bc.head = 0
	} else {
		bc.tail = 0
		bc.head = 0
	}
}

// readConn will read the minimum minLength number of bytes from the connection.
// It will read more if it has extra empty capacity in the buffer.
func (bc *bufferedConn) readConn(minLength int) Error {
	// Corrupted data streams can result in a huge minLength.
	// Do a sanity check here.
	if minLength > MaxBufferSize || minLength <= 0 || minLength > bc.remaining {
		return newError(types.PARSE_ERROR, fmt.Sprintf("Invalid readBytes length: %d", minLength))
	}

	bc.shiftContentToHead(minLength)

	toRead := bc.remaining
	if ec := bc.emptyCap(); toRead > ec {
		toRead = ec
	}

	n, err := bc.conn.Read(bc.buf()[bc.tail:], toRead)
	bc.tail += n
	bc.remaining -= n

	if err != nil {
		logger.Logger.Debug("Requested to read %d bytes, but %d was read. (%v)", minLength, n, err)
		return err
	}

	return nil
}

func (bc *bufferedConn) read(length int) ([]byte, Error) {
	if cl := bc.len(); length > cl {
		if err := bc.readConn(length - cl); err != nil {
			return nil, err
		}
	}

	buf := bc.buf()[bc.head : bc.head+length]
	bc.head += length

	return buf, nil
}

func (bc *bufferedConn) drainConn() Error {
	if !bc.conn.IsConnected() {
		return nil
	}

	toRead := 0
	for bc.remaining > 0 {
		toRead = bc.remaining
		if toRead > len(bc.buf()) {
			toRead = len(bc.buf())
		}

		n, err := bc.conn.Read(bc.conn.dataBuffer, toRead)
		bc.remaining -= n
		if err != nil {
			return err
		}
	}

	return nil
}

func (bc *bufferedConn) reset(total int) Error {
	bc.remaining = total
	bc.head = 0
	bc.tail = 0
	return nil
}
