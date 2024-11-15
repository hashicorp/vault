// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"sync"

	"github.com/golang/snappy"
	"github.com/klauspost/compress/zstd"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// CompressionOpts holds settings for how to compress a payload
type CompressionOpts struct {
	Compressor       wiremessage.CompressorID
	ZlibLevel        int
	ZstdLevel        int
	UncompressedSize int32
}

// mustZstdNewWriter creates a zstd.Encoder with the given level and a nil
// destination writer. It panics on any errors and should only be used at
// package initialization time.
func mustZstdNewWriter(lvl zstd.EncoderLevel) *zstd.Encoder {
	enc, err := zstd.NewWriter(
		nil,
		zstd.WithWindowSize(8<<20), // Set window size to 8MB.
		zstd.WithEncoderLevel(lvl),
	)
	if err != nil {
		panic(err)
	}
	return enc
}

var zstdEncoders = [zstd.SpeedBestCompression + 1]*zstd.Encoder{
	0:                           nil, // zstd.speedNotSet
	zstd.SpeedFastest:           mustZstdNewWriter(zstd.SpeedFastest),
	zstd.SpeedDefault:           mustZstdNewWriter(zstd.SpeedDefault),
	zstd.SpeedBetterCompression: mustZstdNewWriter(zstd.SpeedBetterCompression),
	zstd.SpeedBestCompression:   mustZstdNewWriter(zstd.SpeedBestCompression),
}

func getZstdEncoder(level zstd.EncoderLevel) (*zstd.Encoder, error) {
	if zstd.SpeedFastest <= level && level <= zstd.SpeedBestCompression {
		return zstdEncoders[level], nil
	}
	// The level is outside the expected range, return an error.
	return nil, fmt.Errorf("invalid zstd compression level: %d", level)
}

// zlibEncodersOffset is the offset into the zlibEncoders array for a given
// compression level.
const zlibEncodersOffset = -zlib.HuffmanOnly // HuffmanOnly == -2

var zlibEncoders [zlib.BestCompression + zlibEncodersOffset + 1]sync.Pool

func getZlibEncoder(level int) (*zlibEncoder, error) {
	if zlib.HuffmanOnly <= level && level <= zlib.BestCompression {
		if enc, _ := zlibEncoders[level+zlibEncodersOffset].Get().(*zlibEncoder); enc != nil {
			return enc, nil
		}
		writer, err := zlib.NewWriterLevel(nil, level)
		if err != nil {
			return nil, err
		}
		enc := &zlibEncoder{writer: writer, level: level}
		return enc, nil
	}
	// The level is outside the expected range, return an error.
	return nil, fmt.Errorf("invalid zlib compression level: %d", level)
}

func putZlibEncoder(enc *zlibEncoder) {
	if enc != nil {
		zlibEncoders[enc.level+zlibEncodersOffset].Put(enc)
	}
}

type zlibEncoder struct {
	writer *zlib.Writer
	buf    bytes.Buffer
	level  int
}

func (e *zlibEncoder) Encode(dst, src []byte) ([]byte, error) {
	defer putZlibEncoder(e)

	e.buf.Reset()
	e.writer.Reset(&e.buf)

	_, err := e.writer.Write(src)
	if err != nil {
		return nil, err
	}
	err = e.writer.Close()
	if err != nil {
		return nil, err
	}
	dst = append(dst[:0], e.buf.Bytes()...)
	return dst, nil
}

var zstdBufPool = sync.Pool{
	New: func() interface{} {
		s := make([]byte, 0)
		return &s
	},
}

// CompressPayload takes a byte slice and compresses it according to the options passed
func CompressPayload(in []byte, opts CompressionOpts) ([]byte, error) {
	switch opts.Compressor {
	case wiremessage.CompressorNoOp:
		return in, nil
	case wiremessage.CompressorSnappy:
		return snappy.Encode(nil, in), nil
	case wiremessage.CompressorZLib:
		encoder, err := getZlibEncoder(opts.ZlibLevel)
		if err != nil {
			return nil, err
		}
		return encoder.Encode(nil, in)
	case wiremessage.CompressorZstd:
		encoder, err := getZstdEncoder(zstd.EncoderLevelFromZstd(opts.ZstdLevel))
		if err != nil {
			return nil, err
		}
		ptr := zstdBufPool.Get().(*[]byte)
		b := encoder.EncodeAll(in, *ptr)
		dst := make([]byte, len(b))
		copy(dst, b)
		*ptr = b[:0]
		zstdBufPool.Put(ptr)
		return dst, nil
	default:
		return nil, fmt.Errorf("unknown compressor ID %v", opts.Compressor)
	}
}

var zstdReaderPool = sync.Pool{
	New: func() interface{} {
		r, _ := zstd.NewReader(nil)
		return r
	},
}

// DecompressPayload takes a byte slice that has been compressed and undoes it according to the options passed
func DecompressPayload(in []byte, opts CompressionOpts) ([]byte, error) {
	switch opts.Compressor {
	case wiremessage.CompressorNoOp:
		return in, nil
	case wiremessage.CompressorSnappy:
		l, err := snappy.DecodedLen(in)
		if err != nil {
			return nil, fmt.Errorf("decoding compressed length %w", err)
		} else if int32(l) != opts.UncompressedSize {
			return nil, fmt.Errorf("unexpected decompression size, expected %v but got %v", opts.UncompressedSize, l)
		}
		out := make([]byte, opts.UncompressedSize)
		return snappy.Decode(out, in)
	case wiremessage.CompressorZLib:
		r, err := zlib.NewReader(bytes.NewReader(in))
		if err != nil {
			return nil, err
		}
		out := make([]byte, opts.UncompressedSize)
		if _, err := io.ReadFull(r, out); err != nil {
			return nil, err
		}
		if err := r.Close(); err != nil {
			return nil, err
		}
		return out, nil
	case wiremessage.CompressorZstd:
		buf := make([]byte, 0, opts.UncompressedSize)
		// Using a pool here is about ~20% faster
		// than using a single global zstd.Reader
		r := zstdReaderPool.Get().(*zstd.Decoder)
		out, err := r.DecodeAll(in, buf)
		zstdReaderPool.Put(r)
		return out, err
	default:
		return nil, fmt.Errorf("unknown compressor ID %v", opts.Compressor)
	}
}
