// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package compressutil

import (
	"errors"
	"io/fs"
	"os"
	"testing"
)

var (
	input  []byte
	fname1 = "sdk/helper/compressutil/testdata/table.json"
	fname2 = "testdata/table.json"
	skip   = true
)

func init() {
	var err error
	if _, err = os.Stat(fname1); !errors.Is(err, fs.ErrNotExist) {
		input, err = os.ReadFile(fname1)
		if err != nil {
			return
		}
	} else if _, err = os.Stat(fname2); !errors.Is(err, fs.ErrNotExist) {
		input, err = os.ReadFile(fname2)
		if err != nil {
			return
		}
	}
	// found the input data
	skip = false
}

func Benchmark_Compress_MountTable_gzip_9(b *testing.B) {
	if skip {
		b.Skip("Could not find sample data")
	}
	benchCompress(b, &CompressionConfig{
		Type:                 CompressionTypeGzip,
		GzipCompressionLevel: 9,
	})
}

func Benchmark_Decompress_MountTable_gzip_9(b *testing.B) {
	if skip {
		b.Skip("Could not find sample data")
	}
	benchDecompress(b, &CompressionConfig{
		Type:                 CompressionTypeGzip,
		GzipCompressionLevel: 9,
	})
}

func Benchmark_Compress_MountTable(b *testing.B) {
	if skip {
		b.Skip("Could not find sample data")
	}
	types := []string{
		CompressionTypeGzip,
		CompressionTypeLZ4,
		CompressionTypeLZW,
		CompressionTypeSnappy,
		CompressionTypeZstd,
	}
	for _, t := range types {
		b.Run(t, func(b *testing.B) {
			benchCompress(b, &CompressionConfig{
				Type: t,
			})
		})
	}
}

func Benchmark_Decompress_MountTable(b *testing.B) {
	if skip {
		b.Skip("Could not find sample data")
	}
	types := []string{
		CompressionTypeGzip,
		CompressionTypeLZ4,
		CompressionTypeLZW,
		CompressionTypeSnappy,
		CompressionTypeZstd,
	}
	for _, t := range types {
		b.Run(t, func(b *testing.B) {
			benchDecompress(b, &CompressionConfig{
				Type: t,
			})
		})
	}
}

func benchCompress(b *testing.B, config *CompressionConfig) {
	b.SetBytes(int64(len(input)))
	for i := 0; i < b.N; i++ {
		compressed, err := Compress(input, config)
		if err != nil {
			b.Fatal(err)
		}
		_, _, err = Decompress(compressed)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchDecompress(b *testing.B, config *CompressionConfig) {
	b.StopTimer()
	b.SetBytes(int64(len(input)))
	compressed, err := Compress(input, config)
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _, err = Decompress(compressed)
		if err != nil {
			b.Fatal(err)
		}
	}
}
