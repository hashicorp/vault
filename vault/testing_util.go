// +build !enterprise

package vault

import "github.com/mitchellh/go-testing-interface"

func testGenerateCoreKeys() (interface{}, interface{}, error)                   { return nil, nil, nil }
func testGetLicensingConfig(interface{}) *LicensingConfig                       { return &LicensingConfig{} }
func testAdjustTestCore(*CoreConfig, *TestClusterCore)                          {}
func testExtraClusterCoresTestSetup(testing.T, interface{}, []*TestClusterCore) {}
