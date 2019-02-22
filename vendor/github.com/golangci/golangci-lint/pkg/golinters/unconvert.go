package golinters

import (
	"context"

	unconvertAPI "github.com/golangci/unconvert"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Unconvert struct{}

func (Unconvert) Name() string {
	return "unconvert"
}

func (Unconvert) Desc() string {
	return "Remove unnecessary type conversions"
}

func (lint Unconvert) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	positions := unconvertAPI.Run(lintCtx.Program)
	if len(positions) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(positions))
	for _, pos := range positions {
		res = append(res, result.Issue{
			Pos:        pos,
			Text:       "unnecessary conversion",
			FromLinter: lint.Name(),
		})
	}

	return res, nil
}
