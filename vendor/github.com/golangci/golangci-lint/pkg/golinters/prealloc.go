package golinters

import (
	"context"
	"fmt"
	"go/ast"

	"github.com/golangci/prealloc"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Prealloc struct{}

func (Prealloc) Name() string {
	return "prealloc"
}

func (Prealloc) Desc() string {
	return "Finds slice declarations that could potentially be preallocated"
}

func (lint Prealloc) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var res []result.Issue

	s := &lintCtx.Settings().Prealloc
	for _, f := range lintCtx.ASTCache.GetAllValidFiles() {
		hints := prealloc.Check([]*ast.File{f.F}, s.Simple, s.RangeLoops, s.ForLoops)
		for _, hint := range hints {
			res = append(res, result.Issue{
				Pos:        f.Fset.Position(hint.Pos),
				Text:       fmt.Sprintf("Consider preallocating %s", formatCode(hint.DeclaredSliceName, lintCtx.Cfg)),
				FromLinter: lint.Name(),
			})
		}
	}

	return res, nil
}
