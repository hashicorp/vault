package golinters

import (
	"context"
	"fmt"

	malignedAPI "github.com/golangci/maligned"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Maligned struct{}

func (Maligned) Name() string {
	return "maligned"
}

func (Maligned) Desc() string {
	return "Tool to detect Go structs that would take less memory if their fields were sorted"
}

func (m Maligned) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	issues := malignedAPI.Run(lintCtx.Program)
	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		text := fmt.Sprintf("struct of size %d bytes could be of size %d bytes", i.OldSize, i.NewSize)
		if lintCtx.Settings().Maligned.SuggestNewOrder {
			text += fmt.Sprintf(":\n%s", formatCodeBlock(i.NewStructDef, lintCtx.Cfg))
		}
		res = append(res, result.Issue{
			Pos:        i.Pos,
			Text:       text,
			FromLinter: m.Name(),
		})
	}
	return res, nil
}
