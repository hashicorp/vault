// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"errors"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

func (s *Session) AssertAutopilotHealthy() {
	s.t.Helper()

	// query the autopilot state endpoint and verify that all nodes are healthy according to autopilot
	healthy, err := s.autopilotStateHealthy()
	require.NoError(s.t, err)
	require.Truef(s.t, healthy, "expected autopilot state to be healthy")
}

func (s *Session) autopilotState() (*api.Secret, error) {
	// query the autopilot state endpoint and verify that all nodes are healthy according to autopilot
	var state *api.Secret
	return state, s.Req(
		func(c *api.Client) error {
			var err error
			state, err = c.Logical().Read("sys/storage/raft/autopilot/state")
			return err
		},
		WithClientRootNamespace(),
		WithClientTimeout(2*time.Second),
	)
}

func (s *Session) autopilotStateHealthy() (bool, error) {
	state, err := s.autopilotState()
	if err != nil {
		return false, err
	}

	if state == nil {
		return false, errors.New("no raft data response")
	}

	health, ok := state.Data["healthy"]
	if !ok {
		return false, errors.New("raft data missing 'healthy' key")
	}

	healthy, ok := health.(bool)
	if !ok {
		return false, errors.New("raft data 'healthy' key is unknown type")
	}

	return healthy, nil
}
