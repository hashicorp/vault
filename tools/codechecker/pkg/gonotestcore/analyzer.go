// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package gonotestcore

import (
	"go/ast"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// testCoreFamily matches the TestCore family of test helpers (TestCore,
// TestCoreUnsealed, TestCoreUnsealedWithConfigs, TestCoreWithConfig, ...). It
// intentionally does not match the actual test functions, which use an
// underscore after "TestCore" (e.g. TestCore_Foo), because those are defined
// with, not calls to, the helpers.
var testCoreFamily = regexp.MustCompile(`^TestCore([A-Z].*)?$`)

// bannedWrappers are in-package helper wrappers that call a TestCore* helper in
// their body (directly or transitively). They are flagged by name so the
// report lands on the call site (next to the test) rather than on the far-away
// helper definition.
var bannedWrappers = map[string]struct{}{
	"mockExpiration":                          {},
	"mockPolicyWithCore":                      {},
	"mockRollback":                            {},
	"testActivationFlags_Write_Activate":      {},
	"testActivationFlags_Write_Deactivate":    {},
	"testClusterWithContainerPlugins":         {},
	"testCluster_ForwardRequestsCommon":       {},
	"testControlGroupCore":                    {},
	"testCoreUnsealed":                        {},
	"testCoreWithIdentityTokenGithub":         {},
	"testCoreWithIdentityTokenGithubRoot":     {},
	"testCoreWithPlugins":                     {},
	"testCore_Rekey_Update_Common":            {},
	"testCore_Rekey_Update_Common_Error":      {},
	"testCore_Standby_Common":                 {},
	"testCore_Unmount_Cleanup":                {},
	"testIdentityStoreWithGithubAuth":         {},
	"testIdentityStoreWithGithubAuthRoot":     {},
	"testIdentityStoreWithGithubUserpassAuth": {},
	"testIdentityStoreWithLocalGithubAuth":    {},
	"testNewClientsInternal":                  {},
	"testPoliciesRecoverSourcePath":           {},
	"testTokenStoreHandleRequestLookup":       {},
	"testTokenStore_RevokeTree_NonRecursive":  {},
}

// message explains why the helper is disallowed and what to do instead. The
// leading %s is the name of the offending helper.
const message = "%s is part of the TestCore family of test helpers (or an in-package wrapper around it) and this test should **not** be added as-is. " +
	"Instead, use TestCluster based tests outside of the vault package: they can be parallelized, they keep the " +
	"vault test package from growing, and they are easier to read because they use our own SDK. " +
	"Add the test to a file in the external_tests directory and use vault.NewTestCluster or minimal.NewTestSoloCluster. " +
	"If you are struggling to adapt your test, reach out to the Core team in #team-vault-core."

var Analyzer = &analysis.Analyzer{
	Name:     "gonotestcore",
	Doc:      "Flags usage of the TestCore family of test helpers (and their in-package wrappers) in favor of TestCluster based tests outside the vault package",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// ignore-nil-nil-function-check
func run(pass *analysis.Pass) (interface{}, error) {
	inspctr := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspctr.Preorder(nodeFilter, func(node ast.Node) {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return
		}

		name := calledFuncName(call.Fun)
		if name == "" {
			return
		}

		if isBanned(name) {
			pass.Reportf(call.Pos(), message, name)
		}
	})
	return nil, nil
}

// calledFuncName returns the identifier name of a called function, handling
// both unqualified calls (TestCoreUnsealed(t)) and qualified calls from another
// package (vault.TestCoreUnsealed(t)). It returns "" when the callee is not a
// simple function or method reference.
func calledFuncName(fun ast.Expr) string {
	switch f := fun.(type) {
	case *ast.Ident:
		return f.Name
	case *ast.SelectorExpr:
		return f.Sel.Name
	default:
		return ""
	}
}

// isBanned reports whether a called function name is a TestCore family helper
// or one of the curated in-package wrappers around them.
func isBanned(name string) bool {
	if _, ok := bannedWrappers[name]; ok {
		return true
	}
	return testCoreFamily.MatchString(name)
}
