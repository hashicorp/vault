// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package analysistest provides utilities for testing analyzers.
package analysistest

import (
	"bytes"
	"fmt"
	"go/format"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
	"text/scanner"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/internal/checker"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/internal/diff"
	"golang.org/x/tools/internal/testenv"
	"golang.org/x/tools/txtar"
)

// WriteFiles is a helper function that creates a temporary directory
// and populates it with a GOPATH-style project using filemap (which
// maps file names to contents). On success it returns the name of the
// directory and a cleanup function to delete it.
func WriteFiles(filemap map[string]string) (dir string, cleanup func(), err error) {
	gopath, err := os.MkdirTemp("", "analysistest")
	if err != nil {
		return "", nil, err
	}
	cleanup = func() { os.RemoveAll(gopath) }

	for name, content := range filemap {
		filename := filepath.Join(gopath, "src", name)
		os.MkdirAll(filepath.Dir(filename), 0777) // ignore error
		if err := os.WriteFile(filename, []byte(content), 0666); err != nil {
			cleanup()
			return "", nil, err
		}
	}
	return gopath, cleanup, nil
}

// TestData returns the effective filename of
// the program's "testdata" directory.
// This function may be overridden by projects using
// an alternative build system (such as Blaze) that
// does not run a test in its package directory.
var TestData = func() string {
	testdata, err := filepath.Abs("testdata")
	if err != nil {
		log.Fatal(err)
	}
	return testdata
}

// Testing is an abstraction of a *testing.T.
type Testing interface {
	Errorf(format string, args ...interface{})
}

// RunWithSuggestedFixes behaves like Run, but additionally verifies suggested fixes.
// It uses golden files placed alongside the source code under analysis:
// suggested fixes for code in example.go will be compared against example.go.golden.
//
// Golden files can be formatted in one of two ways: as plain Go source code, or as txtar archives.
// In the first case, all suggested fixes will be applied to the original source, which will then be compared against the golden file.
// In the second case, suggested fixes will be grouped by their messages, and each set of fixes will be applied and tested separately.
// Each section in the archive corresponds to a single message.
//
// A golden file using txtar may look like this:
//
//	-- turn into single negation --
//	package pkg
//
//	func fn(b1, b2 bool) {
//		if !b1 { // want `negating a boolean twice`
//			println()
//		}
//	}
//
//	-- remove double negation --
//	package pkg
//
//	func fn(b1, b2 bool) {
//		if b1 { // want `negating a boolean twice`
//			println()
//		}
//	}
//
// # Conflicts
//
// A single analysis pass may offer two or more suggested fixes that
// (1) conflict but are nonetheless logically composable, (e.g.
// because both update the import declaration), or (2) are
// fundamentally incompatible (e.g. alternative fixes to the same
// statement).
//
// It is up to the driver to decide how to apply such fixes. A
// sophisticated driver could attempt to resolve conflicts of the
// first kind, but this test driver simply reports the fact of the
// conflict with the expectation that the user will split their tests
// into nonconflicting parts.
//
// Conflicts of the second kind can be avoided by giving the
// alternative fixes different names (SuggestedFix.Message) and
// defining the .golden file as a multi-section txtar file with a
// named section for each alternative fix, as shown above.
//
// Analyzers that compute fixes from a textual diff of the
// before/after file contents (instead of directly from syntax tree
// positions) may produce fixes that, although logically
// non-conflicting, nonetheless conflict due to the particulars of the
// diff algorithm. In such cases it may suffice to introduce
// sufficient separation of the statements in the test input so that
// the computed diffs do not overlap. If that fails, break the test
// into smaller parts.
//
// TODO(adonovan): the behavior of RunWithSuggestedFixes as documented
// above is impractical for tests that report multiple diagnostics and
// offer multiple alternative fixes for the same diagnostic, and it is
// inconsistent with the interpretation of multiple diagnostics
// described at Diagnostic.SuggestedFixes.
// We need to rethink the analyzer testing API to better support such
// cases. In the meantime, users of RunWithSuggestedFixes testing
// analyzers that offer alternative fixes are advised to put each fix
// in a separate .go file in the testdata.
func RunWithSuggestedFixes(t Testing, dir string, a *analysis.Analyzer, patterns ...string) []*Result {
	r := Run(t, dir, a, patterns...)

	// If the immediate caller of RunWithSuggestedFixes is in
	// x/tools, we apply stricter checks as required by gopls.
	inTools := false
	{
		var pcs [1]uintptr
		n := runtime.Callers(1, pcs[:])
		frames := runtime.CallersFrames(pcs[:n])
		fr, _ := frames.Next()
		if fr.Func != nil && strings.HasPrefix(fr.Func.Name(), "golang.org/x/tools/") {
			inTools = true
		}
	}

	// Process each result (package) separately, matching up the suggested
	// fixes into a diff, which we will compare to the .golden file.  We have
	// to do this per-result in case a file appears in two packages, such as in
	// packages with tests, where mypkg/a.go will appear in both mypkg and
	// mypkg.test.  In that case, the analyzer may suggest the same set of
	// changes to a.go for each package.  If we merge all the results, those
	// changes get doubly applied, which will cause conflicts or mismatches.
	// Validating the results separately means as long as the two analyses
	// don't produce conflicting suggestions for a single file, everything
	// should match up.
	for _, act := range r {
		// file -> message -> edits
		fileEdits := make(map[*token.File]map[string][]diff.Edit)
		fileContents := make(map[*token.File][]byte)

		// Validate edits, prepare the fileEdits map and read the file contents.
		for _, diag := range act.Diagnostics {
			for _, fix := range diag.SuggestedFixes {

				// Assert that lazy fixes have a Category (#65578, #65087).
				if inTools && len(fix.TextEdits) == 0 && diag.Category == "" {
					t.Errorf("missing Diagnostic.Category for SuggestedFix without TextEdits (gopls requires the category for the name of the fix command")
				}

				for _, edit := range fix.TextEdits {
					start, end := edit.Pos, edit.End
					if !end.IsValid() {
						end = start
					}
					// Validate the edit.
					if start > end {
						t.Errorf(
							"diagnostic for analysis %v contains Suggested Fix with malformed edit: pos (%v) > end (%v)",
							act.Pass.Analyzer.Name, start, end)
						continue
					}
					file, endfile := act.Pass.Fset.File(start), act.Pass.Fset.File(end)
					if file == nil || endfile == nil || file != endfile {
						t.Errorf(
							"diagnostic for analysis %v contains Suggested Fix with malformed spanning files %v and %v",
							act.Pass.Analyzer.Name, file.Name(), endfile.Name())
						continue
					}
					if _, ok := fileContents[file]; !ok {
						contents, err := os.ReadFile(file.Name())
						if err != nil {
							t.Errorf("error reading %s: %v", file.Name(), err)
						}
						fileContents[file] = contents
					}
					if _, ok := fileEdits[file]; !ok {
						fileEdits[file] = make(map[string][]diff.Edit)
					}
					fileEdits[file][fix.Message] = append(fileEdits[file][fix.Message], diff.Edit{
						Start: file.Offset(start),
						End:   file.Offset(end),
						New:   string(edit.NewText),
					})
				}
			}
		}

		for file, fixes := range fileEdits {
			// Get the original file contents.
			orig, ok := fileContents[file]
			if !ok {
				t.Errorf("could not find file contents for %s", file.Name())
				continue
			}

			// Get the golden file and read the contents.
			ar, err := txtar.ParseFile(file.Name() + ".golden")
			if err != nil {
				t.Errorf("error reading %s.golden: %v", file.Name(), err)
				continue
			}

			if len(ar.Files) > 0 {
				// one virtual file per kind of suggested fix

				if len(ar.Comment) != 0 {
					// we allow either just the comment, or just virtual
					// files, not both. it is not clear how "both" should
					// behave.
					t.Errorf("%s.golden has leading comment; we don't know what to do with it", file.Name())
					continue
				}

				for sf, edits := range fixes {
					found := false
					for _, vf := range ar.Files {
						if vf.Name == sf {
							found = true
							// the file may contain multiple trailing
							// newlines if the user places empty lines
							// between files in the archive. normalize
							// this to a single newline.
							golden := append(bytes.TrimRight(vf.Data, "\n"), '\n')

							if err := applyDiffsAndCompare(orig, golden, edits, file.Name()); err != nil {
								t.Errorf("%s", err)
							}
							break
						}
					}
					if !found {
						t.Errorf("no section for suggested fix %q in %s.golden", sf, file.Name())
					}
				}
			} else {
				// all suggested fixes are represented by a single file

				var catchallEdits []diff.Edit
				for _, edits := range fixes {
					catchallEdits = append(catchallEdits, edits...)
				}

				if err := applyDiffsAndCompare(orig, ar.Comment, catchallEdits, file.Name()); err != nil {
					t.Errorf("%s", err)
				}
			}
		}
	}
	return r
}

// applyDiffsAndCompare applies edits to src and compares the results against
// golden after formatting both. fileName is use solely for error reporting.
func applyDiffsAndCompare(src, golden []byte, edits []diff.Edit, fileName string) error {
	out, err := diff.ApplyBytes(src, edits)
	if err != nil {
		return fmt.Errorf("%s: error applying fixes: %v (see possible explanations at RunWithSuggestedFixes)", fileName, err)
	}
	wantRaw, err := format.Source(golden)
	if err != nil {
		return fmt.Errorf("%s.golden: error formatting golden file: %v\n%s", fileName, err, out)
	}
	want := string(wantRaw)

	formatted, err := format.Source(out)
	if err != nil {
		return fmt.Errorf("%s: error formatting resulting source: %v\n%s", fileName, err, out)
	}
	if got := string(formatted); got != want {
		unified := diff.Unified(fileName+".golden", "actual", want, got)
		return fmt.Errorf("suggested fixes failed for %s:\n%s", fileName, unified)
	}
	return nil
}

// Run applies an analysis to the packages denoted by the "go list" patterns.
//
// It loads the packages from the specified
// directory using golang.org/x/tools/go/packages, runs the analysis on
// them, and checks that each analysis emits the expected diagnostics
// and facts specified by the contents of '// want ...' comments in the
// package's source files. It treats a comment of the form
// "//...// want..." or "/*...// want... */" as if it starts at 'want'.
//
// If the directory contains a go.mod file, Run treats it as the root of the
// Go module in which to work. Otherwise, Run treats it as the root of a
// GOPATH-style tree, with package contained in the src subdirectory.
//
// An expectation of a Diagnostic is specified by a string literal
// containing a regular expression that must match the diagnostic
// message. For example:
//
//	fmt.Printf("%s", 1) // want `cannot provide int 1 to %s`
//
// An expectation of a Fact associated with an object is specified by
// 'name:"pattern"', where name is the name of the object, which must be
// declared on the same line as the comment, and pattern is a regular
// expression that must match the string representation of the fact,
// fmt.Sprint(fact). For example:
//
//	func panicf(format string, args interface{}) { // want panicf:"printfWrapper"
//
// Package facts are specified by the name "package" and appear on
// line 1 of the first source file of the package.
//
// A single 'want' comment may contain a mixture of diagnostic and fact
// expectations, including multiple facts about the same object:
//
//	// want "diag" "diag2" x:"fact1" x:"fact2" y:"fact3"
//
// Unexpected diagnostics and facts, and unmatched expectations, are
// reported as errors to the Testing.
//
// Run reports an error to the Testing if loading or analysis failed.
// Run also returns a Result for each package for which analysis was
// attempted, even if unsuccessful. It is safe for a test to ignore all
// the results, but a test may use it to perform additional checks.
func Run(t Testing, dir string, a *analysis.Analyzer, patterns ...string) []*Result {
	if t, ok := t.(testing.TB); ok {
		testenv.NeedsGoPackages(t)
	}

	pkgs, err := loadPackages(a, dir, patterns...)
	if err != nil {
		t.Errorf("loading %s: %v", patterns, err)
		return nil
	}

	if err := analysis.Validate([]*analysis.Analyzer{a}); err != nil {
		t.Errorf("Validate: %v", err)
		return nil
	}

	results := checker.TestAnalyzer(a, pkgs)
	for _, result := range results {
		if result.Err != nil {
			t.Errorf("error analyzing %s: %v", result.Pass, result.Err)
		} else {
			check(t, dir, result.Pass, result.Diagnostics, result.Facts)
		}
	}
	return results
}

// A Result holds the result of applying an analyzer to a package.
type Result = checker.TestAnalyzerResult

// loadPackages uses go/packages to load a specified packages (from source, with
// dependencies) from dir, which is the root of a GOPATH-style project tree.
// loadPackages returns an error if any package had an error, or the pattern
// matched no packages.
func loadPackages(a *analysis.Analyzer, dir string, patterns ...string) ([]*packages.Package, error) {
	env := []string{"GOPATH=" + dir, "GO111MODULE=off", "GOWORK=off"} // GOPATH mode

	// Undocumented module mode. Will be replaced by something better.
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		gowork := filepath.Join(dir, "go.work")
		if _, err := os.Stat(gowork); err != nil {
			gowork = "off"
		}

		env = []string{"GO111MODULE=on", "GOPROXY=off", "GOWORK=" + gowork} // module mode
	}

	// packages.Load loads the real standard library, not a minimal
	// fake version, which would be more efficient, especially if we
	// have many small tests that import, say, net/http.
	// However there is no easy way to make go/packages to consume
	// a list of packages we generate and then do the parsing and
	// typechecking, though this feature seems to be a recurring need.

	mode := packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports |
		packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo |
		packages.NeedDeps | packages.NeedModule
	cfg := &packages.Config{
		Mode:  mode,
		Dir:   dir,
		Tests: true,
		Env:   append(os.Environ(), env...),
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, err
	}

	// If any named package couldn't be loaded at all
	// (e.g. the Name field is unset), fail fast.
	for _, pkg := range pkgs {
		if pkg.Name == "" {
			return nil, fmt.Errorf("failed to load %q: Errors=%v",
				pkg.PkgPath, pkg.Errors)
		}
	}

	// Do NOT print errors if the analyzer will continue running.
	// It is incredibly confusing for tests to be printing to stderr
	// willy-nilly instead of their test logs, especially when the
	// errors are expected and are going to be fixed.
	if !a.RunDespiteErrors {
		if packages.PrintErrors(pkgs) > 0 {
			return nil, fmt.Errorf("there were package loading errors (and RunDespiteErrors is false)")
		}
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages matched %s", patterns)
	}
	return pkgs, nil
}

// check inspects an analysis pass on which the analysis has already
// been run, and verifies that all reported diagnostics and facts match
// specified by the contents of "// want ..." comments in the package's
// source files, which must have been parsed with comments enabled.
func check(t Testing, gopath string, pass *analysis.Pass, diagnostics []analysis.Diagnostic, facts map[types.Object][]analysis.Fact) {
	type key struct {
		file string
		line int
	}

	want := make(map[key][]expectation)

	// processComment parses expectations out of comments.
	processComment := func(filename string, linenum int, text string) {
		text = strings.TrimSpace(text)

		// Any comment starting with "want" is treated
		// as an expectation, even without following whitespace.
		if rest := strings.TrimPrefix(text, "want"); rest != text {
			lineDelta, expects, err := parseExpectations(rest)
			if err != nil {
				t.Errorf("%s:%d: in 'want' comment: %s", filename, linenum, err)
				return
			}
			if expects != nil {
				want[key{filename, linenum + lineDelta}] = expects
			}
		}
	}

	// Extract 'want' comments from parsed Go files.
	for _, f := range pass.Files {
		for _, cgroup := range f.Comments {
			for _, c := range cgroup.List {

				text := strings.TrimPrefix(c.Text, "//")
				if text == c.Text { // not a //-comment.
					text = strings.TrimPrefix(text, "/*")
					text = strings.TrimSuffix(text, "*/")
				}

				// Hack: treat a comment of the form "//...// want..."
				// or "/*...// want... */
				// as if it starts at 'want'.
				// This allows us to add comments on comments,
				// as required when testing the buildtag analyzer.
				if i := strings.Index(text, "// want"); i >= 0 {
					text = text[i+len("// "):]
				}

				// It's tempting to compute the filename
				// once outside the loop, but it's
				// incorrect because it can change due
				// to //line directives.
				posn := pass.Fset.Position(c.Pos())
				filename := sanitize(gopath, posn.Filename)
				processComment(filename, posn.Line, text)
			}
		}
	}

	// Extract 'want' comments from non-Go files.
	// TODO(adonovan): we may need to handle //line directives.
	for _, filename := range pass.OtherFiles {
		data, err := os.ReadFile(filename)
		if err != nil {
			t.Errorf("can't read '// want' comments from %s: %v", filename, err)
			continue
		}
		filename := sanitize(gopath, filename)
		linenum := 0
		for _, line := range strings.Split(string(data), "\n") {
			linenum++

			// Hack: treat a comment of the form "//...// want..."
			// or "/*...// want... */
			// as if it starts at 'want'.
			// This allows us to add comments on comments,
			// as required when testing the buildtag analyzer.
			if i := strings.Index(line, "// want"); i >= 0 {
				line = line[i:]
			}

			if i := strings.Index(line, "//"); i >= 0 {
				line = line[i+len("//"):]
				processComment(filename, linenum, line)
			}
		}
	}

	checkMessage := func(posn token.Position, kind, name, message string) {
		posn.Filename = sanitize(gopath, posn.Filename)
		k := key{posn.Filename, posn.Line}
		expects := want[k]
		var unmatched []string
		for i, exp := range expects {
			if exp.kind == kind && exp.name == name {
				if exp.rx.MatchString(message) {
					// matched: remove the expectation.
					expects[i] = expects[len(expects)-1]
					expects = expects[:len(expects)-1]
					want[k] = expects
					return
				}
				unmatched = append(unmatched, fmt.Sprintf("%#q", exp.rx))
			}
		}
		if unmatched == nil {
			t.Errorf("%v: unexpected %s: %v", posn, kind, message)
		} else {
			t.Errorf("%v: %s %q does not match pattern %s",
				posn, kind, message, strings.Join(unmatched, " or "))
		}
	}

	// Check the diagnostics match expectations.
	for _, f := range diagnostics {
		// TODO(matloob): Support ranges in analysistest.
		posn := pass.Fset.Position(f.Pos)
		checkMessage(posn, "diagnostic", "", f.Message)
	}

	// Check the facts match expectations.
	// Report errors in lexical order for determinism.
	// (It's only deterministic within each file, not across files,
	// because go/packages does not guarantee file.Pos is ascending
	// across the files of a single compilation unit.)
	var objects []types.Object
	for obj := range facts {
		objects = append(objects, obj)
	}
	sort.Slice(objects, func(i, j int) bool {
		// Package facts compare less than object facts.
		ip, jp := objects[i] == nil, objects[j] == nil // whether i, j is a package fact
		if ip != jp {
			return ip && !jp
		}
		return objects[i].Pos() < objects[j].Pos()
	})
	for _, obj := range objects {
		var posn token.Position
		var name string
		if obj != nil {
			// Object facts are reported on the declaring line.
			name = obj.Name()
			posn = pass.Fset.Position(obj.Pos())
		} else {
			// Package facts are reported at the start of the file.
			name = "package"
			posn = pass.Fset.Position(pass.Files[0].Pos())
			posn.Line = 1
		}

		for _, fact := range facts[obj] {
			checkMessage(posn, "fact", name, fmt.Sprint(fact))
		}
	}

	// Reject surplus expectations.
	//
	// Sometimes an Analyzer reports two similar diagnostics on a
	// line with only one expectation. The reader may be confused by
	// the error message.
	// TODO(adonovan): print a better error:
	// "got 2 diagnostics here; each one needs its own expectation".
	var surplus []string
	for key, expects := range want {
		for _, exp := range expects {
			err := fmt.Sprintf("%s:%d: no %s was reported matching %#q", key.file, key.line, exp.kind, exp.rx)
			surplus = append(surplus, err)
		}
	}
	sort.Strings(surplus)
	for _, err := range surplus {
		t.Errorf("%s", err)
	}
}

type expectation struct {
	kind string // either "fact" or "diagnostic"
	name string // name of object to which fact belongs, or "package" ("fact" only)
	rx   *regexp.Regexp
}

func (ex expectation) String() string {
	return fmt.Sprintf("%s %s:%q", ex.kind, ex.name, ex.rx) // for debugging
}

// parseExpectations parses the content of a "// want ..." comment
// and returns the expectations, a mixture of diagnostics ("rx") and
// facts (name:"rx").
func parseExpectations(text string) (lineDelta int, expects []expectation, err error) {
	var scanErr string
	sc := new(scanner.Scanner).Init(strings.NewReader(text))
	sc.Error = func(s *scanner.Scanner, msg string) {
		scanErr = msg // e.g. bad string escape
	}
	sc.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanInts

	scanRegexp := func(tok rune) (*regexp.Regexp, error) {
		if tok != scanner.String && tok != scanner.RawString {
			return nil, fmt.Errorf("got %s, want regular expression",
				scanner.TokenString(tok))
		}
		pattern, _ := strconv.Unquote(sc.TokenText()) // can't fail
		return regexp.Compile(pattern)
	}

	for {
		tok := sc.Scan()
		switch tok {
		case '+':
			tok = sc.Scan()
			if tok != scanner.Int {
				return 0, nil, fmt.Errorf("got +%s, want +Int", scanner.TokenString(tok))
			}
			lineDelta, _ = strconv.Atoi(sc.TokenText())
		case scanner.String, scanner.RawString:
			rx, err := scanRegexp(tok)
			if err != nil {
				return 0, nil, err
			}
			expects = append(expects, expectation{"diagnostic", "", rx})

		case scanner.Ident:
			name := sc.TokenText()
			tok = sc.Scan()
			if tok != ':' {
				return 0, nil, fmt.Errorf("got %s after %s, want ':'",
					scanner.TokenString(tok), name)
			}
			tok = sc.Scan()
			rx, err := scanRegexp(tok)
			if err != nil {
				return 0, nil, err
			}
			expects = append(expects, expectation{"fact", name, rx})

		case scanner.EOF:
			if scanErr != "" {
				return 0, nil, fmt.Errorf("%s", scanErr)
			}
			return lineDelta, expects, nil

		default:
			return 0, nil, fmt.Errorf("unexpected %s", scanner.TokenString(tok))
		}
	}
}

// sanitize removes the GOPATH portion of the filename,
// typically a gnarly /tmp directory, and returns the rest.
func sanitize(gopath, filename string) string {
	prefix := gopath + string(os.PathSeparator) + "src" + string(os.PathSeparator)
	return filepath.ToSlash(strings.TrimPrefix(filename, prefix))
}
