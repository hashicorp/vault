package depguard

import (
	"go/build"
	"go/token"
	"os"
	"sort"
	"strings"

	"github.com/gobwas/glob"
	"golang.org/x/tools/go/loader"
)

//ListType states what kind of list is passed in.
type ListType int

const (
	//LTBlacklist states the list given is a blacklist. (default)
	LTBlacklist ListType = iota
	//LTWhitelist states the list given is a whitelist.
	LTWhitelist
)

//StringToListType makes it easier to turn a string into a ListType.
//It assumes that the string representation is lower case.
var StringToListType = map[string]ListType{
	"whitelist": LTWhitelist,
	"blacklist": LTBlacklist,
}

//Issue with the package with PackageName at the Position.
type Issue struct {
	PackageName string
	Position    token.Position
}

//Depguard checks imports to make sure they follow the given list and constraints.
type Depguard struct {
	ListType       ListType
	Packages       []string
	IncludeGoRoot  bool
	prefixPackages []string
	globPackages   []glob.Glob
	buildCtx       *build.Context
	cwd            string
}

//Run checks for dependencies given the program and validates them against
//Packages.
func (dg *Depguard) Run(config *loader.Config, prog *loader.Program) ([]*Issue, error) {
	//Shortcut execution on an empty blacklist as that means every package is allowed
	if dg.ListType == LTBlacklist && len(dg.Packages) == 0 {
		return nil, nil
	}

	if err := dg.initialize(config, prog); err != nil {
		return nil, err
	}

	directImports, err := dg.createImportMap(prog)
	if err != nil {
		return nil, err
	}
	var issues []*Issue
	for pkg, positions := range directImports {
		if dg.flagIt(pkg) {
			for _, pos := range positions {
				issues = append(issues, &Issue{
					PackageName: pkg,
					Position:    pos,
				})
			}
		}
	}
	return issues, nil
}

func (dg *Depguard) initialize(config *loader.Config, prog *loader.Program) error {
	//Try and get the current working directory
	dg.cwd = config.Cwd
	if dg.cwd == "" {
		var err error
		dg.cwd, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	//Use the &build.Default if one is not specified
	dg.buildCtx = config.Build
	if dg.buildCtx == nil {
		dg.buildCtx = &build.Default
	}

	for _, pkg := range dg.Packages {
		if strings.ContainsAny(pkg, "!?*[]{}") {
			g, err := glob.Compile(pkg, '/')
			if err != nil {
				return err
			}
			dg.globPackages = append(dg.globPackages, g)
		} else {
			dg.prefixPackages = append(dg.prefixPackages, pkg)
		}
	}

	//Sort the packages so we can have a faster search in the array
	sort.Strings(dg.prefixPackages)
	return nil
}

func (dg *Depguard) createImportMap(prog *loader.Program) (map[string][]token.Position, error) {
	importMap := make(map[string][]token.Position)
	//For the directly imported packages
	for _, imported := range prog.InitialPackages() {
		//Go through their files
		for _, file := range imported.Files {
			//And populate a map of all direct imports and their positions
			//This will filter out GoRoot depending on the Depguard.IncludeGoRoot
			for _, fileImport := range file.Imports {
				fileImportPath := cleanBasicLitString(fileImport.Path.Value)
				if !dg.IncludeGoRoot {
					pkg, err := dg.buildCtx.Import(fileImportPath, dg.cwd, 0)
					if err != nil {
						return nil, err
					}
					if pkg.Goroot {
						continue
					}
				}
				position := prog.Fset.Position(fileImport.Pos())
				positions, found := importMap[fileImportPath]
				if !found {
					importMap[fileImportPath] = []token.Position{
						position,
					}
					continue
				}
				importMap[fileImportPath] = append(positions, position)
			}
		}
	}
	return importMap, nil
}

func (dg *Depguard) pkgInList(pkg string) bool {
	if dg.pkgInPrefixList(pkg) {
		return true
	}
	return dg.pkgInGlobList(pkg)
}

func (dg *Depguard) pkgInPrefixList(pkg string) bool {
	//Idx represents where in the package slice the passed in package would go
	//when sorted. -1 Just means that it would be at the very front of the slice.
	idx := sort.Search(len(dg.prefixPackages), func(i int) bool {
		return dg.prefixPackages[i] > pkg
	}) - 1
	//This means that the package passed in has no way to be prefixed by anything
	//in the package list as it is already smaller then everything
	if idx == -1 {
		return false
	}
	return strings.HasPrefix(pkg, dg.prefixPackages[idx])
}

func (dg *Depguard) pkgInGlobList(pkg string) bool {
	for _, g := range dg.globPackages {
		if g.Match(pkg) {
			return true
		}
	}
	return false
}

//InList | WhiteList | BlackList
//   y   |           |     x
//   n   |     x     |
func (dg *Depguard) flagIt(pkg string) bool {
	return dg.pkgInList(pkg) == (dg.ListType == LTBlacklist)
}

func cleanBasicLitString(value string) string {
	return strings.Trim(value, "\"\\")
}
