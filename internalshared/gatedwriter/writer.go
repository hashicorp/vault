package gatedwriter

import (
	"bytes"
	"io"
	"sync"
)

// Writer is an io.Writer implementation that buffers all of its
// data into an internal buffer until it is told to let data through.
type Writer struct {
	writer io.Writer

	buf   bytes.Buffer
	flush bool
	lock  sync.Mutex
}

func NewWriter(underlying io.Writer) *Writer {
	return &Writer{writer: underlying}
}

// Flush tells the Writer to flush any buffered data and to stop
// buffering.
func (w *Writer) Flush() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.flush = true
	_, err := w.buf.WriteTo(w.writer)
	return err
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.flush {
		return w.writer.Write(p)
	}

	return w.buf.Write(p)
}
