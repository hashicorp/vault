package checkers

import (
	"go/ast"

	"github.com/go-lintpack/lintpack"
	"github.com/go-lintpack/lintpack/astwalk"
)

func init() {
	var info lintpack.CheckerInfo
	info.Name = "builtinShadow"
	info.Tags = []string{"style", "opinionated"}
	info.Summary = "Detects when predeclared identifiers shadowed in assignments"
	info.Before = `len := 10`
	info.After = `length := 10`

	collection.AddChecker(&info, func(ctx *lintpack.CheckerContext) lintpack.FileWalker {
		builtins := map[string]bool{
			// Types
			"bool":       true,
			"byte":       true,
			"complex64":  true,
			"complex128": true,
			"error":      true,
			"float32":    true,
			"float64":    true,
			"int":        true,
			"int8":       true,
			"int16":      true,
			"int32":      true,
			"int64":      true,
			"rune":       true,
			"string":     true,
			"uint":       true,
			"uint8":      true,
			"uint16":     true,
			"uint32":     true,
			"uint64":     true,
			"uintptr":    true,

			// Constants
			"true":  true,
			"false": true,
			"iota":  true,

			// Zero value
			"nil": true,

			// Functions
			"append":  true,
			"cap":     true,
			"close":   true,
			"complex": true,
			"copy":    true,
			"delete":  true,
			"imag":    true,
			"len":     true,
			"make":    true,
			"new":     true,
			"panic":   true,
			"print":   true,
			"println": true,
			"real":    true,
			"recover": true,
		}
		c := &builtinShadowChecker{ctx: ctx, builtins: builtins}
		return astwalk.WalkerForLocalDef(c, ctx.TypesInfo)
	})
}

type builtinShadowChecker struct {
	astwalk.WalkHandler
	ctx *lintpack.CheckerContext

	builtins map[string]bool
}

func (c *builtinShadowChecker) VisitLocalDef(name astwalk.Name, _ ast.Expr) {
	if _, isBuiltin := c.builtins[name.ID.String()]; isBuiltin {
		c.warn(name.ID)
	}
}

func (c *builtinShadowChecker) warn(ident *ast.Ident) {
	c.ctx.Warn(ident, "shadowing of predeclared identifier: %s", ident)
}
