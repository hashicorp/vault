package processors

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type SkipDirs struct {
	patterns    []*regexp.Regexp
	log         logutils.Log
	skippedDirs map[string][]string // regexp to dir mapping
	absArgsDirs []string
}

var _ Processor = SkipFiles{}

const goFileSuffix = ".go"

func NewSkipDirs(patterns []string, log logutils.Log, runArgs []string) (*SkipDirs, error) {
	var patternsRe []*regexp.Regexp
	for _, p := range patterns {
		patternRe, err := regexp.Compile(p)
		if err != nil {
			return nil, errors.Wrapf(err, "can't compile regexp %q", p)
		}
		patternsRe = append(patternsRe, patternRe)
	}

	if len(runArgs) == 0 {
		runArgs = append(runArgs, "./...")
	}
	var absArgsDirs []string
	for _, arg := range runArgs {
		base := filepath.Base(arg)
		if base == "..." || strings.HasSuffix(base, goFileSuffix) {
			arg = filepath.Dir(arg)
		}

		absArg, err := filepath.Abs(arg)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to abs-ify arg %q", arg)
		}
		absArgsDirs = append(absArgsDirs, absArg)
	}

	return &SkipDirs{
		patterns:    patternsRe,
		log:         log,
		skippedDirs: map[string][]string{},
		absArgsDirs: absArgsDirs,
	}, nil
}

func (p SkipDirs) Name() string {
	return "skip_dirs"
}

func (p *SkipDirs) Process(issues []result.Issue) ([]result.Issue, error) {
	if len(p.patterns) == 0 {
		return issues, nil
	}

	return filterIssues(issues, p.shouldPassIssue), nil
}

func (p *SkipDirs) shouldPassIssue(i *result.Issue) bool {
	if filepath.IsAbs(i.FilePath()) {
		p.log.Warnf("Got abs path in skip dirs processor, it should be relative")
		return true
	}

	issueRelDir := filepath.Dir(i.FilePath())
	issueAbsDir, err := filepath.Abs(issueRelDir)
	if err != nil {
		p.log.Warnf("Can't abs-ify path %q: %s", issueRelDir, err)
		return true
	}

	for _, absArgDir := range p.absArgsDirs {
		if absArgDir == issueAbsDir {
			// we must not skip issues if they are from explicitly set dirs
			// even if they match skip patterns
			return true
		}
	}

	// We use issueRelDir for matching: it's the relative to the current
	// work dir path of directory of source file with the issue. It can lead
	// to unexpected behavior if we're analyzing files out of current work dir.
	// The alternative solution is to find relative to args path, but it has
	// disadvantages (https://github.com/golangci/golangci-lint/pull/313).

	for _, pattern := range p.patterns {
		if pattern.MatchString(issueRelDir) {
			ps := pattern.String()
			p.skippedDirs[ps] = append(p.skippedDirs[ps], issueRelDir)
			return false
		}
	}

	return true
}

func (p SkipDirs) Finish() {
	for pattern, dirs := range p.skippedDirs {
		p.log.Infof("Skipped by pattern %s dirs: %s", pattern, dirs)
	}
}
