// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// +build !cgo

package driver

import "errors"

func zstdCompress(_ []byte, _ int) ([]byte, error) {
	return nil, errors.New("zstd support requires cgo")
}

func zstdDecompress(in []byte, size int32) ([]byte, error) {
	return nil, errors.New("zstd support requires cgo")
}
