package golinters

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Gochecknoinits struct{}

func (Gochecknoinits) Name() string {
	return "gochecknoinits"
}

func (Gochecknoinits) Desc() string {
	return "Checks that no init functions are present in Go code"
}

func (lint Gochecknoinits) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var res []result.Issue
	for _, f := range lintCtx.ASTCache.GetAllValidFiles() {
		res = append(res, lint.checkFile(f.F, f.Fset)...)
	}

	return res, nil
}

func (lint Gochecknoinits) checkFile(f *ast.File, fset *token.FileSet) []result.Issue {
	var res []result.Issue
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		name := funcDecl.Name.Name
		if name == "init" && funcDecl.Recv.NumFields() == 0 {
			res = append(res, result.Issue{
				Pos:        fset.Position(funcDecl.Pos()),
				Text:       fmt.Sprintf("don't use %s function", formatCode(name, nil)),
				FromLinter: lint.Name(),
			})
		}
	}

	return res
}
