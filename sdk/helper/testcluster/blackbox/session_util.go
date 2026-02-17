// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"time"

	"github.com/hashicorp/vault/api"
)

// Eventually retries the function 'fn' until it returns nil or timeout occurs.
func (s *Session) Eventually(fn func() error) {
	s.t.Helper()

	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	var lastErr error

	for {
		select {
		case <-timeout:
			s.t.Fatalf("Eventually failed after 5s. Last error: %v", lastErr)
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
