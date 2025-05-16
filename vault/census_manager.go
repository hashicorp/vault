// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

// CensusManager provides stub behavior for CE, simplifying the logic between CE
// and ENT. This will always be marked active: false.
type CensusManager struct {
	active bool
	logger hclog.Logger
}

// CensusManagerConfig is empty on CE.
type CensusManagerConfig struct{}

// NewCensusManager sets up the stub CensusManager on CE with active: false.
func NewCensusManager(logger hclog.Logger, conf CensusManagerConfig, storage logical.Storage) (*CensusManager, error) {
	return &CensusManager{
		active: false,
		logger: logger,
	}, nil
}

// setupCensusManager is a stub on CE.
func (c *Core) setupCensusManager(ctx context.Context) error {
	return nil
}

// reloadDevelopmentClusterSetting is a stub on CE.
func (c *Core) reloadDevelopmentClusterSetting() error {
	return nil
}

// persistDevelopmentClusterSetting is a stub on CE.
func (c *Core) persistDevelopmentClusterSetting(ctx context.Context) error {
	return nil
}

// getDevelopmentClusterSetting is a stub on CE.
func (c *Core) getDevelopmentClusterSetting(ctx context.Context) (bool, error) {
	return false, nil
}

// BillingStart is a stub on CE.
func (cm *CensusManager) BillingStart() time.Time {
	return time.Time{}
}

// StartManualReportingSnapshots is a stub for CE.
func (cm *CensusManager) StartManualReportingSnapshots() {}
