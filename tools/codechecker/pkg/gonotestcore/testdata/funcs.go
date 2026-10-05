// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package testdata

import "testing"

// The following are stand-in definitions so that the calls below compile.
// Definitions are not calls, so they must never be flagged by the analyzer.
func TestCore(t testing.TB)             {}
func TestCoreUnsealed(t testing.TB)     {}
func TestCore_NotAHelper(t *testing.T)  {}
func testControlGroupCore(t *testing.T) {}
func mockExpiration(t *testing.T)       {}
func helperOK(t *testing.T)             {}
func NewTestCluster(t *testing.T)       {}

// TestUsesBannedHelpers calls helpers that should be flagged by the analyzer.
func TestUsesBannedHelpers(t *testing.T) {
	TestCore(t)             // want `TestCore is part of the TestCore family`
	TestCoreUnsealed(t)     // want `TestCoreUnsealed is part of the TestCore family`
	testControlGroupCore(t) // want `testControlGroupCore is part of the TestCore family`
	mockExpiration(t)       // want `mockExpiration is part of the TestCore family`
}

// TestUsesAllowedHelpers calls helpers that must not be flagged: an underscore
// after "TestCore" means it is a regular test function, not a helper, and
// helperOK is unrelated.
func TestUsesAllowedHelpers(t *testing.T) {
	TestCore_NotAHelper(t)
	helperOK(t)
}

// TestUsesNewTestCluster calls NewTestCluster and validates it's allowed
func TestUsesNewTestCluster(t *testing.T) {
	NewTestCluster(t)
}
