// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package builder

import (
	"context"
	"errors"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type GenericBackend[CC, C any] struct {
	*framework.Backend
	lock      sync.RWMutex
	client    *C
	newClient func(*CC) (*C, error)
}

var myBackend any

func (gb *GenericBackend[CC, C]) setBackend() {
	myBackend = gb
}

func (gb *GenericBackend[CC, C]) invalidate(_ context.Context, key string) {
	if key == "config" {
		gb.reset()
	}
}

func (gb *GenericBackend[CC, C]) reset() {
	gb.lock.Lock()
	defer gb.lock.Unlock()
	gb.client = nil
}

// getClient locks the backend as it configures and creates a
// a new client for the target API
func (gb *GenericBackend[CC, C]) getClient(ctx context.Context, s logical.Storage) (*C, error) {
	gb.lock.Lock()
	defer gb.lock.Unlock()

	client := gb.client
	if client != nil {
		return client, nil
	}

	config, err := gb.getConfig(ctx, s)
	if err != nil {
		return nil, err
	}

	if gb.client == nil && config == nil {
		config = new(CC)
	}

	gb.client, err = gb.newClient(config)
	if err != nil {
		return nil, err
	}

	return gb.client, nil
}

// Factory returns a new backend as logical.Backend
func factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := myBackend.(logical.Backend)
	if b == nil {
		return nil, errors.New("backend has not been built")
	}
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}
