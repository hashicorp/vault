// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
)

const (
	requestCounterDatePathFormat = "2006/01"

	// This storage path stores both the request counters in this file, and the activity log.
	countersSubPath = "counters/"
)

// ActiveTokens contains the number of active tokens.
type ActiveTokens struct {
	// ServiceTokens contains information about the number of active service
	// tokens.
	ServiceTokens TokenCounter `json:"service_tokens"`
}

// TokenCounter counts the number of tokens
type TokenCounter struct {
	// Total is the total number of tokens
	Total int `json:"total"`
}

// countActiveTokens returns the number of active tokens
func (c *Core) countActiveTokens(ctx context.Context) (*ActiveTokens, error) {
	// Get all of the namespaces
	ns := c.collectNamespaces()

	// Count the tokens under each namespace
	total := 0
	for i := 0; i < len(ns); i++ {
		ids, err := c.tokenStore.idView(ns[i]).List(ctx, "")
		if err != nil {
			return nil, err
		}
		total += len(ids)
	}

	return &ActiveTokens{
		ServiceTokens: TokenCounter{
			Total: total,
		},
	}, nil
}

// ActiveEntities contains the number of active entities.
type ActiveEntities struct {
	// Entities contains information about the number of active entities.
	Entities EntityCounter `json:"entities"`
}

// EntityCounter counts the number of entities
type EntityCounter struct {
	// Total is the total number of entities
	Total int `json:"total"`
}

// countActiveEntities returns the number of active entities
func (c *Core) countActiveEntities(ctx context.Context) (*ActiveEntities, error) {
	count, err := c.identityStore.countEntities()
	if err != nil {
		return nil, err
	}

	return &ActiveEntities{
		Entities: EntityCounter{
			Total: count,
		},
	}, nil
}
