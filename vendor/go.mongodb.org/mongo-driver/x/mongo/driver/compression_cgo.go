// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// +build cgo

package driver

import (
	"bytes"
	"io"

	"github.com/DataDog/zstd"
)

func zstdCompress(in []byte, level int) ([]byte, error) {
	var b bytes.Buffer
	w := zstd.NewWriterLevel(&b, level)
	_, err := w.Write(in)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func zstdDecompress(in []byte, size int32) ([]byte, error) {
	out := make([]byte, size)
	decompressor := zstd.NewReader(bytes.NewReader(in))
	_, err := io.ReadFull(decompressor, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
