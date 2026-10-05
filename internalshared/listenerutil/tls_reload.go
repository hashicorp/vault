// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package listenerutil

import (
	"context"
	"crypto/sha256"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/reloadutil"
)

// tlsReloadFailureEscalation is the number of consecutive failures after which
// reload errors are logged at ERROR rather than WARN.
//
// A rotation writes the certificate and its key as separate operations, so a
// brief mismatch is normal and should not alarm operators. Failures that
// persist past this many polls indicate a pair that will not become usable
// without intervention, such as a rotation that wrote a certificate but never
// its key. The same value throttles repeat ERROR logs so that a permanently
// broken pair does not flood the log.
const tlsReloadFailureEscalation = 5

// certFingerprint is the hash of a certificate/key pair as read from disk. It is
// used to detect content changes without retaining the key material itself.
type certFingerprint struct {
	cert [sha256.Size]byte
	key  [sha256.Size]byte
	// valid reports whether both files were read successfully. A zero-valued
	// fingerprint must never compare equal to a successfully read one.
	valid bool
}

// fingerprintFiles reads the certificate and key files and returns a hash of
// their contents. An error is returned if either file cannot be read, which
// callers are expected to treat as a transient condition.
func fingerprintFiles(certFile, keyFile string) (certFingerprint, error) {
	certPEM, err := os.ReadFile(certFile)
	if err != nil {
		return certFingerprint{}, err
	}

	keyPEM, err := os.ReadFile(keyFile)
	if err != nil {
		return certFingerprint{}, err
	}

	return certFingerprint{
		cert:  sha256.Sum256(certPEM),
		key:   sha256.Sum256(keyPEM),
		valid: true,
	}, nil
}

func (f certFingerprint) equal(other certFingerprint) bool {
	return f.valid && other.valid && f.cert == other.cert && f.key == other.key
}

// TLSReloader polls a certificate and key file on disk and triggers a reload
// when their contents change.
//
// Polling is used rather than filesystem notifications because certificates are
// frequently delivered through indirections that do not produce events on the
// watched paths. On Kubernetes, for example, a mounted Secret is updated by
// atomically swapping a "..data" symlink, so a watch registered against the leaf
// path never fires. Comparing content hashes is agnostic to how the file was
// replaced, and behaves identically on virtual machines and bare metal.
type TLSReloader struct {
	certFile string
	keyFile  string
	interval time.Duration
	reload   reloadutil.ReloadFunc
	logger   hclog.Logger

	// last is the fingerprint of the most recently observed pair. It is only
	// accessed from the polling goroutine.
	last certFingerprint

	// consecutiveFailures counts polls that have failed in a row, and is reset
	// on the next success. It is only accessed from the polling goroutine.
	consecutiveFailures int
}

// NewTLSReloader constructs a reloader for the given certificate pair. The
// reload function is expected to re-read both files and atomically install the
// resulting certificate, such as reloadutil.CertificateGetter.Reload.
func NewTLSReloader(certFile, keyFile string, interval time.Duration, reload reloadutil.ReloadFunc, logger hclog.Logger) *TLSReloader {
	return &TLSReloader{
		certFile: certFile,
		keyFile:  keyFile,
		interval: interval,
		reload:   reload,
		logger:   logger,
	}
}

// Run polls until the context is cancelled. It is intended to be started in its
// own goroutine.
//
// The zero-valued fingerprint is deliberately left in place rather than seeded
// from disk. The listener's certificate was loaded during listener
// initialization, which can precede this goroutine by a significant interval
// while the rest of the server starts up. Seeding here would record whatever is
// on disk at this later moment, so a rotation landing in that gap would be
// treated as already applied and the listener would keep serving the previous
// certificate indefinitely. Starting with an invalid fingerprint forces the
// first poll to reconcile the loaded certificate with the files on disk, at the
// cost of one redundant reload.
//
// Failures are never fatal. A certificate and its key are separate files and may
// be observed mid-update, and a rotation may briefly leave a cert that does not
// match its key. In either case the error is logged and the current certificate
// is left in place, so the listener continues serving the last known good pair
// until a complete one appears.
func (r *TLSReloader) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	r.logger.Debug("started TLS certificate reloader",
		"cert_file", r.certFile, "key_file", r.keyFile, "interval", r.interval)

	for {
		select {
		case <-ctx.Done():
			r.logger.Debug("stopping TLS certificate reloader", "cert_file", r.certFile)
			return
		case <-ticker.C:
			r.poll()
		}
	}
}

// poll checks for a change and reloads if one is found.
func (r *TLSReloader) poll() {
	current, err := fingerprintFiles(r.certFile, r.keyFile)
	if err != nil {
		// Most likely a partially completed update; the next tick will retry.
		r.recordFailure("unable to read TLS certificate, retaining existing certificate", err)
		return
	}

	if current.equal(r.last) {
		return
	}

	if err := r.reload(); err != nil {
		// The pair on disk is inconsistent, such as a new certificate written
		// before its key. Leave the fingerprint unchanged so this is retried.
		r.recordFailure("failed to reload TLS certificate, retaining existing certificate", err)
		return
	}

	if r.consecutiveFailures >= tlsReloadFailureEscalation {
		r.logger.Info("TLS certificate reload recovered after repeated failures",
			"cert_file", r.certFile, "failed_attempts", r.consecutiveFailures)
	}
	r.consecutiveFailures = 0

	r.last = current
	r.logger.Info("reloaded TLS certificate", "cert_file", r.certFile, "key_file", r.keyFile)
}

// recordFailure logs a failed poll, escalating from WARN to ERROR once failures
// have persisted long enough that they are unlikely to resolve on their own.
//
// Vault keeps serving the last known good certificate either way, so this only
// affects how the condition is reported. Escalating lets operators distinguish a
// transient mid-rotation mismatch from a rotation that never completed, and
// throttling the escalated logs keeps a permanently broken pair from flooding
// the log on every tick.
func (r *TLSReloader) recordFailure(msg string, err error) {
	r.consecutiveFailures++

	if r.consecutiveFailures < tlsReloadFailureEscalation {
		r.logger.Warn(msg, "cert_file", r.certFile, "error", err)
		return
	}

	if r.consecutiveFailures%tlsReloadFailureEscalation == 0 {
		r.logger.Error(msg+"; certificate rotation appears incomplete and needs attention",
			"cert_file", r.certFile, "key_file", r.keyFile,
			"consecutive_failures", r.consecutiveFailures, "error", err)
	}
}
