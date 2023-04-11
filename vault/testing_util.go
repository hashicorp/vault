// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package vault

import (
	"crypto/ed25519"
)

func GenerateTestLicenseKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) { return nil, nil, nil }
func testGetLicensingConfig(key ed25519.PublicKey) *LicensingConfig           { return &LicensingConfig{} }
func testAdjustUnderlyingStorage(tcc *TestClusterCore) {
	tcc.UnderlyingStorage = tcc.physical
}
func testApplyEntBaseConfig(coreConfig, base *CoreConfig) {}
