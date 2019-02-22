package golinters

import (
	"context"
	"fmt"
	"strings"

	depguardAPI "github.com/OpenPeeDeeP/depguard"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Depguard struct{}

func (Depguard) Name() string {
	return "depguard"
}

func (Depguard) Desc() string {
	return "Go linter that checks if package imports are in a list of acceptable packages"
}

func (d Depguard) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	dg := &depguardAPI.Depguard{
		Packages:      lintCtx.Settings().Depguard.Packages,
		IncludeGoRoot: lintCtx.Settings().Depguard.IncludeGoRoot,
	}
	listType := lintCtx.Settings().Depguard.ListType
	var found bool
	dg.ListType, found = depguardAPI.StringToListType[strings.ToLower(listType)]
	if !found {
		if listType != "" {
			return nil, fmt.Errorf("unsure what list type %s is", listType)
		}
		dg.ListType = depguardAPI.LTBlacklist
	}

	issues, err := dg.Run(lintCtx.LoaderConfig, lintCtx.Program)
	if err != nil {
		return nil, err
	}
	if len(issues) == 0 {
		return nil, nil
	}
	msgSuffix := "is in the blacklist"
	if dg.ListType == depguardAPI.LTWhitelist {
		msgSuffix = "is not in the whitelist"
	}
	res := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		res = append(res, result.Issue{
			Pos:        i.Position,
			Text:       fmt.Sprintf("%s %s", formatCode(i.PackageName, lintCtx.Cfg), msgSuffix),
			FromLinter: d.Name(),
		})
	}
	return res, nil
}
