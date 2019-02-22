package golinters

import (
	"context"

	"mvdan.cc/interfacer/check"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Interfacer struct{}

func (Interfacer) Name() string {
	return "interfacer"
}

func (Interfacer) Desc() string {
	return "Linter that suggests narrower interface types"
}

func (lint Interfacer) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	c := new(check.Checker)
	c.Program(lintCtx.Program)
	c.ProgramSSA(lintCtx.SSAProgram)

	issues, err := c.Check()
	if err != nil {
		return nil, err
	}
	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		pos := lintCtx.SSAProgram.Fset.Position(i.Pos())
		res = append(res, result.Issue{
			Pos:        pos,
			Text:       markIdentifiers(i.Message()),
			FromLinter: lint.Name(),
		})
	}

	return res, nil
}
