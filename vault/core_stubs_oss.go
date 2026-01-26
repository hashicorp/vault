// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"fmt"
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

func (c *Core) GetCurrentWALHeader() string {
	return ""
}

func (c *Core) IsReplicated(secondaryID, namespacePath, mountPathRelative string) bool {
	return false
}

func (c *Core) SecretsSyncLicensedActivated() bool { return false }

func (c *Core) IsMultisealEnabled() bool { return false }

func (c *Core) SetMultisealEnabled(_ bool) {}

func (c *Core) ReloadReplicationCanaryWriteInterval() {}

func (c *Core) GetReplicationLagMillisIgnoreErrs() int64 { return 0 }

func (c *Core) ReloadOverloadController() {}

func (c *Core) EntSetupUIDefaultAuth(ctx context.Context) error { return nil }

// entGetPluginCacheDir returns empty string and an error indicating that this is an
// enterprise-only feature. This is used to prevent the use of the plugin cache
func (c *Core) entGetPluginCacheDir() (string, error) {
	return "", fmt.Errorf("enterprise only feature")
}

// entGetPluginRuntimeDir returns empty string and an error indicating that this is an
// enterprise-only feature
func (c *Core) entGetPluginRuntimeDir() (string, error) {
	return "", fmt.Errorf("enterprise only feature")
}

// entJoinPluginDir returns empty string and an error indicating that this is an
// enterprise-only feature
func (c *Core) entJoinPluginDir(_ string) (string, error) {
	return "", fmt.Errorf("enterprise only feature")
}

// IsMountTypeAllowed returns true if a given secret engine mount type is permitted.
// Forbidden mount types should be refused in mount requests, and any existing mounts
// of that type should return an error on any routed external requests.
func (c *Core) IsMountTypeAllowed(mountType string) bool {
	return true
}

// IsFlagEnabled returns true if the named flag is set in HCL config feature_flags.
func (c *Core) IsFlagEnabled(name string) bool {
	return false
}
