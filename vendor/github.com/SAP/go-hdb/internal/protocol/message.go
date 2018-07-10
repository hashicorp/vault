/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/bufio"
)

const (
	messageHeaderSize = 32
)

//message header
type messageHeader struct {
	sessionID     int64
	packetCount   int32
	varPartLength uint32
	varPartSize   uint32
	noOfSegm      int16
}

func (h *messageHeader) String() string {
	return fmt.Sprintf("session id %d packetCount %d varPartLength %d, varPartSize %d noOfSegm %d",
		h.sessionID,
		h.packetCount,
		h.varPartLength,
		h.varPartSize,
		h.noOfSegm)
}

func (h *messageHeader) write(wr *bufio.Writer) error {
	wr.WriteInt64(h.sessionID)
	wr.WriteInt32(h.packetCount)
	wr.WriteUint32(h.varPartLength)
	wr.WriteUint32(h.varPartSize)
	wr.WriteInt16(h.noOfSegm)
	wr.WriteZeroes(10) //messageHeaderSize

	if trace {
		outLogger.Printf("write message header: %s", h)
	}

	return nil
}

func (h *messageHeader) read(rd *bufio.Reader) error {
	h.sessionID = rd.ReadInt64()
	h.packetCount = rd.ReadInt32()
	h.varPartLength = rd.ReadUint32()
	h.varPartSize = rd.ReadUint32()
	h.noOfSegm = rd.ReadInt16()
	rd.Skip(10) //messageHeaderSize

	if trace {
		outLogger.Printf("read message header: %s", h)
	}

	return rd.GetError()
}
