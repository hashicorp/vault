// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

func newGRPCObservationsClient(conn *grpc.ClientConn) *GRPCObservationsClient {
	return &GRPCObservationsClient{
		client: pb.NewObservationsClient(conn),
	}
}

type GRPCObservationsClient struct {
	client pb.ObservationsClient
}

var _ logical.ObservationRecorder = (*GRPCObservationsClient)(nil)

func (s *GRPCObservationsClient) RecordObservationFromPlugin(ctx context.Context, observationType string, data map[string]interface{}) error {
	dataAsPb, err := structpb.NewStruct(data)
	if err != nil {
		return err
	}

	_, err = s.client.RecordObservation(ctx, &pb.RecordObservationRequest{
		ObservationType: observationType,
		Data:            dataAsPb,
	})
	return err
}

type GRPCObservationsServer struct {
	pb.UnimplementedObservationsServer
	impl logical.ObservationRecorder
}

func (s *GRPCObservationsServer) RecordObservation(ctx context.Context, req *pb.RecordObservationRequest) (*pb.Empty, error) {
	if s.impl == nil {
		return &pb.Empty{}, nil
	}

	err := s.impl.RecordObservationFromPlugin(ctx, req.ObservationType, req.Data.AsMap())
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
