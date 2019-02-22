package golinters

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Nakedret struct{}

func (Nakedret) Name() string {
	return "nakedret"
}

func (Nakedret) Desc() string {
	return "Finds naked returns in functions greater than a specified function length"
}

type nakedretVisitor struct {
	maxLength int
	f         *token.FileSet
	issues    []result.Issue
}

func (v *nakedretVisitor) processFuncDecl(funcDecl *ast.FuncDecl) {
	file := v.f.File(funcDecl.Pos())
	functionLineLength := file.Position(funcDecl.End()).Line - file.Position(funcDecl.Pos()).Line

	// Scan the body for usage of the named returns
	for _, stmt := range funcDecl.Body.List {
		s, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			continue
		}

		if len(s.Results) != 0 {
			continue
		}

		file := v.f.File(s.Pos())
		if file == nil || functionLineLength <= v.maxLength {
			continue
		}
		if funcDecl.Name == nil {
			continue
		}

		v.issues = append(v.issues, result.Issue{
			FromLinter: Nakedret{}.Name(),
			Text: fmt.Sprintf("naked return in func `%s` with %d lines of code",
				funcDecl.Name.Name, functionLineLength),
			Pos: v.f.Position(s.Pos()),
		})
	}
}

func (v *nakedretVisitor) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return v
	}

	var namedReturns []*ast.Ident

	// We've found a function
	if funcDecl.Type != nil && funcDecl.Type.Results != nil {
		for _, field := range funcDecl.Type.Results.List {
			for _, ident := range field.Names {
				if ident != nil {
					namedReturns = append(namedReturns, ident)
				}
			}
		}
	}

	if len(namedReturns) == 0 || funcDecl.Body == nil {
		return v
	}

	v.processFuncDecl(funcDecl)
	return v
}

func (lint Nakedret) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var res []result.Issue
	for _, f := range lintCtx.ASTCache.GetAllValidFiles() {
		v := nakedretVisitor{
			maxLength: lintCtx.Settings().Nakedret.MaxFuncLines,
			f:         f.Fset,
		}
		ast.Walk(&v, f.F)
		res = append(res, v.issues...)
	}

	return res, nil
}
