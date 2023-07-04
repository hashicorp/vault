package main

import (
	"bufio"
	"errors"
	"fmt"
	"go/types"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"golang.org/x/tools/go/packages"
)

func main() {
	baseFilename := strings.TrimSuffix(os.Getenv("GOFILE"), ".go")
	packageName := os.Getenv("GOPACKAGE")
	stub := os.Args[1]
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

	st, err := wt.Filesystem.Stat("enthelpers")
	onOss := errors.Is(err, os.ErrNotExist)
	onEnt := st != nil
	var repoString string
	var tags []string
	switch {
	case onOss && !onEnt:
		repoString, tags = "oss", nil
	case !onOss && onEnt:
		repoString, tags = "ent", []string{"enterprise"}
	default:
		fatal(err)
	}
	pkg, err := parsePackage(".", tags)
	if err != nil {
		fatal(err)
	}
	//for k, v := range pkg.TypesInfo.Defs {
	//	if v != nil {
	//		_, ok := v.Type().(*types.Signature)
	//		if ok {
	//			fmt.Println(k, v.String())
	//		}
	//	}
	//}
	target := fmt.Sprintf("%s_%s_stubs.go", baseFilename, repoString)
	appendStub(target, packageName, stub, onOss, pkg)
	// fmt.Println(target, packageName, stub, wt.Filesystem.Root(), onOss, onEnt)
}

func appendStub(target string, packageName string, stub string, onOss bool, pkg *packages.Package) error {
	_, err := os.Stat(target)
	var f *os.File
	var lines []string
	if err == nil {
		f, err = os.OpenFile(target, os.O_RDWR, 0o644)
		if err != nil {
			fatal(err)
		}
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	} else {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		f, err = os.Create(target)
		if err != nil {
			return err
		}
		defer func() {
			st, _ := f.Stat()
			if st.Size() == 0 {
				os.Remove(target)
			}
		}()
	}

	// TODO add edit warning
	var add []string
	if onOss {
		add = append(add, "//go:build !enterprise")
	} else {
		parenPos := strings.Index(stub, "(")
		if parenPos == -1 {
			fatal(fmt.Errorf("no paren found"))
		}
		funcName := strings.TrimSpace(stub[:parenPos])
		for name, val := range pkg.TypesInfo.Defs {
			if val == nil {
				continue
			}
			_, ok := val.Type().(*types.Signature)
			if ok && name.Name == funcName {
				//if val.String() != "func "+stub {
				//	fatal(fmt.Errorf("found existing def with different sig: %s", val.String()))
				//}
				return nil
			}
		}
	}
	add = append(add, fmt.Sprintf("package %s", packageName))
	add = append(add, "func "+stub)

	for _, line := range add {
		if stringPresent(lines, line) {
			continue
		}
		_, err = f.WriteString(line + "\n\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func stringPresent(haystack []string, needle string) bool {
	for _, line := range haystack {
		if strings.Contains(line, needle) {
			return true
		}
	}
	return false
}

func fatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func parsePackage(name string, tags []string) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, name)
	if err != nil {
		return nil, err
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("error: %d packages found", len(pkgs))
	}
	return pkgs[0], nil
}
