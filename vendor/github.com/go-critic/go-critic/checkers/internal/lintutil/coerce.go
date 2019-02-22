package lintutil

import (
	"go/ast"
	"go/token"
)

var (
	nilIdent        = &ast.Ident{}
	nilSelectorExpr = &ast.SelectorExpr{}
	nilUnaryExpr    = &ast.UnaryExpr{}
	nilBinaryExpr   = &ast.BinaryExpr{}
	nilCallExpr     = &ast.CallExpr{}
	nilParenExpr    = &ast.ParenExpr{}
	nilAssignStmt   = &ast.AssignStmt{}
)

// IsNil reports whether x is nil.
// Unlike simple nil check, also detects nil AST sentinels.
func IsNil(x ast.Node) bool {
	switch x := x.(type) {
	case *ast.Ident:
		return x == nilIdent || x == nil
	case *ast.SelectorExpr:
		return x == nilSelectorExpr || x == nil
	case *ast.UnaryExpr:
		return x == nilUnaryExpr || x == nil
	case *ast.BinaryExpr:
		return x == nilBinaryExpr || x == nil
	case *ast.CallExpr:
		return x == nilCallExpr || x == nil
	case *ast.ParenExpr:
		return x == nilParenExpr || x == nil
	case *ast.AssignStmt:
		return x == nilAssignStmt || x == nil

	default:
		return x == nil
	}
}

// AsIdent coerces x into non-nil ident.
func AsIdent(x ast.Node) *ast.Ident {
	e, ok := x.(*ast.Ident)
	if !ok {
		return nilIdent
	}
	return e
}

// AsSelectorExpr coerces x into non-nil selector expr.
func AsSelectorExpr(x ast.Node) *ast.SelectorExpr {
	e, ok := x.(*ast.SelectorExpr)
	if !ok {
		return nilSelectorExpr
	}
	return e
}

// AsUnaryExpr coerces x into non-nil unary expr.
func AsUnaryExpr(x ast.Node) *ast.UnaryExpr {
	e, ok := x.(*ast.UnaryExpr)
	if !ok {
		return nilUnaryExpr
	}
	return e
}

// AsUnaryExprOp is like AsUnaryExpr, but also checks for op token.
func AsUnaryExprOp(x ast.Node, op token.Token) *ast.UnaryExpr {
	e, ok := x.(*ast.UnaryExpr)
	if !ok || e.Op != op {
		return nilUnaryExpr
	}
	return e
}

// AsBinaryExpr coerces x into non-nil binary expr.
func AsBinaryExpr(x ast.Node) *ast.BinaryExpr {
	e, ok := x.(*ast.BinaryExpr)
	if !ok {
		return nilBinaryExpr
	}
	return e
}

// AsBinaryExprOp is like AsBinaryExpr, but also checks for op token.
func AsBinaryExprOp(x ast.Node, op token.Token) *ast.BinaryExpr {
	e, ok := x.(*ast.BinaryExpr)
	if !ok || e.Op != op {
		return nilBinaryExpr
	}
	return e
}

// AsCallExpr coerces x into non-nil call expr.
func AsCallExpr(x ast.Node) *ast.CallExpr {
	e, ok := x.(*ast.CallExpr)
	if !ok {
		return nilCallExpr
	}
	return e
}

// AsParenExpr coerces x into non-nil paren expr.
func AsParenExpr(x ast.Node) *ast.ParenExpr {
	e, ok := x.(*ast.ParenExpr)
	if !ok {
		return nilParenExpr
	}
	return e
}

// AsAssignStmt coerces x into non-nil assign stmt.
func AsAssignStmt(x ast.Node) *ast.AssignStmt {
	stmt, ok := x.(*ast.AssignStmt)
	if !ok {
		return nilAssignStmt
	}
	return stmt
}
