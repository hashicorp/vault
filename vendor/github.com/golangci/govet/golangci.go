package govet

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/loader"
)

type Issue struct {
	Pos     token.Position
	Message string
}

var foundIssues []Issue

func Analyze(files []*ast.File, fset *token.FileSet, pkgInfo *loader.PackageInfo, checkShadowing bool, pg astFilePathGetter) ([]Issue, error) {
	foundIssues = nil
	*source = false // import type data for "fmt" from installed packages

	if checkShadowing {
		experimental["shadow"] = false
	}
	for name, setting := range report {
		if *setting == unset && !experimental[name] {
			*setting = setTrue
		}
	}

	initPrintFlags()
	initUnusedFlags()

	filesRun = true
	for _, f := range files {
		name := fset.Position(f.Pos()).Filename
		if !strings.HasSuffix(name, "_test.go") {
			includesNonTest = true
		}
	}
	pkg, err := doPackage(nil, pkgInfo, fset, files, pg)
	if err != nil {
		return nil, err
	}

	if pkg == nil {
		return nil, nil
	}

	return foundIssues, nil
}
