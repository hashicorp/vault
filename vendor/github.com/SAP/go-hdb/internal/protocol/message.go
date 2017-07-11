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
	if err := wr.WriteInt64(h.sessionID); err != nil {
		return err
	}
	if err := wr.WriteInt32(h.packetCount); err != nil {
		return err
	}
	if err := wr.WriteUint32(h.varPartLength); err != nil {
		return err
	}
	if err := wr.WriteUint32(h.varPartSize); err != nil {
		return err
	}
	if err := wr.WriteInt16(h.noOfSegm); err != nil {
		return err
	}

	if err := wr.WriteZeroes(10); err != nil { //messageHeaderSize
		return err
	}

	if trace {
		outLogger.Printf("write message header: %s", h)
	}

	return nil
}

func (h *messageHeader) read(rd *bufio.Reader) error {
	var err error

	if h.sessionID, err = rd.ReadInt64(); err != nil {
		return err
	}
	if h.packetCount, err = rd.ReadInt32(); err != nil {
		return err
	}
	if h.varPartLength, err = rd.ReadUint32(); err != nil {
		return err
	}
	if h.varPartSize, err = rd.ReadUint32(); err != nil {
		return err
	}
	if h.noOfSegm, err = rd.ReadInt16(); err != nil {
		return err
	}

	if err := rd.Skip(10); err != nil { //messageHeaderSize
		return err
	}

	if trace {
		outLogger.Printf("read message header: %s", h)
	}

	return nil
}
