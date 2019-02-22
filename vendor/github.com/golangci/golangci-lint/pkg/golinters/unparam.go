package golinters

import (
	"context"

	"mvdan.cc/unparam/check"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Unparam struct{}

func (Unparam) Name() string {
	return "unparam"
}

func (Unparam) Desc() string {
	return "Reports unused function parameters"
}

func (lint Unparam) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	us := &lintCtx.Settings().Unparam

	c := &check.Checker{}
	c.CallgraphAlgorithm(us.Algo)
	c.CheckExportedFuncs(us.CheckExported)
	c.Packages(lintCtx.Packages)
	c.ProgramSSA(lintCtx.SSAProgram)

	unparamIssues, err := c.Check()
	if err != nil {
		return nil, err
	}

	var res []result.Issue
	for _, i := range unparamIssues {
		res = append(res, result.Issue{
			Pos:        lintCtx.Program.Fset.Position(i.Pos()),
			Text:       markIdentifiers(i.Message()),
			FromLinter: lint.Name(),
		})
	}

	return res, nil
}
