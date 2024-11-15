package magic

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"sync"
)

// A bufio.Reader pool to alleviate problems with memory allocations.
var readerPool = sync.Pool{
	New: func() any {
		// Initiate with empty source reader.
		return bufio.NewReader(nil)
	},
}

func newReader(r io.Reader) *bufio.Reader {
	br := readerPool.Get().(*bufio.Reader)
	br.Reset(r)
	return br
}

// Csv matches a comma-separated values file.
func Csv(raw []byte, limit uint32) bool {
	return sv(raw, ',', limit)
}

// Tsv matches a tab-separated values file.
func Tsv(raw []byte, limit uint32) bool {
	return sv(raw, '\t', limit)
}

func sv(in []byte, comma rune, limit uint32) bool {
	in = dropLastLine(in, limit)

	br := newReader(bytes.NewReader(in))
	defer readerPool.Put(br)
	r := csv.NewReader(br)
	r.Comma = comma
	r.ReuseRecord = true
	r.LazyQuotes = true
	r.Comment = '#'

	lines := 0
	for {
		_, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return false
		}
		lines++
	}

	return r.FieldsPerRecord > 1 && lines > 1
}

// dropLastLine drops the last incomplete line from b.
//
// mimetype limits itself to ReadLimit bytes when performing a detection.
// This means, for file formats like CSV for NDJSON, the last line of the input
// can be an incomplete line.
func dropLastLine(b []byte, readLimit uint32) []byte {
	if readLimit == 0 || uint32(len(b)) < readLimit {
		return b
	}
	for i := len(b) - 1; i > 0; i-- {
		if b[i] == '\n' {
			return b[:i]
		}
	}
	return b
}
