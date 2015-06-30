package s3manager

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/s3"
)

// The default range of bytes to get at a time when using Download().
var DefaultDownloadPartSize int64 = 1024 * 1024 * 5

// The default number of goroutines to spin up when using Download().
var DefaultDownloadConcurrency = 5

// The default set of options used when opts is nil in Download().
var DefaultDownloadOptions = &DownloadOptions{
	PartSize:    DefaultDownloadPartSize,
	Concurrency: DefaultDownloadConcurrency,
}

// DownloadOptions keeps tracks of extra options to pass to an Download() call.
type DownloadOptions struct {
	// The buffer size (in bytes) to use when buffering data into chunks and
	// sending them as parts to S3. The minimum allowed part size is 5MB, and
	// if this value is set to zero, the DefaultPartSize value will be used.
	PartSize int64

	// The number of goroutines to spin up in parallel when sending parts.
	// If this is set to zero, the DefaultConcurrency value will be used.
	Concurrency int

	// An S3 client to use when performing downloads. Leave this as nil to use
	// a default client.
	S3 *s3.S3
}

// NewDownloader creates a new Downloader structure that downloads an object
// from S3 in concurrent chunks. Pass in an optional DownloadOptions struct
// to customize the downloader behavior.
func NewDownloader(opts *DownloadOptions) *Downloader {
	if opts == nil {
		opts = DefaultDownloadOptions
	}
	return &Downloader{opts: opts}
}

// The Downloader structure that calls Download(). It is safe to call Download()
// on this structure for multiple objects and across concurrent goroutines.
type Downloader struct {
	opts *DownloadOptions
}

// Download downloads an object in S3 and writes the payload into w using
// concurrent GET requests.
//
// It is safe to call this method for multiple objects and across concurrent
// goroutines.
func (d *Downloader) Download(w io.WriterAt, input *s3.GetObjectInput) (n int64, err error) {
	impl := downloader{w: w, in: input, opts: *d.opts}
	return impl.download()
}

// downloader is the implementation structure used internally by Downloader.
type downloader struct {
	opts DownloadOptions
	in   *s3.GetObjectInput
	w    io.WriterAt

	wg sync.WaitGroup
	m  sync.Mutex

	pos        int64
	totalBytes int64
	written    int64
	err        error
}

// init initializes the downloader with default options.
func (d *downloader) init() {
	d.totalBytes = -1

	if d.opts.Concurrency == 0 {
		d.opts.Concurrency = DefaultDownloadConcurrency
	}

	if d.opts.PartSize == 0 {
		d.opts.PartSize = DefaultDownloadPartSize
	}

	if d.opts.S3 == nil {
		d.opts.S3 = s3.New(nil)
	}
}

// download performs the implementation of the object download across ranged
// GETs.
func (d *downloader) download() (n int64, err error) {
	d.init()

	// Spin up workers
	ch := make(chan dlchunk, d.opts.Concurrency)
	for i := 0; i < d.opts.Concurrency; i++ {
		d.wg.Add(1)
		go d.downloadPart(ch)
	}

	// Assign work
	for d.geterr() == nil {
		if d.pos != 0 {
			// This is not the first chunk, let's wait until we know the total
			// size of the payload so we can see if we have read the entire
			// object.
			total := d.getTotalBytes()

			if total < 0 {
				// Total has not yet been set, so sleep and loop around while
				// waiting for our first worker to resolve this value.
				time.Sleep(10 * time.Millisecond)
				continue
			} else if d.pos >= total {
				break // We're finished queueing chunks
			}
		}

		// Queue the next range of bytes to read.
		ch <- dlchunk{w: d.w, start: d.pos, size: d.opts.PartSize}
		d.pos += d.opts.PartSize
	}

	// Wait for completion
	close(ch)
	d.wg.Wait()

	// Return error
	return d.written, d.err
}

// downloadPart is an individual goroutine worker reading from the ch channel
// and performing a GetObject request on the data with a given byte range.
//
// If this is the first worker, this operation also resolves the total number
// of bytes to be read so that the worker manager knows when it is finished.
func (d *downloader) downloadPart(ch chan dlchunk) {
	defer d.wg.Done()

	for {
		chunk, ok := <-ch

		if !ok {
			break
		}

		if d.geterr() == nil {
			// Get the next byte range of data
			in := &s3.GetObjectInput{}
			awsutil.Copy(in, d.in)
			rng := fmt.Sprintf("bytes=%d-%d",
				chunk.start, chunk.start+chunk.size-1)
			in.Range = &rng

			resp, err := d.opts.S3.GetObject(in)
			if err != nil {
				d.seterr(err)
			} else {
				d.setTotalBytes(resp) // Set total if not yet set.

				n, err := io.Copy(&chunk, resp.Body)
				resp.Body.Close()

				if err != nil {
					d.seterr(err)
				}
				d.incrwritten(n)
			}
		}
	}
}

// getTotalBytes is a thread-safe getter for retrieving the total byte status.
func (d *downloader) getTotalBytes() int64 {
	d.m.Lock()
	defer d.m.Unlock()

	return d.totalBytes
}

// getTotalBytes is a thread-safe setter for setting the total byte status.
func (d *downloader) setTotalBytes(resp *s3.GetObjectOutput) {
	d.m.Lock()
	defer d.m.Unlock()

	if d.totalBytes >= 0 {
		return
	}

	parts := strings.Split(*resp.ContentRange, "/")
	total, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	if err != nil {
		d.err = err
		return
	}

	d.totalBytes = total
}

func (d *downloader) incrwritten(n int64) {
	d.m.Lock()
	defer d.m.Unlock()

	d.written += n
}

// geterr is a thread-safe getter for the error object
func (d *downloader) geterr() error {
	d.m.Lock()
	defer d.m.Unlock()

	return d.err
}

// seterr is a thread-safe setter for the error object
func (d *downloader) seterr(e error) {
	d.m.Lock()
	defer d.m.Unlock()

	d.err = e
}

// dlchunk represents a single chunk of data to write by the worker routine.
// This structure also implements an io.SectionReader style interface for
// io.WriterAt, effectively making it an io.SectionWriter (which does not
// exist).
type dlchunk struct {
	w     io.WriterAt
	start int64
	size  int64
	cur   int64
}

// Write wraps io.WriterAt for the dlchunk, writing from the dlchunk's start
// position to its end (or EOF).
func (c *dlchunk) Write(p []byte) (n int, err error) {
	if c.cur >= c.size {
		return 0, io.EOF
	}

	n, err = c.w.WriteAt(p, c.start+c.cur)
	c.cur += int64(n)

	return
}
