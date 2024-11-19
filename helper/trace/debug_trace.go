// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package trace

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime/trace"
	"time"
)

func StartDebugTrace(dir string, filePrefix string) (file string, stop func() error, err error) {
	if dir == "" {
		dir = path.Join(os.TempDir(), "vault-traces")
	}

	d, err := os.Stat(dir)
	if err != nil && !os.IsNotExist(err) {
		return "", nil, fmt.Errorf("failed to stat trace directory %q: %s", dir, err)
	}

	if !os.IsNotExist(err) && !d.IsDir() {
		return "", nil, fmt.Errorf("trace directory %q is not a directory", dir)
	}

	if os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0o755); err != nil {
			return "", nil, fmt.Errorf("failed to create trace directory %q: %s", dir, err)
		}
	}

	fileName := fmt.Sprintf("%s-%s.trace", filePrefix, time.Now().Format(time.RFC3339))
	traceFile, err := filepath.Abs(filepath.Join(dir, fileName))
	if err != nil {
		return "", nil, fmt.Errorf("failed to get absolute path for trace file: %s", err)
	}
	f, err := os.Create(traceFile)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create trace file %q: %s", traceFile, err)
	}

	if err := trace.Start(f); err != nil {
		f.Close()
		return "", nil, fmt.Errorf("failed to start trace: %s", err)
	}

	stop = func() error {
		trace.Stop()
		return f.Close()
	}

	return f.Name(), stop, nil
}
