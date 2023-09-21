// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package gonilnilfunctions

import (
	"go/ast"
	"go/types"
	"reflect"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:       "gonilnilfunctions",
	Doc:        "Verifies that every go function with error as one of its two return types cannot return nil, nil",
	Run:        run,
	ResultType: reflect.TypeOf((interface{})(nil)),
	Requires:   []*analysis.Analyzer{inspect.Analyzer},
}

// getNestedReturnStatements searches the AST for return statements, and returns
// them in a tail-call optimized list.
func getNestedReturnStatements(s ast.Stmt, returns []*ast.ReturnStmt) []*ast.ReturnStmt {
	switch s := s.(type) {
	case *ast.BlockStmt:
		statements := make([]*ast.ReturnStmt, 0)
		for _, stmt := range s.List {
			statements = append(statements, getNestedReturnStatements(stmt, make([]*ast.ReturnStmt, 0))...)
		}

		return append(returns, statements...)
	case *ast.BranchStmt:
		return returns
	case *ast.ForStmt:
		return getNestedReturnStatements(s.Body, returns)
	case *ast.IfStmt:
		return getNestedReturnStatements(s.Body, returns)
	case *ast.LabeledStmt:
		return getNestedReturnStatements(s.Stmt, returns)
	case *ast.RangeStmt:
		return getNestedReturnStatements(s.Body, returns)
	case *ast.ReturnStmt:
		return append(returns, s)
	case *ast.SwitchStmt:
		return getNestedReturnStatements(s.Body, returns)
	case *ast.SelectStmt:
		return getNestedReturnStatements(s.Body, returns)
	case *ast.TypeSwitchStmt:
		return getNestedReturnStatements(s.Body, returns)
	case *ast.CommClause:
		statements := make([]*ast.ReturnStmt, 0)
		for _, stmt := range s.Body {
			statements = append(statements, getNestedReturnStatements(stmt, make([]*ast.ReturnStmt, 0))...)
		}

		return append(returns, statements...)
	case *ast.CaseClause:
		statements := make([]*ast.ReturnStmt, 0)
		for _, stmt := range s.Body {
			statements = append(statements, getNestedReturnStatements(stmt, make([]*ast.ReturnStmt, 0))...)
		}

		return append(returns, statements...)
	case *ast.ExprStmt:
		return returns
	}
	return returns
}

// run runs the analysis, failing for functions whose signatures contain two results including one error
// (e.g. (something, error)), that contain multiple nil returns
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

		// If the function has the "Ignore" godoc comment, skip it
		if strings.Contains(funcDecl.Doc.Text(), "ignore-nil-nil-function-check") {
			return
		}

		// The function returns something
		if funcDecl == nil || funcDecl.Type == nil || funcDecl.Type.Results == nil {
			return
		}

		// The function has more than 1 return value
		results := funcDecl.Type.Results.List
		if len(results) < 2 {
			return
		}

		// isError is a helper function to check if a Field is of error type
		isError := func(field *ast.Field) bool {
			if named, ok := pass.TypesInfo.TypeOf(field.Type).(*types.Named); ok {
				namedObject := named.Obj()
				return namedObject != nil && namedObject.Pkg() == nil && namedObject.Name() == "error"
			}
			return false
		}

		// one of the return values is error
		var errorFound bool
		for _, result := range results {
			if isError(result) {
				errorFound = true
				break
			}
		}

		if !errorFound {
			return
		}

		// Since these statements might be e.g. blocks with
		// other statements inside, we need to get the return statements
		// from inside them, first.
		statements := funcDecl.Body.List

		returnStatements := make([]*ast.ReturnStmt, 0)
		for _, statement := range statements {
			returnStatements = append(returnStatements, getNestedReturnStatements(statement, make([]*ast.ReturnStmt, 0))...)
		}

		for _, returnStatement := range returnStatements {
			numResultsNil := 0
			results := returnStatement.Results

			// We only want two-arg functions (something, nil)
			// We can remove this block in the future if we change our mind
			if len(results) != 2 {
				continue
			}

			for _, result := range results {
				// nil is an ident
				ident, isIdent := result.(*ast.Ident)
				if isIdent {
					if ident.Name == "nil" {
						// We found one nil in the return list
						numResultsNil++
					}
				}
			}
			// We found N nils, and our function returns N results, so this fails the check
			if numResultsNil == len(results) {
				// All the return values are nil, so we fail the report
				pass.Reportf(node.Pos(), "Function %s can return an error, and has a statement that returns only nils",
					funcDecl.Name.Name)

				// We break out of the loop of checking return statements, so that we don't repeat ourselves
				break
			}
		}
	})

	var success interface{}
	return success, nil
}
