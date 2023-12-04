// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"fmt"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

// Serve is called from within a plugin and wraps the provided
// Database implementation in a databasePluginRPCServer object and starts a
// RPC server.
func Serve(ev EventSubscriptionPlugin) {
	plugin.Serve(ServeConfig(ev))
}

func ServeConfig(ev EventSubscriptionPlugin) *plugin.ServeConfig {
	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// pluginSets is the map of plugins we can dispense.
	pluginSets := map[int]plugin.PluginSet{
		1: {
			"event_subscriber": &GRPCEventSubscriptionPlugin{
				Impl:                    ev,
				NetRPCUnsupportedPlugin: plugin.NetRPCUnsupportedPlugin{},
			},
		},
	}

	conf := &plugin.ServeConfig{
		HandshakeConfig:  HandshakeConfig,
		VersionedPlugins: pluginSets,
		GRPCServer:       plugin.DefaultGRPCServer,
	}

	return conf
}
