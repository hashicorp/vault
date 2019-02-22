package golinters

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Gochecknoglobals struct{}

func (Gochecknoglobals) Name() string {
	return "gochecknoglobals"
}

func (Gochecknoglobals) Desc() string {
	return "Checks that no globals are present in Go code"
}

func (lint Gochecknoglobals) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var res []result.Issue
	for _, f := range lintCtx.ASTCache.GetAllValidFiles() {
		res = append(res, lint.checkFile(f.F, f.Fset)...)
	}

	return res, nil
}

func (lint Gochecknoglobals) checkFile(f *ast.File, fset *token.FileSet) []result.Issue {
	var res []result.Issue
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)
			for _, vn := range valueSpec.Names {
				if isWhitelisted(vn) {
					continue
				}

				res = append(res, result.Issue{
					Pos:        fset.Position(genDecl.TokPos),
					Text:       fmt.Sprintf("%s is a global variable", formatCode(vn.Name, nil)),
					FromLinter: lint.Name(),
				})
			}
		}
	}

	return res
}

func isWhitelisted(i *ast.Ident) bool {
	return i.Name == "_" || looksLikeError(i)
}

// looksLikeError returns true if the AST identifier starts
// with 'err' or 'Err', or false otherwise.
//
// TODO: https://github.com/leighmcculloch/gochecknoglobals/issues/5
func looksLikeError(i *ast.Ident) bool {
	prefix := "err"
	if i.IsExported() {
		prefix = "Err"
	}
	return strings.HasPrefix(i.Name, prefix)
}
