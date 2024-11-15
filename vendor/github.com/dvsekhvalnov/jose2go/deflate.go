package jose

import (
	"bytes"
	"compress/flate"
	"errors"
	"io"
)

var ErrSizeExceeded = errors.New("Deflate stream size exceeded limit.")

func init() {
	// 250Kb limited decompression buffer
	RegisterJwc(NewDeflate(250 * 1024))
}

// Deflate compression algorithm implementation
type Deflate struct {
	maxBufferSizeBytes int64
}

func NewDeflate(maxBufferSizeBytes int64) JwcAlgorithm {
	return &Deflate{
		maxBufferSizeBytes: maxBufferSizeBytes,
	}
}

func (alg *Deflate) Name() string {
	return DEF
}

func (alg *Deflate) Compress(plainText []byte) []byte {
	var buf bytes.Buffer
	deflate, _ := flate.NewWriter(&buf, 8) //level=DEFLATED

	deflate.Write(plainText)
	deflate.Close()

	return buf.Bytes()
}

func (alg *Deflate) Decompress(compressedText []byte) ([]byte, error) {
	enflated, err := io.ReadAll(
		newMaxBytesReader(alg.maxBufferSizeBytes,
			flate.NewReader(
				bytes.NewReader(compressedText))))

	return enflated, err
}

// Max bytes reader
type maxBytesReader struct {
	reader io.Reader
	limit  int64
}

func newMaxBytesReader(limit int64, r io.Reader) io.Reader {
	return &maxBytesReader{reader: r, limit: limit}
}

func (mbr *maxBytesReader) Read(p []byte) (n int, err error) {
	if mbr.limit <= 0 {
		return 0, ErrSizeExceeded
	}

	if int64(len(p)) > mbr.limit {
		p = p[0:mbr.limit]
	}

	n, err = mbr.reader.Read(p)
	mbr.limit -= int64(n)
	return
}
