package lz4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ReaderLegacy implements the LZ4Demo frame decoder.
// The Header is set after the first call to Read().
type ReaderLegacy struct {
	Header
	// Handler called when a block has been successfully read.
	// It provides the number of bytes read.
	OnBlockDone func(size int)

	lastBlock bool
	buf       [8]byte   // Scrap buffer.
	pos       int64     // Current position in src.
	src       io.Reader // Source.
	zdata     []byte    // Compressed data.
	data      []byte    // Uncompressed data.
	idx       int       // Index of unread bytes into data.
	skip      int64     // Bytes to skip before next read.
	dpos      int64     // Position in dest
}

// NewReaderLegacy returns a new LZ4Demo frame decoder.
// No access to the underlying io.Reader is performed.
func NewReaderLegacy(src io.Reader) *ReaderLegacy {
	r := &ReaderLegacy{src: src}
	return r
}

// readHeader checks the frame magic number and parses the frame descriptoz.
// Skippable frames are supported even as a first frame although the LZ4
// specifications recommends skippable frames not to be used as first frames.
func (z *ReaderLegacy) readLegacyHeader() error {
	z.lastBlock = false
	magic, err := z.readUint32()
	if err != nil {
		z.pos += 4
		if err == io.ErrUnexpectedEOF {
			return io.EOF
		}
		return err
	}
	if magic != frameMagicLegacy {
		return ErrInvalid
	}
	z.pos += 4

	// Legacy has fixed 8MB blocksizes
	// https://github.com/lz4/lz4/blob/dev/doc/lz4_Frame_format.md#legacy-frame
	bSize := blockSize4M * 2

	// Allocate the compressed/uncompressed buffers.
	// The compressed buffer cannot exceed the uncompressed one.
	if n := 2 * bSize; cap(z.zdata) < n {
		z.zdata = make([]byte, n, n)
	}
	if debugFlag {
		debug("header block max size size=%d", bSize)
	}
	z.zdata = z.zdata[:bSize]
	z.data = z.zdata[:cap(z.zdata)][bSize:]
	z.idx = len(z.data)

	z.Header.done = true
	if debugFlag {
		debug("header read: %v", z.Header)
	}

	return nil
}

// Read decompresses data from the underlying source into the supplied buffer.
//
// Since there can be multiple streams concatenated, Header values may
// change between calls to Read(). If that is the case, no data is actually read from
// the underlying io.Reader, to allow for potential input buffer resizing.
func (z *ReaderLegacy) Read(buf []byte) (int, error) {
	if debugFlag {
		debug("Read buf len=%d", len(buf))
	}
	if !z.Header.done {
		if err := z.readLegacyHeader(); err != nil {
			return 0, err
		}
		if debugFlag {
			debug("header read OK compressed buffer %d / %d uncompressed buffer %d : %d index=%d",
				len(z.zdata), cap(z.zdata), len(z.data), cap(z.data), z.idx)
		}
	}

	if len(buf) == 0 {
		return 0, nil
	}

	if z.idx == len(z.data) {
		// No data ready for reading, process the next block.
		if debugFlag {
			debug("  reading block from writer %d %d", z.idx, blockSize4M*2)
		}

		// Reset uncompressed buffer
		z.data = z.zdata[:cap(z.zdata)][len(z.zdata):]

		bLen, err := z.readUint32()
		if err != nil {
			return 0, err
		}
		if debugFlag {
			debug("   bLen %d (0x%x) offset = %d (0x%x)", bLen, bLen, z.pos, z.pos)
		}
		z.pos += 4

		// Legacy blocks are always compressed, even when detrimental
		if debugFlag {
			debug("   compressed block size %d", bLen)
		}

		if int(bLen) > cap(z.data) {
			return 0, fmt.Errorf("lz4: invalid block size: %d", bLen)
		}
		zdata := z.zdata[:bLen]
		if _, err := io.ReadFull(z.src, zdata); err != nil {
			return 0, err
		}
		z.pos += int64(bLen)

		n, err := UncompressBlock(zdata, z.data)
		if err != nil {
			return 0, err
		}

		z.data = z.data[:n]
		if z.OnBlockDone != nil {
			z.OnBlockDone(n)
		}

		z.idx = 0

		// Legacy blocks are fixed to 8MB, if we read a decompressed block smaller than this
		// it means we've reached the end...
		if n < blockSize4M*2 {
			z.lastBlock = true
		}
	}

	if z.skip > int64(len(z.data[z.idx:])) {
		z.skip -= int64(len(z.data[z.idx:]))
		z.dpos += int64(len(z.data[z.idx:]))
		z.idx = len(z.data)
		return 0, nil
	}

	z.idx += int(z.skip)
	z.dpos += z.skip
	z.skip = 0

	n := copy(buf, z.data[z.idx:])
	z.idx += n
	z.dpos += int64(n)
	if debugFlag {
		debug("%v] copied %d bytes to input (%d:%d)", z.lastBlock, n, z.idx, len(z.data))
	}
	if z.lastBlock && len(z.data) == z.idx {
		return n, io.EOF
	}
	return n, nil
}

// Seek implements io.Seeker, but supports seeking forward from the current
// position only. Any other seek will return an error. Allows skipping output
// bytes which aren't needed, which in some scenarios is faster than reading
// and discarding them.
// Note this may cause future calls to Read() to read 0 bytes if all of the
// data they would have returned is skipped.
func (z *ReaderLegacy) Seek(offset int64, whence int) (int64, error) {
	if offset < 0 || whence != io.SeekCurrent {
		return z.dpos + z.skip, ErrUnsupportedSeek
	}
	z.skip += offset
	return z.dpos + z.skip, nil
}

// Reset discards the Reader's state and makes it equivalent to the
// result of its original state from NewReader, but reading from r instead.
// This permits reusing a Reader rather than allocating a new one.
func (z *ReaderLegacy) Reset(r io.Reader) {
	z.Header = Header{}
	z.pos = 0
	z.src = r
	z.zdata = z.zdata[:0]
	z.data = z.data[:0]
	z.idx = 0
}

// readUint32 reads an uint32 into the supplied buffer.
// The idea is to make use of the already allocated buffers avoiding additional allocations.
func (z *ReaderLegacy) readUint32() (uint32, error) {
	buf := z.buf[:4]
	_, err := io.ReadFull(z.src, buf)
	x := binary.LittleEndian.Uint32(buf)
	return x, err
}
