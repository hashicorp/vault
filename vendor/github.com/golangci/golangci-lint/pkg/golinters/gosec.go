package golinters

import (
	"context"
	"fmt"
	"go/token"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/golangci/gosec"
	"github.com/golangci/gosec/rules"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Gosec struct{}

func (Gosec) Name() string {
	return "gosec"
}

func (Gosec) Desc() string {
	return "Inspects source code for security problems"
}

func (lint Gosec) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	gasConfig := gosec.NewConfig()
	enabledRules := rules.Generate()
	logger := log.New(ioutil.Discard, "", 0)
	analyzer := gosec.NewAnalyzer(gasConfig, logger)
	analyzer.LoadRules(enabledRules.Builders())

	analyzer.ProcessProgram(lintCtx.Program)
	issues, _ := analyzer.Report()
	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		text := fmt.Sprintf("%s: %s", i.RuleID, markIdentifiers(i.What)) // TODO: use severity and confidence
		var r *result.Range
		line, err := strconv.Atoi(i.Line)
		if err != nil {
			r = &result.Range{}
			if n, rerr := fmt.Sscanf(i.Line, "%d-%d", &r.From, &r.To); rerr != nil || n != 2 {
				lintCtx.Log.Warnf("Can't convert gosec line number %q of %v to int: %s", i.Line, i, err)
				continue
			}
			line = r.From
		}

		res = append(res, result.Issue{
			Pos: token.Position{
				Filename: i.File,
				Line:     line,
			},
			Text:       text,
			LineRange:  r,
			FromLinter: lint.Name(),
		})
	}

	return res, nil
}
