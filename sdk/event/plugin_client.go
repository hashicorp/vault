// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"
	"errors"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/eventplugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var _ logical.PluginVersioner = (*EventSubscriptionPluginClient)(nil)

type EventSubscriptionPluginClient struct {
	client pluginutil.PluginClient
	EventSubscriptionPlugin
}

func (ec *EventSubscriptionPluginClient) PluginVersion() logical.PluginVersion {
	_, version := ec.Type()
	return logical.PluginVersion{Version: version}
}

// This wraps the Close call and ensures we both close the backend connection
// and kill the plugin.
func (ec *EventSubscriptionPluginClient) Close(ctx context.Context) error {
	err := ec.EventSubscriptionPlugin.Close(ctx)
	_ = ec.client.Close()
	return err
}

// pluginSets is the map of plugins we can dispense.
var pluginSets = map[int]plugin.PluginSet{
	1: {
		"event_subscription": &GRPCEventSubscriptionPlugin{},
	},
}

// NewPluginClient returns a eventSubscriptionPluginClient with a connection to a running
// plugin.
func NewPluginClient(ctx context.Context, sys pluginutil.RunnerUtil, config pluginutil.PluginClientConfig) (EventSubscriptionPlugin, error) {
	pluginClient, err := sys.NewPluginClient(ctx, config)
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := pluginClient.Dispense("event_subscription")
	if err != nil {
		return nil, err
	}

	// We should have an EventSubscriptionPlugin type now. This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	var ev EventSubscriptionPlugin
	switch c := raw.(type) {
	case *gRPCClient:
		// This is an abstraction leak from go-plugin but it is necessary in
		// order to enable multiplexing on multiplexed plugins
		c.client = eventplugin.NewEventSubscribePluginServiceClient(pluginClient.Conn())

		ev = c
	default:
		return nil, errors.New("unsupported client type")
	}

	return &EventSubscriptionPluginClient{
		client:                  pluginClient,
		EventSubscriptionPlugin: ev,
	}, nil
}
