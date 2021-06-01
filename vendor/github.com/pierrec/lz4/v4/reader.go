package lz4

import (
	"io"

	"github.com/pierrec/lz4/v4/internal/lz4errors"
	"github.com/pierrec/lz4/v4/internal/lz4stream"
)

var readerStates = []aState{
	noState:     newState,
	errorState:  newState,
	newState:    readState,
	readState:   closedState,
	closedState: newState,
}

// NewReader returns a new LZ4 frame decoder.
func NewReader(r io.Reader) *Reader {
	zr := &Reader{frame: lz4stream.NewFrame()}
	zr.state.init(readerStates)
	_ = zr.Apply(defaultOnBlockDone)
	zr.Reset(r)
	return zr
}

// Reader allows reading an LZ4 stream.
type Reader struct {
	state   _State
	src     io.Reader        // source reader
	frame   *lz4stream.Frame // frame being read
	data    []byte           // pending data
	idx     int              // size of pending data
	handler func(int)
}

func (*Reader) private() {}

func (r *Reader) Apply(options ...Option) (err error) {
	defer r.state.check(&err)
	switch r.state.state {
	case newState:
	case errorState:
		return r.state.err
	default:
		return lz4errors.ErrOptionClosedOrError
	}
	for _, o := range options {
		if err = o(r); err != nil {
			return
		}
	}
	return
}

// Size returns the size of the underlying uncompressed data, if set in the stream.
func (r *Reader) Size() int {
	switch r.state.state {
	case readState, closedState:
		if r.frame.Descriptor.Flags.Size() {
			return int(r.frame.Descriptor.ContentSize)
		}
	}
	return 0
}

func (r *Reader) init() error {
	return r.frame.InitR(r.src)
}

func (r *Reader) Read(buf []byte) (n int, err error) {
	defer r.state.check(&err)
	switch r.state.state {
	case readState:
	case closedState, errorState:
		return 0, r.state.err
	case newState:
		// First initialization.
		if err = r.init(); r.state.next(err) {
			return
		}
		size := r.frame.Descriptor.Flags.BlockSizeIndex()
		r.data = size.Get()
	default:
		return 0, r.state.fail()
	}
	if len(buf) == 0 {
		return
	}

	var bn int
	if r.idx > 0 {
		// Some left over data, use it.
		goto fillbuf
	}
	// No uncompressed data yet.
	r.data = r.data[:cap(r.data)]
	for len(buf) >= len(r.data) {
		// Input buffer large enough and no pending data: uncompress directly into it.
		switch bn, err = r.frame.Blocks.Block.Uncompress(r.frame, r.src, buf); err {
		case nil:
			r.handler(bn)
			n += bn
			buf = buf[bn:]
		case io.EOF:
			goto close
		default:
			return
		}
	}
	if n > 0 {
		// Some data was read, done for now.
		return
	}
	// Read the next block.
	switch bn, err = r.frame.Blocks.Block.Uncompress(r.frame, r.src, r.data); err {
	case nil:
		r.handler(bn)
		r.data = r.data[:bn]
		goto fillbuf
	case io.EOF:
	default:
		return
	}
close:
	if er := r.frame.CloseR(r.src); er != nil {
		err = er
	}
	r.Reset(nil)
	return
fillbuf:
	bn = copy(buf, r.data[r.idx:])
	n += bn
	r.idx += bn
	if r.idx == len(r.data) {
		// All data read, get ready for the next Read.
		r.idx = 0
	}
	return
}

// Reset clears the state of the Reader r such that it is equivalent to its
// initial state from NewReader, but instead writing to writer.
// No access to reader is performed.
//
// w.Close must be called before Reset.
func (r *Reader) Reset(reader io.Reader) {
	size := r.frame.Descriptor.Flags.BlockSizeIndex()
	size.Put(r.data)
	r.frame.Reset(1)
	r.src = reader
	r.data = nil
	r.idx = 0
	r.state.reset()
}

// WriteTo efficiently uncompresses the data from the Reader underlying source to w.
func (r *Reader) WriteTo(w io.Writer) (n int64, err error) {
	switch r.state.state {
	case closedState, errorState:
		return 0, r.state.err
	case newState:
		if err = r.init(); r.state.next(err) {
			return
		}
	default:
		return 0, r.state.fail()
	}
	defer r.state.nextd(&err)

	var bn int
	block := r.frame.Blocks.Block
	size := r.frame.Descriptor.Flags.BlockSizeIndex()
	data := size.Get()
	defer size.Put(data)
	for {
		switch bn, err = block.Uncompress(r.frame, r.src, data); err {
		case nil:
		case io.EOF:
			err = r.frame.CloseR(r.src)
			return
		default:
			return
		}
		r.handler(bn)
		bn, err = w.Write(data[:bn])
		n += int64(bn)
		if err != nil {
			return
		}
	}
}
