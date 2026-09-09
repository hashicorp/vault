// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"errors"
	"time"

	"github.com/hashicorp/vault/api"
)

// getClusterNodeCount returns the number of nodes in the cluster
func (s *Session) getClusterNodeCount() (int, error) {
	nodes, err := s.haNodes()
	if err != nil {
		return 0, err
	}

	return len(nodes), nil
}

// haNodes returns a slice of nodes from the ha-status endpoint.
func (s *Session) haNodes() ([]map[string]any, error) {
	res := []map[string]any{}

	return res, s.Req(
		func(c *api.Client) error {
			status, err := c.Logical().Read("sys/ha-status")
			if err != nil {
				return err
			}
			if status == nil {
				return errors.New("no ha-status returned")
			}

			nodesAny, ok := status.Data["nodes"]
			if !ok {
				return errors.New("no HA nodes found in ha-status")
			}

			nodes, ok := nodesAny.([]any)
			if !ok {
				return errors.New("invalid ha-status response body")
			}

			for _, node := range nodes {
				nv, ok := node.(map[string]any)
				if !ok {
					return errors.New("malformed node in ha-status response body")
				}
				res = append(res, nv)
			}

			return nil
		},
		WithClientTimeout(2*time.Second),
	)
}

// haActiveNode returns the active node from ha-status.
func (s *Session) haActiveNode() (map[string]any, error) {
	nodes, err := s.haNodes()
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		active, ok := node["active_node"]
		if !ok {
			continue
		}

		activeVal, ok := active.(bool)
		if !ok {
			continue
		}

		if activeVal {
			return node, nil
		}
	}

	return nil, errors.New("no active node in ha-status")
}
