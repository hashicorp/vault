// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"time"

	"github.com/hashicorp/vault/api"
)

// Eventually retries the function 'fn' until it returns nil or timeout occurs.
func (s *Session) Eventually(fn func() error) {
	s.EventuallyWithTimeout(fn, 5*time.Second)
}

// EventuallyWithTimeout retries the function 'fn' until it returns nil or timeout occurs.
// Use this for operations that may take longer than the default 5 seconds.
func (s *Session) EventuallyWithTimeout(fn func() error, timeout time.Duration) {
	s.t.Helper()

	timeoutChan := time.After(timeout)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	var lastErr error

	for {
		select {
		case <-timeoutChan:
			s.t.Fatalf("Eventually failed after %v. Last error: %v", timeout, lastErr)
		case <-ticker.C:
			lastErr = fn()
			if lastErr == nil {
				return
			}
		}
	}
}

func (s *Session) WithRootNamespace(fn func() (*api.Secret, error)) (*api.Secret, error) {
	s.t.Helper()

	oldNamespace := s.Client.Namespace()
	defer s.Client.SetNamespace(oldNamespace)
	s.Client.ClearNamespace()

	return fn()
}
