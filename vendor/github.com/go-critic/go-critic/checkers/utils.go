package checkers

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/go-lintpack/lintpack"
)

// isStdlibPkg reports whether pkg is a package from the Go standard library.
func isStdlibPkg(pkg *types.Package) bool {
	return pkg != nil && pkg.Path() == pkg.Name()
}

// isUnitTestFunc reports whether FuncDecl declares testing function.
func isUnitTestFunc(ctx *lintpack.CheckerContext, fn *ast.FuncDecl) bool {
	if !strings.HasPrefix(fn.Name.Name, "Test") {
		return false
	}
	typ := ctx.TypesInfo.TypeOf(fn.Name)
	if sig, ok := typ.(*types.Signature); ok {
		return sig.Results().Len() == 0 &&
			sig.Params().Len() == 1 &&
			sig.Params().At(0).Type().String() == "*testing.T"
	}
	return false
}

// qualifiedName returns called expr fully-quallified name.
//
// It works for simple identifiers like f => "f" and identifiers
// from other package like pkg.f => "pkg.f".
//
// For all unexpected expressions returns empty string.
func qualifiedName(x ast.Expr) string {
	switch x := x.(type) {
	case *ast.SelectorExpr:
		pkg, ok := x.X.(*ast.Ident)
		if !ok {
			return ""
		}
		return pkg.Name + "." + x.Sel.Name
	case *ast.Ident:
		return x.Name
	default:
		return ""
	}
}

// identOf returns identifier for x that can be used to obtain associated types.Object.
// Returns nil for expressions that yield temporary results, like `f().field`.
func identOf(x ast.Node) *ast.Ident {
	switch x := x.(type) {
	case *ast.Ident:
		return x
	case *ast.SelectorExpr:
		return identOf(x.Sel)
	case *ast.TypeAssertExpr:
		// x.(type) - x may contain ident.
		return identOf(x.X)
	case *ast.IndexExpr:
		// x[i] - x may contain ident.
		return identOf(x.X)
	case *ast.StarExpr:
		// *x - x may contain ident.
		return identOf(x.X)
	case *ast.SliceExpr:
		// x[:] - x may contain ident.
		return identOf(x.X)

	default:
		// Note that this function is not comprehensive.
		return nil
	}
}
