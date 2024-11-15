// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package linkstatus provides a gRPC service that allows HashiCorp Cloud Platform (HCP) to query the status of the
// link library.
package linkstatus

import (
	"context"

	pb "github.com/hashicorp/hcp-link/gen/proto/go/hashicorp/cloud/hcp_link/link_status/v1"
	"github.com/hashicorp/hcp-link/pkg/config"
)

const (
	linkStatusVersion = "0.0.1"
)

// gRPC service to reports the status of the link library
type Service struct {
	// Config contains all dependencies as well as information about the node
	// Link is running on.
	Config *config.Config

	pb.UnimplementedLinkStatusServiceServer
}

// GetLinkStatus will be used to fetch the node’s link specific status.
func (s *Service) GetLinkStatus(_ context.Context, _ *pb.GetLinkStatusRequest) (*pb.GetLinkStatusResponse, error) {
	return &pb.GetLinkStatusResponse{
		NodeId:  s.Config.NodeID,
		Version: linkStatusVersion,
		Features: &pb.Features{
			NodeStatusReporting: &pb.FeatureNodeStatusReporting{
				Enabled: s.Config.NodeStatusReporter != nil,
			},
		},
	}, nil
}
