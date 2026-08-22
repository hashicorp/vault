// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

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
func TestApplyEntBaseConfig(coreConfig, base *CoreConfig) {}
