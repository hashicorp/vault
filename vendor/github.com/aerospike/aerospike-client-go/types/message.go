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

package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type messageType uint8

const (
	MSG_HEADER_SIZE = 8 //sizeof(MessageHeader)

	MSG_INFO    messageType = 1
	MSG_MESSAGE             = 3
)

type MessageHeader struct {
	Version uint8
	Type    uint8
	DataLen [6]byte
}

// Length returns the length of the message
func (msg *MessageHeader) Length() int64 {
	return msgLenFromBytes(msg.DataLen)
}

type Message struct {
	MessageHeader

	Data []byte
}

// NewMessage generates a new Message instance.
func NewMessage(mtype messageType, data []byte) *Message {
	return &Message{
		MessageHeader: MessageHeader{
			Version: uint8(2),
			Type:    uint8(mtype),
			DataLen: msgLenToBytes(int64(len(data))),
		},
		Data: data,
	}
}

const maxAllowedBufferSize = 1024 * 1024

// Resize changes the internal buffer size for the message.
func (msg *Message) Resize(newSize int64) error {
	if newSize > maxAllowedBufferSize || newSize < 0 {
		return fmt.Errorf("Requested new buffer size is invalid. Requested: %d, allowed: 0..%d", newSize, maxAllowedBufferSize)
	}
	if int64(len(msg.Data)) == newSize {
		return nil
	}
	msg.Data = make([]byte, newSize)
	return nil
}

// Serialize returns a byte slice containing the message.
func (msg *Message) Serialize() []byte {
	msg.DataLen = msgLenToBytes(int64(len(msg.Data)))
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, msg.MessageHeader)
	binary.Write(buf, binary.BigEndian, msg.Data[:])

	return buf.Bytes()
}

func msgLenFromBytes(buf [6]byte) int64 {
	nbytes := append([]byte{0, 0}, buf[:]...)
	DataLen := binary.BigEndian.Uint64(nbytes)
	return int64(DataLen)
}

// converts a
func msgLenToBytes(DataLen int64) [6]byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(DataLen))
	res := [6]byte{}
	copy(res[:], b[2:])
	return res
}
