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
	segmentHeaderSize = 24
)

type commandOptions int8

const (
	coNil                    commandOptions = 0x00
	coSelfetchOff            commandOptions = 0x01
	coScrollableCursorOn     commandOptions = 0x02
	coNoResultsetCloseNeeded commandOptions = 0x04
	coHoldCursorOverCommtit  commandOptions = 0x08
	coExecuteLocally         commandOptions = 0x10
)

var commandOptionsText = map[commandOptions]string{
	coSelfetchOff:            "selfetchOff",
	coScrollableCursorOn:     "scrollabeCursorOn",
	coNoResultsetCloseNeeded: "noResltsetCloseNeeded",
	coHoldCursorOverCommtit:  "holdCursorOverCommit",
	coExecuteLocally:         "executLocally",
}

func (k commandOptions) String() string {
	t := make([]string, 0, len(commandOptionsText))

	for option, text := range commandOptionsText {
		if (k & option) != 0 {
			t = append(t, text)
		}
	}
	return fmt.Sprintf("%v", t)
}

//segment header
type segmentHeader struct {
	segmentLength  int32
	segmentOfs     int32
	noOfParts      int16
	segmentNo      int16
	segmentKind    segmentKind
	messageType    messageType
	commit         bool
	commandOptions commandOptions
	functionCode   functionCode
}

func (h *segmentHeader) String() string {
	switch h.segmentKind {

	default: //error
		return fmt.Sprintf(
			"segment length %d segment ofs %d noOfParts %d, segmentNo %d segmentKind %s",
			h.segmentLength,
			h.segmentOfs,
			h.noOfParts,
			h.segmentNo,
			h.segmentKind,
		)
	case skRequest:
		return fmt.Sprintf(
			"segment length %d segment ofs %d noOfParts %d, segmentNo %d segmentKind %s messageType %s commit %t commandOptions %s",
			h.segmentLength,
			h.segmentOfs,
			h.noOfParts,
			h.segmentNo,
			h.segmentKind,
			h.messageType,
			h.commit,
			h.commandOptions,
		)
	case skReply:
		return fmt.Sprintf(
			"segment length %d segment ofs %d noOfParts %d, segmentNo %d segmentKind %s functionCode %s",
			h.segmentLength,
			h.segmentOfs,
			h.noOfParts,
			h.segmentNo,
			h.segmentKind,
			h.functionCode,
		)
	}
}

//  request
func (h *segmentHeader) write(wr *bufio.Writer) error {
	wr.WriteInt32(h.segmentLength)
	wr.WriteInt32(h.segmentOfs)
	wr.WriteInt16(h.noOfParts)
	wr.WriteInt16(h.segmentNo)
	wr.WriteInt8(int8(h.segmentKind))

	switch h.segmentKind {

	default: //error
		wr.WriteZeroes(11) //segmentHeaderLength

	case skRequest:
		wr.WriteInt8(int8(h.messageType))
		wr.WriteBool(h.commit)
		wr.WriteInt8(int8(h.commandOptions))
		wr.WriteZeroes(8) //segmentHeaderSize

	case skReply:

		wr.WriteZeroes(1) //reserved
		wr.WriteInt16(int16(h.functionCode))
		wr.WriteZeroes(8) //segmentHeaderSize

	}

	if trace {
		outLogger.Printf("write segment header: %s", h)
	}

	return nil
}

//  reply || error
func (h *segmentHeader) read(rd *bufio.Reader) error {
	h.segmentLength = rd.ReadInt32()
	h.segmentOfs = rd.ReadInt32()
	h.noOfParts = rd.ReadInt16()
	h.segmentNo = rd.ReadInt16()
	h.segmentKind = segmentKind(rd.ReadInt8())

	switch h.segmentKind {

	default: //error
		rd.Skip(11) //segmentHeaderLength

	case skRequest:
		h.messageType = messageType(rd.ReadInt8())
		h.commit = rd.ReadBool()
		h.commandOptions = commandOptions(rd.ReadInt8())
		rd.Skip(8) //segmentHeaderLength

	case skReply:
		rd.Skip(1) //reserved
		h.functionCode = functionCode(rd.ReadInt16())
		rd.Skip(8) //segmentHeaderLength

	}

	if trace {
		outLogger.Printf("read segment header: %s", h)
	}

	return rd.GetError()
}
