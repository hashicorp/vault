// Copyright 2014 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
	"unicode/utf8"

	"cloud.google.com/go/internal/trace"
)

// A Writer writes a Cloud Storage object.
type Writer struct {
	// ObjectAttrs are optional attributes to set on the object. Any attributes
	// must be initialized before the first Write call. Nil or zero-valued
	// attributes are ignored.
	ObjectAttrs

	// SendCRC32C specifies whether to transmit a CRC32C field. It should be set
	// to true in addition to setting the Writer's CRC32C field, because zero
	// is a valid CRC and normally a zero would not be transmitted.
	// If a CRC32C is sent, and the data written does not match the checksum,
	// the write will be rejected.
	//
	// Note: SendCRC32C must be set to true BEFORE the first call to
	// Writer.Write() in order to send the checksum. If it is set after that
	// point, the checksum will be ignored.
	SendCRC32C bool

	// ChunkSize controls the maximum number of bytes of the object that the
	// Writer will attempt to send to the server in a single request. Objects
	// smaller than the size will be sent in a single request, while larger
	// objects will be split over multiple requests. The value will be rounded up
	// to the nearest multiple of 256K. The default ChunkSize is 16MiB.
	//
	// Each Writer will internally allocate a buffer of size ChunkSize. This is
	// used to buffer input data and allow for the input to be sent again if a
	// request must be retried.
	//
	// If you upload small objects (< 16MiB), you should set ChunkSize
	// to a value slightly larger than the objects' sizes to avoid memory bloat.
	// This is especially important if you are uploading many small objects
	// concurrently. See
	// https://cloud.google.com/storage/docs/json_api/v1/how-tos/upload#size
	// for more information about performance trade-offs related to ChunkSize.
	//
	// If ChunkSize is set to zero, chunking will be disabled and the object will
	// be uploaded in a single request without the use of a buffer. This will
	// further reduce memory used during uploads, but will also prevent the writer
	// from retrying in case of a transient error from the server or resuming an
	// upload that fails midway through, since the buffer is required in order to
	// retry the failed request.
	//
	// ChunkSize must be set before the first Write call.
	ChunkSize int

	// ChunkRetryDeadline sets a per-chunk retry deadline for multi-chunk
	// resumable uploads.
	//
	// For uploads of larger files, the Writer will attempt to retry if the
	// request to upload a particular chunk fails with a transient error.
	// If a single chunk has been attempting to upload for longer than this
	// deadline and the request fails, it will no longer be retried, and the error
	// will be returned to the caller. This is only applicable for files which are
	// large enough to require a multi-chunk resumable upload. The default value
	// is 32s. Users may want to pick a longer deadline if they are using larger
	// values for ChunkSize or if they expect to have a slow or unreliable
	// internet connection.
	//
	// To set a deadline on the entire upload, use context timeout or
	// cancellation.
	ChunkRetryDeadline time.Duration

	// ForceEmptyContentType is an optional parameter that is used to disable
	// auto-detection of Content-Type. By default, if a blank Content-Type
	// is provided, then gax.DetermineContentType is called to sniff the type.
	ForceEmptyContentType bool

	// ProgressFunc can be used to monitor the progress of a large write
	// operation. If ProgressFunc is not nil and writing requires multiple
	// calls to the underlying service (see
	// https://cloud.google.com/storage/docs/json_api/v1/how-tos/resumable-upload),
	// then ProgressFunc will be invoked after each call with the number of bytes of
	// content copied so far.
	//
	// ProgressFunc should return quickly without blocking.
	ProgressFunc func(int64)

	ctx context.Context
	o   *ObjectHandle

	opened bool
	pw     *io.PipeWriter

	donec chan struct{} // closed after err and obj are set.
	obj   *ObjectAttrs

	mu  sync.Mutex
	err error
}

// Write appends to w. It implements the io.Writer interface.
//
// Since writes happen asynchronously, Write may return a nil
// error even though the write failed (or will fail). Always
// use the error returned from Writer.Close to determine if
// the upload was successful.
//
// Writes will be retried on transient errors from the server, unless
// Writer.ChunkSize has been set to zero.
func (w *Writer) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	werr := w.err
	w.mu.Unlock()
	if werr != nil {
		return 0, werr
	}
	if !w.opened {
		if err := w.openWriter(); err != nil {
			return 0, err
		}
	}
	n, err = w.pw.Write(p)
	if err != nil {
		w.mu.Lock()
		werr := w.err
		w.mu.Unlock()
		// Preserve existing functionality that when context is canceled, Write will return
		// context.Canceled instead of "io: read/write on closed pipe". This hides the
		// pipe implementation detail from users and makes Write seem as though it's an RPC.
		if errors.Is(werr, context.Canceled) || errors.Is(werr, context.DeadlineExceeded) {
			return n, werr
		}
	}
	return n, err
}

// Close completes the write operation and flushes any buffered data.
// If Close doesn't return an error, metadata about the written object
// can be retrieved by calling Attrs.
func (w *Writer) Close() error {
	if !w.opened {
		if err := w.openWriter(); err != nil {
			return err
		}
	}

	// Closing either the read or write causes the entire pipe to close.
	if err := w.pw.Close(); err != nil {
		return err
	}

	<-w.donec
	w.mu.Lock()
	defer w.mu.Unlock()
	trace.EndSpan(w.ctx, w.err)
	return w.err
}

func (w *Writer) openWriter() (err error) {
	if err := w.validateWriteAttrs(); err != nil {
		return err
	}
	if w.o.gen != defaultGen {
		return fmt.Errorf("storage: generation not supported on Writer, got %v", w.o.gen)
	}

	isIdempotent := w.o.conds != nil && (w.o.conds.GenerationMatch >= 0 || w.o.conds.DoesNotExist == true)
	opts := makeStorageOpts(isIdempotent, w.o.retry, w.o.userProject)
	params := &openWriterParams{
		ctx:                   w.ctx,
		chunkSize:             w.ChunkSize,
		chunkRetryDeadline:    w.ChunkRetryDeadline,
		bucket:                w.o.bucket,
		attrs:                 &w.ObjectAttrs,
		conds:                 w.o.conds,
		encryptionKey:         w.o.encryptionKey,
		sendCRC32C:            w.SendCRC32C,
		donec:                 w.donec,
		setError:              w.error,
		progress:              w.progress,
		setObj:                func(o *ObjectAttrs) { w.obj = o },
		forceEmptyContentType: w.ForceEmptyContentType,
	}
	if err := w.ctx.Err(); err != nil {
		return err // short-circuit
	}
	w.pw, err = w.o.c.tc.OpenWriter(params, opts...)
	if err != nil {
		return err
	}
	w.opened = true
	go w.monitorCancel()

	return nil
}

// monitorCancel is intended to be used as a background goroutine. It monitors the
// context, and when it observes that the context has been canceled, it manually
// closes things that do not take a context.
func (w *Writer) monitorCancel() {
	select {
	case <-w.ctx.Done():
		w.mu.Lock()
		werr := w.ctx.Err()
		w.err = werr
		w.mu.Unlock()

		// Closing either the read or write causes the entire pipe to close.
		w.CloseWithError(werr)
	case <-w.donec:
	}
}

// CloseWithError aborts the write operation with the provided error.
// CloseWithError always returns nil.
//
// Deprecated: cancel the context passed to NewWriter instead.
func (w *Writer) CloseWithError(err error) error {
	if !w.opened {
		return nil
	}
	return w.pw.CloseWithError(err)
}

// Attrs returns metadata about a successfully-written object.
// It's only valid to call it after Close returns nil.
func (w *Writer) Attrs() *ObjectAttrs {
	return w.obj
}

func (w *Writer) validateWriteAttrs() error {
	attrs := w.ObjectAttrs
	// Check the developer didn't change the object Name (this is unfortunate, but
	// we don't want to store an object under the wrong name).
	if attrs.Name != w.o.object {
		return fmt.Errorf("storage: Writer.Name %q does not match object name %q", attrs.Name, w.o.object)
	}
	if !utf8.ValidString(attrs.Name) {
		return fmt.Errorf("storage: object name %q is not valid UTF-8", attrs.Name)
	}
	if attrs.KMSKeyName != "" && w.o.encryptionKey != nil {
		return errors.New("storage: cannot use KMSKeyName with a customer-supplied encryption key")
	}
	if w.ChunkSize < 0 {
		return errors.New("storage: Writer.ChunkSize must be non-negative")
	}
	return nil
}

// progress is a convenience wrapper that reports write progress to the Writer
// ProgressFunc if it is set and progress is non-zero.
func (w *Writer) progress(p int64) {
	if w.ProgressFunc != nil && p != 0 {
		w.ProgressFunc(p)
	}
}

// error acquires the Writer's lock, sets the Writer's err to the given error,
// then relinquishes the lock.
func (w *Writer) error(err error) {
	w.mu.Lock()
	w.err = err
	w.mu.Unlock()
}
