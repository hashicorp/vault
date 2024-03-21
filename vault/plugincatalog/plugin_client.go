// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// pluginClient represents a connection to a plugin process
type pluginClient struct {
	logger log.Logger

	// id is the connection ID
	id       string
	pluginID string

	// client handles the lifecycle of a plugin process
	// multiplexed plugins share the same client
	client      *plugin.Client
	clientConn  grpc.ClientConnInterface
	cleanupFunc func() error
	reloadFunc  func() error

	plugin.ClientProtocol
}

func (p *pluginClient) Conn() grpc.ClientConnInterface {
	return p.clientConn
}

func (p *pluginClient) Reload() error {
	p.logger.Debug("reload external plugin process")
	return p.reloadFunc()
}

// Close calls the plugin client's cleanupFunc to do any necessary cleanup on
// the plugin client and the PluginCatalog. This implements the
// plugin.ClientProtocol interface.
func (p *pluginClient) Close() error {
	p.logger.Debug("cleaning up plugin client connection", "id", p.id)
	return p.cleanupFunc()
}
