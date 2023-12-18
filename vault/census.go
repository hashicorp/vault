// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import "time"

// CensusAgent is a stub for OSS
type CensusReporter interface{}

// setupCensusAgent is a stub for OSS.
func (c *Core) setupCensusAgent() error      { return nil }
func (c *Core) BillingStart() time.Time      { return time.Time{} }
func (c *Core) CensusLicensingEnabled() bool { return false }
func (c *Core) CensusAgent() CensusReporter  { return nil }
func (c *Core) ReloadCensus() error          { return nil }
func (c *Core) teardownCensusAgent() error   { return nil }
