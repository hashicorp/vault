// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package trace

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/trace"
	"time"
)

func StartDebugTrace(dir string, filePrefix string) (file string, stop func() error, err error) {
	if dir == "" {
		// avoid permission concerns of using /tmp as a default dir
		return "", nil, fmt.Errorf("trace directory is required")
	}

	d, err := os.Stat(dir)
	if err != nil {
		// also fails if dir doesn't already exist
		return "", nil, fmt.Errorf("failed to stat trace directory %q: %s", dir, err)
	}

	if !d.IsDir() {
		return "", nil, fmt.Errorf("trace directory %q is not a directory", dir)
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
