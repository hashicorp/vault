package golinters

import (
	"context"
	"fmt"

	ineffassignAPI "github.com/golangci/ineffassign"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Ineffassign struct{}

func (Ineffassign) Name() string {
	return "ineffassign"
}

func (Ineffassign) Desc() string {
	return "Detects when assignments to existing variables are not used"
}

func (lint Ineffassign) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	issues := ineffassignAPI.Run(getAllFileNames(lintCtx))
	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		res = append(res, result.Issue{
			Pos:        i.Pos,
			Text:       fmt.Sprintf("ineffectual assignment to %s", formatCode(i.IdentName, lintCtx.Cfg)),
			FromLinter: lint.Name(),
		})
	}
	return res, nil
}
