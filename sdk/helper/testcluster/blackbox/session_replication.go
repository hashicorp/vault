// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

// GetDRReplicationMode gets the DR replication mode of the node
func (s *Session) GetDRReplicationMode() (string, error) {
	secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
		return s.Client.Logical().Read("sys/replication/dr/status")
	})
	if err != nil {
		return "", err
	}

	if secret == nil {
		return "", fmt.Errorf("empty status response")
	}

	mode, ok := secret.Data["mode"]
	if !ok {
		return "", fmt.Errorf("no mode field in status response")
	}

	return mode.(string), nil
}

func (s *Session) AssertReplicationDisabled() {
	s.assertReplicationStatus("ce", "disabled")
}

func (s *Session) AssertDRReplicationStatus(expectedMode string) {
	s.assertReplicationStatus("dr", expectedMode)
}

func (s *Session) AssertPerformanceReplicationStatus(expectedMode string) {
	s.assertReplicationStatus("performance", expectedMode)
}

func (s *Session) assertReplicationStatus(which, expectedMode string) {
	s.t.Helper()

	secret, err := s.WithRootNamespace(func() (*api.Secret, error) {
		return s.Client.Logical().Read("sys/replication/status")
	})

	require.NoError(s.t, err)
	require.NotNil(s.t, secret)

	data := s.AssertSecret(secret).Data()

	if which == "ce" {
		data.HasKey("mode", "disabled")
	} else {
		data.GetMap(which).HasKey("mode", expectedMode)
	}
}
