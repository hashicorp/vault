// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func (c *Core) entInitWALPassThrough() func() {
	return nil
}

func (c *Core) entCheckStoredLicense(conf *CoreConfig) error {
	return nil
}

func (c *Core) entIsLicenseAutoloaded() bool {
	return false
}

func (c *Core) entCheckLicenseInit() error {
	return nil
}

func (c *Core) EntGetLicenseState() (*LicenseState, error) {
	return nil, nil
}

func (c *Core) EntGetLicense() (string, error) {
	return "", nil
}

func (c *Core) EntReloadLicense() error {
	return nil
}

func (c *Core) entPostUnseal(isStandby bool) error {
	return nil
}

func (c *Core) entPreSeal() error {
	return nil
}

func (c *Core) entSetupFilteredPaths() error {
	return nil
}

func (c *Core) entSetupQuotas(ctx context.Context) error {
	return nil
}

func (c *Core) entSetupAPILock(ctx context.Context) error {
	return nil
}

func (c *Core) entBlockRequestIfError(nsPath, requestPath string) error {
	return nil
}

func (c *Core) entStartReplication() error {
	return nil
}

func (c *Core) entStopReplication() error {
	return nil
}

func (c *Core) EntLastWAL() uint64 {
	return 0
}

func (c *Core) EntLastPerformanceWAL() uint64 {
	return 0
}

func (c *Core) EntLastDRWAL() uint64 {
	return 0
}

func (c *Core) EntDRMerkleRoot() string {
	return ""
}

func (c *Core) EntPerformanceMerkleRoot() string {
	return ""
}

func (c *Core) EntLastRemoteWAL() uint64 {
	return 0
}

func (c *Core) entLastRemoteUpstreamWAL() uint64 {
	return 0
}

func (c *Core) EntWaitUntilWALShipped(ctx context.Context, index uint64) bool {
	return true
}

func (c *Core) SecretsSyncLicensedActivated() bool { return false }

func (c *Core) IsMultisealEnabled() bool { return false }

func (c *Core) SetMultisealEnabled(_ bool) {}

func (c *Core) ReloadReplicationCanaryWriteInterval() {}

func (c *Core) GetReplicationLagMillisIgnoreErrs() int64 { return 0 }

func (c *Core) ReloadOverloadController() {}
