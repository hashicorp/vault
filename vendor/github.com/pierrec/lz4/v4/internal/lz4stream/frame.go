// Package lz4stream provides the types that support reading and writing LZ4 data streams.
package lz4stream

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pierrec/lz4/v4/internal/lz4block"
	"github.com/pierrec/lz4/v4/internal/lz4errors"
	"github.com/pierrec/lz4/v4/internal/xxh32"
)

//go:generate go run gen.go

const (
	frameMagic     uint32 = 0x184D2204
	frameSkipMagic uint32 = 0x184D2A50
)

func NewFrame() *Frame {
	return &Frame{}
}

type Frame struct {
	buf        [15]byte // frame descriptor needs at most 4(magic)+4+8+1=11 bytes
	Magic      uint32
	Descriptor FrameDescriptor
	Blocks     Blocks
	Checksum   uint32
	checksum   xxh32.XXHZero
}

// Reset allows reusing the Frame.
// The Descriptor configuration is not modified.
func (f *Frame) Reset(num int) {
	f.Magic = 0
	f.Descriptor.Checksum = 0
	f.Descriptor.ContentSize = 0
	_ = f.Blocks.closeW(f, num)
	f.Checksum = 0
}

func (f *Frame) InitW(dst io.Writer, num int) {
	f.Magic = frameMagic
	f.Descriptor.initW()
	f.Blocks.initW(f, dst, num)
	f.checksum.Reset()
}

func (f *Frame) CloseW(dst io.Writer, num int) error {
	if err := f.Blocks.closeW(f, num); err != nil {
		return err
	}
	buf := f.buf[:0]
	// End mark (data block size of uint32(0)).
	buf = append(buf, 0, 0, 0, 0)
	if f.Descriptor.Flags.ContentChecksum() {
		buf = f.checksum.Sum(buf)
	}
	_, err := dst.Write(buf)
	return err
}

func (f *Frame) InitR(src io.Reader) error {
	if f.Magic > 0 {
		// Header already read.
		return nil
	}

newFrame:
	var err error
	if f.Magic, err = f.readUint32(src); err != nil {
		return err
	}
	switch m := f.Magic; {
	case m == frameMagic:
	// All 16 values of frameSkipMagic are valid.
	case m>>8 == frameSkipMagic>>8:
		var skip uint32
		if err := binary.Read(src, binary.LittleEndian, &skip); err != nil {
			return err
		}
		if _, err := io.CopyN(ioutil.Discard, src, int64(skip)); err != nil {
			return err
		}
		goto newFrame
	default:
		return lz4errors.ErrInvalidFrame
	}
	if err := f.Descriptor.initR(f, src); err != nil {
		return err
	}
	f.Blocks.initR(f)
	f.checksum.Reset()
	return nil
}

func (f *Frame) CloseR(src io.Reader) (err error) {
	if !f.Descriptor.Flags.ContentChecksum() {
		return nil
	}
	if f.Checksum, err = f.readUint32(src); err != nil {
		return err
	}
	if c := f.checksum.Sum32(); c != f.Checksum {
		return fmt.Errorf("%w: got %x; expected %x", lz4errors.ErrInvalidFrameChecksum, c, f.Checksum)
	}
	return nil
}

type FrameDescriptor struct {
	Flags       DescriptorFlags
	ContentSize uint64
	Checksum    uint8
}

func (fd *FrameDescriptor) initW() {
	fd.Flags.VersionSet(1)
	fd.Flags.BlockIndependenceSet(true)
}

func (fd *FrameDescriptor) Write(f *Frame, dst io.Writer) error {
	if fd.Checksum > 0 {
		// Header already written.
		return nil
	}

	buf := f.buf[:4+2]
	// Write the magic number here even though it belongs to the Frame.
	binary.LittleEndian.PutUint32(buf, f.Magic)
	binary.LittleEndian.PutUint16(buf[4:], uint16(fd.Flags))

	if fd.Flags.Size() {
		buf = buf[:4+2+8]
		binary.LittleEndian.PutUint64(buf[4+2:], fd.ContentSize)
	}
	fd.Checksum = descriptorChecksum(buf[4:])
	buf = append(buf, fd.Checksum)

	_, err := dst.Write(buf)
	return err
}

func (fd *FrameDescriptor) initR(f *Frame, src io.Reader) error {
	// Read the flags and the checksum, hoping that there is not content size.
	buf := f.buf[:3]
	if _, err := io.ReadFull(src, buf); err != nil {
		return err
	}
	descr := binary.LittleEndian.Uint16(buf)
	fd.Flags = DescriptorFlags(descr)
	if fd.Flags.Size() {
		// Append the 8 missing bytes.
		buf = buf[:3+8]
		if _, err := io.ReadFull(src, buf[3:]); err != nil {
			return err
		}
		fd.ContentSize = binary.LittleEndian.Uint64(buf[2:])
	}
	fd.Checksum = buf[len(buf)-1] // the checksum is the last byte
	buf = buf[:len(buf)-1]        // all descriptor fields except checksum
	if c := descriptorChecksum(buf); fd.Checksum != c {
		return fmt.Errorf("%w: got %x; expected %x", lz4errors.ErrInvalidHeaderChecksum, c, fd.Checksum)
	}
	// Validate the elements that can be.
	if idx := fd.Flags.BlockSizeIndex(); !idx.IsValid() {
		return lz4errors.ErrOptionInvalidBlockSize
	}
	return nil
}

func descriptorChecksum(buf []byte) byte {
	return byte(xxh32.ChecksumZero(buf) >> 8)
}

type Blocks struct {
	Block  *FrameDataBlock
	Blocks chan chan *FrameDataBlock
	err    error
}

func (b *Blocks) initW(f *Frame, dst io.Writer, num int) {
	size := f.Descriptor.Flags.BlockSizeIndex()
	if num == 1 {
		b.Blocks = nil
		b.Block = NewFrameDataBlock(size)
		return
	}
	b.Block = nil
	if cap(b.Blocks) != num {
		b.Blocks = make(chan chan *FrameDataBlock, num)
	}
	// goroutine managing concurrent block compression goroutines.
	go func() {
		// Process next block compression item.
		for c := range b.Blocks {
			// Read the next compressed block result.
			// Waiting here ensures that the blocks are output in the order they were sent.
			// The incoming channel is always closed as it indicates to the caller that
			// the block has been processed.
			block := <-c
			if block == nil {
				// Notify the block compression routine that we are done with its result.
				// This is used when a sentinel block is sent to terminate the compression.
				close(c)
				return
			}
			// Do not attempt to write the block upon any previous failure.
			if b.err == nil {
				// Write the block.
				if err := block.Write(f, dst); err != nil && b.err == nil {
					// Keep the first error.
					b.err = err
					// All pending compression goroutines need to shut down, so we need to keep going.
				}
			}
			close(c)
		}
	}()
}

func (b *Blocks) closeW(f *Frame, num int) error {
	if num == 1 {
		if b.Block == nil {
			// Not initialized yet.
			return nil
		}
		b.Block.CloseW(f)
		return nil
	}
	if b.Blocks == nil {
		// Not initialized yet.
		return nil
	}
	c := make(chan *FrameDataBlock)
	b.Blocks <- c
	c <- nil
	<-c
	err := b.err
	b.err = nil
	return err
}

func (b *Blocks) initR(f *Frame) {
	size := f.Descriptor.Flags.BlockSizeIndex()
	b.Block = NewFrameDataBlock(size)
}

func NewFrameDataBlock(size lz4block.BlockSizeIndex) *FrameDataBlock {
	buf := size.Get()
	return &FrameDataBlock{Data: buf, data: buf}
}

type FrameDataBlock struct {
	Size     DataBlockSize
	Data     []byte // compressed or uncompressed data (.data or .src)
	Checksum uint32
	data     []byte // buffer for compressed data
	src      []byte // uncompressed data
}

func (b *FrameDataBlock) CloseW(f *Frame) {
	if b.data != nil {
		// Block was not already closed.
		size := f.Descriptor.Flags.BlockSizeIndex()
		size.Put(b.data)
		b.Data = nil
		b.data = nil
		b.src = nil
	}
}

// Block compression errors are ignored since the buffer is sized appropriately.
func (b *FrameDataBlock) Compress(f *Frame, src []byte, level lz4block.CompressionLevel) *FrameDataBlock {
	data := b.data[:len(src)] // trigger the incompressible flag in CompressBlock
	var n int
	switch level {
	case lz4block.Fast:
		n, _ = lz4block.CompressBlock(src, data)
	default:
		n, _ = lz4block.CompressBlockHC(src, data, level)
	}
	if n == 0 {
		b.Size.UncompressedSet(true)
		b.Data = src
	} else {
		b.Size.UncompressedSet(false)
		b.Data = data[:n]
	}
	b.Size.sizeSet(len(b.Data))
	b.src = src // keep track of the source for content checksum

	if f.Descriptor.Flags.BlockChecksum() {
		b.Checksum = xxh32.ChecksumZero(src)
	}
	return b
}

func (b *FrameDataBlock) Write(f *Frame, dst io.Writer) error {
	if f.Descriptor.Flags.ContentChecksum() {
		_, _ = f.checksum.Write(b.src)
	}
	buf := f.buf[:]
	binary.LittleEndian.PutUint32(buf, uint32(b.Size))
	if _, err := dst.Write(buf[:4]); err != nil {
		return err
	}

	if _, err := dst.Write(b.Data); err != nil {
		return err
	}

	if b.Checksum == 0 {
		return nil
	}
	binary.LittleEndian.PutUint32(buf, b.Checksum)
	_, err := dst.Write(buf[:4])
	return err
}

func (b *FrameDataBlock) Uncompress(f *Frame, src io.Reader, dst []byte) (int, error) {
	x, err := f.readUint32(src)
	if err != nil {
		return 0, err
	}
	b.Size = DataBlockSize(x)
	if b.Size == 0 {
		// End of frame reached.
		return 0, io.EOF
	}

	isCompressed := !b.Size.Uncompressed()
	size := b.Size.size()
	var data []byte
	if isCompressed {
		// Data is first copied into b.Data and then it will get uncompressed into dst.
		data = b.Data
	} else {
		// Data is directly copied into dst as it is not compressed.
		data = dst
	}
	data = data[:size]
	if _, err := io.ReadFull(src, data); err != nil {
		return 0, err
	}
	if isCompressed {
		n, err := lz4block.UncompressBlock(data, dst)
		if err != nil {
			return 0, err
		}
		data = dst[:n]
	}

	if f.Descriptor.Flags.BlockChecksum() {
		var err error
		if b.Checksum, err = f.readUint32(src); err != nil {
			return 0, err
		}
		if c := xxh32.ChecksumZero(data); c != b.Checksum {
			return 0, fmt.Errorf("%w: got %x; expected %x", lz4errors.ErrInvalidBlockChecksum, c, b.Checksum)
		}
	}
	if f.Descriptor.Flags.ContentChecksum() {
		_, _ = f.checksum.Write(data)
	}
	return len(data), nil
}

func (f *Frame) readUint32(r io.Reader) (x uint32, err error) {
	if _, err = io.ReadFull(r, f.buf[:4]); err != nil {
		return
	}
	x = binary.LittleEndian.Uint32(f.buf[:4])
	return
}
