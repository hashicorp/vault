package golinters

import (
	"context"
	"fmt"
	"go/token"
	"io/ioutil"
	"strings"

	"github.com/golangci/misspell"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Misspell struct{}

func (Misspell) Name() string {
	return "misspell"
}

func (Misspell) Desc() string {
	return "Finds commonly misspelled English words in comments"
}

func (lint Misspell) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	r := misspell.Replacer{
		Replacements: misspell.DictMain,
	}

	// Figure out regional variations
	locale := lintCtx.Settings().Misspell.Locale
	switch strings.ToUpper(locale) {
	case "":
		// nothing
	case "US":
		r.AddRuleList(misspell.DictAmerican)
	case "UK", "GB":
		r.AddRuleList(misspell.DictBritish)
	case "NZ", "AU", "CA":
		return nil, fmt.Errorf("unknown locale: %q", locale)
	}

	r.Compile()

	var res []result.Issue
	for _, f := range getAllFileNames(lintCtx) {
		fileContent, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("can't read file %s: %s", f, err)
		}

		_, diffs := r.ReplaceGo(string(fileContent))
		for _, diff := range diffs {
			text := fmt.Sprintf("`%s` is a misspelling of `%s`", diff.Original, diff.Corrected)
			pos := token.Position{
				Filename: f,
				Line:     diff.Line,
				Column:   diff.Column + 1,
			}
			res = append(res, result.Issue{
				Pos:        pos,
				Text:       text,
				FromLinter: lint.Name(),
			})
		}
	}

	return res, nil
}
