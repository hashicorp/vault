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
	partHeaderSize = 16
)

type requestPart interface {
	kind() partKind
	size() (int, error)
	numArg() int
	write(*bufio.Writer) error
}

type replyPart interface {
	//kind() partKind
	setNumArg(int)
	read(*bufio.Reader) error
}

// PartAttributes is an interface defining methods for reading query resultset parts.
type PartAttributes interface {
	ResultsetClosed() bool
	LastPacket() bool
	NoRows() bool
}

type partAttributes int8

const (
	paLastPacket      partAttributes = 0x01
	paNextPacket      partAttributes = 0x02
	paFirstPacket     partAttributes = 0x04
	paRowNotFound     partAttributes = 0x08
	paResultsetClosed partAttributes = 0x10
)

var partAttributesText = map[partAttributes]string{
	paLastPacket:      "lastPacket",
	paNextPacket:      "nextPacket",
	paFirstPacket:     "firstPacket",
	paRowNotFound:     "rowNotFound",
	paResultsetClosed: "resultsetClosed",
}

func (k partAttributes) String() string {
	t := make([]string, 0, len(partAttributesText))

	for attr, text := range partAttributesText {
		if (k & attr) != 0 {
			t = append(t, text)
		}
	}
	return fmt.Sprintf("%v", t)
}

func (k partAttributes) ResultsetClosed() bool {
	return (k & paResultsetClosed) == paResultsetClosed
}

func (k partAttributes) LastPacket() bool {
	return (k & paLastPacket) == paLastPacket
}

func (k partAttributes) NoRows() bool {
	attrs := paLastPacket | paRowNotFound
	return (k & attrs) == attrs
}

// part header
type partHeader struct {
	partKind         partKind
	partAttributes   partAttributes
	argumentCount    int16
	bigArgumentCount int32
	bufferLength     int32
	bufferSize       int32
}

func (h *partHeader) String() string {
	return fmt.Sprintf("part kind %s partAttributes %s argumentCount %d bigArgumentCount %d bufferLength %d bufferSize %d",
		h.partKind,
		h.partAttributes,
		h.argumentCount,
		h.bigArgumentCount,
		h.bufferLength,
		h.bufferSize,
	)
}

func (h *partHeader) write(wr *bufio.Writer) error {
	wr.WriteInt8(int8(h.partKind))
	wr.WriteInt8(int8(h.partAttributes))
	wr.WriteInt16(h.argumentCount)
	wr.WriteInt32(h.bigArgumentCount)
	wr.WriteInt32(h.bufferLength)
	wr.WriteInt32(h.bufferSize)

	//no filler

	if trace {
		outLogger.Printf("write part header: %s", h)
	}

	return nil
}

func (h *partHeader) read(rd *bufio.Reader) error {
	h.partKind = partKind(rd.ReadInt8())
	h.partAttributes = partAttributes(rd.ReadInt8())
	h.argumentCount = rd.ReadInt16()
	h.bigArgumentCount = rd.ReadInt32()
	h.bufferLength = rd.ReadInt32()
	h.bufferSize = rd.ReadInt32()

	// no filler

	if trace {
		outLogger.Printf("read part header: %s", h)
	}

	return rd.GetError()
}
