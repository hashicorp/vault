package main

import (
	"encoding/json"
	"go/ast"
	"go/parser" // provides methods for parsing source files and generating asts
	"go/token"  // provides types and methods for Go's lexer processor (tokenization)
	"io/ioutil"
	"strings"
)

func main() {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "../helper/builtinplugins/registry.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	cmap := ast.NewCommentMap(fset, file, file.Comments)

	v := Visitor{cmap, make([]Backend, 0), make([]Backend, 0), make([]Backend, 0)}

	ast.Walk(&v, file)

	output, _ := json.MarshalIndent(v, "", " ")
	_ = ioutil.WriteFile("plugins.json", output, 0644)
}

type Backend struct {
	Name       string `json:"name"`
	Deprecated bool   `json:"deprecated"`
}

type Visitor struct {
	// Comment map
	commentMap ast.CommentMap

	// Credential backends
	CredBackends []Backend `json:"CredentialBackends"`

	// Logical backends
	LogBackends []Backend `json:"LogicalBackends"`

	// DB backends
	DbBackends []Backend `json:"DbPlugins"`
}

func (v *Visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	if keyValueExpr, ok := n.(*ast.KeyValueExpr); ok {
		if ident, ok := keyValueExpr.Key.(*ast.Ident); ok {
			val, _ := keyValueExpr.Value.(*ast.CompositeLit)

			expressions := val.Elts
			backends := []Backend{}
			for i := 0; i < len(expressions); i++ {
				ex, _ := expressions[i].(*ast.KeyValueExpr)
				n := ex.Key.(*ast.BasicLit)
				d := false
				comments := v.commentMap.Filter(ex)
				for _, val := range comments {
					for _, c := range val {
						if strings.Contains(c.Text(), "Deprecated") {
							d = true
						}
					}
				}

				name := n.Value[1 : len(n.Value)-1]
				b := Backend{Name: name, Deprecated: d}
				backends = append(backends, b)
			}

			switch ident.Name {
			case "credentialBackends":
				v.CredBackends = append(v.CredBackends, backends...)
			case "databasePlugins":
				v.DbBackends = append(v.DbBackends, backends...)
			case "logicalBackends":
				v.LogBackends = append(v.LogBackends, backends...)
			}
		}
	}
	return v
}
