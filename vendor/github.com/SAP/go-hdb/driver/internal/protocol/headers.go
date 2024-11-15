package protocol

import (
	"fmt"
	"math"

	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
)

// Message header (size: 32 bytes).
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

func (h *messageHeader) encode(enc *encoding.Encoder) error {
	enc.Int64(h.sessionID)
	enc.Int32(h.packetCount)
	enc.Uint32(h.varPartLength)
	enc.Uint32(h.varPartSize)
	enc.Int16(h.noOfSegm)
	enc.Zeroes(10) // size: 32 bytes
	return nil
}

func (h *messageHeader) decode(dec *encoding.Decoder) error {
	h.sessionID = dec.Int64()
	h.packetCount = dec.Int32()
	h.varPartLength = dec.Uint32()
	h.varPartSize = dec.Uint32()
	h.noOfSegm = dec.Int16()
	dec.Skip(10) // size: 32 bytes
	return dec.Error()
}

const (
	segmentHeaderSize = 24
)

type segmentKind int8

const (
	skInvalid segmentKind = 0
	skRequest segmentKind = 1
	skReply   segmentKind = 2
	skError   segmentKind = 5
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

var (
	coList     = []commandOptions{coNil, coSelfetchOff, coScrollableCursorOn, coNoResultsetCloseNeeded, coHoldCursorOverCommtit, coExecuteLocally}
	coListText = []string{"", "selfetchOff", "scrollableCursorOn", "noResltsetCloseNeeded", "holdCursorOverCommit", "executLocally"}
)

func (k commandOptions) String() string {
	var s []string

	for i, option := range coList {
		if (k & option) != 0 {
			s = append(s, coListText[i])
		}
	}
	return fmt.Sprintf("%v", s)
}

// segment header.
type segmentHeader struct {
	segmentLength  int32
	segmentOfs     int32
	noOfParts      int16
	segmentNo      int16
	segmentKind    segmentKind
	messageType    MessageType
	commit         bool
	commandOptions commandOptions
	functionCode   FunctionCode
}

func (h *segmentHeader) String() string {
	switch h.segmentKind {
	default: // error
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

// request.
func (h *segmentHeader) encode(enc *encoding.Encoder) error {
	enc.Int32(h.segmentLength)
	enc.Int32(h.segmentOfs)
	enc.Int16(h.noOfParts)
	enc.Int16(h.segmentNo)
	enc.Int8(int8(h.segmentKind))

	switch h.segmentKind {
	default: // error
		enc.Zeroes(11) // segmentHeaderLength

	case skRequest:
		enc.Int8(int8(h.messageType))
		enc.Bool(h.commit)
		enc.Int8(int8(h.commandOptions))
		enc.Zeroes(8) // segmentHeaderSize

	case skReply:
		enc.Zeroes(1) // reserved
		enc.Int16(int16(h.functionCode))
		enc.Zeroes(8) // segmentHeaderSize
	}
	return nil
}

// reply || error.
func (h *segmentHeader) decode(dec *encoding.Decoder) error {
	h.segmentLength = dec.Int32()
	h.segmentOfs = dec.Int32()
	h.noOfParts = dec.Int16()
	h.segmentNo = dec.Int16()
	h.segmentKind = segmentKind(dec.Int8())

	switch h.segmentKind {
	default: // error
		dec.Skip(11) // segmentHeaderLength

	case skRequest:
		h.messageType = MessageType(dec.Int8())
		h.commit = dec.Bool()
		h.commandOptions = commandOptions(dec.Int8())
		dec.Skip(8) // segmentHeaderLength

	case skReply:
		dec.Skip(1) // reserved
		h.functionCode = FunctionCode(dec.Int16())
		dec.Skip(8) // segmentHeaderLength
	}
	return dec.Error()
}

const (
	partHeaderSize = 16
	bigNumArgInd   = -1
)

// MaxNumArg is the maximum number of arguments allowed to send in a part.
const MaxNumArg = math.MaxInt32

// PartAttributes represents the part attributes.
type PartAttributes int8

const (
	paLastPacket      PartAttributes = 0x01
	paNextPacket      PartAttributes = 0x02
	paFirstPacket     PartAttributes = 0x04
	paRowNotFound     PartAttributes = 0x08
	paResultsetClosed PartAttributes = 0x10
)

var (
	paList     = [...]PartAttributes{paLastPacket, paNextPacket, paFirstPacket, paRowNotFound, paResultsetClosed}
	paListText = [...]string{"lastPacket", "nextPacket", "firstPacket", "rowNotFound", "resultsetClosed"}
)

func (k PartAttributes) String() string {
	var s []string

	for i, attr := range paList {
		if (k & attr) != 0 {
			s = append(s, paListText[i])
		}
	}
	return fmt.Sprintf("%v", s)
}

// ResultsetClosed returns true if the result set is closed, false otherwise.
func (k PartAttributes) ResultsetClosed() bool { return (k & paResultsetClosed) == paResultsetClosed }

// LastPacket returns true if the last packet is sent, false otherwise.
func (k PartAttributes) LastPacket() bool { return (k & paLastPacket) == paLastPacket }

// partHeader represents the part header.
type partHeader struct {
	partKind         PartKind
	partAttributes   PartAttributes
	argumentCount    int16
	bigArgumentCount int32
	bufferLength     int32
	bufferSize       int32
}

func (h *partHeader) String() string {
	return fmt.Sprintf("kind %s partAttributes %s argumentCount %d bigArgumentCount %d bufferLength %d bufferSize %d",
		h.partKind,
		h.partAttributes,
		h.argumentCount,
		h.bigArgumentCount,
		h.bufferLength,
		h.bufferSize,
	)
}

func (h *partHeader) setNumArg(numArg int) error {
	switch {
	default:
		return fmt.Errorf("maximum number of arguments %d exceeded", numArg)
	case numArg <= math.MaxInt16:
		h.argumentCount = int16(numArg) //nolint: gosec
		h.bigArgumentCount = 0
	case numArg <= math.MaxInt32:
		h.argumentCount = bigNumArgInd
		h.bigArgumentCount = int32(numArg)
	}
	return nil
}

func (h *partHeader) numArg() int {
	if h.argumentCount == bigNumArgInd {
		return int(h.bigArgumentCount)
	}
	return int(h.argumentCount)
}

func (h *partHeader) bufLen() int { return int(h.bufferLength) }

func (h *partHeader) encode(enc *encoding.Encoder) error {
	enc.Int8(int8(h.partKind))
	enc.Int8(int8(h.partAttributes))
	enc.Int16(h.argumentCount)
	enc.Int32(h.bigArgumentCount)
	enc.Int32(h.bufferLength)
	enc.Int32(h.bufferSize)
	// no filler
	return nil
}

func (h *partHeader) decode(dec *encoding.Decoder) error {
	h.partKind = PartKind(dec.Int8())
	h.partAttributes = PartAttributes(dec.Int8())
	h.argumentCount = dec.Int16()
	h.bigArgumentCount = dec.Int32()
	h.bufferLength = dec.Int32()
	h.bufferSize = dec.Int32()
	// no filler
	return dec.Error()
}
