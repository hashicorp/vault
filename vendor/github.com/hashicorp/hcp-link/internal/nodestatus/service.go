// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nodestatus

import (
	"context"

	pb "github.com/hashicorp/hcp-link/gen/proto/go/hashicorp/cloud/hcp_link/node_status/v1"
)

type Service struct {
	// Collector implements the logic needed to collect node status information
	Collector *Collector

	pb.UnimplementedNodeStatusServiceServer
}

// GetNodeStatus will be used to regularly fetch the node’s current status.
func (s *Service) GetNodeStatus(ctx context.Context, _ *pb.GetNodeStatusRequest) (*pb.GetNodeStatusResponse, error) {
	status, err := s.Collector.CollectPb(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetNodeStatusResponse{NodeStatus: status}, nil
}
