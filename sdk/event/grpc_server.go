// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"

	"github.com/hashicorp/vault/sdk/eventplugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ eventplugin.EventSubscribePluginServiceServer = (*gRPCServer)(nil)

// Translates GRPC calls to the EventSubscriptionPlugin backend.
// The subscribe interface is itself multiplexed, so we do not need to multiplex to multiple plugin instances.
type gRPCServer struct {
	eventplugin.UnimplementedEventSubscribePluginServiceServer

	instance EventSubscriptionPlugin
}

func (g *gRPCServer) Initialize(ctx context.Context, _ *eventplugin.InitializeRequest) (*eventplugin.InitializeResponse, error) {
	err := g.instance.Initialize(ctx)
	if err != nil {
		return &eventplugin.InitializeResponse{}, status.Errorf(codes.Internal, "failed to initialize: %s", err)
	}

	return &eventplugin.InitializeResponse{}, nil
}

func (g *gRPCServer) Subscribe(ctx context.Context, request *eventplugin.SubscribeRequest) (*eventplugin.SubscribeResponse, error) {
	err := g.instance.Subscribe(ctx, &SubscribeRequest{
		SubscriptionID:   request.SubscriptionId,
		Config:           request.Config.AsMap(),
		VerifyConnection: request.VerifyConnection,
	})
	if err != nil {
		return &eventplugin.SubscribeResponse{}, status.Errorf(codes.Internal, "failed to subscribe: %s", err)
	}
	return &eventplugin.SubscribeResponse{}, nil
}

func (g *gRPCServer) SendSubscriptionEvents(server eventplugin.EventSubscribePluginService_SendSubscriptionEventsServer) error {
	var subscriptionID string
	defer func() {
		_ = server.SendAndClose(nil)
	}()
	for {
		event, err := server.Recv()
		if err != nil {
			return status.Errorf(codes.Internal, "error receiving message: %s", err)
		}
		if event.SubscriptionId != "" {
			subscriptionID = event.SubscriptionId
		}
		if event.EventJson != "" {
			err = g.instance.SendSubscriptionEvent(subscriptionID, event.EventJson)
			if err != nil {
				return status.Errorf(codes.Internal, "error sending event to backend: %s", err)
			}
		}
	}
}

func (g *gRPCServer) Unsubscribe(ctx context.Context, request *eventplugin.UnsubscribeRequest) (*eventplugin.UnsubscribeResponse, error) {
	err := g.instance.Unsubscribe(ctx, request.SubscriptionId)
	if err != nil {
		return &eventplugin.UnsubscribeResponse{}, status.Errorf(codes.Internal, "error unsubscribing: %s", err)
	}
	return &eventplugin.UnsubscribeResponse{}, nil
}

func (g *gRPCServer) Type(_ context.Context, _ *eventplugin.TypeRequest) (*eventplugin.TypeResponse, error) {
	name, version := g.instance.Type()
	return &eventplugin.TypeResponse{
		PluginType:    name,
		PluginVersion: version,
	}, nil
}

func (g *gRPCServer) Close(ctx context.Context, _ *eventplugin.CloseRequest) (*eventplugin.CloseResponse, error) {
	err := g.instance.Close(ctx)
	if err != nil {
		return &eventplugin.CloseResponse{}, status.Errorf(codes.Internal, "error closing: %s", err)
	}
	return &eventplugin.CloseResponse{}, nil
}
