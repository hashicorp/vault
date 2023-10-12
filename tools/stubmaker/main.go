// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-hclog"
	"golang.org/x/tools/go/packages"
)

var logger hclog.Logger

func fatal(err error) {
	logger.Error("fatal error", "error", err)
	os.Exit(1)
}

func main() {
	logger = hclog.New(&hclog.LoggerOptions{
		Name:  "stubmaker",
		Level: hclog.Trace,
	})

	// Setup git, both so we can determine if we're running on enterprise, and
	// so we can make sure we don't clobber a non-transient file.
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		fatal(err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		fatal(err)
	}
	if !isEnterprise(wt) {
		return
	}

	// Read the file and figure out if we need to do anything.
	inputFile := os.Getenv("GOFILE")
	if !strings.HasSuffix(inputFile, "_stubs_oss.go") {
		fatal(fmt.Errorf("stubmaker should only be invoked from files ending in _stubs_oss.go"))
	}

	baseFilename := strings.TrimSuffix(inputFile, "_stubs_oss.go")
	outputFile := baseFilename + "_stubs_ent.go"
	b, err := os.ReadFile(inputFile)
	if err != nil {
		fatal(err)
	}

	inputLines, err := readLines(bytes.NewBuffer(b))
	if err != nil {
		fatal(err)
	}
	funcs := getFuncs(inputLines)
	if needed, err := isStubNeeded(funcs); err != nil {
		fatal(err)
	} else if !needed {
		return
	}

	// We'd like to write the file, but first make sure that we're not going
	// to blow away anyone's work or overwrite a file already in git.
	head, err := repo.Head()
	if err != nil {
		fatal(err)
	}
	obj, err := repo.Object(plumbing.AnyObject, head.Hash())
	if err != nil {
		fatal(err)
	}

	st, err := wt.Status()
	if err != nil {
		fatal(err)
	}

	tracked, err := inGit(wt, st, obj, outputFile)
	if err != nil {
		fatal(err)
	}
	if tracked {
		fatal(fmt.Errorf("output file %s exists in git, not overwriting", outputFile))
	}

	// Now we can finally write the file
	output, err := os.Create(outputFile + ".tmp")
	if err != nil {
		fatal(err)
	}
	_, err = io.WriteString(output, strings.Join(getOutput(inputLines), "\n")+"\n")
	if err != nil {
		// If we don't end up writing to the file, delete it.
		os.Remove(outputFile + ".tmp")
	} else {
		os.Rename(outputFile+".tmp", outputFile)
	}
	if err != nil {
		fatal(err)
	}
}

func inGit(wt *git.Worktree, st git.Status, obj object.Object, path string) (bool, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, fmt.Errorf("path %s can't be made absolute: %w", path, err)
	}
	relPath, err := filepath.Rel(wt.Filesystem.Root(), absPath)
	if err != nil {
		return false, fmt.Errorf("path %s can't be made relative: %w", absPath, err)
	}

	fst := st.File(relPath)
	if fst.Worktree != git.Untracked || fst.Staging != git.Untracked {
		return true, nil
	}

	curwd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	blob, err := resolve(obj, relPath)
	if err != nil && !strings.Contains(err.Error(), "file not found") {
		return false, fmt.Errorf("error resolving path %s from %s: %w", relPath, curwd, err)
	}

	return blob != nil, nil
}

func isEnterprise(wt *git.Worktree) bool {
	st, err := wt.Filesystem.Stat("enthelpers")
	onOss := errors.Is(err, os.ErrNotExist)
	onEnt := st != nil

	switch {
	case onOss && !onEnt:
	case !onOss && onEnt:
	default:
		fatal(err)
	}
	return onEnt
}

// resolve blob at given path from obj. obj can be a commit, tag, tree, or blob.
func resolve(obj object.Object, path string) (*object.Blob, error) {
	switch o := obj.(type) {
	case *object.Commit:
		t, err := o.Tree()
		if err != nil {
			return nil, err
		}
		return resolve(t, path)
	case *object.Tag:
		target, err := o.Object()
		if err != nil {
			return nil, err
		}
		return resolve(target, path)
	case *object.Tree:
		file, err := o.File(path)
		if err != nil {
			return nil, err
		}
		return &file.Blob, nil
	case *object.Blob:
		return o, nil
	default:
		return nil, object.ErrUnsupportedObject
	}
}

func readLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func isStubNeeded(funcs []string) (bool, error) {
	pkg, err := parsePackage(".", []string{"enterprise"})
	if err != nil {
		return false, err
	}

	var found []string
	for name, val := range pkg.TypesInfo.Defs {
		if val == nil {
			continue
		}
		_, ok := val.Type().(*types.Signature)
		if !ok {
			continue
		}
		for _, f := range funcs {
			if name.Name == f {
				found = append(found, f)
			}
		}
	}
	switch {
	case len(found) == len(funcs):
		return false, nil
	case len(found) != 0:
		sort.Strings(found)
		sort.Strings(funcs)
		delta := cmp.Diff(found, funcs)
		return false, fmt.Errorf("funcs partially defined, delta=%s", delta)
	}

	return true, nil
}

var funcRE = regexp.MustCompile("^func *(?:[(][^)]+[)])? *([^(]+)")

func getFuncs(inputLines []string) []string {
	var funcs []string
	for _, line := range inputLines {
		matches := funcRE.FindStringSubmatch(line)
		if len(matches) > 1 {
			funcs = append(funcs, matches[1])
		}
	}
	return funcs
}

func getOutput(inputLines []string) []string {
	warning := "// Code generated by tools/stubmaker; DO NOT EDIT."

	var outputLines []string
	for _, line := range inputLines {
		switch line {
		case "//go:build !enterprise":
			outputLines = append(outputLines, warning, "")
			line = "//go:build enterprise"
		case "//go:generate go run github.com/hashicorp/vault/tools/stubmaker":
			continue
		}
		outputLines = append(outputLines, line)
	}

	return outputLines
}

func parsePackage(name string, tags []string) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, name)
	if err != nil {
		return nil, fmt.Errorf("error parsing package %s: %v", name, err)
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("error: %d packages found", len(pkgs))
	}
	return pkgs[0], nil
}
