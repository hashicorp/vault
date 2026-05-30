// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"fmt"
	"time"
)

// AssertClusterHealthy verifies that the cluster is healthy, with fallback for managed environments
// like HCP where raft APIs may not be accessible. This is the recommended method for general
// cluster health checks in blackbox tests. Uses backoff helper for reasonable retry logic.
//
// Deprecated: Use s.EventuallyClusterHealthy() with an explicit timeout
// retry times.
func (s *Session) AssertClusterHealthy() {
	s.t.Helper()

	s.EventuallyClusterHealthyUnsealed(1 * time.Minute)
}

// EventuallyClusterHealthyUnsealed verifies that the cluster is healthy and
// unsealed.
func (s *Session) EventuallyClusterHealthyUnsealed(timeout time.Duration) {
	s.t.Helper()

	hasActiveNode := func() error {
		// First, wait until we have an active HA node
		_, err := s.haActiveNode()
		return err
	}

	autopilotHealthyIfRaft := func() error {
		// Now, make sure autopilot is healthy if we're using integrated storage
		// and we have the necessary permissions.
		if s.GetParentNamespace() != "" {
			s.t.Log("Skipping autopilot health check because we've been configured with a parent namespace and autopilot needs root namespace access")
			return nil
		}

		storage, err := s.getConfigStorageType()
		if err != nil {
			return err
		}

		if storage != "raft" {
			return nil
		}

		healthy, err := s.autopilotStateHealthy()
		if err != nil {
			return err
		}

		if !healthy {
			return fmt.Errorf("raft autopilot state is not healthy")
		}

		return nil
	}

	clusterUnsealed := func() error {
		sealed, err := s.sealed()
		if err != nil {
			return err
		}

		if sealed {
			return fmt.Errorf("cluster is sealed")
		}

		return nil
	}

	// Go through our checks and wait for each.
	started := time.Now()
	for _, check := range []func() error{
		hasActiveNode,
		autopilotHealthyIfRaft,
		clusterUnsealed,
	} {
		if timeout < 0 {
			s.t.Fatal("timed out waiting for cluster to be healthy")
		}
		s.EventuallyWithTimeout(check, timeout)
		timeout -= time.Since(started)
	}
}
