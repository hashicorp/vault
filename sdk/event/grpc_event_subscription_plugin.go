// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/eventplugin"
	"google.golang.org/grpc"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var HandshakeConfig = plugin.HandshakeConfig{
	MagicCookieKey:   "VAULT_EVENT_SUBSCRIPTION_PLUGIN",
	MagicCookieValue: "cd935c26-17b5-4f24-bca2-a66621a1c944",
}

// Factory is the factory function to create a dbplugin Database.
type Factory func() (interface{}, error)

type GRPCEventSubscriptionPlugin struct {
	Impl EventSubscriptionPlugin

	// Embedding this will disable the netRPC protocol
	plugin.NetRPCUnsupportedPlugin
}

var (
	_ plugin.Plugin     = &GRPCEventSubscriptionPlugin{}
	_ plugin.GRPCPlugin = &GRPCEventSubscriptionPlugin{}
)

func (d GRPCEventSubscriptionPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	server := gRPCServer{instance: d.Impl}
	eventplugin.RegisterEventSubscribePluginServiceServer(s, &server)
	return nil
}

func (GRPCEventSubscriptionPlugin) GRPCClient(doneCtx context.Context, _ *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	client := gRPCClient{
		client:  eventplugin.NewEventSubscribePluginServiceClient(c),
		doneCtx: doneCtx,
	}
	return client, nil
}
