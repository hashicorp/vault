package golinters

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
	"sync"

	gopackages "golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

func formatCode(code string, _ *config.Config) string {
	if strings.Contains(code, "`") {
		return code // TODO: properly escape or remove
	}

	return fmt.Sprintf("`%s`", code)
}

func formatCodeBlock(code string, _ *config.Config) string {
	if strings.Contains(code, "`") {
		return code // TODO: properly escape or remove
	}

	return fmt.Sprintf("```\n%s\n```", code)
}

type replacePattern struct {
	re   string
	repl string
}

type replaceRegexp struct {
	re   *regexp.Regexp
	repl string
}

var replaceRegexps []replaceRegexp
var replaceRegexpsOnce sync.Once

var replacePatterns = []replacePattern{
	// unparam
	{`^(\S+) - (\S+) is unused$`, "`${1}` - `${2}` is unused"},
	{`^(\S+) - (\S+) always receives (\S+) \((.*)\)$`, "`${1}` - `${2}` always receives `${3}` (`${4}`)"},
	{`^(\S+) - (\S+) always receives (.*)$`, "`${1}` - `${2}` always receives `${3}`"},

	// interfacer
	{`^(\S+) can be (\S+)$`, "`${1}` can be `${2}`"},

	// govet
	{`^(\S+) arg list ends with redundant newline$`, "`${1}` arg list ends with redundant newline"},
	{`^(\S+) composite literal uses unkeyed fields$`, "`${1}` composite literal uses unkeyed fields"},

	// gosec
	{`^Blacklisted import (\S+): weak cryptographic primitive$`,
		"Blacklisted import `${1}`: weak cryptographic primitive"},
	{`^TLS InsecureSkipVerify set true.$`, "TLS `InsecureSkipVerify` set true."},

	// gosimple
	{`^should replace loop with (.*)$`, "should replace loop with `${1}`"},

	// megacheck
	{`^this value of (\S+) is never used$`, "this value of `${1}` is never used"},
	{`^should use time.Since instead of time.Now().Sub$`,
		"should use `time.Since` instead of `time.Now().Sub`"},
	{`^(func|const|field|type) (\S+) is unused$`, "${1} `${2}` is unused"},
}

func markIdentifiers(s string) string {
	replaceRegexpsOnce.Do(func() {
		for _, p := range replacePatterns {
			r := replaceRegexp{
				re:   regexp.MustCompile(p.re),
				repl: p.repl,
			}
			replaceRegexps = append(replaceRegexps, r)
		}
	})

	for _, rr := range replaceRegexps {
		rs := rr.re.ReplaceAllString(s, rr.repl)
		if rs != s {
			return rs
		}
	}

	return s
}

func getAllFileNames(ctx *linter.Context) []string {
	var ret []string
	uniqFiles := map[string]bool{} // files are duplicated for test packages
	for _, pkg := range ctx.Packages {
		for _, f := range pkg.GoFiles {
			if uniqFiles[f] {
				continue
			}
			uniqFiles[f] = true
			ret = append(ret, f)
		}
	}
	return ret
}

func getASTFilesForGoPkg(ctx *linter.Context, pkg *gopackages.Package) ([]*ast.File, *token.FileSet, error) {
	var files []*ast.File
	var fset *token.FileSet
	for _, filename := range pkg.GoFiles {
		f := ctx.ASTCache.Get(filename)
		if f == nil {
			return nil, nil, fmt.Errorf("no AST for file %s in cache: %+v", filename, *ctx.ASTCache)
		}

		if f.Err != nil {
			return nil, nil, fmt.Errorf("can't load AST for file %s: %s", f.Name, f.Err)
		}

		files = append(files, f.F)
		fset = f.Fset
	}

	return files, fset, nil
}
