package checkers

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/go-lintpack/lintpack"
	"github.com/go-lintpack/lintpack/astwalk"
	"github.com/go-toolsmith/strparse"
)

func init() {
	var info lintpack.CheckerInfo
	info.Name = "commentedOutCode"
	info.Tags = []string{"diagnostic", "experimental"}
	info.Summary = "Detects commented-out code inside function bodies"
	info.Before = `
// fmt.Println("Debugging hard")
foo(1, 2)`
	info.After = `foo(1, 2)`

	collection.AddChecker(&info, func(ctx *lintpack.CheckerContext) lintpack.FileWalker {
		return astwalk.WalkerForLocalComment(&commentedOutCodeChecker{
			ctx:              ctx,
			notQuiteFuncCall: regexp.MustCompile(`\w+\s+\([^)]*\)\s*$`),
		})
	})
}

type commentedOutCodeChecker struct {
	astwalk.WalkHandler
	ctx *lintpack.CheckerContext

	notQuiteFuncCall *regexp.Regexp
}

func (c *commentedOutCodeChecker) VisitLocalComment(cg *ast.CommentGroup) {
	s := cg.Text() // Collect text once

	// We do multiple heuristics to avoid false positives.
	// Many things can be improved here.

	markers := []string{
		"TODO", // TODO comments with code are permitted.

		// "http://" is interpreted as a label with comment.
		// There are other protocols we might want to include.
		"http://",
		"https://",

		"e.g. ", // Clearly not a "selector expr" (mostly due to extra space)
	}
	for _, m := range markers {
		if strings.Contains(s, m) {
			return
		}
	}

	// Some very short comment that can be skipped.
	// Usually triggering on these results in false positive.
	// Unless there is a very popular call like print/println.
	cond := len(s) < len("quite too short") &&
		!strings.Contains(s, "print") &&
		!strings.Contains(s, "fmt.") &&
		!strings.Contains(s, "log.")
	if cond {
		return
	}

	// Almost looks like a commented-out function call,
	// but there is a whitespace between function name and
	// parameters list. Skip these to avoid false positives.
	if c.notQuiteFuncCall.MatchString(s) {
		return
	}

	stmt := strparse.Stmt(s)
	if stmt == strparse.BadStmt {
		return // Most likely not a code
	}

	if !c.isPermittedStmt(stmt) {
		c.warn(cg)
	}
}

func (c *commentedOutCodeChecker) isPermittedStmt(stmt ast.Stmt) bool {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		return c.isPermittedExpr(stmt.X)
	case *ast.LabeledStmt:
		return c.isPermittedStmt(stmt.Stmt)
	case *ast.DeclStmt:
		decl := stmt.Decl.(*ast.GenDecl)
		return decl.Tok == token.TYPE
	default:
		return false
	}
}

func (c *commentedOutCodeChecker) isPermittedExpr(x ast.Expr) bool {
	// Permit anything except expressions that can be used
	// with complete result discarding.
	switch x := x.(type) {
	case *ast.CallExpr:
		return false
	case *ast.UnaryExpr:
		// "<-" channel receive is not permitted.
		return x.Op != token.ARROW
	default:
		return true
	}
}

func (c *commentedOutCodeChecker) warn(cause ast.Node) {
	c.ctx.Warn(cause, "may want to remove commented-out code")
}
