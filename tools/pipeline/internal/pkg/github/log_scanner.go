// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// LogScanner it a Github Actions workflow job log scanner and decoder.
type LogScanner struct {
	Truncate         bool
	MaxSizeBytes     int
	OnlySteps        []string
	OnlyUnsuccessful bool
}

// LogScanOpts are scanner options.
type LogScanOpts func(*LogScanner)

// LogEntry is a job step log entry.
type LogEntry struct {
	StepName string `json:"step_name,omitempty"`
	SetupLog []byte `json:"setup_log,omitempty"`
	BodyLog  []byte `json:"body_log,omitempty"`
	ErrorLog []byte `json:"error_log,omitempty"`
}

// logSection is a section token in the job log.
type logSection int

const (
	labelGroup    = "##[group]"
	labelEndGroup = "##[endgroup]"
	labelError    = "##[error]"
)

const (
	logSectionNone logSection = iota
	logSectionSetup
	logSectionBody
	logSectionError
)

// NewLogScaner takes none-or-many LogScanOpts and returns a new instance of LogScanner.
func NewLogScaner(opts ...LogScanOpts) *LogScanner {
	scanner := &LogScanner{
		Truncate:         false,
		OnlySteps:        []string{},
		OnlyUnsuccessful: false,
		MaxSizeBytes:     (1 << 20), // 1MiB
	}

	for _, opt := range opts {
		opt(scanner)
	}

	return scanner
}

// WithLogScannerTruncate enables body log truncation.
func WithLogScannerTruncate() LogScanOpts {
	return func(scanner *LogScanner) {
		scanner.Truncate = true
	}
}

// WithLogScannerMaxSize configures the max body log size if truncation is enabled.
func WithLogScannerMaxSize(max int) LogScanOpts {
	return func(scanner *LogScanner) {
		scanner.MaxSizeBytes = max
	}
}

// WithLogScannerOnlySteps takes step names and only returns entries for matching steps.
func WithLogScannerOnlySteps(steps []string) LogScanOpts {
	return func(scanner *LogScanner) {
		if scanner.OnlySteps == nil {
			scanner.OnlySteps = []string{}
		}
		for _, g := range steps {
			scanner.OnlySteps = append(scanner.OnlySteps, strings.TrimSpace(g))
		}
	}
}

// WithLogScannerOnlyUnsuccessful filters any successful log entries.
func WithLogScannerOnlyUnsuccessful() LogScanOpts {
	return func(scanner *LogScanner) {
		scanner.OnlyUnsuccessful = true
	}
}

// Scan scans a Github Actions job raw log file and parses it into individual
// entries for each step that is run.
func (s *LogScanner) Scan(in io.Reader) ([]*LogEntry, error) {
	if s == nil {
		return nil, errors.New("uninitialized scanner")
	}

	scanner := bufio.NewScanner(in)
	logBuffer := newLogBuffer(s.Truncate, s.MaxSizeBytes)
	res := []*LogEntry{}
	logSec := logSectionNone

	for scanner.Scan() {
		// ##[group]
		if strings.Contains(scanner.Text(), labelGroup) {
			logSec = logSectionSetup
			// Start parsing a new log group.

			// Before we begin, persist our last log entry and reset our buffer for
			// our new group.
			if logBuffer.stepName != "" {
				res = append(res, logBuffer.entry())
				logBuffer.reset()
			}

			parts := strings.SplitN(scanner.Text(), labelGroup, 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("malformed log group line expected %s followed by step name, got: %s", labelGroup, scanner.Text())
			}
			logBuffer.stepName = strings.TrimSpace(parts[1])
			continue
		}

		// ##[error]
		if strings.Contains(scanner.Text(), labelError) {
			logSec = logSectionError
			// The error label often preceeds the actual first line of the error log.
			// Make sure we extract it and write it to the log.
			parts := strings.SplitN(scanner.Text(), labelError, 2)
			if len(parts) == 2 {
				logBuffer.write(logSec, parts[0]+parts[1]+"\n")
			}
			continue
		}

		// ##[endgroup]
		if strings.Contains(scanner.Text(), labelEndGroup) {
			logSec = logSectionBody
			continue
		}

		// Write the line to the buffer
		logBuffer.write(logSec, scanner.Text()+"\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Write our last entry
	res = append(res, logBuffer.entry())

	if len(s.OnlySteps) == 0 && !s.OnlyUnsuccessful {
		return res, nil
	}

	// Filter our respons if necessary
	var entries []*LogEntry
	if s.OnlyUnsuccessful {
		for _, entry := range res {
			if len(entry.ErrorLog) > 0 {
				entries = append(entries, entry)
			}
		}
	} else {
		entries = res
	}

	if len(s.OnlySteps) < 1 {
		return entries, nil
	}

	var filtered []*LogEntry
	for _, entry := range entries {
		for _, os := range s.OnlySteps {
			if strings.Contains(strings.TrimSpace(entry.StepName), os) {
				filtered = append(filtered, entry)
			}
		}
	}

	return filtered, nil
}

// logBuffer is a buffer used by the scanner when parsing a raw job log.
type logBuffer struct {
	truncate bool
	maxSize  int
	stepName string
	setup    *strings.Builder
	body     *strings.Builder
	error    *strings.Builder
}

func (b *logBuffer) write(lc logSection, in string) {
	switch lc {
	case logSectionNone:
	case logSectionSetup:
		b.setup.WriteString(in)
	case logSectionBody:
		b.body.WriteString(in)
	case logSectionError:
		b.error.WriteString(in)
	}
}

func (b *logBuffer) entry() *LogEntry {
	var groupBody string
	if b.truncate && b.body.Len() > b.maxSize {
		groupBody = b.body.String()[b.body.Len()-b.maxSize:]
	} else {
		groupBody = b.body.String()
	}

	e := &LogEntry{
		StepName: b.stepName,
	}
	if sl := []byte(strings.TrimSpace(b.setup.String())); len(sl) > 0 {
		e.SetupLog = sl
	}
	if bl := []byte(strings.TrimSpace(groupBody)); len(bl) > 0 {
		e.BodyLog = bl
	}

	if el := []byte(strings.TrimSpace(b.error.String())); len(el) > 0 {
		e.ErrorLog = el
	}

	return e
}

func (b *logBuffer) reset() {
	b.setup.Reset()
	b.body.Reset()
	b.error.Reset()
}

func newLogBuffer(truncate bool, maxSize int) *logBuffer {
	return &logBuffer{
		truncate: truncate,
		maxSize:  maxSize,
		body:     &strings.Builder{},
		setup:    &strings.Builder{},
		error:    &strings.Builder{},
	}
}
