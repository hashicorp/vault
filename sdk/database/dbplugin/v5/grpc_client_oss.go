// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package dbplugin

import (
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

type entGRPCClient struct{}

func (c gRPCClient) Close() error {
	ctx, cancel := getContextWithTimeout(pluginutil.PluginGRPCTimeoutClose)
	defer cancel()

	_, err := c.client.Close(ctx, &proto.Empty{})
	if err != nil {
		if c.doneCtx.Err() != nil {
			return ErrPluginShutdown
		}
		return err
	}
	return nil
}
