// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build !appengine

package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// jsonLogger writes logs in the JSON format required for Flex Logging. It can
// be used concurrently.
type jsonLogger struct {
	mu  sync.Mutex
	enc *json.Encoder
}

type logLine struct {
	Message   string       `json:"message"`
	Timestamp logTimestamp `json:"timestamp"`
	Severity  string       `json:"severity"`
	TraceID   string       `json:"traceId,omitempty"`
}

type logTimestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int   `json:"nanos"`
}

var (
	logger     *jsonLogger
	loggerOnce sync.Once

	logPath      = "/var/log/app_engine/app.json"
	stderrLogger = newJSONLogger(os.Stderr)
	testLogger   = newJSONLogger(ioutil.Discard)

	levels = map[int64]string{
		0: "DEBUG",
		1: "INFO",
		2: "WARNING",
		3: "ERROR",
		4: "CRITICAL",
	}
)

func globalLogger() *jsonLogger {
	loggerOnce.Do(func() {
		f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("failed to open/create log file, logging to stderr: %v", err)
			logger = stderrLogger
			return
		}

		logger = newJSONLogger(f)
	})

	return logger
}

func logf(ctx *context, level int64, format string, args ...interface{}) {
	s := strings.TrimSpace(fmt.Sprintf(format, args...))
	now := time.Now()

	trace := ctx.req.Header.Get(traceHeader)
	if i := strings.Index(trace, "/"); i > -1 {
		trace = trace[:i]
	}

	line := &logLine{
		Message: s,
		Timestamp: logTimestamp{
			Seconds: now.Unix(),
			Nanos:   now.Nanosecond(),
		},
		Severity: levels[level],
		TraceID:  trace,
	}

	if err := ctx.logger.emit(line); err != nil {
		log.Printf("failed to write log line to file: %v", err)
	}

	log.Print(levels[level] + ": " + s)
}

func newJSONLogger(w io.Writer) *jsonLogger {
	return &jsonLogger{
		enc: json.NewEncoder(w),
	}
}

func (l *jsonLogger) emit(line *logLine) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.enc.Encode(line)
}
