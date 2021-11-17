// Package mimetype uses magic number signatures to detect the MIME type of a file.
package mimetype

import (
	"io"
	"io/ioutil"
	"mime"
	"os"
)

// readLimit is the maximum number of bytes from the input used when detecting.
var readLimit uint32 = 3072

// Detect returns the MIME type found from the provided byte slice.
//
// The result is always a valid MIME type, with application/octet-stream
// returned when identification failed.
func Detect(in []byte) *MIME {
	rootMu.RLock()
	defer rootMu.RUnlock()
	if readLimit > 0 && len(in) > int(readLimit) {
		in = in[:readLimit]
	}
	return root.match(in, readLimit)
}

// DetectReader returns the MIME type of the provided reader.
//
// The result is always a valid MIME type, with application/octet-stream
// returned when identification failed with or without an error.
// Any error returned is related to the reading from the input reader.
//
// DetectReader assumes the reader offset is at the start. If the input is an
// io.ReadSeeker you previously read from, it should be rewinded before detection:
//  reader.Seek(0, io.SeekStart)
func DetectReader(r io.Reader) (*MIME, error) {
	rootMu.RLock()
	defer rootMu.RUnlock()
	var in []byte
	var err error

	if readLimit == 0 {
		in, err = ioutil.ReadAll(r)
		if err != nil {
			return root, err
		}
	} else {
		// io.UnexpectedEOF means len(r) < len(in). It is not an error in this case,
		// it just means the input file is smaller than the allocated bytes slice.
		n := 0
		in = make([]byte, readLimit)
		n, err = io.ReadFull(r, in)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return root, err
		}
		in = in[:n]
	}

	return root.match(in, readLimit), nil
}

// DetectFile returns the MIME type of the provided file.
//
// The result is always a valid MIME type, with application/octet-stream
// returned when identification failed with or without an error.
// Any error returned is related to the opening and reading from the input file.
func DetectFile(path string) (*MIME, error) {
	f, err := os.Open(path)
	if err != nil {
		return root, err
	}
	defer f.Close()

	return DetectReader(f)
}

// EqualsAny reports whether s MIME type is equal to any MIME type in mimes.
// MIME type equality test is done on the "type/subtype" section, ignores
// any optional MIME parameters, ignores any leading and trailing whitespace,
// and is case insensitive.
func EqualsAny(s string, mimes ...string) bool {
	s, _, _ = mime.ParseMediaType(s)
	for _, m := range mimes {
		m, _, _ = mime.ParseMediaType(m)
		if s == m {
			return true
		}
	}

	return false
}

// SetLimit sets the maximum number of bytes read from input when detecting the MIME type.
// Increasing the limit provides better detection for file formats which store
// their magical numbers towards the end of the file: docx, pptx, xlsx, etc.
// A limit of 0 means the whole input file will be used.
func SetLimit(limit uint32) {
	rootMu.Lock()
	readLimit = limit
	rootMu.Unlock()
}

// Extend adds detection for other file formats. The detector is a function
// returning true when the raw input file satisfies a  signature.
// The extension should include the leading dot, as in ".html".
func Extend(detector func(raw []byte, limit uint32) bool, mime, extension string, aliases ...string) {
	m := &MIME{
		mime:      mime,
		extension: extension,
		detector:  detector,
		parent:    root,
		aliases:   aliases,
	}

	rootMu.Lock()
	root.children = append([]*MIME{m}, root.children...)
	rootMu.Unlock()
}
