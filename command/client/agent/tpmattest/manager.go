// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package tpmattest

import (
	"context"
	"fmt"
	"math"
	mathrand "math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/backoff"
	sdktpm "github.com/hashicorp/vault/sdk/helper/tpm"
)

// AttestFunc performs TPM attestation and returns the result. The function is
// responsible for all Vault client interaction and TPM configuration; Manager
// treats it as an opaque operation.
type AttestFunc func(ctx context.Context) (*sdktpm.AttestationResult, error)

// ValidateFunc validates a loaded AttestationResult against the local TPM
// hardware (e.g. verifying that the persisted key blobs are still usable).
// A nil ValidateFunc skips hardware validation.
type ValidateFunc func(*sdktpm.AttestationResult) error

type Manager struct {
	attest       AttestFunc
	validate     ValidateFunc
	destPath     string
	umask        os.FileMode
	renewPercent int
	logger       hclog.Logger

	mu           sync.RWMutex
	timer        *time.Timer
	retryBackoff *backoff.Backoff

	stopCh  chan struct{}
	started bool
	stopped atomic.Bool
}

// NewManager creates a Manager that periodically calls attest to obtain and
// renew TPM certificate material. validate may be nil to skip hardware
// validation when loading persisted material.
func NewManager(attest AttestFunc, validate ValidateFunc, destPath string, umask os.FileMode, renewPercent int, logger hclog.Logger) (*Manager, error) {
	if attest == nil {
		return nil, fmt.Errorf("attest func cannot be nil")
	}

	return &Manager{
		attest:       attest,
		validate:     validate,
		destPath:     destPath,
		umask:        umask,
		renewPercent: renewPercent,
		logger:       logger,
		stopCh:       make(chan struct{}),
		retryBackoff: backoff.NewBackoff(math.MaxInt, 1*time.Minute, 1*time.Hour),
	}, nil
}

func (m *Manager) Start(ctx context.Context) {
	m.mu.Lock()
	if m.started {
		m.mu.Unlock()
		return
	}
	m.started = true
	m.mu.Unlock()

	m.logger.Info("starting TPM attestation manager")
	defer m.logger.Info("TPM attestation manager stopped")

	reused, err := m.loadPersistedMaterial()
	if err != nil {
		m.logger.Warn("failed to load persisted TPM certificate material, will reissue", "error", err)
	}

	if !reused {
		if err != nil {
			m.logger.Info("starting initial acquisition after load failure")
		} else {
			m.logger.Info("no reusable persisted TPM certificate material found, starting initial acquisition")
		}
		if err := m.acquireMaterial(ctx); err != nil {
			m.logger.Error("failed to acquire initial TPM certificate material", "error", err)
			m.scheduleRetryAfter()
		}
	} else {
		m.logger.Info("loaded persisted TPM certificate material")
	}

	for {
		m.mu.RLock()
		timerC := m.timerC()
		m.mu.RUnlock()

		select {
		case <-ctx.Done():
			return
		case <-m.stopCh:
			return
		case <-timerC:
			m.logger.Info("renewal timer triggered")
			if err := m.acquireMaterial(ctx); err != nil {
				m.logger.Error("TPM certificate renewal failed", "error", err)
				m.scheduleRetryAfter()
			}
		}
	}
}

func (m *Manager) Stop() {
	if !m.stopped.CompareAndSwap(false, true) {
		return
	}

	m.mu.Lock()
	if m.timer != nil {
		m.timer.Stop()
		m.timer = nil
	}
	m.mu.Unlock()

	close(m.stopCh)
}

func (m *Manager) acquireMaterial(ctx context.Context) error {
	m.logger.Info("acquiring TPM certificate material")

	result, err := m.attest(ctx)
	if err != nil {
		return fmt.Errorf("failed to attest with TPM: %w", err)
	}

	if err := result.SaveToDirectory(m.destPath, m.umask); err != nil {
		return fmt.Errorf("failed to save TPM material: %w", err)
	}

	m.scheduleRenewal(result)

	return nil
}

func (m *Manager) loadPersistedMaterial() (bool, error) {
	result, err := sdktpm.LoadFromDirectory(m.destPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to load persisted material: %w", err)
	}

	if time.Now().After(result.NotAfter) {
		return false, fmt.Errorf("persisted certificate has expired at %s", result.NotAfter)
	}

	if m.validate != nil {
		if err := m.validate(result); err != nil {
			return false, fmt.Errorf("persisted material failed validation: %w", err)
		}
	}

	m.scheduleRenewal(result)
	return true, nil
}

// parseUmask parses a umask string (e.g., "077") into an os.FileMode.
func parseUmask(umaskStr string) (os.FileMode, error) {
	if umaskStr == "" {
		return 0, nil
	}

	umaskVal, err := strconv.ParseUint(umaskStr, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid umask value %q: %w", umaskStr, err)
	}

	return os.FileMode(umaskVal), nil
}

// timerC returns the channel of the current renewal timer, or a nil channel
// (which blocks forever) when no timer is set. Must be called with m.mu held
// for reading.
func (m *Manager) timerC() <-chan time.Time {
	if m.timer == nil {
		return nil
	}
	return m.timer.C
}

// resetTimer arms the renewal timer to fire after duration. If a timer already
// exists its channel is reused (via Reset) so that any goroutine already
// selecting on that channel will wake up. Must be called with m.mu held for
// writing.
func (m *Manager) resetTimer(duration time.Duration) {
	if m.timer != nil {
		// Stop and drain before Reset per the time.Timer docs.
		if !m.timer.Stop() {
			select {
			case <-m.timer.C:
			default:
			}
		}
		m.timer.Reset(duration)
		return
	}
	m.timer = time.NewTimer(duration)
}

func (m *Manager) scheduleRenewal(result *sdktpm.AttestationResult) {
	lifetime := result.NotAfter.Sub(result.NotBefore)
	renewalThreshold := time.Duration(float64(lifetime) * float64(m.renewPercent) / 100.0)
	renewalTime := result.NotAfter.Add(-renewalThreshold)

	duration := time.Until(renewalTime)
	if duration < 0 {
		duration = time.Second
	} else {
		jitterRange := float64(duration) * 0.1
		jitter := time.Duration(mathrand.Float64()*jitterRange*2 - jitterRange)
		duration += jitter
		if duration < 0 {
			duration = time.Second
		}
	}

	m.mu.Lock()
	m.retryBackoff.Reset()
	m.resetTimer(duration)
	m.mu.Unlock()

	m.logger.Info("scheduled TPM certificate renewal", "renewal_time", renewalTime, "duration", duration)
}

func (m *Manager) scheduleRetryAfter() {
	m.mu.Lock()
	// Discard error; can only be ErrMaxRetry and we use math.MaxInt retries.
	duration, _ := m.retryBackoff.Next()
	m.resetTimer(duration)
	m.mu.Unlock()

	m.logger.Warn("scheduled TPM certificate retry after failure", "duration", duration)
}
