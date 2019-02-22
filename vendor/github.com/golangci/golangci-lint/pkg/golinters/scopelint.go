package golinters

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Scopelint struct{}

func (Scopelint) Name() string {
	return "scopelint"
}

func (Scopelint) Desc() string {
	return "Scopelint checks for unpinned variables in go programs"
}

func (lint Scopelint) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var res []result.Issue

	for _, f := range lintCtx.ASTCache.GetAllValidFiles() {
		n := Node{
			fset:          f.Fset,
			dangerObjects: map[*ast.Object]struct{}{},
			unsafeObjects: map[*ast.Object]struct{}{},
			skipFuncs:     map[*ast.FuncLit]struct{}{},
			issues:        &res,
		}
		ast.Walk(&n, f.F)
	}

	return res, nil
}

// The code below is copy-pasted from https://github.com/kyoh86/scopelint

// Node represents a Node being linted.
type Node struct {
	fset          *token.FileSet
	dangerObjects map[*ast.Object]struct{}
	unsafeObjects map[*ast.Object]struct{}
	skipFuncs     map[*ast.FuncLit]struct{}
	issues        *[]result.Issue
}

// Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
//nolint:gocyclo,gocritic
func (f *Node) Visit(node ast.Node) ast.Visitor {
	switch typedNode := node.(type) {
	case *ast.ForStmt:
		switch init := typedNode.Init.(type) {
		case *ast.AssignStmt:
			for _, lh := range init.Lhs {
				switch tlh := lh.(type) {
				case *ast.Ident:
					f.unsafeObjects[tlh.Obj] = struct{}{}
				}
			}
		}

	case *ast.RangeStmt:
		// Memory variables declarated in range statement
		switch k := typedNode.Key.(type) {
		case *ast.Ident:
			f.unsafeObjects[k.Obj] = struct{}{}
		}
		switch v := typedNode.Value.(type) {
		case *ast.Ident:
			f.unsafeObjects[v.Obj] = struct{}{}
		}

	case *ast.UnaryExpr:
		if typedNode.Op == token.AND {
			switch ident := typedNode.X.(type) {
			case *ast.Ident:
				if _, unsafe := f.unsafeObjects[ident.Obj]; unsafe {
					f.errorf(ident, "Using a reference for the variable on range scope %s", formatCode(ident.Name, nil))
				}
			}
		}

	case *ast.Ident:
		if _, obj := f.dangerObjects[typedNode.Obj]; obj {
			// It is the naked variable in scope of range statement.
			f.errorf(node, "Using the variable on range scope %s in function literal", formatCode(typedNode.Name, nil))
			break
		}

	case *ast.CallExpr:
		// Ignore func literals that'll be called immediately.
		switch funcLit := typedNode.Fun.(type) {
		case *ast.FuncLit:
			f.skipFuncs[funcLit] = struct{}{}
		}

	case *ast.FuncLit:
		if _, skip := f.skipFuncs[typedNode]; !skip {
			dangers := map[*ast.Object]struct{}{}
			for d := range f.dangerObjects {
				dangers[d] = struct{}{}
			}
			for u := range f.unsafeObjects {
				dangers[u] = struct{}{}
			}
			return &Node{
				fset:          f.fset,
				dangerObjects: dangers,
				unsafeObjects: f.unsafeObjects,
				skipFuncs:     f.skipFuncs,
				issues:        f.issues,
			}
		}
	}
	return f
}

// The variadic arguments may start with link and category types,
// and must end with a format string and any arguments.
// It returns the new Problem.
//nolint:interfacer
func (f *Node) errorf(n ast.Node, format string, args ...interface{}) {
	pos := f.fset.Position(n.Pos())
	f.errorfAt(pos, format, args...)
}

func (f *Node) errorfAt(pos token.Position, format string, args ...interface{}) {
	*f.issues = append(*f.issues, result.Issue{
		Pos:        pos,
		Text:       fmt.Sprintf(format, args...),
		FromLinter: Scopelint{}.Name(),
	})
}
