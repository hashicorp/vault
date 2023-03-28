// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "godoctests",
	Doc:      "Verifies that every go test has a go doc",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return
		}

		// starts with 'Test'
		if !strings.HasPrefix(funcDecl.Name.Name, "Test") {
			return
		}

		// has one parameter
		params := funcDecl.Type.Params.List
		if len(params) != 1 {
			return
		}

		// parameter is a pointer
		firstParamType, ok := params[0].Type.(*ast.StarExpr)
		if !ok {
			return
		}

		selector, ok := firstParamType.X.(*ast.SelectorExpr)
		if !ok {
			return
		}

		// the pointer comes from package 'testing'
		selectorIdent, ok := selector.X.(*ast.Ident)
		if !ok {
			return
		}
		if selectorIdent.Name != "testing" {
			return
		}

		// the pointer has type 'T'
		if selector.Sel == nil || selector.Sel.Name != "T" {
			return
		}

		// then there must be a godoc
		if funcDecl.Doc == nil {
			pass.Reportf(node.Pos(), "Test %s is missing a go doc",
				funcDecl.Name.Name)
		} else if !strings.HasPrefix(funcDecl.Doc.Text(), funcDecl.Name.Name) {
			pass.Reportf(node.Pos(), "Test %s must have a go doc beginning with the function name",
				funcDecl.Name.Name)
		}
	})
	return nil, nil
}
