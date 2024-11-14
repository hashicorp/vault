package trace

import (
	"fmt"
	"os"
	"runtime/trace"
	"time"
)

func StartDebugTrace(filePrefix string) (file string, stop func() error, err error) {
	path := fmt.Sprintf("%s/%s_%s", os.TempDir(), filePrefix, time.Now().Format(time.RFC3339))
	traceFile, err := os.Create(path)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create trace file: %s", err)
	}

	if err := trace.Start(traceFile); err != nil {
		traceFile.Close()
		return "", nil, fmt.Errorf("failed to start trace: %s", err)
	}

	stop = func() error {
		trace.Stop()
		return traceFile.Close()
	}

	return traceFile.Name(), stop, nil
}
