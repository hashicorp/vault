// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package wiremessage is intended for internal use only. It is made available
// to facilitate use cases that require access to internal MongoDB driver
// functionality and state. The API of this package is not stable and there is
// no backward compatibility guarantee.
//
// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package wiremessage

import (
	"bytes"
	"encoding/binary"
	"strings"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// WireMessage represents a MongoDB wire message in binary form.
type WireMessage []byte

var globalRequestID int32

// NextRequestID returns the next request ID.
func NextRequestID() int32 { return atomic.AddInt32(&globalRequestID, 1) }

// OpCode represents a MongoDB wire protocol opcode.
type OpCode int32

// These constants are the valid opcodes for the version of the wireprotocol
// supported by this library. The skipped OpCodes are historical OpCodes that
// are no longer used.
const (
	OpReply  OpCode = 1
	_        OpCode = 1001
	OpUpdate OpCode = 2001
	OpInsert OpCode = 2002
	_        OpCode = 2003
	// Deprecated: Use OpMsg instead.
	OpQuery        OpCode = 2004
	OpGetMore      OpCode = 2005
	OpDelete       OpCode = 2006
	OpKillCursors  OpCode = 2007
	OpCommand      OpCode = 2010
	OpCommandReply OpCode = 2011
	OpCompressed   OpCode = 2012
	OpMsg          OpCode = 2013
)

// String implements the fmt.Stringer interface.
func (oc OpCode) String() string {
	switch oc {
	case OpReply:
		return "OP_REPLY"
	case OpUpdate:
		return "OP_UPDATE"
	case OpInsert:
		return "OP_INSERT"
	case OpQuery:
		return "OP_QUERY"
	case OpGetMore:
		return "OP_GET_MORE"
	case OpDelete:
		return "OP_DELETE"
	case OpKillCursors:
		return "OP_KILL_CURSORS"
	case OpCommand:
		return "OP_COMMAND"
	case OpCommandReply:
		return "OP_COMMANDREPLY"
	case OpCompressed:
		return "OP_COMPRESSED"
	case OpMsg:
		return "OP_MSG"
	default:
		return "<invalid opcode>"
	}
}

// QueryFlag represents the flags on an OP_QUERY message.
type QueryFlag int32

// These constants represent the individual flags on an OP_QUERY message.
const (
	_ QueryFlag = 1 << iota
	TailableCursor
	SecondaryOK
	OplogReplay
	NoCursorTimeout
	AwaitData
	Exhaust
	Partial
)

// String implements the fmt.Stringer interface.
func (qf QueryFlag) String() string {
	strs := make([]string, 0)
	if qf&TailableCursor == TailableCursor {
		strs = append(strs, "TailableCursor")
	}
	if qf&SecondaryOK == SecondaryOK {
		strs = append(strs, "SecondaryOK")
	}
	if qf&OplogReplay == OplogReplay {
		strs = append(strs, "OplogReplay")
	}
	if qf&NoCursorTimeout == NoCursorTimeout {
		strs = append(strs, "NoCursorTimeout")
	}
	if qf&AwaitData == AwaitData {
		strs = append(strs, "AwaitData")
	}
	if qf&Exhaust == Exhaust {
		strs = append(strs, "Exhaust")
	}
	if qf&Partial == Partial {
		strs = append(strs, "Partial")
	}
	str := "["
	str += strings.Join(strs, ", ")
	str += "]"
	return str
}

// MsgFlag represents the flags on an OP_MSG message.
type MsgFlag uint32

// These constants represent the individual flags on an OP_MSG message.
const (
	ChecksumPresent MsgFlag = 1 << iota
	MoreToCome

	ExhaustAllowed MsgFlag = 1 << 16
)

// ReplyFlag represents the flags of an OP_REPLY message.
type ReplyFlag int32

// These constants represent the individual flags of an OP_REPLY message.
const (
	CursorNotFound ReplyFlag = 1 << iota
	QueryFailure
	ShardConfigStale
	AwaitCapable
)

// String implements the fmt.Stringer interface.
func (rf ReplyFlag) String() string {
	strs := make([]string, 0)
	if rf&CursorNotFound == CursorNotFound {
		strs = append(strs, "CursorNotFound")
	}
	if rf&QueryFailure == QueryFailure {
		strs = append(strs, "QueryFailure")
	}
	if rf&ShardConfigStale == ShardConfigStale {
		strs = append(strs, "ShardConfigStale")
	}
	if rf&AwaitCapable == AwaitCapable {
		strs = append(strs, "AwaitCapable")
	}
	str := "["
	str += strings.Join(strs, ", ")
	str += "]"
	return str
}

// SectionType represents the type for 1 section in an OP_MSG
type SectionType uint8

// These constants represent the individual section types for a section in an OP_MSG
const (
	SingleDocument SectionType = iota
	DocumentSequence
)

// CompressorID is the ID for each type of Compressor.
type CompressorID uint8

// These constants represent the individual compressor IDs for an OP_COMPRESSED.
const (
	CompressorNoOp CompressorID = iota
	CompressorSnappy
	CompressorZLib
	CompressorZstd
)

// String implements the fmt.Stringer interface.
func (id CompressorID) String() string {
	switch id {
	case CompressorNoOp:
		return "CompressorNoOp"
	case CompressorSnappy:
		return "CompressorSnappy"
	case CompressorZLib:
		return "CompressorZLib"
	case CompressorZstd:
		return "CompressorZstd"
	default:
		return "CompressorInvalid"
	}
}

const (
	// DefaultZlibLevel is the default level for zlib compression
	DefaultZlibLevel = 6
	// DefaultZstdLevel is the default level for zstd compression.
	// Matches https://github.com/wiredtiger/wiredtiger/blob/f08bc4b18612ef95a39b12166abcccf207f91596/ext/compressors/zstd/zstd_compress.c#L299
	DefaultZstdLevel = 6
)

// AppendHeaderStart appends a header to the dst slice and returns an index where the wire message
// starts in dst and the updated slice.
func AppendHeaderStart(dst []byte, reqid, respto int32, opcode OpCode) (index int32, b []byte) {
	index, dst = bsoncore.ReserveLength(dst)
	dst = appendi32(dst, reqid)
	dst = appendi32(dst, respto)
	dst = appendi32(dst, int32(opcode))
	return index, dst
}

// AppendHeader appends a header to dst.
func AppendHeader(dst []byte, length, reqid, respto int32, opcode OpCode) []byte {
	dst = appendi32(dst, length)
	dst = appendi32(dst, reqid)
	dst = appendi32(dst, respto)
	dst = appendi32(dst, int32(opcode))
	return dst
}

// ReadHeader reads a wire message header from src.
func ReadHeader(src []byte) (length, requestID, responseTo int32, opcode OpCode, rem []byte, ok bool) {
	if len(src) < 16 {
		return 0, 0, 0, 0, src, false
	}

	length = readi32unsafe(src)
	requestID = readi32unsafe(src[4:])
	responseTo = readi32unsafe(src[8:])
	opcode = OpCode(readi32unsafe(src[12:]))
	return length, requestID, responseTo, opcode, src[16:], true
}

// AppendQueryFlags appends the flags for an OP_QUERY wire message.
func AppendQueryFlags(dst []byte, flags QueryFlag) []byte {
	return appendi32(dst, int32(flags))
}

// AppendMsgFlags appends the flags for an OP_MSG wire message.
func AppendMsgFlags(dst []byte, flags MsgFlag) []byte {
	return appendi32(dst, int32(flags))
}

// AppendReplyFlags appends the flags for an OP_REPLY wire message.
func AppendReplyFlags(dst []byte, flags ReplyFlag) []byte {
	return appendi32(dst, int32(flags))
}

// AppendMsgSectionType appends the section type to dst.
func AppendMsgSectionType(dst []byte, stype SectionType) []byte {
	return append(dst, byte(stype))
}

// AppendQueryFullCollectionName appends the full collection name to dst.
func AppendQueryFullCollectionName(dst []byte, ns string) []byte {
	return appendCString(dst, ns)
}

// AppendQueryNumberToSkip appends the number to skip to dst.
func AppendQueryNumberToSkip(dst []byte, skip int32) []byte {
	return appendi32(dst, skip)
}

// AppendQueryNumberToReturn appends the number to return to dst.
func AppendQueryNumberToReturn(dst []byte, nor int32) []byte {
	return appendi32(dst, nor)
}

// AppendReplyCursorID appends the cursor ID to dst.
func AppendReplyCursorID(dst []byte, id int64) []byte {
	return appendi64(dst, id)
}

// AppendReplyStartingFrom appends the starting from field to dst.
func AppendReplyStartingFrom(dst []byte, sf int32) []byte {
	return appendi32(dst, sf)
}

// AppendReplyNumberReturned appends the number returned to dst.
func AppendReplyNumberReturned(dst []byte, nr int32) []byte {
	return appendi32(dst, nr)
}

// AppendCompressedOriginalOpCode appends the original opcode to dst.
func AppendCompressedOriginalOpCode(dst []byte, opcode OpCode) []byte {
	return appendi32(dst, int32(opcode))
}

// AppendCompressedUncompressedSize appends the uncompressed size of a
// compressed wiremessage to dst.
func AppendCompressedUncompressedSize(dst []byte, size int32) []byte { return appendi32(dst, size) }

// AppendCompressedCompressorID appends the ID of the compressor to dst.
func AppendCompressedCompressorID(dst []byte, id CompressorID) []byte {
	return append(dst, byte(id))
}

// AppendCompressedCompressedMessage appends the compressed wiremessage to dst.
func AppendCompressedCompressedMessage(dst []byte, msg []byte) []byte { return append(dst, msg...) }

// AppendGetMoreZero appends the zero field to dst.
func AppendGetMoreZero(dst []byte) []byte {
	return appendi32(dst, 0)
}

// AppendGetMoreFullCollectionName appends the fullCollectionName field to dst.
func AppendGetMoreFullCollectionName(dst []byte, ns string) []byte {
	return appendCString(dst, ns)
}

// AppendGetMoreNumberToReturn appends the numberToReturn field to dst.
func AppendGetMoreNumberToReturn(dst []byte, numToReturn int32) []byte {
	return appendi32(dst, numToReturn)
}

// AppendGetMoreCursorID appends the cursorID field to dst.
func AppendGetMoreCursorID(dst []byte, cursorID int64) []byte {
	return appendi64(dst, cursorID)
}

// AppendKillCursorsZero appends the zero field to dst.
func AppendKillCursorsZero(dst []byte) []byte {
	return appendi32(dst, 0)
}

// AppendKillCursorsNumberIDs appends the numberOfCursorIDs field to dst.
func AppendKillCursorsNumberIDs(dst []byte, numIDs int32) []byte {
	return appendi32(dst, numIDs)
}

// AppendKillCursorsCursorIDs appends each the cursorIDs field to dst.
func AppendKillCursorsCursorIDs(dst []byte, cursors []int64) []byte {
	for _, cursor := range cursors {
		dst = appendi64(dst, cursor)
	}
	return dst
}

// ReadMsgFlags reads the OP_MSG flags from src.
func ReadMsgFlags(src []byte) (flags MsgFlag, rem []byte, ok bool) {
	i32, rem, ok := readi32(src)
	return MsgFlag(i32), rem, ok
}

// IsMsgMoreToCome returns if the provided wire message is an OP_MSG with the more to come flag set.
func IsMsgMoreToCome(wm []byte) bool {
	return len(wm) >= 20 &&
		OpCode(readi32unsafe(wm[12:16])) == OpMsg &&
		MsgFlag(readi32unsafe(wm[16:20]))&MoreToCome == MoreToCome
}

// ReadMsgSectionType reads the section type from src.
func ReadMsgSectionType(src []byte) (stype SectionType, rem []byte, ok bool) {
	if len(src) < 1 {
		return 0, src, false
	}
	return SectionType(src[0]), src[1:], true
}

// ReadMsgSectionSingleDocument reads a single document from src.
func ReadMsgSectionSingleDocument(src []byte) (doc bsoncore.Document, rem []byte, ok bool) {
	return bsoncore.ReadDocument(src)
}

// ReadMsgSectionDocumentSequence reads an identifier and document sequence from src and returns the document sequence
// data parsed into a slice of BSON documents.
func ReadMsgSectionDocumentSequence(src []byte) (identifier string, docs []bsoncore.Document, rem []byte, ok bool) {
	identifier, rem, ret, ok := ReadMsgSectionRawDocumentSequence(src)
	if !ok {
		return "", nil, src, false
	}

	docs = make([]bsoncore.Document, 0)
	var doc bsoncore.Document
	for {
		doc, rem, ok = bsoncore.ReadDocument(rem)
		if !ok {
			break
		}
		docs = append(docs, doc)
	}
	if len(rem) > 0 {
		return "", nil, src, false
	}

	return identifier, docs, ret, true
}

// ReadMsgSectionRawDocumentSequence reads an identifier and document sequence from src and returns the raw document
// sequence data.
func ReadMsgSectionRawDocumentSequence(src []byte) (identifier string, data []byte, rem []byte, ok bool) {
	length, rem, ok := readi32(src)
	if !ok || int(length) > len(src) || length-4 < 0 {
		return "", nil, src, false
	}

	// After these assignments, rem will be the data containing the identifier string + the document sequence bytes and
	// rest will be the rest of the wire message after this document sequence.
	rem, rest := rem[:length-4], rem[length-4:]

	identifier, rem, ok = readcstring(rem)
	if !ok {
		return "", nil, src, false
	}

	return identifier, rem, rest, true
}

// ReadMsgChecksum reads a checksum from src.
func ReadMsgChecksum(src []byte) (checksum uint32, rem []byte, ok bool) {
	i32, rem, ok := readi32(src)
	return uint32(i32), rem, ok
}

// ReadQueryFlags reads OP_QUERY flags from src.
//
// Deprecated: Construct wiremessages with OpMsg and use the ReadMsg* functions
// instead.
func ReadQueryFlags(src []byte) (flags QueryFlag, rem []byte, ok bool) {
	i32, rem, ok := readi32(src)
	return QueryFlag(i32), rem, ok
}

// ReadQueryFullCollectionName reads the full collection name from src.
//
// Deprecated: Construct wiremessages with OpMsg and use the ReadMsg* functions
// instead.
func ReadQueryFullCollectionName(src []byte) (collname string, rem []byte, ok bool) {
	return readcstring(src)
}

// ReadQueryNumberToSkip reads the number to skip from src.
//
// Deprecated: Construct wiremessages with OpMsg and use the ReadMsg* functions
// instead.
func ReadQueryNumberToSkip(src []byte) (nts int32, rem []byte, ok bool) {
	return readi32(src)
}

// ReadQueryNumberToReturn reads the number to return from src.
//
// Deprecated: Construct wiremessages with OpMsg and use the ReadMsg* functions
// instead.
func ReadQueryNumberToReturn(src []byte) (ntr int32, rem []byte, ok bool) {
	return readi32(src)
}

// ReadQueryQuery reads the query from src.
//
// Deprecated: Construct wiremessages with OpMsg and use the ReadMsg* functions
// instead.
func ReadQueryQuery(src []byte) (query bsoncore.Document, rem []byte, ok bool) {
	return bsoncore.ReadDocument(src)
}

// ReadQueryReturnFieldsSelector reads a return fields selector document from src.
//
// Deprecated: Construct wiremessages with OpMsg and use the ReadMsg* functions
// instead.
func ReadQueryReturnFieldsSelector(src []byte) (rfs bsoncore.Document, rem []byte, ok bool) {
	return bsoncore.ReadDocument(src)
}

// ReadReplyFlags reads OP_REPLY flags from src.
func ReadReplyFlags(src []byte) (flags ReplyFlag, rem []byte, ok bool) {
	i32, rem, ok := readi32(src)
	return ReplyFlag(i32), rem, ok
}

// ReadReplyCursorID reads a cursor ID from src.
func ReadReplyCursorID(src []byte) (cursorID int64, rem []byte, ok bool) {
	return readi64(src)
}

// ReadReplyStartingFrom reads the starting from src.
func ReadReplyStartingFrom(src []byte) (startingFrom int32, rem []byte, ok bool) {
	return readi32(src)
}

// ReadReplyNumberReturned reads the numbered returned from src.
func ReadReplyNumberReturned(src []byte) (numberReturned int32, rem []byte, ok bool) {
	return readi32(src)
}

// ReadReplyDocuments reads as many documents as possible from src
func ReadReplyDocuments(src []byte) (docs []bsoncore.Document, rem []byte, ok bool) {
	rem = src
	for {
		var doc bsoncore.Document
		doc, rem, ok = bsoncore.ReadDocument(rem)
		if !ok {
			break
		}

		docs = append(docs, doc)
	}

	return docs, rem, true
}

// ReadReplyDocument reads a reply document from src.
func ReadReplyDocument(src []byte) (doc bsoncore.Document, rem []byte, ok bool) {
	return bsoncore.ReadDocument(src)
}

// ReadCompressedOriginalOpCode reads the original opcode from src.
func ReadCompressedOriginalOpCode(src []byte) (opcode OpCode, rem []byte, ok bool) {
	i32, rem, ok := readi32(src)
	return OpCode(i32), rem, ok
}

// ReadCompressedUncompressedSize reads the uncompressed size of a
// compressed wiremessage to dst.
func ReadCompressedUncompressedSize(src []byte) (size int32, rem []byte, ok bool) {
	return readi32(src)
}

// ReadCompressedCompressorID reads the ID of the compressor to dst.
func ReadCompressedCompressorID(src []byte) (id CompressorID, rem []byte, ok bool) {
	if len(src) < 1 {
		return 0, src, false
	}
	return CompressorID(src[0]), src[1:], true
}

// ReadKillCursorsZero reads the zero field from src.
func ReadKillCursorsZero(src []byte) (zero int32, rem []byte, ok bool) {
	return readi32(src)
}

// ReadKillCursorsNumberIDs reads the numberOfCursorIDs field from src.
func ReadKillCursorsNumberIDs(src []byte) (numIDs int32, rem []byte, ok bool) {
	return readi32(src)
}

// ReadKillCursorsCursorIDs reads numIDs cursor IDs from src.
func ReadKillCursorsCursorIDs(src []byte, numIDs int32) (cursorIDs []int64, rem []byte, ok bool) {
	var i int32
	var id int64
	for i = 0; i < numIDs; i++ {
		id, src, ok = readi64(src)
		if !ok {
			return cursorIDs, src, false
		}

		cursorIDs = append(cursorIDs, id)
	}
	return cursorIDs, src, true
}

func appendi32(dst []byte, x int32) []byte {
	b := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(b, uint32(x))
	return append(dst, b...)
}

func appendi64(dst []byte, x int64) []byte {
	b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint64(b, uint64(x))
	return append(dst, b...)
}

func appendCString(b []byte, str string) []byte {
	b = append(b, str...)
	return append(b, 0x00)
}

func readi32(src []byte) (int32, []byte, bool) {
	if len(src) < 4 {
		return 0, src, false
	}
	return readi32unsafe(src), src[4:], true
}

func readi32unsafe(src []byte) int32 {
	return int32(binary.LittleEndian.Uint32(src))
}

func readi64(src []byte) (int64, []byte, bool) {
	if len(src) < 8 {
		return 0, src, false
	}
	return int64(binary.LittleEndian.Uint64(src)), src[8:], true
}

func readcstring(src []byte) (string, []byte, bool) {
	idx := bytes.IndexByte(src, 0x00)
	if idx < 0 {
		return "", src, false
	}
	return string(src[:idx]), src[idx+1:], true
}
