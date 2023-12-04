// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hashicorp/vault/sdk/eventplugin"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	_ EventSubscriptionPlugin = (*gRPCClient)(nil)
	_ logical.PluginVersioner = (*gRPCClient)(nil)
)

type gRPCClient struct {
	client  eventplugin.EventSubscribePluginServiceClient
	doneCtx context.Context

	eventSenders     map[string]func(string) error
	eventSendersLock sync.RWMutex
}

func (c *gRPCClient) Initialize(ctx context.Context) error {
	_, err := c.client.Initialize(ctx, &eventplugin.InitializeRequest{})
	if err != nil {
		return fmt.Errorf("unable to initialize: %s", err.Error())
	}
	return nil
}

// TODO: share this with marhsalling.go
func mapToStruct(m map[string]interface{}) (*structpb.Struct, error) {
	// Convert any json.Number typed values to float64, since the
	// type does not have a conversion mapping defined in structpb
	for k, v := range m {
		if n, ok := v.(json.Number); ok {
			nf, err := n.Float64()
			if err != nil {
				return nil, err
			}

			m[k] = nf
		}
	}

	return structpb.NewStruct(m)
}

func (c *gRPCClient) Subscribe(ctx context.Context, request *SubscribeRequest) error {
	config, err := mapToStruct(request.Config)
	if err != nil {
		return err
	}
	_, err = c.client.Subscribe(ctx, &eventplugin.SubscribeRequest{
		Config:           config,
		SubscriptionId:   request.SubscriptionID,
		VerifyConnection: request.VerifyConnection,
	})
	return err
}

func makeSender(server eventplugin.EventSubscribePluginService_SendSubscriptionEventsClient) func(received string) error {
	return func(eventJson string) error {
		return server.Send(&eventplugin.SubscriptionEvent{
			EventJson: eventJson,
		})
	}
}

func (c *gRPCClient) getOrCreateEventSender(subscriptionID string) (func(received string) error, error) {
	c.eventSendersLock.RLock()
	sender, ok := c.eventSenders[subscriptionID]
	c.eventSendersLock.RUnlock()
	if ok {
		return sender, nil
	}

	c.eventSendersLock.Lock()
	defer c.eventSendersLock.Unlock()
	// Check again to avoid a race condition.
	sender, ok = c.eventSenders[subscriptionID]
	if ok {
		return sender, nil
	}
	server, err := c.client.SendSubscriptionEvents(context.Background())
	if err != nil {
		return nil, err
	}
	// send a message with the subscription ID to initialize the subscription
	err = server.Send(&eventplugin.SubscriptionEvent{
		SubscriptionId: subscriptionID,
	})
	if err != nil {
		return nil, err
	}
	f := makeSender(server)
	c.eventSenders[subscriptionID] = f
	return f, nil
}

func (c *gRPCClient) SendSubscriptionEvent(subscriptionID string, eventJson string) error {
	sender, err := c.getOrCreateEventSender(subscriptionID)
	if err != nil {
		return err
	}
	if eventJson != "" {
		return sender(eventJson)
	}
	return nil
}

func (c *gRPCClient) Unsubscribe(ctx context.Context, subscriptionID string) error {
	_, err := c.client.Unsubscribe(ctx, &eventplugin.UnsubscribeRequest{SubscriptionId: subscriptionID})
	return err
}

func (c *gRPCClient) Type() (string, string) {
	resp, err := c.client.Type(context.Background(), &eventplugin.TypeRequest{})
	if err != nil {
		return "", ""
	}
	return resp.PluginType, resp.PluginVersion
}

func (c *gRPCClient) Close(ctx context.Context) error {
	_, err := c.client.Close(ctx, &eventplugin.CloseRequest{})
	return err
}

func (c *gRPCClient) PluginVersion() logical.PluginVersion {
	info, err := c.client.Type(context.Background(), nil)
	if info == nil || err != nil {
		return logical.EmptyPluginVersion
	}
	return logical.PluginVersion{Version: info.PluginVersion}
}
