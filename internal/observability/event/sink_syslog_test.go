// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"bytes"
	"errors"
	"syscall"
	"testing"

	gsyslog "github.com/hashicorp/go-syslog"
	"github.com/stretchr/testify/require"
)

// mockSyslogger records each Write call for inspection in tests.
type mockSyslogger struct {
	writes [][]byte
	err    error
}

func (m *mockSyslogger) Write(b []byte) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	cp := make([]byte, len(b))
	copy(cp, b)
	m.writes = append(m.writes, cp)
	return len(b), nil
}

func (m *mockSyslogger) WriteLevel(_ gsyslog.Priority, b []byte) error {
	_, err := m.Write(b)
	return err
}

func (m *mockSyslogger) Close() error { return nil }

// TestNewSyslogSink ensures that we validate the input arguments and can create
// the SyslogSink if everything goes to plan.
func TestNewSyslogSink(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		format         string
		opts           []Option
		want           *SyslogSink
		wantErr        bool
		expectedErrMsg string
	}{
		"format-empty": {
			format:         "",
			wantErr:        true,
			expectedErrMsg: "format is required: invalid parameter",
		},
		"format-whitespace": {
			format:         "   ",
			wantErr:        true,
			expectedErrMsg: "format is required: invalid parameter",
		},
		"happy": {
			format: "json",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := NewSyslogSink(tc.format, tc.opts...)

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedErrMsg)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
			}
		})
	}
}

// TestSyslogSink_writeChunked tests payload splitting and single-write paths.
func TestSyslogSink_writeChunked(t *testing.T) {
	t.Parallel()

	// Use a fixed chunkSize for tests so they are independent of tag length.
	const testChunkSize = 2010 // syslogSocketLimit - syslogFramingOverheadFixed - len("vault")

	tests := map[string]struct {
		payload       []byte
		syslogErr     error
		wantChunks    int
		wantErr       bool
		wantErrString string
	}{
		"small-payload-single-write": {
			payload:    bytes.Repeat([]byte("x"), 100),
			wantChunks: 1,
		},
		"exact-limit-single-write": {
			payload:    bytes.Repeat([]byte("x"), testChunkSize),
			wantChunks: 1,
		},
		"oversized-splits-into-chunks": {
			// 73 KB payload — matches the real PKI ?help=1 response size.
			payload:    bytes.Repeat([]byte("x"), 73*1024),
			wantChunks: (73*1024 + testChunkSize - 1) / testChunkSize,
		},
		"write-error-propagates": {
			payload:       bytes.Repeat([]byte("x"), 100),
			syslogErr:     errors.New("write: message too long"),
			wantErr:       true,
			wantErrString: "write: message too long",
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mock := &mockSyslogger{err: tc.syslogErr}
			s := &SyslogSink{syslogger: mock, chunkSize: testChunkSize}

			err := s.writeChunked(tc.payload)

			if tc.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.wantErrString)
				return
			}

			require.NoError(t, err)
			require.Len(t, mock.writes, tc.wantChunks)

			// Reassemble chunks and confirm round-trip fidelity.
			var reassembled []byte
			for _, chunk := range mock.writes {
				reassembled = append(reassembled, chunk...)
			}
			require.Equal(t, tc.payload, reassembled)

			// Every chunk except the last must be exactly testChunkSize bytes.
			for i, chunk := range mock.writes[:len(mock.writes)-1] {
				require.Len(t, chunk, testChunkSize,
					"chunk %d should be exactly testChunkSize bytes", i)
			}
		})
	}
}

// TestNewSyslogSink_chunkSize verifies chunkSize accounts for tag length.
func TestNewSyslogSink_chunkSize(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		tag           string
		wantChunkSize int
	}{
		"default-tag": {
			tag:           "vault",
			wantChunkSize: syslogSocketLimit - syslogFramingOverheadFixed - len("vault"),
		},
		"long-tag": {
			tag:           "vault-enterprise-audit",
			wantChunkSize: syslogSocketLimit - syslogFramingOverheadFixed - len("vault-enterprise-audit"),
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := NewSyslogSink("json", WithTag(tc.tag))
			require.NoError(t, err)
			require.Equal(t, tc.wantChunkSize, got.chunkSize)
		})
	}
}

// TestSyslogSink_writeWithRetry verifies the ENOBUFS retry behaviour.
func TestSyslogSink_writeWithRetry(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		failTimes     int   // number of leading ENOBUFS failures before success
		permanentErr  error // non-ENOBUFS error to return on all writes
		wantErr       bool
		wantErrString string
	}{
		"succeeds-first-try": {
			failTimes: 0,
		},
		"succeeds-after-one-enobufs": {
			failTimes: 1,
		},
		"succeeds-after-max-retries": {
			failTimes: 5, // exactly maxRetries
		},
		"fails-after-exhausting-retries": {
			failTimes:     6, // maxRetries + 1
			wantErr:       true,
			wantErrString: syscall.ENOBUFS.Error(),
		},
		"non-enobufs-error-not-retried": {
			permanentErr:  errors.New("write: message too long"),
			wantErr:       true,
			wantErrString: "write: message too long",
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mock := &retryMockSyslogger{
				failTimes:    tc.failTimes,
				permanentErr: tc.permanentErr,
			}
			s := &SyslogSink{syslogger: mock, chunkSize: 2010}

			err := s.writeWithRetry([]byte("test payload"))

			if tc.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.wantErrString)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// retryMockSyslogger returns ENOBUFS for the first failTimes calls, then succeeds.
// If permanentErr is set, it is returned on every call instead.
type retryMockSyslogger struct {
	calls        int
	failTimes    int
	permanentErr error
}

func (r *retryMockSyslogger) Write(b []byte) (int, error) {
	r.calls++
	if r.permanentErr != nil {
		return 0, r.permanentErr
	}
	if r.calls <= r.failTimes {
		return 0, syscall.ENOBUFS
	}
	return len(b), nil
}

func (r *retryMockSyslogger) WriteLevel(_ gsyslog.Priority, b []byte) error {
	_, err := r.Write(b)
	return err
}

func (r *retryMockSyslogger) Close() error { return nil }
