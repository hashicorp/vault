// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/hashicorp/vault/version"
)

func init() {
	// The BuildDate is set as part of the build process in CI so we need to
	// initialize it for testing.
	if version.BuildDate == "" {
		version.BuildDate = time.Now().UTC().AddDate(-1, 0, 0).Format(time.RFC3339)
	}
}

func GenerateTestLicenseKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) { return nil, nil, nil }
func testGetLicensingConfig(key ed25519.PublicKey) *LicensingConfig           { return &LicensingConfig{} }
func testExtraTestCoreSetup(testing.TB, ed25519.PrivateKey, *TestClusterCore) {}
func testAdjustUnderlyingStorage(tcc *TestClusterCore) {
	tcc.UnderlyingStorage = tcc.physical
}
func testApplyEntBaseConfig(coreConfig, base *CoreConfig) {}
