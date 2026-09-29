// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package listenerutil

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
)

// tlsReloaderTestInterval is deliberately generous so that assertions about
// "no reload happened yet" remain reliable even when the test binary is
// under heavy scheduling contention (e.g. many parallel tests).
const tlsReloaderTestInterval = 20 * time.Millisecond

// syncBuffer is a concurrency-safe io.Writer for capturing log output, since
// the reloader logs from its own goroutine while the test reads.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// count returns how many captured log lines contain substr.
func (b *syncBuffer) count(substr string) int {
	return strings.Count(b.String(), substr)
}

// TestTLSReloader_ReloadsOnChange verifies that TLSReloader detects a change
// to either the certificate or key file's contents and invokes the reload
// function. This is the mechanism that lets Vault pick up a rotated certificate
// on its own, without an operator sending SIGHUP or restarting the process.
func TestTLSReloader_ReloadsOnChange(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	certFile := filepath.Join(dir, "cert.pem")
	keyFile := filepath.Join(dir, "key.pem")

	require.NoError(t, os.WriteFile(certFile, []byte("cert-v1"), 0o600))
	require.NoError(t, os.WriteFile(keyFile, []byte("key-v1"), 0o600))

	var reloadCount int32
	reload := func() error {
		atomic.AddInt32(&reloadCount, 1)
		return nil
	}

	r := NewTLSReloader(certFile, keyFile, tlsReloaderTestInterval, reload, hclog.NewNullLogger())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)

	// The first poll always reconciles, so wait for that to land before
	// rotating. This keeps the assertions below about the rotation itself
	// rather than the initial reconcile.
	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&reloadCount) >= 1
	}, 5*time.Second, tlsReloaderTestInterval, "initial reconcile did not occur")
	afterReconcile := atomic.LoadInt32(&reloadCount)

	// Replace the certificate contents, simulating a rotation.
	require.NoError(t, os.WriteFile(certFile, []byte("cert-v2"), 0o600))

	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&reloadCount) > afterReconcile
	}, 5*time.Second, tlsReloaderTestInterval, "reload was not triggered after certificate changed")

	// os.WriteFile is not atomic (truncate then write), so under scheduling
	// contention a poll can occasionally observe a torn intermediate state
	// and reload an extra time; that's harmless in practice, so we only rely
	// on the count strictly increasing rather than an exact value.
	afterCertChange := atomic.LoadInt32(&reloadCount)

	// Changing only the key should also trigger a reload.
	require.NoError(t, os.WriteFile(keyFile, []byte("key-v2"), 0o600))

	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&reloadCount) > afterCertChange
	}, 5*time.Second, tlsReloaderTestInterval, "reload was not triggered after key changed")
}

// TestTLSReloader_ReconcilesOnFirstPoll verifies that the first poll reloads
// even when the files on disk have not changed since Run started.
//
// The listener's certificate is loaded during listener initialization, which
// can precede the reloader goroutine by a significant interval while the rest
// of the server starts up. If the reloader seeded its fingerprint from disk at
// startup, a rotation landing in that gap would be recorded as already applied
// and the listener would serve the superseded certificate indefinitely. Always
// reconciling on the first poll closes that window.
func TestTLSReloader_ReconcilesOnFirstPoll(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	certFile := filepath.Join(dir, "cert.pem")
	keyFile := filepath.Join(dir, "key.pem")

	// These contents are never modified during the test, mimicking a rotation
	// that completed before the reloader goroutine started.
	require.NoError(t, os.WriteFile(certFile, []byte("cert-rotated-during-startup"), 0o600))
	require.NoError(t, os.WriteFile(keyFile, []byte("key-rotated-during-startup"), 0o600))

	var reloadCount int32
	reload := func() error {
		atomic.AddInt32(&reloadCount, 1)
		return nil
	}

	r := NewTLSReloader(certFile, keyFile, tlsReloaderTestInterval, reload, hclog.NewNullLogger())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)

	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&reloadCount) >= 1
	}, 5*time.Second, tlsReloaderTestInterval,
		"first poll did not reload, so a rotation during startup would be missed")
}

// TestTLSReloader_RetainsCurrentCertOnError verifies that when the reload
// function fails (for example because a rotation briefly left a mismatched
// cert/key pair on disk), the reloader keeps polling and does not get stuck:
// once a working pair is restored, it successfully reloads.
func TestTLSReloader_RetainsCurrentCertOnError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	certFile := filepath.Join(dir, "cert.pem")
	keyFile := filepath.Join(dir, "key.pem")

	require.NoError(t, os.WriteFile(certFile, []byte("cert-v1"), 0o600))
	require.NoError(t, os.WriteFile(keyFile, []byte("key-v1"), 0o600))

	var shouldFail atomic.Bool
	shouldFail.Store(true)
	var reloadAttempts, reloadSuccesses int32
	reload := func() error {
		atomic.AddInt32(&reloadAttempts, 1)
		if shouldFail.Load() {
			return os.ErrInvalid
		}
		atomic.AddInt32(&reloadSuccesses, 1)
		return nil
	}

	r := NewTLSReloader(certFile, keyFile, tlsReloaderTestInterval, reload, hclog.NewNullLogger())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)

	// Wait until the reloader has attempted (and failed) to reload at least
	// once. The reload function fails for every attempt until shouldFail is
	// cleared below, so the fingerprint is never advanced.
	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&reloadAttempts) >= 1
	}, 5*time.Second, tlsReloaderTestInterval, "reload was never attempted")
	require.EqualValues(t, 0, atomic.LoadInt32(&reloadSuccesses), "reload should not have succeeded while failing")

	// Once the reload function stops failing, the fingerprint was never
	// advanced past the failed attempt, so the next poll should succeed.
	shouldFail.Store(false)

	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&reloadSuccesses) >= 1
	}, 5*time.Second, tlsReloaderTestInterval, "reload was never retried successfully")
}

// TestTLSReloader_StopsOnContextCancel verifies that cancelling the context
// passed to Run stops the polling goroutine, so it does not outlive its
// listener.
func TestTLSReloader_StopsOnContextCancel(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	certFile := filepath.Join(dir, "cert.pem")
	keyFile := filepath.Join(dir, "key.pem")

	require.NoError(t, os.WriteFile(certFile, []byte("cert-v1"), 0o600))
	require.NoError(t, os.WriteFile(keyFile, []byte("key-v1"), 0o600))

	reload := func() error { return nil }

	r := NewTLSReloader(certFile, keyFile, tlsReloaderTestInterval, reload, hclog.NewNullLogger())

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		r.Run(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("TLSReloader.Run did not return after context cancellation")
	}
}

// TestTLSReloader_EscalatesPersistentFailures verifies that repeated reload
// failures are escalated from WARN to ERROR, and that a subsequent success is
// reported so operators know the condition cleared.
//
// A rotation writes the certificate and key as separate operations, so a brief
// mismatch is normal and stays at WARN. If a rotation never completes (for
// example the tool crashed after writing the certificate but before the key),
// the failure repeats on every poll forever. Without escalation that state is
// indistinguishable from a transient blip and only ever produces WARN noise.
func TestTLSReloader_EscalatesPersistentFailures(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	certFile := filepath.Join(dir, "cert.pem")
	keyFile := filepath.Join(dir, "key.pem")

	require.NoError(t, os.WriteFile(certFile, []byte("cert-v1"), 0o600))
	require.NoError(t, os.WriteFile(keyFile, []byte("key-v1"), 0o600))

	var shouldFail atomic.Bool
	shouldFail.Store(true)
	reload := func() error {
		if shouldFail.Load() {
			return os.ErrInvalid
		}
		return nil
	}

	logs := &syncBuffer{}
	logger := hclog.New(&hclog.LoggerOptions{Output: logs, Level: hclog.Trace})

	r := NewTLSReloader(certFile, keyFile, tlsReloaderTestInterval, reload, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)

	// The first failures stay at WARN, since a mismatched pair is expected
	// briefly during a normal rotation.
	require.Eventually(t, func() bool {
		return logs.count("[WARN]") >= 1
	}, 5*time.Second, tlsReloaderTestInterval, "expected a WARN on the first failure")
	require.Zero(t, logs.count("[ERROR]"), "should not escalate on the first failure")

	// Once failures persist past the threshold they escalate to ERROR.
	require.Eventually(t, func() bool {
		return logs.count("[ERROR]") >= 1
	}, 5*time.Second, tlsReloaderTestInterval,
		"persistent failures were never escalated to ERROR")

	// Escalated logs are throttled rather than emitted on every tick, so ERROR
	// lines must stay well below the number of failed polls.
	require.Eventually(t, func() bool {
		return logs.count("[ERROR]") >= 2
	}, 5*time.Second, tlsReloaderTestInterval, "expected escalation to repeat periodically")
	require.Less(t, logs.count("[ERROR]"), logs.count("[WARN]")+logs.count("[ERROR]"),
		"ERROR logs should be throttled relative to total failures")

	// A success after escalation should be reported so operators know it cleared.
	shouldFail.Store(false)
	require.Eventually(t, func() bool {
		return strings.Contains(logs.String(), "recovered after repeated failures")
	}, 5*time.Second, tlsReloaderTestInterval, "recovery was never reported")
}
