// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"crypto/ed25519"
	"testing"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func GenerateTestLicenseKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) { return nil, nil, nil }
func testGetLicensingConfig(key ed25519.PublicKey) *LicensingConfig           { return &LicensingConfig{} }
func testExtraTestCoreSetup(testing.TB, ed25519.PrivateKey, *TestClusterCore) {}
func testAdjustUnderlyingStorage(tcc *TestClusterCore) {
	tcc.UnderlyingStorage = tcc.physical
}
func testApplyEntBaseConfig(coreConfig, base *CoreConfig) {}
