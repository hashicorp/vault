// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dbplugin

import (
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var _ logical.PluginVersioner = (*DatabasePluginClient)(nil)

type DatabasePluginClient struct {
	client pluginutil.PluginClient
	Database
}

func (dc *DatabasePluginClient) PluginVersion() logical.PluginVersion {
	if versioner, ok := dc.Database.(logical.PluginVersioner); ok {
		return versioner.PluginVersion()
	}
	return logical.EmptyPluginVersion
}

// This wraps the Close call and ensures we both close the database connection
// and kill the plugin.
func (dc *DatabasePluginClient) Close() error {
	err := dc.Database.Close()
	dc.client.Close()

	return err
}

// pluginSets is the map of plugins we can dispense.
var PluginSets = map[int]plugin.PluginSet{
	5: {
		"database": &GRPCDatabasePlugin{},
	},
	6: {
		"database": &GRPCDatabasePlugin{},
	},
}
