// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ipc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/bitutil"
	"github.com/apache/arrow/go/v15/arrow/endian"
	"github.com/apache/arrow/go/v15/arrow/internal"
	"github.com/apache/arrow/go/v15/arrow/internal/dictutils"
	"github.com/apache/arrow/go/v15/arrow/internal/flatbuf"
	"github.com/apache/arrow/go/v15/arrow/memory"
)

// FileReader is an Arrow file reader.
type FileReader struct {
	r ReadAtSeeker

	footer struct {
		offset int64
		buffer *memory.Buffer
		data   *flatbuf.Footer
	}

	// fields dictTypeMap
	memo dictutils.Memo

	schema *arrow.Schema
	record arrow.Record

	irec int   // current record index. used for the arrio.Reader interface
	err  error // last error

	mem            memory.Allocator
	swapEndianness bool
}

// NewFileReader opens an Arrow file using the provided reader r.
func NewFileReader(r ReadAtSeeker, opts ...Option) (*FileReader, error) {
	var (
		cfg = newConfig(opts...)
		err error

		f = FileReader{
			r:    r,
			memo: dictutils.NewMemo(),
			mem:  cfg.alloc,
		}
	)

	if cfg.footer.offset <= 0 {
		cfg.footer.offset, err = f.r.Seek(0, io.SeekEnd)
		if err != nil {
			return nil, fmt.Errorf("arrow/ipc: could retrieve footer offset: %w", err)
		}
	}
	f.footer.offset = cfg.footer.offset

	err = f.readFooter()
	if err != nil {
		return nil, fmt.Errorf("arrow/ipc: could not decode footer: %w", err)
	}

	err = f.readSchema(cfg.ensureNativeEndian)
	if err != nil {
		return nil, fmt.Errorf("arrow/ipc: could not decode schema: %w", err)
	}

	if cfg.schema != nil && !cfg.schema.Equal(f.schema) {
		return nil, fmt.Errorf("arrow/ipc: inconsistent schema for reading (got: %v, want: %v)", f.schema, cfg.schema)
	}

	return &f, err
}

func (f *FileReader) readFooter() error {
	var err error

	if f.footer.offset <= int64(len(Magic)*2+4) {
		return fmt.Errorf("arrow/ipc: file too small (size=%d)", f.footer.offset)
	}

	eof := int64(len(Magic) + 4)
	buf := make([]byte, eof)
	n, err := f.r.ReadAt(buf, f.footer.offset-eof)
	if err != nil {
		return fmt.Errorf("arrow/ipc: could not read footer: %w", err)
	}
	if n != len(buf) {
		return fmt.Errorf("arrow/ipc: could not read %d bytes from end of file", len(buf))
	}

	if !bytes.Equal(buf[4:], Magic) {
		return errNotArrowFile
	}

	size := int64(binary.LittleEndian.Uint32(buf[:4]))
	if size <= 0 || size+int64(len(Magic)*2+4) > f.footer.offset {
		return errInconsistentFileMetadata
	}

	buf = make([]byte, size)
	n, err = f.r.ReadAt(buf, f.footer.offset-size-eof)
	if err != nil {
		return fmt.Errorf("arrow/ipc: could not read footer data: %w", err)
	}
	if n != len(buf) {
		return fmt.Errorf("arrow/ipc: could not read %d bytes from footer data", len(buf))
	}

	f.footer.buffer = memory.NewBufferBytes(buf)
	f.footer.data = flatbuf.GetRootAsFooter(buf, 0)
	return err
}

func (f *FileReader) readSchema(ensureNativeEndian bool) error {
	var (
		err  error
		kind dictutils.Kind
	)

	schema := f.footer.data.Schema(nil)
	if schema == nil {
		return fmt.Errorf("arrow/ipc: could not load schema from flatbuffer data")
	}
	f.schema, err = schemaFromFB(schema, &f.memo)
	if err != nil {
		return fmt.Errorf("arrow/ipc: could not read schema: %w", err)
	}

	if ensureNativeEndian && !f.schema.IsNativeEndian() {
		f.swapEndianness = true
		f.schema = f.schema.WithEndianness(endian.NativeEndian)
	}

	for i := 0; i < f.NumDictionaries(); i++ {
		blk, err := f.dict(i)
		if err != nil {
			return fmt.Errorf("arrow/ipc: could not read dictionary[%d]: %w", i, err)
		}
		switch {
		case !bitutil.IsMultipleOf8(blk.Offset):
			return fmt.Errorf("arrow/ipc: invalid file offset=%d for dictionary %d", blk.Offset, i)
		case !bitutil.IsMultipleOf8(int64(blk.Meta)):
			return fmt.Errorf("arrow/ipc: invalid file metadata=%d position for dictionary %d", blk.Meta, i)
		case !bitutil.IsMultipleOf8(blk.Body):
			return fmt.Errorf("arrow/ipc: invalid file body=%d position for dictionary %d", blk.Body, i)
		}

		msg, err := blk.NewMessage()
		if err != nil {
			return err
		}

		kind, err = readDictionary(&f.memo, msg.meta, bytes.NewReader(msg.body.Bytes()), f.swapEndianness, f.mem)
		if err != nil {
			return err
		}
		if kind == dictutils.KindReplacement {
			return errors.New("arrow/ipc: unsupported dictionary replacement in IPC file")
		}
	}

	return err
}

func (f *FileReader) block(i int) (fileBlock, error) {
	var blk flatbuf.Block
	if !f.footer.data.RecordBatches(&blk, i) {
		return fileBlock{}, fmt.Errorf("arrow/ipc: could not extract file block %d", i)
	}

	return fileBlock{
		Offset: blk.Offset(),
		Meta:   blk.MetaDataLength(),
		Body:   blk.BodyLength(),
		r:      f.r,
		mem:    f.mem,
	}, nil
}

func (f *FileReader) dict(i int) (fileBlock, error) {
	var blk flatbuf.Block
	if !f.footer.data.Dictionaries(&blk, i) {
		return fileBlock{}, fmt.Errorf("arrow/ipc: could not extract dictionary block %d", i)
	}

	return fileBlock{
		Offset: blk.Offset(),
		Meta:   blk.MetaDataLength(),
		Body:   blk.BodyLength(),
		r:      f.r,
		mem:    f.mem,
	}, nil
}

func (f *FileReader) Schema() *arrow.Schema {
	return f.schema
}

func (f *FileReader) NumDictionaries() int {
	if f.footer.data == nil {
		return 0
	}
	return f.footer.data.DictionariesLength()
}

func (f *FileReader) NumRecords() int {
	return f.footer.data.RecordBatchesLength()
}

func (f *FileReader) Version() MetadataVersion {
	return MetadataVersion(f.footer.data.Version())
}

// Close cleans up resources used by the File.
// Close does not close the underlying reader.
func (f *FileReader) Close() error {
	if f.footer.data != nil {
		f.footer.data = nil
	}

	if f.footer.buffer != nil {
		f.footer.buffer.Release()
		f.footer.buffer = nil
	}

	if f.record != nil {
		f.record.Release()
		f.record = nil
	}
	return nil
}

// Record returns the i-th record from the file.
// The returned value is valid until the next call to Record.
// Users need to call Retain on that Record to keep it valid for longer.
func (f *FileReader) Record(i int) (arrow.Record, error) {
	record, err := f.RecordAt(i)
	if err != nil {
		return nil, err
	}

	if f.record != nil {
		f.record.Release()
	}

	f.record = record
	return record, nil
}

// Record returns the i-th record from the file. Ownership is transferred to the
// caller and must call Release() to free the memory. This method is safe to
// call concurrently.
func (f *FileReader) RecordAt(i int) (arrow.Record, error) {
	if i < 0 || i > f.NumRecords() {
		panic("arrow/ipc: record index out of bounds")
	}

	blk, err := f.block(i)
	if err != nil {
		return nil, err
	}
	switch {
	case !bitutil.IsMultipleOf8(blk.Offset):
		return nil, fmt.Errorf("arrow/ipc: invalid file offset=%d for record %d", blk.Offset, i)
	case !bitutil.IsMultipleOf8(int64(blk.Meta)):
		return nil, fmt.Errorf("arrow/ipc: invalid file metadata=%d position for record %d", blk.Meta, i)
	case !bitutil.IsMultipleOf8(blk.Body):
		return nil, fmt.Errorf("arrow/ipc: invalid file body=%d position for record %d", blk.Body, i)
	}

	msg, err := blk.NewMessage()
	if err != nil {
		return nil, err
	}
	defer msg.Release()

	if msg.Type() != MessageRecordBatch {
		return nil, fmt.Errorf("arrow/ipc: message %d is not a Record", i)
	}

	return newRecord(f.schema, &f.memo, msg.meta, bytes.NewReader(msg.body.Bytes()), f.swapEndianness, f.mem), nil
}

// Read reads the current record from the underlying stream and an error, if any.
// When the Reader reaches the end of the underlying stream, it returns (nil, io.EOF).
//
// The returned record value is valid until the next call to Read.
// Users need to call Retain on that Record to keep it valid for longer.
func (f *FileReader) Read() (rec arrow.Record, err error) {
	if f.irec == f.NumRecords() {
		return nil, io.EOF
	}
	rec, f.err = f.Record(f.irec)
	f.irec++
	return rec, f.err
}

// ReadAt reads the i-th record from the underlying stream and an error, if any.
func (f *FileReader) ReadAt(i int64) (arrow.Record, error) {
	return f.Record(int(i))
}

func newRecord(schema *arrow.Schema, memo *dictutils.Memo, meta *memory.Buffer, body ReadAtSeeker, swapEndianness bool, mem memory.Allocator) arrow.Record {
	var (
		msg   = flatbuf.GetRootAsMessage(meta.Bytes(), 0)
		md    flatbuf.RecordBatch
		codec decompressor
	)
	initFB(&md, msg.Header)
	rows := md.Length()

	bodyCompress := md.Compression(nil)
	if bodyCompress != nil {
		codec = getDecompressor(bodyCompress.Codec())
		defer codec.Close()
	}

	ctx := &arrayLoaderContext{
		src: ipcSource{
			meta:  &md,
			r:     body,
			codec: codec,
			mem:   mem,
		},
		memo:    memo,
		max:     kMaxNestingDepth,
		version: MetadataVersion(msg.Version()),
	}

	pos := dictutils.NewFieldPos()
	cols := make([]arrow.Array, schema.NumFields())
	for i := 0; i < schema.NumFields(); i++ {
		data := ctx.loadArray(schema.Field(i).Type)
		defer data.Release()

		if err := dictutils.ResolveFieldDict(memo, data, pos.Child(int32(i)), mem); err != nil {
			panic(err)
		}

		if swapEndianness {
			swapEndianArrayData(data.(*array.Data))
		}

		cols[i] = array.MakeFromData(data)
		defer cols[i].Release()
	}

	return array.NewRecord(schema, cols, rows)
}

type ipcSource struct {
	meta  *flatbuf.RecordBatch
	r     ReadAtSeeker
	codec decompressor
	mem   memory.Allocator
}

func (src *ipcSource) buffer(i int) *memory.Buffer {
	var buf flatbuf.Buffer
	if !src.meta.Buffers(&buf, i) {
		panic("arrow/ipc: buffer index out of bound")
	}

	if buf.Length() == 0 {
		return memory.NewBufferBytes(nil)
	}

	raw := memory.NewResizableBuffer(src.mem)
	if src.codec == nil {
		raw.Resize(int(buf.Length()))
		_, err := src.r.ReadAt(raw.Bytes(), buf.Offset())
		if err != nil {
			panic(err)
		}
	} else {
		sr := io.NewSectionReader(src.r, buf.Offset(), buf.Length())
		var uncompressedSize uint64

		err := binary.Read(sr, binary.LittleEndian, &uncompressedSize)
		if err != nil {
			panic(err)
		}

		var r io.Reader = sr
		// check for an uncompressed buffer
		if int64(uncompressedSize) != -1 {
			raw.Resize(int(uncompressedSize))
			src.codec.Reset(sr)
			r = src.codec
		} else {
			raw.Resize(int(buf.Length() - 8))
		}

		if _, err = io.ReadFull(r, raw.Bytes()); err != nil {
			panic(err)
		}
	}

	return raw
}

func (src *ipcSource) fieldMetadata(i int) *flatbuf.FieldNode {
	var node flatbuf.FieldNode
	if !src.meta.Nodes(&node, i) {
		panic("arrow/ipc: field metadata out of bound")
	}
	return &node
}

func (src *ipcSource) variadicCount(i int) int64 {
	return src.meta.VariadicBufferCounts(i)
}

type arrayLoaderContext struct {
	src       ipcSource
	ifield    int
	ibuffer   int
	ivariadic int
	max       int
	memo      *dictutils.Memo
	version   MetadataVersion
}

func (ctx *arrayLoaderContext) field() *flatbuf.FieldNode {
	field := ctx.src.fieldMetadata(ctx.ifield)
	ctx.ifield++
	return field
}

func (ctx *arrayLoaderContext) buffer() *memory.Buffer {
	buf := ctx.src.buffer(ctx.ibuffer)
	ctx.ibuffer++
	return buf
}

func (ctx *arrayLoaderContext) variadic() int64 {
	v := ctx.src.variadicCount(ctx.ivariadic)
	ctx.ivariadic++
	return v
}

func (ctx *arrayLoaderContext) loadArray(dt arrow.DataType) arrow.ArrayData {
	switch dt := dt.(type) {
	case *arrow.NullType:
		return ctx.loadNull()

	case *arrow.DictionaryType:
		indices := ctx.loadPrimitive(dt.IndexType)
		defer indices.Release()
		return array.NewData(dt, indices.Len(), indices.Buffers(), indices.Children(), indices.NullN(), indices.Offset())

	case *arrow.BooleanType,
		*arrow.Int8Type, *arrow.Int16Type, *arrow.Int32Type, *arrow.Int64Type,
		*arrow.Uint8Type, *arrow.Uint16Type, *arrow.Uint32Type, *arrow.Uint64Type,
		*arrow.Float16Type, *arrow.Float32Type, *arrow.Float64Type,
		*arrow.Decimal128Type, *arrow.Decimal256Type,
		*arrow.Time32Type, *arrow.Time64Type,
		*arrow.TimestampType,
		*arrow.Date32Type, *arrow.Date64Type,
		*arrow.MonthIntervalType, *arrow.DayTimeIntervalType, *arrow.MonthDayNanoIntervalType,
		*arrow.DurationType:
		return ctx.loadPrimitive(dt)

	case *arrow.BinaryType, *arrow.StringType, *arrow.LargeStringType, *arrow.LargeBinaryType:
		return ctx.loadBinary(dt)

	case arrow.BinaryViewDataType:
		return ctx.loadBinaryView(dt)

	case *arrow.FixedSizeBinaryType:
		return ctx.loadFixedSizeBinary(dt)

	case *arrow.ListType:
		return ctx.loadList(dt)

	case *arrow.LargeListType:
		return ctx.loadList(dt)

	case *arrow.ListViewType:
		return ctx.loadListView(dt)

	case *arrow.LargeListViewType:
		return ctx.loadListView(dt)

	case *arrow.FixedSizeListType:
		return ctx.loadFixedSizeList(dt)

	case *arrow.StructType:
		return ctx.loadStruct(dt)

	case *arrow.MapType:
		return ctx.loadMap(dt)

	case arrow.ExtensionType:
		storage := ctx.loadArray(dt.StorageType())
		defer storage.Release()
		return array.NewData(dt, storage.Len(), storage.Buffers(), storage.Children(), storage.NullN(), storage.Offset())

	case *arrow.RunEndEncodedType:
		field, buffers := ctx.loadCommon(dt.ID(), 1)
		defer memory.ReleaseBuffers(buffers)

		runEnds := ctx.loadChild(dt.RunEnds())
		defer runEnds.Release()
		values := ctx.loadChild(dt.Encoded())
		defer values.Release()

		return array.NewData(dt, int(field.Length()), buffers, []arrow.ArrayData{runEnds, values}, int(field.NullCount()), 0)

	case arrow.UnionType:
		return ctx.loadUnion(dt)

	default:
		panic(fmt.Errorf("arrow/ipc: array type %T not handled yet", dt))
	}
}

func (ctx *arrayLoaderContext) loadCommon(typ arrow.Type, nbufs int) (*flatbuf.FieldNode, []*memory.Buffer) {
	buffers := make([]*memory.Buffer, 0, nbufs)
	field := ctx.field()

	var buf *memory.Buffer

	if internal.HasValidityBitmap(typ, flatbuf.MetadataVersion(ctx.version)) {
		switch field.NullCount() {
		case 0:
			ctx.ibuffer++
		default:
			buf = ctx.buffer()
		}
	}
	buffers = append(buffers, buf)

	return field, buffers
}

func (ctx *arrayLoaderContext) loadChild(dt arrow.DataType) arrow.ArrayData {
	if ctx.max == 0 {
		panic("arrow/ipc: nested type limit reached")
	}
	ctx.max--
	sub := ctx.loadArray(dt)
	ctx.max++
	return sub
}

func (ctx *arrayLoaderContext) loadNull() arrow.ArrayData {
	field := ctx.field()
	return array.NewData(arrow.Null, int(field.Length()), nil, nil, int(field.NullCount()), 0)
}

func (ctx *arrayLoaderContext) loadPrimitive(dt arrow.DataType) arrow.ArrayData {
	field, buffers := ctx.loadCommon(dt.ID(), 2)

	switch field.Length() {
	case 0:
		buffers = append(buffers, nil)
		ctx.ibuffer++
	default:
		buffers = append(buffers, ctx.buffer())
	}

	defer memory.ReleaseBuffers(buffers)

	return array.NewData(dt, int(field.Length()), buffers, nil, int(field.NullCount()), 0)
}

func (ctx *arrayLoaderContext) loadBinary(dt arrow.DataType) arrow.ArrayData {
	field, buffers := ctx.loadCommon(dt.ID(), 3)
	buffers = append(buffers, ctx.buffer(), ctx.buffer())
	defer memory.ReleaseBuffers(buffers)

	return array.NewData(dt, int(field.Length()), buffers, nil, int(field.NullCount()), 0)
}

func (ctx *arrayLoaderContext) loadBinaryView(dt arrow.DataType) arrow.ArrayData {
	nVariadicBufs := ctx.variadic()
	field, buffers := ctx.loadCommon(dt.ID(), 2+int(nVariadicBufs))
	buffers = append(buffers, ctx.buffer())
	for i := 0; i < int(nVariadicBufs); i++ {
		buffers = append(buffers, ctx.buffer())
	}
	defer memory.ReleaseBuffers(buffers)

	return array.NewData(dt, int(field.Length()), buffers, nil, int(field.NullCount()), 0)
}

func (ctx *arrayLoaderContext) loadFixedSizeBinary(dt *arrow.FixedSizeBinaryType) arrow.ArrayData {
	field, buffers := ctx.loadCommon(dt.ID(), 2)
	buffers = append(buffers, ctx.buffer())
	defer memory.ReleaseBuffers(buffers)

	return array.NewData(dt, int(field.Length()), buffers, nil, int(field.NullCount()), 0)
}

func (ctx *arrayLoaderContext) loadMap(dt *arrow.MapType) arrow.ArrayData {
	field, buffers := ctx.loadCommon(dt.ID(), 2)
	buffers = append(buffers, ctx.buffer())
	defer memory.ReleaseBuffers(buffers)

	sub := ctx.loadChild(dt.Elem())
	defer sub.Release()

	return array.NewData(dt, int(field.Length()), buffers, []arrow.ArrayData{sub}, int(field.NullCount()), 0)
}

func (ctx *arrayLoaderContext) loadList(dt arrow.ListLikeType) arrow.ArrayData {
	field, buffers := ctx.loadCommon(dt.ID(), 2)
	buffers = append(buffers, ctx.buffer())
	defer memory.ReleaseBuffers(buffers)

	sub := ctx.loadChild(dt.Elem())
	defer sub.Release()

	return array.NewData(dt, int(field.Length()), buffers, []arrow.ArrayData{sub}, int(field.NullCount()), 0)
}

func (ctx *arrayLoaderContext) loadListView(dt arrow.VarLenListLikeType) arrow.ArrayData {
	field, buffers := ctx.loadCommon(dt.ID(), 3)
	buffers = append(buffers, ctx.buffer(), ctx.buffer())
	defer memory.ReleaseBuffers(buffers)

	sub := ctx.loadChild(dt.Elem())
	defer sub.Release()

	return array.NewData(dt, int(field.Length()), buffers, []arrow.ArrayData{sub}, int(field.NullCount()), 0)
}

func (ctx *arrayLoaderContext) loadFixedSizeList(dt *arrow.FixedSizeListType) arrow.ArrayData {
	field, buffers := ctx.loadCommon(dt.ID(), 1)
	defer memory.ReleaseBuffers(buffers)

	sub := ctx.loadChild(dt.Elem())
	defer sub.Release()

	return array.NewData(dt, int(field.Length()), buffers, []arrow.ArrayData{sub}, int(field.NullCount()), 0)
}

func (ctx *arrayLoaderContext) loadStruct(dt *arrow.StructType) arrow.ArrayData {
	field, buffers := ctx.loadCommon(dt.ID(), 1)
	defer memory.ReleaseBuffers(buffers)

	subs := make([]arrow.ArrayData, dt.NumFields())
	for i, f := range dt.Fields() {
		subs[i] = ctx.loadChild(f.Type)
	}
	defer func() {
		for i := range subs {
			subs[i].Release()
		}
	}()

	return array.NewData(dt, int(field.Length()), buffers, subs, int(field.NullCount()), 0)
}

func (ctx *arrayLoaderContext) loadUnion(dt arrow.UnionType) arrow.ArrayData {
	// Sparse unions have 2 buffers (a nil validity bitmap, and the type ids)
	nBuffers := 2
	// Dense unions have a third buffer, the offsets
	if dt.Mode() == arrow.DenseMode {
		nBuffers = 3
	}

	field, buffers := ctx.loadCommon(dt.ID(), nBuffers)
	if field.NullCount() != 0 && buffers[0] != nil {
		panic("arrow/ipc: cannot read pre-1.0.0 union array with top-level validity bitmap")
	}

	switch field.Length() {
	case 0:
		buffers = append(buffers, memory.NewBufferBytes([]byte{}))
		ctx.ibuffer++
		if dt.Mode() == arrow.DenseMode {
			buffers = append(buffers, nil)
			ctx.ibuffer++
		}
	default:
		buffers = append(buffers, ctx.buffer())
		if dt.Mode() == arrow.DenseMode {
			buffers = append(buffers, ctx.buffer())
		}
	}

	defer memory.ReleaseBuffers(buffers)
	subs := make([]arrow.ArrayData, dt.NumFields())
	for i, f := range dt.Fields() {
		subs[i] = ctx.loadChild(f.Type)
	}
	defer func() {
		for i := range subs {
			subs[i].Release()
		}
	}()
	return array.NewData(dt, int(field.Length()), buffers, subs, 0, 0)
}

func readDictionary(memo *dictutils.Memo, meta *memory.Buffer, body ReadAtSeeker, swapEndianness bool, mem memory.Allocator) (dictutils.Kind, error) {
	var (
		msg   = flatbuf.GetRootAsMessage(meta.Bytes(), 0)
		md    flatbuf.DictionaryBatch
		data  flatbuf.RecordBatch
		codec decompressor
	)
	initFB(&md, msg.Header)

	md.Data(&data)
	bodyCompress := data.Compression(nil)
	if bodyCompress != nil {
		codec = getDecompressor(bodyCompress.Codec())
		defer codec.Close()
	}

	id := md.Id()
	// look up the dictionary value type, which must have been added to the
	// memo already before calling this function
	valueType, ok := memo.Type(id)
	if !ok {
		return 0, fmt.Errorf("arrow/ipc: no dictionary type found with id: %d", id)
	}

	ctx := &arrayLoaderContext{
		src: ipcSource{
			meta:  &data,
			codec: codec,
			r:     body,
			mem:   mem,
		},
		memo: memo,
		max:  kMaxNestingDepth,
	}

	dict := ctx.loadArray(valueType)
	defer dict.Release()

	if swapEndianness {
		swapEndianArrayData(dict.(*array.Data))
	}

	if md.IsDelta() {
		memo.AddDelta(id, dict)
		return dictutils.KindDelta, nil
	}
	if memo.AddOrReplace(id, dict) {
		return dictutils.KindNew, nil
	}
	return dictutils.KindReplacement, nil
}
