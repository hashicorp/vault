package lz4

import (
	"encoding/binary"
	"io"
)

// WriterLegacy implements the LZ4Demo frame decoder.
type WriterLegacy struct {
	Header
	// Handler called when a block has been successfully read.
	// It provides the number of bytes read.
	OnBlockDone func(size int)

	dst       io.Writer    // Destination.
	data      []byte       // Data to be compressed + buffer for compressed data.
	idx       int          // Index into data.
	hashtable [winSize]int // Hash table used in CompressBlock().
}

// NewWriterLegacy returns a new LZ4 encoder for the legacy frame format.
// No access to the underlying io.Writer is performed.
// The supplied Header is checked at the first Write.
// It is ok to change it before the first Write but then not until a Reset() is performed.
func NewWriterLegacy(dst io.Writer) *WriterLegacy {
	z := new(WriterLegacy)
	z.Reset(dst)
	return z
}

// Write compresses data from the supplied buffer into the underlying io.Writer.
// Write does not return until the data has been written.
func (z *WriterLegacy) Write(buf []byte) (int, error) {
	if !z.Header.done {
		if err := z.writeHeader(); err != nil {
			return 0, err
		}
	}
	if debugFlag {
		debug("input buffer len=%d index=%d", len(buf), z.idx)
	}

	zn := len(z.data)
	var n int
	for len(buf) > 0 {
		if z.idx == 0 && len(buf) >= zn {
			// Avoid a copy as there is enough data for a block.
			if err := z.compressBlock(buf[:zn]); err != nil {
				return n, err
			}
			n += zn
			buf = buf[zn:]
			continue
		}
		// Accumulate the data to be compressed.
		m := copy(z.data[z.idx:], buf)
		n += m
		z.idx += m
		buf = buf[m:]
		if debugFlag {
			debug("%d bytes copied to buf, current index %d", n, z.idx)
		}

		if z.idx < len(z.data) {
			// Buffer not filled.
			if debugFlag {
				debug("need more data for compression")
			}
			return n, nil
		}

		// Buffer full.
		if err := z.compressBlock(z.data); err != nil {
			return n, err
		}
		z.idx = 0
	}

	return n, nil
}

// writeHeader builds and writes the header to the underlying io.Writer.
func (z *WriterLegacy) writeHeader() error {
	// Legacy has fixed 8MB blocksizes
	// https://github.com/lz4/lz4/blob/dev/doc/lz4_Frame_format.md#legacy-frame
	bSize := 2 * blockSize4M

	buf := make([]byte, 2*bSize, 2*bSize)
	z.data = buf[:bSize] // Uncompressed buffer is the first half.

	z.idx = 0

	// Header consists of one mageic number, write it out.
	if err := binary.Write(z.dst, binary.LittleEndian, frameMagicLegacy); err != nil {
		return err
	}
	z.Header.done = true
	if debugFlag {
		debug("wrote header %v", z.Header)
	}

	return nil
}

// compressBlock compresses a block.
func (z *WriterLegacy) compressBlock(data []byte) error {
	bSize := 2 * blockSize4M
	zdata := z.data[bSize:cap(z.data)]
	// The compressed block size cannot exceed the input's.
	var zn int

	if level := z.Header.CompressionLevel; level != 0 {
		zn, _ = CompressBlockHC(data, zdata, level)
	} else {
		zn, _ = CompressBlock(data, zdata, z.hashtable[:])
	}

	if debugFlag {
		debug("block compression %d => %d", len(data), zn)
	}
	zdata = zdata[:zn]

	// Write the block.
	if err := binary.Write(z.dst, binary.LittleEndian, uint32(zn)); err != nil {
		return err
	}
	written, err := z.dst.Write(zdata)
	if err != nil {
		return err
	}
	if h := z.OnBlockDone; h != nil {
		h(written)
	}
	return nil
}

// Flush flushes any pending compressed data to the underlying writer.
// Flush does not return until the data has been written.
// If the underlying writer returns an error, Flush returns that error.
func (z *WriterLegacy) Flush() error {
	if debugFlag {
		debug("flush with index %d", z.idx)
	}
	if z.idx == 0 {
		return nil
	}

	data := z.data[:z.idx]
	z.idx = 0
	return z.compressBlock(data)
}

// Close closes the WriterLegacy, flushing any unwritten data to the underlying io.Writer, but does not close the underlying io.Writer.
func (z *WriterLegacy) Close() error {
	if !z.Header.done {
		if err := z.writeHeader(); err != nil {
			return err
		}
	}
	if err := z.Flush(); err != nil {
		return err
	}

	if debugFlag {
		debug("writing last empty block")
	}

	return nil
}

// Reset clears the state of the WriterLegacy z such that it is equivalent to its
// initial state from NewWriterLegacy, but instead writing to w.
// No access to the underlying io.Writer is performed.
func (z *WriterLegacy) Reset(w io.Writer) {
	z.Header.Reset()
	z.dst = w
	z.idx = 0
	// reset hashtable to ensure deterministic output.
	for i := range z.hashtable {
		z.hashtable[i] = 0
	}
}
