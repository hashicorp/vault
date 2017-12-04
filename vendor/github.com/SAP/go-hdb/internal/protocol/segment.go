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
	if err := wr.WriteInt32(h.segmentLength); err != nil {
		return err
	}
	if err := wr.WriteInt32(h.segmentOfs); err != nil {
		return err
	}
	if err := wr.WriteInt16(h.noOfParts); err != nil {
		return err
	}
	if err := wr.WriteInt16(h.segmentNo); err != nil {
		return err
	}
	if err := wr.WriteInt8(int8(h.segmentKind)); err != nil {
		return err
	}

	switch h.segmentKind {

	default: //error
		if err := wr.WriteZeroes(11); err != nil { //segmentHeaderLength
			return err
		}

	case skRequest:

		if err := wr.WriteInt8(int8(h.messageType)); err != nil {
			return err
		}
		if err := wr.WriteBool(h.commit); err != nil {
			return err
		}
		if err := wr.WriteInt8(int8(h.commandOptions)); err != nil {
			return err
		}

		if err := wr.WriteZeroes(8); err != nil { //segmentHeaderSize
			return err
		}

	case skReply:

		if err := wr.WriteZeroes(1); err != nil { //reerved
			return err
		}
		if err := wr.WriteInt16(int16(h.functionCode)); err != nil {
			return err
		}

		if err := wr.WriteZeroes(8); err != nil { //segmentHeaderSize
			return err
		}

	}

	if trace {
		outLogger.Printf("write segment header: %s", h)
	}

	return nil
}

//  reply || error
func (h *segmentHeader) read(rd *bufio.Reader) error {
	var err error

	if h.segmentLength, err = rd.ReadInt32(); err != nil {
		return err
	}
	if h.segmentOfs, err = rd.ReadInt32(); err != nil {
		return err
	}
	if h.noOfParts, err = rd.ReadInt16(); err != nil {
		return err
	}
	if h.segmentNo, err = rd.ReadInt16(); err != nil {
		return err
	}
	if sk, err := rd.ReadInt8(); err == nil {
		h.segmentKind = segmentKind(sk)
	} else {
		return err
	}

	switch h.segmentKind {

	default: //error
		if err := rd.Skip(11); err != nil { //segmentHeaderLength
			return err
		}

	case skRequest:
		if mt, err := rd.ReadInt8(); err == nil {
			h.messageType = messageType(mt)
		} else {
			return err
		}
		if h.commit, err = rd.ReadBool(); err != nil {
			return err
		}
		if co, err := rd.ReadInt8(); err == nil {
			h.commandOptions = commandOptions(co)
		} else {
			return err
		}
		if err := rd.Skip(8); err != nil { //segmentHeaderLength
			return err
		}

	case skReply:
		if err := rd.Skip(1); err != nil { //reserved
			return err
		}
		if fc, err := rd.ReadInt16(); err == nil {
			h.functionCode = functionCode(fc)
		} else {
			return err
		}
		if err := rd.Skip(8); err != nil { //segmentHeaderLength
			return err
		}
	}

	if trace {
		outLogger.Printf("read segment header: %s", h)
	}

	return nil
}
