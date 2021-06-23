// +build !enterprise

package vault

import (
	"crypto/ed25519"

	testing "github.com/mitchellh/go-testing-interface"
)

func GenerateTestLicenseKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) { return nil, nil, nil }
func testGetLicensingConfig(key ed25519.PublicKey) *LicensingConfig           { return &LicensingConfig{} }
func testExtraTestCoreSetup(testing.T, ed25519.PrivateKey, *TestClusterCore)  {}
func testAdjustUnderlyingStorage(tcc *TestClusterCore) {
	tcc.UnderlyingStorage = tcc.physical
}
func testApplyEntBaseConfig(coreConfig, base *CoreConfig) {}
