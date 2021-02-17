// +build !enterprise

package vault

import (
	testing "github.com/mitchellh/go-testing-interface"
)

func testGenerateCoreKeys() (interface{}, interface{}, error)          { return nil, nil, nil }
func testGetLicensingConfig(interface{}) *LicensingConfig              { return &LicensingConfig{} }
func testExtraTestCoreSetup(testing.TB, interface{}, *TestClusterCore) {}
func testAdjustUnderlyingStorage(tcc *TestClusterCore) {
	tcc.UnderlyingStorage = tcc.physical
}
func testApplyEntBaseConfig(coreConfig, base *CoreConfig) {}
