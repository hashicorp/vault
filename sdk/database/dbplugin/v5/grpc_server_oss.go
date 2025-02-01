// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package dbplugin

import (
	"sync"

	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"github.com/hashicorp/vault/sdk/logical"
)

type gRPCServer struct {
	proto.UnimplementedDatabaseServer
	logical.UnimplementedPluginVersionServer

	// holds the non-multiplexed Database
	// when this is set the plugin does not support multiplexing
	singleImpl Database

	// instances holds the multiplexed Databases
	instances   map[string]Database
	factoryFunc func() (interface{}, error)

	sync.RWMutex
}
