// (c) Copyright 2016 Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Modifications copyright (C) 2018 GolangCI

// Package gosec holds the central scanning logic used by gosec security scanner
package gosec

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/tools/go/loader"
)

// The Context is populated with data parsed from the source code as it is scanned.
// It is passed through to all rule functions as they are called. Rules may use
// this data in conjunction withe the encoutered AST node.
type Context struct {
	FileSet  *token.FileSet
	Comments ast.CommentMap
	Info     *types.Info
	Pkg      *types.Package
	Root     *ast.File
	Config   map[string]interface{}
	Imports  *ImportTracker
	Ignores  []map[string]bool
}

// Metrics used when reporting information about a scanning run.
type Metrics struct {
	NumFiles int `json:"files"`
	NumLines int `json:"lines"`
	NumNosec int `json:"nosec"`
	NumFound int `json:"found"`
}

// Analyzer object is the main object of gosec. It has methods traverse an AST
// and invoke the correct checking rules as on each node as required.
type Analyzer struct {
	ignoreNosec bool
	ruleset     RuleSet
	context     *Context
	config      Config
	logger      *log.Logger
	issues      []*Issue
	stats       *Metrics
}

// NewAnalyzer builds a new anaylzer.
func NewAnalyzer(conf Config, logger *log.Logger) *Analyzer {
	ignoreNoSec := false
	if setting, err := conf.GetGlobal("nosec"); err == nil {
		ignoreNoSec = setting == "true" || setting == "enabled"
	}
	if logger == nil {
		logger = log.New(os.Stderr, "[gosec]", log.LstdFlags)
	}
	return &Analyzer{
		ignoreNosec: ignoreNoSec,
		ruleset:     make(RuleSet),
		context:     &Context{},
		config:      conf,
		logger:      logger,
		issues:      make([]*Issue, 0, 16),
		stats:       &Metrics{},
	}
}

// LoadRules instantiates all the rules to be used when analyzing source
// packages
func (gosec *Analyzer) LoadRules(ruleDefinitions map[string]RuleBuilder) {
	for id, def := range ruleDefinitions {
		r, nodes := def(id, gosec.config)
		gosec.ruleset.Register(r, nodes...)
	}
}

// Process kicks off the analysis process for a given package
func (gosec *Analyzer) Process(buildTags []string, packagePaths ...string) error {
	ctx := build.Default
	ctx.BuildTags = append(ctx.BuildTags, buildTags...)
	packageConfig := loader.Config{
		Build:       &ctx,
		ParserMode:  parser.ParseComments,
		AllowErrors: true,
	}
	for _, packagePath := range packagePaths {
		abspath, err := GetPkgAbsPath(packagePath)
		if err != nil {
			gosec.logger.Printf("Skipping: %s. Path doesn't exist.", abspath)
			continue
		}
		gosec.logger.Println("Searching directory:", abspath)

		basePackage, err := build.Default.ImportDir(packagePath, build.ImportComment)
		if err != nil {
			return err
		}

		var packageFiles []string
		for _, filename := range basePackage.GoFiles {
			packageFiles = append(packageFiles, path.Join(packagePath, filename))
		}

		packageConfig.CreateFromFilenames(basePackage.Name, packageFiles...)
	}

	builtPackage, err := packageConfig.Load()
	if err != nil {
		return err
	}

	gosec.ProcessProgram(builtPackage)
	return nil
}

// ProcessProgram kicks off the analysis process for a given program
func (gosec *Analyzer) ProcessProgram(builtPackage *loader.Program) {
	for _, pkg := range builtPackage.InitialPackages() {
		gosec.logger.Println("Checking package:", pkg.String())
		for _, file := range pkg.Files {
			gosec.logger.Println("Checking file:", builtPackage.Fset.File(file.Pos()).Name())
			gosec.context.FileSet = builtPackage.Fset
			gosec.context.Config = gosec.config
			gosec.context.Comments = ast.NewCommentMap(gosec.context.FileSet, file, file.Comments)
			gosec.context.Root = file
			gosec.context.Info = &pkg.Info
			gosec.context.Pkg = pkg.Pkg
			gosec.context.Imports = NewImportTracker()
			gosec.context.Imports.TrackPackages(gosec.context.Pkg.Imports()...)
			ast.Walk(gosec, file)
			gosec.stats.NumFiles++
			gosec.stats.NumLines += builtPackage.Fset.File(file.Pos()).LineCount()
		}
	}
}

// ignore a node (and sub-tree) if it is tagged with a "#nosec" comment
func (gosec *Analyzer) ignore(n ast.Node) ([]string, bool) {
	if groups, ok := gosec.context.Comments[n]; ok && !gosec.ignoreNosec {
		for _, group := range groups {
			if strings.Contains(group.Text(), "#nosec") {
				gosec.stats.NumNosec++

				// Pull out the specific rules that are listed to be ignored.
				re := regexp.MustCompile("(G\\d{3})")
				matches := re.FindAllStringSubmatch(group.Text(), -1)

				// If no specific rules were given, ignore everything.
				if matches == nil || len(matches) == 0 {
					return nil, true
				}

				// Find the rule IDs to ignore.
				var ignores []string
				for _, v := range matches {
					ignores = append(ignores, v[1])
				}
				return ignores, false
			}
		}
	}
	return nil, false
}

// Visit runs the gosec visitor logic over an AST created by parsing go code.
// Rule methods added with AddRule will be invoked as necessary.
func (gosec *Analyzer) Visit(n ast.Node) ast.Visitor {
	// If we've reached the end of this branch, pop off the ignores stack.
	if n == nil {
		if len(gosec.context.Ignores) > 0 {
			gosec.context.Ignores = gosec.context.Ignores[1:]
		}
		return gosec
	}

	// Get any new rule exclusions.
	ignoredRules, ignoreAll := gosec.ignore(n)
	if ignoreAll {
		return nil
	}

	// Now create the union of exclusions.
	ignores := make(map[string]bool, 0)
	if len(gosec.context.Ignores) > 0 {
		for k, v := range gosec.context.Ignores[0] {
			ignores[k] = v
		}
	}

	for _, v := range ignoredRules {
		ignores[v] = true
	}

	// Push the new set onto the stack.
	gosec.context.Ignores = append([]map[string]bool{ignores}, gosec.context.Ignores...)

	// Track aliased and initialization imports
	gosec.context.Imports.TrackImport(n)

	for _, rule := range gosec.ruleset.RegisteredFor(n) {
		if _, ok := ignores[rule.ID()]; ok {
			continue
		}
		issue, err := rule.Match(n, gosec.context)
		if err != nil {
			file, line := GetLocation(n, gosec.context)
			file = path.Base(file)
			gosec.logger.Printf("Rule error: %v => %s (%s:%d)\n", reflect.TypeOf(rule), err, file, line)
		}
		if issue != nil {
			gosec.issues = append(gosec.issues, issue)
			gosec.stats.NumFound++
		}
	}
	return gosec
}

// Report returns the current issues discovered and the metrics about the scan
func (gosec *Analyzer) Report() ([]*Issue, *Metrics) {
	return gosec.issues, gosec.stats
}

// Reset clears state such as context, issues and metrics from the configured analyzer
func (gosec *Analyzer) Reset() {
	gosec.context = &Context{}
	gosec.issues = make([]*Issue, 0, 16)
	gosec.stats = &Metrics{}
}
