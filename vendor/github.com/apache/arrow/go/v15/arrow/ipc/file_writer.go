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
	"encoding/binary"
	"fmt"
	"io"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/bitutil"
	"github.com/apache/arrow/go/v15/arrow/internal/dictutils"
	"github.com/apache/arrow/go/v15/arrow/internal/flatbuf"
	"github.com/apache/arrow/go/v15/arrow/memory"
)

// PayloadWriter is an interface for injecting a different payloadwriter
// allowing more reusability with the Writer object with other scenarios,
// such as with Flight data
type PayloadWriter interface {
	Start() error
	WritePayload(Payload) error
	Close() error
}

type pwriter struct {
	w   io.WriteSeeker
	pos int64

	schema *arrow.Schema
	dicts  []fileBlock
	recs   []fileBlock
}

func (w *pwriter) Start() error {
	var err error

	err = w.updatePos()
	if err != nil {
		return fmt.Errorf("arrow/ipc: could not update position while in start: %w", err)
	}

	// only necessary to align to 8-byte boundary at the start of the file
	_, err = w.Write(Magic)
	if err != nil {
		return fmt.Errorf("arrow/ipc: could not write magic Arrow bytes: %w", err)
	}

	err = w.align(kArrowIPCAlignment)
	if err != nil {
		return fmt.Errorf("arrow/ipc: could not align start block: %w", err)
	}

	return err
}

func (w *pwriter) WritePayload(p Payload) error {
	blk := fileBlock{Offset: w.pos, Meta: 0, Body: p.size}
	n, err := writeIPCPayload(w, p)
	if err != nil {
		return err
	}

	blk.Meta = int32(n)

	err = w.updatePos()
	if err != nil {
		return fmt.Errorf("arrow/ipc: could not update position while in write-payload: %w", err)
	}

	switch flatbuf.MessageHeader(p.msg) {
	case flatbuf.MessageHeaderDictionaryBatch:
		w.dicts = append(w.dicts, blk)
	case flatbuf.MessageHeaderRecordBatch:
		w.recs = append(w.recs, blk)
	}

	return nil
}

func (w *pwriter) Close() error {
	var err error

	// write file footer
	err = w.updatePos()
	if err != nil {
		return fmt.Errorf("arrow/ipc: could not update position while in close: %w", err)
	}

	pos := w.pos
	err = writeFileFooter(w.schema, w.dicts, w.recs, w)
	if err != nil {
		return fmt.Errorf("arrow/ipc: could not write file footer: %w", err)
	}

	// write file footer length
	err = w.updatePos() // not strictly needed as we passed w to writeFileFooter...
	if err != nil {
		return fmt.Errorf("arrow/ipc: could not compute file footer length: %w", err)
	}

	size := w.pos - pos
	if size <= 0 {
		return fmt.Errorf("arrow/ipc: invalid file footer size (size=%d)", size)
	}

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(size))
	_, err = w.Write(buf)
	if err != nil {
		return fmt.Errorf("arrow/ipc: could not write file footer size: %w", err)
	}

	_, err = w.Write(Magic)
	if err != nil {
		return fmt.Errorf("arrow/ipc: could not write Arrow magic bytes: %w", err)
	}

	return nil
}

func (w *pwriter) updatePos() error {
	var err error
	w.pos, err = w.w.Seek(0, io.SeekCurrent)
	return err
}

func (w *pwriter) align(align int32) error {
	remainder := paddedLength(w.pos, align) - w.pos
	if remainder == 0 {
		return nil
	}

	_, err := w.Write(paddingBytes[:int(remainder)])
	return err
}

func (w *pwriter) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.pos += int64(n)
	return n, err
}

func writeIPCPayload(w io.Writer, p Payload) (int, error) {
	n, err := writeMessage(p.meta, kArrowIPCAlignment, w)
	if err != nil {
		return n, err
	}

	// now write the buffers
	for _, buf := range p.body {
		var (
			size    int64
			padding int64
		)

		// the buffer might be null if we are handling zero row lengths.
		if buf != nil {
			size = int64(buf.Len())
			padding = bitutil.CeilByte64(size) - size
		}

		if size > 0 {
			_, err = w.Write(buf.Bytes())
			if err != nil {
				return n, fmt.Errorf("arrow/ipc: could not write payload message body: %w", err)
			}
		}

		if padding > 0 {
			_, err = w.Write(paddingBytes[:padding])
			if err != nil {
				return n, fmt.Errorf("arrow/ipc: could not write payload message padding: %w", err)
			}
		}
	}

	return n, err
}

// Payload is the underlying message object which is passed to the payload writer
// for actually writing out ipc messages
type Payload struct {
	msg  MessageType
	meta *memory.Buffer
	body []*memory.Buffer
	size int64 // length of body
}

// Meta returns the buffer containing the metadata for this payload,
// callers must call Release on the buffer
func (p *Payload) Meta() *memory.Buffer {
	if p.meta != nil {
		p.meta.Retain()
	}
	return p.meta
}

// SerializeBody serializes the body buffers and writes them to the provided
// writer.
func (p *Payload) SerializeBody(w io.Writer) error {
	for _, data := range p.body {
		if data == nil {
			continue
		}

		size := int64(data.Len())
		padding := bitutil.CeilByte64(size) - size
		if size > 0 {
			if _, err := w.Write(data.Bytes()); err != nil {
				return fmt.Errorf("arrow/ipc: could not write payload message body: %w", err)
			}

			if padding > 0 {
				if _, err := w.Write(paddingBytes[:padding]); err != nil {
					return fmt.Errorf("arrow/ipc: could not write payload message padding bytes: %w", err)
				}
			}
		}
	}
	return nil
}

func (p *Payload) Release() {
	if p.meta != nil {
		p.meta.Release()
		p.meta = nil
	}
	for i, b := range p.body {
		if b == nil {
			continue
		}
		b.Release()
		p.body[i] = nil
	}
}

type payloads []Payload

func (ps payloads) Release() {
	for i := range ps {
		ps[i].Release()
	}
}

// FileWriter is an Arrow file writer.
type FileWriter struct {
	w io.WriteSeeker

	mem memory.Allocator

	header struct {
		started bool
		offset  int64
	}

	footer struct {
		written bool
	}

	pw PayloadWriter

	schema          *arrow.Schema
	mapper          dictutils.Mapper
	codec           flatbuf.CompressionType
	compressNP      int
	minSpaceSavings *float64

	// map of the last written dictionaries by id
	// so we can avoid writing the same dictionary over and over
	// also needed for correctness when writing IPC format which
	// does not allow replacements or deltas.
	lastWrittenDicts map[int64]arrow.Array
}

// NewFileWriter opens an Arrow file using the provided writer w.
func NewFileWriter(w io.WriteSeeker, opts ...Option) (*FileWriter, error) {
	var (
		cfg = newConfig(opts...)
		err error
	)

	f := FileWriter{
		w:               w,
		pw:              &pwriter{w: w, schema: cfg.schema, pos: -1},
		mem:             cfg.alloc,
		schema:          cfg.schema,
		codec:           cfg.codec,
		compressNP:      cfg.compressNP,
		minSpaceSavings: cfg.minSpaceSavings,
	}

	pos, err := f.w.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("arrow/ipc: could not seek current position: %w", err)
	}
	f.header.offset = pos

	return &f, err
}

func (f *FileWriter) Close() error {
	err := f.checkStarted()
	if err != nil {
		return fmt.Errorf("arrow/ipc: could not write empty file: %w", err)
	}

	if f.footer.written {
		return nil
	}

	err = f.pw.Close()
	if err != nil {
		return fmt.Errorf("arrow/ipc: could not close payload writer: %w", err)
	}
	f.footer.written = true

	return nil
}

func (f *FileWriter) Write(rec arrow.Record) error {
	schema := rec.Schema()
	if schema == nil || !schema.Equal(f.schema) {
		return errInconsistentSchema
	}

	if err := f.checkStarted(); err != nil {
		return fmt.Errorf("arrow/ipc: could not write header: %w", err)
	}

	const allow64b = true
	var (
		data = Payload{msg: MessageRecordBatch}
		enc  = newRecordEncoder(f.mem, 0, kMaxNestingDepth, allow64b, f.codec, f.compressNP, f.minSpaceSavings)
	)
	defer data.Release()

	err := writeDictionaryPayloads(f.mem, rec, true, false, &f.mapper, f.lastWrittenDicts, f.pw, enc)
	if err != nil {
		return fmt.Errorf("arrow/ipc: failure writing dictionary batches: %w", err)
	}

	enc.reset()
	if err := enc.Encode(&data, rec); err != nil {
		return fmt.Errorf("arrow/ipc: could not encode record to payload: %w", err)
	}

	return f.pw.WritePayload(data)
}

func (f *FileWriter) checkStarted() error {
	if !f.header.started {
		return f.start()
	}
	return nil
}

func (f *FileWriter) start() error {
	f.header.started = true
	err := f.pw.Start()
	if err != nil {
		return err
	}

	f.mapper.ImportSchema(f.schema)
	f.lastWrittenDicts = make(map[int64]arrow.Array)

	// write out schema payloads
	ps := payloadFromSchema(f.schema, f.mem, &f.mapper)
	defer ps.Release()

	for _, data := range ps {
		err = f.pw.WritePayload(data)
		if err != nil {
			return err
		}
	}

	return nil
}
