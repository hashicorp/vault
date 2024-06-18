// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"crypto/sha256"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// testAuditBackend will create an audit.Backend (which expects to use the eventlogger).
// NOTE: this will create the backend, it does not care whether Enterprise only options are in place.
func testAuditBackend(t *testing.T, path string, config map[string]string) Backend {
	t.Helper()

	headersCfg := &HeadersConfig{
		headerSettings: make(map[string]*HeaderSettings),
		view:           nil,
	}

	view := &logical.InmemStorage{}
	se := &logical.StorageEntry{Key: "salt", Value: []byte("juan")}
	err := view.Put(context.Background(), se)
	require.NoError(t, err)

	cfg := &BackendConfig{
		SaltView: view,
		SaltConfig: &salt.Config{
			HMAC:     sha256.New,
			HMACType: "hmac-sha256",
		},
		Logger:    corehelpers.NewTestLogger(t),
		Config:    config,
		MountPath: path,
	}

	be, err := NewSyslogBackend(cfg, headersCfg)
	require.NoError(t, err)
	require.NotNil(t, be)

	return be
}

// TestAuditBroker_Deregister_Multiple ensures that we can call deregister multiple
// times without issue if is no matching backend registered.
func TestAuditBroker_Deregister_Multiple(t *testing.T) {
	t.Parallel()

	l := corehelpers.NewTestLogger(t)
	a, err := NewBroker(l)
	require.NoError(t, err)
	require.NotNil(t, a)

	err = a.Deregister(context.Background(), "foo")
	require.NoError(t, err)

	err = a.Deregister(context.Background(), "foo2")
	require.NoError(t, err)
}

// TestAuditBroker_Register_MultipleFails checks for failure when we try to
// re-register an audit backend.
func TestAuditBroker_Register_MultipleFails(t *testing.T) {
	t.Parallel()

	l := corehelpers.NewTestLogger(t)
	a, err := NewBroker(l)
	require.NoError(t, err)
	require.NotNil(t, a)

	path := "b2-no-filter"
	noFilterBackend := testAuditBackend(t, path, map[string]string{})

	err = a.Register(noFilterBackend, false)
	require.NoError(t, err)

	err = a.Register(noFilterBackend, false)
	require.Error(t, err)
	require.EqualError(t, err, "backend already registered 'b2-no-filter': invalid configuration")
}

// BenchmarkAuditBroker_File_Request_DevNull Attempts to register a single `file`
// audit device on the broker, which points at /dev/null.
// It will then attempt to benchmark how long it takes Vault to complete logging
// a request, this really only shows us how Vault can handle lots of calls to the
// broker to trigger the eventlogger pipelines that audit devices are configured as.
// Since we aren't writing anything to file or doing any I/O.
// This test used to live in the file package for the file backend, but once the
// move to eventlogger was complete, there wasn't a way to create a file backend
// and manually just write to the underlying file itself, the old code used to do
// formatting and writing all together, but we've split this up with eventlogger
// with different nodes in a pipeline (think 1 audit device:1 pipeline) each
// handling a responsibility, for example:
// filter nodes filter events, so you can select which ones make it to your audit log
// formatter nodes format the events (to JSON/JSONX and perform HMACing etc)
// sink nodes handle sending the formatted data to a file, syslog or socket.
func BenchmarkAuditBroker_File_Request_DevNull(b *testing.B) {
	backendConfig := &BackendConfig{
		Config: map[string]string{
			"path": "/dev/null",
		},
		MountPath:  "test",
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
		Logger:     hclog.NewNullLogger(),
	}

	sink, err := NewFileBackend(backendConfig, nil)
	require.NoError(b, err)

	broker, err := NewBroker(nil)
	require.NoError(b, err)

	err = broker.Register(sink, false)
	require.NoError(b, err)

	in := &logical.LogInput{
		Auth: &logical.Auth{
			ClientToken:     "foo",
			Accessor:        "bar",
			EntityID:        "foobarentity",
			DisplayName:     "testtoken",
			NoDefaultPolicy: true,
			Policies:        []string{"root"},
			TokenType:       logical.TokenTypeService,
		},
		Request: &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "/foo",
			Connection: &logical.Connection{
				RemoteAddr: "127.0.0.1",
			},
			WrapInfo: &logical.RequestWrapInfo{
				TTL: 60 * time.Second,
			},
			Headers: map[string][]string{
				"foo": {"bar"},
			},
		},
	}

	ctx := namespace.RootContext(context.Background())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := broker.LogRequest(ctx, in); err != nil {
				panic(err)
			}
		}
	})
}

// TestBroker_isContextViable_basics checks the expected result of isContextViable
// for basic inputs such as nil, cancelled context and a never-ending context.
func TestBroker_isContextViable_basics(t *testing.T) {
	t.Parallel()

	require.False(t, isContextViable(nil))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.False(t, isContextViable(ctx))
	require.True(t, isContextViable(context.Background()))
}

// TestBroker_isContextViable_timeouts checks the expected result of isContextViable
// for various timeout durations.
func TestBroker_isContextViable_timeouts(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		timeout  time.Duration
		expected bool
	}{
		"2s-smaller-deadline": {
			timeout:  timeout - 2*time.Second,
			expected: false,
		},
		"same-deadline": {
			timeout:  timeout,
			expected: false, // Expected as a near miss
		},
		"same-deadline-plus": {
			timeout:  timeout + 100*time.Millisecond,
			expected: true,
		},
		"2x-longer-deadline": {
			timeout:  timeout * 2,
			expected: true,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			t.Cleanup(func() { cancel() })
			require.Equal(t, tc.expected, isContextViable(ctx))
		})
	}
}
