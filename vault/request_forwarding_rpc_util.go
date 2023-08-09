// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package vault

import (
	"context"
)

func (s *forwardedRequestRPCServer) PerformanceStandbyElectionRequest(in *PerfStandbyElectionInput, reqServ RequestForwarding_PerformanceStandbyElectionRequestServer) error {
	return nil
}

type ReplicationTokenInfo struct{}

func (c *forwardingClient) PerformanceStandbyElection(ctx context.Context) (*ReplicationTokenInfo, error) {
	return nil, nil
}
