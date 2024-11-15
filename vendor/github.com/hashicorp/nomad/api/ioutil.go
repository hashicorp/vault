// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"hash"
	"io"
	"strings"
)

var errMismatchChecksum = errors.New("mismatch checksum")

// checksumValidatingReader is a wrapper reader that validates
// the checksum of the underlying reader.
type checksumValidatingReader struct {
	r io.ReadCloser

	// algo is the hash algorithm (e.g. `sha-256`)
	algo string

	// checksum is the base64 component of checksum
	checksum string

	// hash is the hashing function used to compute the checksum
	hash hash.Hash
}

// newChecksumValidatingReader returns a checksum-validating wrapper reader, according
// to a digest received in HTTP header
//
// The digest must be in the format "<algo>=<base64 of hash>" (e.g. "sha-256=gPelGB7...").
//
// When the reader is fully consumed (i.e. EOT is encountered), if the checksum don't match,
// `Read` returns a checksum mismatch error.
func newChecksumValidatingReader(r io.ReadCloser, digest string) (io.ReadCloser, error) {
	parts := strings.SplitN(digest, "=", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid digest format")
	}

	algo := parts[0]
	var hash hash.Hash
	switch algo {
	case "sha-256":
		hash = sha256.New()
	case "sha-512":
		hash = sha512.New()
	case "md5":
		hash = md5.New()
	}

	return &checksumValidatingReader{
		r:        r,
		algo:     algo,
		checksum: parts[1],
		hash:     hash,
	}, nil
}

func (r *checksumValidatingReader) Read(b []byte) (int, error) {
	n, err := r.r.Read(b)
	if n != 0 {
		r.hash.Write(b[:n])
	}

	if err == io.EOF || err == io.ErrClosedPipe {
		found := base64.StdEncoding.EncodeToString(r.hash.Sum(nil))
		if found != r.checksum {
			return n, errMismatchChecksum
		}
	}

	return n, err
}

func (r *checksumValidatingReader) Close() error {
	return r.r.Close()
}
