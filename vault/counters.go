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
