package golinters

import (
	"context"

	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	libpackages "github.com/golangci/golangci-lint/pkg/packages"
	"github.com/golangci/golangci-lint/pkg/result"
)

type TypeCheck struct{}

func (TypeCheck) Name() string {
	return "typecheck"
}

func (TypeCheck) Desc() string {
	return "Like the front-end of a Go compiler, parses and type-checks Go code"
}

func (lint TypeCheck) parseError(srcErr packages.Error) (*result.Issue, error) {
	pos, err := libpackages.ParseErrorPosition(srcErr.Pos)
	if err != nil {
		return nil, err
	}

	return &result.Issue{
		Pos:        *pos,
		Text:       srcErr.Msg,
		FromLinter: lint.Name(),
	}, nil
}

func (lint TypeCheck) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	uniqReportedIssues := map[string]bool{}

	var res []result.Issue
	for _, pkg := range lintCtx.NotCompilingPackages {
		errors := libpackages.ExtractErrors(pkg, lintCtx.ASTCache)
		for _, err := range errors {
			i, perr := lint.parseError(err)
			if perr != nil { // failed to parse
				if uniqReportedIssues[err.Msg] {
					continue
				}
				uniqReportedIssues[err.Msg] = true
				lintCtx.Log.Errorf("typechecking error: %s", err.Msg)
			} else {
				res = append(res, *i)
			}
		}
	}

	return res, nil
}
