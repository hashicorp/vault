// Copyright (c) 2019, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

// Package format exposes gofumpt's formatting in an API similar to go/format.
// In general, the APIs are only guaranteed to work well when the input source
// is in canonical gofmt format.
package format

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/mod/semver"
	"golang.org/x/tools/go/ast/astutil"
)

type Options struct {
	// LangVersion corresponds to the Go language version a piece of code is
	// written in. The version is used to decide whether to apply formatting
	// rules which require new language features. When inside a Go module,
	// LangVersion should generally be specified as the result of:
	//
	//     go list -m -f {{.GoVersion}}
	//
	// LangVersion is treated as a semantic version, which might start with
	// a "v" prefix. Like Go versions, it might also be incomplete; "1.14"
	// is equivalent to "1.14.0". When empty, it is equivalent to "v1", to
	// not use language features which could break programs.
	LangVersion string

	ExtraRules bool
}

// Source formats src in gofumpt's format, assuming that src holds a valid Go
// source file.
func Source(src []byte, opts Options) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	File(fset, file, opts)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// File modifies a file and fset in place to follow gofumpt's format. The
// changes might include manipulating adding or removing newlines in fset,
// modifying the position of nodes, or modifying literal values.
func File(fset *token.FileSet, file *ast.File, opts Options) {
	if opts.LangVersion == "" {
		opts.LangVersion = "v1"
	} else if opts.LangVersion[0] != 'v' {
		opts.LangVersion = "v" + opts.LangVersion
	}
	if !semver.IsValid(opts.LangVersion) {
		panic(fmt.Sprintf("invalid semver string: %q", opts.LangVersion))
	}
	f := &fumpter{
		File:    fset.File(file.Pos()),
		fset:    fset,
		astFile: file,
		Options: opts,
	}
	pre := func(c *astutil.Cursor) bool {
		f.applyPre(c)
		if _, ok := c.Node().(*ast.BlockStmt); ok {
			f.blockLevel++
		}
		return true
	}
	post := func(c *astutil.Cursor) bool {
		if _, ok := c.Node().(*ast.BlockStmt); ok {
			f.blockLevel--
		}
		return true
	}
	astutil.Apply(file, pre, post)
}

// Multiline nodes which could fit on a single line under this many
// bytes may be collapsed onto a single line.
const shortLineLimit = 60

var rxOctalInteger = regexp.MustCompile(`\A0[0-7_]+\z`)

type fumpter struct {
	Options

	*token.File
	fset *token.FileSet

	astFile *ast.File

	blockLevel int
}

func (f *fumpter) commentsBetween(p1, p2 token.Pos) []*ast.CommentGroup {
	comments := f.astFile.Comments
	i1 := sort.Search(len(comments), func(i int) bool {
		return comments[i].Pos() >= p1
	})
	comments = comments[i1:]
	i2 := sort.Search(len(comments), func(i int) bool {
		return comments[i].Pos() >= p2
	})
	comments = comments[:i2]
	return comments
}

func (f *fumpter) inlineComment(pos token.Pos) *ast.Comment {
	comments := f.astFile.Comments
	i := sort.Search(len(comments), func(i int) bool {
		return comments[i].Pos() >= pos
	})
	if i >= len(comments) {
		return nil
	}
	line := f.Line(pos)
	for _, comment := range comments[i].List {
		if f.Line(comment.Pos()) == line {
			return comment
		}
	}
	return nil
}

// addNewline is a hack to let us force a newline at a certain position.
func (f *fumpter) addNewline(at token.Pos) {
	offset := f.Offset(at)

	field := reflect.ValueOf(f.File).Elem().FieldByName("lines")
	n := field.Len()
	lines := make([]int, 0, n+1)
	for i := 0; i < n; i++ {
		cur := int(field.Index(i).Int())
		if offset == cur {
			// This newline already exists; do nothing. Duplicate
			// newlines can't exist.
			return
		}
		if offset >= 0 && offset < cur {
			lines = append(lines, offset)
			offset = -1
		}
		lines = append(lines, cur)
	}
	if offset >= 0 {
		lines = append(lines, offset)
	}
	if !f.SetLines(lines) {
		panic(fmt.Sprintf("could not set lines to %v", lines))
	}
}

// removeLines removes all newlines between two positions, so that they end
// up on the same line.
func (f *fumpter) removeLines(fromLine, toLine int) {
	for fromLine < toLine {
		f.MergeLine(fromLine)
		toLine--
	}
}

// removeLinesBetween is like removeLines, but it leaves one newline between the
// two positions.
func (f *fumpter) removeLinesBetween(from, to token.Pos) {
	f.removeLines(f.Line(from)+1, f.Line(to))
}

type byteCounter int

func (b *byteCounter) Write(p []byte) (n int, err error) {
	*b += byteCounter(len(p))
	return len(p), nil
}

func (f *fumpter) printLength(node ast.Node) int {
	var count byteCounter
	if err := format.Node(&count, f.fset, node); err != nil {
		panic(fmt.Sprintf("unexpected print error: %v", err))
	}

	// Add the space taken by an inline comment.
	if c := f.inlineComment(node.End()); c != nil {
		fmt.Fprintf(&count, " %s", c.Text)
	}

	// Add an approximation of the indentation level. We can't know the
	// number of tabs go/printer will add ahead of time. Trying to print the
	// entire top-level declaration would tell us that, but then it's near
	// impossible to reliably find our node again.
	return int(count) + (f.blockLevel * 8)
}

// rxCommentDirective covers all common Go comment directives:
//
//   //go:         | standard Go directives, like go:noinline
//   //some-words: | similar to the syntax above, like lint:ignore or go-sumtype:decl
//   //line        | inserted line information for cmd/compile
//   //export      | to mark cgo funcs for exporting
//   //extern      | C function declarations for gccgo
//   //sys(nb)?    | syscall function wrapper prototypes
//   //nolint      | nolint directive for golangci
//
// Note that the "some-words:" matching expects a letter afterward, such as
// "go:generate", to prevent matching false positives like "https://site".
var rxCommentDirective = regexp.MustCompile(`^([a-z-]+:[a-z]+|line\b|export\b|extern\b|sys(nb)?\b|nolint\b)`)

// visit takes either an ast.Node or a []ast.Stmt.
func (f *fumpter) applyPre(c *astutil.Cursor) {
	switch node := c.Node().(type) {
	case *ast.File:
		var lastMulti bool
		var lastEnd token.Pos
		for _, decl := range node.Decls {
			pos := decl.Pos()
			comments := f.commentsBetween(lastEnd, pos)
			if len(comments) > 0 {
				pos = comments[0].Pos()
			}

			// multiline top-level declarations should be separated
			multi := f.Line(pos) < f.Line(decl.End())
			if multi && lastMulti && f.Line(lastEnd)+1 == f.Line(pos) {
				f.addNewline(lastEnd)
			}

			lastMulti = multi
			lastEnd = decl.End()
		}

		// Join contiguous lone var/const/import lines; abort if there
		// are empty lines or comments in between.
		newDecls := make([]ast.Decl, 0, len(node.Decls))
		for i := 0; i < len(node.Decls); {
			newDecls = append(newDecls, node.Decls[i])
			start, ok := node.Decls[i].(*ast.GenDecl)
			if !ok || isCgoImport(start) {
				i++
				continue
			}
			lastPos := start.Pos()
			for i++; i < len(node.Decls); {
				cont, ok := node.Decls[i].(*ast.GenDecl)
				if !ok || cont.Tok != start.Tok || cont.Lparen != token.NoPos ||
					f.Line(lastPos) < f.Line(cont.Pos())-1 || isCgoImport(cont) {
					break
				}
				start.Specs = append(start.Specs, cont.Specs...)
				if c := f.inlineComment(cont.End()); c != nil {
					// don't move an inline comment outside
					start.Rparen = c.End()
				}
				lastPos = cont.Pos()
				i++
			}
		}
		node.Decls = newDecls

		// Comments aren't nodes, so they're not walked by default.
	groupLoop:
		for _, group := range node.Comments {
			for _, comment := range group.List {
				body := strings.TrimPrefix(comment.Text, "//")
				if body == comment.Text {
					// /*-style comment
					continue groupLoop
				}
				if rxCommentDirective.MatchString(body) {
					// this line is a directive
					continue groupLoop
				}
				r, _ := utf8.DecodeRuneInString(body)
				if !unicode.IsLetter(r) && !unicode.IsNumber(r) && !unicode.IsSpace(r) {
					// this line could be code like "//{"
					continue groupLoop
				}
			}
			// If none of the comment group's lines look like a
			// directive or code, add spaces, if needed.
			for _, comment := range group.List {
				body := strings.TrimPrefix(comment.Text, "//")
				r, _ := utf8.DecodeRuneInString(body)
				if !unicode.IsSpace(r) {
					comment.Text = "// " + strings.TrimPrefix(comment.Text, "//")
				}
			}
		}

	case *ast.DeclStmt:
		decl, ok := node.Decl.(*ast.GenDecl)
		if !ok || decl.Tok != token.VAR || len(decl.Specs) != 1 {
			break // e.g. const name = "value"
		}
		spec := decl.Specs[0].(*ast.ValueSpec)
		if spec.Type != nil {
			break // e.g. var name Type
		}
		tok := token.ASSIGN
		names := make([]ast.Expr, len(spec.Names))
		for i, name := range spec.Names {
			names[i] = name
			if name.Name != "_" {
				tok = token.DEFINE
			}
		}
		c.Replace(&ast.AssignStmt{
			Lhs: names,
			Tok: tok,
			Rhs: spec.Values,
		})

	case *ast.GenDecl:
		if node.Tok == token.IMPORT && node.Lparen.IsValid() {
			f.joinStdImports(node)
		}

		// Single var declarations shouldn't use parentheses, unless
		// there's a comment on the grouped declaration.
		if node.Tok == token.VAR && len(node.Specs) == 1 &&
			node.Lparen.IsValid() && node.Doc == nil {
			specPos := node.Specs[0].Pos()
			specEnd := node.Specs[0].End()

			if len(f.commentsBetween(node.TokPos, specPos)) > 0 {
				// If the single spec has any comment, it must
				// go before the entire declaration now.
				node.TokPos = specPos
			} else {
				f.removeLines(f.Line(node.TokPos), f.Line(specPos))
			}
			f.removeLines(f.Line(specEnd), f.Line(node.Rparen))

			// Remove the parentheses. go/printer will automatically
			// get rid of the newlines.
			node.Lparen = token.NoPos
			node.Rparen = token.NoPos
		}

	case *ast.BlockStmt:
		f.stmts(node.List)
		comments := f.commentsBetween(node.Lbrace, node.Rbrace)
		if len(node.List) == 0 && len(comments) == 0 {
			f.removeLinesBetween(node.Lbrace, node.Rbrace)
			break
		}

		var sign *ast.FuncType
		var cond ast.Expr
		switch parent := c.Parent().(type) {
		case *ast.FuncDecl:
			sign = parent.Type
		case *ast.FuncLit:
			sign = parent.Type
		case *ast.IfStmt:
			cond = parent.Cond
		case *ast.ForStmt:
			cond = parent.Cond
		}

		if len(node.List) > 1 && sign == nil {
			// only if we have a single statement, or if
			// it's a func body.
			break
		}
		var bodyPos, bodyEnd token.Pos

		if len(node.List) > 0 {
			bodyPos = node.List[0].Pos()
			bodyEnd = node.List[len(node.List)-1].End()
		}
		if len(comments) > 0 {
			if pos := comments[0].Pos(); !bodyPos.IsValid() || pos < bodyPos {
				bodyPos = pos
			}
			if pos := comments[len(comments)-1].End(); !bodyPos.IsValid() || pos > bodyEnd {
				bodyEnd = pos
			}
		}

		f.removeLinesBetween(bodyEnd, node.Rbrace)

		if cond != nil && f.Line(cond.Pos()) != f.Line(cond.End()) {
			// The body is preceded by a multi-line condition, so an
			// empty line can help readability.
			return
		}
		if sign != nil {
			var lastParam *ast.Field
			if l := sign.Results; l != nil && len(l.List) > 0 {
				lastParam = l.List[len(l.List)-1]
			} else if l := sign.Params; l != nil && len(l.List) > 0 {
				lastParam = l.List[len(l.List)-1]
			}
			endLine := f.Line(sign.End())
			if lastParam != nil && f.Line(sign.Pos()) != endLine && f.Line(lastParam.Pos()) == endLine {
				// The body is preceded by a multi-line function
				// signature, and the empty line helps readability.
				return
			}
		}

		f.removeLinesBetween(node.Lbrace, bodyPos)

	case *ast.CompositeLit:
		if len(node.Elts) == 0 {
			// doesn't have elements
			break
		}
		openLine := f.Line(node.Lbrace)
		closeLine := f.Line(node.Rbrace)
		if openLine == closeLine {
			// all in a single line
			break
		}

		newlineAroundElems := false
		newlineBetweenElems := false
		lastLine := openLine
		for i, elem := range node.Elts {
			if f.Line(elem.Pos()) > lastLine {
				if i == 0 {
					newlineAroundElems = true
				} else {
					newlineBetweenElems = true
				}
			}
			lastLine = f.Line(elem.End())
		}
		if closeLine > lastLine {
			newlineAroundElems = true
		}

		if newlineBetweenElems || newlineAroundElems {
			first := node.Elts[0]
			if openLine == f.Line(first.Pos()) {
				// We want the newline right after the brace.
				f.addNewline(node.Lbrace + 1)
				closeLine = f.Line(node.Rbrace)
			}
			last := node.Elts[len(node.Elts)-1]
			if closeLine == f.Line(last.End()) {
				// We want the newline right before the brace.
				f.addNewline(node.Rbrace)
			}
		}

		// If there's a newline between any consecutive elements, there
		// must be a newline between all composite literal elements.
		if !newlineBetweenElems {
			break
		}
		for i1, elem1 := range node.Elts {
			i2 := i1 + 1
			if i2 >= len(node.Elts) {
				break
			}
			elem2 := node.Elts[i2]
			// TODO: do we care about &{}?
			_, ok1 := elem1.(*ast.CompositeLit)
			_, ok2 := elem2.(*ast.CompositeLit)
			if !ok1 && !ok2 {
				continue
			}
			if f.Line(elem1.End()) == f.Line(elem2.Pos()) {
				f.addNewline(elem1.End())
			}
		}

	case *ast.CaseClause:
		f.stmts(node.Body)
		openLine := f.Line(node.Case)
		closeLine := f.Line(node.Colon)
		if openLine == closeLine {
			// nothing to do
			break
		}
		if len(f.commentsBetween(node.Case, node.Colon)) > 0 {
			// don't move comments
			break
		}
		if f.printLength(node) > shortLineLimit {
			// too long to collapse
			break
		}
		f.removeLines(openLine, closeLine)

	case *ast.CommClause:
		f.stmts(node.Body)

	case *ast.FieldList:
		if node.NumFields() == 0 && f.inlineComment(node.Pos()) == nil {
			// Empty field lists should not contain a newline.
			// Do not join the two lines if the first has an inline
			// comment, as that can result in broken formatting.
			openLine := f.Line(node.Pos())
			closeLine := f.Line(node.End())
			f.removeLines(openLine, closeLine)
		}

		// Merging adjacent fields (e.g. parameters) is disabled by default.
		if !f.ExtraRules {
			break
		}
		switch c.Parent().(type) {
		case *ast.FuncDecl, *ast.FuncType, *ast.InterfaceType:
			node.List = f.mergeAdjacentFields(node.List)
			c.Replace(node)
		case *ast.StructType:
			// Do not merge adjacent fields in structs.
		}

	case *ast.BasicLit:
		// Octal number literals were introduced in 1.13.
		if semver.Compare(f.LangVersion, "v1.13") >= 0 {
			if node.Kind == token.INT && rxOctalInteger.MatchString(node.Value) {
				node.Value = "0o" + node.Value[1:]
				c.Replace(node)
			}
		}
	}
}

func (f *fumpter) stmts(list []ast.Stmt) {
	for i, stmt := range list {
		ifs, ok := stmt.(*ast.IfStmt)
		if !ok || i < 1 {
			continue // not an if following another statement
		}
		as, ok := list[i-1].(*ast.AssignStmt)
		if !ok || as.Tok != token.DEFINE ||
			!identEqual(as.Lhs[len(as.Lhs)-1], "err") {
			continue // not "..., err := ..."
		}
		be, ok := ifs.Cond.(*ast.BinaryExpr)
		if !ok || ifs.Init != nil || ifs.Else != nil {
			continue // complex if
		}
		if be.Op != token.NEQ || !identEqual(be.X, "err") ||
			!identEqual(be.Y, "nil") {
			continue // not "err != nil"
		}
		f.removeLinesBetween(as.End(), ifs.Pos())
	}
}

func identEqual(expr ast.Expr, name string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == name
}

// isCgoImport returns true if the declaration is simply:
//
//   import "C"
//
// Note that parentheses do not affect the result.
func isCgoImport(decl *ast.GenDecl) bool {
	if decl.Tok != token.IMPORT || len(decl.Specs) != 1 {
		return false
	}
	spec := decl.Specs[0].(*ast.ImportSpec)
	return spec.Path.Value == `"C"`
}

// joinStdImports ensures that all standard library imports are together and at
// the top of the imports list.
func (f *fumpter) joinStdImports(d *ast.GenDecl) {
	var std, other []ast.Spec
	firstGroup := true
	lastEnd := d.Pos()
	needsSort := false
	for i, spec := range d.Specs {
		spec := spec.(*ast.ImportSpec)
		if coms := f.commentsBetween(lastEnd, spec.Pos()); len(coms) > 0 {
			lastEnd = coms[len(coms)-1].End()
		}
		if i > 0 && firstGroup && f.Line(spec.Pos()) > f.Line(lastEnd)+1 {
			firstGroup = false
		} else {
			// We're still in the first group, update lastEnd.
			lastEnd = spec.End()
		}

		path, _ := strconv.Unquote(spec.Path.Value)
		switch {
		// Imports with a period are definitely third party.
		case strings.Contains(path, "."):
			fallthrough
		// "test" and "example" are reserved as per golang.org/issue/37641.
		// "internal" is unreachable.
		case strings.HasPrefix(path, "test/") ||
			strings.HasPrefix(path, "example/") ||
			strings.HasPrefix(path, "internal/"):
			fallthrough
		// To be conservative, if an import has a name or an inline
		// comment, and isn't part of the top group, treat it as non-std.
		case !firstGroup && (spec.Name != nil || spec.Comment != nil):
			other = append(other, spec)
			continue
		}

		// If we're moving this std import further up, reset its
		// position, to avoid breaking comments.
		if !firstGroup || len(other) > 0 {
			setPos(reflect.ValueOf(spec), d.Pos())
			needsSort = true
		}
		std = append(std, spec)
	}
	// Ensure there is an empty line between std imports and other imports.
	if len(std) > 0 && len(other) > 0 && f.Line(std[len(std)-1].End())+1 >= f.Line(other[0].Pos()) {
		// We add two newlines, as that's necessary in some edge cases.
		// For example, if the std and non-std imports were together and
		// without indentation, adding one newline isn't enough. Two
		// empty lines will be printed as one by go/printer, anyway.
		f.addNewline(other[0].Pos() - 1)
		f.addNewline(other[0].Pos())
	}
	// Finally, join the imports, keeping std at the top.
	d.Specs = append(std, other...)

	// If we moved any std imports to the first group, we need to sort them
	// again.
	if needsSort {
		ast.SortImports(f.fset, f.astFile)
	}
}

// mergeAdjacentFields returns fields with adjacent fields merged if possible.
func (f *fumpter) mergeAdjacentFields(fields []*ast.Field) []*ast.Field {
	// If there are less than two fields then there is nothing to merge.
	if len(fields) < 2 {
		return fields
	}

	// Otherwise, iterate over adjacent pairs of fields, merging if possible,
	// and mutating fields. Elements of fields may be mutated (if merged with
	// following fields), discarded (if merged with a preceeding field), or left
	// unchanged.
	i := 0
	for j := 1; j < len(fields); j++ {
		if f.shouldMergeAdjacentFields(fields[i], fields[j]) {
			fields[i].Names = append(fields[i].Names, fields[j].Names...)
		} else {
			i++
			fields[i] = fields[j]
		}
	}
	return fields[:i+1]
}

func (f *fumpter) shouldMergeAdjacentFields(f1, f2 *ast.Field) bool {
	if len(f1.Names) == 0 || len(f2.Names) == 0 {
		// Both must have names for the merge to work.
		return false
	}
	if f.Line(f1.Pos()) != f.Line(f2.Pos()) {
		// Trust the user if they used separate lines.
		return false
	}

	// Only merge if the types are equal.
	opt := cmp.Comparer(func(x, y token.Pos) bool { return true })
	return cmp.Equal(f1.Type, f2.Type, opt)
}

var posType = reflect.TypeOf(token.NoPos)

// setPos recursively sets all position fields in the node v to pos.
func setPos(v reflect.Value, pos token.Pos) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return
	}
	if v.Type() == posType {
		v.Set(reflect.ValueOf(pos))
	}
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			setPos(v.Field(i), pos)
		}
	}
}
