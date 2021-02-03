// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
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
			"segmentLength %d segmentOfs %d noOfParts %d, segmentNo %d segmentKind %s",
			h.segmentLength,
			h.segmentOfs,
			h.noOfParts,
			h.segmentNo,
			h.segmentKind,
		)
	case skRequest:
		return fmt.Sprintf(
			"segmentLength %d segmentOfs %d noOfParts %d, segmentNo %d segmentKind %s messageType %s commit %t commandOptions %s",
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
			"segmentLength %d segmentOfs %d noOfParts %d, segmentNo %d segmentKind %s functionCode %s",
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
func (h *segmentHeader) encode(enc *encoding.Encoder) error {
	enc.Int32(h.segmentLength)
	enc.Int32(h.segmentOfs)
	enc.Int16(h.noOfParts)
	enc.Int16(h.segmentNo)
	enc.Int8(int8(h.segmentKind))

	switch h.segmentKind {

	default: //error
		enc.Zeroes(11) //segmentHeaderLength

	case skRequest:
		enc.Int8(int8(h.messageType))
		enc.Bool(h.commit)
		enc.Int8(int8(h.commandOptions))
		enc.Zeroes(8) //segmentHeaderSize

	case skReply:
		enc.Zeroes(1) //reserved
		enc.Int16(int16(h.functionCode))
		enc.Zeroes(8) //segmentHeaderSize
	}
	return nil
}

//  reply || error
func (h *segmentHeader) decode(dec *encoding.Decoder) error {
	h.segmentLength = dec.Int32()
	h.segmentOfs = dec.Int32()
	h.noOfParts = dec.Int16()
	h.segmentNo = dec.Int16()
	h.segmentKind = segmentKind(dec.Int8())

	switch h.segmentKind {

	default: //error
		dec.Skip(11) //segmentHeaderLength

	case skRequest:
		h.messageType = messageType(dec.Int8())
		h.commit = dec.Bool()
		h.commandOptions = commandOptions(dec.Int8())
		dec.Skip(8) //segmentHeaderLength

	case skReply:
		dec.Skip(1) //reserved
		h.functionCode = functionCode(dec.Int16())
		dec.Skip(8) //segmentHeaderLength
	}
	return dec.Error()
}
